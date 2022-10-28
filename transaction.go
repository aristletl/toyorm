package toyorm

import (
	"context"
	"database/sql"

	"github.com/aristletl/toyorm/internal/model"
)

type Tx struct {
	core
	tx *sql.Tx
}

func (t *Tx) getCore() core {
	return t.core
}

func (t *Tx) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *Tx) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func (t *Tx) RollBack() error {
	return t.tx.Rollback()
}

func (t *Tx) RollbackIfNotCommit() error {
	err := t.RollBack()
	if err == sql.ErrTxDone {
		return err
	}
	return nil
}

type Session interface {
	getCore() core
	queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	execContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type core struct {
	r       *model.Registry
	ms      []Middleware
	dialect Dialect
}
