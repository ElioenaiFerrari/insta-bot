FROM golang:1.20-alpine

WORKDIR /app
COPY go.* . 
RUN go mod download
COPY . .
RUN apk add --update chromium
RUN go build -ldflags '-s -w' -o bin/main ./cmd/insta-bot
CMD ["./bin/main"]