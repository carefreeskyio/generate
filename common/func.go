package common

import (
	"fmt"
	"os"
	"strings"
)

func GetFileContentToSlice(filePath string) (result []string, err error) {
	result = make([]string, 0)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("open file %v failed: err=%v \n", filePath, err)
		return result, err
	}

	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			fmt.Printf("file.Close failed: fileName=%v err=%v \n", file.Name(), err)
		}
	}(file)

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file %v stat failed: err=%v \n", filePath, err)
		return result, err
	}

	buffer := make([]byte, fileInfo.Size())
	if _, err = file.Read(buffer); err != nil {
		fmt.Printf("read file %v failed: err=%v \n", filePath, err)
		return result, err
	}

	content := strings.Split(string(buffer), "\n")

	for _, line := range content {
		if line != "" {
			result = append(result, line)
		}
	}

	return result, nil
}

func ParseFunc(line string) (result FuncInfo) {
	sliWithCloseParenthesis := strings.Split(line, "(")
	lineSlice := make([]string, 0)
	for _, item := range sliWithCloseParenthesis {
		if !strings.Contains(item, ")") {
			lineSlice = append(lineSlice, item)
			continue
		}
		itemSli := strings.Split(item, ")")
		for _, i := range itemSli {
			if i == " " {
				continue
			}
			lineSlice = append(lineSlice, i)
		}
	}
	params := strings.Split(lineSlice[3], ",")
	filteredParams := make([]string, 0)
	for _, param := range params {
		if strings.Contains(param, "context.Context") {
			continue
		}
		filteredParams = append(filteredParams, strings.Trim(param, " "))
	}

	result = FuncInfo{
		Name:   strings.Trim(lineSlice[2], " "),
		Params: strings.Join(filteredParams, ", "),
	}

	return result
}

func GetPassThroughParam(params string) (result string) {
	passThrough := make([]string, 0)
	paramsSli := strings.Split(params, ",")
	for _, param := range paramsSli {
		param = strings.Trim(param, " ")
		paramSli := strings.Split(param, " ")
		passThrough = append(passThrough, paramSli[0])
	}
	result = strings.Join(passThrough, ", ")

	return result
}
