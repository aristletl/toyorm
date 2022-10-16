package valuer

import (
	"database/sql"
	"github.com/aristletl/toyorm/internal/model"
)

// Value 是对结构体实例的内部抽象
type Value interface {
	// SetColumns 设置新值
	SetColumns(rows *sql.Rows) error
}

type Creator func(val interface{}, m *model.Model) Value
