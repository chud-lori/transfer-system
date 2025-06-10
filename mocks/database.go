package mocks

import (
	"context"
	"database/sql"
	"transfer-system/domain/ports"

	"github.com/stretchr/testify/mock"
)

// type MockTransaction struct {
// 	mock.Mock
// }

// func (m *MockTransaction) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
// 	argsM := m.Called(ctx, query, args)
// 	return argsM.Get(0).(sql.Result), argsM.Error(1)
// }
// func (m *MockTransaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
// 	argsM := m.Called(ctx, query, args)
// 	return argsM.Get(0).(*sql.Row)
// }
// func (m *MockTransaction) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
// 	return nil, nil
// }
// func (m *MockTransaction) Commit() error {
// 	return m.Called().Error(0)
// }
// func (m *MockTransaction) Rollback() error {
// 	return m.Called().Error(0)
// }

// type MockDatabase struct {
// 	mock.Mock
// }

// func (m *MockDatabase) BeginTx(ctx context.Context) (ports.Transaction, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).(ports.Transaction), args.Error(1)
// }
// func (m *MockDatabase) Close() error {
// 	return nil
// }

// MockRow helps us mock sql.Row behavior
type MockRow struct {
	scanFunc func(dest ...interface{}) error
}

func (m *MockRow) Scan(dest ...interface{}) error {
	if m.scanFunc != nil {
		return m.scanFunc(dest...)
	}
	return nil
}

// Create a mock row that will scan specific values
func NewMockRowWithScan(scanFunc func(dest ...interface{}) error) *sql.Row {
	// We can't directly create sql.Row, so we'll work with the interface
	// This is a limitation of the current approach
	return &sql.Row{}
}

type MockTransaction struct {
	mock.Mock
}

func (m *MockTransaction) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	// Convert individual args to match the actual call signature
	callArgs := make([]interface{}, 0, len(args)+2)
	callArgs = append(callArgs, ctx, query)
	callArgs = append(callArgs, args...)

	argsM := m.Called(callArgs...)
	return argsM.Get(0).(sql.Result), argsM.Error(1)
}

func (m *MockTransaction) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	// The key fix: match the actual call signature
	// The actual call passes args as a slice, so we need to handle it correctly
	callArgs := make([]interface{}, 0, len(args)+2)
	callArgs = append(callArgs, ctx, query)
	callArgs = append(callArgs, args...)

	argsM := m.Called(callArgs...)
	return argsM.Get(0).(*sql.Row)
}

func (m *MockTransaction) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	callArgs := make([]interface{}, 0, len(args)+2)
	callArgs = append(callArgs, ctx, query)
	callArgs = append(callArgs, args...)

	argsM := m.Called(callArgs...)
	return argsM.Get(0).(*sql.Rows), argsM.Error(1)
}

func (m *MockTransaction) Commit() error {
	return m.Called().Error(0)
}

func (m *MockTransaction) Rollback() error {
	return m.Called().Error(0)
}

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) BeginTx(ctx context.Context) (ports.Transaction, error) {
	args := m.Called(ctx)
	return args.Get(0).(ports.Transaction), args.Error(1)
}

func (m *MockDatabase) Close() error {
	return m.Called().Error(0)
}

// MockResult implements sql.Result for testing
type MockResult struct {
	mock.Mock
}

func (m *MockResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockResult) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
