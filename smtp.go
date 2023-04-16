package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func SendEmail(subject, content string) error {
	log.Println("sending notifications for", subject)

	if CONF.SMTP.Host == "" {
		log.Println("error sending email: invalid SMTP configuration")
		return nil
	}

	recipients, err := ReadEmails()
	if err != nil {
		return fmt.Errorf("error reading email recipients: %w", err)
	}

	for _, r := range recipients {
		from := "From: " + CONF.SMTP.Sender + "\r\n"
		to := "To: " + r + "\r\n"
		subj := "Subject: " + subject + "\r\n"
		mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		mesg := []byte(from + to + subj + mime + content)

		addr := fmt.Sprintf("%s:%d", CONF.SMTP.Host, CONF.SMTP.Port)
		auth := smtp.PlainAuth("", CONF.SMTP.Username, CONF.SMTP.Password, CONF.SMTP.Host)
		if err := smtp.SendMail(addr, auth, CONF.SMTP.Sender, []string{r}, mesg); err != nil {
			return fmt.Errorf("error sending email: %w", err)
		}

		log.Println("notification sent to", r)
	}

	return nil
}
