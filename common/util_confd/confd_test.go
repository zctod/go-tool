package util_confd

import (
	"fmt"
	"testing"
)

type Configure struct {
	Port        int    `config:"default:8080;comment:监听端口"`
	Mode        string `config:"default:debug;comment:运行模式debug调试、test测试、release生产"`
	ModelPlates []ItemPlates
	Normal      ItemNormal
	Arr        []float64
}

type ItemNormal struct {
	MaxNum        int
	LevelDiff     float64
	SizeUnCompare bool
}

type ItemPlates struct {
	Num    int
	Str    string
	Amount float64
	Id     uint
}

const PATH = "../../testdata/env"
var cfg = &Configure{}

func TestInitConfig(t *testing.T) {

	err := InitConfig(cfg, PATH)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)
}

func TestReadConfig(t *testing.T) {

	err := ReadConfig(cfg, PATH)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)
}

func TestWriteConfig(t *testing.T) {

	err := WriteConfig(cfg, PATH)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)
}