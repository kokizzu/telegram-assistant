.PHONY: build

build:
	go build -o bin github.com/adityathebe/telegram-assistant/cmd/telegram

run:
	go run github.com/adityathebe/telegram-assistant/cmd/telegram