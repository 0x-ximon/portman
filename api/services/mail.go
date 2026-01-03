package services

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/resend/resend-go/v3"
)

type MailService struct {
	client *resend.Client
}

func NewMailService() (*MailService, error) {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("RESEND_API_KEY is not set")
	}

	return &MailService{
		client: resend.NewClient(apiKey),
	}, nil
}

func (m *MailService) SendOTP(to string, otp string) error {
	data := struct {
		OTP string
	}{
		OTP: otp,
	}

	temp, err := template.ParseFiles("templates/otp.tmpl")
	if err != nil {
		return err
	}

	var b bytes.Buffer
	err = temp.Execute(&b, data)
	if err != nil {
		return err
	}

	emailParams := &resend.SendEmailRequest{
		From:    "Portman <onboarding@mail.ximon.dev>",
		To:      []string{to},
		Html:    b.String(),
		Subject: "OTP Verification",
	}

	_, err = m.client.Emails.Send(emailParams)
	if err != nil {
		return err
	}

	return nil
}
