package toyorm

import (
	"context"
	"github.com/aristletl/toyorm/internal/errs"
	"github.com/aristletl/toyorm/internal/valuer"
)

type Updater[T any] struct {
	SQLBuilder
	sess       Session
	valCreator valuer.Creator
	val        *T
	assigns    []Assignable
	where      []Predicate
}

func (u *Updater[T]) Exec(ctx context.Context) Result {
	q, err := u.Build()
	if err != nil {
		return Result{
			err: err,
		}
	}
	res := Result{}
	res.res, res.err = u.sess.execContext(ctx, q.SQL, q.Args...)
	return res
}

func NewUpdater[T any](sess Session) *Updater[T] {
	c := sess.getCore()
	return &Updater[T]{
		SQLBuilder: SQLBuilder{
			core: c,
		},
		valCreator: valuer.NewUnsafeValue,
	}
}

func (u *Updater[T]) Build() (*Query, error) {
	if len(u.assigns) == 0 {
		return nil, errs.ErrNoUpdatedColumns
	}

	var err error
	u.builder.WriteString(SQLUpdate)
	u.model, err = u.r.Get(u.val)
	if err != nil {
		return nil, err
	}
	u.Quota(u.model.TableName)

	u.Margin(SQLSet)
	val := u.valCreator(u.val, u.model)
	for i, assign := range u.assigns {
		if i > 0 {
			u.Comma()
		}
		switch expr := assign.(type) {
		case Column:
			arg, err := val.FieldByName(expr.name)
			if err != nil {
				return nil, err
			}
			if err = u.buildColumn(expr.name); err != nil {
				return nil, err
			}
			u.builder.WriteString("=?")
			u.AddArgs(arg)
		case Assignment:
			if err = u.buildAssignment(expr); err != nil {
				return nil, err
			}
		}
	}

	if len(u.where) != 0 {
		u.Margin(SQLWhere)
		if err = u.buildPredicates(u.where); err != nil {
			return nil, err
		}
	}

	return &Query{
		SQL:  u.string(),
		Args: u.args,
	}, nil
}

func (u *Updater[T]) Update(t *T) *Updater[T] {
	u.val = t
	return u
}

func (u *Updater[T]) Set(assigns ...Assignable) *Updater[T] {
	u.assigns = assigns
	return u
}

func (u *Updater[T]) Where(ps ...Predicate) *Updater[T] {
	u.where = ps
	return u
}
