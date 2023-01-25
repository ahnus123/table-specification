package message

import "table-specification/model"

type Spec struct {
	TableSpec map[string]([]*model.TableSpec)
	IndexSpec map[string]([]*model.IndexSpec)
}
