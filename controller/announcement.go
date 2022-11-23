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

func (c *AnnouncementController) GetByCollectionID(callbackQuery *tgBot.CallbackQuery) {
	logger.Info("[callback|announcement] handling, message is %s", callbackQuery.Data)

	userID := uint(callbackQuery.From.ID)
	params := strings.Split(callbackQuery.Data, "`")
	collectionID, _ := strconv.ParseUint(params[1], 10, 64)

	// 判断用户是否订阅了这个collection
	if !c.collectionService.HasAlreadyWatched(userID, uint(collectionID)) {
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "Sorry, you haven't added this collection to your watchlist.")
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[callback|announcement] send message err, %v", err)
		}
		return
	}

	collection := c.collectionService.GetByID(uint(collectionID))
	announcements := c.announcementService.ListByCollectionID(uint(collectionID))
	if len(announcements) == 0 {
		// 该项目没有发布过announcements
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "Sorry, there is no announcements of "+collection.Name)
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[callback|announcement] send message err, %v", err)
		}
		return
	}
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
