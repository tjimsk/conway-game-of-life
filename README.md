# CONWAY'S GAME OF LIFE

## Rules

1. Any live cell with fewer than two live neighbours dies, as if caused by under-population.
2. Any live cell with two or three live neighbours lives on to the next generation.
3. Any live cell with more than three live neighbours dies, as if by overcrowding.
4. Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

## Description

### Server

The server is written in Go.  It is made of 3 endpoints: root (/), assets (/dist), websocket (/websocket).  Some parameters are read and set from environment variables, such as PORT, HEIGHT, WIDTH.  Each Websocket connection is tied to a user who is assigned a random color on creation.  Go channels are used to concurrently pass grid evolution updates to each registered user.  They are also used in order to transmit updates from a single client message to every registered users' websocket connection.

The grid is seeded with 2 shapes.

### Client

The client uses the Reactjs framework.  The main components are App, Grid, Cell, Toolbox, Status.  In order to optimize the rendering performance, only updated cells are re-rendered.  Hence, the Grid component does no longer update once the initial grid cells are rendered.  Instead, the App component owns references to a wrapper of setState of each cell and only calls that.

The Toolbox component allows users to select a drawing tool to place on the grid.  It is implemented by passing the clicked location's coordinate to a function that returns an array of cell coordinates that should be activated based on the selected tool.

### Deployment

The app is deployed on a single core droplet with 1GB RAM and should be accessible at http://178.128.221.183.

The Dockerfile can be pulled from dockerhub:

`docker pull docker.io/jimlearnstofly/terminal1`

To run the container:

`docker run --rm --name=conway --publish=8080:8080 jimlearnstofly/terminal1`

