package models

import (
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
	"log"
	"main/pkg/setting"
)

var db *sql.DB

// Setup 初始化数据库
func Setup() {
	var err error
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Shanghai",
			setting.DatabaseSetting.Host,
			setting.DatabaseSetting.Port,
			setting.DatabaseSetting.User,
			setting.DatabaseSetting.Name,
			setting.DatabaseSetting.Password,
		)
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("数据库配置错误: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	CreateTable([]interface{}{
		&Role{},
		&Account{},
		&Questionnaire{},
		&Question{},
		&Option{},
		&File{},
	}, true)
}
