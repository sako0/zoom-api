# go
FROM golang:1.18.0-alpine
# apkのアップデートとgitのインストール
RUN apk add --update && apk add git && apk add build-base

# ワーキングディレクトリの設定
WORKDIR /go/src/zoom-api
# ホストのファイルをコンテナの作業ディレクトリに移行
ADD ../ /go/src/zoom-api/

# AIRのインストール
RUN go install github.com/cosmtrek/air@latest

# sql-migrateのインストール
RUN go install github.com/rubenv/sql-migrate/...

# gRPC周りのインストール
RUN go install github.com/bufbuild/buf/cmd/buf@latest
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest

# Connectのセットアップ
RUN go get github.com/bufbuild/connect-go

# net/http2のセットアップ
RUN go get golang.org/x/net/http2

# GORMのセットアップ
RUN go get github.com/jinzhu/gorm

RUN go mod tidy

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
