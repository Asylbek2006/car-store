/* =========================================
   1. GLOBAL LOGOUT FUNCTION
   ========================================= */
// This allows onclick="logout()" in HTML to work immediately
window.logout = function() {
    localStorage.clear();
    sessionStorage.clear();
    window.location.href = 'index.html';
};

document.addEventListener('DOMContentLoaded', () => {
    // Admin Check
    const role = localStorage.getItem('role');
    if (role !== 'admin') {
        window.location.href = 'index.html';
        return;
    }

    // --- FIX: Attach Logout Listener ---
    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', (e) => {
            e.preventDefault();
            window.logout();
        });
    }

    setupTabs();
    loadStats();
    
    // Lazy Loading Listeners
    const userTab = document.querySelector('[data-tab="users"]');
    if(userTab) userTab.addEventListener('click', loadUsers);

    const carTab = document.querySelector('[data-tab="all-cars"]');
    if(carTab) carTab.addEventListener('click', loadAllCars);

    const rentTab = document.querySelector('[data-tab="all-rentals"]');
    if(rentTab) rentTab.addEventListener('click', loadAllRentals);

    const saleTab = document.querySelector('[data-tab="all-sales"]');
    if(saleTab) saleTab.addEventListener('click', loadAllSales);

    const payTab = document.querySelector('[data-tab="all-payments"]');
    if(payTab) payTab.addEventListener('click', loadAllPayments);
});

/* =========================================
   2. DATA LOADING
   ========================================= */

async function loadStats() {
    try {
        // Ensure your API endpoint exists, otherwise this will fail silently
        const stats = await api.get('/api/admin/stats'); 
        
        const container = document.getElementById('stats-container');
        if(container) {
            container.innerHTML = `
                <div class="stat-card"><h3>Users</h3><div class="stat-val">${stats.users || 0}</div></div>
                <div class="stat-card"><h3>Cars</h3><div class="stat-val">${stats.cars || 0}</div></div>
                <div class="stat-card"><h3>Rentals</h3><div class="stat-val">${stats.total_rentals || 0}</div></div>
                <div class="stat-card"><h3>Revenue</h3><div class="stat-val">$${stats.total_revenue || 0}</div></div>
            `;
            drawChart(stats);
        }
    } catch (e) { console.error("Stats load error:", e); }
}

function drawChart(stats) {
    const canvas = document.getElementById('adminChart');
    if (!canvas || !canvas.getContext) return;
    const ctx = canvas.getContext('2d');
    
    // Fallback if data is missing
    const u = stats.users || 0;
    const c = stats.cars || 0;
    const r = stats.total_rentals || 0;

    const data = [u, c, r];
    const labels = ["Users", "Cars", "Rentals"];
    const max = Math.max(...data) + 5;
    
    ctx.clearRect(0,0,600,300);
    
    // Draw Axis
    ctx.beginPath();
    ctx.moveTo(30, 20); ctx.lineTo(30, 250); ctx.lineTo(500, 250); 
    ctx.strokeStyle = '#333';
    ctx.stroke();

    // Draw Bars
    data.forEach((val, i) => {
        const h = max > 0 ? (val / max) * 200 : 0;
        const x = 50 + (i * 120);
        const y = 250 - h;
        
        ctx.fillStyle = '#3b82f6';
        ctx.fillRect(x, y, 60, h);
        
        ctx.fillStyle = '#000';
        ctx.font = '14px Arial';
        ctx.fillText(labels[i], x + 10, 270);
        ctx.fillText(val, x + 20, y - 5);
    });
}

async function loadUsers() {
    try {
        const users = await api.get('/api/admin/users');
        const tbody = document.querySelector('#users-table tbody');
        if(tbody) tbody.innerHTML = users.map(u => `<tr><td>${u.user_id}</td><td>${u.username || u.full_name}</td><td>${u.email}</td><td>${u.role}</td></tr>`).join('');
    } catch(e) { console.error(e); }
}

async function loadAllCars() {
    try {
        const cars = await api.get('/api/admin/cars');
        const tbody = document.querySelector('#cars-table tbody');
        if(tbody) tbody.innerHTML = cars.map(c => `<tr><td>${c.car_id}</td><td>${c.owner_id}</td><td>${c.brand}</td><td>${c.model}</td><td>$${c.price}</td></tr>`).join('');
    } catch(e) { console.error(e); }
}

async function loadAllRentals() {
    try {
        const rentals = await api.get('/api/admin/rentals');
        const tbody = document.querySelector('#rentals-table tbody');
        if(tbody) tbody.innerHTML = rentals.map(r => `<tr><td>${r.rental_id}</td><td>${r.car_id}</td><td>${r.user_id || r.renter_id}</td><td>${r.start_date.split('T')[0]}</td><td>$${r.total_price}</td></tr>`).join('');
    } catch(e) { console.error(e); }
}

async function loadAllSales() {
    try {
        const sales = await api.get('/api/admin/sales');
        const tbody = document.querySelector('#sales-table tbody');
        if(tbody) tbody.innerHTML = sales.map(s => `<tr><td>${s.sale_id}</td><td>${s.car_id}</td><td>${s.buyer_id}</td><td>$${s.sale_price}</td><td>${s.sale_date.split('T')[0]}</td></tr>`).join('');
    } catch(e) { console.error(e); }
}

async function loadAllPayments() {
    try {
        const payments = await api.get('/api/admin/payments');
        const tbody = document.querySelector('#payments-table tbody');
        if(tbody) tbody.innerHTML = payments.map(p => `<tr><td>${p.payment_id}</td><td>${p.user_id}</td><td>$${p.amount}</td><td>${p.payment_type}</td><td>${p.created_at.split('T')[0]}</td></tr>`).join('');
    } catch(e) { console.error(e); }
}

function setupTabs() {
    document.querySelectorAll('.tab-link').forEach(btn => {
        btn.addEventListener('click', () => {
            document.querySelectorAll('.tab-link').forEach(b => b.classList.remove('active'));
            document.querySelectorAll('.tab-content').forEach(c => c.classList.add('hidden'));
            btn.classList.add('active');
            const targetId = btn.getAttribute('data-tab');
            const target = document.getElementById(targetId);
            if(target) target.classList.remove('hidden');
        });
    });
}