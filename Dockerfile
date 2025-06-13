FROM ubuntu:latest
LABEL authors="evgen"

ENTRYPOINT ["top", "-b"]

FROM golang:1.24-alpine

WORKDIR /app

# Копируем файлы проекта
COPY . .

# Устанавливаем зависимости и собираем приложение
RUN go mod tidy && go build -o main .

# Открываем порт 8080
EXPOSE 8080

CMD ["./main"]