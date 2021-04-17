package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database!")
	}

	// Migrates schema
	err = db.AutoMigrate(&Product{})
	if err != nil {
		panic(err)
	}

	// Create record
	db.Create(&Product{
		Code:  "D42",
		Price: 100,
	})

	// Read record
	var product Product
	db.First(&product, 1) // Gets first product with pk=1
	// SELECT * FROM product WHERE id = 1 ORDER BY id LIMIT 1;
	// fmt.Printf("Error: %s, Rows: %d, Statement: %v\n", db.Error, db.RowsAffected, db.Statement.Statement)

	db.First(&product, "code = ?", "D42") // Gets first product with code D42
	// SELECT * FROM product  WHERE code = 'D42' ORDER BY id LIMIT 1;

	// when querying, pointer to variable provided is used to modify variable
	fmt.Printf("%+v\n", product)
	// Need to check if DeletedAt has a value
	if !product.DeletedAt.Valid {
		fmt.Printf("Record deleted\n")

	}

	// Fetch multiple records into slice
	var products []Product
	db.Find(&products, "code = ?", "D42")
	for _, p := range products {
		if p.DeletedAt.Valid {
			continue // skip deleted records
		}
		fmt.Printf("id: %v, created: %s, updated: %s, code: %s, price: %v\n", p.ID, p.CreatedAt, p.UpdatedAt, p.Code, p.Price)

	}

	// Fetch single record into map
	result := make(map[string]interface{})
	db.Model(&Product{}).First(&result, "id = ?", 2)
	fmt.Printf("%s\n", result)

	// Fetch multiple records into slice of maps
	results := make([]map[string]interface{}, 1)
	db.Model(&Product{}).Find(&results, "code = ?", "D42")
	fmt.Printf("%s\n", results)

	// Update product's price to 200
	db.Model(&product).Update("Price", "200")
	// Update multiple fields
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	db.Model(&product).Updates(map[string]interface{}{"Price": 300, "Code": "G42"})

	// Multiple updates
	db.Model(Product{}).Where("price = ?", 300).Updates(Product{Price: 310})
	for _, p := range products {
		if p.DeletedAt.Valid {
			continue
		}
		db.Model(&product).Update("Price", p.Price*2)
	}

	// Delete product
	db.Delete(&product, 1)

}
