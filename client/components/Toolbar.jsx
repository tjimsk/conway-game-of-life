import React from "react"
import styles from "Styles/Toolbar.scss"

const Toolbar = props => {
	let {
		interval,
		toolId, 
		onChangeIntervalHandler, 
		onClickResetHandler,
		onClickToolsMenuHandler,
		onClickPauseHandler,
		onClickPattern
	} = props
	
	let speeds = ["0.5x", "1x", "2x", "4x", "8x", "16x"]
	let intervals = [2000,1000,500,250,120,60]
	let speed = speeds[intervals.indexOf(interval)]
	let nextInterval = interval == intervals[intervals.length - 1] ? intervals[0] : intervals[intervals.indexOf(interval) + 1]

	return 	<div className={styles.toolbar}>
				<div className={styles.toolbarContainer}>
					<div className={styles.interval}
						onClick={onChangeIntervalHandler.bind(this, nextInterval)}>{speed}</div>

					<div className={interval == -1 ? styles.resume : styles.pause}
					 	onClick={onChangeIntervalHandler.bind(this, interval == -1 ? 1000 : -1)}>
					 	{interval == -1 ? "Resume" : "Pause"}
					</div>

					<div className={styles.reset}
					 	onClick={onClickResetHandler.bind(this)}>Reset</div>

					<div className={styles.toolsMenu}
						onClick={onClickToolsMenuHandler.bind(this)}>Tools</div>

					<div className={styles.tools}>
						{tools.map((tool, id) => {
							let className = `${styles.tool} ${toolId == id ? `${styles.selected}` : ""}`
							return 	<div className={className}
										onClick={onClickPattern.bind(this, id)}
										key={id} >{tool.name}</div>
						})}
					</div>
				</div>
			</div>
}

const ToolbarDrawer = props => {
	let {
		toolId,
		onClickToolsMenuHandler,
		onClickToolsMenuPattern
	} = props

	return 	<div className={styles.toolbarDrawer}>
				<div className={styles.toolbarDrawerClose}
					onClick={onClickToolsMenuHandler.bind(this)}>{"\u2715"}</div>

				{tools.map((tool, id) => {
					let className = `${styles.toolbarDrawerTool} ${toolId == id ? `${styles.selected}` : ""}`
					return 	<div className={className}
								onClick={onClickToolsMenuPattern.bind(this, id)}
								key={id} ><span>{tool.name}</span></div>
				})}
			</div>
}

const tools = [
	{
		name: "Point",
		func: ({x,y}) => [{x, y}]
	},
	{
		name: "Block",
		func: ({x,y}) => [{x:x,y:y},{x:x,y:y+1},{x:x+1,y:y},{x:x+1,y:y+1}]
	},
	{
		name: "Beehive",
		func: ({x,y}) => [{x:x,y:y},{x:x+1,y:y-1},{x:x+1,y:y+1},{x:x+2,y:y-1},{x:x+2,y:y+1},{x:x+3,y:y}]
	},
	{
		name: "Loaf",
		func: ({x,y}) => [{x:x,y:y},{x:x+1,y:y-1},{x:x+1,y:y+1},{x:x+2,y:y-1},{x:x+2,y:y+2},{x:x+3,y:y},{x:x+3,y:y+1},{x:x+3,y:y+2}]
	},
	{
		name: "Boat",
		func: ({x,y}) => [{x:x,y:y},{x:x,y:y+1},{x:x+1,y:y},{x:x+1,y:y+2},{x:x+2,y:y+1}]
	},
	{
		name: "Tub",
		func: ({x,y}) => [{x:x,y:y},{x:x+1,y:y-1},{x:x+1,y:y+1},{x:x+2,y:y}]
	},
	{
		name: "Blinker",
		func: ({x,y}) => [{x:x,y:y},{x:x+1,y:y},{x:x+2,y:y}]
	},
	{
		name: "Toad",
		func: ({x,y}) => [{x:x,y:y},{x:x,y:y+1},{x:x+1,y:y},{x:x+1,y:y+1},{x:x+2,y:y},{x:x-1 ,y:y+1}]
	},
	{
		name: "Beacon",
		func: ({x,y}) => [{x:x,y:y},{x:x,y:y+1},{x:x+1,y:y},{x:x+1,y:y+1},{x:x+2,y:y+2},{x:x+2,y:y+3},{x:x+3,y:y+2},{x:x+3 ,y:y+3}]
	},
	{
		name: "Pulsar",
		func: ({x,y}) => [{x:x,y:y},{x:x,y:y+1},{x:x,y:y+2},{x:x,y:y+6},{x:x,y:y+7},{x:x,y:y+8},{x:x+2,y:y-2},{x:x+2,y:y+3},{x:x+2,y:y+5},{x:x+2,y:y+10},{x:x+3,y:y-2},{x:x+3,y:y+3},{x:x+3,y:y+5},{x:x+3,y:y+10},{x:x+4,y:y-2},{x:x+4,y:y+3},{x:x+4,y:y+5},{x:x+4,y:y+10},{x:x+5,y:y},{x:x+5,y:y+1},{x:x+5,y:y+2},{x:x+5,y:y+6},{x:x+5,y:y+7},{x:x+5,y:y+8},{x:x+7,y:y},{x:x+7,y:y+1},{x:x+7,y:y+2},{x:x+7,y:y+6},{x:x+7,y:y+7},{x:x+7,y:y+8},{x:x+8,y:y-2},{x:x+8,y:y+3},{x:x+8,y:y+5},{x:x+8,y:y+10},{x:x+9,y:y-2},{x:x+9,y:y+3},{x:x+9,y:y+5},{x:x+9,y:y+10},{x:x+10,y:y-2},{x:x+10,y:y+3},{x:x+10,y:y+5},{x:x+10,y:y+10},{x:x+12,y:y},{x:x+12,y:y+1},{x:x+12,y:y+2},{x:x+12,y:y+6},{x:x+12,y:y+7},{x:x+12,y:y+8}]
	},
	{
		name: "Glider",
		func: ({x,y}) => [{x:x,y:y},{x:x+1,y:y},{x:x+1,y:y-2},{x:x+2,y:y},{x:x+2,y:y-1}]
	},
	{
		name: "LWSS",
		func: ({x,y}) => [{x:x,y:y},{x:x,y:y+1},{x:x+1,y:y-1},{x:x+1,y:y},{x:x+1,y:y+1},{x:x+2,y:y-1},{x:x+2,y:y},{x:x+2,y:y+2},{x:x+3,y:y},{x:x+3,y:y+1},{x:x+3,y:y+2},{x:x+4,y:y+1}]
	},
	{
		name: "MWSS",
		func: ({x,y}) => [{x:x,y:y},{x:x,y:y+1},{x:x+1,y:y-1},{x:x+1,y:y},{x:x+1,y:y+1},{x:x+2,y:y-1},{x:x+2,y:y},{x:x+2,y:y+1},{x:x+3,y:y-1},{x:x+3,y:y},{x:x+3,y:y+2},{x:x+4,y:y},{x:x+4,y:y+1},{x:x+4,y:y+2},{x:x+5,y:y+1}]
	}
]


export {
	Toolbar,
	ToolbarDrawer,
	tools
}