package db

import (
	"context"
)

type Storer interface {
	ListUsers(context.Context) ([]User, error)
	ListCompleteProducts(context.Context) ([]CompleteProduct, error)
	CreateNewProduct(context.Context, Product) (Product, error)
	GetProductByID(context.Context, int) (Product, error)
	//GetProductByIdWithCategory(ctx context.Context, Id int)(Product, Category, error)
	DeleteProductById(context.Context, int) error
	UpdateProductById(context.Context, Product, int) (CompleteProduct, error)
	GetProductImagesByID(context.Context, int) ([]ProductImage, error)
	GetCompleteProductByID(context.Context, int) (CompleteProduct, error)
}
