// tool 小工具包
package tool

import (
	"encoding/json"
	"math/rand"
	"os"
	"regexp"
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