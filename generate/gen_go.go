package generate

import (
	"configCodeGen/reader"
	"fmt"
	"regexp"
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/stoewer/go-strcase"
)

const (
	ConfPrefix = "cfg"
	ConfSuffix = "Cfg"
	pattern    = "[a-zA-Z]+"
	All        = "All"
	Decode     = "Decode"
)

const (
	ConfigInfo      = "ConfigInfo"
	PointConfigInfo = "*ConfigInfo"
	PointConfigItem = "*ConfigItem"
)

var (
	regex = regexp.MustCompile(pattern)
)

func GetStructName(sheetName string) string {
	return fmt.Sprintf("%v%v", strings.Trim(sheetName, " "), ConfSuffix)
}

// SheetGenerate 工作表转golang代码
func SheetGenerate(excelName, outPutPath string) {
	excelReader := reader.NewExcelReader(excelName)
	sheets := excelReader.GetSheetNames()
	// ConfigGenerate(sheets, outPutPath)
	for _, sheet := range sheets {
		// 获取cols
		cols := excelReader.GetCols(sheet)
		sheet = strcase.UpperCamelCase(regex.FindString(sheet))
		sourceSheet := strings.ReplaceAll(sheet, "FH", "")
		structName := GetStructName(sheet)
		f := NewFile(ConfPrefix)
		statement := f.Type().Id(structName)
		allVarMapName := All + sheet + "Map"
		allVarListName := All + sheet + "List"
		// 生成全局变量
		f.Var().Id(allVarMapName).Op("=").Make(Map(Int32()).Id("*" + structName))
		f.Var().Id(allVarListName).Op("=").Make(Index().Id("*" + structName).Id(",0"))
		// step:生成结构体
		fields := make([]Code, 0, len(cols))
		for _, col := range cols {
			if len(col) >= 2 {
				fieldName := strcase.UpperCamelCase(col[0]) // 字段名称 大驼峰形式
				fieldType := strings.ToLower(col[1])        // 字段类型
				fieldComment := col[2]                      // 字段描述
				var field Code
				switch fieldType {
				case "int":
					field = Id(fieldName).Int32().Tag(map[string]string{"json": strcase.SnakeCase(fieldName)}).Comment(fieldComment)
				case "long":
					field = Id(fieldName).Int64().Tag(map[string]string{"json": strcase.SnakeCase(fieldName)}).Comment(fieldComment)
				case "list<string>":
					field = Id(fieldName).Index().String().Tag(map[string]string{"json": strcase.SnakeCase(fieldName)}).Comment(fieldComment)
				case "list<int>":
					field = Id(fieldName).Index().Int32().Tag(map[string]string{"json": strcase.SnakeCase(fieldName)}).Comment(fieldComment)
				case "list<long>":
					field = Id(fieldName).Index().Int64().Tag(map[string]string{"json": strcase.SnakeCase(fieldName)}).Comment(fieldComment)
				case "int[][]":
					field = Id(fieldName).Index().Index().Int32().Tag(map[string]string{"json": strcase.SnakeCase(fieldName)}).Comment(fieldComment)
				case "long[][]":
					field = Id(fieldName).Index().Index().Int64().Tag(map[string]string{"json": strcase.SnakeCase(fieldName)}).Comment(fieldComment)
				case "string":
					field = Id(fieldName).String().Tag(map[string]string{"json": strcase.SnakeCase(fieldName)}).Comment(fieldComment)
				}
				fields = append(fields, field)
			}
		}
		statement.Struct(fields...).Line()
		// data := value
		// 生成函数
		f.Func().Id(fmt.Sprintf("%s%s", Decode, sheet)).Params(Id("items "+PointConfigInfo)).Id("error").Block(
			For(Id("_,value").Op(":=").Range().Id("items."+sheet)).Block(
				Id("data").Op(":=").Id("value"),
				Id(allVarMapName).Id("[value.Id]").Op("=").Id("&data"),
				Id(allVarListName).Op("=").Append(Id(allVarListName), Id("&data"))),
			Return(Id("nil")),
		).Line()
		f.Func().Id(fmt.Sprintf("Get%sById", sheet)).Params(Id("id").Int32()).Id("*" + GetStructName(sheet)).Block(
			Return(Id(allVarMapName).Id(fmt.Sprintf("[id]"))),
		).Line()
		f.Func().Id(fmt.Sprintf("GetAll%s", sheet)).Params().Id("[]*" + GetStructName(sheet)).Block(
			Return(Id(allVarListName)),
		).Line()

		// RName := fmt.Sprintf("%v%v", structName, "Register")
		// statementR := f.Type().Id(RName)
		// fieldsR := make([]Code, 0, len(sheets))
		// fieldsR = append(fieldsR, Id(PointConfigItem))
		//
		// register := make([]Code, 0, len(sheets))
		// for _, s := range sheets {
		// 	pureSheet := regex.FindString(s)
		// 	fieldName := strcase.UpperCamelCase(pureSheet)
		// 	field := Id(fieldName).Id("[]" + fieldName + "Cfg").Tag(map[string]string{"json": pureSheet}).Comment(s)
		// 	fieldsR = append(fieldsR, field)
		// 	register = append(register, Id("FileCfgRegister").Call(Lit(pureSheet), Id(Decode+fieldName)))
		// }
		// statementR.Struct(fieldsR...)
		// // 定义init方法
		// f.Func().Id("init").Params().Block(register...)

		_ = f.Save(outPutPath + strcase.SnakeCase("Config_"+sourceSheet) + ".go")
	}
}

// ConfigGenerate ConfigInfo 生成
func ConfigGenerate(sheets []string, outPutPath string) {
	f := NewFile(ConfPrefix)
	statement := f.Type().Id(ConfigInfo)
	fields := make([]Code, 0, len(sheets))
	fields = append(fields, Id(PointConfigItem))

	register := make([]Code, 0, len(sheets))
	for _, sheet := range sheets {
		pureSheet := regex.FindString(sheet)
		fieldName := strcase.UpperCamelCase(pureSheet)
		field := Id(fieldName).Id("[]" + fieldName + "Cfg").Tag(map[string]string{"json": pureSheet}).Comment(sheet)
		fields = append(fields, field)
		register = append(register, Id("FileCfgRegister").Call(Lit(pureSheet), Id(Decode+fieldName)))
	}
	statement.Struct(fields...)
	// 定义init方法
	f.Func().Id("init").Params().Block(register...)
	_ = f.Save(outPutPath + "config_mgr" + ".go")
}
