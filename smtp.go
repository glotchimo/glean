package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func SendEmail(subject, content string) error {
	log.Printf("sending email: %s\n", subject)

	if CONF.SMTP.Host == "" {
		log.Println("error sending email: invalid SMTP configuration")
		return nil
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	mesg := []byte(subject + "\r\n" + mime + content)

	addr := fmt.Sprintf("%s:%d", CONF.SMTP.Host, CONF.SMTP.Port)
	auth := smtp.PlainAuth("", CONF.SMTP.Username, CONF.SMTP.Password, CONF.SMTP.Host)

	recipients, err := ReadEmails()
	if err != nil {
		return fmt.Errorf("error reading email recipients: %w", err)
	}

	if err := smtp.SendMail(addr, auth, CONF.SMTP.Sender, recipients, mesg); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	log.Println("sent!")

	return nil
}
