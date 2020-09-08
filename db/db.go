package db

import (
	"context"
)

type Storer interface {
	ListUsers(context.Context) ([]User, error)
	ListProducts(context.Context) ([]Product, error)
	CreateNewProduct(context.Context, Product)(Product, error)
	GetProductById(context.Context, int) (Product, error)
	//Create(context.Context, User) error
	//GetUser(context.Context) (User, error)
	//Delete(context.Context, string) error
}
