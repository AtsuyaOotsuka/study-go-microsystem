package csrf_svc

import (
	"microservices/auth/tests/mocks/pkg/csrf"
	"testing"
)

func TestCreateCSRFToken(t *testing.T) {
	csrfMock := &csrf.CsrfPkgMockStruct{}

	cvs := CsrfSvcStruct{}

	token := cvs.CreateCSRFToken(csrfMock, 1234567890)

	if token != "mocked_csrf_token" {
		t.Errorf("expected 'mocked_csrf_token', got '%s'", token)
	}
}
