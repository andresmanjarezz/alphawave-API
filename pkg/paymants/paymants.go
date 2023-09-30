package paymants

import (
	"fmt"

	"github.com/Coke15/AlphaWave-BackEnd/pkg/logger"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

type PaymentProvider struct {
	StripeAPIKey string
}

func NewPaymentProvider(stripeAPIKey string) *PaymentProvider {
	stripe.Key = stripeAPIKey
	return &PaymentProvider{
		StripeAPIKey: stripeAPIKey,
	}
}

type PaymantPayload struct {
	Amount   int64
	Currency string
}

func (p *PaymentProvider) Paymant(input PaymantPayload) {

	params := &stripe.ChargeParams{
		Amount:      stripe.Int64(input.Amount),
		Currency:    stripe.String(input.Currency),
		Description: stripe.String("Test Pay"),
	}

	params.SetSource("tok_visa")

	ch, err := charge.New(params)

	if err != nil {
		logger.Error(err)
	}

	fmt.Printf("ID: %s, Amount: %d, Descr: %s", ch.ID, ch.Amount, ch.Description)
}
