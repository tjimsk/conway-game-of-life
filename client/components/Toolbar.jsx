import React from "react"

const styles = require("Components/Toolbar.scss")

class Toolbar extends React.Component {
	render() {
		return (
		<div className={styles.toolbar}>
			<div className={styles.play}>
				<div className={styles.label}>Play:</div>
				<div className={styles.interval}>
					<input type="range" 
						min={this.props.minInterval}
						max={this.props.maxInterval}
						step="100"
						value={this.state.interval} 
						onChange={this.props.onChangeIntervalHandler.bind(this)} />
				</div>				
				<div className={this.state.paused ? styles.resume : styles.pause}
				 	onClick={this.props.onClickPauseHandler.bind(this)}>
				</div>
			</div>

			<div className={styles.pattern}>
				<div className={styles.label}>Patterns:</div>
				<div className={`${styles.block}${this.state.pattern == 0 ? ` ${styles.selected}` : ""}`} title="block" 
					onClick={this.props.onClickPattern.bind(this, 0)}></div>
				<div className={`${styles.beehive}${this.state.pattern == 1 ? ` ${styles.selected}` : ""}`} title="beehive"
					onClick={this.props.onClickPattern.bind(this, 1)}></div>
				<div className={`${styles.loaf}${this.state.pattern == 2 ? ` ${styles.selected}` : ""}`} title="loaf"
					onClick={this.props.onClickPattern.bind(this, 2)}></div>
				<div className={`${styles.boat}${this.state.pattern == 3 ? ` ${styles.selected}` : ""}`} title="boat"
					onClick={this.props.onClickPattern.bind(this, 3)}></div>
				<div className={`${styles.tub}${this.state.pattern == 4 ? ` ${styles.selected}` : ""}`} title="tub"
					onClick={this.props.onClickPattern.bind(this, 4)}></div>
				<div className={`${styles.blinker}${this.state.pattern == 5 ? ` ${styles.selected}` : ""}`} title="blinker"
					onClick={this.props.onClickPattern.bind(this, 5)}></div>
				<div className={`${styles.toad}${this.state.pattern == 6 ? ` ${styles.selected}` : ""}`} title="toad"
					onClick={this.props.onClickPattern.bind(this, 6)}></div>
				<div className={`${styles.beacon}${this.state.pattern == 7 ? ` ${styles.selected}` : ""}`} title="beacon"
					onClick={this.props.onClickPattern.bind(this, 7)}></div>
				<div className={`${styles.pulsar}${this.state.pattern == 8 ? ` ${styles.selected}` : ""}`} title="pulsar"
					onClick={this.props.onClickPattern.bind(this, 8)}></div>
				<div className={`${styles.glider}${this.state.pattern == 9 ? ` ${styles.selected}` : ""}`} title="glider"
					onClick={this.props.onClickPattern.bind(this, 9)}></div>
				<div className={`${styles.lwss}${this.state.pattern == 10 ? ` ${styles.selected}` : ""}`} title="lwss"
					onClick={this.props.onClickPattern.bind(this, 10)}></div>
				<div className={`${styles.mwss}${this.state.pattern == 11 ? ` ${styles.selected}` : ""}`} title="mwss"
					onClick={this.props.onClickPattern.bind(this, 11)}></div>
			</div>
		</div>)
	}

	constructor(props) {
		super(props)
		this.state = {
			paused: props.paused,
			interval: 1100,
			pattern: -1
		}
		props.toolbarRef.component = this
	}

	patternPoints(p) {
		switch (this.state.pattern) {
		case 0:
			return this.blockPoints(p)
		case 1:
			return this.beehivePoints(p)
		case 2:
			return this.loafPoints(p)
		case 3:
			return this.boatPoints(p)
		case 4:
			return this.tubPoints(p)
		case 5:
			return this.blinkerPoints(p)
		case 6:
			return this.toadPoints(p)
		case 7:
			return this.beaconPoints(p)
		case 8:
			return this.pulsarPoints(p)
		case 9:
			return this.gliderPoints(p)
		case 10:
			return this.lwssPoints(p)
		case 11:
			return this.mwssPoints(p)
		default:
			return [p]
		}
	}

	blockPoints(p) {
		var x = p.x
		var y = p.y
		return [
			{x: x, 		y: y},
			{x: x, 		y: y + 1},
			{x: x + 1, 	y: y},
			{x: x + 1, 	y: y + 1}
		]
	}

	beehivePoints(p) {
		var x = p.x
		var y = p.y
		return [
			{x: x, 		y: y},
			{x: x + 1, 	y: y - 1},
			{x: x + 1, 	y: y + 1},
			{x: x + 2, 	y: y - 1},
			{x: x + 2, 	y: y + 1},
			{x: x + 3, 	y: y}
		]
	}

	loafPoints(p) {
		var x = p.x
		var y = p.y
		return [
			{x: x, 		y: y},
			{x: x + 1, 	y: y - 1},
			{x: x + 1, 	y: y + 1},
			{x: x + 2, 	y: y - 1},
			{x: x + 2, 	y: y + 2},
			{x: x + 3, 	y: y},
			{x: x + 3, 	y: y + 1},
			{x: x + 3, 	y: y + 2}
		]
	}

	boatPoints(p) {
		var x = p.x
		var y = p.y
		return [
			{x: x, 		y: y},
			{x: x, 		y: y + 1},
			{x: x + 1, 	y: y},
			{x: x + 1, 	y: y + 2},
			{x: x + 2, 	y: y + 1}
		]
	}

	tubPoints(p) {
		var x = p.x
		var y = p.y
		return [
			{x: x, 		y: y},
			{x: x + 1, 	y: y - 1},
			{x: x + 1, 	y: y + 1},
			{x: x + 2, 	y: y}
		]
	}

	blinkerPoints(p) {
		var x = p.x
		var y = p.y
		return [
			{x: x, 		y: y},
			{x: x + 1, 	y: y},
			{x: x + 2, 	y: y}
		]
	}

	toadPoints(p) {
		var x = p.x
		var y = p.y
		return [
			{x: x, 		y: y},
			{x: x, 		y: y + 1},
			{x: x + 1, 	y: y},
			{x: x + 1, 	y: y + 1},
			{x: x + 2, 	y: y},
			{x: x - 1 , y: y + 1}
		]
	}

	beaconPoints(p) {
		var x = p.x
		var y = p.y
		return [
			{x: x, 		y: y},
			{x: x, 		y: y + 1},
			{x: x + 1, 	y: y},
			{x: x + 1, 	y: y + 1},
			{x: x + 2, 	y: y + 2},
			{x: x + 2, 	y: y + 3},
			{x: x + 3, 	y: y + 2},
			{x: x + 3 , y: y + 3}
		]
	}

	pulsarPoints(p) {
		var x = p.x
		var y = p.y
		return [
			{x: x, 		y: y},
			{x: x, 		y: y + 1},
			{x: x, 		y: y + 2},
			{x: x, 		y: y + 6},
			{x: x, 		y: y + 7},
			{x: x, 		y: y + 8},

			{x: x + 2, 	y: y - 2},
			{x: x + 2, 	y: y + 3},
			{x: x + 2, 	y: y + 5},
			{x: x + 2, 	y: y + 10},

			{x: x + 3, 	y: y - 2},
			{x: x + 3, 	y: y + 3},
			{x: x + 3, 	y: y + 5},
			{x: x + 3, 	y: y + 10},

			{x: x + 4, 	y: y - 2},
			{x: x + 4, 	y: y + 3},
			{x: x + 4, 	y: y + 5},
			{x: x + 4, 	y: y + 10},

			{x: x + 5, 	y: y},
			{x: x + 5, 	y: y + 1},
			{x: x + 5, 	y: y + 2},
			{x: x + 5, 	y: y + 6},
			{x: x + 5, 	y: y + 7},
			{x: x + 5, 	y: y + 8},

			{x: x + 7, 	y: y},
			{x: x + 7, 	y: y + 1},
			{x: x + 7, 	y: y + 2},
			{x: x + 7, 	y: y + 6},
			{x: x + 7, 	y: y + 7},
			{x: x + 7, 	y: y + 8},

			{x: x + 8, 	y: y - 2},
			{x: x + 8, 	y: y + 3},
			{x: x + 8, 	y: y + 5},
			{x: x + 8, 	y: y + 10},

			{x: x + 9, 	y: y - 2},
			{x: x + 9, 	y: y + 3},
			{x: x + 9, 	y: y + 5},
			{x: x + 9, 	y: y + 10},

			{x: x + 10, y: y - 2},
			{x: x + 10, y: y + 3},
			{x: x + 10, y: y + 5},
			{x: x + 10, y: y + 10},

			{x: x + 12, y: y},
			{x: x + 12, y: y + 1},
			{x: x + 12, y: y + 2},
			{x: x + 12, y: y + 6},
			{x: x + 12, y: y + 7},
			{x: x + 12, y: y + 8}
		]
	}

	gliderPoints(p) {
		var x = p.x
		var y = p.y
		return [
			{x: x, 		y: y},
			{x: x + 1, 	y: y},
			{x: x + 1, 	y: y - 2},
			{x: x + 2, 	y: y},
			{x: x + 2, 	y: y - 1},
		]
	}

	lwssPoints(p) {
		var x = p.x
		var y = p.y
		return [
			{x: x, y: y},
			{x: x, y: y + 1},
			{x: x + 1, y: y - 1},
			{x: x + 1, y: y},
			{x: x + 1, y: y + 1},
			{x: x + 2, y: y - 1},
			{x: x + 2, y: y},
			{x: x + 2, y: y + 2},
			{x: x + 3, y: y},
			{x: x + 3, y: y + 1},
			{x: x + 3, y: y + 2},
			{x: x + 4, y: y + 1}
		]
	}

	mwssPoints(p) {
		var x = p.x
		var y = p.y
		return [
			{x: x, y: y},
			{x: x, y: y + 1},
			{x: x + 1, y: y - 1},
			{x: x + 1, y: y},
			{x: x + 1, y: y + 1},
			{x: x + 2, y: y - 1},
			{x: x + 2, y: y},
			{x: x + 2, y: y + 1},
			{x: x + 3, y: y - 1},
			{x: x + 3, y: y},
			{x: x + 3, y: y + 2},
			{x: x + 4, y: y},
			{x: x + 4, y: y + 1},
			{x: x + 4, y: y + 2},
			{x: x + 5, y: y + 1}
		]
	}


}

export {
	Toolbar
}