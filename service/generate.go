package service

import (
	"github.com/carefreex-io/generate/common"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

type (
	Service struct {
		ServicePath string
	}

	FileContent struct {
		ServiceName        string
		ServiceNewFuncName string
		FuncList           []common.FuncInfo
	}

	PkgName = string

	PkgInfo struct {
		Pkg         string
		DefaultName string
		OldName     string
		Name        string
	}
)

var (
	pkgList      = make(map[PkgName]PkgInfo)
	pkgNameCount = make(map[PkgName]int)
)

func NewService(servicePath string) *Service {
	return &Service{
		ServicePath: servicePath,
	}
}

func (g *Service) GenService() {
	rd, err := ioutil.ReadDir(g.ServicePath)
	if err != nil {
		log.Fatalf("read dir %v failed: err=%v \n", g.ServicePath, err)
	}

	serviceStruct := "type Service struct {"
	newService := "return &Service{"
	serviceFunc := ""
	for _, fileInfo := range rd {
		if strings.Contains(fileInfo.Name(), "_test") || fileInfo.Name() == "service.go" {
			continue
		}
		filePath := path.Join(g.ServicePath, fileInfo.Name())
		content, err := common.GetFileContentToSlice(filePath)
		if err != nil {
			log.Fatalf("get file content to slice failed: filePath=%v err=%v", filePath, err)
		}
		parseResult := g.parseContent(content)
		if parseResult.ServiceName == "" {
			log.Printf("parse file failed: name=%v", fileInfo.Name())
			continue
		}
		serviceStruct += "\n\t" + parseResult.ServiceName + " *service." + parseResult.ServiceName
		newService += "\n\t\t" + parseResult.ServiceName + ": service." + parseResult.ServiceNewFuncName + "(),"
		for _, funcInfo := range parseResult.FuncList {
			serviceFunc += g.generateFunc(funcInfo, parseResult.ServiceName)
		}
	}
	serviceStruct = serviceStruct + "\n}"
	newService = newService + "\n\t}"

	content := fileTemp
	content = strings.ReplaceAll(content, "{import_list}", g.generateImportList())
	content = strings.ReplaceAll(content, "{service_struct}", serviceStruct)
	content = strings.ReplaceAll(content, "{new_service}", newService)
	content = strings.ReplaceAll(content, "{service_func}", serviceFunc)

	serviceFilePath := path.Join("./", fileName)
	serviceFile, err := os.Create(serviceFilePath)
	if err != nil {
		log.Fatalf("create file %v failed: err=%v", serviceFilePath, err)
	}
	if _, err = io.WriteString(serviceFile, content); err != nil {
		log.Fatalf("write file %v failed: err=%v", serviceFilePath, err)
	}

	log.Printf("%v generate service file successful", common.SuccessStr)
}

func (g *Service) getModName() string {
	slice, err := common.GetFileContentToSlice("./go.mod")
	if err != nil {
		log.Fatalf("%v get mod name failed: err=%v", common.ErrorStr, err)
	}

	return strings.Split(slice[0], " ")[1]
}

func (g *Service) generateFunc(funcInfo common.FuncInfo, serviceName string) (result string) {
	result = funcTemp
	result = strings.ReplaceAll(result, "{name}", funcInfo.Name)
	result = strings.ReplaceAll(result, "{param}", funcInfo.Params)
	result = strings.ReplaceAll(result, "{service_name}", serviceName)
	result = strings.ReplaceAll(result, "{pass_through}", common.GetPassThroughParam(funcInfo.Params))

	return result
}

func (g *Service) generateImportList() (result string) {
	result = "import (\n\t\"" + g.getModName() + "/app/service\""

	for _, pkgInfo := range pkgList {
		pkg := pkgInfo.Pkg
		if pkgInfo.Name != pkgInfo.DefaultName {
			pkg = pkgInfo.Name + " \"" + pkgInfo.Pkg + "\""
		} else {
			pkg = "\"" + pkgInfo.Pkg + "\""
		}
		result += "\n\t" + pkg
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
			oldProtoName, newProtoName := g.parsePkg(line)
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

func (g *Service) getPkgInfo(line string) (result PkgInfo) {
	line = strings.Trim(line, "\"")
	line = strings.Trim(line, "\t")
	line = strings.ReplaceAll(line, "\"", "")
	lineSli := strings.Split(line, " ")
	if len(lineSli) == 1 {
		pkg := lineSli[0]
		pkgSli := strings.Split(pkg, "/")

		result = PkgInfo{
			Pkg:         pkg,
			Name:        pkgSli[len(pkgSli)-1],
			DefaultName: pkgSli[len(pkgSli)-1],
			OldName:     pkgSli[len(pkgSli)-1],
		}
	} else {
		pkg := lineSli[1]
		pkgSli := strings.Split(pkg, "/")

		result = PkgInfo{
			Pkg:         pkg,
			Name:        pkgSli[len(pkgSli)-1],
			DefaultName: pkgSli[len(pkgSli)-1],
			OldName:     lineSli[0],
		}
	}

	return result
}

func (g *Service) parsePkg(line string) (oldPkgName string, newPkgName string) {
	pkgInfo := g.getPkgInfo(line)

	if _, ok := pkgList[pkgInfo.DefaultName]; !ok {
		pkgNameCount[pkgInfo.DefaultName] = 0
		pkgList[pkgInfo.DefaultName] = pkgInfo

		return pkgInfo.OldName, pkgInfo.Name
	}
	if pkgList[pkgInfo.DefaultName].Pkg == pkgInfo.Pkg {
		return pkgInfo.OldName, pkgInfo.DefaultName
	}
	pkgNameCount[pkgInfo.DefaultName]++
	pkgInfo.Name = pkgInfo.DefaultName + strconv.Itoa(pkgNameCount[pkgInfo.DefaultName])
	pkgList[pkgInfo.Name] = pkgInfo

	return pkgInfo.OldName, pkgInfo.Name
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
