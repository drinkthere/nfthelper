package model

import "time"

type Indicator struct {
	Value          string
	LastUpdateTime time.Time
}
