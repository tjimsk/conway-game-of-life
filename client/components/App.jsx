import React from "react"
import axios from "axios"

import {Toolbox} from "Components/Toolbox"
import {Toolbar} from "Components/Toolbar"
import {Grid} from "Components/Grid"
import {Status} from "Components/Status"
import {ConnectBackground} from "Components/ConnectBackground"

import {newClient} from "Util/websocket"

const styles = require("Components/App.scss")

class App extends React.Component {
	render() {
		return (
			<div className={styles.app}>
				{this.state.readyState !== 1 ? <ConnectBackground /> : null}

				<Toolbar
					minInterval={100}
					maxInterval={2000}
					paused={this.state.paused}
					toolbarRef={this.toolbarRef}
					onClickPauseHandler={this.onClickPauseHandler.bind(this)}
					onChangeIntervalHandler={this.onChangeIntervalHandler.bind(this)}
					onClickResetHandler={this.onClickResetHandler.bind(this)}
					onClickPattern={this.onClickPattern.bind(this)} />

				<Grid 
					width={this.state.width}
					height={this.state.height}
					cellRefs={this.cellRefs}
					onClickCell={this.onClickCell.bind(this)} />

				<Status
					player={this.state.player}
					readyState={this.state.readyState}
					generation={this.state.generation}
					dataSizeReceived={this.state.dataSizeReceived} />
			</div>)
	}

	constructor(props) {
		super(props)
		this.state = this.getInitialState()
		this.cellRefs = {}
		this.toolbarRef = {}
		this.activeCells = []
		this.createWebsocketConnection()
	}

	getInitialState() {
		return {
			width: 0,
			height: 0,
			generation: 0,
			player: {name: "", color: {r: 0, g: 0, b: 0}},
			dataSizeReceived: 0,
			readyState: -1,

			selectedToolIndex: 0,
		}
	}

	createWebsocketConnection() {
		this.client = newClient()
		this.client.onopen = this.onWebSocketOpen.bind(this)
		this.client.onclose = this.onWebSocketClose.bind(this)
		this.client.onmessage = this.onWebSocketMessage.bind(this)
	}

	onWebSocketOpen() {
		setTimeout(() => {
			this.setState({
				readyState: this.client.readyState
			}, this.onWebSocketOpen.bind(this))
		}, 1000)
	}

	onWebSocketClose() {
		this.setState(this.getInitialState())
		setTimeout(this.createWebsocketConnection.bind(this), 500)
	}

	onWebSocketMessage(e) {
        var msg = JSON.parse(e.data)
	    var size = e.data.length / 1000

	    if (msg.player != null && Object.keys(msg).length == 1) {
			this.setState({
				player: msg.player
			})    		
	    } else if (msg.activeCells != null) {
	    	this.setState({
	    		width: msg.width,
	    		height: msg.height,
	    		generation: msg.generation,
	    		dataSizeReceived: size,
	    		paused: msg.paused,
	    		interval: msg.interval
	    	}, this.didReceiveNewState.bind(this, msg))
	    } else if (msg.paused != null) {
	    	this.setState({
	    		paused: msg.paused,
	    		pausedBy: msg.player
	    	}, this.renderToolbar.bind(this))
	    } else if (msg.interval != null) {
	    	this.setState({
	    		interval: msg.interval,
	    		intervalSetBy: msg.player
	    	}, this.renderToolbar.bind(this))
	    }
	}

	didReceiveNewState(msg) {
		this.renderToolbar()
		this.renderGrid(msg.activeCells)
	}

	renderToolbar() {
		this.toolbarRef.component.setState({
			interval: this.state.interval,
			paused: this.state.paused
		})
	}

	renderGrid(nextActiveCellsArr) {
		var currentActiveCells = {}
		for (var i = 0; i < this.activeCells.length; i++) {
			var cell = this.activeCells[i]
			currentActiveCells[`${cell.x}:${cell.y}`] = cell
		}
		var nextActiveCells = {}
		for (var i = 0; i < nextActiveCellsArr.length; i++) {
			var cell = nextActiveCellsArr[i]
			nextActiveCells[`${cell.x}:${cell.y}`] = cell
		}
		Object.keys(currentActiveCells).map((cellId) => {
			if (nextActiveCells[cellId] == null) {
				var cell = this.cellRefs[cellId]
				if (cell != null) {
					cell.setState({
						active: false,
						color: null
					})
				}
			}
		})
		Object.keys(nextActiveCells).map((cellId) => {
			var cell = this.cellRefs[cellId]
			if (cell != null) {
				cell.setState({
					active: true,
					red: nextActiveCells[cellId].c.R,
					green: nextActiveCells[cellId].c.G,
					blue: nextActiveCells[cellId].c.B
				})
			}
		})
		this.activeCells = nextActiveCellsArr
	}

	onClickCell(x, y) {
		var cell = this.cellRefs[`${x}:${y}`]
		if (cell.state.active && this.toolbarRef.component.state.pattern == -1) {
			console.log("deactivate")
			axios.post("/deactivate", {
				point: {x, y},
				player: this.state.player
			}).then((response) => {
				console.log(response)
			}).catch((err) => {
				console.log(err.response)
			})
		} else {
			var points = this.toolbarRef.component.patternPoints({x, y})
			axios.post("/activate", {
				points: points,
				player: this.state.player
			}).then((response) => {
				console.log(response)
			}).catch((err) => {
				console.log(err.response)
			})
		}
	}

	onClickPauseHandler(e) {
		var currentPauseState = this.toolbarRef.component.state.paused
		var nextPauseState = !currentPauseState
		this.toolbarRef.component.setState({
			paused: nextPauseState
		})
		axios.post("/pause", {
			player: this.state.player,
			pause: nextPauseState
		}).then((response) => {
			console.log(response)
		}).catch((err) => {
			console.log(err.response)
		})
	}

	onChangeIntervalHandler(e) {
		var interval = e.target.value
		this.toolbarRef.component.setState({
			interval: e.target.value
		})
		axios.post("/interval", {
			player: this.state.player,
			interval: parseInt(interval)
		}).then((response) => {
			console.log(response)
		}).catch((err) => {
			console.log(err.response)
		})
	}

	onClickResetHandler() {
		axios.post("/reset", {
			player: this.state.player
		}).then((response) => {
			console.log(response)
		}).catch((err) => {
			console.log(err.response)
		})
	}

	onClickPattern(patternId) {
		this.toolbarRef.component.setState({
			pattern: this.toolbarRef.component.state.pattern == patternId ? -1 : patternId
		})
	}
}

export {
	App
}