package conf

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// mysql数据库到golang 语言映射
const DefaultDbTypeToLangMapperRule string = `#mysql 到golang 之间的映射
mysql-golang:
  TINYINT: int8
  SMALLINT: int16
  MEDIUMINT: int32
  INT: int32
  INTEGER: int32
  BIGINT: int64
  FLOAT: float32
  DOUBLE: float64
  DECIMAL: float32
  DATE: time.Time
  TIME: time.Time
  YEAR: time.Time
  DATETIME: time.Time
  TIMESTAMP: time.Time
  CHAR: string
  VARCHAR: string
  TINYTEXT: string
  TEXT: string
  MEDIUMTEXT: string
  LONGTEXT: string`

// map[数据库类型-目标语言][数据库数据类型][目标语言数据类型]
//
// map[mysql-go][varchar][]
var DataTypeTOLangMapperRule map[string]map[string]string = make(map[string]map[string]string, 0)

func init() {
	vp1 := viper.New()
	vp1.SetConfigType("yml")
	in := strings.NewReader(DefaultDbTypeToLangMapperRule)
	err := vp1.ReadConfig(in)
	if err != nil {
		fmt.Println("✖", err.Error())
		return
	}
	err = vp1.Unmarshal(&DataTypeTOLangMapperRule)
	if err != nil {
		fmt.Println("✖", err.Error())
		return
	}
	fmt.Println("✅Load DataTypeTOLangMapperRule success")
}

func ResetMapperRuleFromYaml(filepath string) (r map[string]map[string]string, e error) {
	m := make(map[string]map[string]string, 0)
	vp1 := viper.New()
	vp1.SetConfigType("yml")
	vp1.SetConfigFile(filepath)
	e = vp1.ReadInConfig()
	vp1.Unmarshal(&m)
	if e != nil {
		return
	}
	for sec1, map1 := range m {
		for datatype, langtype := range map1 {
			DataTypeTOLangMapperRule[sec1][datatype] = langtype
		}
	}
	return DataTypeTOLangMapperRule, nil
}

// 通过map重置规则
func ResetMapperRuleFromMap(m map[string]map[string]string) (r map[string]map[string]string, e error) {
	for sec1, map1 := range m {
		for datatype, langtype := range map1 {
			DataTypeTOLangMapperRule[sec1][datatype] = langtype
		}
	}
	return DataTypeTOLangMapperRule, nil
}

// 重置映射规则
//
// ResetMapperRule("./test.yml")或  ResetMapperRule(map[string]map[string]string)
//
// 返回最新规则
func ResetMapperRule(input any) (r map[string]map[string]string, e error) {
	switch input := input.(type) {
	case string:
		return ResetMapperRuleFromYaml(input)
	case map[string]map[string]string:
		return ResetMapperRuleFromMap(input)
	default:
		return DataTypeTOLangMapperRule, fmt.Errorf("not suport")
	}
}
