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
}
