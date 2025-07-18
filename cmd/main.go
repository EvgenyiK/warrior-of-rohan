package main

import (
	"log"
	"net/http"
	"sync"
	"warrior-of-rohan/internal/game"
	"warrior-of-rohan/internal/handlers"
)

func main() {
	var wg sync.WaitGroup

	// Запуск веб-сервера
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.HandleFunc("/", handlers.HandleGame)
		log.Println("Starting web server at :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Запуск игры на Ebiten
	wg.Add(1)
	go func() {
		defer wg.Done()
		game.RunEbiten()
	}()

	wg.Wait()
}
