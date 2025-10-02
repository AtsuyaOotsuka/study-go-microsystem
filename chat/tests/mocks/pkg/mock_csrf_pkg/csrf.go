package mock_csrf_pkg

import "github.com/stretchr/testify/mock"

type CsrfPkgMockStruct struct {
	mock.Mock
}

func (m *CsrfPkgMockStruct) GenerateNonceString() string {
	args := m.Called()
	return args.String(0)
}
func (m *CsrfPkgMockStruct) GetSecret() string {
	args := m.Called()
	return args.String(0)
}
func (m *CsrfPkgMockStruct) GenerateCSRFCookieToken(secret string, timestamp int64, nonceStr string) string {
	args := m.Called(secret, timestamp, nonceStr)
	return args.String(0)
}
func (m *CsrfPkgMockStruct) ValidateCSRFCookieToken(token string, secret string, timestamp int64) error {
	args := m.Called(token, secret, timestamp)
	return args.Error(0)
}
