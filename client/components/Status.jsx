import React from "react"

const styles = require("Components/Status.scss")

const READYSTATES = {
	0: "CONNECTING",
	1: "CONNECTED",
	2: "CLOSING",
	3: "CLOSED"
}

class Status extends React.Component {
	render() {
		let playerName = this.props.player.name ? `as ${this.props.player.name}` : ""

		return (
		<div className={styles.status}>
			<span className={styles.readyState}>
				{READYSTATES[this.props.readyState] || "DISCONNECTED"} {playerName.toUpperCase()}
			</span>
			<span className={styles.generation}>
				GENERATION #{this.props.generation}
			</span>
			<span className={styles.dataSizeReceived}>
				{this.props.dataSizeReceived} KB
			</span>
		</div>)
	}
}

export {
	Status
}