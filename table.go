package toyorm

type TableReference interface {
	tableAlias() string
}

// Table 普通表
type Table struct {
	entity any
	alias  string
}

func (t Table) tableAlias() string {
	return t.alias
}

func TableOf(entity any) Table {
	return Table{}
}

func (t Table) Join() *JoinBuilder {
	return &JoinBuilder{
		left:  t,
		typ:   "JOIN",
		right: Join{},
	}
}

// Join Join查询
type Join struct {
}

type JoinBuilder struct {
	left  TableReference
	typ   string
	right TableReference
}

func (j Join) tableAlias() string {
	return ""
}

// SubQuery 子查询
type SubQuery struct {
}

func (s SubQuery) tableAlias() string {
	//TODO implement me
	panic("implement me")
}
