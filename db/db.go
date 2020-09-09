package db

import (
	"context"
)

type Storer interface {
	ListUsers(context.Context) ([]User, error)
	ListProducts(context.Context) ([]Product, error)
	CreateNewProduct(context.Context, Product)(Product, error)
	GetProductById(context.Context, int) (Product, error)
	DeleteProductById(context.Context, int) error
	UpdateProductById(context.Context, Product, int) (Product, error)
}
