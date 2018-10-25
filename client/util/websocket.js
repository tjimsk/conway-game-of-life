var W3CWebSocket = require('websocket').w3cwebsocket;

function newClient() {
    var websocketAddr = 'ws://'+window.location.host+'/websocket'
    var client = new W3CWebSocket(websocketAddr, 'echo-protocol');

    return client
}

export {
	newClient
}