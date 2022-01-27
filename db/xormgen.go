package db

import (
	"github.com/xormplus/xorm"
	"github.com/xormplus/xorm/names"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

type XormDB struct {
	OutPut string
	DB     *xorm.Engine
}

func NewXormDB(dns string, tablePrefix string, outPut string) *XormDB {
	db, err := xorm.NewEngine("mysql", dns)
	if err != nil {
		log.Fatalf("xorm.NewEngine failed: err=%v", err)
		return nil
	}
	tbMapper := names.NewPrefixMapper(names.SnakeMapper{}, tablePrefix)
	db.SetTableMapper(tbMapper)
	return &XormDB{
		OutPut: outPut,
		DB:     db,
	}
}

func (g *XormDB) Gen() {
	tableNames, err := g.GetAllTable()
	if err != nil {
		log.Fatalf("get all table failed: err=%v", err)
		return
	}
	for _, tableName := range tableNames {
		createSql, err := g.GetTableCreateSql(tableName)
		if err != nil {
			log.Fatalf("get tabel create sql failed: tableName=%v err=%v", tableName, err)
		}
		table := ParseTableCreateSql(createSql)
		fileContent := XormModelTemp
		fileContent = strings.ReplaceAll(fileContent, "{name}", table.Name)
		fileContent = strings.ReplaceAll(fileContent, "{source_name}", table.SourceName)
		fieldList := make([]string, 0)
		for _, fieldSourceName := range table.FieldOrder {
			fieldList = append(fieldList, g.getFiledStr(table.Field[fieldSourceName], table.MaxFieldLen, table.MaxGoTypeLen))
		}
		fileContent = strings.ReplaceAll(fileContent, "{field_list}", strings.Join(fieldList, "\n"))

		if !table.IsHaveTime {
			fileContent = strings.ReplaceAll(fileContent, "\n\t\"time\"", "")
		}
		g.Write(table.SourceName, fileContent)
	}

	log.Println("\u001B[32m[SUCCESS]\u001B[0m generate xorm model file successful")
}

func (g *XormDB) Write(tableName string, content string) {
	filePath := path.Join(g.OutPut, tableName+".go")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("create file %v failed: err=%v", filePath, err)
	}
	if _, err = io.WriteString(file, content); err != nil {
		log.Fatalf("write file %v failed: err=%v", file, err)
	}
}

func (g *XormDB) getFiledStr(fieldInfo TableField, maxFieldLen int, maxGoType int) string {
	field := XormFieldTemp
	if fieldInfo.IsPrimaryKey {
		field = XormPKFieldTemp
	}
	field = strings.ReplaceAll(field, "{name}", fieldInfo.Name)
	field = strings.ReplaceAll(field, "{name_space}", strings.Repeat(" ", maxFieldLen-len(fieldInfo.Name)))
	field = strings.ReplaceAll(field, "{type}", fieldInfo.GoType)
	field = strings.ReplaceAll(field, "{type_space}", strings.Repeat(" ", maxGoType-len(fieldInfo.GoType)))
	field = strings.ReplaceAll(field, "{source_name}", fieldInfo.SourceName)

	return field
}

func (g *XormDB) GetAllTable() (result []string, err error) {
	result = make([]string, 0)

	res, err := g.DB.QueryString("show tables")
	if err != nil {
		return result, err
	}
	for _, item := range res {
		result = append(result, item["Tables_in_test"])
	}

	return result, err
}

func (g *XormDB) GetTableCreateSql(tableName string) (result string, err error) {
	res, err := g.DB.QueryString("show create table " + tableName)
	if err != nil {
		return result, err
	}
	result = res[0]["Create Table"]
	//if res := g.DB.Raw("show create table " + tableName).Scan(&result); res.Error != nil {
	//	return result, res.Error
	//}

	return result, err
}
