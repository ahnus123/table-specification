package message

import "table-specification/model"

type TableInfo struct {
	TableName string
	TableSpec map[string]([]*model.TableSpec)
	IndexSpec map[string]([]*model.IndexSpec)
}
