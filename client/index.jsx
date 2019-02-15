import React from "react"
import ReactDOM from "react-dom"
import {App} from "Components/App"
import "Styles/App.scss"

document.addEventListener("DOMContentLoaded", function(e) {
	var mount = document.getElementById("react-mount")
	ReactDOM.render(<App />, mount)
})
