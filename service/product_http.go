package service

import(
	"encoding/json"
	"net/http"

	"santoshkavhar/go-e-commerce/db"
	logger "github.com/sirupsen/logrus"
	//"github.com/gorilla/mux"
)

// @Title listProducts
// @Description list all Products
// @Router /products [GET]
// @Accept	json
// @Success 200 {object}
// @Failure 400 {object}
func listProductsHandler(deps Dependencies) http.HandlerFunc{
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request){

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

		rw.Header().Add("Content-Type","application/json")
		rw.Write(respBytes)
	})
}

// @Title createProduct
// @Description create a Product, insert into DB
// @Router /createProduct [POST]
// @Accept	json
// @Success 200 {object}
// @Failure 400 {object}
func createProductHandler(deps Dependencies) http.HandlerFunc{
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request){
		
		var product db.Product
		err := json.NewDecoder(req.Body).Decode(&product)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			logger.WithField("err", err.Error()).Error("Error while decoding product")
			response(rw, http.StatusBadRequest, errorResponse{	
				Error : messageObject {
					Message:"Invalid json body",
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