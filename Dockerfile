# デプロイ用コンテナに含めるバイナリを作成するコンテナ
FROM golang:1.19.2 as deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags "-s" -o=./api ./cmd/api

# -----------------------------------------

# デプロイ用のコンテナ
FROM ubuntu:22.04 as deploy

RUN apt-get update && apt install -y ca-certificates

COPY --from=deploy-builder /app/api .

CMD ["./api"]

# -----------------------------------------

# ローカル開発環境で利用するホットリロード環境
FROM golang:1.19.2 as development

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

CMD ["air", "-c", ".air.toml"]
