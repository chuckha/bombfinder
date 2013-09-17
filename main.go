package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"html/template"
	"minesweeper/game"
	"net/http"
)

const (
	port = "8080"
)

var (
	indexTmpl = template.Must(template.ParseFiles("templates/index.html"))

//	games     = make(map[*websocket.Conn]*game.Board)
)

func websocketHandler(ws *websocket.Conn) {
	newGame, _ := game.NewBoard(10, 10, 10)
	//games[ws] = newGame
	websocket.JSON.Send(ws, newGame)
	// make a new game
	// send to client
	for {
		var msg game.Message
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			break
		}
		if msg.Click == "left" {
			newGame.LeftClick(msg.Row, msg.Col)
		}
		if msg.Click == "right" {
			newGame.RightClick(msg.Row, msg.Col)
		}
		fmt.Println("Message Got: ", msg)
		websocket.JSON.Send(ws, newGame)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexTmpl.Execute(w, nil)
}

func main() {
	http.Handle("/sweep", websocket.Handler(websocketHandler))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)

	fmt.Println("Serving on port", port)
	http.ListenAndServe(":"+port, nil)
}
