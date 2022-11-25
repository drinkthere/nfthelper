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

	userID := uint(message.From.ID)
	collections := c.collectionService.ListByUserID(userID)
	if len(collections) == 0 {
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
		text += "\nClick to see latest announcements:"
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
	status.SetIndicator(userID, status.Start)
}

func (c *NFTController) AddNFT(message *tgBot.Message) {
	logger.Info("[command|add nft] handling, message is %s", message.Text)

	// 获取用户的subscription方案
	userID := uint(message.From.ID)
	subscription := c.subscriptionService.GetByUserID(userID)
	logger.Info("===user's subscription is %+v", subscription)
	collections := c.collectionService.ListByUserID(userID)
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
	status.SetIndicator(userID, status.AddNFT)
}

func (c *NFTController) SearchNFT(message *tgBot.Message) {
	logger.Info("[text|search nft] handling, message is %s", message.Text)

	userID := uint(message.From.ID)
	if status.GetIndicator(userID) == status.AddNFT {
		var msg tgBot.MessageConfig
		if strings.HasPrefix(message.Text, "0x") && len(message.Text) == 42 {
			// 合约地址
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

			collections := c.collectionService.Search(message.Text)
			if len(collections) == 0 {
				// 没有找到对应的NFT Collection
				msg := tgBot.NewMessage(message.Chat.ID, "No NFTs were found.")
				if _, err := c.TgBotAPI.Send(msg); err != nil {
					logger.Error("[text|search nft] send message err, %v", err)
					return
				}
				return
			}

			text := "Choose your NFT:\n"
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
		status.SetIndicator(userID, status.ConfirmNFT)
	} else {
		// 如果不是在添加NFT的时候，用户输入内容无效
		msg := tgBot.NewMessage(message.Chat.ID, "Sorry, I don't understand. Please use /menu")
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[text|search nft] send message err, %v", err)
			return
		}
		status.SetIndicator(userID, status.Start)
	}
}

func (c *NFTController) Cancel(callbackQuery *tgBot.CallbackQuery) {
	logger.Info("[callback|cancel] handling, message is %s", callbackQuery.Data)

	userID := uint(callbackQuery.From.ID)
	msg := tgBot.NewDeleteMessage(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID)
	if _, err := c.TgBotAPI.Send(msg); err != nil {
		if err.Error() == "json: cannot unmarshal bool into Go value of type tgbotapi.Message" {
			status.SetIndicator(userID, status.Start)
		} else {
			logger.Error("[command|cancel] send message err, %v", err)
		}
		return
	}
	status.SetIndicator(userID, status.Start)
}

func (c *NFTController) ConfirmAddingNFT(callbackQuery *tgBot.CallbackQuery) {
	logger.Info("[callback|confirm adding nft] handling, message is %s", callbackQuery.Data)

	userID := uint(callbackQuery.From.ID)
	if status.GetIndicator(userID) != status.ConfirmNFT {
		// 如果不是在添加NFT的时候，用户输入内容无效
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "Sorry, I don't understand. Please use /menu")
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[callback|confirm adding nft] send message err, %v", err)
			return
		}
		status.SetIndicator(userID, status.Start)
		return
	}

	collectionIDStr := strings.Split(callbackQuery.Data, "`")[1]
	collectionID, _ := strconv.ParseUint(collectionIDStr, 10, 64)
	logger.Info("[callback|confirm adding nft] collection ID is %d", collectionID)

	// 判断用户是不是已经添加过这个NFT collection到watchlist
	if c.collectionService.HasAlreadyWatched(userID, uint(collectionID)) {
		// 如果不是在添加NFT的时候，用户输入内容无效
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "This token is already in your watchlist.")
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[callback|confirm adding nft] send message err, %v", err)
			return
		}
		status.SetIndicator(userID, status.Start)
		return
	}

	collection := c.collectionService.GetByID(uint(collectionID))
	c.collectionService.AddUserCollection(userID, collection)

	msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, collection.Name+" was added to your watchlist!")
	// 发送inline button
	inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
		tgBot.NewInlineKeyboardRow(
			tgBot.NewInlineKeyboardButtonData("🗑️ Delete "+collection.Name, "Delete NFT`"+fmt.Sprint(collection.ID)),
		),
		tgBot.NewInlineKeyboardRow(
			tgBot.NewInlineKeyboardButtonData("➕ Add more to watchlist", "➕ Add"),
		),
	)
	msg.ReplyMarkup = inlineKeyboard
	if _, err := c.TgBotAPI.Send(msg); err != nil {
		logger.Error("[callback|confirm adding nft] send message err, %v", err)
		return
	}
	status.SetIndicator(userID, status.Start)

	return
}

func (c *NFTController) EditNFTs(callbackQuery *tgBot.CallbackQuery) {
	logger.Info("[callback|edit nfts] handling, message is %s", callbackQuery.Data)

	userID := uint(callbackQuery.From.ID)
	collections := c.collectionService.ListByUserID(userID)
	if len(collections) == 0 {
		// 如果用户没有订阅NFT，就鼓励用户订阅
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "You haven't add any NFT to your watchlist.\n\n"+
			"It's a great time to add your first NFT")

		// 发送inline button
		inlineKeyboard := tgBot.NewInlineKeyboardMarkup(
			tgBot.NewInlineKeyboardRow(
				tgBot.NewInlineKeyboardButtonData("Add NFT", "➕ Add"),
			),
		)
		msg.ReplyMarkup = inlineKeyboard
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[callback|edit nfts] send message err, %v", err)
			return
		}
		return
	}

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

	userID := uint(callbackQuery.From.ID)
	collectionIDStr := strings.Split(callbackQuery.Data, "`")[1]
	collectionID, _ := strconv.ParseUint(collectionIDStr, 10, 64)
	collection := c.collectionService.GetByID(uint(collectionID))

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
	status.SetIndicator(userID, status.DeleteNFT)
}

func (c *NFTController) ConfirmDeleteNFT(callbackQuery *tgBot.CallbackQuery) {
	logger.Info("[callback|confirm deleting nft] handling, message is %s", callbackQuery.Data)

	userID := uint(callbackQuery.From.ID)
	if status.GetIndicator(userID) != status.DeleteNFT {
		// 如果不存在就报错
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "Sorry, I don't understand. Please use /menu")
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[text|confirm deleting nft] send message err, %v", err)
			return
		}
		status.SetIndicator(userID, status.Start)
		return
	}

	collectionIDStr := strings.Split(callbackQuery.Data, "`")[1]
	collectionID, _ := strconv.ParseUint(collectionIDStr, 10, 64)
	if !c.collectionService.HasAlreadyWatched(userID, uint(collectionID)) {
		// 如果不是在添加NFT的时候，用户输入内容无效
		msg := tgBot.NewMessage(callbackQuery.Message.Chat.ID, "Sorry, this collection doesn't exist or you haven't added it to your watchlist")
		if _, err := c.TgBotAPI.Send(msg); err != nil {
			logger.Error("[text|confirm deleting nft] send message err, %v", err)
			return
		}
		status.SetIndicator(userID, status.Start)
		return
	}

	collection := c.collectionService.GetByID(uint(collectionID))
	logger.Info("[callback|confirm deleting nft] collection ID is %d, collection is %+v", collectionID, collection)
	c.collectionService.DeleteUserCollection(userID, uint(collectionID))

	msg := tgBot.NewDeleteMessage(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID)
	if _, err := c.TgBotAPI.Send(msg); err != nil {
		if err.Error() == "json: cannot unmarshal bool into Go value of type tgbotapi.Message" {
			status.SetIndicator(userID, status.Start)
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
	// 如果不做设置，From会是bot的id
	callbackQuery.Message.From = callbackQuery.From
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
