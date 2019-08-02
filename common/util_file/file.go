package util_file

import (
	"net/http"
)

// 获取文件类型
func GetFileContentType(fileByte []byte) string {
	buffer := make([]byte, 512)
	buffer = append(buffer, fileByte[:512]...)
	return http.DetectContentType(fileByte)
}