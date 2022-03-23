package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/austin-rey/go-first-web-service/cors"
)

const productsBasePath = "products"

func SetupRoutes(apiBasePath string) {
	handleProducts := http.HandlerFunc(productsHandler)
	handleProduct := http.HandlerFunc(productHandler)

	http.Handle(
		fmt.Sprintf("%s/%s", apiBasePath, productsBasePath),
		cors.Middleware(handleProducts),
	)

	http.Handle(
		fmt.Sprintf("%s/%s/", apiBasePath, productsBasePath), 
		cors.Middleware(handleProduct),
	)
}

// Handler Function -	/products/
// GET - 				Fetch product by ID
// PUT - 				Update product by ID
func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "products/")
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments) - 1])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	product := getProduct(productID)
	if(product == nil) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
		case http.MethodGet:
			productJson, err := json.Marshal(product)
				
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write(productJson)

		case http.MethodPut:
			var updatedProduct Product
			bodyBytes, err := ioutil.ReadAll(r.Body)

			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			err = json.Unmarshal(bodyBytes, &updatedProduct )
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			if updatedProduct.ProductID != productID {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			addOrUpdateProduct(updatedProduct)
			w.WriteHeader(http.StatusAccepted)
			return

		case http.MethodOptions:
			return

		case http.MethodDelete:
			removeProduct(productID)
			
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Handler Function -   /products
// GET - 				Fetch all products
// POST - 				Create product
func productsHandler(w http.ResponseWriter, r *http.Request){
	switch r.Method {
		case http.MethodGet:
			productList := getProductList()
			productsJson, err := json.Marshal(productList)
			
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(productsJson)

		case http.MethodPost:
			var newProduct Product
			bodyBytes, err := ioutil.ReadAll(r.Body)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err = json.Unmarshal(bodyBytes, &newProduct )
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if newProduct.ProductID != 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			
			_, err = addOrUpdateProduct(newProduct)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusCreated)
			return
		case http.MethodOptions:
			return
	}
}