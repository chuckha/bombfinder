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

/*
	Given a websocket connection that has a pointer to a game
	look up all players given a game
	map[*game][]*players
	send messages to all players.
*/

var (
	indexTmpl = template.Must(template.ParseFiles("templates/index.html"))
	games     = make([]*game.Game, 0)

//	games     = make(map[*websocket.Conn]*game.Board)
)

func init() {
	/*
		for i := 0; i < 10; i++ {
			ng := game.NewGame(10, 10, 10)
			games[ng] = make([]*game.Player, 0)
		}
	*/
}

func websocketHandler(ws *websocket.Conn) {
	defer ws.Close()

	// After you connect you have to send a msg about your username
	username := ""
	// Just don't ever accept
	for username == "" {
		err := websocket.Message.Receive(ws, &username)
		if err != nil {
			break
		}
		if username == "" {
			websocket.JSON.Send(ws, &game.ErrorMessage{
				Type:  "InvalidUsername",
				Value: "Username must not be blank.",
			})
		}
	}
	websocket.JSON.Send(ws, &game.ErrorMessage{
		Type:  "ok",
		Value: username,
	})

	// A player joins
	//   If there is a game waiting for more players they join it
	// Otherwise they get a new game
	// If there are currently no games

	var currentGame *game.Game
	// If there are no games make a new one
	// Only add newly created games to the list of games.
	if len(games) == 0 {
		currentGame = game.NewGame(1, 10, 10, 10)
		games = append(games, currentGame)
	} else {
		// If there are some games
		for i := 0; i < len(games); i++ {
			// If the game has fewer players than required that is the one we want
			if len(games[i].Players) < games[i].NumPlayers {
				currentGame = games[i]
			}
		}
		// If all games have been filled we need a new game
		if currentGame == nil {
			currentGame = game.NewGame(2, 10, 10, 10)
			games = append(games, currentGame)
		}
	}
	// Create a new player
	player := game.NewPlayer(ws, username)
	// holy fuck. This line is so epic.
	defer currentGame.RemovePlayer(player)

	// Add them to the game
	currentGame.AddPlayer(player)

	if currentGame.IsFull() {
		currentGame.SendBoard()
	} else {
		currentGame.SendInfo()
	}
	currentGame.SendPlayerInfo()

	for {
		var msg game.InMessage
		// Block on receiving a message
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			break
		}

		switch msg.Type {
		case "move":
			// Don't do anything to a finished board
			if currentGame.Board.Finished {
				continue
			}
			val := msg.Value
			switch val.Click {
			case "left":
				currentGame.Board.LeftClick(player, val.Row, val.Col)
			case "right":
				currentGame.Board.RightClick(player, val.Row, val.Col)
			default:
				// A weird move will just be ignored
				continue
			}
			currentGame.SendBoard()
			if currentGame.Board.Finished {
				message := "Game over :("
				if currentGame.Board.Won {
					message = "Victory!"
				}
				fmt.Println("board is finished")
				websocket.JSON.Send(ws, &game.OutMessage{
					Type:  "status",
					Value: message,
				})
			}
		default:
			continue
		}
	}
}

// Handle the different messages from the JS client.
// newGame should give them a new game
func handleMessage(msg *game.InMessage, ws *websocket.Conn, board *game.Board) {
	fmt.Println("Message Got: ", msg)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexTmpl.Execute(w, games)
}

func main() {
	http.Handle("/sweep", websocket.Handler(websocketHandler))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)

	fmt.Println("Serving on port", port)
	http.ListenAndServe(":"+port, nil)
}
