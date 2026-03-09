package di

import (
	"log"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	addr := viper.GetString("mysql.addr")
	db, err := gorm.Open(mysql.Open(addr), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.K8sCluster{},
		&model.PromScrapePool{},
		&model.PromAlertRule{},
		&model.PromAlertEvent{},
		&model.PromSendGroup{},
		&model.WoTemplate{},
		&model.WoInstance{},
		&model.WoFlow{},
		&model.WoComment{},
		&model.AiKnowledgeChunk{},
	); err != nil {
		log.Printf("数据库迁移失败: %v", err)
	}

	return db
}
