# ローカル開発環境で利用するホットリロード環境
FROM golang:1.19.2 as development

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

CMD ["air", "-c", ".air.toml"]
