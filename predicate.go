package toyorm

type Op string

func (o Op) String() string {
	return string(o)
}

const (
	opEQ = "="
	opLT = "<"
	opGT = ">"

	opNOT = "NOT"
	opAND = "AND"
	opOR  = "OR"
)

// Predicate 表达式
type Predicate struct {
	left  Expression
	op    Op
	right Expression
}

func (p Predicate) Expr() {}

type Value struct {
	val any
}

func (v Value) Expr() {}

func (p Predicate) AND(p1 Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opAND,
		right: p1,
	}
}

func (p Predicate) OR(p1 Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opOR,
		right: p1,
	}
}

func NOT(p1 Predicate) Predicate {
	return Predicate{
		op:    opNOT,
		right: p1,
	}
}
