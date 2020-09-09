package db

import (
	"context"
)

type Storer interface {
	ListUsers(context.Context) ([]User, error)
	ListProducts(context.Context) ([]Product, error)
	CreateNewProduct(context.Context, Product) (Product, error)
	//GetProductByIdWithCategory(ctx context.Context, Id int)(Product, Category, error)
	DeleteProductById(context.Context, int) error
	UpdateProductById(context.Context, Product, int) (Product, error)
	GetProductImagesByID(context.Context, int) ([]ProductImage, error)
	GetProductByID(context.Context, int) (Product, error)
}
