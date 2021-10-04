package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"

	"github.com/danisbagus/golang-messaging-rabbitmq/common/messaging"
	"github.com/joho/godotenv"
)

type TemplateEmail struct {
	Name          string
	TransactionID string
}

type MailTransactionPaidData struct {
	Name          string
	Email         string
	TransactionID string
}

type Message map[string]interface{}

func main() {
	// load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// messaging client driver
	messagingClient, err := messaging.GetMessagingConnection("amqp://guest:guest@localhost")
	if err != nil {
		fmt.Println("Error while connect to broker", err)
	}

	defer messagingClient.Close()

	msgs, err := messagingClient.ConsumeQueue("sendMailQueue")
	if err != nil {
		fmt.Println("Error while cosume queue", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			msg, err := Deserialize(d.Body)
			if err != nil {
				log.Fatal("Error while deserialize message", err.Error())
			}

			customerName := fmt.Sprintf("%v", msg["customer_name"])
			customerEmail := fmt.Sprintf("%v", msg["customer_email"])
			transactionID := fmt.Sprintf("%v", msg["transaction_id"])

			mailData := MailTransactionPaidData{
				Name:          customerName,
				Email:         customerEmail,
				TransactionID: transactionID,
			}

			errSend := SendMail(mailData)
			if errSend != nil {
				log.Fatal("Error while send email", errSend.Error())
			}
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-forever

}

func SendMail(data MailTransactionPaidData) error {
	targetEmail := []string{data.Email}
	subjectEmail := "Transaction Paid"

	templateData := TemplateEmail{
		Name:          data.Name,
		TransactionID: data.TransactionID,
	}

	t, err := template.ParseFiles("mail-service/transaction_paid.html")
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, templateData); err != nil {
		return err
	}

	emailBody := buf.String()

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	from := "From: " + os.Getenv("CONFIG_SENDER_NAME") + "\n"
	subject := "Subject: " + subjectEmail + "!\n"
	msg := []byte(from + subject + mime + "\n" + emailBody)

	smtpAddr := fmt.Sprintf("%s:%s", os.Getenv("CONFIG_SMTP_HOST"), os.Getenv("CONFIG_SMTP_PORT"))
	auth := smtp.PlainAuth("", os.Getenv("CONFIG_AUTH_EMAIL"), os.Getenv("CONFIG_AUTH_PASSWORD"), os.Getenv("CONFIG_SMTP_HOST"))

	if err := smtp.SendMail(smtpAddr, auth, os.Getenv("CONFIG_AUTH_EMAIL"), targetEmail, msg); err != nil {
		return err
	}

	log.Println("Successfully to send email")

	return nil

}

func Deserialize(b []byte) (Message, error) {
	var msg Message
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}
