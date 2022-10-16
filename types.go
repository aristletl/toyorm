package toyorm

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
)

type QueryBuilder interface {
	Build() (*Query, error)
}

type Query struct {
	SQL  string
	Args []any
}
