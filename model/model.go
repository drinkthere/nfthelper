package model

import "time"

type Indicator struct {
	Value          string
	LastUpdateTime time.Time
}

type Collection struct {
	ID      int64
	Name    string
	Address string
	Price   float64
	OsURL   string
}

type Subscription struct {
	ID     int64
	Name   string
	Price  float64 // USD
	MaxNFT int64   // max NFT support
}

type Announcement struct {
	ID  int64
	URL string
}
