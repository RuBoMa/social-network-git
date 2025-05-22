let socket;

export function initWebSocket(userId, onMessage) {
    if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) {
        console.log("WebSocket already connected");
        return;
    }
    console.log("Initializing WebSocket connection...");
    socket = new WebSocket(`ws://localhost:8080/ws`);

    socket.onopen = () => {
        console.log("WebSocket connected");
        socket.send(JSON.stringify({ type: "connect", sender_id: userId }));
    };

    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (onMessage) onMessage(data);
        console.log("Message received:", data);
    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };

    socket.onclose = (event) => {
        console.log("WebSocket closed:", event.reason);
    };

    console.log("WebSocket state:", socket.readyState);
}

export function sendMessage(message) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(message));
        console.log("Message sent:", message);
    } else {
        console.warn("WebSocket is not open. Message not sent:", message);
    }
}

export function closeWebSocket() {
    if (socket) {
        socket.close();
        socket = null;
        console.log("WebSocket connection closed");
    }
}
