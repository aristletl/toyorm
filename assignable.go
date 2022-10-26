package toyorm

// Assignable 用于指派某个具体列的一个抽象接口，无实际意义
type Assignable interface {
	assign()
}

type Assignment struct {
	column string
	val    any
}

func Assign(col string, val any) Assignment {
	if _, ok := val.(Expression); !ok {
		val = Value{val: val}
	}
	return Assignment{
		column: col,
		val:    val,
	}
}

func (a Assignment) assign() {}
