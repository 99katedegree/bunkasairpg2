# syntax=docker/dockerfile:1

# --- build stage ---
# go:embed でモンスター画像（api/assets/images/）をバイナリに同梱するため、
# ビルドコンテキストには api/ 配下一式が必要。イメージ自体もその分大きくなる。
FROM golang:1.26-bookworm AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/server ./cmd/server

# --- runtime stage ---
FROM debian:bookworm-slim
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=build /out/server ./server
# golang-migrate が file://migrations という相対パスで読むため同梱する。
COPY migrations ./migrations

EXPOSE 8085
USER nobody
ENTRYPOINT ["./server"]
