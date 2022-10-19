package toyorm

import (
	"github.com/aristletl/toyorm/internal/errs"
	"github.com/aristletl/toyorm/internal/model"
	"github.com/aristletl/toyorm/internal/valuer"
)

type Inserter[T any] struct {
	SQLBuilder
	db         *DB
	valCreator valuer.Creator
	values     []*T
	columns    []string
}

func NewInserter[T any](db *DB) *Inserter[T] {
	return &Inserter[T]{
		db:         db,
		valCreator: valuer.NewUnsafeValue,
	}
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}
	var err error
	i.model, err = i.db.r.Get(i.values[0])
	if err != nil {
		return nil, err
	}

	i.builder.WriteString(SQLInsert)
	i.margin(SQLInto)
	i.quota(i.model.TableName)

	if err = i.buildColumns(); err != nil {
		return nil, err
	}

	if err = i.buildValues(); err != nil {
		return nil, err
	}

	return &Query{
		SQL:  i.string(),
		Args: i.args,
	}, nil
}

func (i *Inserter[T]) Columns(cols ...string) *Inserter[T] {
	i.columns = cols
	return i
}

// Values 指定 INSERT INTO XXX VALUES 的 VALUES 部分
func (i *Inserter[T]) Values(vals ...*T) *Inserter[T] {
	i.values = vals
	return i
}

func (i *Inserter[T]) buildColumns() error {
	i.builder.WriteString("(")
	defer i.builder.WriteString(")")
	if len(i.columns) == 0 {
		for idx, c := range i.model.Columns {
			if idx > 0 {
				i.comma()
			}
			i.quota(c.ColName)
		}
	} else {
		for idx, c := range i.columns {
			if idx > 0 {
				i.comma()
			}
			if err := i.buildColumn(c); err != nil {
				return err
			}
		}
	}
	return nil
}

func (i *Inserter[T]) buildValues() error {
	fields := i.model.Columns
	if len(i.columns) != 0 {
		fields = make([]*model.Field, 0, len(i.columns))
		for _, colName := range i.columns {
			fd, ok := i.model.FieldMap[colName]
			if !ok {
				return errs.NewErrUnknownField(colName)
			}
			fields = append(fields, fd)
		}
	}
	i.builder.WriteString(" VALUES")
	for j := 0; j < len(i.values); j++ {
		if j > 0 {
			i.comma()
		}
		val := i.valCreator(i.values[j], i.model)
		i.builder.WriteString("(")
		for k, meta := range fields {
			if k > 0 {
				i.comma()
			}
			i.builder.WriteString("?")
			fdVal, err := val.Field(meta.Index)
			if err != nil {
				return err
			}
			i.addArgs(fdVal)
		}
		i.builder.WriteString(")")
	}
	return nil
}
