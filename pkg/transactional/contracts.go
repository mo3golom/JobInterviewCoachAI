package transactional

import (
	"context"
	"database/sql"
)

type Template interface {
	Execute(ctx context.Context, callback func(tx Tx) error) error
}

type Tx interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
