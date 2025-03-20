# Базовый образ
FROM golang:1.23-alpine

# Установка зависимостей
RUN apk add --no-cache bash git


# Установка air
RUN go install github.com/air-verse/air@latest

# Создание рабочей директории
WORKDIR /app

# Копирование go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Команда для запуска air
CMD ["air"]