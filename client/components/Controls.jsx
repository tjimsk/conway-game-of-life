import React from "react"

const styles = require("Components/Controls.scss")

class Controls extends React.Component {
	render() {
		return (
			<div className={styles.controls}>
				<div className={styles.header}>Conway's Game of Life</div>

				<table>
					<tbody className={styles.section}>
						<tr className={styles.title}>
							<td colSpan={2}>User</td>
						</tr>

						<tr className={styles.field}>
							<td className={styles.label}>Name</td>
							<td><input type="text" /></td>
						</tr>

						<tr className={styles.field}>
							<td className={styles.label}>Color</td>
							<td><input type="text" /></td>
						</tr>
					</tbody>

					<tbody className={styles.section}>
						<tr className={styles.title}>
							<td colSpan={2}>Display</td>
						</tr>

						<tr className={styles.field}>
							<td className={styles.label}>Cell Size</td>
							<td>
								<select defaultValue={8}>
									<option value={8}>8px</option>
									<option value={12}>12px</option>
									<option value={16}>16px</option>
								</select>
							</td>
						</tr>
					</tbody>

					<tbody className={styles.section}>
						<tr className={styles.title}>
							<td colSpan={2}>Settings</td>
						</tr>

						<tr className={styles.field}>
							<td className={styles.label}>Automatic Evolution</td>
							<td>
								<select defaultValue={false}>
									<option value={true}>Yes</option>
									<option value={false}>No</option>
								</select>
							</td>
						</tr>

						<tr className={styles.field}>
							<td className={styles.label}>Interval</td>
							<td>
								<select defaultValue={1}>
									<option value={0.5}>0.5s</option>
									<option value={1}>1s</option>
									<option value={2.5}>2.5s</option>
									<option value={5}>5s</option>
									<option value={10}>10s</option>
								</select>
							</td>
						</tr>
					</tbody>
				</table>

				<div className={styles.log}>
					<div className={styles.output}>
					</div>

					<div className={styles.input}>
						<input type="text" /><input type="button" value="Send" />
					</div>
				</div>
			</div>)
	}
}

export {
	Controls
}