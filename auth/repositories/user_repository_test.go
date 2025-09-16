package repositories

import (
	"database/sql"
	"microservices/auth/tests/mocks/global_mock"
	"microservices/auth/tests/mocks/models_mock"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestGetByEmail(t *testing.T) {

	mockUser := models_mock.CreateUserMock()

	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	rows := sqlmock.NewRows([]string{"id", "password", "email"}).
		AddRow(1, mockUser.Password, mockUser.Email)
	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE email = \\?").
		WithArgs(mockUser.Email, sqlmock.AnyArg()).
		WillReturnRows(rows)

	defer cleanup()
	repo := &UserRepositoryStruct{Db: gdb}
	user, err := repo.GetByEmail(mockUser.Email)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if user.Email != mockUser.Email {
		t.Errorf("expected email %s, but got %s", mockUser.Email, user.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestGetByEmail_DBError(t *testing.T) {

	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE email = \\?").
		WithArgs("dberror@example.com", sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	defer cleanup()
	repo := &UserRepositoryStruct{Db: gdb}
	user, err := repo.GetByEmail("dberror@example.com")
	if err == nil {
		t.Fatalf("expected error, but got user: %v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestGetByEmail_NotFound(t *testing.T) {

	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE email = \\?").
		WithArgs("notfound@example.com", sqlmock.AnyArg()).
		WillReturnError(gorm.ErrRecordNotFound)

	defer cleanup()
	repo := &UserRepositoryStruct{Db: gdb}
	user, err := repo.GetByEmail("notfound@example.com")
	if err == nil {
		t.Fatalf("expected error, but got user: %v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestGetByRefreshToken(t *testing.T) {
	mockUser := models_mock.CreateUserMock()

	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	rows := sqlmock.NewRows([]string{"id", "email", "refresh_token"}).
		AddRow(1, mockUser.Email, mockUser.RefreshToken)
	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE refresh_token = \\?").
		WithArgs(mockUser.RefreshToken, sqlmock.AnyArg()).
		WillReturnRows(rows)

	defer cleanup()
	repo := &UserRepositoryStruct{Db: gdb}
	user, err := repo.GetByRefreshToken(mockUser.RefreshToken)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if user.RefreshToken != mockUser.RefreshToken {
		t.Errorf("expected refresh token %s, but got %s", mockUser.RefreshToken, user.RefreshToken)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestGetByRefreshToken_DBError(t *testing.T) {

	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE refresh_token = \\?").
		WithArgs("dberror_refresh_token", sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	defer cleanup()
	repo := &UserRepositoryStruct{Db: gdb}
	user, err := repo.GetByRefreshToken("dberror_refresh_token")
	if err == nil {
		t.Fatalf("expected error, but got user: %v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestGetByRefreshToken_NotFound(t *testing.T) {

	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE refresh_token = \\?").
		WithArgs("notfound_refresh_token", sqlmock.AnyArg()).
		WillReturnError(gorm.ErrRecordNotFound)

	defer cleanup()
	repo := &UserRepositoryStruct{Db: gdb}
	user, err := repo.GetByRefreshToken("notfound_refresh_token")
	if err == nil {
		t.Fatalf("expected error, but got user: %v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
