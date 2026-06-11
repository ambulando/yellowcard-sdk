package yellowcard

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// sign produces the HMAC-SHA256 signature expected by the YellowCard API.
// Signature covers: timestamp + HTTP method + path + raw request body.
func sign(secretKey, method, path, date string, body interface{}) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(date))
	h.Write([]byte(path))
	h.Write([]byte(method))
	if body != nil {
		bodyJSON, _ := json.Marshal(body)
		bodyHmac := sha256.Sum256(bodyJSON)
		bodyBase64 := base64.StdEncoding.EncodeToString(bodyHmac[:])
		h.Write([]byte(bodyBase64[:]))
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func authHeaders(apiKey, secretKey, method, path, body string) map[string]string {
	date := time.Now().UTC().Format(time.RFC3339)
	signature := sign(secretKey, method, path, date, body)
	return map[string]string{
		"X-YC-Timestamp": date,
		"Authorization":  fmt.Sprintf("YcHmacV1 %s:%s", apiKey, signature),
	}
}
