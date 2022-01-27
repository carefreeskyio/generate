package db

var GormModelTemp = `package db

import (
	"context"
	"github.com/carefreex-io/dbdao/gormdb"
	"gorm.io/gorm"
	"time"
)

type {name} struct {
{field_list}
}

type {name}Db struct {
	DB *gorm.DB
}

func New{name}Db(ctx context.Context, arg ...bool) (db *{name}Db) {
	db = &{name}Db{
		DB: gormdb.Read,
	}
	if len(arg) != 0 && arg[0] {
		db.DB = gormdb.Write
	}
	db.DB.WithContext(ctx)

	return db
}

func (d *{name}Db) TableName() string {
	return "{source_name}"
}

`

var GormPKFieldTemp = "\t{name}{name_space} {type}{type_space} `gorm:\"primaryKey\" json:\"{source_name}\"`"
var GormFieldTemp = "\t{name}{name_space} {type}{type_space} `json:\"{source_name}\"`"
