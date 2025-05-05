// Function to refresh session by sending a request every 10 minutes
const refreshSession = async () => {
    try {
        const response = await fetch('/api/refresh-session', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        });

        const result = await response.json();
        if (response.ok) {
            console.log(result.message);
        }
    } catch (error) {
        console.error('Error refreshing session:', error);
    }
};

// Refresh the session every 10 minutes
setInterval(refreshSession, 10 * 60 * 1000);
