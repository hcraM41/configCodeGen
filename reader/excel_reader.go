package reader

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ExcelReader struct {
	file *excelize.File
}

func NewExcelReader(fileName string) *ExcelReader {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		panic(fmt.Sprintf("读取文件失败,文件：%v,err:%v", fileName, err))
	}
	return &ExcelReader{file: f}

}

func (er *ExcelReader) Close() {
	if err := er.file.Close(); err != nil {
		fmt.Println(err)
	}
}

func (er *ExcelReader) GetSheetNames() []string {
	return er.file.GetSheetList()
}

func (er *ExcelReader) GetRows(sheetName string) [][]string {
	rows, err := er.file.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return rows
}

func (er *ExcelReader) GetCols(sheetName string) [][]string {
	rows, err := er.file.GetCols(sheetName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return rows
}
