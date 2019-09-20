package util_confd

import (
	"fmt"
	"testing"
)

/**
 * Created by zc on 2019-09-06.
 */

const YAMLPATH = "../../testdata/config.yml"

func TestInitYamlConfig(t *testing.T) {

	err := InitConfig(cfg, YAMLPATH)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)
}

func TestReadYamlConfig(t *testing.T) {

	err := ReadConfig(cfg, YAMLPATH)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)
}

func TestWriteYamlConfig(t *testing.T) {

	err := WriteConfig(cfg, YAMLPATH)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)
}