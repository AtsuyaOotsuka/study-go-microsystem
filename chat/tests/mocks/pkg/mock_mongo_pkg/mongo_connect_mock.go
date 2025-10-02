package mock_mongo_pkg

import (
	"microservices/chat/pkg/mongo_pkg"

	"github.com/stretchr/testify/mock"
)

type MongoPkgMock struct {
	mock.Mock
}

func (m *MongoPkgMock) NewMongoConnect(dbName string) (*mongo_pkg.MongoPkgStruct, error) {
	args := m.Called(dbName)
	return args.Get(0).(*mongo_pkg.MongoPkgStruct), args.Error(1)
}

type MongoPkgWithErrorMock struct {
	mock.Mock
}

func (m *MongoPkgWithErrorMock) NewMongoConnect(dbName string) (*mongo_pkg.MongoPkgStruct, error) {
	args := m.Called(dbName)
	return nil, args.Error(1)
}
