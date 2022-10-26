package toyorm

import "context"

const (
	SQLSelect  = "SELECT "
	SQLFrom    = "FROM"
	SQLWhere   = "WHERE"
	SQLGroupBy = "GROUP BY"
	SQLHaving  = "HAVING"
	SQLOrderBy = "ORDER BY"
	SQLLimit   = "LIMIT"
	SQLOffset  = "OFFSET"

	SQLInsert = "INSERT"
	SQLInto   = "INTO"
	SQLValues = "VALUES"

	SQLUpdate = "UPDATE "
	SQLSet    = "SET"
)

type Executor interface {
	Exec(ctx context.Context) Result
}

type QueryBuilder interface {
	Build() (*Query, error)
}

type Query struct {
	SQL  string
	Args []any
}
