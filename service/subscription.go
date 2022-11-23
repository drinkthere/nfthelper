package service

import (
	"fmt"
	"nfthelper/database"
	"nfthelper/logger"
	"nfthelper/model"
	"nfthelper/model/dbmodel"
)

type SubscriptionService struct {
}

func (s *SubscriptionService) List() []model.Subscription {
	// todo 从redis中 获取subscription方案
	return []model.Subscription{
		{ID: 1, Name: "Advanced", MaxNFT: 100, Price: 9},
	}
}

func (s *SubscriptionService) GetByID(id uint) model.Subscription {
	return model.Subscription{
		ID: 1, Name: "Advanced", MaxNFT: 100, Price: 9,
	}
}

func (s *SubscriptionService) GetByUserID(uid uint) (subscription model.Subscription) {
	msg := fmt.Sprintf("get user subscription uid=%d", uid)
	logger.Info(msg)
	var subscriptions []model.Subscription
	result := database.DB.Table("user_subscription").Select("subscription.id, subscription.name, subscription.price, subscription.max_nft").Joins("left join subscription on user_subscription.subscription_id = subscription.id where user_id=?", uid).Scan(&subscriptions)
	if result.Error != nil {
		logger.Error("%s, error is %+v", msg, result.Error)
	}
	if len(subscriptions) > 0 {
		subscription = subscriptions[0]
	}
	return
}

func (s *SubscriptionService) SetBasicSubscription(uid uint, isFirst bool) {
	msg := "set user subscription to basic plan"
	logger.Info(msg)
	if isFirst {
		// 添加新纪录
		us := dbmodel.UserSubscription{
			UserID:         uid,
			SubscriptionID: 1, // basic subscription id is 1
		}
		result := database.DB.Create(&us)
		if result.Error != nil {
			logger.Error("%s, error is %+v", msg, result.Error)
		}
	} else {
		// 更新记录
		var us dbmodel.UserSubscription
		result := database.DB.First(&us, "user_id = ?", uid)
		if result.Error != nil {
			logger.Error("%s, error is %+v", msg, result.Error)
			return
		}

		us.SubscriptionID = 1
		result = database.DB.Save(&us)
		if result.Error != nil {
			logger.Error("%s, error is %+v", msg, result.Error)
		}
	}
	return
}
