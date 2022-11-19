package status

import (
	"nfthelper/model"
	"sync"
	"time"
)

var (
	indicatorMap  map[int64]model.Indicator
	indicatorLock sync.RWMutex
)

const (
	Start            = "start"
	AddNFT           = "add nft"
	ConfirmNFT       = "confirm nft"
	DeleteNFT        = "delete nft"
	ConfirmDeleteNFT = "confirm delete nft"
)

// 初始化日志
func InitStatus() {
	indicatorMap = make(map[int64]model.Indicator)
}

func GetIndicator(userID int64) string {
	indicatorLock.RLock()
	defer indicatorLock.RUnlock()
	if indicator, ok := indicatorMap[userID]; ok {
		return indicator.Value
	}
	return ""
}

func SetIndicator(userID int64, indicator string) {
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
