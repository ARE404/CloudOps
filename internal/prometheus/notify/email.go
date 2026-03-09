package notify

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/spf13/viper"
)

type EmailNotifier struct{}

func NewEmailNotifier() *EmailNotifier {
	return &EmailNotifier{}
}

func (n *EmailNotifier) Send(ctx context.Context, subject, content string, channels []NotifyChannel) error {
	host := viper.GetString("notification.email.smtp_host")
	port := viper.GetInt("notification.email.smtp_port")
	username := viper.GetString("notification.email.username")
	password := viper.GetString("notification.email.password")

	if host == "" || username == "" {
		return nil // email not configured, skip
	}

	auth := smtp.PlainAuth("", username, password, host)
	addr := fmt.Sprintf("%s:%d", host, port)

	for _, ch := range channels {
		if ch.Type != "email" {
			continue
		}
		msg := strings.Join([]string{
			"From: " + username,
			"To: " + ch.Target,
			"Subject: " + subject,
			"MIME-Version: 1.0",
			"Content-Type: text/plain; charset=UTF-8",
			"",
			content,
		}, "\r\n")
		if err := smtp.SendMail(addr, auth, username, []string{ch.Target}, []byte(msg)); err != nil {
			return err
		}
	}
	return nil
}
