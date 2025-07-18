package handlers

import (
	"html/template"
	"net/http"
	"warrior-of-rohan/internal/models"
)

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.New("game").Parse(`
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
    `)
	if err != nil {
		panic(err)
	}
}

func HandleGame(w http.ResponseWriter, r *http.Request) {
	// Можно добавить логику получения состояния из базы или сессии
	content := "Вы находитесь в темной комнате. Что вы делаете? (введите 'осмотреться' или 'идти дальше')"

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
	}

	data := models.PageData{
		Title:   "Текстовый квест",
		Content: content,
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
