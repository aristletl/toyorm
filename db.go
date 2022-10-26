package toyorm

import (
	"context"
	"database/sql"

	"go.uber.org/multierr"

	"github.com/aristletl/toyorm/internal/model"
)

type DBOption func(*DB)

type DB struct {
	core
	db *sql.DB
}

func (d *DB) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

func (d *DB) getCore() core {
	return d.core
}

func (d *DB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

func Open(driver string, dns string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dns)
	if err != nil {
		return nil, err
	}

	return OpenDB(db, opts...)
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	r, err := model.NewRegistry()
	if err != nil {
		return nil, err
	}

	res := &DB{
		core: core{
			r:       r,
			dialect: &mysqlDialect{},
		},
		db: db,
	}

	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func (d *DB) Begin(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx}, nil
}

func (d *DB) DoTx(ctx context.Context, opts *sql.TxOptions, task func(ctx context.Context, tx *Tx) error) (err error) {
	tx, err := d.Begin(ctx, opts)
	if err != nil {
		return err
	}

	panicked := true
	defer func() {
		if panicked {
			e := tx.RollBack()
			err = multierr.Combine(err, e)
		} else {
			err = multierr.Combine(err, tx.Commit())
		}
	}()

	err = task(ctx, tx)
	panicked = false
	return
}

func DBWithRegistry(r *model.Registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}

func DBWithDialect(d Dialect) DBOption {
	return func(db *DB) {
		db.dialect = d
	}
}
