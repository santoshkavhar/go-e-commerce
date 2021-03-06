package db

import (
	"context"
	"fmt"

	//"database/sql"
	logger "github.com/sirupsen/logrus"
)

const (
	getProductIDQuery = `SELECT id FROM products`
	// id is PRIMARY KEY, so no need to limit
	getProductByIDQuery = `SELECT * FROM products WHERE id=$1`
	insertProductQuery  = `INSERT INTO products (
		id, name, description, price, discount, quantity, category_id) VALUES ( :id, :name, :description, :price, :discount, :quantity, :category_id)`
	updateProductQuery = `UPDATE products SET (
		name, description, price, discount, quantity, category_id
			) =  ($1, $2, $3, $4, $5, $6) where id = $7`
	deleteProductIdQuery = `DELETE FROM products WHERE id = $1`

	getCategoryByID = `SELECT name FROM category WHERE id = $1`
)

type Product struct {
	Id           int      `db:"id" json:"product_id"`
	Name         string   `db:"name" json:"product_name"`
	Description  string   `db:"description" json:"product_description"`
	Price        float32  `db:"price" json:"price"`
	Discount     float32  `db:"discount" json:"discount"`
	Quantity     int      `db:"quantity" json:"available_quantity"`
	CategoryId   int      `db:"category_id" json:"category_id"`
	CategoryName string   `json:"category_name,omitempty"`
	URLs         []string `json:"productimage_urls,omitempty"`
}

func (product *Product) Validate() (errorResponse map[string]ErrorResponse, valid bool) {
	fieldErrors := make(map[string]string)

	if product.Id == 0 {
		fieldErrors["product_id"] = "Can't be blank"
	}
	if product.Name == "" {
		fieldErrors["product_name"] = "Can't be blank"
	}
	if product.Description == "" {
		fieldErrors["product_description"] = "Can't be blank"
	}
	if product.Price <= 0 {
		fieldErrors["price"] = "Can't be blank"
	}
	if product.Discount < 0 {
		fieldErrors["discount"] = "Can't be blank"
	}
	if product.Quantity < 0 {
		fieldErrors["available_quantity"] = "Can't be blank"
	}
	if product.CategoryId == 0 {
		fieldErrors["category_id"] = "Can't be blank"
	}

	if len(fieldErrors) == 0 {
		valid = true
		return
	}

	errorResponse = map[string]ErrorResponse{
		"error": ErrorResponse{
			Code:    "Invalid_data",
			Message: "Please Provide valid Product data",
			Fields:  fieldErrors,
		},
	}
	// TODO Other Validations
	return
}

func (s *pgStore) GetProductByID(ctx context.Context, Id int) (product Product, err error) {

	err = s.db.Get(&product, getProductByIDQuery, Id)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting product from database by id: " + string(Id))
		return
	}

	var category string
	err = s.db.Get(&category, getCategoryByID, product.CategoryId)

	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching category from database by product_id: " + string(Id))
		return
	}

	product.CategoryName = category

	productImage, err := s.GetProductImagesByID(ctx, Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting productImage from database by id " + string(Id))
		return
	}

	for j := 0; j < len(productImage); j++ {
		product.URLs = append(product.URLs, productImage[j].URL)
	}

	return

}

func (s *pgStore) ListProducts(ctx context.Context) (products []Product, err error) {

	// idArr stores id's of all products
	var idArr []int

	result, err := s.db.Query(getProductIDQuery)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error fetching Product Ids from database")
		return
	}

	for result.Next() {
		var Id int
		err = result.Scan(&Id)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Couldn't Scan Product ids")
			return
		}
		idArr = append(idArr, Id)
	}

	for i := 0; i < len(idArr); i++ {
		var product Product
		product, err = s.GetProductByID(ctx, int(idArr[i]))
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error selecting Product from database by id " + string(idArr[i]))
			return
		}
		products = append(products, product)
	}

	return
}

/*
func (s *pgStore) GetProductByIdWithCategory(ctx context.Context, Id int) (product Product, category Category,  err error) {
	err = s.db.Get(&product, getProductByIDQuery, Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting product from database by id " + string(Id))
		return
	}

	err = s.db.Get(&category, getCategoryByID, product.CategoryId)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting category from database by id " + string(Id))
		return
	}
	return
}
*/

// CreateNewProduct
func (s *pgStore) CreateNewProduct(ctx context.Context, p Product) (createdProduct Product, err error) {
	// First, make sure Product isn't already in db, if Product is present, just return the it
	createdProduct, err = s.GetProductByID(ctx, p.Id)
	if err == nil {
		// If there's already a product, err wil be nil, so no new Product is populated.
		err = fmt.Errorf("Product Already exists!")
		return
	}
	tx, err := s.db.Beginx() // Use Beginx instead of MustBegin so process doesn't die if there is an error
	if err != nil {
		// FAIL : Could not begin database transaction
		logger.WithField("err", err.Error()).Error("Error beginning product insert transaction in db.CreateNewProduct with Id: " + string(p.Id))
		return
	}

	_, err = tx.NamedExec(insertProductQuery, p)
	// p.Id, p.Name, p.Description, p.Price, p.Discount, p.Quantity, p.CategoryId

	if err != nil {
		// FAIL : Could not run insert Query
		logger.WithField("err", err.Error()).Error("Error inserting product to database: " + p.Name)
		return
	}
	err = tx.Commit()
	if err != nil {
		// FAIL : transaction commit failed.Will Automatically rollback
		logger.WithField("err", err.Error()).Error("Error commiting transaction inserting product into database: " + string(p.Id))
		return
	}

	// Re-select Product and return it
	createdProduct, err = s.GetProductByID(ctx, p.Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting from database with id: " + string(p.Id))
		return
	}
	return
}

func (s *pgStore) UpdateProductById(ctx context.Context, product Product, Id int) (updatedProduct Product, err error) {

	var dbProduct Product
	err = s.db.Get(&dbProduct, getProductByIDQuery, Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while fetching product ")
		return
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.WithField("err:", err.Error()).Error("Error while initiating update transaction")
		return
	}

	_, err = tx.Exec(updateProductQuery,
		product.Name,
		product.Description,
		product.Price,
		product.Discount,
		product.Quantity,
		product.CategoryId,
		Id,
	)

	if err != nil {
		// FAIL : Could not Update Product
		logger.WithField("err", err.Error()).Error("Error updating product attribute(s) to database :" + string(Id))
		return
	}

	err = tx.Commit()
	if err != nil {
		// FAIL : transaction commit failed. Will Automatically rollback
		logger.WithField("err", err.Error()).Error("Error commiting transaction updating product into database: " + string(Id))
		return
	}

	updatedProduct, err = s.GetProductByID(ctx, Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error while getting updated product ")
		return
	}
	return
}

func (s *pgStore) DeleteProductById(ctx context.Context, Id int) (err error) {

	rows, err := s.db.Exec(deleteProductIdQuery, Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error deleting product" + string(Id))
		return
	}

	rows_affected, err := rows.RowsAffected()
	if rows_affected == 0 {
		err = fmt.Errorf("Product doesn't exist in db, goodluck deleting it")
	}
	return
}
