import React from "react"

const styles = require("Components/Grid.scss")

class Grid extends React.Component {
	render() {
		return (
			<div className={styles.grid}>
				<div className={styles.scrollView}>

					{Object.keys(this.props.cells).map((k) => {
						var props = {
							key: k,
							user: this.props.user,
							cell: this.props.cells[k],
							onClickCell: this.props.onClickCell,
							setActiveColorFunc: this.props.setActiveColorFunc
						}

						return (
							<Cell {...props} />)
					})}
				</div>
			</div>)
	}

	// disable updates from parent propagation
	shouldComponentUpdate(nextProps, nextState) {
		if (nextProps.receivedActiveCells) {
			return false
		} else {
			return true
		}
	}
}

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

		if (c.g > this.state.clickGeneration) {
			if (c.a && c.c != null) {
				style.backgroundColor = `rgb(${c.c.r},${c.c.g},${c.c.b})`
			} else {
				style.backgroundColor = "transparent"
			}
		} else {
			if (this.state.active) {
				var uc = this.props.user.color
				style.backgroundColor = `rgb(${uc.r},${uc.g},${uc.b})`
			} else {
				style.backgroundColor = "transparent"
			}
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

		this.setActiveColor(this.state.active ? null : this.props.user.color)
	}
}

export {
	Grid
}

