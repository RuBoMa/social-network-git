let socket;

export default function initWebSocket(onMessage) {
    console.log("Initializing WebSocket connection...");

    if (!socket || socket.readyState === WebSocket.CLOSED) {
        // Get token from local storage
        const token = localStorage.getItem("token");
        if (!token) {
            console.error("No token found in local storage. Cannot establish WebSocket connection.");
            return;
        }
        console.log("Creating new WebSocket instance");
        socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

        socket.addEventListener("open", () => {
            console.log("WebSocket connected");
        });

        socket.addEventListener("message", (event) => {
            const data = JSON.parse(event.data);
            if (onMessage) onMessage(data);
            console.log("Message received:", data);
        });

        socket.addEventListener("error", (error) => {
            console.error("WebSocket error:", error);
        });

        socket.addEventListener("close", (event) => {
            console.log("WebSocket closed:", event.reason);
            closeWebSocket();
        });
    }

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
