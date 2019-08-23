package util_confd

import (
	"reflect"
	"strings"
)

const (
	SPACENUM  = 2
	TYPE_YAML = "yaml"
	TYPE_ENV  = "env"
)

type ConfigItem struct {
	Name    string       // 配置名
	Value   string       // 配置值
	Comment string       // 备注
	Kind    reflect.Kind // 配置字段类型
	Set     ConfigSet    // 子配置项
}

type ConfigSet []*ConfigItem

func (s ConfigSet) Len() int      { return len(s) }
func (s ConfigSet) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ConfigSet) Less(i, j int) bool {
	if s[i].Kind != reflect.Struct && s[i].Kind != reflect.Slice {
		return true
	}
	if s[i].Kind == reflect.Struct && s[j].Kind == reflect.Slice {
		return true
	}
	return false
}

// 生成配置文本
func (s ConfigSet) CreateConfigStr(n int) string {

	var str, spaceStr string
	for {
		if n == len(spaceStr) {
			break
		}
		spaceStr += " "
	}

	for k, item := range s {
		switch item.Kind {
		case reflect.Struct:
			if k != 0 {
				str += "\n"
			}
			if item.Comment != "" {
				comments := strings.Split(item.Comment, "\n")
				for _, comment := range comments {
					str += spaceStr + "// " + comment + "\n"
				}
			}
			str += spaceStr + "[" + item.Name + "]\n"
			str += item.Set.CreateConfigStr(n + SPACENUM)
			str += spaceStr + "[" + item.Name + " END]\n"
			break
		case reflect.Slice:
			if k != 0 {
				str += "\n"
			}
			if item.Comment != "" {
				comments := strings.Split(item.Comment, "\n")
				for _, comment := range comments {
					str += spaceStr + "// " + comment + "\n"
				}
			}
			str += spaceStr + "[" + item.Name + "]\n"
			str += item.Set.CreateConfigStr(n + SPACENUM)
			str += spaceStr + "[" + item.Name + " END]\n"
			break
		default:
			str += spaceStr + item.Name + "=" + item.Value
			if item.Comment != "" {
				str += " // " + strings.Replace(item.Comment, "\n", " ", -1)
			}
			str += "\n"
			break
		}
	}
	return str
}