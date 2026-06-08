package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	// MySQL 数据库连接配置
	user := "crrc_user"  // 数据库用户名
	password := "123456" // 数据库密码
	host := "127.0.0.1"  // 数据库地址
	port := 3306         // 数据库端口
	dbname := "crrc_pm"  // 数据库名

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 设置连接池
	DB.SetMaxOpenConns(50) // 最大打开连接数
	DB.SetMaxIdleConns(10) // 最大空闲连接数

	// 测试连接
	if err = DB.Ping(); err != nil {
		log.Fatalf("数据库无法访问: %v", err)
	}

	log.Println("✅ 数据库连接成功！")
}
