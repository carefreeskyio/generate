package rpc

import (
	"github.com/carefreex-io/generate/common"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

type Rpc struct {
	ServiceName string
	BasePath    string
}

func NewRpc(serviceName string, basePath string) *Rpc {
	return &Rpc{
		ServiceName: serviceName,
		BasePath:    basePath,
	}
}

func (r *Rpc) GenRpc() {
	rpcFileContent := fileTemp

	rpcFileContent = strings.ReplaceAll(rpcFileContent, "{base_path}", r.BasePath)
	rpcFileContent = strings.ReplaceAll(rpcFileContent, "{service_name}", r.ServiceName)
	servicePath := path.Join("./", "service.go")
	serviceContent, err := common.GetFileContentToSlice(servicePath)
	if err != nil {
		log.Fatalf("get service file content to slice failed: path=%v err=%v", servicePath, err)
	}

	isImportContent := false
	rpcFunc := ""
	for _, line := range serviceContent {
		if line == "" {
			continue
		}
		if isImportContent && line == ")" {
			isImportContent = false
		}
		if isImportContent && !strings.Contains(line, "service\"") && !strings.HasSuffix(line, "validation\"") {
			importList[line] = 0
			continue
		}
		if strings.Contains(line, "import (") {
			isImportContent = true
		}
		if strings.HasPrefix(line, "func NewService") {
			continue
		}
		if strings.HasPrefix(line, "func ") {
			rpcFunc += r.generateFunc(line)
		}
	}
	rpcImportListStr := r.generateImportList()

	rpcFileContent = strings.ReplaceAll(rpcFileContent, "{import_list}", rpcImportListStr)
	rpcFileContent = strings.ReplaceAll(rpcFileContent, "{rpc_func}", rpcFunc)

	r.generateDir()
	fullFilePath := path.Join(filePath, strings.ToLower(r.ServiceName)+".go")
	rpcFile, err := os.Create(fullFilePath)
	if err != nil {
		log.Fatalf("create file %v failed: err=%v", fullFilePath, err)
	}
	if _, err = io.WriteString(rpcFile, rpcFileContent); err != nil {
		log.Fatalf("write file %v failed: err=%v", fullFilePath, err)
	}

	log.Println("\u001B[32m[SUCCESS]\u001B[0m generate rpc file successful")
}

func (r *Rpc) generateDir() {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(filePath, os.ModePerm)
			if err != nil {
				log.Fatalf("create rpc path failed: path=%v err=%v", filePath, err)
			}
		}
	}
}

func (r *Rpc) generateImportList() (result string) {
	result = "import (\n"

	for item := range importList {
		result += item + "\n"
	}
	result += ")\n"

	return result
}

func (r *Rpc) generateFunc(line string) (result string) {
	funcInfo := common.ParseFunc(line)
	result = funcTemp
	result = strings.ReplaceAll(result, "{name}", funcInfo.Name)
	result = strings.ReplaceAll(result, "{param}", funcInfo.Params)
	result = strings.ReplaceAll(result, "{pass_through}", common.GetPassThroughParam(funcInfo.Params))

	return result
}
