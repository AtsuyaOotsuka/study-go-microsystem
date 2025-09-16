package encrypt_pkg

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestCreatePasswordHash(t *testing.T) {
	password := "my_secure_password"
	encryptor := &EncryptPkgStruct{}
	hashedPassword, err := encryptor.CreatePasswordHash(password)
	if err != nil {
		t.Fatalf("failed to create password hash: %v", err)
	}

	if hashedPassword == "" {
		t.Fatal("hashed password should not be empty")
	}

	if hashedPassword == password {
		t.Fatal("hashed password should not be the same as the original password")
	}

	// bcryptのCompareHashAndPasswordを使ってハッシュが正しいことを確認
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		t.Fatalf("hashed password does not match the original password: %v", err)
	}
}

func TestCreatePasswordHash_withError(t *testing.T) {
	encryptor := &EncryptPkgStruct{}
	randomString := strings.Repeat("A", 73)
	// bcryptのパスワード長制限は72文字なので、オーバーフローさせてエラーを発生させる
	hashedPassword, err := encryptor.CreatePasswordHash(randomString)

	if hashedPassword != "" {
		t.Fatal("hashed password should be empty")
	}

	if err == nil {
		t.Fatal("error should not be nil")
	}
}
