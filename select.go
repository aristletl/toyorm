package toyorm

import (
	"context"

	"github.com/aristletl/toyorm/internal/errs"
	"github.com/aristletl/toyorm/internal/valuer"
)

type Selector[T any] struct {
	SQLBuilder
	core
	sess       Session
	valCreator valuer.Creator

	tableName string
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
		core: sess.getCore(),
		SQLBuilder: SQLBuilder{
			quoter: sess.getCore().dialect.Quoter(),
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
func (s *Selector[T]) From(table string) *Selector[T] {
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
	q, err := s.Build()
	if err != nil {
		return nil, err
	}
	// s.db 是我们定义的 DB
	// s.db.db 则是 sql.DB
	// 使用 QueryContext，从而和 GetMulti 能够复用处理结果集的代码
	rows, err := s.sess.queryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, errs.ErrNoRows
	}

	tp := new(T)
	//meta, err := s.db.r.Get(tp)
	//if err != nil {
	//	return nil, err
	//}
	val := s.valCreator(tp, s.model)
	err = val.SetColumns(rows)
	return tp, err
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
		s.builder.WriteString(" LIMIT ?")
		s.addArgs(s.limit)
	}

	if s.offset > 0 {
		s.builder.WriteString(" OFFSET ?")
		s.addArgs(s.offset)
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
			s.comma()
		}
		switch col := c.(type) {
		case Column:
			if err := s.buildColumn(col.name); err != nil {
				return err
			}
			s.as(col.alias)
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

func (s *Selector[T]) buildPredicates(pres []Predicate) error {
	pred := pres[0]
	for i := 1; i < len(pres); i++ {
		pred = pred.AND(pres[i])
	}
	return s.buildExpression(pred)
}

func (s *Selector[T]) buildWhere() error {
	if len(s.where) != 0 {
		s.margin(SQLWhere)
		return s.buildPredicates(s.where)
	}
	return nil
}

func (s *Selector[T]) buildFrom() {
	s.margin(SQLFrom)
	if s.tableName == "" {
		s.tableName = s.model.TableName
		s.quota(s.tableName)
	} else {
		s.tableName = s.r.UnderscoreName(s.tableName)
		s.builder.WriteString(s.tableName)
	}
}

func (s *Selector[T]) buildGroupBy() error {
	if len(s.groupBy) != 0 {
		s.margin(SQLGroupBy)
		for i, c := range s.groupBy {
			if i > 0 {
				s.comma()
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
		s.margin(SQLHaving)
		return s.buildPredicates(s.having)
	}
	return nil
}

func (s *Selector[T]) buildOrderBy() error {
	if len(s.orderBy) != 0 {
		s.margin(SQLOrderBy)
		for i, o := range s.orderBy {
			if i > 0 {
				s.comma()
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
