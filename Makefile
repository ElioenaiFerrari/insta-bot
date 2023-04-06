include .env

default: dev


.PHONY: dev 
dev:
	ENV=dev go run ./cmd/insta-bot

.PHONY: prd 
prd:
	go build -ldflags '-s -w' -o bin/ cmd/insta-bot/main.go
	ENV=prd ./bin/main

