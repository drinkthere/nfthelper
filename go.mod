module nfthelper

go 1.18

require (
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/joho/godotenv v1.4.0
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	go.uber.org/zap v1.23.0
)

require (
	github.com/jonboulle/clockwork v0.3.0 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
)

replace github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 => github.com/drinkthere/telegram-bot-api/v5 v5.0.0-20221117064003-aaf47ebfb6ee
