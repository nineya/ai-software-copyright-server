package initialize

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"path/filepath"
	"tool-server/internal/global"
	"tool-server/internal/utils/jwt"
)

func readCommandLineArgs() string {
	pflag.StringP("workDir", "d", "./", "The directory where the working data is stored")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	path, _ := filepath.Abs(viper.GetString("workDir"))
	return path
}

func InitSystemConfig() {
	viperConfig := viper.New()
	workDir := readCommandLineArgs()
	// 设置配置文件名，没有后缀
	viperConfig.SetConfigName("application")
	// 设置读取文件格式为: yaml
	viperConfig.SetConfigType("yaml")
	// 设置配置文件目录(可以设置多个,优先级根据添加顺序来)
	viperConfig.AddConfigPath(workDir)
	if err := viperConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("Configuration file not found: %v \n", err))
		} else {
			panic(fmt.Errorf("Configuration file parsing error: %v \n", err))
		}
	}
	if err := viperConfig.Unmarshal(&global.CONFIG); err != nil {
		panic(fmt.Errorf("Description Failed to configure item mapping: %v \n", err))
	}
	global.JWT = jwt.NewJWT(global.CONFIG.JWT.SigningKey)
	global.WORK_DIR = workDir
}
