package querylog

import (
	"context"
	"github.com/aristletl/toyorm"
)

type MiddlewareBuilder struct {
	logFunc func(sql string, args ...any)
}

func (m *MiddlewareBuilder) LogFunc(logFunc func(sql string, args ...any)) *MiddlewareBuilder {
	m.logFunc = logFunc
	return m
}

func (m MiddlewareBuilder) Build() toyorm.Middleware {
	return func(next toyorm.Handler) toyorm.Handler {
		return func(ctx context.Context, qc *toyorm.QueryContext) *toyorm.QueryResult {
			q, err := qc.Builder.Build()
			if err != nil {
				return &toyorm.QueryResult{
					Err: err,
				}
			}

			m.logFunc(q.SQL, q.Args...)
			return next(ctx, qc)
		}
	}
}
