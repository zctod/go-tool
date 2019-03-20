package util_confd

import (
	"errors"
	"fmt"
	"github.com/zctod/tool/common/utils"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// 初始化配置
func InitConfig(config interface{}, path string) error {

	var err error
	body, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = WriteConfig(config, path)
			if err == nil {
				return InitConfig(config, path)
			}
		}
		return err
	}

	t := reflect.TypeOf(config).Elem()
	c := reflect.ValueOf(config).Elem()

	// 先根据\n截取字符串数组
	var arr = strings.Split(string(body), "\n")
	m, err := arrHandle(arr)
	if err != nil {
		return err
	}
	return SetVal(t, c, m)
}

// 读取配置信息
func ReadConfig(config interface{}, path string) error {
	var err error
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(config).Elem()
	c := reflect.ValueOf(config).Elem()

	// 先根据\n截取字符串数组
	var arr = strings.Split(string(body), "\n")
	m, err := arrHandle(arr)
	if err != nil {
		return err
	}
	fmt.Println(m)
	return SetVal(t, c, m)
}

// 自动写配置
func WriteConfig(config interface{}, path string) error {

	t := reflect.TypeOf(config).Elem()
	set, err := Parse(t)
	if err != nil {
		return err
	}
	return set.Write(path)
}

const SPACENUM = 2

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

func (s ConfigSet) Write(path string) error {

	// 自动生成文件
	file, err := utils.CreateFile(path)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.WriteString(s.CreateConfigStr(0))
	return err
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
			fmt.Println(item.Set[0])
			if k != 0 {
				str += "\n"
			}
			if item.Comment != "" {
				comments := strings.Split(item.Comment, "\n")
				for _, comment := range comments {
					str += spaceStr + "// " + comment + "\n"
				}
			}
			str += spaceStr + "-[" + item.Name + "]\n"
			str += item.Set.CreateConfigStr(n + SPACENUM)
			str += spaceStr + "-[" + item.Name + " END]\n"
			break
		default:
			str += spaceStr + item.Name + "=" + item.Value
			if item.Comment != "" {
				str += " // " + strings.ReplaceAll(item.Comment, "\n", " ")
			}
			str += "\n"
			break
		}
	}
	return str
}

// 设置配置值
func SetVal(t reflect.Type, c reflect.Value, m map[string]interface{}) (err error) {

	for k, v := range m {
		if f, ok := t.FieldByName(k); ok {
			cv := c.FieldByName(k)
			switch reflect.TypeOf(v).Kind() {
			case reflect.String:
				err := setValue(f, cv, v.(string))
				if err != nil {
					goto quit
				}
				break
			case reflect.Map:
				if fs, ok := t.FieldByName(k); ok {
					err = SetVal(fs.Type, cv, v.(map[string]interface{}))
					if err != nil {
						goto quit
					}
				}
				break
			case reflect.Slice:
				elType := cv.Type().Elem()
				slice := reflect.MakeSlice(cv.Type(), cv.Len(), cv.Cap())
				for _, vm := range v.([]map[string]interface{}) {
					el := reflect.New(elType).Elem()
					err := SetVal(elType, el, vm)
					if err != nil {
						goto quit
					}
					slice = reflect.Append(slice, el)
				}
				cv.Set(slice)
				break
			}
		}
	}
quit:
	m = nil
	return
}

// 按类型设置值
func setValue(f reflect.StructField, v reflect.Value, val string) error {

	var err error
	switch f.Type.Kind() {
	case reflect.String:
		v.SetString(val)
		break
	case reflect.Bool:
		if val == "true" || val == "1" {
			v.SetBool(true)
		} else {
			v.SetBool(false)
		}
		break
	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(val, 64)
		if err == nil {
			v.SetFloat(value)
			return nil
		}
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.Atoi(val)
		if err == nil {
			v.SetInt(int64(value))
			return nil
		}
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.ParseUint(val, 10, 64)
		if err == nil {
			v.SetUint(value)
			return nil
		}
		break
	default:
		return errors.New("field type is not supported")
	}
	return err
}

// 处理单层配置
func arrHandle(arr []string) (map[string]interface{}, error) {

	var err error
	var m = make(map[string]interface{})
	var index = -1
	regObjStart, _ := regexp.Compile(`^\s*\[[\s\S]*]\s*$`)
	regObjPre, _ := regexp.Compile(`^\s*\[`)
	regPost, _ := regexp.Compile(`]\s*$`)
	regArrStart, _ := regexp.Compile(`^\s*-\[[\s\S]*]\s*$`)
	regArrPre, _ := regexp.Compile(`^\s*-\[`)
	regComment, _ := regexp.Compile(`\s+//[\s\S]*$|^//[\s\S]*$`)

	for k, v := range arr {
		// 跳过被截取的部分
		if index >= k {
			continue
		}
		// 结构体解析
		if regObjStart.MatchString(v) {
			v = regObjPre.ReplaceAllString(v, "")
			v = regPost.ReplaceAllString(v, "")
			regEnd, _ := regexp.Compile(`^\s*\[` + v + ` END]\s*$`)
			index = 0
			for kk, vv := range arr {
				if regEnd.MatchString(vv) {
					index = kk
					break
				}
			}
			if k >= index {
				err = errors.New("env format error")
				break
			}
			mc, err := arrHandle(arr[k+1 : index])
			if err != nil {
				break
			}
			field := utils.UnderlineToCamel(v)
			m[field] = mc
			continue
		}
		// 数组解析
		if regArrStart.MatchString(v) {
			v = regArrPre.ReplaceAllString(v, "")
			v = regPost.ReplaceAllString(v, "")
			regEnd, _ := regexp.Compile(`^\s*-\[` + v + ` END]\s*$`)
			index = 0
			for kk, vv := range arr {
				if regEnd.MatchString(vv) {
					index = kk
					if k >= index {
						continue
					}
					break
				}
			}
			if k >= index {
				err = errors.New("env format error")
				break
			}
			mc, err := arrHandle(arr[k+1 : index])
			if err != nil {
				break
			}
			field := utils.UnderlineToCamel(v)
			if m[field] == nil {
				m[field] = make([]map[string]interface{}, 0)
			}
			m[field] = append(m[field].([]map[string]interface{}), mc)
			continue
		}
		// 过滤备注
		v = regComment.ReplaceAllString(v, "")
		// 普通配置截取
		v = strings.Replace(v, " ", "", -1)
		strArr := strings.Split(v, "=")
		if len(strArr) != 2 {
			continue
		}
		field := utils.UnderlineToCamel(strArr[0])
		m[field] = strArr[1]
	}

	return m, err
}

// 解析配置
func Parse(t reflect.Type) (ConfigSet, error) {

	if t.Kind() != reflect.Struct {
		return nil, errors.New("config must be a struct")
	}
	var set = make(ConfigSet, 0)
	var err error
	for i := 0; i < t.NumField(); i++ {

		field := t.Field(i)
		var item *ConfigItem

		switch field.Type.Kind() {
		case reflect.String, reflect.Bool, reflect.Float32, reflect.Float64,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

			item, err = parseValue(field)
			if err != nil {
				goto quit
			}
			break
		case reflect.Struct:
			item, err = parseValue(field)
			if err != nil {
				goto quit
			}
			itemSet, err := Parse(field.Type)
			if err != nil {
				goto quit
			}
			item.Set = itemSet
			break
		case reflect.Slice:
			item, err = parseValue(field)
			if err != nil {
				goto quit
			}
			if field.Type.Elem().Kind() == reflect.Slice {
				err = errors.New("slice is support single layer only")
				goto quit
			}
			itemSet, err := Parse(field.Type.Elem())
			if err != nil {
				goto quit
			}
			item.Set = itemSet
			break
		//case reflect.Map:
		//	item, err = ParseValue(field)
		//	break
		default:
			err = errors.New("field type is not supported")
			goto quit
		}
		set = append(set, item)
	}

quit:
	sort.Sort(set)
	return set, err
}

// 解析单个配置参数
func parseValue(f reflect.StructField) (*ConfigItem, error) {

	var item = &ConfigItem{
		Name: utils.CamelToUnderline(f.Name),
		Kind: f.Type.Kind(),
	}
	var err error

	tag := f.Tag.Get("config")
	if tag != "" {
		tagArr := strings.Split(tag, ";")
		for _, v := range tagArr {
			if v != "" {
				sts, err := regexp.MatchString(`^\s*default:`, v)
				if err != nil {
					goto quit
				}
				if sts == true {
					defValArr := strings.Split(v, "default:")
					item.Value = defValArr[1]
					continue
				}
				sts, err = regexp.MatchString(`^\s*comment:`, v)
				if err != nil {
					goto quit
				}
				if sts == true {
					ctValArr := strings.Split(v, "comment:")
					item.Comment = ctValArr[1]
					continue
				}
			}
		}
	}

quit:
	return item, err
}
