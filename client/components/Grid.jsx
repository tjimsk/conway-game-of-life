import React from "react"
import {Cell} from "Components/Cell"

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
							cellRefs: this.props.cellRefs,
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

export {
	Grid
}

