package endpoint

import (
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/scorredoira/email"
)

func (e *endpoint) sendEmail(data []byte) error {
	log.Printf("%s %q Sending email.\n", e.uid, e.name)
	m := email.NewMessage(e.c.Get("MessageSubject"),
		e.c.Get("MessageBody"))
	m.From = mail.Address{
		Name:    e.c.Get("MessageFromName"),
		Address: e.c.Get("MessageFromAddress"),
	}
	m.To = []string{e.name}

	err := m.AttachBuffer("scan.jpg", data, false)
	if err != nil {
		return err
	}

	smtpServerPort := strings.Split(e.c.Get("SmtpServerPort"), ":")

	auth := smtp.PlainAuth("", e.c.Get("SmtpAuthUser"),
		e.c.Get("SmtpAuthPassword"), smtpServerPort[0])

	fmt.Println(e.c.Get("SmtpServerPort"), auth, m)
	err = email.Send(e.c.Get("SmtpServerPort"), auth, m)
	if err != nil {
		return err
	}

	log.Printf("%s %q Email sent.\n", e.uid, e.name)

	return e.register()
}
