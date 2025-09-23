package csrf_svc

import (
	"context"
	"microservices/auth/internal/svc/clock_svc"
	"microservices/auth/tests/mocks/pkg/csrf"
	"testing"
)

func TestVerify(t *testing.T) {
	verifier := NewVerifier(
		&csrf.CsrfPkgMockStruct{},
		"secrets",
		clock_svc.RealClockStruct{},
	)

	err := verifier.Verify(context.Background(), "valid_token")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestVerifyMissingToken(t *testing.T) {
	verifier := NewVerifier(
		&csrf.CsrfPkgMockStruct{},
		"",
		clock_svc.RealClockStruct{},
	)

	err := verifier.Verify(context.Background(), "")
	if err == nil {
		t.Errorf("expected error, got none")
	}
}
