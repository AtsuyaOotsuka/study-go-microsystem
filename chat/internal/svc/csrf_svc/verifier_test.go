package csrf_svc

import (
	"context"
	"microservices/chat/internal/svc/clock_svc"
	"microservices/chat/tests/mocks/pkg/csrf_pkg"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestVerify(t *testing.T) {
	csrfPkgMock := &csrf_pkg.CsrfPkgMockStruct{}
	csrfPkgMock.On("GetSecret").Return("secrets")
	csrfPkgMock.On("ValidateCSRFCookieToken", "valid_token", "secrets", mock.AnythingOfType("int64")).Return(nil)

	verifier := NewVerifier(
		csrfPkgMock,
		"secrets",
		clock_svc.RealClockStruct{},
	)

	err := verifier.Verify(context.Background(), "valid_token")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestVerifyMissingToken(t *testing.T) {
	csrfPkgMock := &csrf_pkg.CsrfPkgMockStruct{}
	csrfPkgMock.On("GetSecret").Return("secrets")

	verifier := NewVerifier(
		csrfPkgMock,
		"secrets",
		clock_svc.RealClockStruct{},
	)

	err := verifier.Verify(context.Background(), "")
	if err == nil {
		t.Errorf("expected error, got none")
	}
}
