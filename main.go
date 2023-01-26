package main

import (
	"fmt"
	"log"
	"table-specification/db"
	"table-specification/excel"
	"table-specification/services"
)

func main() {
	fmt.Println("-----------------------------------------------------------------------")
	fmt.Println("------------------------ Table Description ----------------------------")
	fmt.Println("-----------------------------------------------------------------------")

	var err error

	// DB 연결
	err = db.ConnectDB()
	if err != nil {
		log.Fatalln(err)
	}
	subsDB := db.DB()

	// 테이블 명세 조회
	schemaList := []string{
		"mydb",
	}
	specList, err := services.GetTableSpecs(subsDB, schemaList)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(specList)

	// 테이블 명세서 엑셀로 추출
	err = excel.ExportExcel(specList)
	if err != nil {
		log.Fatalln(err)
	}

	// DB 연결 종료
	db.CloseDB()
	if err != nil {
		log.Fatalln(err)
	}

}
