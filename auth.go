package yellowcard

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

// sign produces the HMAC-SHA256 signature expected by the YellowCard API.
// Signature covers: timestamp + HTTP method + path + raw request body.
func sign(secretKey, method, path, body string, ts time.Time) string {
	timestamp := strconv.FormatInt(ts.UnixMilli(), 10)
	payload := timestamp + method + path + body
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func authHeaders(apiKey, secretKey, method, path, body string) map[string]string {
	ts := time.Now()
	return map[string]string{
		"YC-API-Key":   apiKey,
		"YC-Signature": sign(secretKey, method, path, body, ts),
		"YC-Timestamp": fmt.Sprintf("%d", ts.UnixMilli()),
	}
}
