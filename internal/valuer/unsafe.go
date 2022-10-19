package valuer

import (
	"database/sql"
	"reflect"
	"unsafe"

	"github.com/aristletl/toyorm/internal/errs"
	"github.com/aristletl/toyorm/internal/model"
)

type UnsafeValue struct {
	addr  unsafe.Pointer
	model *model.Model
}

func NewUnsafeValue(val any, m *model.Model) Value {
	return UnsafeValue{
		addr:  reflect.ValueOf(val).UnsafePointer(),
		model: m,
	}
}

func (u UnsafeValue) Field(index int) (any, error) {
	if index < 0 || index >= len(u.model.Columns) {
		return nil, errs.NewErrUnknownField("")
	}
	val := reflect.NewAt(u.model.Columns[index].Type, unsafe.Pointer(uintptr(u.addr)+u.model.Columns[index].Offset)).Elem()
	return val.Interface(), nil
}

func (u UnsafeValue) FieldByName(name string) (any, error) {
	field, ok := u.model.FieldMap[name]
	if !ok {
		return nil, errs.NewErrUnknownField(name)
	}
	val := reflect.NewAt(field.Type, unsafe.Pointer(uintptr(u.addr)+field.Offset)).Elem()
	return val.Interface(), nil
}

func (u UnsafeValue) SetColumns(rows *sql.Rows) error {
	cs, err := rows.Columns()
	if err != nil {
		return err
	}

	if len(cs) > len(u.model.ColMap) {
		return errs.ErrTooManyReturnedColumns
	}

	colValues := make([]any, len(cs))
	for i, c := range cs {
		fd, ok := u.model.ColMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}

		ptr := unsafe.Pointer(uintptr(u.addr) + fd.Offset)
		val := reflect.NewAt(fd.Type, ptr)
		colValues[i] = val.Interface()
	}
	return rows.Scan(colValues...)
}
