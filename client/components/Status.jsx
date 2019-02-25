import React from "react"
import styles from "Styles/Status.scss"

const READYSTATES = {
	0: "Connecting",
	1: "Connected",
	2: "Closing",
	3: "Closed"
}

export default props => {
	let {
		// player,
		readyState,
		generation,
		interval,
		dataSizeReceived
	} = props

	let kbps = parseFloat(dataSizeReceived) / parseFloat(interval) * 1000
	kbps = kbps ? `${kbps.toFixed(1)} KB/s` : ""

	return	<div className={styles.status}>
				<span className={styles.label}>{READYSTATES[readyState] || "Disconnected"}</span>
				{/*<span className={styles.value}>{player.name}</span>*/}
				<span className={styles.label}>Generation</span>
				<span className={styles.value}>#{generation}</span>
				<span className={styles.value}>{(interval || 0) / 1000}s</span>
				<span className={styles.value}>{kbps}</span>
			</div>
}