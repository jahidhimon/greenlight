package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

// Below we declare a new variable with the type embed.FS (embedded file system)
// to hold our email templates. This has a comment directive in the format
// `//go:embed <path>` Immediately above it which indicates to Go that we
// want to store the contents of the ./templates directory in the templateFS
// embedded file system variable

//go:embed "templates"
var templateFS embed.FS

// Define mailer struct which contains a mail.Dialer instance
// (used to connect to a SMTP server) and the sender information
// for you emails (the name and address you want to the mail to be from)

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	// Initialize a new mailer dialer instance with the given SMTP settings.
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

// Define Send() method on the Mailer type. This takes the recipent email
// address as the first parameter, the name of the file containing the
// templates, and any dynamic data for the templates as an interface{}
// parameter
func (m Mailer) Send(recipent, templateFile string, data interface{}) error {
	// Use the parseFS() method to parse the required template file from
	// the embedded file system.
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		fmt.Println("fucking plainbody");
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err!= nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	// Use the mail.NewMessage() function to initialize a new mail.Message instance.
	// Then we use the SetHeader() method to set email recipent, sender and subject
	// SetBody() method to set the plain-text-body, and the AddAlternative() method
	// to set the HTML body. It's important to note that AddAlternative() should
	// alsways be called *after* SetBody.
	msg := mail.NewMessage()
	msg.SetHeader("To", recipent)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
