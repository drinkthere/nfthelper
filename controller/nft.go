package controller

import (
	"fmt"
	"nfthelper/common"
	"nfthelper/logger"
	"nfthelper/service"
	"nfthelper/status"
	"strconv"
	"strings"

	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NFTController struct {
	TgBotAPI            *tgBot.BotAPI
	collectionService   *service.CollectionService
	subscriptionService *service.SubscriptionService
}

func (c *NFTController) Init(botAPI *tgBot.BotAPI) {
	c.TgBotAPI = botAPI
	c.collectionService = new(service.CollectionService)
}

func (c *NFTController) ListNFT(message *tgBot.Message) {
	logger.Info("[command|list nft] handling, message is %s", message.Text)
	collections := c.collectionService.ListByUserID(message.From.ID)
	if collections == nil {
		// 如果用户没有订阅NFT，就鼓励用户订阅
		msg := tgBot.NewMessage(message.Chat.ID, "It's a great time to add your first NFT")

		// 发送inline button
		inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
			tgBot.NewInlineKeyboardRow(
				tgBot.NewInlineKeyboardButtonData("Add NFT", "➕ Add"),
			),
		)
		msg.ReplyMarkup = inlineKeyboard
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[command|list nft] send message err, %v", err)
			return
		}
	} else {
		// 用户已订阅就展示用户定义的NFT
		var options []tgBot.InlineKeyboardButton

		text := "Your watchlist:\n"
		for index, collection := range collections {
			id := index + 1
			text = text + fmt.Sprintf("<b>%d</b> %s\n", id, collection.Name)
			option := tgBot.NewInlineKeyboardButtonData(fmt.Sprint(id), "Get NFT announcement`"+fmt.Sprint(collection.ID))
			options = append(options, option)
		}
		// 分成多行
		optionsRow := splitToMultiRows(5, options)
		text += "\nClick to see latest announcement:"
		msg := tgBot.NewMessage(message.Chat.ID, text)
		msg.ParseMode = tgBot.ModeHTML
		// 发送inline button
		inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
			optionsRow...,
		)
		msg.ReplyMarkup = inlineKeyboard
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[command|list nft] send message err, %v", err)
			return
		}
	}
	status.SetIndicator(message.From.ID, status.Start)
}

func (c *NFTController) AddNFT(message *tgBot.Message) {
	logger.Info("[command|add nft] handling, message is %s", message.Text)

	// 获取用户的subscription方案
	subscription := c.subscriptionService.GetByUserID(message.From.ID)

	collections := c.collectionService.ListByUserID(message.From.ID)
	if len(collections) >= subscription.MaxNFT {
		// 提示用户 删除 或 更新 subscription 方案
		text := fmt.Sprintf("⚠️ ️You have reached <b>%d NFTs</b> limits:\n\n"+
			"🖼️️ <b>NFT</b> %d/%d ⚠️\n\n"+
			"<b>Edit NFTs</b> or <b>Upgrade</b> your Subscription Plan", subscription.MaxNFT, len(collections), subscription.MaxNFT)
		msg := tgBot.NewMessage(message.Chat.ID, text)
		msg.ParseMode = tgBot.ModeHTML
		// 发送inline button
		inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
			tgBot.NewInlineKeyboardRow(
				tgBot.NewInlineKeyboardButtonData("Edit NFTs", "Edit NFTs"),
			),
			tgBot.NewInlineKeyboardRow(
				tgBot.NewInlineKeyboardButtonData("🛎️ Upgrade subscription plan", "🛎️ Upgrade subscription plan"),
			),
		)
		msg.ReplyMarkup = inlineKeyboard
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[command|add nft] send message err, %v", err)
			return
		}
		return
	}

	// 发送add信息
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
		logger.Error("[command|add nft] send message err, %v", err)
		return
	}
	status.SetIndicator(message.From.ID, status.AddNFT)
}

func (c *NFTController) SearchNFT(message *tgBot.Message) {
	logger.Info("[text|search nft] handling, message is %s", message.Text)
	//indicator := status.GetIndicator(message.From.ID)
	//logger.Info("in search the indicator is %s, user id is %d", indicator, message.From.ID)

	if status.GetIndicator(message.From.ID) == status.AddNFT {
		var msg tgBot.MessageConfig
		if strings.HasPrefix(message.Text, "0x") && len(message.Text) == 42 {
			// 合约地址
			// todo 按照合约地址搜索合约
			collection := c.collectionService.GetByAddr(message.Text)
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
			collections := c.collectionService.Search(message.Text)
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
			logger.Error("[text|search nft] send message err, %v", err)
			return
		}
		status.SetIndicator(message.From.ID, status.ConfirmNFT)
	} else {
		// 如果不是在添加NFT的时候，用户输入内容无效
		msg := tgBot.NewMessage(message.Chat.ID, "Sorry, I don't understand. Please use /menu")
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[text|search nft] send message err, %v", err)
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

func (c *NFTController) ConfirmAddingNFT(callbackQuery *tgBot.CallbackQuery) {
	logger.Info("[callback|confirm adding nft] handling, message is %s", callbackQuery.Data)
	//indicator := status.GetIndicator(callbackQuery.From.ID)
	//logger.Info("in add the indicator is %s, user id is %d", indicator, callbackQuery.From.ID)
	if status.GetIndicator(callbackQuery.From.ID) != status.ConfirmNFT {
		// 如果不是在添加NFT的时候，用户输入内容无效
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "Sorry, I don't understand. Please use /menu")
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[callback|confirm adding nft] send message err, %v", err)
			return
		}
		status.SetIndicator(callbackQuery.From.ID, status.Start)
		return
	}

	collectionIDStr := strings.Split(callbackQuery.Data, "`")[1]
	collectionID, _ := strconv.ParseInt(collectionIDStr, 10, 64)
	logger.Info("[callback|confirm adding nft] collection ID is %d", collectionID)

	collection := c.collectionService.GetByID(collectionID)
	// todo 数据库存储

	msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, collection.Name+" was added to Main watchlist!")
	// 发送inline button
	inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
		tgBot.NewInlineKeyboardRow(
			tgBot.NewInlineKeyboardButtonData("🗑️ Delete "+collection.Name, "Delete NFT`"+fmt.Sprint(collection.ID)),
		),
		tgBot.NewInlineKeyboardRow(
			tgBot.NewInlineKeyboardButtonData("➕ Add more to main", "➕ Add"),
		),
	)
	msg.ReplyMarkup = inlineKeyboard
	if _, err := c.TgBotAPI.Send(msg); err != nil {
		logger.Error("[callback|confirm adding nft] send message err, %v", err)
		return
	}
	status.SetIndicator(callbackQuery.From.ID, status.Start)

	return
}

func (c *NFTController) EditNFTs(callbackQuery *tgBot.CallbackQuery) {
	logger.Info("[callback|edit nfts] handling, message is %s", callbackQuery.Data)
	collections := c.collectionService.ListByUserID(callbackQuery.From.ID)
	// 发送list 信息
	msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "Select NFT to <b>delete</b>:\n\n")
	msg.ParseMode = tgBot.ModeHTML

	// 发送inline button
	var options []tgBot.InlineKeyboardButton
	for _, collection := range collections {
		option := tgBot.NewInlineKeyboardButtonData(fmt.Sprintf("Delete %s", collection.Name), "Delete NFT`"+fmt.Sprint(collection.ID))
		options = append(options, option)
	}
	// 分成多行
	optionsRow := splitToMultiRows(1, options)
	inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
		optionsRow...,
	)
	msg.ReplyMarkup = inlineKeyboard
	if _, err := c.TgBotAPI.Send(msg); err != nil {
		logger.Error("[callback|edit nfts] send message err, %v", err)
		return
	}
}

func (c *NFTController) DeleteNFT(callbackQuery *tgBot.CallbackQuery) {
	logger.Info("[callback|delete nft] handling, message is %s", callbackQuery.Data)

	collectionIDStr := strings.Split(callbackQuery.Data, "`")[1]
	collectionID, _ := strconv.ParseInt(collectionIDStr, 10, 64)
	collection := c.collectionService.GetByID(collectionID)
	logger.Info("[callback|delete nft] collection ID is %d, collection is %+v", collectionID, collection)

	msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "Do you want to remove token <b>"+collection.Name+"</b>?\n\n")
	msg.ParseMode = tgBot.ModeHTML
	// 发送inline button
	inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
		tgBot.NewInlineKeyboardRow(
			tgBot.NewInlineKeyboardButtonData("🗑️ Delete", "Confirm deleting NFT`"+fmt.Sprint(collection.ID)),
			tgBot.NewInlineKeyboardButtonData("❌ Cancel", "Cancel deleting NFT"),
		),
	)
	msg.ReplyMarkup = inlineKeyboard
	if _, err := c.TgBotAPI.Send(msg); err != nil {
		logger.Error("[callback|delete nft] send message err, %v", err)
		return
	}
	status.SetIndicator(callbackQuery.From.ID, status.DeleteNFT)
}

func (c *NFTController) ConfirmDeleteNFT(callbackQuery *tgBot.CallbackQuery) {
	logger.Info("[callback|confirm deleting nft] handling, message is %s", callbackQuery.Data)
	//indicator := status.GetIndicator(callbackQuery.From.ID)
	//logger.Info("in add the indicator is %s, user id is %d", indicator, callbackQuery.From.ID)
	if status.GetIndicator(callbackQuery.From.ID) != status.DeleteNFT {
		// 如果不是在添加NFT的时候，用户输入内容无效
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "Sorry, I don't understand. Please use /menu")
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[text|confirm deleting nft] send message err, %v", err)
			return
		}
		status.SetIndicator(callbackQuery.From.ID, status.Start)
		return
	}

	collectionIDStr := strings.Split(callbackQuery.Data, "`")[1]
	collectionID, _ := strconv.ParseInt(collectionIDStr, 10, 64)
	collection := c.collectionService.GetByID(collectionID)
	logger.Info("[callback|confirm deleting nft] collection ID is %d, collection is %+v", collectionID, collection)

	// todo 数据库删除订阅

	msg := tgBot.NewDeleteMessage(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID)
	if _, err := c.TgBotAPI.Send(msg); err != nil {
		if err.Error() == "json: cannot unmarshal bool into Go value of type tgbotapi.Message" {
			status.SetIndicator(callbackQuery.From.ID, status.Start)
		} else {
			logger.Error("[callback|confirm deleting nft, delete msg] send message err, %v", err)
		}
	}

	textMsg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "<b>"+collection.Name+"</b> has been successfully removed from your watchlist")
	textMsg.ParseMode = tgBot.ModeHTML
	if _, err := c.TgBotAPI.Send(textMsg); err != nil {
		logger.Error("[callback|confirm deleting nft] send message err, %v", err)
		return
	}

	// 发送更新之后的订阅信息给用户
	c.ListNFT(callbackQuery.Message)
	return
}

func splitToMultiRows(numPerRow int, options []tgBot.InlineKeyboardButton) (optionsRow [][]tgBot.InlineKeyboardButton) {
	for i := 0; i <= len(options)/numPerRow; i += 1 {
		start := i * numPerRow
		end := (i + 1) * numPerRow
		end = common.Min(end, len(options))
		opts := options[start:end]
		optionsRow = append(optionsRow, opts)
	}
	return
}
