import React from "react"

const styles = require("Components/Toolbox.scss")

class Toolbox extends React.Component {
	render() {
		return (
		<div className={styles.toolbox}>
			{Object.keys(this.tools).map((k) => {
				var tool = this.tools[k]
				return (
				<Tool key={k}
					{...tool}
					selected={k == this.props.selectedToolIndex}
					onClickTool={this.props.onClickTool} />)
			})}
		</div>)
	}

	constructor(props) {
		super(props)

		this.tools = {
			0: {
				index: 0,
				name: "Point",
				func: (x, y) => {
					return [
						{x: x, y: y}
					]
				}
			},
			1: {
				index: 1,
				name: "Block",
				func: (x, y) => {
					return [
						{x: x, y: y},
						{x: x + 1, y: y},
						{x: x, y: y + 1},
						{x: x + 1, y: y + 1}
					]
				}
			},
			2: {
				index: 2,
				name: "Beehive",
				func: (x, y) => {
					return [
						{x: x, y: y},
						{x: x + 1, y: y - 1},
						{x: x + 1, y: y + 1},
						{x: x + 2, y: y - 1},
						{x: x + 2, y: y + 1},
						{x: x + 3, y: y}
					]
				}
			},
			3: {
				index: 3,
				name: "Toad",
				func: (x, y) => {
					return [
						{x: x, y: y},
						{x: x, y: y + 1},
						{x: x + 1, y: y},
						{x: x + 1, y: y + 1},
						{x: x + 2, y: y},
						{x: x - 1 , y: y + 1}
					]
				}
			},
			4: {
				index: 4,
				name: "Pulsar",
				func: (x, y) => {
					return [
						{x: x, y: y},
						{x: x, y: y + 1},
						{x: x, y: y + 2},
						{x: x, y: y + 6},
						{x: x, y: y + 7},
						{x: x, y: y + 8},

						{x: x + 5, y: y},
						{x: x + 5, y: y + 1},
						{x: x + 5, y: y + 2},
						{x: x + 5, y: y + 6},
						{x: x + 5, y: y + 7},
						{x: x + 5, y: y + 8},

						{x: x + 7, y: y},
						{x: x + 7, y: y + 1},
						{x: x + 7, y: y + 2},
						{x: x + 7, y: y + 6},
						{x: x + 7, y: y + 7},
						{x: x + 7, y: y + 8},

						{x: x + 12, y: y},
						{x: x + 12, y: y + 1},
						{x: x + 12, y: y + 2},
						{x: x + 12, y: y + 6},
						{x: x + 12, y: y + 7},
						{x: x + 12, y: y + 8},

						{x: x + 2, y: y - 2},
						{x: x + 3, y: y - 2},
						{x: x + 4, y: y - 2},
						{x: x + 8, y: y - 2},
						{x: x + 9, y: y - 2},
						{x: x + 10, y: y - 2},

						{x: x + 2, y: y + 3},
						{x: x + 3, y: y + 3},
						{x: x + 4, y: y + 3},
						{x: x + 8, y: y + 3},
						{x: x + 9, y: y + 3},
						{x: x + 10, y: y + 3},

						{x: x + 2, y: y + 5},
						{x: x + 3, y: y + 5},
						{x: x + 4, y: y + 5},
						{x: x + 8, y: y + 5},
						{x: x + 9, y: y + 5},
						{x: x + 10, y: y + 5},

						{x: x + 2, y: y + 10},
						{x: x + 3, y: y + 10},
						{x: x + 4, y: y + 10},
						{x: x + 8, y: y + 10},
						{x: x + 9, y: y + 10},
						{x: x + 10, y: y + 10}
					]
				}
			},
			5: {
				index: 5,
				name: "Lightweight Spaceship",
				func: (x, y) => {
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
			}
		}
	}

	selectedToolName() {
		var tool = this.tools[this.props.selectedToolIndex]

		return tool.name
	}

	getCellsFromSelectedTool(x, y) {
		var tool = this.tools[this.props.selectedToolIndex]

		return tool.func(x, y)
	}
}

class Tool extends React.Component {
	render() {
		return (
		<div className={this.props.selected ? styles.selectedTool : styles.tool}
		 onClick={this.props.onClickTool.bind(this, this.props.index)} >
			<span className={styles.label}>{this.props.name}</span>
		</div>)
	}
}

export {
	Toolbox
}