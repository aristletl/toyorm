package toyorm

// Expression 一个抽象接口，无实际意义
// 主要是解决 Predicate 表达式的问题，Predicate
// 操作符两侧，既可以是and语句连接的两个类似于 `id = 1` 的表达式
// 也有可能是单纯的比较语句 `id > 2` ，两侧分别是 `columns` 与 `value`
// 因此需要一个抽象层将两边统一起来，所以出现了 `Expression`
type Expression interface {
	Expr()
}

type RawExpr struct {
	raw  string
	args []any
}

func (r RawExpr) selectable() {}

func Raw(expr string, args ...any) RawExpr {
	return RawExpr{
		raw:  expr,
		args: args,
	}
}
