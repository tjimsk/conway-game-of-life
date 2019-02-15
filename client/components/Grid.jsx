import React from "react"
import styles from "Styles/Grid.scss"

export default class Grid extends React.Component {
	render() {
		var [rows, cols] = [[], []]
		for (var y = 1; y <= this.props.height; y++) rows.push(y)			
		for (var x = 1; x <= this.props.width; x++) cols.push(x)

		return	<div className={styles.grid}>
					<div className={styles.gridContainer}>
						{rows.map((rowId) => {
							return 	<Row key={rowId} 
										rowId={rowId}
										cols={cols} 
										cellRefs={this.props.cellRefs}
										clickHandler={this.props.onClickCell} />
						})}
					</div>
				</div>
	}

	shouldComponentUpdate(nextProps, nextState) {
		return this.props.width * this.props.height == 0
	}
}

const Row = (props) => {
	let {cols, rowId, clickHandler, cellRefs} = props

	return 	<div className={styles.row} key={rowId}>
				{cols.map((colId) => {
					return 	<Cell key={`${colId}:${rowId}`}
								cellId={`${colId}:${rowId}`}
								cellRefs={cellRefs}
								x={colId}
								y={rowId}
								clickHandler={clickHandler.bind(this)} />
				})}
			</div>
}

class Cell extends React.Component {
	render() {
		var style = this.state.active ? {
			backgroundColor: `rgb(${this.state.red},${this.state.green},${this.state.blue})`
		} : {}
		return 	<div className={styles.cell}>
					<div className={styles.cellContainer} 
						style={style} 
						onClick={this.props.clickHandler.bind(this, this.props.x, this.props.y)}></div>
				</div>
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
