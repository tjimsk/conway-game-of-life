var W3CWebSocket = require('websocket').w3cwebsocket;

var client = new W3CWebSocket('ws://localhost:8080/websocket', 'echo-protocol');

client.onerror = function() {
    console.log('Connection Error');
};

client.onopen = function() {
    console.log('WebSocket Client Connected');

    // function sendNumber() {
    //     if (client.readyState === client.OPEN) {
    //         var number = Math.round(Math.random() * 0xFFFFFF);
    //         client.send(number.toString());
    //         setTimeout(sendNumber, 1000);
    //     }
    // }

    // sendNumber();
};

client.onclose = function() {
    console.log('echo-protocol Client Closed');

    // retry

};

// client.onmessage = function(e) {
// };

class Message {
    constructor(type, content) {
        this.type = type
        this.content = content
    }
}

function sendMessage(msg) {
    client.send(msg)
}

export {
	client
}