package main

import (
	"log"
	"net/http"

	"github.com/austin-rey/go-first-web-service/database"
	"github.com/austin-rey/go-first-web-service/product"
	receipts "github.com/austin-rey/go-first-web-service/receipt"
	_ "github.com/go-sql-driver/mysql"
)

const basePath = "/api"

func main() {
	database.SetupDatabase()
	product.SetupRoutes(basePath)
	receipts.SetupRoutes(basePath)
	log.Fatal(http.ListenAndServe(":5000", nil))
}