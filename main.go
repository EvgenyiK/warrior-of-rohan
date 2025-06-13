package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Структура для передачи данных в шаблон
type PageData struct {
	Title   string
	Content string
}

func main() {
	// Подключение к базе данных PostgreSQL
	connStr := os.Getenv("POSTGRES_CONN")
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Обработка маршрутов
	http.HandleFunc("/", handleGame)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	// Получение текущего состояния игры из базы данных или создание нового
	var content string

	// Для простоты: всегда начинаем с одного и того же состояния
	content = "Вы находитесь в темной комнате. Что вы делаете? (введите 'осмотреться' или 'идти дальше')"

	if r.Method == http.MethodPost {
		action := r.FormValue("action")
		switch action {
		case "осмотреться":
			content = "Вы увидели старую дверь и окно."
		case "идти дальше":
			content = "Вы вышли из комнаты и нашли выход!"
		default:
			content = "Неизвестное действие."
		}
		// Можно сохранять прогресс в базу данных здесь (опционально)
	}

	tmpl := `
        <html>
        <head><title>Текстовый квест</title></head>
        <body>
            <h1>{{.Title}}</h1>
            <p>{{.Content}}</p>
            <form method="post">
                <input type="text" name="action" placeholder="Ваш выбор" />
                <button type="submit">Действие</button>
            </form>
        </body>
        </html>
    `

	t, err := template.New("game").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title:   "Текстовый квест",
		Content: content,
	}

	t.Execute(w, data)
}
