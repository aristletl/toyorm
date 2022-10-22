package toyorm

type Column struct {
	name  string
	alias string
}

func (c Column) assign() {}

func (c Column) selectable() {}

func (c Column) Expr() {}

func Col(name string) Column {
	return Column{name: name}
}

func (c Column) EQ(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opEQ,
		right: Value{val: val},
	}
}

func (c Column) GT(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opGT,
		right: Value{val: val},
	}
}

func (c Column) LT(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opLT,
		right: Value{val: val},
	}
}

func (c Column) AS(alias string) Column {
	return Column{
		name:  c.name,
		alias: alias,
	}
}

type OrderBy struct {
	col   string
	order string
}

func Asc(col string) OrderBy {
	return OrderBy{
		col:   col,
		order: "ASC",
	}
}

func Desc(col string) OrderBy {
	return OrderBy{
		col:   col,
		order: "DESC",
	}
}
