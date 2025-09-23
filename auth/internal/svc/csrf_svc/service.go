package csrf_svc

import (
	"microservices/auth/pkg/csrf_pkg"
)

type CsrfSvcInterface interface {
	CreateCSRFToken(
		csrf csrf_pkg.CsrfPkgInterface,
		timestamp int64,
	) string
}

type CsrfSvcStruct struct{}

func (s *CsrfSvcStruct) CreateCSRFToken(csrf csrf_pkg.CsrfPkgInterface, timestamp int64) string {
	secret := csrf.GetSecret()
	nonceStr := csrf.GenerateNonceString()
	return csrf.GenerateCSRFCookieToken(secret, timestamp, nonceStr)
}
