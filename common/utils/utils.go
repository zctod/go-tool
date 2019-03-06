package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 生成随机字符串
func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := make([]byte, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//自动创建文件夹
func PathCreate(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		mkErr := os.MkdirAll(path, os.ModePerm)
		if mkErr != nil {
			return mkErr
		}
	}
	return nil
}

//自动生成全路径文件
func CreateFile(path string) (*os.File, error) {

	var pathArr = strings.Split(path, "/")
	var pathLen = len(pathArr)
	if pathLen > 1 {
		dir := strings.Join(pathArr[:pathLen-1], "/")
		if err := PathCreate(dir); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
}

// json转map(兼容，弃用)
func JsonToMap(s interface{}) (data map[string]interface{}) {
	b, _ := json.Marshal(&s)
	re := regexp.MustCompile(`[^\{]*(\{.*\})[^\}]*`)
	jsonStr := []byte(re.ReplaceAllString(string(b), "$1"))

	_ = json.Unmarshal(jsonStr, &data)
	return data
}

// 结构体转map
func StrcutToMap(s interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(&s)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	return data, err
}

// map转结构体
// s map
// res 结构体指针地址
func MapToStruct(s interface{}, res interface{}) error {
	b, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, res)
}

// 结构体数组转map数组
func ArrayStructToMap(s interface{}) ([]map[string]interface{}, error) {
	jsonStr, err := json.Marshal(&s)
	if err != nil {
		return nil, err
	}
	var data []map[string]interface{}
	err = json.Unmarshal(jsonStr, &data)
	return data, err
}

// 字符串驼峰转下划线
func CamelToUnderline(s string) string {
	num := len(s)
	data := make([]byte, 0, num*2)
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

// 字符串下划线转驼峰
func UnderlineToCamel(s string) string {
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
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// 判断元素是否存在数组中
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

// md5加密并转string
func MD5(str string) string {
	var m = md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// 通用map转字符串map
func MapToStringMap(data map[string]interface{}) map[string]string {

	var res = make(map[string]string)
	for k, v := range data {
		var val string
		switch v.(type) {
		case int:
			val = strconv.Itoa(v.(int))
			break
		case int8:
			val = strconv.Itoa(int(v.(int8)))
			break
		case int16:
			val = strconv.Itoa(int(v.(int16)))
			break
		case int32:
			val = strconv.Itoa(int(v.(int32)))
			break
		case int64:
			val = strconv.Itoa(int(v.(int64)))
			break
		case float32:
			val = strconv.FormatFloat(float64(v.(float32)), 'f', -1, 64)
			break
		case float64:
			val = strconv.FormatFloat(v.(float64), 'f', -1, 64)
			break
		default:
			val = v.(string)
			break
		}
		res[k] = val
	}
	return res
}

// map转xml
func MapToXml(data map[string]string) []byte {

	var buf bytes.Buffer
	buf.WriteString(`<xml>`)
	for k, v := range data {
		buf.WriteString(`<`)
		buf.WriteString(k)
		buf.WriteString(`><![CDATA[`)
		buf.WriteString(v)
		buf.WriteString(`]]></`)
		buf.WriteString(k)
		buf.WriteString(`>`)
	}
	buf.WriteString(`</xml>`)

	return buf.Bytes()
}

// xml转map
func XmlToMap(b []byte) map[string]string {

	params := make(map[string]string)
	decoder := xml.NewDecoder(bytes.NewReader(b))

	var key, value string
	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement: // 开始标签
			key = token.Name.Local
		case xml.CharData: // 标签内容
			content := string([]byte(token))
			value = content
		}
		if key != "xml" {
			if value != "\n" {
				params[key] = value
			}
		}
	}
	return params
}