package jwtinfo_svc

import (
	"context"
	"testing"
)

func TestNewJwtInfo(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserIDKey, 12345)
	ctx = context.WithValue(ctx, EmailKey, "test@example.com")
	jwtinfo := NewJwtInfo(ctx)

	if jwtinfo.UserID != 12345 {
		t.Errorf("expected UserID to be 12345, got %d", jwtinfo.UserID)
	}
	if jwtinfo.Email != "test@example.com" {
		t.Errorf("expected Email to be test@example.com, got %s", jwtinfo.Email)
	}
}
