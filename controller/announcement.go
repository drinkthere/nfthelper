package controller

import (
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"nfthelper/logger"
	"nfthelper/service"
	"nfthelper/status"
	"strconv"
	"strings"
)

type AnnouncementController struct {
	TgBotAPI            *tgBot.BotAPI
	announcementService *service.AnnouncementService
	collectionService   *service.CollectionService
}

func (c *AnnouncementController) Init(botAPI *tgBot.BotAPI) {
	c.TgBotAPI = botAPI
	c.announcementService = new(service.AnnouncementService)
	c.collectionService = new(service.CollectionService)
}

func (c *AnnouncementController) GetByCollectionIDAndUserID(callbackQuery *tgBot.CallbackQuery) {
	logger.Info("[|subscription] handling, message is %s", callbackQuery.Data)

	userID := uint(callbackQuery.From.ID)
	params := strings.Split(callbackQuery.Data, "`")
	collectionID, _ := strconv.ParseUint(params[1], 10, 64)
	collection := c.collectionService.GetByID(uint(collectionID))
	announcements := c.announcementService.GetByCollectionIDAndUserID(uint(collectionID), userID)

	for index, announcement := range announcements {
		text := ""
		if index == 0 {
			text += "The latest 3 announcements of " + collection.Name + ": \n\n"
		}
		text += announcement.URL
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, text)
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[command|subscription] send message err, %v", err)
			return
		}
	}
	// 设置indicator
	status.SetIndicator(userID, status.Start)
}
