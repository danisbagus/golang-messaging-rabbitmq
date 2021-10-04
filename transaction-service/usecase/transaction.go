package usecase

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/danisbagus/golang-messaging-rabbitmq/transaction-service/dto"
	"github.com/danisbagus/golang-messaging-rabbitmq/transaction-service/model"
	"github.com/danisbagus/golang-messaging-rabbitmq/transaction-service/repo"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type ITransactionUsecase interface {
	Create(*dto.NewTransactionRequest) (*dto.NewTransactionResponse, error)
}

type TransactionUsecase struct {
	MailRepo repo.IMailRepo
}

func NewTransactionUsecase(mailRepo repo.IMailRepo) ITransactionUsecase {
	return &TransactionUsecase{
		MailRepo: mailRepo,
	}
}

func (r TransactionUsecase) Create(data *dto.NewTransactionRequest) (*dto.NewTransactionResponse, error) {
	transactionID := fmt.Sprintf("TR%v", String(6))

	form := model.TransactionModel{
		TransactionID:   transactionID,
		ProductID:       data.ProductID,
		Quantity:        data.Quantity,
		TransactionDate: time.Now(),
	}

	sendMailData := model.MailModel{
		TransactionID: transactionID,
		CustomerName:  "danisbagus22",
		CustomerEmail: "danishter22@gmail.com",
	}

	go r.MailRepo.SendMail(&sendMailData)

	response := dto.NewNewTransactionResponse(&form)

	return response, nil
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}
