package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"santoshkavhar/go-e-commerce/db"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProductsHandlerTestSuite struct {
	suite.Suite

	dbMock *db.DBMockStore
}

func (suite *ProductsHandlerTestSuite) SetupTest() {
	suite.dbMock = &db.DBMockStore{}
}

func TestProductTestSuite(t *testing.T) {
	suite.Run(t, new(ProductsHandlerTestSuite))
}

func (suite *ProductsHandlerTestSuite) TestListProductsSucces() {
	suite.dbMock.On("ListProducts", mock.Anything).Return(

		[]db.Product{
			db.Product{Id: 0, Name: "test-Product", Description: "test-descp", Price: 12.01, Discount: 0.0, Quantity: 0, CategoryId: 1},
		},
		nil,
	)

	recorder := makeHTTPCalls(
		http.MethodGet,
		"/products",
		"",
		listProductsHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), `[{"Id": 0, "Name": "test-Product", "Description": "test-descp", "Price": 12.01, "Discount": 0.0, "Quantity": 0, "CategoryId": 1 }
	}]`, recorder.Body.String())
	suite.dbMock.AssertExpectations(suite.T())
}

func (suite *ProductsHandlerTestSuite) TestListProductsWhenDBFailure() {
	suite.dbMock.On("ListProducts", mock.Anything).Return(
		[]db.Product{},
		errors.New("error fetching user records"),
	)

	recorder := makeHTTPCalls(
		http.MethodGet,
		"/products",
		"",
		listProductsHandler(Dependencies{Store: suite.dbMock}),
	)

	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
	suite.dbMock.AssertExpectations(suite.T())
}

func makeHTTPCalls(method, path, body string, handlerFunc http.HandlerFunc) (record *httptest.ResponseRecorder) {

	req, _ := http.NewRequest(method, path, strings.NewReader(body))

	recorder := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc(path, handlerFunc).Methods(method)

	router.ServeHTTP(recorder, req)
	return
}
