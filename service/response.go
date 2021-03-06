package service

import(
	"encoding/json"
	"net/http"

	logger "github.com/sirupsen/logrus"
)

type successResponse struct {
	Data interface{} `json:"data"`
}

type errorResponse struct{
	Error interface{} `json:"error"`
}

type messageObject struct{
	Message string `json:"message"`
}

type errorObject struct {
	Code string `json:"code"`
	messageObject
	Fields map[string]string `json:"fields"`
}

func response(rw http.ResponseWriter, status int, responseBody interface{}){
	respBytes, err := json.Marshal(responseBody)
	rw.Header().Add("Content-Type","application/json")
	if err != nil{
		logger.WithField("err", err.Error()).Error("Error while marshaling core values data")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(status)
	rw.Write(respBytes)
}