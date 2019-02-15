import React from "react"
import {Toolbar, ToolbarDrawer, tools} from "Components/Toolbar"
import Grid from "Components/Grid"
import Status from "Components/Status"
import Spinner from "Components/Spinner"
import websocket from "Util/websocket"
import styles from "Styles/App.scss"
import {TransitionGroup, CSSTransition} from "react-transition-group"

class App extends React.Component {
	render() {
		let toolbarDrawerTransition = {
			classNames: {
				enter: 			styles.toolbarDrawerTransitionEnter,
				enterActive: 	styles.toolbarDrawerTransitionEnterActive,
				enterDone: 		styles.toolbarDrawerTransitionEnterDone,
				exit: 			styles.toolbarDrawerTransitionExit,
				exitActive: 	styles.toolbarDrawerTransitionExitActive,
				exitDone: 		styles.toolbarDrawerTransitionExitDone
			},
			timeout: 150,
			in: true
		}

		return (
			<div className={styles.app}>
				<div className={styles.appContainer}>
					{this.state.readyState !== 1 ? <Spinner /> : null}

					<Toolbar
						minInterval={200}
						maxInterval={4000}
						increment={400}
						interval={this.state.interval}
						paused={this.state.paused}
						toolId={this.state.selectedToolIndex}
						toolbarRef={this.toolbarRef}
						onClickPauseHandler={this.onClickPauseHandler.bind(this)}
						onChangeIntervalHandler={this.onChangeIntervalHandler.bind(this)}
						onClickResetHandler={this.onClickResetHandler.bind(this)}
						onClickToolsMenuHandler={this.onClickToolsMenuHandler.bind(this)}
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
						interval={this.state.interval}
						dataSizeReceived={this.state.dataSizeReceived} />
				</div>

				<TransitionGroup className={styles.transitionDiv}>
					{!this.state.showToolbarDrawer ? null :
					<CSSTransition {...toolbarDrawerTransition}>
						<ToolbarDrawer
							toolId={this.state.selectedToolIndex}
							onClickToolsMenuHandler={this.onClickToolsMenuHandler.bind(this)}
							onClickToolsMenuPattern={this.onClickToolsMenuPattern.bind(this)} />
					</CSSTransition>}
				</TransitionGroup>
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
			showToolbarDrawer: false,
			paused: true,
			selectedToolIndex: 0,
		}
	}

	// websocket
	createWebsocketConnection() {
		this.client = websocket()
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
	    	}, this.renderGrid.bind(this, msg.activeCells))
	    } else if (msg.paused != null) {
	    	this.setState({
	    		paused: msg.paused,
	    		pausedBy: msg.player
	    	})
	    } else if (msg.interval != null) {
	    	this.setState({
	    		interval: msg.interval,
	    		intervalSetBy: msg.player
	    	})
	    }
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
		if (cell.state.active && this.state.selectedToolIndex == 0) {
			console.log("deactivate")
			fetch("/deactivate", {
				method: "POST",
				body: JSON.stringify({
					point: {x, y},
					player: this.state.player
				})
			}).catch((err) => console.log(err))
		} else {
			var points = tools[this.state.selectedToolIndex].func({x,y})
			fetch("/activate", {
				method: "POST",
				body: JSON.stringify({
					points: points,
					player: this.state.player
				})
			}).catch((err) => console.log(err))
		}
	}

	onClickPauseHandler(e) {
		let pauseState = !this.state.paused
		this.setState({
			paused: pauseState
		})
		fetch("/pause", {
			method: "POST",
			body: JSON.stringify({
				player: this.state.player,
				pause: pauseState
			})
		}).catch((err) => console.log(err))
	}

	onChangeIntervalHandler(interval) {
		interval = interval <= 200 ? 200 : interval
		interval = interval >= 4000 ? 4000 : interval

		fetch("/interval", {
			method: "POST",
			body: JSON.stringify({
				player: this.state.player,
				interval: interval
			})
		}).catch((err) => console.log(err))
	}

	onClickResetHandler() {
		fetch("/reset", {
			method: "POST",
			body: JSON.stringify({
				player: this.state.player
			})
		}).catch((err) => console.log(err))
	}

	onClickToolsMenuHandler() {
		this.setState({
			showToolbarDrawer: !this.state.showToolbarDrawer
		}, () => console.log(this.state.showToolbarDrawer))
	}

	onClickPattern(patternId) {
		this.setState({
			selectedToolIndex: patternId
		})
	}

	onClickToolsMenuPattern(patternId) {
		this.setState({
			selectedToolIndex: patternId,
			showToolbarDrawer: false
		})
	}
}

export {
	App
}