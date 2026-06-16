package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	user := getEnv("DB_USER", "crrc_user")
	password := getEnv("DB_PASSWORD", "123456")
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnvInt("DB_PORT", 3306)
	dbname := getEnv("DB_NAME", "crrc_pm")

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

	for i := 1; i <= 30; i++ {
		if err = DB.Ping(); err == nil {
			log.Println("✅ 数据库连接成功！")
			return
		}
		log.Printf("数据库暂不可访问，正在重试(%d/30): %v", i, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("数据库无法访问: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return result
}
