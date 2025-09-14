package contract

import "time"

type OcrRequest struct {
	Source         string        `json:"source"`
	UserID         int64         `json:"user_id"`
	Username       *string       `json:"username,omitempty"`
	ChatID         int64         `json:"chat_id"`
	Timestamp      time.Time     `json:"timestamp"`
	Media          []MediaObject `json:"media,omitempty"`
	OriginalUpdate any           `json:"original_update"`
}

type MediaObject struct {
	Type           string  `json:"type"`
	OriginalFileID string  `json:"original_file_id"`
	Filename       *string `json:"filename,omitempty"`
	MimeType       *string `json:"mime_type,omitempty"`
	S3URL          string  `json:"s3_url"`
}

type NormalizedResponse struct {
	ChatID         int64     `json:"chat_id"`
	Text           string    `json:"text,omitempty"`
	Silent         bool      `json:"silent,omitempty"`
	Context        any       `json:"context,omitempty"`
	Source         string    `json:"source"`
	UserID         int64     `json:"user_id"`
	Username       *string   `json:"username,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
	OriginalUpdate any       `json:"original_update"`
}
