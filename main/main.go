package main

import (
	"fmt"
	"nfthelper/config"
	"nfthelper/logger"
	"nfthelper/router"
	"nfthelper/status"
	"os"

	tgBot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

var (
	cfg    config.Config
	botAPI *tgBot.BotAPI
	rt     router.Router
)

func Init() {
	// 加载配置文件
	cfg = *config.LoadConfig(os.Args[1])

	// 初始化日志
	logLevel := zapcore.InfoLevel
	if cfg.LogLevel == "DEBUG" {
		logLevel = zapcore.DebugLevel
	}
	logger.InitLogger(cfg.LogPath, logLevel)

	// 加载.env文件
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// 初始化telegram bot
	botAPI, err = tgBot.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		panic(err)
	}
	if cfg.IsTgDebug {
		botAPI.Debug = true
	}
	// 设置menuButton
	commands := tgBot.NewSetMyCommands(
		tgBot.BotCommand{
			Command:     "/start",
			Description: "start to subscribe NFT announcement",
		},
	)
	_, err = botAPI.Request(commands)
	if err != nil {
		panic(err)
	}

	// 初始化Status
	status.InitStatus()

	// 初始化Router
	rt.Init(botAPI)
}

func Start() {
	// 监听并处理电报消息
	u := tgBot.NewUpdate(0)
	u.Timeout = 60

	updates := botAPI.GetUpdatesChan(u)
	for update := range updates {
		rt.Route(update)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s config_file\n", os.Args[0])
		os.Exit(1)
	}

	Init()

	Start()
}
