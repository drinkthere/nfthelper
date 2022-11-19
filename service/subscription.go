package service

import "nfthelper/model"

type SubscriptionService struct {
}

func (s *SubscriptionService) List() []model.Subscription {
	// todo 从redis中 获取subscription方案
	return []model.Subscription{
		{ID: 1, Name: "Advanced", MaxNFT: 100, Price: 9},
	}
}

func (s *SubscriptionService) GetByID(id int64) model.Subscription {
	return model.Subscription{
		ID: 1, Name: "Advanced", MaxNFT: 100, Price: 9,
	}
}

func (s *SubscriptionService) GetByUserID(uid int64) model.Subscription {
	return model.Subscription{
		ID: 1, Name: "Basic", MaxNFT: 5, Price: 9,
	}
}
