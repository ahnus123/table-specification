package excel

import (
	"errors"
	"fmt"
	"table-specification/message"

	"github.com/xuri/excelize/v2"
)

var (
	listTitleColNum = 1
	listColNum      = 2
	listFinalColNum = -1
)

// CreateListSheet : 목차 시트 생성
func CreateListSheet(file *excelize.File, specList map[string][]*message.TableInfo) error {
	if specList == nil {
		return errors.New("there is no list")
	}

	sheetName := "List"
	idx, err := file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	setListSheetLayout(file, sheetName)
	tlEndIdx, err := setTableList(file, sheetName, specList, listColNum+1)
	if err != nil {
		return err
	}
	listFinalColNum = tlEndIdx
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
	setListSheetStyle(file, sheetName, *commonStyle, listFinalColNum)
	setListTitleCellStyle(file, sheetName, *commonStyle)

	return nil
}

// setListSheetLayout : 목차 시트 레이아웃 설정
func setListSheetLayout(file *excelize.File, sheetName string) {
	// title column
	file.MergeCell(sheetName, fmt.Sprint("A", listTitleColNum), fmt.Sprint("D", listTitleColNum))
	file.SetCellValue(sheetName, fmt.Sprint("A", listTitleColNum), "Table List")

	// info column
	file.SetCellValue(sheetName, fmt.Sprint("A", listColNum), "no")
	file.SetCellValue(sheetName, fmt.Sprint("B", listColNum), "Schema")
	file.SetCellValue(sheetName, fmt.Sprint("C", listColNum), "Table Name")
	file.SetCellValue(sheetName, fmt.Sprint("D", listColNum), "Table Comment")
}

// setTableList : 테이블 목차 설정
func setTableList(file *excelize.File, sheetName string, specList map[string][]*message.TableInfo, startIdx int) (int, error) {
	i := 0
	for schema, tableList := range specList {
		for _, ti := range tableList {
			if len(ti.TableSpec) > 0 {
				file.SetCellValue(sheetName, fmt.Sprint("A", (startIdx+i)), (i + 1))
				file.SetCellValue(sheetName, fmt.Sprint("B", (startIdx+i)), schema)                                     // table schema
				file.SetCellValue(sheetName, fmt.Sprint("C", (startIdx+i)), ti.TableName)                               // table name
				file.SetCellValue(sheetName, fmt.Sprint("D", (startIdx+i)), ti.TableSpec[ti.TableName][0].TableComment) // table comment
				i++
			}
		}
	}
	endIdx := startIdx + i

	return endIdx, nil
}

// setListSheetStyle : 목차 시트 스타일 설정
func setListSheetStyle(file *excelize.File, sheetName string, commonStyle excelize.Style, endIdx int) error {
	// 폰트 설정
	file.SetDefaultFont("Calibri")

	// border, 가운데 정렬
	style, err := file.NewStyle(&commonStyle)
	if err != nil {
		return err
	}
	file.SetCellStyle(sheetName, "A1", fmt.Sprint("D", (endIdx-1)), style)

	// Row Height
	file.SetRowHeight(sheetName, listTitleColNum, 27)
	file.SetRowHeight(sheetName, listColNum, 21)
	height := float64(15.75)
	for i := listColNum; i < listFinalColNum; i++ {
		file.SetRowHeight(sheetName, i, height)
	}

	// Column Width
	file.SetColWidth(sheetName, "A", "A", 8)
	file.SetColWidth(sheetName, "B", "B", 16)
	file.SetColWidth(sheetName, "C", "C", 25)
	file.SetColWidth(sheetName, "D", "D", 45)

	return nil
}

// setListTitleCellStyle : Title셀 스타일 설정
func setListTitleCellStyle(file *excelize.File, sheetName string, commonStyle excelize.Style) error {

	titleCellStyle := commonStyle
	titleCellStyle.Fill = excelize.Fill{Type: "pattern", Color: []string{"#121f4a"}, Pattern: 1}
	titleCellStyle.Font = &excelize.Font{Bold: true, Size: 11, Color: "#F0F0F0"}

	// DarkBlue, 가운데 정렬, 글자 굵게
	darkBlueFillStyle, err := file.NewStyle(&titleCellStyle)
	if err != nil {
		return err
	}
	file.SetCellStyle(sheetName, fmt.Sprint("A", listTitleColNum), fmt.Sprint("D", listTitleColNum), darkBlueFillStyle)
	file.SetCellStyle(sheetName, fmt.Sprint("A", listColNum), fmt.Sprint("D", listColNum), darkBlueFillStyle)

	return nil
}
