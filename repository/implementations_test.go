package repository

import (
	"database/sql"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"context"
	"fmt"
	"github.com/go-test/deep"
	"github.com/google/uuid"
	"regexp"
	"testing"
)

var err error

type TestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	db   *sql.DB
	r    *Repository
	ctx  context.Context

	user                      User
	userInfo                  UserInfo
	userSecret                UserSecret
	getUserByPhoneNumberInput GetUserByPhoneNumberInput
	getUserByIDInput          GetUserByIDInput
	getUserByPhoneNumber      GetUserByPhoneNumberInput
	getUserByFullName         GetUserByFullNameInput
}

func (s *TestSuite) SetupSuite() {
	s.db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.r = &Repository{
		s.db,
	}
	s.ctx = context.TODO()
	s.userInfo = UserInfo{
		FullName:    "Test test",
		PhoneNumber: "+62123456789",
	}
	s.userSecret = UserSecret{
		Password:     "$2a$12$0cRbL1Sy8/WayYcISedNyee90dMyhwQXtDZHRppECb28SUTNAHbRO",
		PasswordSalt: "test123",
	}

	s.user = User{
		ID:         uuid.New(),
		UserInfo:   s.userInfo,
		UserSecret: s.userSecret,
	}
	s.getUserByPhoneNumberInput = GetUserByPhoneNumberInput{
		PhoneNumber: s.user.PhoneNumber,
	}
	s.getUserByIDInput = GetUserByIDInput{
		ID: s.user.ID,
	}
	s.getUserByFullName = GetUserByFullNameInput{
		FullName: s.user.FullName,
	}
}

func (s *TestSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestSuiteInit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestInsertUserSuccess() {
	prepare := s.mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO users(phone_number, full_name, password, password_salt) VALUES($1,$2,$3,$3) RETURNING id"))
	prepare.ExpectQuery().
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(
				s.user.ID,
			),
		).
		WithArgs(
			s.user.PhoneNumber,
			s.user.FullName,
			s.user.Password,
			s.user.PasswordSalt,
		)
	output, err := s.r.InsertUser(s.ctx, s.user)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(output.ID, s.user.ID))
}

func (s *TestSuite) TestInsertUserFailedToPrepareQuery() {
	s.mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO users(phone_number, full_name, password, password_salt) VALUES($1,$2,$3,$3) RETURNING id")).
		WillReturnError(fmt.Errorf("faield to prepare query"))
	output, err := s.r.InsertUser(s.ctx, s.user)
	require.Error(s.T(), err)
	require.Nil(s.T(), deep.Equal(output, InsertUserOutput{}))
}

func (s *TestSuite) TestInsertUserFailed() {
	prepare := s.mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO users(phone_number, full_name, password, password_salt) VALUES($1,$2,$3,$3) RETURNING id"))
	prepare.ExpectQuery().
		WillReturnError(fmt.Errorf("failed to insert user")).
		WithArgs(
			s.user.PhoneNumber,
			s.user.FullName,
			s.user.Password,
			s.user.PasswordSalt,
		)
	output, err := s.r.InsertUser(s.ctx, s.user)
	require.Error(s.T(), err)
	require.Nil(s.T(), deep.Equal(output, InsertUserOutput{}))
}

func (s *TestSuite) TestGetUserByPhoneNumberSuccess() {
	prepare := s.mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, phone_number, full_name, password, password_salt FROM users WHERE phone_number = $1	"))
	prepare.ExpectQuery().
		WillReturnRows(sqlmock.NewRows([]string{"id", "phone_number", "full_name", "password", "password_salt"}).AddRow(
			s.user.ID,
			s.user.PhoneNumber,
			s.user.FullName,
			s.user.Password,
			s.user.PasswordSalt,
		)).
		WithArgs(
			s.user.PhoneNumber,
		)
	output, err := s.r.GetUserByPhoneNumber(s.ctx, s.getUserByPhoneNumberInput)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(output, s.user))
}

func (s *TestSuite) TestGEtUserByPhoneNumberFailedPrepareQuery() {
	s.mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, phone_number, full_name, password, password_salt FROM users WHERE phone_number = $1	")).
		WillReturnError(fmt.Errorf("internal server error"))
	output, err := s.r.GetUserByPhoneNumber(s.ctx, s.getUserByPhoneNumberInput)
	require.Error(s.T(), err)
	require.Nil(s.T(), deep.Equal(output, User{}))
}

func (s *TestSuite) TestGetUserByPhoneNumberFailed() {
	prepare := s.mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, phone_number, full_name, password, password_salt FROM users WHERE phone_number = $1	"))
	prepare.ExpectQuery().
		WillReturnError(fmt.Errorf("internal server error")).
		WithArgs(
			s.user.PhoneNumber,
		)
	output, err := s.r.GetUserByPhoneNumber(s.ctx, s.getUserByPhoneNumberInput)
	require.Error(s.T(), err)
	require.Nil(s.T(), deep.Equal(output, User{}))
}

func (s *TestSuite) TestGetUserByPhoneNumberNotFound() {
	prepare := s.mock.ExpectPrepare(regexp.QuoteMeta("SELECT id, phone_number, full_name, password, password_salt FROM users WHERE phone_number = $1	"))
	prepare.ExpectQuery().
		WillReturnError(fmt.Errorf("sql: no rows in result set")).
		WithArgs(
			s.user.PhoneNumber,
		)
	output, err := s.r.GetUserByPhoneNumber(s.ctx, s.getUserByPhoneNumberInput)
	require.Error(s.T(), err)
	require.Nil(s.T(), deep.Equal(output, User{}))
}

func (s *TestSuite) TestGetUserByIDSuccess() {
	prepare := s.mock.ExpectPrepare(regexp.QuoteMeta("SELECT phone_number, full_name FROM users WHERE id = $1"))
	prepare.ExpectQuery().
		WillReturnRows(sqlmock.NewRows([]string{"phone_number", "full_name"}).AddRow(
			s.user.PhoneNumber,
			s.user.FullName,
		)).
		WithArgs(
			s.user.ID,
		)
	output, err := s.r.GetUserByID(s.ctx, s.getUserByIDInput)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(output, s.userInfo))
}

func (s *TestSuite) TestGetUserByIDFailedPrepareQuery() {
	s.mock.ExpectPrepare(regexp.QuoteMeta("SELECT phone_number, full_name FROM users WHERE id = $1")).
		WillReturnError(fmt.Errorf("sql: internal server error"))
	output, err := s.r.GetUserByID(s.ctx, s.getUserByIDInput)
	require.Error(s.T(), err)
	require.Nil(s.T(), deep.Equal(output, UserInfo{}))
}

func (s *TestSuite) TestGetUserByIDFailed() {
	prepare := s.mock.ExpectPrepare(regexp.QuoteMeta("SELECT phone_number, full_name FROM users WHERE id = $1"))
	prepare.ExpectQuery().
		WillReturnError(fmt.Errorf("sql: internal server error")).
		WithArgs(
			s.user.ID,
		)
	output, err := s.r.GetUserByID(s.ctx, s.getUserByIDInput)
	require.Error(s.T(), err)
	require.Nil(s.T(), deep.Equal(output, UserInfo{}))
}

func (s *TestSuite) TestGetUserByFullNameSuccess() {
	prepare := s.mock.ExpectPrepare(regexp.QuoteMeta("SELECT phone_number, full_name FROM users WHERE LOWER(full_name) = LOWER($1)"))
	prepare.ExpectQuery().
		WillReturnRows(sqlmock.NewRows([]string{"phone_number", "full_name"}).AddRow(
			s.user.PhoneNumber,
			s.user.FullName,
		)).
		WithArgs(
			s.user.FullName,
		)
	output, err := s.r.GetUserByFullName(s.ctx, s.getUserByFullName)
	require.NoError(s.T(), err)
	require.Nil(s.T(), deep.Equal(output, s.userInfo))
}

func (s *TestSuite) TestGetUserByFullNameFailedPrepareQuery() {
	s.mock.ExpectPrepare(regexp.QuoteMeta("SELECT phone_number, full_name FROM users WHERE LOWER(full_name) = LOWER($1)")).
		WillReturnError(fmt.Errorf("sql: internal server error"))
	output, err := s.r.GetUserByFullName(s.ctx, s.getUserByFullName)
	require.Error(s.T(), err)
	require.Nil(s.T(), deep.Equal(output, UserInfo{}))
}

func (s *TestSuite) TestGetUserByFullNameFailed() {
	prepare := s.mock.ExpectPrepare(regexp.QuoteMeta("SELECT phone_number, full_name FROM users WHERE LOWER(full_name) = LOWER($1)"))
	prepare.ExpectQuery().
		WillReturnError(fmt.Errorf("sql: internal server error")).
		WithArgs(
			s.user.FullName,
		)
	output, err := s.r.GetUserByFullName(s.ctx, s.getUserByFullName)
	require.Error(s.T(), err)
	require.Nil(s.T(), deep.Equal(output, UserInfo{}))
}

func (s *TestSuite) TestUpdateUserByIDSuccess() {
	prepare := s.mock.ExpectPrepare(regexp.QuoteMeta("UPDATE users SET phone_number = $1, full_name = $2 WHERE id = $3"))
	prepare.ExpectExec().
		WithArgs(
			s.user.PhoneNumber,
			s.user.FullName,
			s.user.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err := s.r.UpdateUserByID(s.ctx, s.user)
	require.NoError(s.T(), err)
}

func (s *TestSuite) TestUpdateUserByIDFailedPrepareQuery() {
	s.mock.ExpectPrepare(regexp.QuoteMeta("UPDATE users SET phone_number = $1, full_name = $2 WHERE id = $3")).
		WillReturnError(fmt.Errorf("sql: internal server error"))
	err := s.r.UpdateUserByID(s.ctx, s.user)
	require.Error(s.T(), err)
}

func (s *TestSuite) TestUpdateUserByIDFailed() {
	prepare := s.mock.ExpectPrepare(regexp.QuoteMeta("UPDATE users SET phone_number = $1, full_name = $2 WHERE id = $3"))
	prepare.ExpectExec().
		WillReturnError(fmt.Errorf("sql: internal server error")).
		WithArgs(
			s.user.PhoneNumber,
			s.user.FullName,
			s.user.ID,
		)
	err := s.r.UpdateUserByID(s.ctx, s.user)
	require.Error(s.T(), err)
}
