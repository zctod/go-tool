package util_confd

import (
	"github.com/zctod/go-tool/common/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// 初始化yaml配置
func InitYamlConfig(config interface{}, path string) error {

	var err error
	body, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = WriteConfig(config, path)
			if err == nil {
				return InitYamlConfig(config, path)
			}
		}
		return err
	}
	return yaml.Unmarshal(body, config)
}

// 读取yaml配置信息
func ReadYamlConfig(config interface{}, path string) error {

	var err error
	body, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(body, config)
}

// 自动写yaml配置
func WriteYamlConfig(config interface{}, path string) error {

	b, err :=  yaml.Marshal(config)
	if err != nil {
		return err
	}
	// 自动生成文件
	file, err := utils.CreateFile(path)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.Write(b)
	return err
}