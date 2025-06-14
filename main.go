package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
)

var db *sql.DB

// Структура для передачи данных в шаблон
type PageData struct {
	Title   string
	Content string
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	host := "localhost"
	port := 5432
	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASS")
	dbname := os.Getenv("PGDB")

	// Формируем строку подключения
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

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
