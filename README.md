# msg-responder-ocr-ocr

## Recreate gRPC

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

export PATH="$PATH:$(go env GOPATH)/bin"

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/proto/ocr/v1/ocr.proto
```

## Command Guide

### Run with exported .env (one‑liner)

Exports all variables from `.env` into the current shell and runs the service.

```
export $(cat .env | xargs) && go run ./cmd/msg-responder-ocr
```

### Run with `source` (safer for complex values)

Loads `.env` preserving quotes and special characters, then runs the service.

```
set -a && source .env && set +a && go run ./cmd/msg-responder-ocr
```

### Fetch/clean module deps

Resolves dependencies and prunes unused ones.

```
go mod tidy
```

### Verbose build (diagnostics)

Builds the binary with verbose and command tracing. Removes old binary after build to keep the tree clean.

```
go build -v -x ./cmd/msg-responder-ocr && rm -f msg-responder-ocr
```

### Docker build (Buildx)

Builds the image with detailed progress logs and without cache.

```
docker buildx build --no-cache --progress=plain .
```

### Create and push tag

Cuts a release tag and pushes it to remote.

```
git tag v0.0.1
git push origin v0.0.1
```

### Manage tags

List all tags, delete a tag locally and remotely, verify deletion.

```
git tag -l
git tag -d vX.Y.Z
git push --delete origin vX.Y.Z
git ls-remote --tags origin | grep 'refs/tags/vX.Y.Z$'
```
