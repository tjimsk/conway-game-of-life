import React from "react"
import {Cell} from "Components/Cell"

const styles = require("Components/Grid.scss")

class Grid extends React.Component {
	render() {
		var rows = []
		for (var y = 1; y <= this.props.height; y++) {
			rows.push(y)
		}
		var cols = []
		for (var x = 1; x <= this.props.width; x++) {
			cols.push(x)
		}

		return (
			<div className={styles.grid}>
				<table className={styles.table}>
					<tbody>
						{rows.map((rowId) => {
							return (
							<tr key={rowId}>
								{cols.map((colId) => {
									return (
									<Cell key={`${colId}:${rowId}`} 
										cellId={`${colId}:${rowId}`}
										cellRefs={this.props.cellRefs}
										x={colId}
										y={rowId}
										onClickCell={this.props.onClickCell.bind(this)} />)
								})}
							</tr>)
						})}
					</tbody>
				</table>
			</div>)
	}

	// disable updates from parent propagation
	shouldComponentUpdate(nextProps, nextState) {
		if (this.props.width * this.props.height == 0) {
			return true
		} else {
			return false
		}
	}
}

export {
	Grid
}

