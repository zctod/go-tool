package logs

import (
	"fmt"
	"github.com/zctod/tool/common/utils"
	"log"
	"os"
	"strings"
	"time"
)

var (
	Normal *log.Logger
	Info   *log.Logger
	Warn   *log.Logger
	Error  *log.Logger
)

func init() {
	var dirPath = "log/"
	var path = time.Now().Format("20060102") + ".log"
	Normal = CreateLog(dirPath + path)
	Info = CreateLog(dirPath + "info_" + path)
	Warn = CreateLog(dirPath + "warn_" + path)
	Error = CreateLog(dirPath + "error_" + path)
}

func CreateLog(path string) *log.Logger {
	var pathArr = strings.Split(path, "/")
	var pathLen = len(pathArr)

	var file *os.File
	var err error
	for i := 0; i < pathLen; i++ {
		if pathArr[i] != "" {
			if i == pathLen - 1 {
				file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
				if nil != err {
					fmt.Println(err)
				}
			} else {
				if err = utils.PathCreate(pathArr[i]); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	return log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}