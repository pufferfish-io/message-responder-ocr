# msg-responder-ocr

## Что делает

1. Подписывается на Kafka-топик запросов и десериализует `OcrRequest` (см. `internal/contract`).
2. Берёт первый медиа-файл, получает OAuth-токен через `doc2text` и вызывает gRPC `Parse`, передавая URL (S3).
3. Собирает `NormalizedResponse`, добавляет в начало источника (`source` из запроса) и публикует текст в Kafka-ответный топик.
4. Ошибки логируются и пользователю отсылается шаблонное сообщение `⚠️…`.

## Запуск

1. Задайте необходимые окружения (см. ниже).
2. Соберите и запустите локально:
   ```bash
   go run ./cmd/msg-responder-ocr
   ```
3. Или соберите Docker-образ и запустите его с переменными:
   ```bash
   docker build -t msg-responder-ocr .
   docker run --rm -e ... msg-responder-ocr
   ```

## Переменные окружения

Каждая переменная обязательна, кроме SASL, если Kafka не требует.

- `DOC3TEXT_ACCESS_TOKEN_URL` — URL провайдера OAuth для `doc2text`.
- `DOC3TEXT_CLIENT_ID` и `DOC3TEXT_CLIENT_SECRET` — данные клиента `doc2text`.
- `DOC3TEXT_G_RPC_URL` — адрес gRPC-сервиса распознавания текста.
- `KAFKA_BOOTSTRAP_SERVERS_VALUE` — список брокеров (`host:port[,host:port]`).
- `KAFKA_GROUP_ID_MESSAGE_RESPONDER_OCR` — идентификатор consumer group.
- `KAFKA_TOPIC_NAME_OCR_REQUEST` — входной топик с запросами.
- `KAFKA_TOPIC_NAME_TEMP_RESPONSE_PREPARER` — суффикс топика ответа; итоговый топик `source + response`.
- `KAFKA_CLIENT_ID_MESSAGE_RESPONDER_OCR` — идентификатор Kafka-клиента (продюсер и консьюмер).
- `KAFKA_SASL_USERNAME` и `KAFKA_SASL_PASSWORD` — по необходимости для SASL/PLAIN.

## Примечания

- `doc2text` должен возвращать текст по `ParseRequest` при получении объекта по URL.
- Логика сообщений описана в `internal/contract`, сам gRPC-клиент в `internal/processor`.


kubectl exec -n app <pod-name> -- /bin/sh -c 'apk add --no-cache curl >/dev/null 2>&1; grpcurl -help || true; nc -vz doc2text.app.svc.cluster.local 50052'. nc