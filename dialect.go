package toyorm

// Dialect 方言， 构造个性部分
type Dialect interface {
	// Quoter 方言中的引号不太一样
	Quoter() byte
	BuildOnDuplicateKey(sb *SQLBuilder, odk *Upsert) error
}

// SQL 标准实现
type standardSQL struct {
}

// MySQL 方言实现
type mysqlDialect struct {
	standardSQL
}

func (m *mysqlDialect) Quoter() byte {
	return '`'
}

func (m *mysqlDialect) BuildOnDuplicateKey(b *SQLBuilder, odk *Upsert) error {
	if odk != nil {
		b.Margin("ON DUPLICATE KEY UPDATE")
		for idx, assign := range odk.assigns {
			if idx > 0 {
				b.Comma()
			}
			switch expr := assign.(type) {
			case Assignment:
				if err := b.buildAssignment(expr); err != nil {
					return err
				}
			case Column:
				if err := b.buildColumn(expr.name); err != nil {
					return err
				}
				b.builder.WriteString("=VALUES(")
				_ = b.buildColumn(expr.name)
				b.builder.WriteString(")")
			}
		}
	}
	return nil
}

// sqlite 方言实现
type sqliteDialect struct {
	standardSQL
}
