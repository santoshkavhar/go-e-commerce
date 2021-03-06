package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"santoshkavhar/go-e-commerce/db"

	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
)

// @Title listProducts
// @Description list all Products
// @Router /products [GET]
// @Accept	json
// @Success 200 {object}
// @Failure 400 {object}
func listProductsHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		products, err := deps.Store.ListProducts(req.Context())
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error fetching data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBytes, err := json.Marshal(products)
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error mashaling products data")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Add("Content-Type", "application/json")
		rw.Write(respBytes)
	})
}

// @Title createProduct
// @Description create a Product, insert into DB
// @Router /createProduct [POST]
// @Accept	json
// @Success 200 {object}
// @Failure 400 {object}
func createProductHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		var product db.Product
		err := json.NewDecoder(req.Body).Decode(&product)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			logger.WithField("err", err.Error()).Error("Error while decoding product")
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Invalid json body",
				},
			})
			return
		}

		errRes, valid := product.Validate()
		if !valid {
			_, err := json.Marshal(errRes)
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error marshalling Product's data")
				response(rw, http.StatusBadRequest, errorResponse{
					Error: messageObject{
						Message: "Invalid json body",
					},
				})
				return
			}
			response(rw, http.StatusBadRequest, errRes)
			return
		}

		var createdProduct db.Product
		createdProduct, err = deps.Store.CreateNewProduct(req.Context(), product)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Error inserting the product, possibly not new",
				},
			})
			logger.WithField("err", err.Error()).Error("Error while inserting product")
			return
		}

		response(rw, http.StatusOK, successResponse{Data: createdProduct})
		return
	})
}

// @ Title getCompleteProductById
// @ Description get single product by its id
// @ Router /product/product_id [get]
// @ Accept json
// @ Success 200 {object}
// @ Failure 400 {object}

func getProductByIdHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["product_id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error id key is missing")
			rw.WriteHeader(http.StatusBadRequest)
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error product_id is invalid",
				},
			})
			return
		}

		product, err := deps.Store.GetProductByID(req.Context(), id)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			logger.WithField("err", err.Error()).Error("Error fetching data")
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error feching data No Row Found",
				},
			})
			return
		}

		_, err = json.Marshal(product)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			logger.WithField("err", err.Error()).Error("Error marshaling products data")
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Internal server error",
				},
			})
			return
		}
		response(rw, http.StatusOK, product)
		return
	})
}

// @ Title deleteProductById
// @ Description delete product by its id
// @ Router /product/product_id [delete]
// @ Accept json
// @ Success 200 {object}
// @ Failure 400 {object}

func deleteProductByIdHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["product_id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error id key is missing")
			rw.WriteHeader(http.StatusBadRequest)
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error id is missing/invalid",
				},
			})
			return
		}

		err = deps.Store.DeleteProductById(req.Context(), id)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			logger.WithField("err", err.Error()).Error("Error fetching data no row found")
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Internal server error  (Error feching data)",
				},
			})
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Header().Add("Content-Type", "application/json")
		return
	})
}

// @ Title updateProductById
// @ Description update product by its id
// @ Router /product/product_id [put]
// @ Accept json
// @ Success 200 {object}
// @ Failure 400 {object}

func updateProductByIdHandler(deps Dependencies) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		vars := mux.Vars(req)
		id, err := strconv.Atoi(vars["product_id"])
		if err != nil {
			logger.WithField("err", err.Error()).Error("Error id key is missing")
			rw.WriteHeader(http.StatusBadRequest)
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Error id is missing/invalid",
				},
			})
			return
		}

		var product db.Product
		err = json.NewDecoder(req.Body).Decode(&product)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			logger.WithField("err", err.Error()).Error("Error while decoding user")
			response(rw, http.StatusBadRequest, errorResponse{
				Error: messageObject{
					Message: "Internal server error",
				},
			})
			return
		}

		errRes, valid := product.Validate()
		if !valid {
			respBytes, err := json.Marshal(errRes)
			if err != nil {
				logger.WithField("err", err.Error()).Error("Error marshaling product data")
				response(rw, http.StatusBadRequest, errorResponse{
					Error: messageObject{
						Message: "Invalid json body",
					},
				})
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			rw.Header().Add("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(respBytes)
			return
		}

		var updatedProduct db.Product
		updatedProduct, err = deps.Store.UpdateProductById(req.Context(), product, id)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response(rw, http.StatusInternalServerError, errorResponse{
				Error: messageObject{
					Message: "Internal server error",
				},
			})
			logger.WithField("err", err.Error()).Error("Error while updating product attribute")
			return
		}

		response(rw, http.StatusOK, successResponse{Data: updatedProduct})

		return
	})
}
