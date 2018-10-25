import React from "react"

import {Toolbox} from "Components/Toolbox"
import {Grid} from "Components/Grid"
import {Status} from "Components/Status"
import {ConnectBackground} from "Components/ConnectBackground"

import {newClient} from "Util/websocket"

const styles = require("Components/App.scss")

const MESSAGE_TYPE_USER_DETAILS = 0
const MESSAGE_TYPE_GRID_DETAILS = 1
const MESSAGE_TYPE_GRID_ACTIVE_CELLS = 2
const MESSAGE_TYPE_GRID_UPDATE = 3
const MESSAGE_TYPE_UPDATE_CELLS = 4

class App extends React.Component {
	render() {
		return (
			<div className={styles.app}>
				{this.state.readyState !== 1 ? <ConnectBackground /> : null}

				{this.state.readyState !== 1 ? null :
				<Toolbox
					ref={ref => {this.toolbox = ref}}
					grid={this.state.grid}
					user={this.state.user}
					selectedToolIndex={this.state.selectedToolIndex}
					onClickTool={this.onClickTool.bind(this)} />}

				<Grid grid={this.state.grid}
					user={this.state.user}
					cells={this.state.cells}
					onClickCell={this.onClickCell.bind(this)}
					receivedActiveCells={this.state.receivedActiveCells}
					setActiveColorFunc={this.setActiveColorFunc} />

				<Status
					grid={this.state.grid}
					readyState={this.state.readyState}
					dataSizeReceived={this.state.dataSizeReceived} />
			</div>)
	}

	constructor(props) {
		super(props)
		this.state = this.getInitialState()
		this.setActiveColorFunc = {}
		this.createWebsocketConnection()
	}

	getInitialState() {
		return {
			cells: {},
			grid: {height: 0, width: 0, generation: 0},
			readyState: -1,
			dataSizeReceived: 0,
			receivedActiveCells: false,
			user: {name: "", color: {r: 0, g: 0, b: 0}},
			selectedToolIndex: 0
		}
	}

	createWebsocketConnection() {
		this.client = newClient()

		this.client.onopen = this.onWebSocketOpen.bind(this)
		this.client.onclose = this.onWebSocketClose.bind(this)
		this.client.onmessage = this.onWebSocketMessage.bind(this)
	}

	onWebSocketOpen() {
		this.checkWebSocketState()
	}

	checkWebSocketState() {
		setTimeout(() => {
			this.setState({
				readyState: this.client.readyState
			}, this.checkWebSocketState.bind(this))
		}, 500)
	}

	onWebSocketClose() {
		this.setState(this.getInitialState())

		setTimeout(this.createWebsocketConnection.bind(this), 4000)
	}

	onWebSocketMessage(e) {
	    if (typeof e.data !== "string") {return}

        var msg = JSON.parse(e.data)
	    var size = e.data.length / 1000
        // console.log("WebSocket message:", msg)

        switch(msg.t) {
    	case MESSAGE_TYPE_USER_DETAILS:
    		this.setUserDetails(msg, size); break
    	case MESSAGE_TYPE_GRID_DETAILS:
    		this.renderGrid(msg, size); break
		case MESSAGE_TYPE_GRID_ACTIVE_CELLS:
			this.setInitialActiveCells(msg, size); break
		case MESSAGE_TYPE_GRID_UPDATE:
			this.applyGridUpdates(msg, size); break
		default:
			break
        }
	}

	setUserDetails(msg) {
		this.setState({user: msg.c})
	}

	renderGrid(msg) {
		var g = Object.assign({}, this.state.grid, msg.c)

		// create `cells` object
		// all cells inactive with null color
		var cells = {}
		for (var i = 1; i <= g.width; i++) {
			for (var j = 1; j <= g.height; j++) {
				var k = `${i};${j}`
				cells[k] = {x: i, y: j, k: k}
			}
		}

		this.setState({
			grid: g,
			cells: cells
		})
	}

	setInitialActiveCells(msg) {
		this.setState({
			receivedActiveCells: true // prevent Grid update
		}, this.updateCells.bind(this, msg.c))
	}

	applyGridUpdates(msg, size) {
		var g = Object.assign({}, this.state.grid)
		g.generation = msg.c.generation

		this.setState({
			grid: g,
			dataSizeReceived: size
		}, this.updateCells.bind(this, msg.c.cells))
	}

	updateCells(cells) {
		for (var i = 0; i < cells.length; i++) {
			var c = cells[i]
			var k = `${c.x};${c.y}`

			var setActiveColorFunc = this.setActiveColorFunc[k]
			if (setActiveColorFunc != null) {
				setActiveColorFunc(c.a ? c.c : null)
			}
		}
	}

	onClickTool(index) {
		this.setState({
			selectedToolIndex: index
		})
	}

	onClickCell(c) {
		var msg = {
			t: MESSAGE_TYPE_UPDATE_CELLS,
			c: [],
			u: this.state.user
		}

		msg.c = this.toolbox.getCellsFromSelectedTool(c.x, c.y)

		var cellsArr = []
		Object.keys(msg.c).map((k) => {
			var c = msg.c[k]
			c.a = true
			c.c = this.state.user.color

			cellsArr.push(c)
		})

		this.updateCells(cellsArr)

		this.client.send(JSON.stringify(msg))
	}
}

export {
	App
}