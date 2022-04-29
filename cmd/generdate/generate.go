package main

import (
	"strings"

	"gorm.io/gen"
	"gorm.io/gen/examples/dal"
)

const MySQLDSN = "root:Qwer1234@tcp(tencent:3306)/crontab?charset=utf8mb4&parseTime=True"

func init() {
	dal.DB = dal.ConnectDB(MySQLDSN).Debug()

	// prepare(dal.DB) // prepare table for generate
}

// dataMap mapping relationship
var dataMap = map[string]func(detailType string) (dataType string){
	// int mapping
	"int": func(detailType string) (dataType string) { return "int32" },

	// bool mapping
	"tinyint": func(detailType string) (dataType string) {
		if strings.HasPrefix(detailType, "tinyint(1)") {
			return "bool"
		}
		return "byte"
	},
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "../../app/master/repository/dal/query",
		ModelPkgPath: "../../app/master/repository/dal/model",

		// generate model global configuration
		FieldNullable:     true, // generate pointer when field is nullable
		FieldCoverable:    true, // generate pointer when field has default value
		FieldWithIndexTag: true, // generate with gorm index tag
		FieldWithTypeTag:  true, // generate with gorm column type tag
	})

	g.UseDB(dal.DB)

	// specify diy mapping relationship
	g.WithDataTypeMap(dataMap)

	// generate all field with json tag end with "_example"
	g.WithJSONTagNameStrategy(func(c string) string { return c })

	mytable := g.GenerateModel("job")
	g.ApplyBasic(mytable)
	// g.ApplyBasic(g.GenerateAllTable()...) // generate all table in db server

	g.Execute()
}
