package processor

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"

	"msg-responder-ocr/internal/contract"
	"msg-responder-ocr/internal/logger"

	ocrv1 "msg-responder-ocr/internal/proto/ocr/v1"

	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	gstatus "google.golang.org/grpc/status"
)

const errorMsg = "⚠️ Произошла ошибка — мы уже решаем проблему ⚠️"

type Producer interface {
	Send(_ context.Context, topic string, data []byte) error
}

type messageResponder struct {
	logger     logger.Logger
	kafkaTopic string
	producer   Producer
	doc2text   string
	auth       AuthOptions
}

type Option struct {
	Doc2textURL string
	KafkaTopic  string
	Producer    Producer
	Logger      logger.Logger
	Auth        AuthOptions
}

func NewMessageResponder(opt Option) *messageResponder {
	return &messageResponder{
		logger:     opt.Logger,
		kafkaTopic: opt.KafkaTopic,
		producer:   opt.Producer,
		doc2text:   opt.Doc2textURL,
		auth:       opt.Auth,
	}
}

func (t *messageResponder) Handle(ctx context.Context, raw []byte) error {
	var requestMessage contract.OcrRequest
	if err := json.Unmarshal(raw, &requestMessage); err != nil {
		t.logger.Error("unmarshal request error: %v", err)
		return err
	}

	response := contract.NormalizedResponse{
		ChatID:         requestMessage.ChatID,
		Source:         requestMessage.Source,
		UserID:         requestMessage.UserID,
		Username:       requestMessage.Username,
		Timestamp:      requestMessage.Timestamp,
		OriginalUpdate: requestMessage.OriginalUpdate,
	}
	if len(requestMessage.Media) == 0 {
		t.logger.Error("empty media array in request")
		response.Text = errorMsg
	} else {
		text, err := processObject(ctx, t.doc2text, requestMessage.Media[0].S3URL, t.auth)
		if err != nil {
			code := gstatus.Code(err).String()
			t.logger.Error("process object error: code=%s err=%v", code, err)
			response.Text = errorMsg
		} else {
			response.Text = fmt.Sprintf("🖼️➜📝 Готово!\n\n🔷🔷🔷🔷🔷🔷\n%s\n🔷🔷🔷🔷🔷🔷", text)
		}
	}
	out, err := json.Marshal(response)
	if err != nil {
		t.logger.Error("marshal response error: %v", err)
		return err
	}

	t.logger.Debug("sending response: topic=%s bytes=%d", t.kafkaTopic, len(out))
	if err := t.producer.Send(ctx, t.kafkaTopic, out); err != nil {
		t.logger.Error("producer send error: %v", err)
		return err
	}
	t.logger.Info("response sent: topic=%s", t.kafkaTopic)
	return nil
}

type AuthOptions struct {
	AccessTokenURL string
	ClientID       string
	ClientSecret   string
}

func processObject(ctx context.Context, addr, objectKeyOrURL string, auth AuthOptions) (string, error) {
	var dialOpts []grpc.DialOption

	dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))

	conf := clientcredentials.Config{
		ClientID:     auth.ClientID,
		ClientSecret: auth.ClientSecret,
		TokenURL:     auth.AccessTokenURL,
	}
	ts := conf.TokenSource(ctx)
	perRPC := oauth.TokenSource{TokenSource: ts}
	dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(perRPC))

	conn, err := grpc.NewClient(addr, dialOpts...)
	if err != nil {
		return "", fmt.Errorf("grpc dial: %w", err)
	}
	defer conn.Close()

	client := ocrv1.NewOcrServiceClient(conn)

	resp, err := client.Process(ctx, &ocrv1.ParseRequest{Objectkey: objectKeyOrURL})
	if err != nil {
		return "", fmt.Errorf("rpc Process: %w", err)
	}
	return resp.GetText(), nil
}
