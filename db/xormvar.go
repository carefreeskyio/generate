package db

var XormModelTemp = `package db

import (
	"context"
	"github.com/carefreex-io/dbdao/xormdb"
	"github.com/xormplus/xorm"
	"time"
)

type {name} struct {
{field_list}
}

type {name}Db struct {
	DB *xorm.Engine
}

func New{name}Db(ctx context.Context, arg ...bool) (db *{name}Db) {
	db = &{name}Db{
		DB: xormdb.Read,
	}
	if len(arg) != 0 && arg[0] {
		db.DB = xormdb.Write
	}
	db.DB.SetDefaultContext(ctx)

	return db
}

func (d *{name}Db) TableName() string {
	return "{source_name}"
}

`

var XormPKFieldTemp = "\t{name}{name_space} {type}{type_space} `xorm:\"pk\" json:\"{source_name}\"`"
var XormFieldTemp = "\t{name}{name_space} {type}{type_space} `json:\"{source_name}\"`"
