package db

import "testing"

func TestXormDB_Gen(t *testing.T) {
	dns := "root:1@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	NewXormDB(dns, "", "./").Gen()
}
