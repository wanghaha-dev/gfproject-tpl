// {{.StructName}}
group.GET("/{{.StructName}}/:id/detail", {{.PackageName}}.{{.StructNameUpper}}.DetailById)
group.DELETE("/{{.StructName}}/:id/delete", {{.PackageName}}.{{.StructNameUpper}}.DeleteById)
group.GET("/{{.StructName}}/list", {{.PackageName}}.{{.StructNameUpper}}.List)
group.GET("/{{.StructName}}/all", {{.PackageName}}.{{.StructNameUpper}}.All)
group.POST("/{{.StructName}}/add", {{.PackageName}}.{{.StructNameUpper}}.Add)
group.POST("/{{.StructName}}/:id/update", {{.PackageName}}.{{.StructNameUpper}}.UpdateById)