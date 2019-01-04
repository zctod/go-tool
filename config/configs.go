// 配置包
package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
	"tool/common"
)

// 写入配置项
func writeConfig(config interface{}, path string) {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	t := reflect.TypeOf(config).Elem()
	for i := 0; i < t.NumField(); i++ {
		var defValue string
		field := tool.CamelToUnderline(t.Field(i).Name)
		tag := t.Field(i).Tag.Get("config")
		if tag != "" {
			tagArr := strings.Split(tag, ";")
			for _, v := range tagArr {
				if v != "" {
					sts, err := regexp.MatchString("^default:", v)
					if err != nil {
						fmt.Println(err)
					}
					if sts == true {
						defVauleArr := strings.Split(v, "default:")
						defValue = defVauleArr[1]
						break
					}
				}
			}
		}
		_, err := file.WriteString(field + "=" + defValue + "\n")
		if err != nil {
			fmt.Println(err)
		}
	}
}

// 初始化配置
func InitConfig(config interface{}, path string) {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			writeConfig(config, path)
			panic(errors.New("env is not exist, automatically generated and filled in"))
		} else {
			panic(err)
		}
	}

	c := reflect.ValueOf(config).Elem()

	arr := strings.Split(string(body), "\n")
	for _, v := range arr {
		v = strings.Replace(v, " ", "", -1)
		strArr := strings.Split(v, "=")
		if len(strArr) != 2 {
			continue
		}
		value := reflect.ValueOf(strArr[1])
		field := tool.UnderlineToCamel(strArr[0])
		c.FieldByName(field).Set(value)
	}
}
