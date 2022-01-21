package generate

import (
	"github.com/carefreeskyio/generate/common"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

type Rpc struct {
	ServicePath string
	ServiceName string
	BasePath    string
}

var (
	rpcImportList = map[string]byte{
		"\t\"context\"": 0,
		"\t\"github.com/carefreeskyio/rpcxclient\"": 0,
		"\t\"sync\"":    0,
		"\t\"strings\"": 0,
	}
	rpcFilePath = "./rpc"
)

func NewRpc(serviceName string, servicePath string, basePath string) *Rpc {
	return &Rpc{
		ServicePath: servicePath,
		ServiceName: serviceName,
		BasePath:    basePath,
	}
}

func (g *Rpc) GenRpc() {
	rpcFileContent := common.RpcTemp

	rpcFileContent = strings.ReplaceAll(rpcFileContent, "{base_path}", g.BasePath)
	rpcFileContent = strings.ReplaceAll(rpcFileContent, "{service_name}", g.ServiceName)
	serviceFilePath := path.Join(g.ServicePath, common.ServiceFileName)
	serviceContent, err := common.GetFileContentToSlice(serviceFilePath)
	if err != nil {
		log.Fatalf("get service file content to slice failed: path=%v err=%v", serviceFilePath, err)
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
		if isImportContent {
			rpcImportList[line] = 0
			continue
		}
		if strings.Contains(line, "import (") {
			isImportContent = true
		}
		if strings.HasPrefix(line, "func NewService") {
			continue
		}
		if strings.HasPrefix(line, "func ") {
			rpcFunc += g.generateFunc(line)
		}
	}
	rpcImportListStr := g.generateImportList()

	rpcFileContent = strings.ReplaceAll(rpcFileContent, "{import_list}", rpcImportListStr)
	rpcFileContent = strings.ReplaceAll(rpcFileContent, "{rpc_func}", rpcFunc)

	g.generateDir()
	filePath := path.Join(rpcFilePath, common.RpcFileName)
	rpcFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("create file %v failed: err=%v", filePath, err)
	}
	if _, err = io.WriteString(rpcFile, rpcFileContent); err != nil {
		log.Fatalf("write file %v failed: err=%v", filePath, err)
	}

	log.Println("generate rpc file successful")
}

func (g *Rpc) generateDir() {
	filePath := "./rpc"
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(filePath, os.ModePerm)
			if err != nil {
				log.Fatalf("create rpc path failed: path=%v err=%v", filePath, err)
			}
		}
	}
}

func (g *Rpc) generateImportList() (result string) {
	result = "import (\n"

	for item, _ := range rpcImportList {
		result += item + "\n"
	}
	result += ")\n"

	return result
}

func (g *Rpc) generateFunc(line string) (result string) {
	funcInfo := common.ParseFunc(line)
	result = common.RpcFuncTemp
	result = strings.ReplaceAll(result, "{name}", funcInfo.Name)
	result = strings.ReplaceAll(result, "{param}", funcInfo.Params)
	result = strings.ReplaceAll(result, "{pass_through}", common.GetPassThroughParam(funcInfo.Params))

	return result
}
