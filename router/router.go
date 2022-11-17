package router

import (
	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"nfthelper/controller"
	"nfthelper/logger"
)

type Router struct {
	startController *controller.StartController
	nftController   *controller.NFTController
}

func (r *Router) Init(botAPI *tgBot.BotAPI) {
	r.startController = &controller.StartController{
		TgBotAPI: botAPI,
	}

	r.nftController = &controller.NFTController{
		TgBotAPI: botAPI,
	}
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
	switch callbackQuery.Data {
	case "Cancel Add":
		r.nftController.Cancel(callbackQuery)
	}
}

func (r *Router) RouteCommand(message *tgBot.Message) {
	switch message.Command() {
	case "start":
		r.startController.Handle(message)
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
	switch message.Text {
	case "➕ Add":
		r.nftController.Add(message)
	}
}
