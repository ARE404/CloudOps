package di

import (
	"log"

	"github.com/spf13/viper"
)

// InitConfig loads configuration from file
func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("警告: 配置文件加载失败: %v，使用默认值", err)
	}
}
