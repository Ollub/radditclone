package repo

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang-stepik-2022q1/reditclone/pkg/users"
	"testing"
)

type Suite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock

	repo *RepoSql
}

func (s *Suite) SetupSuite() {
	var err error
	s.db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)
	s.repo = NewSql(s.db)
}

func (s *Suite) TearDownSuite() {
	s.db.Close()
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

type MockRows struct {
	columns []string
	rows    [][]driver.Value
}
type MockOpts struct {
	query string
	args  []driver.Value
	rows  *MockRows
	err   error
}

func (s *Suite) SetupMock(opts MockOpts) {
	query := s.mock.ExpectQuery(opts.query)
	if len(opts.args) != 0 {
		query.WithArgs(opts.args...)
	}
	if opts.rows != nil {
		rows := sqlmock.NewRows(opts.rows.columns)
		for _, row := range opts.rows.rows {
			rows.AddRow(row...)
		}
		query.WillReturnRows(rows)
	}
	if opts.err != nil {
		query.WillReturnError(opts.err)
	}
}

func (s *Suite) TestGetByName() {
	var name = "John"
	var dbErr = errors.New("Some db error")
	john := &users.User{1, name, "hashedPass"}

	for _, tt := range [...]struct {
		opts          MockOpts
		expected      *users.User
		expectedError error
	}{
		// single users returned
		{
			opts: MockOpts{
				query: "SELECT id, name, pass_hash FROM users",
				args:  []driver.Value{name},
				rows: &MockRows{
					[]string{"id", "name", "pass_hash"},
					[][]driver.Value{
						{john.Id, john.Name, john.PassHash},
					},
				},
			},
			expected: john,
		},
		// no rows not lead to repo error
		{
			opts: MockOpts{
				query: "SELECT id, name, pass_hash FROM users",
				args:  []driver.Value{name},
				err:   sql.ErrNoRows,
			},
		},
		// unexpected error proxied
		{
			opts: MockOpts{
				query: "SELECT id, name, pass_hash FROM users",
				args:  []driver.Value{name},
				err:   dbErr,
			},
			expectedError: dbErr,
		},
	} {
		s.SetupMock(tt.opts)
		item, err := s.repo.GetByName(name)
		s.Equal(tt.expectedError, err)
		s.Equal(tt.expected, item)
	}
}

func (s *Suite) TestAdd() {
	var dbErr = errors.New("Some db error")
	var insertId int64 = 123
	johnIn := &users.User{Name: "John", PassHash: "hashedPass"}

	for _, tt := range [...]struct {
		name          string
		dbErr         error
		expected      int64
		expectedError error
	}{
		{"OK", nil, insertId, nil},
		{"Unexpected error", dbErr, 0, dbErr},
	} {
		s.mock.
			ExpectQuery("INSERT INTO users").
			WithArgs(johnIn.Name, johnIn.PassHash).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(insertId)).
			WillReturnError(tt.expectedError)

		s.Run(tt.name, func() {
			item, err := s.repo.Add(johnIn)
			s.Equal(tt.expectedError, err)
			s.Equal(tt.expected, item)
		})
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

////////////////////////
// Splitted suite tests
// I left them just to compare complexity and readability
/////////////////////////

//func (s *Suite) TestGetByNameOk() {
//	var name = "John"
//	users := &users.UserId{"1", name, "hashedPass"}
//
//	s.SetupMock(MockOpts{
//		query: "SELECT id, name, pass_hash FROM users",
//		args:  []driver.Value{name},
//		rows: &MockRows{
//			[]string{"id", "name", "pass_hash"},
//			[][]driver.Value{
//				{users.Id, users.Name, users.PassHash},
//			},
//		},
//	})
//
//	item, err := s.repo.GetByName(name)
//	s.NoError(err)
//	s.Equal(item, users)
//}
//
//func (s *Suite) TestGetByNameEmptyResult() {
//	var name = "John"
//
//	s.SetupMock(MockOpts{
//		query: "SELECT id, name, pass_hash FROM users",
//		args:  []driver.Value{name},
//		err:   sql.ErrNoRows,
//	})
//
//	item, err := s.repo.GetByName(name)
//	s.Nil(err)
//	s.Nil(item)
//}
//
//func (s *Suite) TestGetByNameUnexpectedError() {
//	var name = "John"
//	var dbErr = errors.New("Some db error")
//
//	s.SetupMock(MockOpts{
//		query: "SELECT id, name, pass_hash FROM users",
//		args:  []driver.Value{name},
//		err:   dbErr,
//	})
//
//	item, err := s.repo.GetByName(name)
//	s.Equal(err, dbErr)
//	s.Nil(item)
//}

//////////////////
// Standard lib based tests
/////////////////

//
//func TestGetByName(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("cant create mock: %s", err)
//	}
//	defer db.Close()
//
//	repo := NewSql(db)
//
//	var name string = "John"
//
//	// good query
//	rows := sqlmock.
//		NewRows([]string{"id", "name", "pass_hash"})
//	expect := []*users.UserId{
//		{"1", name, "hashedPass"},
//	}
//	for _, item := range expect {
//		rows = rows.AddRow(item.Id, item.Name, item.PassHash)
//	}
//
//	mock.
//		ExpectQuery("SELECT id, name, pass_hash FROM users").
//		WithArgs(name).
//		WillReturnRows(rows)
//
//	item, err := repo.GetByName(name)
//	if err != nil {
//		t.Errorf("unexpected err: %s", err)
//		return
//	}
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//		return
//	}
//	if !reflect.DeepEqual(item, expect[0]) {
//		t.Errorf("results not match, want %v, have %v", expect[0], item)
//		return
//	}
//}
//
//// sql.ErrNoRows should not cause an error
//// empty result with nil error should be returned
//func TestGetByNameEmptyResult(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("cant create mock: %s", err)
//	}
//	defer db.Close()
//
//	repo := NewSql(db)
//
//	var name string = "John"
//
//	mock.
//		ExpectQuery("SELECT id, name, pass_hash FROM users").
//		WithArgs(name).
//		WillReturnError(sql.ErrNoRows)
//
//	item, err := repo.GetByName(name)
//	if err != nil {
//		t.Errorf("unexpected err: %s", err)
//		return
//	}
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//		return
//	}
//	if item != nil {
//		t.Errorf("Expected empty result, got %v", item)
//		return
//	}
//}
//
//// should return error if db raise error
//func TestGetByNameUnexpectedError(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("cant create mock: %s", err)
//	}
//	defer db.Close()
//
//	repo := NewSql(db)
//
//	var (
//		name  = "John"
//		dbErr = errors.New("Some db error")
//	)
//
//	mock.
//		ExpectQuery("SELECT id, name, pass_hash FROM users").
//		WithArgs(name).
//		WillReturnError(dbErr)
//
//	item, err := repo.GetByName(name)
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//		return
//	}
//	if item != nil {
//		t.Errorf("Expected empty result, got %v", item)
//		return
//	}
//	if err != dbErr {
//		t.Errorf("Wrong error returned. Expected: %s. Got: %s", dbErr, err)
//		return
//	}
//}
