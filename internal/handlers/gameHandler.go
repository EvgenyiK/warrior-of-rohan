package handlers

import (
	"fmt"
	"net/http"
)

func HandleGame(dataChan chan<- string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.URL.Query().Get("data")

		if data != "" {
			select {
			case dataChan <- data:
				fmt.Println("send chanal")
			default:
				fmt.Println("wait")
			}
		}

		w.Write([]byte("Data received"))
	}
}
