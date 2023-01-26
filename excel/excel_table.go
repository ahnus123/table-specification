package excel

import (
	"errors"
	"fmt"
	"table-specification/model"
	"time"

	"github.com/xuri/excelize/v2"
)

var (
	titleColNum  = 1
	infoColNum1  = 2
	infoColNum2  = 3
	columnColNum = 4
	indexColNum  = -1
	footerColNum = -1
	finalColNum  = -1
)

// CreateNewTableSheet : 새 테이블 시트 생성
func CreateNewTableSheet(file *excelize.File, schema, tableName string, tableSpec []*model.TableSpec, indexSpec []*model.IndexSpec) error {
	if len(tableSpec) < 1 {
		return errors.New("there is no table spec")
	}

	tableComment := tableSpec[0].TableComment
	sheetName := tableName
	idx, err := file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// 테이블 명세 데이터 셋팅
	setTableSheetLayout(file, schema, sheetName, tableName, tableComment)
	cEndIdx, err := setTableColumn(file, sheetName, tableSpec, columnColNum+1)
	// _, err = SetSheetColumn(file, sheetName, tableSpec, 6)
	if err != nil {
		return err
	}
	indexColNum = cEndIdx
	iEndIdx, err := setTableIndex(file, sheetName, indexSpec, cEndIdx)
	if err != nil {
		return err
	}
	footerColNum = iEndIdx
	fEndIdx, err := setTableSheetFooter(file, sheetName, iEndIdx)
	if err != nil {
		return err
	}
	finalColNum = fEndIdx
	file.SetActiveSheet(idx)

	// 스타일 설정
	commonStyle := &excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Font:      &excelize.Font{Size: 10},
	}
	setTableSheetStyle(file, sheetName, *commonStyle, finalColNum)
	setTitleCellStyle(file, sheetName, *commonStyle, indexColNum, finalColNum)

	return nil
}

// setTableColumn : 테이블의 컬럼 설정
func setTableColumn(file *excelize.File, sheetName string, tableSpec []*model.TableSpec, startIdx int) (int, error) {
	for i, ts := range tableSpec {
		file.SetCellValue(sheetName, fmt.Sprint("A", (startIdx+i)), (i + 1))
		file.SetCellValue(sheetName, fmt.Sprint("B", (startIdx+i)), ts.ColumnName)    // column name
		file.SetCellValue(sheetName, fmt.Sprint("C", (startIdx+i)), ts.ColumnComment) // column comment
		file.SetCellValue(sheetName, fmt.Sprint("D", (startIdx+i)), ts.ColumnType)    // column type
		file.SetCellValue(sheetName, fmt.Sprint("E", (startIdx+i)), ts.ColumnLength)  // column length
		file.SetCellValue(sheetName, fmt.Sprint("F", (startIdx+i)), ts.IsNullable)    // column is nullable
		file.SetCellValue(sheetName, fmt.Sprint("G", (startIdx+i)), ts.ColumnKey)     // column key
		file.SetCellValue(sheetName, fmt.Sprint("H", (startIdx+i)), ts.ColumnDefault) // column default
	}
	endIdx := startIdx + len(tableSpec)

	return endIdx, nil
}

// setTableIndex : 테이블의 인덱스 설정
func setTableIndex(file *excelize.File, sheetName string, indexSpec []*model.IndexSpec, startIdx int) (int, error) {
	// 스타일 설정
	file.MergeCell(sheetName, fmt.Sprint("A", (startIdx)), fmt.Sprint("B", (startIdx)))
	file.MergeCell(sheetName, fmt.Sprint("C", (startIdx)), fmt.Sprint("I", (startIdx)))
	file.SetCellValue(sheetName, fmt.Sprint("A", (startIdx)), "Index")
	file.SetCellValue(sheetName, fmt.Sprint("C", (startIdx)), "Index Key")

	// 회색 채우기
	darkBlueFillStyle, err := file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#C0C0C0"}, Pattern: 1},
	})
	if err != nil {
		return startIdx, err
	}
	file.SetCellStyle(sheetName, fmt.Sprint("A", (startIdx)), fmt.Sprint("I", (startIdx)), darkBlueFillStyle)

	startIdx++
	for i, is := range indexSpec {
		file.MergeCell(sheetName, fmt.Sprint("A", (startIdx+i)), fmt.Sprint("B", (startIdx+i)))
		file.SetCellValue(sheetName, fmt.Sprint("A", (startIdx+i)), is.IndexName) // index name
		file.MergeCell(sheetName, fmt.Sprint("C", (startIdx+i)), fmt.Sprint("I", (startIdx+i)))
		file.SetCellValue(sheetName, fmt.Sprint("C", (startIdx+i)), is.Columns) // index column
	}
	endIdx := startIdx + len(indexSpec)

	return endIdx, nil
}

// setTableSheetFooter : 테이블 시트의 Footer 설정
func setTableSheetFooter(file *excelize.File, sheetName string, startIdx int) (int, error) {
	file.SetCellValue(sheetName, fmt.Sprint("A", (startIdx)), "etc")
	file.MergeCell(sheetName, fmt.Sprint("B", (startIdx)), fmt.Sprint("I", (startIdx)))

	return startIdx + 1, nil
}

// setTableSheetLayout : 테이블 시트 레이아웃 설정
func setTableSheetLayout(file *excelize.File, schema, sheetName, tableName, tableComment string) {
	// title column
	file.MergeCell(sheetName, fmt.Sprint("A", titleColNum), fmt.Sprint("I", titleColNum))
	file.SetCellValue(sheetName, fmt.Sprint("A", titleColNum), "Table Specification")

	createDate := time.Now().UTC().Local().Format("2006-01-02")
	createdBy := "박선하"

	// info column
	file.MergeCell(sheetName, fmt.Sprint("A", infoColNum1), fmt.Sprint("B", infoColNum1))
	file.SetCellValue(sheetName, fmt.Sprint("A", infoColNum1), "Schema")
	file.SetCellValue(sheetName, fmt.Sprint("C", infoColNum1), schema)
	file.SetCellValue(sheetName, fmt.Sprint("D", infoColNum1), "Date")
	file.MergeCell(sheetName, fmt.Sprint("E", infoColNum1), fmt.Sprint("F", infoColNum1))
	file.SetCellValue(sheetName, fmt.Sprint("E", infoColNum1), createDate)
	file.MergeCell(sheetName, fmt.Sprint("G", infoColNum1), fmt.Sprint("H", infoColNum1))
	file.SetCellValue(sheetName, fmt.Sprint("G", infoColNum1), "Creator")
	file.SetCellValue(sheetName, fmt.Sprint("I", infoColNum1), createdBy)

	// info column
	file.MergeCell(sheetName, fmt.Sprint("A", infoColNum2), fmt.Sprint("B", infoColNum2))
	file.SetCellValue(sheetName, fmt.Sprint("A", infoColNum2), "Table Name")
	file.SetCellValue(sheetName, fmt.Sprint("C", infoColNum2), tableName)
	file.MergeCell(sheetName, fmt.Sprint("D", infoColNum2), fmt.Sprint("E", infoColNum2))
	file.SetCellValue(sheetName, fmt.Sprint("D", infoColNum2), "Table Comment")
	file.MergeCell(sheetName, fmt.Sprint("F", infoColNum2), fmt.Sprint("I", infoColNum2))
	file.SetCellValue(sheetName, fmt.Sprint("F", infoColNum2), tableComment)

	// column column
	file.SetCellValue(sheetName, fmt.Sprint("A", columnColNum), "no")
	file.SetCellValue(sheetName, fmt.Sprint("B", columnColNum), "Coulmn Name")
	file.SetCellValue(sheetName, fmt.Sprint("C", columnColNum), "Column Comment")
	file.SetCellValue(sheetName, fmt.Sprint("D", columnColNum), "type")
	file.SetCellValue(sheetName, fmt.Sprint("E", columnColNum), "Length")
	file.SetCellValue(sheetName, fmt.Sprint("F", columnColNum), "Nullable")
	file.SetCellValue(sheetName, fmt.Sprint("G", columnColNum), "Key")
	file.SetCellValue(sheetName, fmt.Sprint("H", columnColNum), "Default")
	file.SetCellValue(sheetName, fmt.Sprint("I", columnColNum), "Note")
}

// setTableSheetStyle : 테이블 시트 스타일 설정
func setTableSheetStyle(file *excelize.File, sheetName string, commonStyle excelize.Style, endIdx int) error {
	// 폰트 설정
	file.SetDefaultFont("Calibri")

	// border, 가운데 정렬
	style, err := file.NewStyle(&commonStyle)
	if err != nil {
		return err
	}
	file.SetCellStyle(sheetName, "A1", fmt.Sprint("I", (endIdx-1)), style)

	// Row Height
	file.SetRowHeight(sheetName, titleColNum, 27)
	height2 := float64(21)
	for i := infoColNum1; i <= columnColNum; i++ {
		file.SetRowHeight(sheetName, i, height2)
	}

	height3 := float64(15.75)
	for i := columnColNum; i < indexColNum; i++ {
		file.SetRowHeight(sheetName, i, height3)
	}

	// Column Width
	file.SetColWidth(sheetName, "A", "A", 8)
	file.SetColWidth(sheetName, "B", "B", 18)
	file.SetColWidth(sheetName, "C", "C", 22)
	file.SetColWidth(sheetName, "D", "D", 10)
	file.SetColWidth(sheetName, "E", "E", 7)
	file.SetColWidth(sheetName, "F", "F", 10)
	file.SetColWidth(sheetName, "G", "G", 4)
	file.SetColWidth(sheetName, "H", "H", 8)
	file.SetColWidth(sheetName, "I", "I", 33)

	return nil
}

// setTitleCellStyle : Title셀 스타일 설정
func setTitleCellStyle(file *excelize.File, sheetName string, commonStyle excelize.Style, cEndIdx, fEndIdx int) error {

	titleCellStyle := commonStyle
	titleCellStyle.Fill = excelize.Fill{Type: "pattern", Color: []string{"#121f4a"}, Pattern: 1}
	titleCellStyle.Font = &excelize.Font{Bold: true, Size: 11, Color: "#F0F0F0"}

	// DarkBlue, 가운데 정렬, 글자 굵게
	darkBlueFillStyle, err := file.NewStyle(&titleCellStyle)
	if err != nil {
		return err
	}
	file.SetCellStyle(sheetName, fmt.Sprint("A", titleColNum), fmt.Sprint("A", titleColNum), darkBlueFillStyle)
	file.SetCellStyle(sheetName, fmt.Sprint("A", infoColNum1), fmt.Sprint("A", infoColNum1), darkBlueFillStyle)
	file.SetCellStyle(sheetName, fmt.Sprint("D", infoColNum1), fmt.Sprint("D", infoColNum1), darkBlueFillStyle)
	file.SetCellStyle(sheetName, fmt.Sprint("G", infoColNum1), fmt.Sprint("G", infoColNum1), darkBlueFillStyle)
	file.SetCellStyle(sheetName, fmt.Sprint("A", infoColNum2), fmt.Sprint("A", infoColNum2), darkBlueFillStyle)
	file.SetCellStyle(sheetName, fmt.Sprint("D", infoColNum2), fmt.Sprint("D", infoColNum2), darkBlueFillStyle)
	file.SetCellStyle(sheetName, fmt.Sprint("A", columnColNum), fmt.Sprint("I", columnColNum), darkBlueFillStyle)
	file.SetCellStyle(sheetName, fmt.Sprint("A", (cEndIdx)), fmt.Sprint("I", (cEndIdx)), darkBlueFillStyle)
	file.SetCellStyle(sheetName, fmt.Sprint("A", (fEndIdx-1)), fmt.Sprint("A", (fEndIdx-1)), darkBlueFillStyle)
	file.SetRowHeight(sheetName, (fEndIdx - 1), 28)

	return nil
}
