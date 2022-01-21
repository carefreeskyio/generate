package main

import (
	"flag"
	"github.com/carefreeskyio/generate/generate"
	"log"
)

func main() {
	genType := flag.String("type", "service", "the generate type, optional: rpc、service")
	servicePath := flag.String("service_path", "./service", "the service dir path, required when generate rpc")
	basePath := flag.String("base_path", "/carefreesky", "the registry base path, required when generate rpc")
	serviceName := flag.String("service_name", "CarefreeSky", "the service name, required when generate rpc")
	flag.Parse()

	log.Printf("generate type=%v\n", *genType)
	log.Printf("service dir path=%v\n", *servicePath)
	log.Printf("service name=%v\n", *serviceName)
	log.Printf("registry base path=%v\n", *basePath)

	switch *genType {
	case "service":
		generate.NewService(*servicePath).GenService()
	case "rpc":
		generate.NewRpc(*serviceName, *servicePath, *basePath).GenRpc()
	default:
		log.Fatalf("unsupported genType's value: %v", *genType)
	}
}
