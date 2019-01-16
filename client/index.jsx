import React from "react"
import ReactDOM from "react-dom"

import {App} from "Components/App"

import pauseImg from "Assets/pause.png"
import resumeImg from "Assets/resume.png"
import blockImg from "Assets/block.png"
import beehiveImg from "Assets/beehive.png"
import loafImg from "Assets/loaf.png"
import boatImg from "Assets/boat.png"
import tubImg from "Assets/tub.png"
import blinkerImg from "Assets/blinker.png"
import toadImg from "Assets/toad.png"
import beaconImg from "Assets/beacon.png"
import pulsarImg from "Assets/pulsar.png"
import gliderImg from "Assets/glider.png"
import lwssImg from "Assets/lwss.png"
import mwssImg from "Assets/mwss.png"

document.addEventListener("DOMContentLoaded", function(e) {
	var mount = document.getElementById("react-mount")
	ReactDOM.render(<App />, mount)
})
