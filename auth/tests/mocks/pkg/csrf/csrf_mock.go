package csrf

type CsrfPkgMockStruct struct{}

func (c *CsrfPkgMockStruct) GenerateNonceString() string {
	return "mocked_nonce"
}

func (c *CsrfPkgMockStruct) GetSecret() string {
	return "mocked_secret"
}

func (c *CsrfPkgMockStruct) GenerateCSRFCookieToken(secret string, timestamp int64, nonceStr string) string {
	return "mocked_csrf_token"
}

func (c *CsrfPkgMockStruct) ValidateCSRFCookieToken(token string, secret string, timestamp int64) error {
	return nil
}
