package main

import (
	"log"

	"github.com/adityathebe/telegram-assistant/cmd/telegram/commands"
	"github.com/adityathebe/telegram-assistant/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	TG_OWNER_ID      int
	TG_API_KEY       string // Telegram API Key
	SA_API_KEY       string // Simple Analytics API Key
	SA_SITE_NAME     string // Simple Analytics site name
	HUBSTAFF_SESSION string // Hubstaff session cookie
	HUBSTAFF_ORG_ID  string // Hubstaff organization id
)

func initTelegram() {
	bot, err := tgbotapi.NewBotAPI(TG_API_KEY)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false
	log.Printf("Logged in as [@%s]\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	services := services.NewService(SA_API_KEY, SA_SITE_NAME, HUBSTAFF_SESSION, HUBSTAFF_ORG_ID)
	commandHandler := commands.NewHandler(bot, services, TG_OWNER_ID)
	if err := commandHandler.RegisterHandlers(); err != nil {
		log.Fatal(err)
	}
	commandHandler.Handle(updates)
}

func main() {
	readKeys()
	initTelegram()
}
