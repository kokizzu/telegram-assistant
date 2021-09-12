package nepsealpha

import (
	"fmt"

	"github.com/adityathebe/telegram-assistant/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	Portfolio = "/portfolio"
)

type HandlersFunc map[string]func(tgbotapi.Update)

type NepseAlphaHandler struct {
	services *services.NepseAlpha
	cmdMap   map[string]func(tgbotapi.Update)
	bot      *tgbotapi.BotAPI
}

func NewHandler(bot *tgbotapi.BotAPI, service *services.NepseAlpha) *NepseAlphaHandler {
	return &NepseAlphaHandler{
		bot:      bot,
		services: service,
		cmdMap:   make(HandlersFunc),
	}
}

func (t *NepseAlphaHandler) Handlers() HandlersFunc {
	t.cmdMap[Portfolio] = t.portfolio
	return t.cmdMap
}

func (t *NepseAlphaHandler) portfolio(update tgbotapi.Update) {
	data, err := t.services.PortfolioDailySummary()

	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
		t.bot.Send(msg)
		return
	}

	data = fmt.Sprintf("<pre>%s</pre>", data)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, data)
	msg.ParseMode = "HTML"
	t.bot.Send(msg)
}
