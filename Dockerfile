FROM golang:1.23-alpine as builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o /app/main ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache \
    yt-dlp \
    ca-certificates \
    ffmpeg \
    bash

RUN mkdir /app
RUN mkdir /app/downloads

COPY --from=builder /app/main /app/main
COPY --from=builder /app/cookies.txt cookies.txt

RUN chmod 644 cookies.txt

#ENV API_KEY_YOUTUBE=your_api_key
#ENV CORS_ORIGIN=*

EXPOSE 8080

CMD ["/app/main"]
