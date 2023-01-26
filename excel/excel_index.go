package excel

import (
	"errors"

	"github.com/xuri/excelize/v2"
)

// CreateIntroSheet : 표지 시트 생성
func CreateIntroSheet(file *excelize.File) error {
	if file == nil {
		return errors.New("there is no file")
	}

	sheetName := "Index"
	idx, err := file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	setIndexSheetLayout(file, sheetName)
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
		Font:      &excelize.Font{Size: 40},
	}
	setIndexSheetStyle(file, sheetName, *commonStyle)

	return nil
}

// setIndexSheetLayout : 표지 시트 레이아웃 설정
func setIndexSheetLayout(file *excelize.File, sheetName string) {
	file.MergeCell(sheetName, "A1", "I40")
	file.SetCellValue(sheetName, "A1", "Table \nSpecification")
}

// setIndexSheetStyle : 표지 시트 스타일 설정
func setIndexSheetStyle(file *excelize.File, sheetName string, commonStyle excelize.Style) error {
	// 폰트 설정
	file.SetDefaultFont("Calibri")

	// border, 가운데 정렬
	style, err := file.NewStyle(&commonStyle)
	if err != nil {
		return err
	}
	file.SetCellStyle(sheetName, "A1", "A1", style)

	return nil
}
