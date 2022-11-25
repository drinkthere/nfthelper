package service

import (
	"context"
	"fmt"
	"github.com/NdoleStudio/coinpayments-go"
	"net/http"
	"nfthelper/logger"
	"os"
	"strconv"
)

type PaymentService struct {
}

func (p *PaymentService) GeneratePaymentLink(network, currency string, price float64) (link string) {
	apiKey := os.Getenv("CP_APIKEY")
	apiSecret := os.Getenv("CP_APISECRET")
	originalCurrency := os.Getenv("ORIGINAL_CURRENCY")
	sendingCurrency := getSendingCurrency(network, currency)
	amount := strconv.FormatFloat(price, 'f', 5, 64)

	client := coinpayments.New(
		coinpayments.WithAPIKey(apiKey),
		coinpayments.WithAPISecret(apiSecret),
	)
	status, response, err := client.Payment.CreateTransaction(context.Background(), &coinpayments.CreatePaymentRequest{
		Amount:           fmt.Sprintf("%s", amount),
		OriginalCurrency: originalCurrency,
		SendingCurrency:  sendingCurrency,
		BuyerEmail:       "drinkthere@gmail.com",
	})
	if err != nil {
		logger.Error("generate payment link failed, network=%s, currency=%s, amount=%s, baseCurrency=%s, quoteCurrency=%s, error:%+v",
			network, currency, amount, originalCurrency, sendingCurrency, err)
		return
	}

	if response.HTTPResponse.StatusCode != http.StatusOK {
		logger.Error("generate payment link failed, network=%s, currency=%s, amount=%s, baseCurrency=%s, quoteCurrency=%s, error:%+v",
			network, currency, amount, originalCurrency, sendingCurrency, err)
		return
	}
	if status.Error != "ok" {
		logger.Error("generate payment link failed, network=%s, currency=%s, amount=%s, baseCurrency=%s, quoteCurrency=%s, error:%+v",
			network, currency, amount, originalCurrency, sendingCurrency, status.Error)
		return
	}
	logger.Info("generate payment link succ, network=%s, currency=%s, amount=%s, baseCurrency=%s, quoteCurrency=%s, checkout%s", status.Result.CheckoutURL)
	// todo 将交易写入数据库
	return status.Result.CheckoutURL
}

func getSendingCurrency(network, currency string) string {
	switch currency {
	case "USDT":
		switch network {
		case "ETH":
			return "USDT.ERC20"
		case "BSC":
			return "USDT.BEP20"
		case "Tron":
			return "USDT.TRC20"
		}
	case "USDC":
		switch network {
		case "ETH":
			return "USDC"
		case "BSC":
			return "USDC.BEP20"
		case "Tron":
			return "USDC.TRC20"
		}
	}
	return "USDC"
}
