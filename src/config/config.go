package config

import (
	"os"
	"path"

	"github.com/Kor-SVS/cocoa/src/log"
	"github.com/Kor-SVS/cocoa/src/util"
	"github.com/spf13/viper"
)

var (
	RootConfig *Config
	logger     *log.Logger
)

func init() {
	RootConfig = NewConfig(viper.GetViper())

	logOption := log.NewLoggerOption()
	logOption.Prefix = "[config]"
	logWriter := log.NewLogWriter(nil, nil, nil, nil)
	logger = log.NewLogger(logOption, logWriter)

	logger.Trace("Config init...")

	currentPath := util.GetExecutablePath()
	configPath := path.Join(currentPath, "config.yaml")

	if !util.FileExists(configPath) {
		os.Create(configPath)
	}

	RootConfig.v.SetConfigName("config")
	RootConfig.v.SetConfigType("yaml")
	RootConfig.v.AddConfigPath(currentPath)

	ReadConfig()
}

func ReadConfig() {
	if err := RootConfig.v.ReadInConfig(); err != nil {
		logger.Warningf("설정 파일 로드 실패 (err=%v)", err)
	}
}

func WriteConfig() {
	err := RootConfig.v.WriteConfig()
	if err != nil {
		logger.Errorf("설정 파일 저장 실패 (err=%v)", err)
	}
}
