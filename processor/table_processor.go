package processor

import (
	"log"
	"strings"
	"table-specification/model"

	"github.com/jinzhu/gorm"
)

// GetTableListBySchema : 스키마별 테이블명 조회
func GetTableListBySchema(db *gorm.DB, schemas []string) (map[string][]string, error) {
	tableList := make(map[string][]string)

	if len(schemas) > 0 {
		for _, schema := range schemas {
			tables := []string{}
			rows, err := db.Table("tables").Select("table_name").Where("table_schema = ?", schema).Rows()
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			for rows.Next() {
				var table string
				rows.Scan(&table)
				tables = append(tables, table)
			}
			tableList[schema] = tables
		}
	}

	return tableList, nil
}

// GetTableSpec : 테이블별 명세 조회(컬럼, 인덱스)
func GetTableSpec(db *gorm.DB, tableSchema, tableName string) (map[string]([]*model.TableSpec), map[string]([]*model.IndexSpec), error) {

	tableSpec := make(map[string]([]*model.TableSpec))
	indexSpec := make(map[string]([]*model.IndexSpec))

	tSpec, err := getTableSpecByTableName(db, tableSchema, tableName)
	if err != nil {
		log.Fatalln(err)
		return nil, nil, err
	}
	tableSpec[tableName] = tSpec

	iSpec, err := getTableIndexByTableName(db, tableSchema, tableName)
	if err != nil {
		log.Fatalln(err)
		return nil, nil, err
	}
	indexSpec[tableName] = iSpec

	return tableSpec, indexSpec, nil
}

// getTableSpecByTableName : 테이블명으로 테이블 정보 조회
func getTableSpecByTableName(db *gorm.DB, tableSchema, tableName string) ([]*model.TableSpec, error) {
	tableInfo := []*model.TableSpec{}

	query := db.Table("tables as tb").
		Joins("INNER JOIN columns col ON tb.table_schema = col.table_schema AND tb.table_name = col.table_name").
		Select(`tb.table_name
				, tb.table_comment
				, col.column_name
				, col.column_comment
				, col.column_type
				, case when col.is_nullable='YES' then '' when col.IS_NULLABLE='NO' then 'NotNull' else '' end as is_nullable
				, case when col.column_key='PRI' then 'PK' else '' end as column_key
				, col.column_default`).
		Where("tb.table_schema = ? AND tb.table_name = ?", tableSchema, tableName)
	err := query.Scan(&tableInfo).Error
	if err != nil {
		return nil, err
	}

	// 타입, 길이 설정
	for _, ti := range tableInfo {
		if ti.ColumnType != "" {
			cType, after, found := strings.Cut(ti.ColumnType, "(")
			if found {
				cLength, _, _ := strings.Cut(after, ")")
				ti.ColumnLength = cLength
			}
			ti.ColumnType = cType
		}
	}

	return tableInfo, nil
}

// getTableIndexByTableName : 테이블명으로 인덱스 조회
func getTableIndexByTableName(db *gorm.DB, tableSchema, tableName string) ([]*model.IndexSpec, error) {
	indexSpec := []*model.IndexSpec{}

	query := db.Table("statistics").
		Select(`index_name
			, GROUP_CONCAT(column_name ORDER BY seq_in_index SEPARATOR ', ') AS 'columns'
			, IF(max(non_unique)=1, 'No', 'Yes') AS 'unique'`).
		Where("table_schema = ?", tableSchema).
		Where("table_name = ?", tableName).
		Group("index_name").
		Group("if(index_name='PRIMARY', '1', index_name)")
	err := query.Scan(&indexSpec).Error
	if err != nil {
		return nil, err
	}

	return indexSpec, nil
}
