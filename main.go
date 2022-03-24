package main

import (
	"net/http"

	"github.com/austin-rey/go-first-web-service/database"
	"github.com/austin-rey/go-first-web-service/product"
	_ "github.com/go-sql-driver/mysql"
)


const apiBasePath = "/api"

// Creates HTTP server, registers handlers to pattern
func main() {
	database.SetupDatabase()
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}