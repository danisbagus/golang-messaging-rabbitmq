package repo

import (
	"encoding/json"
	"fmt"

	"github.com/danisbagus/golang-messaging-rabbitmq/common/messaging"
	"github.com/danisbagus/golang-messaging-rabbitmq/transaction-service/model"
)

type IMailRepo interface {
	SendMail(*model.MailModel)
}

type MailRepo struct {
	mc messaging.Connection
}

func NewMailRepo(mc messaging.Connection) IMailRepo {
	return &MailRepo{
		mc: mc,
	}
}

func (r MailRepo) SendMail(data *model.MailModel) {
	fmt.Println("Running SendMail")

	mailData, _ := json.Marshal(&data)

	err := r.mc.PublishQueue([]byte(mailData), "sendMailQueue")
	if err != nil {
		fmt.Println("Error while publish queue on SendMail: ", err)
	}
}
