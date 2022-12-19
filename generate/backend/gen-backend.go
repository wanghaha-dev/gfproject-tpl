package main

import (
	"fmt"
	"html/template"
	"os"

	"github.com/wanghaha-dev/gf/os/gfile"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type vars struct {
	PackageName     string
	StructName      string
	StructNameUpper string
	DBName          string
	TableName       string
	Cache           bool
}

func main() {
	packageName := "artlib"
	structName := "alIp"
	dbName := "default"
	tableName := "al_ip"
	structNameUpper := cases.Title(language.English).String(structName)
	cache := true

	v := vars{
		PackageName:     packageName,
		StructName:      structName,
		StructNameUpper: structNameUpper,
		TableName:       tableName,
		DBName:          dbName,
		Cache:           cache,
	}

	dist := "generate/backend/dist"
	tplDir := "generate/backend"

	// check dir exists
	createFileNotExists(dist)

	// generate api file
	apiContent := readFileToString(fmt.Sprintf("%s/api_template.tpl", tplDir))

	apiT := template.Must(template.New("api").Parse(apiContent))
	apiFileName := fmt.Sprintf("%s/api_%v%v", dist, v.StructName, ".go")
	apiFile := getFileObject(apiFileName)
	err := apiT.Execute(apiFile, v)
	getErr(err)

	// generate router
	routerContent := readFileToString(fmt.Sprintf("%s/router_template.tpl", tplDir))
	routerT := template.Must(template.New("api").Parse(routerContent))
	routerFileName := fmt.Sprintf("%s/router_%v%v", dist, v.StructName, ".go")

	routerFile := getFileObject(routerFileName)
	err = routerT.Execute(routerFile, v)
	getErr(err)
	fmt.Println("finish.")
}

func getErr(err error) {
	if err != nil {
		panic(err)
	}
}

// createFileNotExists 不存在则创建文件活目录
func createFileNotExists(fPath string) {
	if !gfile.Exists(fPath) {
		err := os.MkdirAll(fPath, os.ModePerm)
		getErr(err)
	}
}

// readFileToString 读取文件字符串
func readFileToString(fPath string) string {
	bytes, err := os.ReadFile(fPath)
	getErr(err)
	return string(bytes)
}

// getFileObject 获取文件对象
func getFileObject(fPath string) *os.File {
	file, err := os.OpenFile(fPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	getErr(err)
	return file
}
