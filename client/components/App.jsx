import React from "react"
import {Grid} from "Components/Grid"
import {Controls} from "Components/Controls"

import {client} from "Util/websocket"

const styles = require("Components/App.scss")

class App extends React.Component {
	render() {
		return (
			<div className={styles.app}>
				<Grid />
				<Controls />
			</div>)
	}

	constructor(props) {
		super(props)

		client.onopen = this.onWebSocketOpen
		client.onclose = this.onWebSocketClose
		client.onmessage = this.onWebSocketMessage
		client.onerror = this.onWebSocketError
	}

	onWebSocketOpen() {
		console.log("websocket client connected")
	}

	onWebSocketClose() {
		console.log("websocket client closed")
	}

	onWebSocketMessage(e) {
	    if (typeof e.data === "string") {
	        var msg = JSON.parse(e.data)
	        console.log("WebSocket message:", msg)
	    }
	}

	onWebSocketError() {
		console.log("websocket connection error")
	}
}

export {
	App
}