package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"

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

	// set template data
	templateData := TemplateEmail{
		Name:          "Coder Pemula",
		TransactionID: "PT10001",
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
