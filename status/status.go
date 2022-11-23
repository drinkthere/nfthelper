package status

import (
	"nfthelper/model"
	"sync"
	"time"
)

var (
	indicatorMap  map[uint]model.Indicator
	indicatorLock sync.RWMutex
)

const (
	Start            = "start"
	AddNFT           = "add nft"
	ConfirmNFT       = "confirm nft"
	DeleteNFT        = "delete nft"
	Subscription     = "subscription"
	ListSubscription = "list subscription"
	ChooseCurrency   = "choose currency"
	ChooseNetwork    = "choose currency"
)

// 初始化日志
func InitStatus() {
	indicatorMap = make(map[uint]model.Indicator)
}

func GetIndicator(userID uint) string {
	indicatorLock.RLock()
	defer indicatorLock.RUnlock()
	if indicator, ok := indicatorMap[userID]; ok {
		return indicator.Value
	}
	return ""
}

func SetIndicator(userID uint, indicator string) {
	indicatorLock.Lock()
	defer indicatorLock.Unlock()

	indicatorMap[userID] = model.Indicator{
		Value:          indicator,
		LastUpdateTime: time.Now(),
	}
}

func DelIndicator(duration float64) {
	indicatorLock.Lock()
	defer indicatorLock.Unlock()

	for userID, indicator := range indicatorMap {
		if time.Now().Sub(indicator.LastUpdateTime).Seconds() > duration {
			delete(indicatorMap, userID)
		}
	}
}
