# yellowcard-sdk

Go client for the [YellowCard](https://yellowcard.io) API — a crypto exchange and payment platform serving African markets.

## Installation

```bash
go get github.com/ambulando/yellowcard-sdk
```

Requires Go 1.22 or later.

## Quick start

```go
import yellowcard "github.com/ambulando/yellowcard-sdk"

client := yellowcard.New("your-api-key", "your-secret-key")
```

Use `WithSandbox()` during development to target the sandbox environment:

```go
client := yellowcard.New("your-api-key", "your-secret-key", yellowcard.WithSandbox())
```

## Services

### Rates

Look up exchange rates for a given currency.

```go
ctx := context.Background()

rates, err := client.Rates.List(ctx, "USD")
for _, r := range rates {
    fmt.Printf("%s: buy=%.4f sell=%.4f\n", r.Code, r.Buy, r.Sell)
}
```

### Networks and channels

Discover supported networks and payment channels (mobile money, bank transfer, etc.), optionally filtered by ISO country code.

```go
// Supported networks
networks, err := client.Networks.List(ctx, "GH")

// Payment channels
channels, err := client.Networks.Channels(ctx, "GH")
for _, ch := range channels {
    fmt.Printf("%s — %s (%s)\n", ch.Name, ch.Type, ch.Currency)
}
```

### Payments

```go
// Create a receive payment
payment, err := client.Payments.Create(ctx, yellowcard.ReceivePaymentRequest{
    SequenceId:   "order-12345", // idempotency key you assign
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

// Fetch by ID
payment, err = client.Payments.Get(ctx, payment.Id)

// Fetch by your own sequence ID
payment, err = client.Payments.GetBySequenceID(ctx, "order-12345")

// List with filters
payments, err := client.Payments.List(ctx, yellowcard.ReceiveQuery{
    PerPage: 20,
    OrderBy: yellowcard.Desc,
    SortBy:  yellowcard.CreatedAt,
})

// Lifecycle actions
payment, err = client.Payments.Accept(ctx, payment.Id)
payment, err = client.Payments.Deny(ctx, payment.Id)
payment, err = client.Payments.Cancel(ctx, payment.Id)
payment, err = client.Payments.Refund(ctx, payment.Id)
```

### Accounts

```go
accounts, err := client.Accounts.List(ctx)
for _, a := range accounts {
    fmt.Printf("%s: %.2f %s\n", a.CurrencyType, a.Available, a.Currency)
}
```

## Error handling

API errors are returned as `*yellowcard.APIError`, which carries the HTTP status code plus the API error code and message.

```go
payment, err := client.Payments.Get(ctx, "unknown-id")
if err != nil {
    if yellowcard.IsNotFound(err) {
        // 404 — payment does not exist
    }
    if yellowcard.IsUnauthorized(err) {
        // 401 — check your API key and secret
    }
    log.Fatal(err)
}
```

## Configuration options

| Option | Description |
|--------|-------------|
| `WithSandbox()` | Target `https://sandbox.yellowcard.io` |
| `WithBaseURL(url)` | Override the base URL (e.g. for a local mock) |
| `WithHTTPClient(hc)` | Supply a custom `*http.Client` (timeouts, tracing, etc.) |

```go
hc := &http.Client{Timeout: 10 * time.Second}
client := yellowcard.New(apiKey, secretKey, yellowcard.WithHTTPClient(hc))
```

## Authentication

The client handles authentication automatically. Each request is signed with HMAC-SHA256 and sent via the `X-YC-Timestamp` and `Authorization` headers. You only need to supply your API key and secret to `New()`.

## Running the example

```bash
cd examples/quickstart
go run main.go
```

## License

See [LICENSE](LICENSE.md).
