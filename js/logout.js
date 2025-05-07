async function LogoutUser() {

    try {
        const response = await fetch('/api/logout', { 
            method: 'POST', 
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include' 
        });

        const data = await response.json();

        if (!response.ok) {
            showError(data.message);
            return
        } else {
        // Update UI
        if (socket !== null) socket.close(); socket = null;
        document.getElementById('logout-button').style.display = 'none';
        document.getElementById('chat-sidebar').style.display = 'none';
        localStorage.removeItem('username');
        const loggedInUserElement = document.getElementById('logged-in-user');
        if (loggedInUserElement) {
            loggedInUserElement.textContent = '';
        }
        history.pushState({}, '', '/login');
        }
        return
    } catch (error) {
        showError(data.message);
    }
    
}