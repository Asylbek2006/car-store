// frontend/js/api.js
const BASE_URL = 'http://localhost:8000';

class Api {
    getToken() {
        return localStorage.getItem('token');
    }

    async request(endpoint, method = 'GET', body = null) {
        const headers = {
            'Content-Type': 'application/json'
        };

        const token = this.getToken();
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const config = { method, headers };
        if (body) config.body = JSON.stringify(body);

        try {
            const response = await fetch(`${BASE_URL}${endpoint}`, config);
            
            if (response.status === 401) {
                localStorage.clear();
                window.location.href = 'index.html';
                return;
            }

            // CRITICAL CHANGE: Handle both JSON and Text responses
            const responseText = await response.text();
            let data;
            try {
                data = JSON.parse(responseText);
            } catch (e) {
                // If not JSON, use the raw text (Gin validation errors often come here)
                data = { error: responseText };
            }
            
            if (!response.ok) {
                // Throw the actual message from the server
                throw new Error(data.error || data.message || responseText || 'Request failed');
            }

            return data;
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    }

    get(url) { return this.request(url, 'GET'); }
    post(url, data) { return this.request(url, 'POST', data); }
    put(url, data) { return this.request(url, 'PUT', data); }
    delete(url) { return this.request(url, 'DELETE'); }
}

const api = new Api();

function showAlert(msg, type = 'error') {
    const box = document.getElementById('alert-box');
    if(box) {
        box.textContent = msg;
        box.className = `alert ${type}`;
        box.classList.remove('hidden');
        setTimeout(() => box.classList.add('hidden'), 5000); // Increased time to read error
    } else {
        alert(msg);
    }
}