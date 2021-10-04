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

type Request struct {
	to      []string
	subject string
	body    string
}

type TemplateEmail struct {
	Name          string
	TransactionID string
}

type Message map[string]interface{}

func main() {
	// load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// set request data
	reqMail := &Request{
		to:      []string{"danishter22@gmail.com"},
		subject: "Transaction Paid",
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
			transactionID := fmt.Sprintf("%v", msg["transaction_id"])

			// set template data
			templateData := TemplateEmail{
				Name:          customerName,
				TransactionID: transactionID,
			}

			// parse template
			errParse := reqMail.ParseTemplate("mail-service/transaction_paid.html", templateData)
			if errParse != nil {
				log.Fatal("Error while parse the template", err.Error())
			}

			// send email
			errSend := reqMail.SendEmail()
			if errSend != nil {
				log.Fatal("Error while send email", errSend.Error())
			}

			log.Println("Successfully to send email")

		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

func (r *Request) SendEmail() error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	from := "From: " + os.Getenv("CONFIG_SENDER_NAME") + "\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(from + subject + mime + "\n" + r.body)

	smtpAddr := fmt.Sprintf("%s:%s", os.Getenv("CONFIG_SMTP_HOST"), os.Getenv("CONFIG_SMTP_PORT"))
	auth := smtp.PlainAuth("", os.Getenv("CONFIG_AUTH_EMAIL"), os.Getenv("CONFIG_AUTH_PASSWORD"), os.Getenv("CONFIG_SMTP_HOST"))

	if err := smtp.SendMail(smtpAddr, auth, os.Getenv("CONFIG_AUTH_EMAIL"), r.to, msg); err != nil {
		return err
	}
	return nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

func Deserialize(b []byte) (Message, error) {
	var msg Message
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}
