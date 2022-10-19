package valuer

import (
	"database/sql"
	"reflect"

	"github.com/aristletl/toyorm/internal/errs"
	"github.com/aristletl/toyorm/internal/model"
)

type ReflectValue struct {
	val   reflect.Value
	model *model.Model
}

func (r ReflectValue) Field(index int) (any, error) {
	//TODO implement me
	panic("implement me")
}

func (r ReflectValue) FieldByName(name string) (any, error) {
	fd, ok := r.model.FieldMap[name]
	if !ok {
		return nil, errs.NewErrUnknownField(name)
	}
	res := r.val.Field(fd.Index)
	return res.Interface(), nil
}

func (r ReflectValue) SetColumns(rows *sql.Rows) error {
	cs, err := rows.Columns()
	if err != nil {
		return err
	}

	if len(cs) > len(r.model.ColMap) {
		return errs.ErrTooManyReturnedColumns
	}

	colValues := make([]any, len(cs))
	colEleValues := make([]reflect.Value, len(cs))
	for i, c := range cs {
		fd, ok := r.model.ColMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		val := reflect.New(fd.Type)
		colValues[i] = val.Interface()
		colEleValues[i] = val.Elem()
	}
	if err = rows.Scan(colValues...); err != nil {
		return err
	}

	for i, c := range cs {
		cm := r.model.ColMap[c]
		fd := r.val.Field(cm.Index)
		fd.Set(colEleValues[i])
	}
	return nil
}

func NewReflectValue(val any, m *model.Model) Value {
	return ReflectValue{
		val:   reflect.ValueOf(val).Elem(),
		model: m,
	}
}
