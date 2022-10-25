package toyorm

import (
	"context"

	"github.com/aristletl/toyorm/internal/errs"
	"github.com/aristletl/toyorm/internal/model"
	"github.com/aristletl/toyorm/internal/valuer"
)

type Inserter[T any] struct {
	SQLBuilder
	core
	sess        Session
	valCreator  valuer.Creator
	values      []*T
	columns     []string
	onDuplicate *Upsert
}

func NewInserter[T any](sess Session) *Inserter[T] {
	return &Inserter[T]{
		sess: sess,
		core: sess.getCore(),
		SQLBuilder: SQLBuilder{
			quoter: sess.getCore().dialect.Quoter(),
		},
		valCreator: valuer.NewUnsafeValue,
	}
}

func (i *Inserter[T]) Exec(ctx context.Context) Result {
	query, err := i.Build()
	if err != nil {
		return Result{res: nil,
			err: err,
		}
	}
	res, err := i.sess.execContext(ctx, query.SQL, query.Args...)
	return Result{
		res: res,
		err: err,
	}
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}
	var err error
	i.model, err = i.r.Get(i.values[0])
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

	if err = i.dialect.BuildOnDuplicateKey(&i.SQLBuilder, i.onDuplicate); err != nil {
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

func (i *Inserter[T]) Upsert() *UpsertBuilder[T] {
	return &UpsertBuilder[T]{
		i: i,
	}
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

func (i *Inserter[T]) buildOnDuplicateKey() error {
	if i.onDuplicate != nil {
		i.builder.WriteString(" ON DUPLICATE KEY UPDATE ")
		for idx, assign := range i.onDuplicate.assigns {
			if idx > 0 {
				i.comma()
			}
			switch expr := assign.(type) {
			case Assignment:
				if err := i.buildColumn(expr.column); err != nil {
					return err
				}
				i.builder.WriteString("=?")
				i.addArgs(expr.val)
			case Column:
				if err := i.buildColumn(expr.name); err != nil {
					return err
				}
				i.builder.WriteString("=VALUES(")
				_ = i.buildColumn(expr.name)
				i.builder.WriteString(")")
			}
		}
	}
	return nil
}

type UpsertBuilder[T any] struct {
	i *Inserter[T]
}

func (u *UpsertBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	u.i.onDuplicate = &Upsert{assigns: assigns}
	return u.i
}

type Upsert struct {
	assigns []Assignable
}
