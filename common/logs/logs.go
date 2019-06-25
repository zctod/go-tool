package logs

import (
	"fmt"
	"github.com/zctod/go-tool/common/utils"
	"log"
	"os"
	"strings"
)

func CreateLog(path string) *log.Logger {
	
	return log.New(ReadPath(path), "", log.Ldate|log.Ltime|log.Lshortfile)
}

func ReadPath(path string) *os.File {

	var pathArr = strings.Split(path, "/")
	var pathLen = len(pathArr)

	dir := strings.Join(pathArr[:pathLen-1], "/")

	if err := utils.PathCreate(dir); err != nil {
		fmt.Println(err)
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if nil != err {
		fmt.Println(err)
	}
	return file
}
