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
	DefaultEvoInterval = 400
	DefaultGridWidth   = 120
	DefaultGridHeight  = 70
	DefaultSeedGrid    = true
	DefaultLogEvoTime  = false
	chanTimeout        = 2 * time.Second
)

var (
	_        fmt.Stringer
	_        log.Logger
	wd       string
	grid     *Grid
	users    = map[string]User{}
	mu       = &sync.Mutex{}
	upgrader = websocket.Upgrader{
		EnableCompression: true,
		Subprotocols:      []string{"Sec-WebSocket-Protocol", "echo-protocol"},
		CheckOrigin:       func(r *http.Request) bool { return true },
	}
)

func main() {
	wd, _ = os.Getwd()

	// set default environment variables
	viper.SetDefault("evoInterval", DefaultEvoInterval)
	viper.SetDefault("seedGrid", DefaultSeedGrid)
	viper.SetDefault("port", DefaultPort)
	viper.SetDefault("height", DefaultGridHeight)
	viper.SetDefault("width", DefaultGridWidth)

	// bind environment variables
	viper.BindEnv("evoInterval")
	viper.BindEnv("seedGrid")
	viper.BindEnv("port")
	viper.BindEnv("height")
	viper.BindEnv("width")

	// initialize grid
	grid = NewGrid(viper.GetInt("width"), viper.GetInt("height"))
	if viper.GetBool("seedGrid") {
		grid.seed()
	}

	go grid.StartEvolutions()

	// listen on evoChan for evolution signals
	go func() {
		for {
			e := <-grid.evoChan

			for _, u := range users {
				go func(u User) {
					// time out and clear connection if channel data is not received
					select {
					case u.evoChan <- e:
					case <-time.After(chanTimeout):
					}
				}(u)
			}
		}
	}()

	// listen on updateChan for cells update signals
	go func() {
		for {
			cells := <-grid.updateChan

			for _, u := range users {
				go func(u User) {
					// time out and clear connection if channel data is not received
					select {
					case u.updateChan <- cells:
					case <-time.After(chanTimeout):
					}
				}(u)
			}
		}
	}()

	// listen and serve http
	http.HandleFunc("/", HandleRoot)
	http.HandleFunc("/dist/", HandleAssets)
	http.HandleFunc("/websocket", HandleWebSocket)

	log.Println("listening on", viper.GetString("port"))
	log.Fatal(http.ListenAndServe(viper.GetString("port"), nil))
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
		log.Print("upgrader.Upgrade:", err)
		return
	}
	conn.EnableWriteCompression(true)
	defer conn.Close()

	// create user or resume existing user
	u := NewUser(conn)
	RegisterUser(u)

	if err := WriteUserDetails(conn, u); err != nil {
		log.Println("WriteUserDetails:", err)
		return
	}

	if err := WriteGridDetails(conn, grid); err != nil {
		log.Println("WriteGridDetails:", err)
		return
	}

	if err := WriteGridActiveCells(conn, grid.activeCells()); err != nil {
		log.Println("WriteGridActiveCells:", err)
		return
	}

	// receive grid evolution signals
	// write grid evolution updates to websocket
	go func(u User) {
		for {
			evo := <-u.evoChan

			if err := WriteEvolution(conn, evo); err != nil {
				log.Println("WriteEvolution:", err)
				u.closeChan <- true
				break
			}
		}
	}(u)

	// receive grid cells update signals
	// write grid cell updates to websocket
	go func(u User) {
		for {
			cells := <-u.updateChan

			if err := WriteCellsUpdate(conn, cells); err != nil {
				log.Println("WriteCellsUpdate:", err)
				u.closeChan <- true
				break
			}
		}
	}(u)

	// receive grid cells update websocket messages
	// send cells update signal to grid
	go func(u User) {
		for {
			msg, err := ReadCellsUpdate(conn)
			if err != nil {
				log.Println("ReadCellsUpdate:", err)
				u.closeChan <- true
				break
			}

			// apply received cells update
			cells := []*Cell{}
			for _, mc := range msg.Cells {
				c, err := grid.CellAtPoint(Point{mc.X, mc.Y})
				if err != nil {
					continue
				}

				c.active = mc.Active

				if mc.Active {
					c.color = u.Color
				} else {
					c.color = Color{}
				}

				c.Flush()
				cells = append(cells, c)
			}

			// send cells update signal to grid
			select {
			case grid.updateChan <- cells:
			case <-time.After(chanTimeout):
				UnregisterUser(u.Name)
				log.Println("grid.updateChan timeout")
				u.closeChan <- true
			}
		}
	}(u)

	<-u.closeChan
}
