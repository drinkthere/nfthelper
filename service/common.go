package service

import "nfthelper/model"

type CommonService struct {
}

func (c *CommonService) ListSubscriptions() []model.Subscription {
	// todo 从redis中 获取subscription方案
	return []model.Subscription{
		{ID: 1, Name: "Advanced", MaxNFT: 100, Price: 9},
	}
}

func (c *CommonService) GetSubscription(id int64) model.Subscription {
	return model.Subscription{
		ID: 1, Name: "Advanced", MaxNFT: 100, Price: 9,
	}
}

func (c *CommonService) ListCurrencies() []string {
	return []string{"USDT", "USDC"}
}

func (c *CommonService) ListNetworks() []string {
	return []string{"ETH", "BSC", "Tron"}
}
