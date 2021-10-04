package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/danisbagus/golang-messaging-rabbitmq/common/messaging"
	"github.com/danisbagus/golang-messaging-rabbitmq/transaction-service/handler"
	"github.com/danisbagus/golang-messaging-rabbitmq/transaction-service/repo"
	"github.com/danisbagus/golang-messaging-rabbitmq/transaction-service/usecase"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// message broker driver
	// messagingClient := messaging.NewMessagingClient()
	messagingClient, err := messaging.GetMessagingConnection("amqp://guest:guest@localhost")
	if err != nil {
		fmt.Println("Error while connect to messaging client", err)
	}

	defer messagingClient.Close()

	// multiplexer
	router := mux.NewRouter()

	// injenction
	mailRepo := repo.NewMailRepo(messagingClient)
	transactionUsecase := usecase.NewTransactionUsecase(mailRepo)
	transactionHandler := handler.NewTransactionHanldler(transactionUsecase)

	// routing
	router.HandleFunc("/api/transactions", transactionHandler.NewTransaction).Methods(http.MethodPost)

	// starting server
	fmt.Println("Starting transaction service on port: 9050")
	log.Fatal(http.ListenAndServe("localhost:9050", router))

}
