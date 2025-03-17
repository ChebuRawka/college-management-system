# Используем официальный образ Go
FROM golang:1.23-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Устанавливаем air для горячей перезагрузки
RUN go install github.com/air-verse/air@latest

# Копируем исходный код
COPY . .

# Команда для запуска приложения с горячей перезагрузкой
CMD ["air"]