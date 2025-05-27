let socket;
let messageHandlers = [];

export default function initWebSocket() {
    console.log("Initializing WebSocket connection...");

    if (!socket || socket.readyState === WebSocket.CLOSED) {
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
            console.log("Raw WebSocket message received:", data);
            
            // Notify all registered handlers
            messageHandlers.forEach(handler => {
                try {
                    handler(data);
                } catch (error) {
                    console.error("Error in message handler:", error);
                }
            });
        });

        socket.addEventListener("error", (error) => {
            console.error("WebSocket error:", error);
        });

        socket.addEventListener("close", (event) => {
            console.log("WebSocket closed:", event.reason);
            socket = null;
        });
    }

    console.log("WebSocket state:", socket.readyState);
    
    // Return cleanup function
    return () => {
        messageHandlers = [];
    };
}

// Function to add message handlers
export function addMessageHandler(handler) {
    if (!messageHandlers.includes(handler)) {
        messageHandlers.push(handler);
        console.log("Message handler added. Total handlers:", messageHandlers.length);
    }
    
    // Return cleanup function to remove this specific handler
    return () => {
        const index = messageHandlers.indexOf(handler);
        if (index > -1) {
            messageHandlers.splice(index, 1);
            console.log("Message handler removed. Total handlers:", messageHandlers.length);
        }
    };
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