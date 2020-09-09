package db

import (
	"context"

	logger "github.com/sirupsen/logrus"
)

const (
	getProductIDsQuery = `SELECT id from products`
)

type CompleteProduct struct {
	Product Product
	URLs    []string `json:"image_urls"`
	// []ProductImage.URL `db:"url" json:"image_url"`
}

func (s *pgStore) GetCompleteProductByID(ctx context.Context, Id int) (cp CompleteProduct, err error) {

	// _ is used as errors are handled in called functions
	Product, err := s.GetProductByID(ctx, Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting product from database by id " + string(Id))
		return
	}

	productImage, err := s.GetProductImagesByID(ctx, Id)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error selecting productImage from database by id " + string(Id))
		return
	}

	var Urls []string

	for j := 0; j < len(productImage); j++ {
		Urls = append(Urls, productImage[j].URL)
	}

	cp = CompleteProduct{
		Product: Product,
		URLs:    Urls,
	}

	return

}

func (s *pgStore) ListCompleteProducts(ctx context.Context) (completeProducts []CompleteProduct, err error) {

	// idArr stores id's of all products
	var idArr []int

	result, err := s.db.Query(getProductIDsQuery)
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
		var cp CompleteProduct
		cp, err = s.GetCompleteProductByID(ctx, int(idArr[i]))
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error selecting CompleteProduct from database by id " + string(idArr[i]))
			return
		}
		completeProducts = append(completeProducts, cp)
	}

	return
}
