package csrf_svc

import (
	"context"
	"fmt"
	"microservices/auth/internal/clock_svc"
	"microservices/auth/pkg/csrf_pkg"
)

type Verifier struct {
	Csrf    csrf_pkg.CsrfPkgInterface
	Secrets string
	Clock   clock_svc.ClockInterface
}

func NewVerifier(csrf csrf_pkg.CsrfPkgInterface, secrets string, clock clock_svc.ClockInterface) *Verifier {
	return &Verifier{Csrf: csrf, Secrets: secrets, Clock: clock}
}

func (v *Verifier) Verify(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("missing csrf token")
	}
	secret := v.Csrf.GetSecret()
	return v.Csrf.ValidateCSRFCookieToken(token, secret, v.Clock.Now().Unix())
}
