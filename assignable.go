package toyorm

type Assignable interface {
	assign()
}

type Assignment struct {
	column string
	val    any
}

func Assign(col string, val any) Assignment {
	return Assignment{
		column: col,
		val:    val,
	}
}

func (a Assignment) assign() {}
