package toyorm

import (
	"database/sql"
	"github.com/aristletl/toyorm/internal/model"
)

type DBOption func(*DB)

type DB struct {
	r  *model.Registry
	db *sql.DB
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
		r:  r,
		db: db,
	}

	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func DBWithRegistry(r *model.Registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}
