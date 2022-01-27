package db

import "strings"

type (
	Table struct {
		Name         string
		SourceName   string
		Field        map[string]TableField
		FieldOrder   []string
		MaxFieldLen  int
		MaxGoTypeLen int
		IsHaveTime   bool
	}

	TableField struct {
		Name          string
		SourceName    string
		GoType        string
		DBType        string
		Default       string
		IsNotNull     bool
		IsIndex       bool
		IsUniqueIndex bool
		IsPrimaryKey  bool
		IndexName     string
	}
)

var FieldTypeMap = map[string]string{
	"int":        "int",
	"bigint":     "int64",
	"smallint":   "int",
	"mediumint":  "int",
	"tinyint":    "int",
	"float":      "float32",
	"double":     "float64",
	"decimal":    "float64",
	"date":       "time.Time",
	"time":       "time.Time",
	"year":       "time.Time",
	"datetime":   "time.Time",
	"timestamp":  "time.Time",
	"char":       "string",
	"varchar":    "string",
	"tinyblob":   "string",
	"tinytext":   "string",
	"blob":       "string",
	"test":       "string",
	"mediumblob": "string",
	"mediumtext": "string",
	"longblob":   "string",
	"longtext":   "string",
}

func ParseTableCreateSql(sql string) Table {
	table := Table{
		Field:        make(map[string]TableField),
		FieldOrder:   make([]string, 0),
		MaxFieldLen:  0,
		MaxGoTypeLen: 0,
		IsHaveTime:   false,
	}

	sqlSli := strings.Split(sql, "\n")
	for _, item := range sqlSli {
		if strings.HasPrefix(item, "CREATE") {
			table.SourceName, table.Name = ParseTableName(item)
			continue
		}
		if strings.HasPrefix(item, "  `") {
			field := ParseTableField(item)
			table.Field[field.SourceName] = field
			table.FieldOrder = append(table.FieldOrder, field.SourceName)
			if len(field.Name) > table.MaxFieldLen {
				table.MaxFieldLen = len(field.Name)
			}
			if len(field.GoType) > table.MaxGoTypeLen {
				table.MaxGoTypeLen = len(field.GoType)
			}
			if field.GoType == "time.Time" {
				table.IsHaveTime = true
			}
		}
		if strings.HasPrefix(item, "  PRIMARY KEY") {
			fieldNames := ParsePrimaryKey(item)
			for _, fieldName := range fieldNames {
				field := table.Field[fieldName]
				field.IsPrimaryKey = true
				table.Field[fieldName] = field
			}
		}
		if strings.HasPrefix(item, "  UNIQUE KEY") {
			fieldNames, indexName := ParseUniqueIndex(item)
			for _, fieldName := range fieldNames {
				field := table.Field[fieldName]
				field.IsUniqueIndex = true
				field.IndexName = indexName
				table.Field[fieldName] = field
			}
		}
		if strings.HasPrefix(item, "  KEY") {
			fieldNames, indexName := ParseIndex(item)
			for _, fieldName := range fieldNames {
				field := table.Field[fieldName]
				field.IsUniqueIndex = true
				field.IndexName = indexName
				table.Field[fieldName] = field
			}
		}
	}

	return table
}

func ParsePrimaryKey(str string) (fieldName []string) {
	fieldName = make([]string, 0)
	fieldSli := strings.Split(str, "(")
	fieldSli = strings.Split(fieldSli[1], ")")
	fieldSli = strings.Split(fieldSli[0], ",")
	for _, item := range fieldSli {
		fieldName = append(fieldName, strings.Trim(item, "`"))
	}

	return fieldName
}

func ParseIndex(str string) (fieldName []string, indexName string) {
	fieldName = make([]string, 0)
	fieldSli := strings.Split(str, "(")
	fieldSli = strings.Split(fieldSli[1], ")")
	fieldSli = strings.Split(fieldSli[0], ",")
	for _, item := range fieldSli {
		fieldName = append(fieldName, strings.Trim(item, "`"))
	}
	indexSli := strings.Split(str, " ")
	indexName = strings.Trim(indexSli[3], "`")

	return fieldName, indexName
}

func ParseUniqueIndex(str string) (fieldName []string, indexName string) {
	fieldName = make([]string, 0)
	fieldSli := strings.Split(str, "(")
	fieldSli = strings.Split(fieldSli[1], ")")
	fieldSli = strings.Split(fieldSli[0], ",")
	for _, item := range fieldSli {
		fieldName = append(fieldName, strings.Trim(item, "`"))
	}
	indexSli := strings.Split(str, " ")
	indexName = strings.Trim(indexSli[4], "`")

	return fieldName, indexName
}

func ParseTableName(str string) (sourceName string, name string) {
	strSli := strings.Split(str, "`")

	sourceName = strSli[1]
	name = SnakeToUpperCamelCase(sourceName)

	return sourceName, name
}

func ParseTableField(str string) (field TableField) {
	strSli := strings.Split(str, "`")
	field.SourceName = strSli[1]
	field.Name = SnakeToUpperCamelCase(field.SourceName)

	optionsStr := strings.Trim(strSli[2], " ")
	options := strings.Split(optionsStr, " ")
	field.DBType = options[0]
	field.GoType = FieldTypeMap[options[0]]
	if strings.Contains(options[0], "(") {
		field.GoType = FieldTypeMap[strings.Split(options[0], "(")[0]]
	}
	field.Default = GetFieldDefaultVal(options)
	field.IsNotNull = strings.Contains(optionsStr, "NOT NULL")

	return field
}

func GetFieldDefaultVal(fieldOptions []string) (result string) {
	isDefaultVal := false
	for _, option := range fieldOptions {
		if isDefaultVal {
			result = strings.Trim(option, ",")
			break
		}
		if option == "DEFAULT" {
			isDefaultVal = true
		}
	}

	return result
}

func SnakeToUpperCamelCase(str string) string {
	strSli := strings.Split(str, "_")
	for i, item := range strSli {
		res := strFirstToUpper(item)
		strSli[i] = res
	}

	return strings.Join(strSli, "")
}

func strFirstToUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
