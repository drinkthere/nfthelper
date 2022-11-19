package service

type PaymentService struct {
}

func (p *PaymentService) GeneratePaymentLink(network, currency string, price float64) string {
	return "https://www.coinpayments.net/index.php?cmd=checkout&id="
}
