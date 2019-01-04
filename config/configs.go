package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
)

// writeConfig 写入配置项
func writeConfig(config interface{}, path string) {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	t := reflect.TypeOf(config).Elem()
	for i := 0; i < t.NumField(); i++ {
		var defValue string
		field := camelToUnderline(t.Field(i).Name)
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

// InitConfig 初始化配置
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
		field := underlineToCamel(strArr[0])
		c.FieldByName(field).Set(value)
	}
}

// camelToUnderline 字符串驼峰转下划线
func camelToUnderline(s string) string {
	num := len(s)
	data := make([]byte, 0, num * 2)
	j := false
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// underlineToCamel 字符串下划线转驼峰
func underlineToCamel(s string) string {
	data := make([]byte, 0, len(s))
	j, k := false, false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j, k = false, true
		}
		if k && d == '_' && num > i && s[i + 1] >= 'a' && s[i + 1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
