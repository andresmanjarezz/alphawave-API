package service

import (
	"fmt"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/config"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/email"
)

type EmailService struct {
	sender email.Sender
	config config.EmailConfig
}

func NewEmailService(sender email.Sender, config config.EmailConfig) *EmailService {
	return &EmailService{
		sender: sender,
		config: config,
	}
}

type VerificationEmailInput struct {
	Name  string
	Email string
	URL   string
}

type ForgotPasswordInput struct {
	Email            string
	TokenExpiresTime time.Duration
	URL              string
}

func (e *EmailService) SendUserVerificationEmail(input VerificationEmailInput) error {
	subject := fmt.Sprintf(e.config.Subjects.Verification, input.Name)
	sendInput := email.SendEmailInput{To: input.Email, Subject: subject}

	templateInput := VerificationEmailInput{Name: input.Name, URL: input.URL}

	err := sendInput.GenerateBodyFromHTML(e.config.Templates.Verification, templateInput)
	if err != nil {
		return err
	}
	err = e.sender.Send(sendInput)
	return err
}

func (e *EmailService) SendUserForgotPassword(input ForgotPasswordInput) error {
	subject := fmt.Sprintf(e.config.Subjects.ForgotPassword)
	sendInput := email.SendEmailInput{To: input.Email, Subject: subject}

	templateInput := ForgotPasswordInput{Email: input.Email, TokenExpiresTime: input.TokenExpiresTime, URL: input.URL}

	err := sendInput.GenerateBodyFromHTML(e.config.Templates.ForgotPassword, templateInput)
	if err != nil {
		return err
	}

	err = e.sender.Send(sendInput)
	return err
}
