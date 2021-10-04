package handler

import (
	"encoding/json"
	"net/http"

	"github.com/danisbagus/golang-messaging-rabbitmq/transaction-service/dto"
	"github.com/danisbagus/golang-messaging-rabbitmq/transaction-service/usecase"
)

type TransactionHandler struct {
	usecase usecase.ITransactionUsecase
}

func NewTransactionHanldler(usecase usecase.ITransactionUsecase) *TransactionHandler {
	return &TransactionHandler{
		usecase: usecase,
	}
}

func (rc TransactionHandler) NewTransaction(w http.ResponseWriter, r *http.Request) {
	var request dto.NewTransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := rc.usecase.Create(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeResponse(w, http.StatusCreated, data)
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
