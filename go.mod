module nfthelper

go 1.18

require (
	github.com/NdoleStudio/coinpayments-go v0.0.3
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/joho/godotenv v1.4.0
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	go.uber.org/zap v1.23.0
	gorm.io/driver/mysql v1.4.4
	gorm.io/gorm v1.24.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jonboulle/clockwork v0.3.0 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
)

replace github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1 => github.com/drinkthere/telegram-bot-api/v5 v5.0.0-20221117064003-aaf47ebfb6ee

replace github.com/NdoleStudio/coinpayments-go v0.0.3 => ./coinpayments-go
