package main

import (
	"flag"
	"fmt"
	"github.com/carefreex-io/generate/common"
	"github.com/carefreex-io/generate/db"
	"github.com/carefreex-io/generate/rpc"
	"github.com/carefreex-io/generate/service"
	"log"
)

const help = `CarefreeX Generate

	You can use this program to generate RPC service file or RPC client file or Gorm model file or Xorm model file

The parameter description:
-t string
	the generate type, option values: rpc、service、gorm、xorm
-base_path string
	the registry base path, required when generate rpc client file
-service_name string
	the service name, required when generate rpc client file
-db_dns string
	the database dns, required when generate gorm model file or xorm model file
-db_tp string
	the database's table prefix, required when generate gorm model file or xorm model file
-db_out string
	the table model file output dir
`

func main() {
	genType := flag.String("t", "", "the generate type, optional: rpc、service、gorm、xorm. see detail:./generate rpc help")
	basePath := flag.String("base_path", "", "the registry base path, required when generate rpc file")
	serviceName := flag.String("service_name", "", "the service name, required when generate rpc file")
	dns := flag.String("db_dns", "", "the database dns, required when generate gorm model file or xorm model file")
	tablePrefix := flag.String("db_tp", "", "the database's table prefix, required when generate gorm model file or xorm model file")
	outPut := flag.String("db_out", "./app/dao/db/", "the table model file output dir")

	flag.Usage = func() {
		fmt.Print(help)
	}
	flag.Parse()

	log.Printf("the generate type=%v\n", *genType)

	switch *genType {
	case "service":
		servicePath := "./app/service"
		log.Printf("the service path=%v", servicePath)
		service.NewService(servicePath).GenService()
	case "rpc":
		if *basePath == "" || *serviceName == "" {
			log.Fatalf("%v the base_path and service_name was required when generate rpc client file", common.ErrorStr)
		}
		log.Printf("the registry base path=%v\n", *basePath)
		log.Printf("the service name=%v\n", *serviceName)
		rpc.NewRpc(*serviceName, *basePath).GenRpc()
	case "gorm":
		if *dns == "" {
			log.Fatalf("%v the db_dns was required when generate gorm model file", common.ErrorStr)
		}
		log.Printf("the database dns=%v", *dns)
		log.Printf("the database's table prefix=%v", *tablePrefix)
		db.NewGormDB(*dns, *tablePrefix, *outPut).Gen()
	case "xorm":
		if *dns == "" {
			log.Fatalf("%v the db_dns was required when generate gorm model file", common.ErrorStr)
		}
		log.Printf("the database dns=%v", *dns)
		log.Printf("the database's table prefix=%v", *tablePrefix)
		db.NewXormDB(*dns, *tablePrefix, *outPut).Gen()
	default:
		log.Fatalf("%v unsupported genType's value: %v", *genType, common.ErrorStr)
	}
}
