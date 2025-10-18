package test_funcs

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func GenerateCSRFCookieToken(
	secret string,
	timestamp int64,
) string {
	nonce := make([]byte, 32)
	rand.Read(nonce)

	data := fmt.Sprintf("%d:%s", timestamp, base64.StdEncoding.EncodeToString(nonce))

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sig := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s:%s", data, sig)
}
