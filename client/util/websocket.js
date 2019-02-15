var W3CWebSocket = require('websocket').w3cwebsocket;

export default function websocket() {
    var websocketAddr = 'ws://'+window.location.host+'/websocket'
    var client = new W3CWebSocket(websocketAddr, 'echo-protocol');

    return client
}
