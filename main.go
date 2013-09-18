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
	board, _ := game.NewBoard(10, 10, 10)
	for {
		var msg game.InMessage
		// Block on receiving a message
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			break
		}

		switch msg.Type {
		case "newGame":
			fmt.Println("Received newGame")
			board, _ = game.NewBoard(10, 10, 10)
			websocket.JSON.Send(ws, &game.OutMessage{
				Type:  "board",
				Value: board,
			})
		case "move":
			// Don't do anything to a finished board
			if board.Finished {
				continue
			}
			val := msg.Value
			if val.Click == "left" {
				board.LeftClick(val.Row, val.Col)
			}
			if val.Click == "right" {
				board.RightClick(val.Row, val.Col)
			}
			websocket.JSON.Send(ws, &game.OutMessage{
				Type:  "board",
				Value: board,
			})
			if board.Finished {
				message := "Game over :("
				if board.Won {
					message = "Victory!"
				}
				fmt.Println("board is finished")
				websocket.JSON.Send(ws, &game.OutMessage{
					Type:  "status",
					Value: message,
				})
			}
		default:
			return
		}
	}
}

// Handle the different messages from the JS client.
// newGame should give them a new game
func handleMessage(msg *game.InMessage, ws *websocket.Conn, board *game.Board) {
	fmt.Println("Message Got: ", msg)

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
