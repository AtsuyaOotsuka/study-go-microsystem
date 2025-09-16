package encrypt_pkg_mock

import "fmt"

type EncryptPkgMockStruct struct{}

func (e *EncryptPkgMockStruct) CreatePasswordHash(password string) (string, error) {
	return "mocked_hashed_password", nil
}

type EncryptPkgMockErrorStruct struct{}

func (e *EncryptPkgMockErrorStruct) CreatePasswordHash(password string) (string, error) {
	return "", fmt.Errorf("mocked error creating password hash")
}
