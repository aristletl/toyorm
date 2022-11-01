package toyorm

import "context"

type RawQuery[T any] struct {
	sql  string
	args []any
}

func NewRawQuery[T any](sql string, args ...any) *RawQuery[T] {
	return &RawQuery[T]{
		sql:  sql,
		args: args,
	}
}

func (r *RawQuery[T]) Get(ctx context.Context) (*T, error) {
	//rows, err := r.sess
	panic("implement me")
}
