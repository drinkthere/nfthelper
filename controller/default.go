package controller

import (
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"nfthelper/logger"
	"nfthelper/model"
)

type DefaultController struct {
	tgBotAPI *tgBot.BotAPI
}

func (c *DefaultController) Init(botAPI *tgBot.BotAPI) {
	c.tgBotAPI = botAPI
}

func (c *DefaultController) Handle(message *tgBot.Message, indicatorMap *map[int64]model.Indicator) {
	logger.Info("[command|default] handling, message is %+v", message)
}
