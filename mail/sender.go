package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpGmailHost    = "smtp.gmail.com"
	smtpGmailAddress = "smtp.gmail.com:587"
)

type Sender interface {
	SendEmail(
		toEmailAdress []string,
		subject string,
		content string, // HTML format
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	FromEmailName     string
	FromEmailAddress  string
	FromEmailPassword string
}

func NewGmailSender(fromEmailName, fromEmailAddress, fromEmailPassword string) Sender {
	return &GmailSender{
		FromEmailName:     fromEmailName,
		FromEmailAddress:  fromEmailAddress,
		FromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	toEmailAdress []string,
	subject string,
	content string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.FromEmailName, sender.FromEmailAddress)
	e.To = toEmailAdress
	e.Subject = subject
	e.HTML = []byte(content)
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return err
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.FromEmailAddress, sender.FromEmailPassword, smtpGmailHost)

	return e.Send(smtpGmailAddress, smtpAuth)
}
