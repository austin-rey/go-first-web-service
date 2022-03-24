package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

// Read/Write mutex of product map
// We wrap a mutex around the map to handle rw controls
var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func init() {
	fmt.Println("Loading products...")
	prodMap, err := loadProductMap()
	productMap.m = prodMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("%d products loaded...\n", len(productMap.m))
}

// Loads data from a file into a slice (index is identifier) and transform data into a map where productID is identifier
func loadProductMap() (map[int]Product, error) {
	fileName := "products.json"
	_, err := os.Stat(fileName)

	// Check if file exists
	if os.IsNotExist(err){
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}

	// Parse file for contents and deserialize products into slice
	file, _ := ioutil.ReadFile(fileName)
	productList := make([]Product,0)
	err = json.Unmarshal([]byte(file), &productList)
	if err != nil {
		log.Fatal(err)
	}

	// Init product map and add products from map to it
	prodMap := make(map[int]Product)
	for i := 0; i < len(productList); i++ {
		prodMap[productList[i].ProductID] = productList[i]
	}

	return prodMap, nil
}

// Return product by ID
func getProduct(productID int) *Product {
	// Protects the productMap struct from being accessed by other thread operations
	productMap.RLock()
	defer productMap.RUnlock()

	if product, ok := productMap.m[productID]; ok {
		return &product
	}

	return nil
}

// Remove product by ID
func removeProduct(productID int) {
	productMap.Lock()
	defer productMap.Unlock()
	delete(productMap.m, productID)
}

// Returns slice of products
func getProductList() []Product {
	productMap.RLock()
	products := make([]Product,0,len(productMap.m))
	for _,value := range productMap.m {
		products = append (products, value)
	}
	productMap.RUnlock()
	return products
}

// Returns sorted int slice of product Ids
func getProductIds() []int {
	productMap.RLock()
	productIds := []int{}
	for key := range productMap.m {
		productIds = append(productIds, key)
	}
	productMap.RUnlock()
	sort.Ints(productIds)
	return productIds
}

// Helper Function - Returns highest product ID + 1
func getNextProductID() int {
	productIDs := getProductIds()
	return productIDs[len(productIDs)-1]+1
}

// Returns a product ID and error (if applicable) when a product is created or updated
func addOrUpdateProduct(product Product) (int, error) {
	// If the product id is set, update, otherwise add
	addOrUpdateID := -1

	if product.ProductID > 0 {
		// Update
		oldProduct := getProduct(product.ProductID)
		// Check of product exists
		if oldProduct == nil {
			return 0, fmt.Errorf("product id [%d] does not exist", product.ProductID)
		}
		addOrUpdateID = product.ProductID
	} else {
		// Add
		addOrUpdateID = getNextProductID()
		product.ProductID = addOrUpdateID
	}
	productMap.Lock()
	productMap.m[addOrUpdateID] = product
	productMap.Unlock()
	return addOrUpdateID, nil
}