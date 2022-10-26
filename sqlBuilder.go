package toyorm

import (
	"strings"

	"github.com/aristletl/toyorm/internal/errs"
	"github.com/aristletl/toyorm/internal/model"
)

type SQLBuilder struct {
	core
	builder strings.Builder
	model   *model.Model
	args    []any
}

func (s *SQLBuilder) Comma() {
	s.builder.WriteString(", ")
}

func (s *SQLBuilder) Quota(str string) {
	s.builder.WriteByte(s.dialect.Quoter())
	s.builder.WriteString(str)
	s.builder.WriteByte(s.dialect.Quoter())
}

func (s *SQLBuilder) Margin(str string) {
	if str == "" {
		return
	}
	s.builder.WriteString(" ")
	s.builder.WriteString(str)
	s.builder.WriteString(" ")
}

func (s *SQLBuilder) As(alias string) {
	if alias != "" {
		s.Margin("AS")
		s.Quota(alias)
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
	s.Quota(fd.ColName)
	return nil
}

func (s *SQLBuilder) buildPredicates(pres []Predicate) error {
	pred := pres[0]
	for i := 1; i < len(pres); i++ {
		pred = pred.AND(pres[i])
	}
	return s.buildExpression(pred)
}

func (s *SQLBuilder) buildPredicate(e Predicate) error {
	if err := s.buildSubExpr(e.left); err != nil {
		return err
	}

	s.Margin(e.op.String())

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

func (s *SQLBuilder) buildExpression(e any) error {
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
	case RawExpr:
		return s.buildRawExpr(expr)
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
		s.As(a.alias)
	}
	return nil
}

func (s *SQLBuilder) buildAssignment(assign Assignment) error {
	if err := s.buildColumn(assign.column); err != nil {
		return err
	}
	s.builder.WriteString("=")
	return s.buildExpression(assign.val)
}

func (s *SQLBuilder) buildRawExpr(raw RawExpr) error {
	s.builder.WriteString(raw.raw)
	if len(raw.args) != 0 {
		s.addArgs(raw.args...)
	}
	return nil
}

func (s *SQLBuilder) string() string {
	s.builder.WriteString(";")
	return s.builder.String()
}
