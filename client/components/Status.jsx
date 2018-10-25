import React from "react"

const styles = require("Components/Status.scss")

const STATUS_CONNECTING = 0
const STATUS_CONNECTED = 1
const STATUS_CLOSING = 2
const STATUS_CLOSED = 3

class Status extends React.Component {
	render() {
		let status = "NOT CONNECTED"
		switch(this.props.readyState) {
		case STATUS_CONNECTING:
			status = "CONNECTING"; break
		case STATUS_CONNECTED:
			status = "CONNECTED"; break
		case STATUS_CLOSING:
			status = "CLOSING"; break
		case STATUS_CLOSED:
			status = "CLOSED"; break
		}

		return (
		<div className={styles.status}>
			<span className={styles.readyState}>{status}</span>
			<span className={styles.generation}>GENERATION #{this.props.grid.generation}</span>
			<span className={styles.dataSizeReceived}>RECEIVED {this.props.dataSizeReceived} KB</span>
		</div>)
	}
}

export {
	Status
}