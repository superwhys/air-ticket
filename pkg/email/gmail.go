package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const (
	subject    = "机票监控"
	smtpServer = "smtp.gmail.com"
	smtpPort   = "465"
)

type EmailConf struct {
	Sender   string
	Password string
}

type GmailSender struct {
	conf    *EmailConf
	auth    smtp.Auth
	tlsConf *tls.Config
}

func NewGmailSender(conf *EmailConf) *GmailSender {
	auth := smtp.PlainAuth("", conf.Sender, conf.Password, smtpServer)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer,
	}
	return &GmailSender{
		conf:    conf,
		auth:    auth,
		tlsConf: tlsConfig,
	}
}

func (e *GmailSender) getClient(target string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", smtpServer+":"+smtpPort, e.tlsConf)
	if err != nil {
		return nil, err
	}

	client, err := smtp.NewClient(conn, smtpServer)
	if err != nil {
		return nil, err
	}

	if err = client.Auth(e.auth); err != nil {
		return nil, err
	}

	if err = client.Mail(e.conf.Sender); err != nil {
		return nil, err
	}
	if err = client.Rcpt(target); err != nil {
		return nil, err
	}

	return client, nil
}

func (e *GmailSender) wrapMsg(subject, target string, msg []byte) []byte {
	subject = fmt.Sprintf("Subject: %s!\r\n,", subject)
	fromHeader := fmt.Sprintf("From: %s\r\n", e.conf.Sender)
	toHeader := fmt.Sprintf("To: %s\r\n", target)
	contentTypeHeader := "Content-Type: text/plain; charset=\"UTF-8\"\r\n"
	mimeVersionHeader := "MIME-Version: 1.0\r\n"

	data := []byte(fromHeader + toHeader + subject + mimeVersionHeader + contentTypeHeader + "\r\n")
	return append(data, msg...)
}

func (e *GmailSender) SentMsg(ctx context.Context, targets []string, msg []byte) error {
	grp, ctx := errgroup.WithContext(ctx)
	grp.SetLimit(len(targets))

	for _, target := range targets {
		target := target
		grp.Go(func() error {
			client, err := e.getClient(target)
			if err != nil {
				return errors.Wrapf(err, "getSmtpClient: %v", target)
			}
			defer client.Quit()

			w, err := client.Data()
			if err != nil {
				return errors.Wrapf(err, "getWriter: %v", target)
			}
			defer w.Close()

			_, err = w.Write(e.wrapMsg(subject, target, msg))
			if err != nil {
				return errors.Wrapf(err, "writeData: %v", target)
			}

			return nil
		})
	}
	return grp.Wait()
}
