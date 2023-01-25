package model

type Tables struct {
}

type TableSpec struct {
	TableName     string `gorm:"column:table_name"`
	TableComment  string `gorm:"column:table_comment"`
	ColumnName    string `gorm:"column:column_name"`
	ColumnComment string `gorm:"column:column_comment"`
	ColumnType    string `gorm:"column:column_type"`
	ColumnLength  string `gorm:"-"`
	IsNullable    string `gorm:"column:is_nullable"`
	ColumnKey     string `gorm:"column:column_key"`
	ColumnDefault string `gorm:"column:column_default"`
}

type IndexSpec struct {
	IndexName string `gorm:"column:index_name"`
	Columns   string `gorm:"column:columns"`
	Unique    string `gorm:"column:unique"`
}
