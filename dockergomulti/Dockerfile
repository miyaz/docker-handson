# multi stage build のサンプル
# dockergo/と同じアプリですがイメージサイズが約1/100になります

FROM golang:latest as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /echo
COPY . .
RUN go build -o goecho main.go

# runtime image
FROM busybox
#FROM alpine # TLS/SSLを使う場合はalpineを使いRUNをコメントインすること
#RUN apk add --no-cache ca-certificates
COPY --from=builder /echo/goecho /goecho
ENTRYPOINT ["/goecho"]
