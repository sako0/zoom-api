# go
FROM golang:1.18.0-alpine
# apkのアップデートとgitのインストール
RUN apk add --update && apk add git && apk add build-base
# appディレクトリの作成
RUN mkdir /src

# ワーキングディレクトリの設定
WORKDIR /src
# ホストのファイルをコンテナの作業ディレクトリに移行
ADD ../src /src

# AIRのインストール
RUN go install github.com/cosmtrek/air@latest

# sql-migrateのインストール
RUN go install github.com/rubenv/sql-migrate/...

# ECHOのセットアップ
RUN go get -u github.com/labstack/echo

# GORMのセットアップ
RUN go get github.com/jinzhu/gorm

RUN go mod tidy

CMD ["air", "-c", ".air.toml"]
