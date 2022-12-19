package main

import (
	"fmt"
	"github.com/wanghaha-dev/gf/database/gdb"
	"github.com/wanghaha-dev/gf/frame/g"
	"github.com/wanghaha-dev/gf/os/gfile"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
)

type vars struct {
	Fields     map[string]*gdb.TableField
	AllFields  map[string]*gdb.TableField
	TypeName   string
	RouterName string
	ViewName   string
}

func TypeCheck(t string) bool {
	return strings.Contains(t, "varchar") || strings.Contains(t, "datetime")
}

func main() {
	typeName := "链接管理"
	routerName := "objs"
	tableName := "objs"
	viewName := cases.Title(language.English).String(routerName)
	var mapData = make(map[string]*gdb.TableField)

	fields, err := g.Model(tableName).TableFields(tableName)
	if err != nil {
		panic(err)
	}

	for k, v := range fields {
		// 不需要的字段跳过
		a := strings.Contains(k, "created_")
		b := strings.Contains(k, "updated_")
		c := strings.Contains(k, "deleted_")
		d := k == "id"

		if a || b || c || d {
			continue
		}

		mapData[k] = v
	}

	v := vars{
		Fields:     mapData,
		AllFields:  fields,
		TypeName:   typeName,
		RouterName: routerName,
		ViewName:   viewName,
	}

	dist := "generate/frontend/dist"
	tplDir := "generate/frontend"

	// check dir exists
	if !gfile.Exists(dist) {
		err := os.MkdirAll(dist, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// generate view file
	viewData, err := ioutil.ReadFile(fmt.Sprintf("%s/view.tpl", tplDir))
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

	// ##################################
	// generate router file
	routerData, err := ioutil.ReadFile(fmt.Sprintf("%s/router.tpl", tplDir))
	if err != nil {
		panic(err)
	}

	routerContent := string(routerData)
	tpl = template.New("vue")
	// 修改分隔符, 不然会和vue的分割符冲突
	newTpl = tpl.Delims("[[", "]]")
	newTpl = newTpl.Funcs(template.FuncMap{"TypeCheck": TypeCheck})

	tpl = template.Must(newTpl.Parse(routerContent))
	routerFileName := fmt.Sprintf("%s/%v_router%v", dist, v.RouterName, ".vue")

	// write file
	routerFile, err := os.OpenFile(routerFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	err = tpl.Execute(routerFile, v)
	if err != nil {
		panic(err)
	}

	// ##################################
	// generate router file
	linkData, err := ioutil.ReadFile(fmt.Sprintf("%s/link.tpl", tplDir))
	if err != nil {
		panic(err)
	}

	linkContent := string(linkData)
	tpl = template.New("vue")
	// 修改分隔符, 不然会和vue的分割符冲突
	newTpl = tpl.Delims("[[", "]]")
	newTpl = newTpl.Funcs(template.FuncMap{"TypeCheck": TypeCheck})

	tpl = template.Must(newTpl.Parse(linkContent))
	linkFileName := fmt.Sprintf("%s/%v_link%v", dist, v.RouterName, ".vue")

	// write file
	linkFile, err := os.OpenFile(linkFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	err = tpl.Execute(linkFile, v)
	if err != nil {
		panic(err)
	}
	fmt.Println("finish.")
}
