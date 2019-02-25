package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"

	"github.com/tjimsk/life"
)

var (
	grid     *life.Grid
	upgrader = websocket.Upgrader{
		EnableCompression: true,
		Subprotocols:      []string{"Sec-WebSocket-Protocol", "echo-protocol"},
		CheckOrigin:       func(r *http.Request) bool { return true },
	}
)

func main() {
	log.SetFlags(log.Lshortfile)

	viper.SetDefault("interval", 1000)
	viper.SetDefault("port", ":8080")
	viper.SetDefault("static", "../../client/dist/development")

	viper.BindEnv("interval")
	viper.BindEnv("port")
	viper.BindEnv("static")

	// initialize a grid and start evolution loop
	grid = life.NewGrid(viper.GetInt("interval"))
	go func() {
		for {
			for grid.Interval == -1 || grid.NoConnectedUser() {
				time.Sleep(100 * time.Millisecond)
			}

			grid.Evolve()
			grid.PushStateChange()

			time.Sleep(time.Duration(grid.Interval) * time.Millisecond)
		}
	}()

	http.HandleFunc("/", HandleStatic)
	http.HandleFunc("/activate", HandleActivate)
	http.HandleFunc("/deactivate", HandleDeactivate)
	http.HandleFunc("/interval", HandleInterval)
	http.HandleFunc("/reset", HandleReset)
	http.HandleFunc("/websocket", HandleWebSocket)

	log.Println("listening on", viper.GetString("port"))
	log.Fatal(http.ListenAndServe(viper.GetString("port"), nil))
}

func HandleStatic(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleStatic: %v %v\n", r.Method, r.URL.Path)
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(w, r, path.Join(viper.GetString("static"), r.URL.Path))
}

func HandleActivate(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleActivate: %v %v\n", r.Method, r.URL.Path)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	msg := life.RequestActivateMessage{}
	if err := json.Unmarshal(b, &msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if grid.PlayerConnected(msg.Player) {
		grid.Activate(msg.Points, msg.Player)
		grid.PushStateChange()
	}
}

func HandleDeactivate(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleDeactivate: %v %v\n", r.Method, r.URL.Path)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	msg := life.RequestDeactivateMessage{}
	if err := json.Unmarshal(b, &msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if grid.PlayerConnected(msg.Player) {
		grid.Deactivate(msg.Point, msg.Player)
		grid.PushStateChange()
	}
}

func HandleInterval(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleInterval: %v %v\n", r.Method, r.URL.Path)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	msg := life.RequestIntervalMessage{}
	if err := json.Unmarshal(b, &msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("HandleInterval:", msg)

	if grid.PlayerConnected(msg.Player) {
		grid.SetInterval(msg.Interval, msg.Player)
		grid.PushIntervalChange(msg.Player)
	}
}

func HandleReset(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleReset: %v %v\n", r.Method, r.URL.Path)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	msg := life.RequestResetMessage{}
	if err := json.Unmarshal(b, &msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("HandleReset:", msg)

	if grid.PlayerConnected(msg.Player) {
		grid.Reset(msg.Player)
		grid.PushStateChange()
	}
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleWebSocket:%v:%v:%v", r.Method, r.URL.Path, websocket.Subprotocols(r))

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrader.Upgrade:", err)
		return
	}
	defer conn.Close()

	player := grid.AddPlayer(conn)
	defer grid.RemovePlayer(player)

	go player.ListenMessages()

	grid.PushPlayer(player)
	grid.PushState(player)

	<-player.DisconnectChan
}
