import React from "react"
import ReactDOM from "react-dom"

import {App} from "Components/App"

document.addEventListener("DOMContentLoaded", function(e) {
	var mount = document.getElementById("react-mount")
	ReactDOM.render(<App />, mount)
})
