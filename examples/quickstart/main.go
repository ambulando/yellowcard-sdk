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

	// List exchange rates for USD
	rates, err := client.Rates.List(ctx, "USD")
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range rates {
		fmt.Printf("%s: buy=%.4f sell=%.4f\n", r.Code, r.Buy, r.Sell)
	}

	// List payment channels for Ghana
	channels, err := client.Networks.Channels(ctx, "GH")
	if err != nil {
		log.Fatal(err)
	}
	for _, ch := range channels {
		fmt.Printf("channel: %s — %s (%s)\n", ch.Name, ch.Type, ch.Currency)
	}

	// Create a receive payment
	payment, err := client.Payments.Create(ctx, yellowcard.ReceivePaymentRequest{
		SequenceId:   "order-12345",
		Amount:       100,
		Currency:     "USD",
		Country:      "GH",
		ChannelId:    "channel-id",
		CustomerType: yellowcard.Retail,
		Reason:       "salary",
		Recipient: yellowcard.Recipient{
			Name:  "Jane Doe",
			Phone: "0241234567",
			Email: "jane@example.com",
		},
		Source: yellowcard.Source{
			AccountNumber: "0241234567",
			NetworkId:     "network-id",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Payment %s status: %s\n", payment.Id, payment.Status)

	// List payments with filters
	payments, err := client.Payments.List(ctx, yellowcard.ReceiveQuery{
		PerPage: 10,
		OrderBy: yellowcard.Desc,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range payments {
		fmt.Printf("%s: %s %d %s\n", p.Id, p.Status, p.Amount, p.Currency)
	}
}
