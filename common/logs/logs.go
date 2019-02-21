package logs

import (
	"fmt"
	"github.com/zctod/tool/common/utils"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

const (
	LOG_LEVEL_NORMAL = "normal"
	LOG_LEVEL_INFO   = "info"
	LOG_LEVEL_WARN   = "warn"
	LOG_LEVEL_ERROR  = "error"
)

var dirPath = "logs/"
var path = time.Now().Format("20060102") + ".log"

var w io.Writer
var (
	normalLog = log.New(w, "", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog   = log.New(w, "[INFO]", log.Ldate|log.Ltime|log.Lshortfile)
	warnLog   = log.New(w, "[WARN]", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog  = log.New(w, "[ERROR]", log.Ldate|log.Ltime|log.Lshortfile)
)
var (
	normalPath = dirPath + path
	infoPath   = dirPath + "info_" + path
	warnPath   = dirPath + "warn_" + path
	errorPath  = dirPath + "error_" + path

	normalPathOrigin = ""
	infoPathOrigin   = ""
	warnPathOrigin   = ""
	errorPathOrigin  = ""
)
var (
	Normal = SetPath(normalLog, normalPath, LOG_LEVEL_NORMAL)
	Info   = SetPath(infoLog, infoPath, LOG_LEVEL_INFO)
	Warn   = SetPath(warnLog, warnPath, LOG_LEVEL_WARN)
	Error  = SetPath(errorLog, errorPath, LOG_LEVEL_ERROR)
)

func SetPath(l *log.Logger, path string, levelType string) *log.Logger {
	switch levelType {
	case LOG_LEVEL_NORMAL:
		if normalPath == normalPathOrigin {
			return l
		}
		break;
	case LOG_LEVEL_INFO:
		if infoPath == infoPathOrigin {
			return l
		}
		break;
	case LOG_LEVEL_WARN:
		if warnPath == warnPathOrigin {
			return l
		}
		break;
	case LOG_LEVEL_ERROR:
		if errorPath == errorPathOrigin {
			return l
		}
		break;
	default:
		break;
	}
	l.SetOutput(ReadPath(path))
	return l
}

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
