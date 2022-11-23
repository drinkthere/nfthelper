package controller

import (
	"fmt"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"nfthelper/logger"
	"nfthelper/service"
	"nfthelper/status"
)

type CommonController struct {
	TgBotAPI            *tgBot.BotAPI
	subscriptionService *service.SubscriptionService
	paymentService      *service.PaymentService
}

func (c *CommonController) Init(botAPI *tgBot.BotAPI) {
	c.TgBotAPI = botAPI
	c.subscriptionService = new(service.SubscriptionService)
	c.paymentService = new(service.PaymentService)
}

func (c *CommonController) Start(message *tgBot.Message) {
	logger.Info("[command|start] handling, message is %+v", message)
	userID := uint(message.From.ID)
	subscription := c.subscriptionService.GetByUserID(userID)
	if subscription.ID == 0 {
		// 如果有用户就获取用户的subscription信息, 如果没有就注册用户，并且设置成basic plan
		c.subscriptionService.SetBasicSubscription(userID, true)
	}

	// 发送onboard 信息
	userName := message.From.FirstName
	text := fmt.Sprintf("🤖 <b>Hey hey, %s, welcome onboard!</b>\n\n"+
		"Here I could help you catch the <b>latest announcements</b> of NFT collections.\n\n", userName)
	msg := tgBot.NewMessage(message.Chat.ID, text)
	msg.ParseMode = tgBot.ModeHTML

	// 设置keyboard
	replyKeyboard := tgBot.NewReplyKeyboard(
		tgBot.NewKeyboardButtonRow(
			tgBot.NewKeyboardButton("➕ Add"),
			tgBot.NewKeyboardButton("🖼 NFT"),
		),
		tgBot.NewKeyboardButtonRow(
			tgBot.NewKeyboardButton("🛎️ Subscription"),
		),
	)
	msg.ReplyMarkup = replyKeyboard

	if _, err := c.TgBotAPI.Send(msg); err != nil {
		logger.Error("[command|start] send message err, %v", err)
		return
	}

	// 设置indicator
	status.SetIndicator(userID, status.Start)
}

func (c *CommonController) Menu(message *tgBot.Message) {
	logger.Info("[command|menu] handling, message is %+v", message)

	userID := uint(message.From.ID)
	subscription := c.subscriptionService.GetByUserID(userID)
	if subscription.ID == 0 {
		// 如果有用户就获取用户的subscription信息, 如果没有就注册用户，并且设置成basic plan
		c.subscriptionService.SetBasicSubscription(userID, true)
	}

	// 发送onboard 信息
	userName := message.From.FirstName
	text := fmt.Sprintf("🤖 <b>Hey hey, %s, welcome onboard!</b>\n\n"+
		"Here I could help you catch the <b>latest announcements</b> of NFT collections.\n\n", userName)
	msg := tgBot.NewMessage(message.Chat.ID, text)
	msg.ParseMode = tgBot.ModeHTML

	// 发送inline button
	inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
		tgBot.NewInlineKeyboardRow(
			tgBot.NewInlineKeyboardButtonData("➕ Add", "➕ Add"),
			tgBot.NewInlineKeyboardButtonData("🖼 NFT", "🖼 NFT"),
		),
		tgBot.NewInlineKeyboardRow(
			tgBot.NewInlineKeyboardButtonData("🛎️ Subscription", "🛎️ Subscription"),
		),
	)
	msg.ReplyMarkup = inlineKeyboard

	if _, err := c.TgBotAPI.Send(msg); err != nil {
		logger.Error("[command|menu] send message err, %v", err)
		return
	}

	// 设置indicator
	status.SetIndicator(userID, status.Start)
}
