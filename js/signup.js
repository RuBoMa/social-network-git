let signupFormListenerAdded = false;  // Track if the listener has been added

function renderSignupPage() {
    const signupPage = document.getElementById('signup-form');
    signupPage.innerHTML = '';
    
    signupPage.innerHTML = `
        <label for="username-signup">Username
            <div class="hover-icon">
                <span class="material-symbols-outlined" style="font-size: 20px; vertical-align: middle;">info</span>
                <span class="tooltip">Username must be 3-20 characters, letters, numbers, or _</span>
            </div>
        </label>
        <input type="text" id="username-signup" name="username" placeholder="Enter your username" required>
        <label for="age">Age</label>
        <input type="number" id="age" name="age" placeholder="Enter your age" required>
        <label for="signup-gender">Gender</label>
        <select id="signup-gender" name="gender">
            <option value="">Select Gender</option>
            <option value="male">Male</option>
            <option value="female">Female</option>
            <option value="non-binary">Non-binary</option>
            <option value="other">Other</option>
            <option value="prefer not to say">Prefer not to say</option>
        </select> 
        <label for="first-name">First Name</label>
        <input type="text" id="first-name" name="first-name" placeholder="Enter your first name" required>
        <label for="last-name">Last Name</label>
        <input type="text" id="last-name" name="last-name" placeholder="Enter your last name" required>
        <label for="email">E-mail</label>
        <input type="email" id="email" name="email" placeholder="Enter your e-mail" required>
        <label for="password-signup">Password</label>
        <input type="password" id="password-signup" name="password" placeholder="Enter your password" required>
        <label for="confirm-password">Re-enter Password</label>
        <input type="password" id="confirm-password" name="confirm-password" placeholder="Re-enter your password" required>
        <p id="password-error" style="color: red; display: none;">Passwords don't match</p>
        <p id="signup-error" style="color: red; display: none;"></p>
        <button type="submit">Sign Up</button>
    `;
    


const signupForm = document.getElementById("signup-page");
const passwordError = document.getElementById("password-error"); 

if (!signupFormListenerAdded) {
signupForm.addEventListener("submit", function(event) {
    event.preventDefault();  

    const formData = {
        username: document.getElementById("username-signup").value,
        age: document.getElementById("age").value,
        gender: document.getElementById("signup-gender").value,
        firstName: document.getElementById("first-name").value,
        lastName: document.getElementById("last-name").value,
        email: document.getElementById("email").value,
        password: document.getElementById("password-signup").value,
        confirmPassword: document.getElementById("confirm-password").value,
    };

    // check if passwords match
    if (formData.password !== formData.confirmPassword) {
        passwordError.style.display = 'block'; // Show the error message
        document.getElementById("confirm-password").classList.add('error'); // Add error class to confirm password field
        document.getElementById("confirm-password").classList.remove('success'); // Remove success class if exists
        return;  
    } else {
        // passwordError.style.display = 'none'; // Hide the error message
        document.getElementById("confirm-password").classList.remove('error');
        document.getElementById("confirm-password").classList.add('success'); // Add success class
    }

    apiPOST('/api/signup', 'signup', formData) 
    signupFormListenerAdded = true; 
});
}


// Handle the link to sign-up page
document.getElementById('login-link').addEventListener('click', (event) => {
    event.preventDefault();
    history.pushState({}, '', '/login');
    init();
});
}