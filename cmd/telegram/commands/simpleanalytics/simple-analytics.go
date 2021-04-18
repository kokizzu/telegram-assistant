package simpleanalytics

import (
	"github.com/adityathebe/telegram-assistant/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	DailySummary = "/simpleanalytics"
)

type HandlersFunc map[string]func(tgbotapi.Update)

type SimpleAnalyticsHandler struct {
	services *services.SimpleAnalytics
	cmdMap   map[string]func(tgbotapi.Update)
	bot      *tgbotapi.BotAPI
}

func NewHandler(bot *tgbotapi.BotAPI, service *services.SimpleAnalytics) *SimpleAnalyticsHandler {
	return &SimpleAnalyticsHandler{
		bot:      bot,
		services: service,
		cmdMap:   make(HandlersFunc),
	}
}

func (t *SimpleAnalyticsHandler) Handlers() HandlersFunc {
	t.cmdMap[DailySummary] = t.simpleAnalytics
	return t.cmdMap
}

func (t *SimpleAnalyticsHandler) simpleAnalytics(update tgbotapi.Update) {
	data, err := t.services.DailySummary()

	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
		t.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, data)
	msg.ParseMode = "HTML"
	t.bot.Send(msg)
}
