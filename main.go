// @title           CloudOps API
// @version         1.0
// @description     CloudOps — K8s多集群管理与AI智能运维平台
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"

	"github.com/GoSimplicity/CloudOps/pkg/di"
	"github.com/spf13/viper"
)

func main() {
	// Load config first (before wire providers run)
	di.InitConfig()

	engine, bootstrap, err := di.InitApp()
	if err != nil {
		log.Fatalf("应用初始化失败: %v", err)
	}

	// Run startup tasks
	bootstrap.Start()

	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	log.Printf("CloudOps 启动，监听 :%s", port)
	if err := engine.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
