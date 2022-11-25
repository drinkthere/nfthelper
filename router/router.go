package router

import (
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"nfthelper/controller"
	"nfthelper/logger"
	"strings"
)

type Router struct {
	commonController       *controller.CommonController
	nftController          *controller.NFTController
	subscriptionController *controller.SubscriptionController
	announcementController *controller.AnnouncementController
}

func (r *Router) Init(botAPI *tgBot.BotAPI) {
	r.commonController = new(controller.CommonController)
	r.commonController.Init(botAPI)

	r.nftController = new(controller.NFTController)
	r.nftController.Init(botAPI)

	r.subscriptionController = new(controller.SubscriptionController)
	r.subscriptionController.Init(botAPI)

	r.announcementController = new(controller.AnnouncementController)
	r.announcementController.Init(botAPI)
}

func (r *Router) Route(update tgBot.Update) {
	if update.CallbackQuery != nil {
		logger.Info("[callback|%s] arrive", update.CallbackQuery.Data)
		// 处理keyboard 回调
		r.RouteCallback(update.CallbackQuery)
	} else if update.Message != nil {
		if update.Message.IsCommand() {
			logger.Info("[command|%s] arrive", update.Message.Command())
			// 处理命令
			r.RouteCommand(update.Message)
		} else {
			// 处理文本
			logger.Info("[text|%s] arrive", update.Message.Text)
			r.RouteText(update.Message)
		}
	}
}

func (r *Router) RouteCallback(callbackQuery *tgBot.CallbackQuery) {
	dataSlice := strings.Split(callbackQuery.Data, "`")
	data := dataSlice[0]
	logger.Info("[callback] data is %s", data)
	switch data {
	case "🖼 NFT":
		callbackQuery.Message.From.ID = callbackQuery.From.ID
		r.nftController.ListNFT(callbackQuery.Message)
	case "Cancel adding NFT before inputting", "Cancel adding NFT after inputting",
		"Cancel deleting NFT":
		r.nftController.Cancel(callbackQuery)
	case "Confirm adding NFT":
		r.nftController.ConfirmAddingNFT(callbackQuery)
	case "➕ Add":
		callbackQuery.Message.From.ID = callbackQuery.From.ID
		r.nftController.AddNFT(callbackQuery.Message)
	case "Edit NFTs":
		callbackQuery.Message.From.ID = callbackQuery.From.ID
		r.nftController.EditNFTs(callbackQuery)
	case "Delete NFT":
		callbackQuery.Message.From.ID = callbackQuery.From.ID
		r.nftController.DeleteNFT(callbackQuery)
	case "Confirm deleting NFT":
		r.nftController.ConfirmDeleteNFT(callbackQuery)
	case "🛎️ Subscription":
		callbackQuery.Message.From.ID = callbackQuery.From.ID
		r.subscriptionController.Subscription(callbackQuery.Message)
	case "🛎️ Choose subscription plan", "🛎️ Upgrade subscription plan":
		r.subscriptionController.ListSubscription(callbackQuery)
	case "Choose subscription":
		r.subscriptionController.ChooseSubscription(callbackQuery)
	case "Choose currency":
		r.subscriptionController.ChooseCurrency(callbackQuery)
	case "Choose network":
		r.subscriptionController.ChooseNetwork(callbackQuery)
	case "Get NFT announcement":
		r.announcementController.GetByCollectionID(callbackQuery)
	}

}

func (r *Router) RouteCommand(message *tgBot.Message) {
	switch message.Command() {
	case "start":
		r.commonController.Start(message)
	case "menu":
		r.commonController.Menu(message)
	}
}

func (r *Router) RouteText(message *tgBot.Message) {
	dataSlice := strings.Split(message.Text, "`")
	data := dataSlice[0]
	switch data {
	case "➕ Add":
		r.nftController.AddNFT(message)
	case "🛎️ Subscription":
		r.subscriptionController.Subscription(message)
	case "🖼 NFT":
		r.nftController.ListNFT(message)
	default:
		r.nftController.SearchNFT(message)
	}
}
