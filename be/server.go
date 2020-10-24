package main

import (
	"log"
	server "net/http"
)

func setEndpoints() {
	// Game endpoints
	server.HandleFunc("/game/new", newGame)
	server.HandleFunc("/game/make_move", makeMove)
	server.HandleFunc("/game/get", getGame)

	// Score endpoint
	server.HandleFunc("/highscores", highScores)
}

func startServer() {
	setEndpoints()

	log.Fatal(server.ListenAndServe(":8080", nil))
}
