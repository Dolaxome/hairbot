.PHONY:
.SILENT:
build:
	go build -o ./.bin/bot cmd/bot/main.go
run: build
	./.bin/bot

build-image:
	docker build -t hair-bot:v0.1 .

start-container:
	docker run --name telegram-bot -p 80:80 --env-file .env hair-bot:v0.1