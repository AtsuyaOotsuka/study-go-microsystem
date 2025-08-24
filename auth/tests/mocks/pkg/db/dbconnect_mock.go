package db

import (
	"database/sql"
)

type DBConnectMock struct {
}

func NewDBConnectMock() *DBConnectMock {
	return &DBConnectMock{}
}

func (m *DBConnectMock) ConnectDB() (*sql.DB, error) {
	// モックの DB 接続を返す
	return &sql.DB{}, nil
}
