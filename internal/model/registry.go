package model

import (
	"github.com/aristletl/toyorm/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

type Option func(r *Registry) error

type Underscore func(name string) string

type Registry struct {
	UnderscoreName Underscore
	models         sync.Map
}

func NewRegistry(opts ...Option) (*Registry, error) {
	r := &Registry{
		UnderscoreName: underscoreName,
	}

	for _, opt := range opts {
		err := opt(r)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

// Get 查找元数据模型
func (r *Registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	return r.Register(val)
}

func (r *Registry) Register(val any) (*Model, error) {
	m, err := r.parseModel(val)
	if err != nil {
		return nil, err
	}
	typ := reflect.TypeOf(val)
	r.models.Store(typ, m)
	return m, nil
}

// parseModel 支持从标签中提取自定义设置
// 标签形式 orm:"key1=value1,key2=value2"
func (r *Registry) parseModel(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()

	// 获得字段的数量
	numField := typ.NumField()
	fieldMap := make(map[string]*Field, numField)
	colMap := make(map[string]*Field, numField)
	cols := make([]*Field, numField)
	for i := 0; i < numField; i++ {
		fdType := typ.Field(i)
		colName := r.UnderscoreName(fdType.Name)
		f := &Field{
			Index:   i,
			ColName: colName,
			Type:    fdType.Type,
			Offset:  fdType.Offset,
		}
		fieldMap[fdType.Name] = f
		colMap[colName] = f
		cols[i] = f
	}

	return &Model{
		TableName: r.UnderscoreName(typ.Name()),
		Columns:   cols,
		FieldMap:  fieldMap,
		ColMap:    colMap,
	}, nil
}

// underscoreName 驼峰转字符串命名
func underscoreName(name string) string {
	var builder strings.Builder
	for i, v := range name {
		if unicode.IsUpper(v) {
			if i != 0 {
				builder.WriteString("_")
			}
			builder.WriteRune(unicode.ToLower(v))
		} else {
			builder.WriteRune(v)
		}
	}
	return builder.String()
}

func WithUnderscoreName(fn Underscore) Option {
	return func(r *Registry) error {
		r.UnderscoreName = fn
		return nil
	}
}
