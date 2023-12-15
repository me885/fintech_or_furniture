package main

import (
	"log"
	"me885/fintech-or-furniture/handlers"
	"me885/fintech-or-furniture/quiz/database"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	db := database.InitDatabase("sqlite.db")
	handlersContext := &handlers.Context{DB: db}

	http.HandleFunc("/", handlers.RootPage)
	http.HandleFunc("/new-game/", handlersContext.NewGame)
	http.HandleFunc("/answer/", handlersContext.Answer)
	http.HandleFunc("/next-question/", handlersContext.NextQuestion)
	http.HandleFunc("/leaderboard/", handlersContext.Leaderboard)
	http.HandleFunc("/leaderboard-content/", handlersContext.LeaderboardTable)
	http.HandleFunc("/result/", handlersContext.EndPage)

	log.Print("Now running on http://localhost:8002")
	log.Fatal(http.ListenAndServe(":8002", nil))
}
