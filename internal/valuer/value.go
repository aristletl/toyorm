package valuer

import (
	"database/sql"

	"github.com/aristletl/toyorm/internal/model"
)

// Value 是对结构体实例的内部抽象
type Value interface {
	FieldByName(name string) (any, error)
	Field(index int) (any, error)
	// SetColumns 设置新值
	SetColumns(rows *sql.Rows) error
}

// Creator 用于创建
type Creator func(val interface{}, m *model.Model) Value
