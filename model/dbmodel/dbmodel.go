package dbmodel

import "gorm.io/gorm"

type Subscription struct {
	gorm.Model
	Name   string  `gorm:"type:varchar(20) not null;unique"`
	Price  float64 `gorm:"type:float"`
	MaxNFT int     // max NFT support
}

type Collection struct {
	gorm.Model
	Name    string  `gorm:"type:varchar(100);unique"`
	Address string  `gorm:"type:char(50);unique"`
	Price   float64 `gorm:"type:float"`
	OsURL   string  `gorm:"type:varchar(255)"`
}

type UserSubscription struct {
	gorm.Model
	UserID         uint `gorm:"unique"`
	SubscriptionID uint
}

type UserCollection struct {
	gorm.Model
	UserID       uint `gorm:"uniqueIndex:idx_user_col"`
	CollectionID uint `gorm:"uniqueIndex:idx_user_col"`
}

type Announcement struct {
	gorm.Model
	CollectionID uint
	URL          string `gorm:"unique"`
}
