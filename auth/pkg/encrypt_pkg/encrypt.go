package encrypt_pkg

import "golang.org/x/crypto/bcrypt"

type EncryptPkgInterface interface {
	CreatePasswordHash(password string) (string, error)
}

type EncryptPkgStruct struct{}

func (e *EncryptPkgStruct) CreatePasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
