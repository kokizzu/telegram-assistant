.PHONY: build

build:
	go build -o bin \
		github.com/adityathebe/telegram-assistant/cmd/telegram \
		github.com/adityathebe/telegram-assistant/cmd/cli

run:
	go run github.com/adityathebe/telegram-assistant/cmd/telegram

install:
	go build -i -o ~/go/bin/pa github.com/adityathebe/telegram-assistant/cmd/cli
