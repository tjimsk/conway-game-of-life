import React from "react"

const styles = require("Components/ConnectBackground.scss")

class ConnectBackground extends React.PureComponent {
	render() {
		return (
		<div className={styles.connectBackground}>
			<div className={styles.container}>
				<span className={styles.message}>Connecting...</span>
			</div>
		</div>)
	}
}

export {
	ConnectBackground
}