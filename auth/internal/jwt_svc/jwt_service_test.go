package jwt_svc

import (
	"crypto/rand"
	"crypto/rsa"
	"microservices/auth/models"
	"microservices/auth/tests/mocks/svc_internal/clock"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewJwtService(t *testing.T) {
	jwtService := NewJwtService()
	assert.IsType(t, &JwtServiceStruct{}, jwtService)
}

func TestCreateJwt(t *testing.T) {
	mockUser := &models.User{
		ID:    1,
		Email: "test@example.com",
	}

	jwt_svc := JwtServiceStruct{
		Clock:  clock.FixedClock{FixedTime: time.Now()},
		Method: jwt.SigningMethodHS256,
		Key:    []byte(os.Getenv("JWT_SECRET")),
	}

	// JWTトークンを作成
	token, err := jwt_svc.CreateJwt(mockUser)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	// トークンが正しい形式であることを確認
	assert.NotEmpty(t, token)

	// トークンを検証
	claims, err := jwt_svc.ValidateJwt(token)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}

	// クレームが正しいことを確認
	assert.Equal(t, (int)(mockUser.ID), claims.UserID)
	assert.Equal(t, (string)(mockUser.Email), claims.Email)
}

func TestValidateJwt_InvalidToken(t *testing.T) {
	jwt_svc := JwtServiceStruct{
		Clock:  clock.FixedClock{FixedTime: time.Now()},
		Method: jwt.SigningMethodHS256,
		Key:    []byte(os.Getenv("JWT_SECRET")),
	}

	// 無効なトークンを検証
	claims, err := jwt_svc.ValidateJwt("invalid.token.string")
	assert.Nil(t, claims)
	assert.Error(t, err)
}

func TestValidateJwt_ExpiredToken(t *testing.T) {
	jwt_svc := JwtServiceStruct{
		Clock:  clock.FixedClock{FixedTime: time.Now().Add(-2 * time.Hour)},
		Method: jwt.SigningMethodHS256,
		Key:    []byte(os.Getenv("JWT_SECRET")),
	}

	mockUser := &models.User{
		ID:    1,
		Email: "test@example.com",
	}

	// トークンを作成
	token, err := jwt_svc.CreateJwt(mockUser)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	// トークンを検証
	claims, err := jwt_svc.ValidateJwt(token)
	assert.Nil(t, claims)
	assert.Error(t, err)
}

func TestValidateJwt_UnexpectedSigningMethod(t *testing.T) {
	// 1) RS256の鍵をその場で生成（外部依存なし）
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	// 2) RS256 で署名された JWT を作る（ペイロード内容は何でもOK）
	claims := jwt.RegisteredClaims{
		Subject:   "1",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	}
	rsTok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	rsTokenString, err := rsTok.SignedString(priv)
	if err != nil {
		t.Fatal(err)
	}

	// 3) サービス経由で検証 → keyfunc 内で alg チェックに引っかかるはず
	svc := JwtServiceStruct{Clock: clock.FixedClock{FixedTime: time.Now()}}

	got, verr := svc.ValidateJwt(rsTokenString)

	assert.Nil(t, got)    // クレームは返らない
	assert.Error(t, verr) // エラーになる
	// ライブラリや実装によってはメッセージをラップしてることもあるので、Contains は任意
	// assert.Contains(t, verr.Error(), "unexpected signing method")
}

func TestCreateJwt_SignedString_Error_ByNilKey(t *testing.T) {
	svc := &JwtServiceStruct{
		Clock:  clock.FixedClock{FixedTime: time.Now()},
		Method: jwt.SigningMethodHS256,
		Key:    nil, // ← nilならエラー
	}

	_, err := svc.CreateJwt(&models.User{ID: 1, Email: "test@example.com"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
