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

	viper.SetDefault("interval", 200)
	viper.SetDefault("seed", true)
	viper.SetDefault("port", ":8080")
	viper.SetDefault("height", 64)
	viper.SetDefault("width", 120)
	viper.SetDefault("static", "../../client/dist/")

	viper.BindEnv("interval")
	viper.BindEnv("seed")
	viper.BindEnv("port")
	viper.BindEnv("height")
	viper.BindEnv("width")
	viper.BindEnv("static")

	// set up and start grid
	grid = life.NewGrid(viper.GetInt("width"), viper.GetInt("height"))
	if viper.GetBool("seed") {
		seed()
	}

	go func() {
		for {
			for grid.Paused {
				time.Sleep(10 * time.Millisecond)
			}
			grid.Evolve()
			grid.PushStateChange()
			time.Sleep(time.Duration(viper.GetInt("interval")) * time.Millisecond)
		}
	}()

	http.HandleFunc("/", HandleStatic)
	http.HandleFunc("/pause", HandlePause)
	http.HandleFunc("/activate", HandleActivate)
	http.HandleFunc("/deactivate", HandleDeactivate)
	http.HandleFunc("/interval", HandleInterval)
	http.HandleFunc("/websocket", HandleWebSocket)

	log.Println("listening on", viper.GetString("port"))
	log.Fatal(http.ListenAndServe(viper.GetString("port"), nil))
}

func seed() {
	ps := []life.Point{
		// blinker
		life.Point{X: 2, Y: 3},
		life.Point{X: 3, Y: 3},
		life.Point{X: 4, Y: 3},
		// beacon
		life.Point{X: 7, Y: 2},
		life.Point{X: 7, Y: 3},
		life.Point{X: 8, Y: 2},
		life.Point{X: 9, Y: 5},
		life.Point{X: 10, Y: 4},
		life.Point{X: 10, Y: 5},
	}
	grid.Activate(ps, life.NewPlayer(nil))
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

func HandlePause(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandlePause: %v %v\n", r.Method, r.URL.Path)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	msg := life.RequestPauseMessage{}
	if err := json.Unmarshal(b, &msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if grid.PlayerConnected(msg.Player) {
		grid.SetPause(msg.Pause, msg.Player)
		grid.PushPauseChange(msg.Player)
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

	grid.PushPlayer(player)
	grid.PushState(player)

	<-player.DisconnectChan
}
