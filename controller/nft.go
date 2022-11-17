package controller

import (
	"nfthelper/logger"
	"nfthelper/model"
	"nfthelper/status"

	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NFTController struct {
	TgBotAPI     *tgBot.BotAPI
	IndicatorMap *map[int64]model.Indicator
}

func (c *NFTController) Add(message *tgBot.Message) {
	indicator := status.GetIndicator(message.From.ID)
	logger.Info("in add the indicator is %s", indicator)
	// logger.Info("[command|add] handling, message is %+v", message)
	// todo 获取用户plan和已订阅数，如果超过了plan的最大订阅数，则返回edit或者update的inlineKeyword， message.From.ID

	// 判断是添加的第几步
	logger.Info("[command|add] handling, message is %+v", message.Text)
	// 发送onboard 信息
	text := "Enter NFT <b>token name</b> (e.g Azuki) or <b>contract address</b> (Ethereum network only):"
	msg := tgBot.NewMessage(message.Chat.ID, text)
	msg.ParseMode = tgBot.ModeHTML
	// 发送inline button
	inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
		tgBot.NewInlineKeyboardRow(
			tgBot.NewInlineKeyboardButtonData("❌ Cancel", "Cancel Add"),
		),
	)
	msg.ReplyMarkup = inlineKeyboard
	if _, err := c.TgBotAPI.Send(msg); err != nil {
		logger.Error("[command|add] send message err, %v", err)
		return
	}
	status.SetIndicator(message.From.ID, status.AddNFT)
}

func (c *NFTController) Search(message *tgBot.Message, indicatorMap *map[int64]model.Indicator) {
	// logger.Info("[command|add] handling, message is %+v", message)
}

func (c *NFTController) Cancel(callbackQuery *tgBot.CallbackQuery) {
	msg := tgBot.NewDeleteMessage(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID)
	if _, err := c.TgBotAPI.Send(msg); err != nil {
		logger.Error("[command|add] send message err, %v", err)
		return
	}
	status.SetIndicator(callbackQuery.From.ID, status.Start)
}
