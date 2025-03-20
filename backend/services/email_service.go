package services

import (
    "gopkg.in/gomail.v2"
    "backend/config"
)

type EmailService struct{}

func NewEmailService() *EmailService {
    return &EmailService{}
}

func (s *EmailService) SendEmail(to, subject, body string) error {
    cfg := config.GetEmailConfig()

    m := gomail.NewMessage()
    m.SetHeader("From", cfg.From)
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)

    d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

    if err := d.DialAndSend(m); err != nil {
        return err
    }

    return nil
}