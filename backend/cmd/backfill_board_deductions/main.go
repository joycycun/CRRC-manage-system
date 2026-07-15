package main

import (
	"crrc_pm_backend/config"
	"crrc_pm_backend/handler"
)

func main() {
	config.InitDB()
	handler.BackfillBoardCompositionDeductionsForFactoryInventory()
}
