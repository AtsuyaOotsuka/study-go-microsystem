package csrf_pkg_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"microservices/auth/pkg/csrf_pkg"
	"microservices/auth/tests/test_funcs"
	"strconv"
	"testing"
)

func TestGenerateCSRFCookieToken(t *testing.T) {

	csrf := csrf_pkg.CsrfPkgStruct{}

	csrfResult := csrf.GenerateCSRFCookieToken("my_secret_key", 1735689600, "aaaaaaaaaaa")
	expectedCsrf := "1735689600:aaaaaaaaaaa:f1c47b8855d489cca451a1ae8be2cfa46d413a5386354519ebf245df9cf11258"

	if csrfResult != expectedCsrf {
		t.Errorf("Expected %s, got %s", expectedCsrf, csrfResult)
	}
}

func TestGenerateNonceString_Length(t *testing.T) {
	csrf := csrf_pkg.CsrfPkgStruct{}

	nonce := csrf.GenerateNonceString()
	if len(nonce) < 44 { // Base64(32バイト) = 約44文字
		t.Errorf("Expected nonce length >= 44, got %d", len(nonce))
	}
}

func TestGetSecret(t *testing.T) {
	test_funcs.WithEnv("CSRF_TOKEN", "test_secret", t, func() {

		csrf := csrf_pkg.CsrfPkgStruct{}
		got := csrf.GetSecret()
		if got != "test_secret" {
			t.Errorf("Expected test_secret, got %s", got)
		}
	})
}

func TestGetSecret_Empty(t *testing.T) {

	csrf := csrf_pkg.CsrfPkgStruct{}
	test_funcs.WithEnv("CSRF_TOKEN", "", t, func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic when CSRF_TOKEN is empty")
			}
		}()
		csrf.GetSecret()
	})
}

func TestValidateCSRFCookieToken(t *testing.T) {

	csrf := csrf_pkg.CsrfPkgStruct{}
	secret := "my_secret_key"
	timestampStr := "1735689600"
	nonceStr := "aaaaaaaaaaa"
	token := timestampStr + ":" + nonceStr + ":f1c47b8855d489cca451a1ae8be2cfa46d413a5386354519ebf245df9cf11258"

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		t.Fatalf("Failed to parse timestamp: %v", err)
	}

	if err := csrf.ValidateCSRFCookieToken(token, secret, int64(timestamp)); err != nil {
		t.Errorf("Expected token to be valid, but it was not: %v", err)
	}
}

func TestValidateCSRFCookieToken_InvalidSignature(t *testing.T) {
	csrf := csrf_pkg.CsrfPkgStruct{}

	secret := "my_secret_key"
	timestampStr := "1735689600"
	nonceStr := "aaaaaaaaaaa"
	token := timestampStr + ":" + nonceStr + ":invalid_signature"
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		t.Fatalf("Failed to parse timestamp: %v", err)
	}
	if csrf.ValidateCSRFCookieToken(token, secret, int64(timestamp)) == nil {
		t.Errorf("Expected token to be invalid due to incorrect signature, but it was valid")
	}
}

func TestValidateCSRFCookieToken_ExpiryBoundary(t *testing.T) {
	impl := csrf_pkg.CsrfPkgStruct{}
	secret := "my_secret_key"
	nonce := "n"

	cases := []struct {
		name    string
		tsStr   string
		now     int64
		wantErr bool
	}{
		{"valid_at_+600", "1000", 1600, false},  // 差分=600 → 有効（>600 で失効）
		{"expired_at_+601", "1000", 1601, true}, // 差分=601 → 失効
	}

	for _, cse := range cases {
		t.Run(cse.name, func(t *testing.T) {
			data := cse.tsStr + ":" + nonce
			h := hmac.New(sha256.New, []byte(secret))
			h.Write([]byte(data))
			sig := hex.EncodeToString(h.Sum(nil))
			token := cse.tsStr + ":" + nonce + ":" + sig

			err := impl.ValidateCSRFCookieToken(token, secret, cse.now)
			if (err != nil) != cse.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, cse.wantErr)
			}
		})
	}
}

func TestValidateCSRFCookieToken_InvalidFormat(t *testing.T) {
	csrf := csrf_pkg.CsrfPkgStruct{}

	secret := "my_secret_key"
	token := "invalid_format_token"
	if csrf.ValidateCSRFCookieToken(token, secret, 0) == nil {
		t.Errorf("Expected token to be invalid due to incorrect format, but it was valid")
	}
}

func TestValidateCSRFCookieToken_InvalidTimestamp(t *testing.T) {
	csrf := csrf_pkg.CsrfPkgStruct{}

	secret := "my_secret_key"
	timestampStr := "not_a_number" // これが原因で ParseInt が失敗！
	nonceStr := "aaaaaaaaaaa"
	data := timestampStr + ":" + nonceStr
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sig := hex.EncodeToString(h.Sum(nil))

	token := data + ":" + sig

	// 引数として渡す timestamp は適当でOK（使われないので）
	if csrf.ValidateCSRFCookieToken(token, secret, 9999999999) == nil {
		t.Errorf("Expected token to be invalid due to non-numeric timestamp, but it was valid")
	}
}
