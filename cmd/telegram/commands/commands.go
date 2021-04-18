package commands

import (
	"fmt"
	"log"

	"github.com/adityathebe/telegram-assistant/cmd/telegram/commands/hubstaff"
	"github.com/adityathebe/telegram-assistant/cmd/telegram/commands/simpleanalytics"
	"github.com/adityathebe/telegram-assistant/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Handler struct {
	bot      *tgbotapi.BotAPI
	services *services.Services
	cmdMap   HandlersFunc
}

type HandlersFunc map[string]func(tgbotapi.Update)

func NewHandler(bot *tgbotapi.BotAPI, services *services.Services) *Handler {
	return &Handler{
		bot:      bot,
		services: services,
		cmdMap:   make(HandlersFunc),
	}
}

func (t *Handler) registerHandler(hmap HandlersFunc) error {
	for k, v := range hmap {
		if t.cmdMap[k] != nil {
			return fmt.Errorf("duplicate command: %q", k)
		}
		t.cmdMap[k] = v
	}
	return nil
}

func (t *Handler) RegisterHandlers() error {
	hubstaff := hubstaff.NewHandler(t.bot, t.services.Hubstaff)
	if err := t.registerHandler(HandlersFunc(hubstaff.Handlers())); err != nil {
		return err
	}

	sa := simpleanalytics.NewHandler(t.bot, t.services.SimpleAnalytics)
	if err := t.registerHandler(HandlersFunc(sa.Handlers())); err != nil {
		return err
	}

	return nil
}

func (t *Handler) getHandler(cmd string) func(tgbotapi.Update) {
	for k, v := range t.cmdMap {
		if k == cmd {
			return v
		}
	}
	return nil
}

func (t *Handler) Handle(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		handlerFunc := t.getHandler(update.Message.Text)
		if handlerFunc != nil {
			go handlerFunc(update)
		}
	}
}

func (t *Handler) Commands() (cmds []string) {
	for k := range t.cmdMap {
		cmds = append(cmds, k)
	}
	return
}