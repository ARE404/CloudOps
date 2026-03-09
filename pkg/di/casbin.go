package di

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/spf13/viper"
)

func InitCasbin() *casbin.Enforcer {
	modelPath := viper.GetString("casbin.model_path")
	if modelPath == "" {
		modelPath = "./config/rbac_model.conf"
	}
	enforcer, err := casbin.NewEnforcer(modelPath)
	if err != nil {
		log.Printf("Casbin 初始化失败: %v，RBAC 功能将被禁用", err)
		return nil
	}
	return enforcer
}
