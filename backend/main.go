package main

import (
	"crrc_pm_backend/config"
	"crrc_pm_backend/handler"
	"crrc_pm_backend/router"
	"log"
	"net/http"
	"os"
)

func main() {
	config.InitDB()
	handler.BackfillBoardCompositionDeductionsForFactoryInventory()

	r := router.InitRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("后端启动成功：http://localhost:%s", port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("后端启动失败: %v", err)
	}
}
