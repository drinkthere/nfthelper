package service

import (
	"fmt"
	"nfthelper/database"
	"nfthelper/logger"
	"nfthelper/model/dbmodel"
)

type SubscriptionService struct {
}

func (s *SubscriptionService) List() (subscriptions []dbmodel.Subscription) {
	msg := fmt.Sprintf("list subscriptions")
	logger.Info(msg)

	// list 订阅方案，不包含basic plan
	result := database.DB.Where("id != ?", 1).Find(&subscriptions)
	if result.Error != nil {
		logger.Error("%s, error is %+v", msg, result.Error)
		return
	}
	return
}

func (s *SubscriptionService) GetByID(id uint) (subscription dbmodel.Subscription) {
	msg := fmt.Sprintf("get subscription by id=%d", id)
	logger.Info(msg)

	// list 订阅方案，不包含basic plan
	result := database.DB.Where("id=", 1).First(&subscription)
	if result.Error != nil {
		logger.Error("%s, error is %+v", msg, result.Error)
		return
	}
	return
}

func (s *SubscriptionService) GetByUserID(uid uint) (subscription dbmodel.Subscription) {
	msg := fmt.Sprintf("get user subscription uid=%d", uid)
	logger.Info(msg)
	var subscriptions []dbmodel.Subscription
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
