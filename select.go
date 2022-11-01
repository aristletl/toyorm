package toyorm

import (
	"context"
	"errors"

	"github.com/aristletl/toyorm/internal/errs"
	"github.com/aristletl/toyorm/internal/valuer"
)

type Selector[T any] struct {
	SQLBuilder
	sess       Session
	valCreator valuer.Creator

	tableName TableReference
	where     []Predicate
	columns   []Selectable
	groupBy   []Column
	orderBy   []OrderBy
	having    []Predicate
	offset    int
	limit     int
}

// NewSelector 泛型T不支持指针
func NewSelector[T any](sess Session) *Selector[T] {
	return &Selector[T]{
		sess: sess,
		SQLBuilder: SQLBuilder{
			core: sess.getCore(),
		},
		valCreator: valuer.NewUnsafeValue,
	}
}

// Select select语句指定列名，因为这里既可以是列名又可以是聚合函数
// 因此，将参数设计为 selectable 接口
func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.columns = cols
	return s
}

// Where select语句的where
func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
	return s
}

// From select 语句的 from 指定表名
func (s *Selector[T]) From(table TableReference) *Selector[T] {
	s.tableName = table
	return s
}

func (s *Selector[T]) GroupBy(cols ...Column) *Selector[T] {
	s.groupBy = cols
	return s
}

func (s *Selector[T]) Having(ps ...Predicate) *Selector[T] {
	s.having = ps
	return s
}

func (s *Selector[T]) OrderBy(os ...OrderBy) *Selector[T] {
	s.orderBy = os
	return s
}

func (s *Selector[T]) Offset(offset int) *Selector[T] {
	s.offset = offset
	return s
}

func (s *Selector[T]) Limit(limit int) *Selector[T] {
	s.limit = limit
	return s
}

// Get 数据库查询
func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	var root Handler = func(ctx context.Context, qc *QueryContext) *QueryResult {
		q, err := qc.Builder.Build()
		if err != nil {
			return &QueryResult{Err: err}
		}
		// 使用 QueryContext，从而和 GetMulti 能够复用处理结果集的代码
		rows, err := s.sess.queryContext(ctx, q.SQL, q.Args...)
		if err != nil {
			return &QueryResult{Err: err}
		}

		if !rows.Next() {
			return &QueryResult{Err: errs.ErrNoRows}
		}

		tp := new(T)
		val := s.valCreator(tp, s.model)
		err = val.SetColumns(rows)
		return &QueryResult{
			Result: tp,
			Err:    err,
		}
	}

	for i := len(s.ms) - 1; i >= 0; i-- {
		root = s.ms[i](root)
	}

	res := root(ctx, &QueryContext{
		Type:    SQLSelect,
		Builder: s,
	})

	if res.Err != nil {
		return nil, res.Err
	}

	if t, ok := res.Result.(*T); ok {
		return t, nil
	}

	return nil, errors.New("ORM: 非正常格式")
}

func (s *Selector[T]) Build() (*Query, error) {
	var (
		err error
		t   T
	)

	s.model, err = s.r.Get(&t)
	if err != nil {
		return nil, err
	}

	s.builder.WriteString(SQLSelect)
	if err = s.buildColumns(); err != nil {
		return nil, err
	}

	s.buildFrom()

	err = s.buildWhere()
	if err != nil {
		return nil, err
	}

	err = s.buildGroupBy()
	if err != nil {
		return nil, err
	}

	err = s.buildHaving()
	if err != nil {
		return nil, err
	}

	err = s.buildOrderBy()
	if err != nil {
		return nil, err
	}

	if s.limit > 0 {
		s.Margin(SQLLimit)
		s.builder.WriteString("?")
		s.AddArgs(s.limit)
	}

	if s.offset > 0 {
		s.Margin(SQLOffset)
		s.builder.WriteString("?")
		s.AddArgs(s.offset)
	}

	return &Query{
		SQL:  s.string(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildColumns() error {
	if len(s.columns) == 0 {
		s.builder.WriteString("*")
		return nil
	}

	for i, c := range s.columns {
		if i > 0 {
			s.Comma()
		}
		switch col := c.(type) {
		case Column:
			if err := s.buildColumn(col.name); err != nil {
				return err
			}
			s.As(col.alias)
		case Aggregate:
			if err := s.buildAggregate(col, true); err != nil {
				return err
			}
		case RawExpr:
			s.builder.WriteString(col.raw)
			if len(col.args) != 0 {
				s.args = append(s.args, col.args...)
			}
		}
	}

	return nil
}

func (s *Selector[T]) buildWhere() error {
	if len(s.where) != 0 {
		s.Margin(SQLWhere)
		return s.buildPredicates(s.where)
	}
	return nil
}

func (s *Selector[T]) buildFrom() {
	s.Margin(SQLFrom)
	switch s.tableName.(type) {
	case nil:
		s.Quota(s.model.TableName)
	case Table:

	case Join:
	}
}

func (s *Selector[T]) buildGroupBy() error {
	if len(s.groupBy) != 0 {
		s.Margin(SQLGroupBy)
		for i, c := range s.groupBy {
			if i > 0 {
				s.Comma()
			}
			if err := s.buildColumn(c.name); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Selector[T]) buildHaving() error {
	if len(s.having) != 0 {
		s.Margin(SQLHaving)
		return s.buildPredicates(s.having)
	}
	return nil
}

func (s *Selector[T]) buildOrderBy() error {
	if len(s.orderBy) != 0 {
		s.Margin(SQLOrderBy)
		for i, o := range s.orderBy {
			if i > 0 {
				s.Comma()
			}
			if err := s.buildColumn(o.col); err != nil {
				return err
			}
			s.builder.WriteString(" ")
			s.builder.WriteString(o.order)
		}
	}
	return nil
}

type Selectable interface {
	selectable()
}
