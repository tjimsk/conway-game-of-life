package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

var (
	autoEvoInterval = 5 * time.Second    // use mutex
	autoEvoEnabled  = true               // use mutex
	grid            *Grid                // use mutex
	events          = []*Event{}         // use mutex
	users           = map[string]*User{} // use mutex
	mu              = &sync.Mutex{}
	upgrader        = websocket.Upgrader{
		Subprotocols: []string{"Sec-Websocket-Protocol"},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	viper.SetDefault("port", ":8080")
	viper.BindEnv("port")
	viper.SetDefault("height", 360)
	viper.BindEnv("height")
	viper.SetDefault("width", 600)
	viper.BindEnv("width")

	http.HandleFunc("/websocket", HandleWebSocket)

	grid = NewGrid(viper.GetInt("width"), viper.GetInt("height"))
	grid.seed()

	go automaticUpdateGrid(grid)

	go broadcastEvents()

	log.Println("listening on", viper.GetString("port"))
	log.Printf("grid is %vx%v\n", viper.GetInt("width"), viper.GetInt("height"))

	http.ListenAndServe(viper.GetString("port"), nil)
}

func automaticUpdateGrid(g *Grid) {
	for {
		update := g.nextGeneration()

		log.Printf("Updating %v cells after %v\n", len(update.Cells), autoEvoInterval)

		// notify users
		for _, u := range users {
			u.gridChan <- update
		}

		time.Sleep(autoEvoInterval)
	}
}

func broadcastEvents() {
	for {

	}
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleWebSocket:%v:%v", r.Method, r.URL.Path)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer conn.Close()

	// create user
	u := NewUser(fmt.Sprintf("player%v", len(users)+1))
	users[u.Name] = u

	// send user
	conn.WriteJSON(Message{
		Type:    OUTBOUND_MESSAGE_TYPE_USER,
		Content: u,
	})

	// send grid
	conn.WriteJSON(Message{
		Type:    OUTBOUND_MESSAGE_TYPE_GRID,
		Content: grid,
	})

	// send active cells
	conn.WriteJSON(Message{
		Type:    OUTBOUND_MESSAGE_TYPE_GRID_ACTIVE_CELLS,
		Content: grid.activeCells(),
	})

	// send periodic cell updates
	go func() {
		for {
			updates := <-u.gridChan

			conn.WriteJSON(Message{
				Type:    OUTBOUND_MESSAGE_TYPE_GRID_UPDATE,
				Content: updates,
			})
		}
	}()

	// send events updates
	go func() {
		for <-u.eventChan {
			conn.WriteJSON(Message{
				Type:    OUTBOUND_MESSAGE_TYPE_EVENT_UPDATE,
				Content: events[len(events)-1],
			})
		}
	}()

	for {
		// mt, message, err := conn.ReadMessage()
		// if err != nil {
		// 	// log.Println("read:", err)
		// 	break
		// }
		// // log.Printf("recv: %s", message)
		// err = conn.WriteMessage(mt, message)
		// if err != nil {
		// 	// log.Println("write:", err)
		// 	break
		// }
		time.Sleep(100 * time.Millisecond)
	}
}
