package config

type EmailConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    From     string
}

func GetEmailConfig() *EmailConfig {
    return &EmailConfig{
        Host:     "smtp.example.com", // SMTP-сервер (например, Gmail: smtp.gmail.com)
        Port:     587,               // Порт SMTP (обычно 587 для TLS)
        Username: "your-email@example.com", // Ваш email
        Password: "your-email-password",   // Пароль от email
        From:     "your-email@example.com", // Отправитель
    }
}