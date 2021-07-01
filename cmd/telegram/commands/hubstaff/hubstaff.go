package hubstaff

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/adityathebe/telegram-assistant/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	HubstaffWeekly = "^/hsw ?[0-9]*$"
	HubstaffDaily  = "/hsd"
)

type HandlersFunc map[string]func(tgbotapi.Update)

type HubstaffHandler struct {
	services *services.Hubstaff
	cmdMap   map[string]func(tgbotapi.Update)
	bot      *tgbotapi.BotAPI
}

func NewHandler(bot *tgbotapi.BotAPI, service *services.Hubstaff) *HubstaffHandler {
	return &HubstaffHandler{
		bot:      bot,
		services: service,
		cmdMap:   make(HandlersFunc),
	}
}

func (t *HubstaffHandler) Handlers() HandlersFunc {
	t.cmdMap[HubstaffWeekly] = t.hubstaffWeekly
	t.cmdMap[HubstaffDaily] = t.hubstaffDaily
	return t.cmdMap
}

func (t *HubstaffHandler) hubstaffWeekly(update tgbotapi.Update) {
	re := regexp.MustCompile(HubstaffWeekly)
	matches := re.FindAllString(update.Message.Text, -1)
	var count int
	var err error
	if len(matches) == 1 {
		tokens := strings.Fields(update.Message.Text)
		if len(tokens) == 2 {
			count, err = strconv.Atoi(tokens[1])
			if err != nil {
				log.Println(err)
			}
		}
	}

	data, err := t.services.WeeklyStats(count)
	data = "<pre>" + data + "</pre>"

	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
		t.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, data)
	msg.ParseMode = "HTML"
	t.bot.Send(msg)
}

func (t *HubstaffHandler) hubstaffDaily(update tgbotapi.Update) {
	data, err := t.services.DailyStats()
	data = "<pre>" + data + "</pre>"

	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
		t.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, data)
	msg.ParseMode = "HTML"
	t.bot.Send(msg)
}
