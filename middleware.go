package toyorm

import "context"

type QueryContext struct {
	// 用在 UPDATE, DELETE, SELECT 以及 INSERT 语句上的，
	// 用来给用户标识语句到底是上述哪个类型
	Type string
	// 可以提供给用户用于篡改 builder 本身
	Builder QueryBuilder
}

type QueryResult struct {
	Result any
	Err    error
}

type Handler func(ctx context.Context, qc *QueryContext) *QueryResult

type Middleware func(next Handler) Handler
