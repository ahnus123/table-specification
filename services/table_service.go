package services

import (
	"table-specification/message"
	"table-specification/processor"

	"github.com/jinzhu/gorm"
)

// GetTableSpecs : 스키마별 테이블 명세 조회
func GetTableSpecs(db *gorm.DB, schemaList []string) (map[string][]*message.Spec, error) {
	specList := map[string][]*message.Spec{}

	tableList, err := processor.GetTableListBySchema(db, schemaList)
	if err != nil {
		return nil, err
	}

	if len(tableList) > 0 {
		for schema, tables := range tableList {

			specs := []*message.Spec{}
			if len(tables) > 0 {
				for _, table := range tables {
					tSpec, iSpec, err := processor.GetTableSpec(db, schema, table)
					if err != nil {
						return nil, err
					}
					specs = append(specs, &message.Spec{TableSpec: tSpec, IndexSpec: iSpec})
				}
			}
			specList[schema] = specs
		}
	}

	return specList, nil
}
