package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

type GormDB struct {
	OutPut string
	DB     *gorm.DB
}

func NewGormDB(dns string, tablePrefix string, outPut string) *GormDB {
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("gorm.Open failed: err=%v", err)
		return nil
	}
	return &GormDB{
		OutPut: outPut,
		DB:     db,
	}
}

func (g *GormDB) Gen() {
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
		fileContent := GormModelTemp
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
	log.Println("\u001B[32m[SUCCESS]\u001B[0m generate gorm model file successful")
}

func (g *GormDB) Write(tableName string, content string) {
	filePath := path.Join(g.OutPut, tableName+".go")
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("create file %v failed: err=%v", filePath, err)
	}
	if _, err = io.WriteString(file, content); err != nil {
		log.Fatalf("write file %v failed: err=%v", file, err)
	}
}

func (g *GormDB) getFiledStr(fieldInfo TableField, maxFieldLen int, maxGoType int) string {
	field := GormFieldTemp
	if fieldInfo.IsPrimaryKey {
		field = GormPKFieldTemp
	}
	field = strings.ReplaceAll(field, "{name}", fieldInfo.Name)
	field = strings.ReplaceAll(field, "{name_space}", strings.Repeat(" ", maxFieldLen-len(fieldInfo.Name)))
	field = strings.ReplaceAll(field, "{type}", fieldInfo.GoType)
	field = strings.ReplaceAll(field, "{type_space}", strings.Repeat(" ", maxGoType-len(fieldInfo.GoType)))
	field = strings.ReplaceAll(field, "{source_name}", fieldInfo.SourceName)

	return field
}

func (g *GormDB) GetAllTable() (result []string, err error) {
	if res := g.DB.Raw("show tables").Scan(&result); res.Error != nil {
		return result, res.Error
	}

	return result, err
}

func (g *GormDB) GetTableCreateSql(tableName string) (result string, err error) {
	var sql []string
	if res := g.DB.Raw("show create table " + tableName).Scan(&sql); res.Error != nil {
		return result, res.Error
	}
	result = sql[0]

	return result, err
}
