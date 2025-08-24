package csrf

import "microservices/auth/pkg/csrf_pkg"

type CsrfSvcMockStruct struct{}

func (s *CsrfSvcMockStruct) CreateCSRFToken(csrf csrf_pkg.CsrfPkgInterface, timestamp int64) string {
	return "mocked_csrf_token"
}
