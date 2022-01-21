package generate

import (
	"github.com/carefreeskyio/generate/common"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

type Service struct {
	WorkPath string
}

type FileContent struct {
	ServiceName        string
	ServiceNewFuncName string
	FuncList           []common.FuncInfo
}

var (
	protoNameNum = 1
	protoNameLog = make(map[string]byte)
	importList   = make(map[string]string)
)

func NewService(workPath string) *Service {
	return &Service{
		WorkPath: workPath,
	}
}

func (g *Service) GenService() {
	rd, err := ioutil.ReadDir(g.WorkPath)
	if err != nil {
		log.Fatalf("read dir %v failed: err=%v \n", g.WorkPath, err)
	}

	serviceStruct := "type Service struct {"
	newService := "return &Service{"
	serviceFunc := ""
	for _, fileInfo := range rd {
		if strings.Contains(fileInfo.Name(), "_test") || fileInfo.Name() == "service.go" {
			continue
		}
		filePath := path.Join(g.WorkPath, fileInfo.Name())
		content, err := common.GetFileContentToSlice(filePath)
		if err != nil {
			log.Fatalf("get file content to slice failed: filePath=%v err=%v", filePath, err)
		}
		parseResult := g.parseContent(content)
		if parseResult.ServiceName == "" {
			log.Printf("parse file failed: name=%v", fileInfo.Name())
			continue
		}
		serviceStruct += "\n\t" + parseResult.ServiceName + " *" + parseResult.ServiceName
		newService += "\n\t\t" + parseResult.ServiceName + ": " + parseResult.ServiceNewFuncName + "(),"
		for _, funcInfo := range parseResult.FuncList {
			serviceFunc += g.generateFunc(funcInfo, parseResult.ServiceName)
		}
	}
	serviceStruct = serviceStruct + "\n}"
	newService = newService + "\n\t}"

	content := common.ServiceTemp
	content = strings.ReplaceAll(content, "{import_list}", g.generateImportList())
	content = strings.ReplaceAll(content, "{service_struct}", serviceStruct)
	content = strings.ReplaceAll(content, "{new_service}", newService)
	content = strings.ReplaceAll(content, "{service_func}", serviceFunc)

	serviceFilePath := path.Join(g.WorkPath, common.ServiceFileName)
	serviceFile, err := os.Create(serviceFilePath)
	if err != nil {
		log.Fatalf("create file %v failed: err=%v", serviceFilePath, err)
	}
	if _, err = io.WriteString(serviceFile, content); err != nil {
		log.Fatalf("write file %v failed: err=%v", serviceFilePath, err)
	}

	log.Println("generate service file successful")
}

func (g *Service) generateFunc(funcInfo common.FuncInfo, serviceName string) (result string) {
	result = common.ServiceFuncTemp
	result = strings.ReplaceAll(result, "{name}", funcInfo.Name)
	result = strings.ReplaceAll(result, "{param}", funcInfo.Params)
	result = strings.ReplaceAll(result, "{service_name}", serviceName)
	result = strings.ReplaceAll(result, "{pass_through}", common.GetPassThroughParam(funcInfo.Params))

	return result
}

func (g *Service) generateImportList() (result string) {
	result = "import ("

	for item, alias := range importList {
		if alias != "" {
			item = alias + " \"" + item + "\""
		} else {
			item = "\"" + item + "\""
		}
		result += "\n\t" + item
	}

	result += "\n)"

	return result
}

func (g *Service) parseContent(content []string) (result FileContent) {
	result.FuncList = make([]common.FuncInfo, 0)
	isImportContent := false
	fileImportList := make(map[string]string)
	for _, line := range content {
		if strings.Contains(line, "type") && strings.Contains(line, "Service") {
			result.ServiceName = strings.Split(line, " ")[1]
		}
		if isImportContent && line == ")" {
			isImportContent = false
		}
		if isImportContent {
			oldProtoName, newProtoName := g.parseImport(line)
			fileImportList[oldProtoName] = newProtoName
			continue
		}
		if strings.Contains(line, "import (") {
			isImportContent = true
		}
		if strings.HasPrefix(line, "func New") {
			result.ServiceNewFuncName = g.parseNewFunc(line)
			continue
		}
		if strings.HasPrefix(line, "func ") {
			result.FuncList = append(result.FuncList, g.parseFunc(line, fileImportList))
		}
	}

	return result
}

func (g *Service) parseNewFunc(line string) (result string) {
	lineSli := strings.Split(line, "(")
	funcSli := strings.Split(lineSli[0], " ")
	result = funcSli[1]

	return result
}

func (g *Service) parseImport(line string) (oldProtoName string, newProtoName string) {
	line = strings.Trim(line, "\"")
	line = strings.Trim(line, "\t")
	line = strings.ReplaceAll(line, "\"", "")
	lineSli := strings.Split(line, "/")
	oldProtoName = lineSli[len(lineSli)-1]
	isContainsSpace := false
	if strings.Contains(line, " ") {
		isContainsSpace = true
		sli := strings.Split(line, " ")
		oldProtoName = sli[0]
		line = sli[1]
	}
	newProtoName, ok := importList[line]
	if ok {
		return oldProtoName, newProtoName
	}

	newProtoName = oldProtoName
	if _, ok = protoNameLog[newProtoName]; !ok {
		importList[line] = ""
		if isContainsSpace {
			importList[line] = newProtoName
		}
		protoNameLog[newProtoName] = 0
		return oldProtoName, ""
	}
	newProtoName += strconv.Itoa(protoNameNum)
	importList[line] = newProtoName
	protoNameLog[newProtoName] = 0
	protoNameNum++

	return oldProtoName, newProtoName
}

func (g *Service) parseFunc(line string, fileImportList map[string]string) (result common.FuncInfo) {
	result = common.ParseFunc(line)

	for oldProtoName, newProtoName := range fileImportList {
		if newProtoName == "" {
			continue
		}
		if !strings.Contains(result.Params, oldProtoName) {
			continue
		}
		result.Params = strings.Replace(result.Params, oldProtoName+".", newProtoName+".", -1)
	}

	return result
}
