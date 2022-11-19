package service

type CommonService struct {
}

func (s *CommonService) ListCurrencies() []string {
	return []string{"USDT", "USDC"}
}

func (s *CommonService) ListNetworks() []string {
	return []string{"ETH", "BSC", "Tron"}
}
