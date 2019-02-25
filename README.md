## CONWAY'S GAME OF LIFE

#### General

The Game of Life is a cellular automaton devised by John Horton Conway in 1970 where you create cellular patterns on a grid and observe them evolve over time, following 4 basic rules:

- Any live cell with fewer than two live neighbors dies, as if by underpopulation.
- Any live cell with two or three live neighbors lives on to the next generation.
- Any live cell with more than three live neighbors dies, as if by overpopulation.
- Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction.

Conway's conjecture was that no pattern can grow without limit.

The result can be found [here](http://life.jimsk.com).  This is the design notes of building this application.

#### Design Notes

![Design](/assets/conway-design.jpg)

##### Frontend

The frontend consists of a wrapper component and 3 components.  The breakdown of each component is as follows:

Component    | Responsibilities
-------------|------------------
`App` | <ul><li>Maintain server websocket connection alive</li><li>Receive evolution updates and update `Grid`.</li><li>Synchronize playback state with server.</li><li>Define handlers for wrapped components.</li><li>Pass state to wrapped components.</li></ul>
`Toolbar` | <ul><li>Render list of patterns.</li><li>Render evolution playback controls.</li></ul>
`Grid` | <ul><li>Render cells.</li></ul>
`Status` | <ul><li>Render connection status, generation ID, evolution rate, transfer rate, etc</li></ul>

##### Backend

The backend server endpoints are as follows:

Endpoint    | Description
------------|-------------
`GET /` | <ul><li>Serve static files.</li></ul>
`POST /activate` | <ul><li>Activate a set of points on the grid.</li><li>Push changes to all connected users via `/websocket`.</li></ul>
`POST /deactivate` | <ul><li>Deactivate a set of points on the grid.</li><li>Push changes to all connected users via `/websocket`.</li></ul>
`POST /interval` | <ul><li>Update the evolution interval.</li><li>Push changes to all connected users via `/websocket`.</li></ul>
`POST /reset` | <ul><li>Reset the grid.</li><li>Push changes to all connected users via `/websocket`.</li></ul>
`GET /websocket` | <ul><li>Push messages to connected users.</li></ul>


#### Discussion

For discussion sake some of the snippets below are simplified from the actual code.

###### Evolution Algorithm

The evolution algorithm from state `A` to `B` can be expressed simply as:

`f(A) = B`

where `A` and `B` are the states of cells on the grid, respectively before and after an evolution, and `f` is the evolution function.  Note that we can define another function `f'(A)` and `g(X)` such that:

`f(A') = f(g(A)) = f(A) = B <=> g(A) = A'`

where `g(A)` is the function that returns `A'` which is a subset of `A` with only its live cells and their neighbors.  We note that in order to evolve any given state we only need to apply the evolution function on the live cells and their neighbors of that state.  

We name the `g(A)` function `unstableCells()`, which takes a set of live cells, appends all their adjacents cells and then returns it:

<pre class="prettyprint lang-go">
type Point struct {
	X int
	Y int
}
func unstableCells(liveCells map[Point]struct{}) map[Point]struct{} {
	unstables := map[Point]struct{}{}

	for p, _ := range liveCells {
		unstables[p] = struct{}{}
		for _, ap := range adjacentPoints(p) {
			unstables[ap] = struct{}{}
		}
	}
	return unstables	
}
func adjacentPoints(p Point) (pts []Point) {
	return []Point{
		Point{p.X - 1, p.Y - 1},
		Point{p.X, p.Y - 1},
		Point{p.X + 1, p.Y - 1},
		// etc...
	}
}
</pre>

We note that for each cell in `A'`we must compute the number of neighboring live cell along with that cell's live status in order to compute that cell's next generation status.  We call that function `nextState()`:

<pre class="prettyprint lang-go">
func nextState(p Point, liveCells map[Point]struct{}) bool {
	count := 0
	for _, ap := range adjacentPoints(p) {
		if _, ok := liveCells[ap]; ok {
			count++
		}
	}
	if _, alive := liveCells[p]; alive {
		// condition #1 & #3: cell dies
		if count < 2 || count > 3 {
			return false

		//condition #2: cell stays alive
		} else if count == 2 || count == 3 {
			return true
		}
	} else {
		// condition #4: cell comes to life
		if count == 3 {
			return true
		}
	}
	return false	
}
</pre>

We note further that `unstableCells()` has time complexity `O(n)` for `n` number of live cells.  

In a usability standpoint, we must set an evolution rate limit proportional to the grid size.  At this time, with 60ms per evolution and a 75x90 cells grid configuration the app behaves without any problem.

###### Pushing the State to All Websocket Connections

The implementation of the websocket connection handler using packages `net/http` and `github.com/gorilla/websocket` looks like this:

<pre class="prettyprint lang-go">
var connections = map[*websocker.Conn](chan struct{}){}

func (h *ConnectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// upgrade the http connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// keep reference of connection using a map of type
	// map[*websocket.Conn](chan struct{}) so we can push to it
	// use custom struct type for message to pass.
	msgChan := make(chan struct{})
	endChan := make(chan bool)
	connections[conn] = msgChan
	defer delete(connections, conn)

	// listening loop to push new messages on
	go func(c *websocket.Conn, mChan chan struct{}, eChan chan bool) {
		for err := error(nil); err == nil; {
			err = conn.WriteJSON(<-mChan)
		}
		// error handling omitted here
		eChan <- true
	}(conn, msgChan, endChan)

	// block until terminated
	<-endChan
}
</pre>

To push updates to all connections we iterate over the `connections` map whose values conveniently are message channels, which accept any type we want to define.  We make sure to run each message on a separate goroutine and implement a timeout so it doesn't block and it ensures the goroutine always returns.

<pre class="prettyprint lang-go">
func PushState(state struct{}) {
	for _, msgChan := range connections {
		go func() {
			select {
			case msgChan <- state:
			case <-time.After(2 * time.Second):
			}
		}()
	}
}
</pre>

#### Conclusion

While there could be more improvements regarding performance and scalability to make, when designing an app or a game it is more often than not that performance is secondary to it being fun and interactive.  So I won't be getting into more details for this particular project as, frankly, I have just grown tired of this game. :)


