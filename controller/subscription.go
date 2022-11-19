package controller

import (
	"fmt"
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"nfthelper/logger"
	"nfthelper/service"
	"nfthelper/status"
	"strconv"
	"strings"
)

type SubscriptionController struct {
	TgBotAPI       *tgBot.BotAPI
	commonService  *service.CommonService
	paymentService *service.PaymentService
}

func (c *SubscriptionController) Init(botAPI *tgBot.BotAPI) {
	c.TgBotAPI = botAPI
	c.commonService = new(service.CommonService)
	c.paymentService = new(service.PaymentService)
}

func (c *SubscriptionController) Subscription(message *tgBot.Message) {
	logger.Info("[command|subscription] handling, message is %s", message.Text)

	// todo 用户注册 message.From.ID
	// 如果有用户就获取用户的subscription信息
	// 发送subscription 信息
	text := "Your current subscription: ✅ <b>Basic</b>\n\n" +
		"....................................\n\n" +
		"🖼️️ <b>NFT</b> <i>3/5</i> ⚠️"
	msg := tgBot.NewMessage(message.Chat.ID, text)
	msg.ParseMode = tgBot.ModeHTML

	// 发送inline button
	inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
		tgBot.NewInlineKeyboardRow(
			tgBot.NewInlineKeyboardButtonData("🛎️ Choose subscription plan", "🛎️ Choose subscription plan"),
		),
	)
	msg.ReplyMarkup = inlineKeyboard

	if _, err := c.TgBotAPI.Send(msg); err != nil {
		logger.Error("[command|subscription] send message err, %v", err)
		return
	}

	// 设置indicator
	status.SetIndicator(message.From.ID, status.Subscription)
}

func (c *SubscriptionController) ListSubscription(callbackQuery *tgBot.CallbackQuery) {
	if status.GetIndicator(callbackQuery.From.ID) != status.Subscription {
		// 如果不是在添加NFT的时候，用户输入内容无效
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "Sorry, I don't understand. Please use /menu")
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[callback|list subscription] send message err, %v", err)
			return
		}
		status.SetIndicator(callbackQuery.From.ID, status.Start)
	} else {
		// 根据用户id获取当前plan
		// current
		text := "✅ <b>Basic</b><i>(Current Plan)</i>\n\n" +
			"☑️️ Up to 5 NFTs"
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, text)
		msg.ParseMode = tgBot.ModeHTML
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[command|list subscription:current] send message err, %v", err)
			return
		}

		subscriptions := c.commonService.ListSubscriptions()
		for _, subscription := range subscriptions {
			// advanced
			text = fmt.Sprintf("💎️ <b>%s monthly</b>\n\n"+
				"☑️️ Up to %d NFTs", subscription.Name, subscription.MaxNFT)
			msg = tgBot.NewMessage(callbackQuery.Message.Chat.ID, text)
			msg.ParseMode = tgBot.ModeHTML

			// 发送inline button
			inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
				tgBot.NewInlineKeyboardRow(
					tgBot.NewInlineKeyboardButtonData(fmt.Sprintf("Choose for $%.2f/month", subscription.Price), fmt.Sprintf("Choose subscription`%d", subscription.ID)),
				),
			)
			msg.ReplyMarkup = inlineKeyboard

			if _, err := c.TgBotAPI.Send(msg); err != nil {
				logger.Error("[command|list subscription: %s] send message err, %v", subscription.Name, err)
				return
			}
		}
		status.SetIndicator(callbackQuery.From.ID, status.ListSubscription)
	}
}

func (c *SubscriptionController) ChooseSubscription(callbackQuery *tgBot.CallbackQuery) {
	if status.GetIndicator(callbackQuery.From.ID) != status.ListSubscription {
		// 如果不是在添加NFT的时候，用户输入内容无效
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "Sorry, I don't understand. Please use /menu")
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[callback|choose subscription] send message err, %v", err)
			return
		}
		status.SetIndicator(callbackQuery.From.ID, status.Start)
	} else {
		// 获取subscription
		subscriptionIDStr := strings.Split(callbackQuery.Data, "`")[1]
		subscriptionID, _ := strconv.ParseInt(subscriptionIDStr, 10, 64)

		text := "Choose payment currency:"
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, text)

		// 获取currency
		currencies := c.commonService.ListCurrencies()
		var currencyKeyboards []tgBot.InlineKeyboardButton
		for _, currency := range currencies {
			kb := tgBot.NewInlineKeyboardButtonData(currency, fmt.Sprintf("Choose currency`%s`%d", currency, subscriptionID))
			currencyKeyboards = append(currencyKeyboards, kb)
		}
		// 发送inline button
		inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
			tgBot.NewInlineKeyboardRow(
				currencyKeyboards...,
			),
		)
		msg.ReplyMarkup = inlineKeyboard

		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[command|choose subscription] send message err, %v", err)
			return
		}
		status.SetIndicator(callbackQuery.From.ID, status.ChooseCurrency)
	}
}

func (c *SubscriptionController) ChooseCurrency(callbackQuery *tgBot.CallbackQuery) {
	if status.GetIndicator(callbackQuery.From.ID) != status.ChooseCurrency {
		// 如果不是在添加NFT的时候，用户输入内容无效
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "Sorry, I don't understand. Please use /menu")
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[callback|choose currency] send message err, %v", err)
			return
		}
		status.SetIndicator(callbackQuery.From.ID, status.Start)
	} else {
		// 获取subscription
		paramsStr := strings.Split(callbackQuery.Data, "`")
		currency := paramsStr[1]
		subscriptionID, _ := strconv.ParseInt(paramsStr[2], 10, 64)

		text := fmt.Sprintf("Choose network for %s coin:", currency)
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, text)

		// 获取currency
		networks := c.commonService.ListNetworks()
		var currencyKeyboards []tgBot.InlineKeyboardButton
		for _, network := range networks {
			kb := tgBot.NewInlineKeyboardButtonData(network, fmt.Sprintf("Choose network`%s`%s`%d", network, currency, subscriptionID))
			currencyKeyboards = append(currencyKeyboards, kb)
		}
		// 发送inline button
		inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
			tgBot.NewInlineKeyboardRow(
				currencyKeyboards...,
			),
		)
		msg.ReplyMarkup = inlineKeyboard

		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[command|choose currency] send message err, %v", err)
			return
		}
		status.SetIndicator(callbackQuery.From.ID, status.ChooseNetwork)
	}
}

func (c *SubscriptionController) ChooseNetwork(callbackQuery *tgBot.CallbackQuery) {
	if status.GetIndicator(callbackQuery.From.ID) != status.ChooseNetwork {
		// 如果不是在添加NFT的时候，用户输入内容无效
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "Sorry, I don't understand. Please use /menu")
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[callback|choose network] send message err, %v", err)
			return
		}
		status.SetIndicator(callbackQuery.From.ID, status.Start)
	} else {
		// 获取subscription
		paramsStr := strings.Split(callbackQuery.Data, "`")
		network := paramsStr[1]
		currency := paramsStr[2]
		subscriptionID, _ := strconv.ParseInt(paramsStr[3], 10, 64)
		subscription := c.commonService.GetSubscription(subscriptionID)

		paymentLink := c.paymentService.GeneratePaymentLink(network, currency, subscription.Price)
		text := fmt.Sprintf("💎️ <b>%s monthly</b>\n\n"+
			"<b>Please follow the link to proceed with your $%.2f payment. Once it is completed, you will receive a notification status.</b>\n\n"+
			"❗️Send <b>%s</b> using <b>%s</b> network ❗\n\n"+
			"Your personal payment link:\n"+
			"%s\n", subscription.Name, subscription.Price, currency, network, paymentLink)
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, text)
		msg.ParseMode = tgBot.ModeHTML

		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[callback|choose network] send message err, %v", err)
			return
		}
		// status.SetIndicator(callbackQuery.From.ID, status.Payment)
	}
}
