package main

import (
	"configCodeGen/parse"
	"flag"
	"fmt"
	"os"
)

var (
	Help    bool
	Example string
)

var (
	ExcelFilePath string // 文件路径
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	// step:参数定义
	flag.StringVar(&ExcelFilePath, "excel_file_path", "", "Excel文件路径[选填]")
	flag.BoolVar(&Help, "h", false, "帮助")
	flag.StringVar(&Example, "example", "./config_gen.exe -excel_file_path 'E:\\goProject\\fishingdoc\\策划文档\\Config\\Excel\\activeGif活动礼包.xlsx'", "示例")
}

func main() {
	flag.Parse()
	if Help {
		Usage()
		return
	}
	if len(ExcelFilePath) == 0 {
		fmt.Fprintf(os.Stderr, "%s", "请指定excel文件路径")
		Usage()
		return
	}
	// step:参数解析
	parse.CmdParas(ExcelFilePath)

}
