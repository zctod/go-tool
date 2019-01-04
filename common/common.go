// 小工具包
package common

import (
	"encoding/json"
	"math/rand"
	"os"
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

// json结果转map
func JsonToMap(s interface{}) (data map[string]interface{}) {
	jsonStr, e := json.Marshal(&s)
	if e != nil {
		panic(e)
	}
	re := regexp.MustCompile(`[^\{]*(\{.*\})[^\}]*`)
	jsonStr = []byte(re.ReplaceAllString(string(jsonStr), "$1"))

	_ = json.Unmarshal(jsonStr, &data)
	return
}

// 字符串驼峰转下划线
func CamelToUnderline(s string) string {
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
		if k && d == '_' && num > i && s[i + 1] >= 'a' && s[i + 1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}