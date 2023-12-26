let conn;

document.addEventListener("DOMContentLoaded", function() {
    // Connect to WebSocket server
    conn = new WebSocket("ws://localhost:8000/ws");

    // Handle incoming messages
    conn.onmessage = function(evt) {
        const messages = document.getElementById('messages');
        const msg = document.createElement('li');

        msg.appendChild(document.createTextNode(evt.data));
        messages.appendChild(msg);
    };

    conn.onerror = function (evt) {
        console.error('WebSocket error observed:', evt);
    }

    // Log WebSocket connection status
    conn.onopen = function(evt) {
        console.log('Connected to websocket');
    }

    conn.onclose = function(evt) {
        console.log('Disconnected from websocket');
    }

    // Send message to server
    function sendMessage() {
        const messageInput = document.getElementById('messageInput');

        if (messageInput.value) {
            conn.send(messageInput.value);
            messageInput.value = '';
        }
    }

    const messageButton = document.getElementById('sendMessage');
    messageButton.addEventListener('click', function() {
        sendMessage();
    });
})