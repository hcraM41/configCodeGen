package parse

import (
	"configCodeGen/generate"
	"log"
	"os"
	"path"
)

// GetAllFileByExtensions 通过扩展名获取文件名
func GetAllFileByExtensions(folder string, extensions []string) ([]string, error) {
	if len(folder) == 0 {
		return nil, nil
	}
	files := make([]string, 0, 10)
	rd, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	for _, fi := range rd {
		if !fi.IsDir() {
			ext := path.Ext(fi.Name())
			if len(extensions) == 0 {
				files = append(files, folder+fi.Name())
			} else {
				for _, extension := range extensions {
					if extension == ext {
						files = append(files, folder+fi.Name())
					}
				}
			}
		}
	}
	return files, nil
}

// CmdParas 参数解析
func CmdParas(ExcelFilePath string) {
	files := make([]string, 0, 0)
	if len(ExcelFilePath) > 0 {
		files = append(files, ExcelFilePath)
	}
	packagePath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	outPath := packagePath + "\\cfg\\"
	for _, file := range files {
		generate.SheetGenerate(file, outPath)
	}
}
