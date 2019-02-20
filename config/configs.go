// 配置包
package config

import (
	"github.com/zctod/tool/common/util_confd"
)

// 写入配置项
// Deprecated: please use package common/util_confd.
func writeConfig(config interface{}, path string) error {

	return util_confd.WriteConfig(config, path)
}

// 初始化配置
// Deprecated: please use package common/util_confd.
func InitConfig(config interface{}, path string) error {

	return util_confd.InitConfig(config, path)
}
