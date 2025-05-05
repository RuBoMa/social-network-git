function renderLoginPage() {
    const loginPage = document.getElementById('login-form');
    loginPage.innerHTML = '';

    loginPage.innerHTML = `
            <label for="username-login">Username/e-mail</label>
            <input type="text" id="username-login" name="username" placeholder="Enter your username or e-mail" required>
            <label for="password-login">Password</label>
            <input type="password" id="password-login" name="password" placeholder="Enter your password" required>
            <p id="login-error" style="color: red; display: none;"></p>
            <button type="submit">Login</button>
    `;

        document.getElementById('login-form').addEventListener('submit', (event) => {
            event.preventDefault();


            const username = document.getElementById('username-login').value;
            const password = document.getElementById('password-login').value;

            const loginData = {
                username: username,
                password: password
            };

            apiPOST('/api/login', 'login', loginData)

});


    // Handle the link to sign-up page
    document.getElementById('signup-link').addEventListener('click', (event) => {
        event.preventDefault();
        history.pushState({}, '', '/signup');
        loadPage();
    });

}