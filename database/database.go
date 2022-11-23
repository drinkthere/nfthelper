package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"nfthelper/model/dbmodel"
)

var (
	DB *gorm.DB
)

type Config struct {
	Host   string
	Port   int64
	User   string
	Pass   string
	DBName string
}

func Init(cfg *Config) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DBName)
	//连接MYSQL
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	// 迁移 schema
	DB.AutoMigrate(&dbmodel.Subscription{}, &dbmodel.Collection{}, &dbmodel.UserSubscription{}, &dbmodel.UserCollection{})
	/*
		DB.Create([]dbmodel.Subscription{
			{Name: "Basic", MaxNFT: 5, Price: 0},
			{Name: "Advanced", MaxNFT: 100, Price: 9},
		})

		DB.Create([]dbmodel.Collection{
			{
				Name:    "CryptoPunks",
				Address: "0xb47e3cd837ddf8e4c57f05d70ab865de6e193bbb",
				Price:   0,
				OsURL:   "https://opensea.io/zh-CN/collection/cryptopunks",
			},
			{
				Name:    "Bored Ape Yacht Club\n",
				Address: "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d",
				Price:   59.969,
				OsURL:   "https://opensea.io/zh-CN/collection/boredapeyachtclub",
			},
			{
				Name:    "Azuki\n",
				Address: "0xed5af388653567af2f388e6224dc7c4b3241c544",
				Price:   10.479,
				OsURL:   "https://opensea.io/zh-CN/collection/azuki",
			},
			{
				Name:    "CLONE X - X TAKASHI MURAKAMI\n\n",
				Address: "0x49cf6f5d44e70224e2e23fdcdd2c053f30ada28b",
				Price:   8.8,
				OsURL:   "https://opensea.io/zh-CN/collection/clonex",
			},
			{
				Name:    "Moonbirds\n",
				Address: "0x23581767a106ae21c074b2276D25e5C3e136a68b",
				Price:   7.4,
				OsURL:   "https://opensea.io/zh-CN/collection/proof-moonbirds",
			},
		})
	*/
	return
}
