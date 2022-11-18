package router

import (
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"nfthelper/controller"
	"nfthelper/logger"
	"strings"
)

type Router struct {
	commonController *controller.CommonController
	nftController    *controller.NFTController
}

func (r *Router) Init(botAPI *tgBot.BotAPI) {
	r.commonController = new(controller.CommonController)
	r.commonController.Init(botAPI)

	r.nftController = new(controller.NFTController)
	r.nftController.Init(botAPI)
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
	case "Cancel adding NFT before inputting", "Cancel adding NFT after inputting":
		r.nftController.Cancel(callbackQuery)
	case "Confirm adding NFT":
		r.nftController.Confirm(callbackQuery)
	case "➕ Add":
		callbackQuery.Message.From.ID = callbackQuery.From.ID
		r.nftController.Add(callbackQuery.Message)
	}

}

func (r *Router) RouteCommand(message *tgBot.Message) {
	switch message.Command() {
	case "start":
		r.commonController.Start(message)
	case "menu":
		r.commonController.Menu(message)
	}
	//
	//if commandController, ok := controller.CommandControllersMap[message.Command()]; ok {
	//	logger.Info("[command] controller is %+v", commandController)
	//	commandController.Handle(message, indicatorMap)
	//} else {
	//	defaultController, _ := controller.CommandControllersMap["default"]
	//	defaultController.Handle(message, indicatorMap)
	//}
}

func (r *Router) RouteText(message *tgBot.Message) {
	dataSlice := strings.Split(message.Text, "`")
	data := dataSlice[0]
	switch data {
	case "➕ Add":
		r.nftController.Add(message)
	default:
		r.nftController.Search(message)
	}
}
