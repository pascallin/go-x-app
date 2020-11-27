package cli

import (
	"github.com/pascallin/go-x-app/internal"
	"log"
	"os"
)

func main() {
	err := internal.Insert(&internal.InsertOptions{
		ExcelFilePath: "/Users/pascal/Documents/Book1.xlsx",
		SheetName: "Sheet1",
		ImageDir: "/Users/pascal/Pictures/",
		KeyColumn: "A",
	})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}