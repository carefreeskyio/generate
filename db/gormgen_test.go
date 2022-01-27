package db

import "testing"

func TestGormDB_Gen(t *testing.T) {
	dns := "root:1@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	NewGormDB(dns, "", "./").Gen()
}
