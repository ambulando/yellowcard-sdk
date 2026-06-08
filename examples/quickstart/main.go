package main

import (
	"context"
	"fmt"
	"log"

	yellowcard "github.com/ambulando/yellowcard-sdk"
)

func main() {
	client := yellowcard.New(
		"your-api-key",
		"your-secret-key",
		yellowcard.WithSandbox(),
	)

	ctx := context.Background()

	// List exchange rates
	rates, err := client.Rates.List(ctx, yellowcard.RatesParams{From: "USD"})
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range rates {
		fmt.Printf("%s → %s: %.4f\n", r.From, r.To, r.Rate)
	}

	// Create a payment
	payment, err := client.Payments.Create(ctx, yellowcard.CreatePaymentRequest{
		SequenceID: "order-12345",
		Amount:     100,
		Currency:   "USD",
		ChannelID:  "channel-id",
		Destination: yellowcard.Destination{
			AccountName:   "Jane Doe",
			AccountNumber: "0241234567",
			Country:       "GH",
		},
		Reason: "salary",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Payment %s status: %s\n", payment.ID, payment.Status)
}
