package message

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"

	"github.com/jordan-wright/email"
	"github.com/sergeyzalunin/go-replication-loader/argsp"
	"github.com/sergeyzalunin/go-replication-loader/logger"
	rep "github.com/sergeyzalunin/go-replication-loader/replication"
)

// EmailMessage sends email by using inputs via ArgumentOptions
type EmailMessage struct {
	logger                logger.Log
	args                  argsp.ArgumentOptions
	deleteDescriptionFile bool
}

// New is a constructor for EmailMessageType
func New(args argsp.ArgumentOptions, logger logger.Log) EmailMessage {
	e := EmailMessage{
		args:   args,
		logger: logger,
	}
	return e
}

// Send message via email
func (em *EmailMessage) Send() {
	em.deleteDescriptionFile = true
	em.send(nil)
}

// SendFailed message via email if programm catches exception
func (em *EmailMessage) SendFailed(err interface{}) {
	em.deleteDescriptionFile = false
	em.send(err.(error))
}

func (em EmailMessage) send(err error) {
	if hasAnyEmailCommandLineParameters(em.args) {
		var e *email.Email

		if err == nil {
			e = em.getEmail()
		} else {
			e = em.getErrorEmail(err)
		}

		err = em.sendWithTLS(e)
		if err != nil {
			panic(err)
		}
	}
}

func hasAnyEmailCommandLineParameters(args argsp.ArgumentOptions) bool {
	result := strings.TrimSpace(args.SMTPServer) == "" ||
		strings.TrimSpace(args.SMTPLogin) == "" ||
		strings.TrimSpace(args.From) == "" ||
		hasAnyEmptyEmail(args.ToEmailList)

	return !result
}

func hasAnyEmptyEmail(toEmails []string) bool {
	if len(toEmails) == 0 {
		return true
	}

	for _, email := range toEmails {
		if strings.TrimSpace(email) == "" {
			return true
		}
	}

	return false
}

func (em EmailMessage) getEmail() *email.Email {
	e := email.Email{
		From:    em.args.From,
		To:      em.args.ToEmailList,
		Subject: em.getSubject(),
		Text:    em.getMessageBody(),
		Headers: textproto.MIMEHeader{},
	}
	_, err := e.AttachFile(em.logger.GetFileName())
	if err != nil {
		em.logger.Error("Couldn't attach log file due to error", err)
	}

	return &e
}

func (em EmailMessage) getErrorEmail(err error) *email.Email {
	e := email.Email{
		From:    em.args.From,
		To:      em.args.ToEmailList,
		Subject: em.getErrorSubject(),
		Text:    em.getErrorMessageBody(err),
		Headers: textproto.MIMEHeader{},
	}

	return &e
}

func (em EmailMessage) getSubject() string {
	eventTime := time.Now().Format("02.01.2006 15:04:05")
	return fmt.Sprintf("Replication on %s Base Completed Successfully at %s", em.args.ProjectName, eventTime)
}

func (em EmailMessage) getErrorSubject() string {
	eventTime := time.Now().Format("02.01.2006 15:04:05")
	return fmt.Sprintf("Replication on %s Base Failed at %s", em.args.ProjectName, eventTime)
}

func (em EmailMessage) getMessageBody() []byte {
	loader := rep.DescriptionLoader{}
	desc := loader.GetDescriptionContent(em.deleteDescriptionFile)
	result := fmt.Sprintf("%s\n\n%s", em.args.Body, desc)
	return []byte(result)
}

func (em EmailMessage) getErrorMessageBody(err error) []byte {
	result := fmt.Sprintf("The replication failed with next exception: %s\n"+
		"See the attached log file for details.", err.Error())
	return []byte(result)
}

func (em *EmailMessage) sendWithTLS(e *email.Email) error {
	addr := fmt.Sprintf("%s:%d", em.args.SMTPServer, em.args.SMTPPort)
	auth := smtp.PlainAuth(
		"",
		em.args.SMTPLogin,
		em.args.SMTPPassword,
		em.args.SMTPServer,
	)

	return e.SendWithTLS(addr, auth, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         em.args.SMTPServer,
	})
}
