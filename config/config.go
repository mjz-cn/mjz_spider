package config

import (
	"io/ioutil"
	"path"
	"runtime"

	"gopkg.in/yaml.v2"
)

type Config struct {
	MysqlConn string `yaml:"mysql"`
	Tesseract string `yaml:"tesseract"`
	ExecuteDir string `yaml:"executeDir"`
}

var GlobalConfig *Config

func init() {
	_, curFilename, _, _ := runtime.Caller(1)
	configDir := path.Dir(curFilename) // 读取当前文件所在路径

	filename := configDir + "/config.yml"

	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)  // 必须有配置文件
	}

	err = yaml.Unmarshal(yamlFile, &GlobalConfig)
	if err != nil {
		panic(err)  // 配置文件必须读取成功
	}

	if GlobalConfig.ExecuteDir == "" {
	 	GlobalConfig.ExecuteDir = path.Dir(configDir)
	}
}