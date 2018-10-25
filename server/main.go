package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

const (
	DefaultPort        = ":8080"
	DefaultAutoEvo     = true
	DefaultEvoInterval = 500 * time.Millisecond
	DefaultGridWidth   = 120
	DefaultGridHeight  = 70
)

const (
	chanTimeout = 2 * time.Second
)

var (
	_  fmt.Stringer
	_  log.Logger
	wd string

	grid     *Grid
	logs     = []*Message{}
	users    = map[string]*User{}
	mu       = &sync.Mutex{}
	upgrader = websocket.Upgrader{
		EnableCompression: true,
		Subprotocols:      []string{"Sec-WebSocket-Protocol", "echo-protocol"},
		CheckOrigin:       func(r *http.Request) bool { return true },
	}
)

func main() {
	wd, _ = os.Getwd()

	parseEnv()
	initGrid()

	go startEvolve(grid)

	http.HandleFunc("/", HandleRoot)
	http.HandleFunc("/dist/", HandleAssets)
	http.HandleFunc("/websocket", HandleWebSocket)
	log.Println("listening on", viper.GetString("port"))
	log.Fatal(http.ListenAndServe(viper.GetString("port"), nil))
}

func parseEnv() {
	viper.SetDefault("port", DefaultPort)
	viper.SetDefault("height", DefaultGridHeight)
	viper.SetDefault("width", DefaultGridWidth)

	viper.BindEnv("port")
	viper.BindEnv("height")
	viper.BindEnv("width")
}

func initGrid() {
	grid = NewGrid(viper.GetInt("width"), viper.GetInt("height"), DefaultAutoEvo, DefaultEvoInterval)
	grid.seed()

	log.Printf("grid is %vx%v\n", viper.GetInt("width"), viper.GetInt("height"))
}

func startEvolve(g *Grid) {
	for {
		// start := time.Now()
		gu := g.Evolve()
		// end := time.Now()
		// diff := end.Sub(start)

		// log.Printf("\tGENERATION=%v;LIVE=%v;UPDATES=%v;TIME=%v;INTERVAL=%v\n",
		// g.Generation, len(g.activeCells()), len(gu.Cells), diff, g.EvoInterval)

		broadcastGridUpdate(gu, users)

		time.Sleep(g.EvoInterval)
	}
}

func broadcastGridUpdate(gu GridUpdate, _users map[string]*User) {
	for _, u := range _users {
		go func(_u *User) {
			select {
			case _u.gridUpdateChan <- gu:
			case <-time.After(chanTimeout):
				UnregisterUser(_u.Name)
				log.Printf("connection timeout for user %v; %v users remaining\n", _u.Name, len(users))
				_u.endChan <- true
			}
		}(u)
	}
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleRoot: %v %v\n", r.Method, r.URL.Path)

	p := path.Join(wd, "../client/dist/index.html")

	http.ServeFile(w, r, p)
}

func HandleAssets(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleAssets: %v %v\n", r.Method, r.URL.Path)

	p := path.Join(wd, "../client", r.URL.Path)

	http.ServeFile(w, r, p)
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleWebSocket:%v:%v:%v", r.Method, r.URL.Path, websocket.Subprotocols(r))

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("connection upgrade error:", err)
		return
	}

	// experimental feature in github.com/gorilla/websocket
	// may impact performance
	// TODO: benchmarks
	conn.EnableWriteCompression(true)

	defer conn.Close()

	// create user or resume existing user
	u := NewUser(conn)
	RegisterUser(u)

	// send user to client
	if err := u.SendUserDetails(); err != nil {
		log.Println("error sending user details:", err)
		return
	}

	// send grid to client
	if err := u.SendGridDetails(grid); err != nil {
		log.Println("error sending grid details:", err)
		return
	}

	// send grid active cells to websocket client
	if err := u.SendGridActiveCells(grid); err != nil {
		log.Println("error sending grid active cells:", err)
		return
	}

	// receive and send grid updates loop
	go func(_u *User) {
		for {
			if err := u.SendGridUpdate(<-u.gridUpdateChan); err != nil {
				log.Printf("error writing message json: %v [closing connection user:%v]", err, u.Name)
				u.endChan <- true
				break
			}
		}
	}(u)

	// receive cell updates from websocket client
	go func(_u *User, g *Grid) {
		for {
			msg := ActivateCellsMessage{}
			if err := conn.ReadJSON(&msg); err != nil {
				log.Printf("error reading message json: %v [closing connection user:%v]", err, u.Name)
				u.endChan <- true
				break
			}

			cellsMap := map[string]Cell{}
			for _, c := range msg.Cells {
				k := fmt.Sprintf(`%v;%v`, c.X, c.Y)
				cellsMap[k] = c
			}

			cells := g.ActivateCells(cellsMap, _u)

			broadcastGridUpdate(g.UpdateFromCells(cells), users)
		}
	}(u, grid)

	<-u.endChan
}
