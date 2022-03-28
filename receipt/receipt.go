package receipts

import (
	"io/ioutil"
	"path/filepath"
	"time"
)

var ReceiptDirectory string = filepath.Join("uploads")

type Receipt struct {
	ReceiptName string `json:"name"`
	UploadDate time.Time `json:"upload"`
}

// Constructs slice of receipts from files in the upload dir and returns it and any errors
func GetReceipts() ([]Receipt, error) {
	receipts := make([]Receipt, 0)
	files, err := ioutil.ReadDir(ReceiptDirectory)
	if err != nil {
		return nil, err
	}
	for _,f := range files {
		receipts = append(receipts, Receipt{ReceiptName: f.Name(), UploadDate: f.ModTime()})
	}

	return receipts, nil
}