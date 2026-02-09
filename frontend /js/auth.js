document.addEventListener('DOMContentLoaded', () => {
    // 1. Check if already logged in
    const token = localStorage.getItem('token');
    const role = localStorage.getItem('role');
    
    if (token) {
        // FIX: Removed '/static/' prefix. 
        // If your Go server serves files at root, use 'admin.html' or '/admin.html'
        if (role === 'admin') window.location.href = 'admin.html';
        else window.location.href = 'dashboard.html';
        return;
    }

    // Toggle Forms
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');

    if(document.getElementById('show-register')) {
        document.getElementById('show-register').addEventListener('click', (e) => {
            e.preventDefault();
            loginForm.classList.add('hidden');
            registerForm.classList.remove('hidden');
        });
    }

    if(document.getElementById('show-login')) {
        document.getElementById('show-login').addEventListener('click', (e) => {
            e.preventDefault();
            registerForm.classList.add('hidden');
            loginForm.classList.remove('hidden');
        });
    }

    // 2. Login Handler
    if(loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const email = document.getElementById('login-email').value;
            const password = document.getElementById('login-password').value;

            try {
                // Ensure this endpoint matches your Go Backend routes exactly
                const res = await api.post('/api/user/signIn', { email, password });
                
                // Decode JWT to find role securely
                let role = 'user'; // Default
                if(res.token) {
                    try {
                        const base64Url = res.token.split('.')[1];
                        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
                        const payload = JSON.parse(window.atob(base64));
                        role = payload.role;
                        
                        // Store detailed info
                        localStorage.setItem('user_id', payload.user_id);
                        localStorage.setItem('role', payload.role);
                    } catch (decodeErr) {
                        console.error("Token decode error", decodeErr);
                    }
                    localStorage.setItem('token', res.token);
                }

                alert('Login successful!'); // Simple alert for immediate feedback
                
                // FIX: Redirect without /static/
                if (role === 'admin') window.location.href = 'admin.html';
                else window.location.href = 'dashboard.html';

            } catch (err) {
                console.error(err);
                alert('Login Failed: ' + (err.message || "Unknown error"));
            }
        });
    }

    // 3. Register Handler
    if(registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            // Note: Ensure IDs match your HTML (reg-fullname vs reg-name)
            const full_name = document.getElementById('reg-fullname') ? document.getElementById('reg-fullname').value : document.getElementById('reg-name').value;
            const email = document.getElementById('reg-email').value;
            const password = document.getElementById('reg-password').value;

            try {
                await api.post('/api/user/signUp', { full_name, email, password });
                alert('Registration successful! Please login.');
                document.getElementById('show-login').click();
            } catch (err) {
                alert('Registration Failed: ' + err.message);
            }
        });
    }
});