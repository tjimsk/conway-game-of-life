# CONWAY'S GAME OF LIFE

## Description
This project includes a server/client implementation of Conway's Game of Life.  It is a multiplayer Web app where players can interact with a Conway GoL grid and see real-time changes.
### Basic Requirements
1. A multiplayer web app allowing players to interact with a Conway GoL grid
2. Real-time updates sent to all connected players
3. Each player assigned a random colour 
4. A set of common GoL cell patterns insertable on the grid 
5. 4 rules of Conway GoL
6. Cells coming to life following rule #4 get a colour value averaged from colour(s) of the 3 composing cells
### Rules
1. Any live cell with fewer than two live neighbours dies, as if caused by under-population.
2. Any live cell with two or three live neighbours lives on to the next generation.
3. Any live cell with more than three live neighbours dies, as if by overcrowding.
4. Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
## How to Build
The project is separated into server and client.

To build the server app, assuming `go` is already installed, you will need to get the dependencies by running `go get -d` inside the `./server` folder.  To build, run `go build` from the same folder.

To build the client app, assuming `yarn` is installed, run `yarn` to get the list of dependencies and then `yarn build`.

In development, you can build and run the server on default port `8080`, and run `yarn start` from the `./client` folder to run a hot reload version of the client that proxies Websocket requests to the server.
## How to Deploy
For deployment, the application is containerised using Docker.  Assuming `docker` is installed, you can run `docker build --tag=conway --file=Dockerfile .` from the project root.  To run the application, simply run `docker run --publish=80:8080 --name=conway conway`
## Technical/Architectural Decisions
### Basics
The primary element of this app is the `Grid`.  It has the following basic properties:

- `Generation`
- `Height`
- `Width`
- `cells` = to track cells' statuses
- `users` = to track of connected users
- `evoChan` = passes evolution updates between the grid and users
- `updateChan` = passes user driven updates between users, the grid, back to all the other users

A `Grid` may evolve and return an `Evolution`, which literally describes the evolution thru the following:

- `Cells` = list of updated cells – contains colour and active status
- `Duration` = evolution time
- `Interval` = time elapsed since last evolution
- `Generation` = current generation

The `Grid` introduces the concept of `unstable cells`, which denotes cells that are adjacent to active cells.  Those are the only cells capable of dying or coming to life.
### Server
- The server app is written in Go and makes great use of go's concurrency model using `goroutines` to synchronise communication between user connections.
- The server/client communication is done in `HTTP` and `Websocket` is used in order to avoid polling and keep the solution simple.
- The server consists of the following 3 endpoints: 
	- `/` serves the `index.html`
	- `/dist` serves assets `*.js`, `*.css`
	- `/websocket` creates a two-way Websocket connection for all updates
- The server supports the following configuration environment variables in real-time:
	- `PORT`: server listening port (default `8080`)
	- `HEIGHT`: grid height (requires app restart, default `70`)
	- `WIDTH`: grid width (requires app restart, default `120`)
	- `EVO_INTERVAL`: time interval between evolution (default `400`, in `ms`)
	- `SEED_GRID`: controls whether to seed the grid with some active data points (defaults to `true`)
- There are 5 types of messages between server/client:
	1. `MESSAGE_TYPE_USER_DETAILS`: informs client of user's color
	2. `MESSAGE_TYPE_GRID_DETAILS`: informs client of grid's dimensions and generation at request time
	3. `MESSAGE_TYPE_GRID_ACTIVE_CELLS`: informs client of all the active cells at request time
	4. `MESSAGE_TYPE_NEW_EVOLUTION`: informs client of all cell updates after an evolution; that includes cells colours as well.
	5. `MESSAGE_TYPE_CELLS_UPDATE`: informs server of a user's changes to the grid by including an array of cells and their colours; may also inform client of updates resulting from other users' actions.
### Client
- The client app is written using `Reactjs` and `SASS` for its structural clarity.  It is bundled using `Webpack` for ease of development and versatility.
- In order to prevent the `Grid` component from re-rendering entirely each time an update is sent to the client, `App` is handed a reference to each `Cell` directly.  When the `MESSAGE_TYPE_NEW_EVOLUTION` or `MESSAGE_TYPE_CELLS_UPDATE` messages are received, the `App` sets the state of updated `Cell`'s directly.  
## Technical Tradeoffs
### Deployment
- The app is deployed on DigitalOcean in a Docker container instead of a PaaS.  While a PaaS may have its advantages (simplicity, CD, etc), installing `Docker` on a server and running 2 commands to pull the image and run the container is a good tradeoff for the additional control on configuration, redundancy, restart policy, etc.  
- Having access to the underlying OS also allows runtime configuration through environment variables.
### Server
- Additional time would be spent on the following:
	- Adding an extra integration test to cover active cells following rule #4 and their assigned averaged colour.
	- The use of goroutines for concurrency is very advantageous for performance and code clarity, but its typical tradeoff is testability.  There would need to be some additional refactoring to ensure that everything can be tested, up to a satisfactory coverage level.
	- Implementing mutex while accessing cells and other grid properties.
	- Better error handling/logging: currently the server kills the socket connection and lets the client reconnect whenever any message fails to send/receive.  While this is okay in a session-less version of the app, that would need to be handled with more care otherwise.
### Client
- Additional time would be spent on the following:
	- Implementing React PropTypes on every component 
	- Declaring classes for server messages would clarify the structure of messages.
