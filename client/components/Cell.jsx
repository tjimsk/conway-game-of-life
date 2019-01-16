import React from "react"

const styles = require("Components/Cell.scss")

class Cell extends React.Component {
	render() {
		var style = this.state.active ? {
			backgroundColor: `rgb(${this.state.red},${this.state.green},${this.state.blue})`
		} : {}
		return (
		<td className={styles.cell}>
			<div className={styles.content} 
				style={style} 
				onClick={this.props.onClickCell.bind(this, this.props.x, this.props.y)}></div>
		</td>)
	}

	constructor(props) {
		super(props)
		this.state = {
			active: false,
			red: 255, 
			green: 255,
			blue: 255
		}		
		props.cellRefs[props.cellId] = this
	}
}

export {
	Cell
}