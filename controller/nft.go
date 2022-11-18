package controller

import (
	"fmt"
	"nfthelper/logger"
	"nfthelper/model"
	"nfthelper/service"
	"nfthelper/status"
	"strconv"
	"strings"

	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NFTController struct {
	TgBotAPI          *tgBot.BotAPI
	collectionService *service.CollectionService
}

func (c *NFTController) Init(botAPI *tgBot.BotAPI) {
	c.TgBotAPI = botAPI
	c.collectionService = new(service.CollectionService)
}

func (c *NFTController) Add(message *tgBot.Message) {
	logger.Info("[command|add] handling, message is %+v", message)
	// todo 获取用户plan和已订阅数，如果超过了plan的最大订阅数，则返回edit或者update的inlineKeyword， message.From.ID

	// 发送onboard 信息
	text := "Enter NFT <b>token name</b> (e.g Azuki) or <b>contract address</b> (Ethereum network only):"
	msg := tgBot.NewMessage(message.Chat.ID, text)
	msg.ParseMode = tgBot.ModeHTML
	// 发送inline button
	inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
		tgBot.NewInlineKeyboardRow(
			tgBot.NewInlineKeyboardButtonData("❌ Cancel", "Cancel adding NFT before inputting"),
		),
	)
	msg.ReplyMarkup = inlineKeyboard
	if _, err := c.TgBotAPI.Send(msg); err != nil {
		logger.Error("[command|add] send message err, %v", err)
		return
	}
	status.SetIndicator(message.From.ID, status.AddNFT)
}

func (c *NFTController) Search(message *tgBot.Message) {
	indicator := status.GetIndicator(message.From.ID)
	logger.Info("in search the indicator is %s, user id is %d", indicator, message.From.ID)
	// logger.Info("[text|search] handling, message is %+v", message)
	if status.GetIndicator(message.From.ID) == status.AddNFT {
		var msg tgBot.MessageConfig
		if strings.HasPrefix(message.Text, "0x") && len(message.Text) == 42 {
			// 合约地址
			// todo 按照合约地址搜索合约
			collection := model.Collection{
				ID:      1,
				Name:    "Homa Gang - Valentine (Homa Gang - Valentine)",
				Address: "0x6C869A43A9D362eF870d75daE56A01887578421d",
				Price:   6.1,
			}
			// 发送onboard 信息
			msg = tgBot.NewMessage(message.Chat.ID, "<b>"+collection.Name+"</b>")
			msg.ParseMode = tgBot.ModeHTML
			// 发送inline button
			inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
				tgBot.NewInlineKeyboardRow(
					tgBot.NewInlineKeyboardButtonData("✔️ Confirm", "Confirm adding NFT`"+fmt.Sprint(collection.ID)),
					tgBot.NewInlineKeyboardButtonData("❌ Cancel", "Cancel adding NFT after inputting"),
				),
			)
			msg.ReplyMarkup = inlineKeyboard
		} else {
			var options []tgBot.InlineKeyboardButton

			// todo 按照合约名称搜索合约
			collections := []model.Collection{
				{ID: 2, Name: "Azuki"},
				{ID: 3, Name: "AzukiApeSocialClub"},
				{ID: 4, Name: "OkayAzukis"},
			}
			text := "Choose you NFT:\n"
			for index, collection := range collections {
				id := index + 1
				text = text + fmt.Sprintf("<b>%d</b> %s\n", id, collection.Name)
				option := tgBot.NewInlineKeyboardButtonData(fmt.Sprint(id), "Confirm adding NFT`"+fmt.Sprint(collection.ID))
				options = append(options, option)
			}
			msg = tgBot.NewMessage(message.Chat.ID, text)
			msg.ParseMode = tgBot.ModeHTML
			// 发送inline button
			inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
				tgBot.NewInlineKeyboardRow(options...),
			)
			msg.ReplyMarkup = inlineKeyboard
		}
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[command|search] send message err, %v", err)
			return
		}
		status.SetIndicator(message.From.ID, status.ConfirmNFT)
	} else {
		// 如果不是在添加NFT的时候，用户输入内容无效
		msg := tgBot.NewMessage(message.Chat.ID, "Sorry, I don't understand. Please use /menu")
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[text|search] send message err, %v", err)
			return
		}
		status.SetIndicator(message.From.ID, status.Start)
	}
}

func (c *NFTController) Cancel(callbackQuery *tgBot.CallbackQuery) {
	// todo 增加status 判断
	logger.Info("[callback|cancel] handling, message is %s", callbackQuery.Data)
	msg := tgBot.NewDeleteMessage(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID)
	if _, err := c.TgBotAPI.Send(msg); err != nil {
		if err.Error() == "json: cannot unmarshal bool into Go value of type tgbotapi.Message" {
			status.SetIndicator(callbackQuery.From.ID, status.Start)
		} else {
			logger.Error("[command|cancel] send message err, %v", err)
		}
		return
	}
	status.SetIndicator(callbackQuery.From.ID, status.Start)
}

func (c *NFTController) Confirm(callbackQuery *tgBot.CallbackQuery) {
	logger.Info("[callback|confirm] handling, message is %s", callbackQuery.Data)
	//indicator := status.GetIndicator(callbackQuery.From.ID)
	//logger.Info("in add the indicator is %s, user id is %d", indicator, callbackQuery.From.ID)
	if status.GetIndicator(callbackQuery.From.ID) != status.ConfirmNFT {
		// 如果不是在添加NFT的时候，用户输入内容无效
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "Sorry, I don't understand. Please use /menu")
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[text|search] send message err, %v", err)
			return
		}
		status.SetIndicator(callbackQuery.From.ID, status.Start)
		return
	}

	collectionIDStr := strings.Split(callbackQuery.Data, "`")[1]
	collectionID, _ := strconv.ParseInt(collectionIDStr, 10, 64)
	logger.Info("[callback|confirm] collection ID is %d", collectionID)

	collection := c.collectionService.GetCollectionByID(collectionID)
	// todo 数据库存储

	msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, collection.Name+" was added to Main watchlist!")
	// 发送inline button
	inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
		tgBot.NewInlineKeyboardRow(
			tgBot.NewInlineKeyboardButtonData("🔧 Edit "+collection.Name, "Edit "+collection.Name+"`"+fmt.Sprint(collection.ID)),
		),
		tgBot.NewInlineKeyboardRow(
			tgBot.NewInlineKeyboardButtonData("➕ Add more to main", "➕ Add"),
		),
	)
	msg.ReplyMarkup = inlineKeyboard
	if _, err := c.TgBotAPI.Send(msg); err != nil {
		logger.Error("[callback|confirm] send message err, %v", err)
		return
	}
	status.SetIndicator(callbackQuery.From.ID, status.Start)

	return
}
