import React from "react"
import {Toolbar, ToolbarDrawer, tools} from "Components/Toolbar"
import Grid from "Components/Grid"
import Status from "Components/Status"
import wsClient from "Util/websocket"
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
					<Toolbar
						interval={this.state.interval}
						toolId={this.state.selectedToolIndex}
						onChangeIntervalHandler={this.onChangeIntervalHandler.bind(this)}
						onClickResetHandler={this.onClickResetHandler.bind(this)}
						onClickToolsMenuHandler={this.onClickToolsMenuHandler.bind(this)}
						onClickPattern={this.onClickPattern.bind(this)} />

					<Grid 
						width={this.state.width}
						height={this.state.height}
						cellRefs={this.cellRefs}
						onClickCell={this.onClickCell.bind(this)}
						transitionTime={this.state.interval / 10} />

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
		this.activeCells = {}
		this.createWebsocketConnection()
	}

	getInitialState() {
		return {
			width: 90,
			height: 75,
			generation: 0,
			interval: 1000,
			player: {color: {r: 0, g: 0, b: 0}},
			dataSizeReceived: 0,
			readyState: -1,
			showToolbarDrawer: false,
			selectedToolIndex: 0,
		}
	}

	updateActiveCells(activeCells) {
		var activeCellsObj = {}
		activeCells.map((cell) => {
			let cellId = `${cell.p.X}:${cell.p.Y}`
			this.activeCells[cellId] = cell
			activeCellsObj[cellId] = cell
		})

		Object.keys(this.activeCells).map((cellId) => {
			let cell = activeCellsObj[cellId]
			let cellRef = this.cellRefs[cellId]
			if (cellRef == null) return
			if (cell != null) {
				cellRef.setState({
					active: true,
					red: cell.c.R,
					green: cell.c.G,
					blue: cell.c.B
				})
			} else {
				cellRef.setState({active: false})
				delete this.activeCells[cellId]
			}
		})
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

		this.setState({showToolbarDrawer: false})
	}

	onChangeIntervalHandler(interval) {
		fetch("/interval", {
			method: "POST",
			body: JSON.stringify({
				player: this.state.player,
				interval: interval
			})
		}).catch((err) => console.log(err))
		this.setState({showToolbarDrawer: false})
	}

	onClickResetHandler() {
		fetch("/reset", {
			method: "POST",
			body: JSON.stringify({
				player: this.state.player
			})
		}).catch((err) => console.log(err))
		this.setState({showToolbarDrawer: false})
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

	// websocket
	createWebsocketConnection() {
		this.client = new wsClient()
		this.client.onopen = this.onWebSocketOpen.bind(this)
		this.client.onclose = this.onWebSocketClose.bind(this)
		this.client.onmessage = this.onWebSocketMessage.bind(this)
	}

	onWebSocketOpen() {
		setTimeout(() => {
			this.setState({readyState: this.client.readyState}, 
				() => this.onWebSocketOpen())
		}, 1000)
	}

	onWebSocketClose() {
		this.setState(this.getInitialState())
		setTimeout(() => this.createWebsocketConnection(), 500)
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
	    		generation: msg.generation,
	    		dataSizeReceived: size,
	    		interval: msg.interval
	    	}, this.updateActiveCells.bind(this, msg.activeCells))
	    } else if (msg.interval != null) {
	    	this.setState({
	    		interval: msg.interval,
	    		intervalSetBy: msg.player
	    	})
	    }
	}
}

export {
	App
}