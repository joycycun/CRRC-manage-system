package main

import (
	"crrc_pm_backend/config"
	"crrc_pm_backend/router"
	"log"
	"net/http"
)

func main() {
	config.InitDB()

	r := router.InitRouter()

	log.Println("后端启动成功：http://localhost:8080")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("后端启动失败: %v", err)
	}
}