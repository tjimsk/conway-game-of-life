import React from "react"

const styles = require("Components/Cell.scss")

class Cell extends React.Component {
	render() {
		var props = {
			onClick: this.props.onClickCell.bind(this, this.props.cell)
		}

		var style = {
			left: this.state.left,
			top: this.state.top,
			backgroundColor: this.color()
		}

		return <div className={styles.cell} {...props} style={style}></div>
	}

	constructor(props) {
		super(props)

		this.state = {
			left: 16 * (props.cell.x-1),
			top: 16 * (props.cell.y-1),
			active: props.cell.a, // initial active
			color: props.cell.c // initial color
		}

		// bind setState function to allow direct call from App so
		// Grid does not need to update
		props.setActiveColorFunc[props.cell.k] = this.setActiveColor.bind(this)
		props.cellRefs[props.cell.k] = this
	}

	color() {
		if (this.state.active) {
			var r = this.state.color.r
			var g = this.state.color.g
			var b = this.state.color.b

			return `rgb(${r},${g},${b})`
		} else {
			return "transparent"
		}

		if (this.state.active) {
			var uc = this.props.user.c
			style.backgroundColor = `rgb(${uc.r},${uc.g},${uc.b})`
		} else {
			style.backgroundColor = "transparent"
		}

		return c
	}

	setActiveColor(color) {
		this.setState({
			active: color != null,
			color: color
		})
	}

	onClickCell() {
		this.props.onClickCell(this.props.cell)

		this.setActiveColor(this.state.active ? null : this.props.user.c)
	}
}

export {
	Cell
}