package main

import (
	"net/http"

	"github.com/austin-rey/go-first-web-service/product"
)


const apiBasePath = "/api"

// Creates HTTP server, registers handlers to pattern
func main() {
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}