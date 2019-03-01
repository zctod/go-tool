package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"math/rand"
	"os"
	"reflect"
	"regexp"
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
	existSts := true
	if err != nil {
		existSts = false
	}
	if !existSts {
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

	dir := strings.Join(pathArr[:pathLen-1], "/")
	if err := PathCreate(dir); err != nil {
		return nil, err
	}

	return os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
}

// 兼容原始方法
func JsonToMap(s interface{}) (data map[string]interface{}) {
	data, _ = StrcutToMap(s)
	return
}

// 结构体转map
func StrcutToMap(s interface{}) (map[string]interface{}, error) {
	jsonStr, err := json.Marshal(&s)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`[^\{]*(\{.*\})[^\}]*`)
	jsonStr = []byte(re.ReplaceAllString(string(jsonStr), "$1"))

	var data map[string]interface{}
	err = json.Unmarshal(jsonStr, &data)
	return data, err
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

// map转xml
func MapToXml(data XmlMap) ([]byte, error) {

	return xml.Marshal(XmlMap(data))
}

// xml转map
func XmlToMap(b []byte) (XmlMap, error) {

	var mp = make(XmlMap)
	err := xml.Unmarshal(b, (*XmlMap)(&mp))
	return mp, err
}