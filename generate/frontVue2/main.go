package main

import (
	"fmt"
	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/os/gfile"
	"os"
	"sort"
	"strings"
	"text/template"
)

type DataStruct struct {
	ViewName   string
	Data       map[int]interface{}
	RouterName string
	BaseURL    string
}

func main() {
	tableName := "m_project"
	fields, err := g.Model(tableName).TableFields(tableName)
	if err != nil {
		panic(err)
	}

	fieldMap := make(map[int]interface{})
	for _, v := range fields {
		//v.Type
		//v.Comment
		//v.Index
		//v.Default
		//v.Name
		//v.Extra
		//v.Key
		//v.Null
		fieldMap[v.Index] = v
	}

	var keys []int
	for key, _ := range fieldMap {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	for item := range keys {
		fmt.Println(fieldMap[item])
	}

	v := DataStruct{
		ViewName:   "project",
		Data:       fieldMap,
		RouterName: "project",
		BaseURL:    "http://127.0.0.1:8199",
	}

	dist := "generate/front2/dist"
	tplDir := "generate/front2"

	// check dir exists
	if !gfile.Exists(dist) {
		err := os.MkdirAll(dist, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// generate view file
	viewData, err := os.ReadFile(fmt.Sprintf("%s/view.vue", tplDir))
	if err != nil {
		panic(err)
	}

	viewContent := string(viewData)
	tpl := template.New("vue")
	// 修改分隔符, 不然会和vue的分割符冲突
	newTpl := tpl.Delims("[[", "]]")
	newTpl = newTpl.Funcs(template.FuncMap{"TypeCheck": TypeCheck})

	tpl = template.Must(newTpl.Parse(viewContent))
	viewFileName := fmt.Sprintf("%s/%vView%v", dist, v.ViewName, ".vue")

	// write file
	viewFile, err := os.OpenFile(viewFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	err = tpl.Execute(viewFile, v)
	if err != nil {
		panic(err)
	}
}

func TypeCheck(t string) bool {
	return strings.Contains(t, "varchar") || strings.Contains(t, "datetime")
}
