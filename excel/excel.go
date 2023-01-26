package excel

import (
	"fmt"
	"table-specification/message"
	"time"

	"github.com/xuri/excelize/v2"
)

func ExportExcel(specList map[string][]*message.TableInfo) error {
	var err error

	// 파일 생성
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("close file ERROR : ", err.Error())
		}
	}()

	// TODO : 표지 생성
	// TODO : 목차 생성

	// 스키마별로 테이블 시트 생성
	for schema, tableList := range specList {
		if len(tableList) > 0 {
			for _, ts := range tableList {
				tableName := ts.TableName
				tSpec := ts.TableSpec[tableName]
				iSpec := ts.IndexSpec[tableName]

				err = CreateNewTableSheet(file, schema, ts.TableName, tSpec, iSpec)
				if err != nil {
					return err
				}
			}
		}
	}

	// 파일 저장
	date := time.Now().UTC().Local()
	fileName := "TEST" + fmt.Sprintf("%d", date.Hour()) + fmt.Sprintf("%d", date.Minute()) + fmt.Sprintf("%d", date.Second())

	err = file.SaveAs(fileName + ".xlsx")
	if err != nil {
		return err
	}

	return nil
}
