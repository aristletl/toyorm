package toyorm

import (
	"github.com/aristletl/toyorm/internal/errs"
	"github.com/aristletl/toyorm/internal/model"
	"strings"
)

type SQLBuilder struct {
	builder strings.Builder
	model   *model.Model
	args    []any
}

func (s *SQLBuilder) comma() {
	s.builder.WriteString(", ")
}

func (s *SQLBuilder) quota(str string) {
	s.builder.WriteString("`")
	s.builder.WriteString(str)
	s.builder.WriteString("`")
}

func (s *SQLBuilder) margin(str string) {
	if str == "" {
		return
	}
	s.builder.WriteString(" ")
	s.builder.WriteString(str)
	s.builder.WriteString(" ")
}

func (s *SQLBuilder) as(alias string) {
	if alias != "" {
		s.margin("AS")
		s.quota(alias)
	}
}

func (s *SQLBuilder) brackets(fn any) error {
	s.builder.WriteString("(")
	defer s.builder.WriteString(")")
	switch expr := fn.(type) {
	case Predicate:
		return s.buildPredicate(expr)
	case Aggregate:
		return s.buildColumn(expr.arg)
	}
	return nil
}

func (s *SQLBuilder) addArgs(vals ...any) {
	s.args = append(s.args, vals...)
}

// 创建列
func (s *SQLBuilder) buildColumn(col string) error {
	fd, ok := s.model.FieldMap[col]
	if !ok {
		return errs.NewErrUnknownField(col)
	}
	s.quota(fd.ColName)
	return nil
}

func (s *SQLBuilder) buildPredicate(e Predicate) error {
	if err := s.buildSubExpr(e.left); err != nil {
		return err
	}

	s.margin(e.op.String())

	if err := s.buildSubExpr(e.right); err != nil {
		return err
	}

	return nil
}

func (s *SQLBuilder) buildSubExpr(e Expression) error {
	if e == nil {
		return nil
	}

	switch expr := e.(type) {
	case Predicate:
		return s.brackets(expr)
	default:
		return s.buildExpression(expr)
	}
}

func (s *SQLBuilder) buildExpression(e Expression) error {
	if e == nil {
		return nil
	}

	switch expr := e.(type) {
	case Column:
		return s.buildColumn(expr.name)
	case Value:
		s.builder.WriteString("?")
		s.args = append(s.args, expr.val)
	case Predicate:
		return s.buildPredicate(expr)
	case Aggregate:
		return s.buildAggregate(expr, false)
	default:
		return errs.NewErrUnsupportedExpressionType(e)
	}
	return nil
}

func (s *SQLBuilder) buildAggregate(a Aggregate, useAlisa bool) error {
	s.builder.WriteString(a.fn)
	if err := s.brackets(a); err != nil {
		return err
	}
	if useAlisa {
		s.as(a.alias)
	}
	return nil
}

func (s *SQLBuilder) string() string {
	s.builder.WriteString(";")
	return s.builder.String()
}
