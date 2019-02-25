import React from "react"
import styles from "Styles/Grid.scss"

export default class Grid extends React.Component {
	render() {
		var [rows, cols] = [[], []]
		for (var y = 1; y <= this.props.height; y++) rows.push(y)			
		for (var x = 1; x <= this.props.width; x++) cols.push(x)

		let {cellRefs, transitionTime, onClickCell} = this.props

		return	<div className={styles.grid}>
					<div className={styles.gridContainer}>
						{rows.map((rowId) => {
							return 	<Row key={rowId} rowId={rowId} cols={cols} cellRefs={cellRefs} transitionTime={transitionTime} clickHandler={onClickCell} />
						})}
					</div>
				</div>
	}

	shouldComponentUpdate(nextProps, nextState) {
		return this.props.width * this.props.height == 0
	}
}

const Row = (props) => {
	let {cols, rowId, clickHandler, cellRefs, transitionTime} = props

	return 	<div className={styles.row} key={rowId}>
				{cols.map((colId) => {
					return 	<Cell key={`${colId}:${rowId}`}
								cellId={`${colId}:${rowId}`}
								cellRefs={cellRefs}
								x={colId}
								y={rowId}
								transitionTime={transitionTime}
								clickHandler={clickHandler.bind(this)} />
				})}
			</div>
}

class Cell extends React.Component {
	render() {
		let {x, y, clickHandler, transitionTime} = this.props
		let style = !this.state.active ? {} : {
			backgroundColor: `rgb(${this.state.red},${this.state.green},${this.state.blue})`,
			// transition: `background-color ${transitionTime}ms ease-in-out`
		}

		return 	<div className={styles.cell} onClick={clickHandler.bind(this, x, y)}>
					<div className={styles.cellContainer} style={style}></div>
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
