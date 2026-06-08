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

Look up exchange rates between currency pairs.

```go
ctx := context.Background()

// All rates, optionally filtered
rates, err := client.Rates.List(ctx, yellowcard.RatesParams{From: "USD"})

// Single pair
rate, err := client.Rates.Get(ctx, "USD", "GHS")
fmt.Printf("1 USD = %.4f GHS (expires %s)\n", rate.Rate, rate.ExpiresAt)
```

### Networks and channels

Discover supported blockchain networks and payment channels (mobile money, bank transfer, etc.).

```go
// Supported networks
networks, err := client.Networks.List(ctx)

// Payment channels, optionally filtered by ISO country code
channels, err := client.Networks.Channels(ctx, "GH")
for _, ch := range channels {
    fmt.Printf("%s — %s (%s)\n", ch.Name, ch.Type, ch.Currency)
}
```

### Payments

```go
// Create a payment
payment, err := client.Payments.Create(ctx, yellowcard.CreatePaymentRequest{
    SequenceID: "order-12345",   // idempotency key you assign
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

// Fetch by API-assigned ID
payment, err = client.Payments.Get(ctx, payment.ID)

// Fetch by your own sequence ID
payment, err = client.Payments.GetBySequenceID(ctx, "order-12345")

// Cancel a pending payment
err = client.Payments.Cancel(ctx, payment.ID)
```

Payment lifecycle statuses: `pending` → `processing` → `completed` / `failed` / `expired`.

### Accounts

```go
accounts, err := client.Accounts.List(ctx)
for _, a := range accounts {
    fmt.Printf("%s: %.2f %s (%s)\n", a.Label, a.Balance, a.Currency, a.Status)
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
    // full error detail
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
import "net/http"
import "time"

hc := &http.Client{Timeout: 10 * time.Second}
client := yellowcard.New(apiKey, secretKey, yellowcard.WithHTTPClient(hc))
```

## Authentication

The client handles authentication automatically. Each request is signed with HMAC-SHA256 over `timestamp + method + path + body` and sent via the `YC-API-Key`, `YC-Signature`, and `YC-Timestamp` headers. You only need to supply your API key and secret to `New()`.

## Running the example

```bash
cd examples/quickstart
go run main.go
```

## License

See [LICENSE](LICENSE.md).
