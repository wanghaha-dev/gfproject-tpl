package main

import (
	"fmt"
	"github.com/wanghaha-dev/gf/os/gfile"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html/template"
	"io/ioutil"
	"os"
)

type vars struct {
	PackageName      string
	PackageFuncName  string
	PackageNameNotes string
}

func main() {
	packageName := "artwork"
	packageFuncName := "artwork"
	packageNameNotes := "艺术品"
	packageFuncNameUpper := cases.Title(language.English).String(packageFuncName)

	v := vars{
		PackageName:      packageName,
		PackageFuncName:  packageFuncNameUpper,
		PackageNameNotes: packageNameNotes,
	}

	dist := "generate/frontendApi/dist"
	tplDir := "generate/frontendApi"

	// check dir exists
	if !gfile.Exists(dist) {
		err := os.MkdirAll(dist, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// generate api file
	apiData, err := ioutil.ReadFile(fmt.Sprintf("%s/api.tpl", tplDir))
	if err != nil {
		panic(err)
	}

	apiContent := string(apiData)
	apiT := template.Must(template.New("api").Parse(apiContent))
	apiFileName := fmt.Sprintf("%v/api.%v", dist, ".js")
	apiFile, err := os.OpenFile(apiFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	err = apiT.Execute(apiFile, v)
	if err != nil {
		panic(err)
	}
	fmt.Println("finish.")
}
