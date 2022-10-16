// Package model ORM框架需要解析模型以获得模型的元数据，这
// 些元数据将被用于构建 SQL、执行校验，以及用于处理结果集。
package model

import "reflect"

// Model 用于定义数据到数据库表的映射关系
type Model struct {
	TableName string
	// 当insert遍历属性的时候，map遍历的返回结果会变得无序，因此需要一个 slice 切片保证属性的有序性
	Columns  []*Field
	FieldMap map[string]*Field
	ColMap   map[string]*Field
}

// Field field字段
type Field struct {
	Index   int
	ColName string
	Type    reflect.Type
	Offset  uintptr
}
