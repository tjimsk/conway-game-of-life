import React from "react"
import ReactDOM from "react-dom"

import {App} from "Components/App"

import "Assets/pause.png"
import "Assets/resume.png"
import "Assets/reset.png"
import "Assets/block.png"
import "Assets/beehive.png"
import "Assets/loaf.png"
import "Assets/boat.png"
import "Assets/tub.png"
import "Assets/blinker.png"
import "Assets/toad.png"
import "Assets/beacon.png"
import "Assets/pulsar.png"
import "Assets/glider.png"
import "Assets/lwss.png"
import "Assets/mwss.png"

document.addEventListener("DOMContentLoaded", function(e) {
	var mount = document.getElementById("react-mount")
	ReactDOM.render(<App />, mount)
})
