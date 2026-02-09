/* =========================================
   GLOBAL CONFIG & STATE
   ========================================= */
   const API = {
    CARS: '/api/cars',
    RENTALS: '/api/user/rentals', 
    PAYMENTS: '/api/payments',
    SALES: '/api/sales', 
    BUY_CAR: '/api/car/buy', // Correct endpoint for buying
    RECOMMEND: '/api/recommendation' 
};

let currentRentPrice = 0; 
let recQuestions = [];
let userAnswers = {};

const GENERIC_SUGGESTIONS = {
    "SUV": { desc: "Great for families/off-road.", examples: "Toyota RAV4, Honda CR-V" },
    "Sedan": { desc: "Efficient city driving.", examples: "Toyota Camry, Honda Accord" },
    "Hatchback": { desc: "Compact and versatile.", examples: "Volkswagen Golf, Ford Focus" },
    "Electric": { desc: "Eco-friendly tech.", examples: "Tesla Model 3, Nissan Leaf" }
};

/* =========================================
   INITIALIZATION
   ========================================= */
document.addEventListener('DOMContentLoaded', () => {
    console.log("Dashboard JS: Initializing...");

    // 1. Auth Check
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = 'index.html';
        return;
    }

    // 2. Setup UI
    setupTabs();
    
    // 3. Logout
    const logoutBtn = document.getElementById('logout-btn');
    if(logoutBtn) logoutBtn.onclick = window.logout;

    // 4. Load Data
    loadMarketplace();
    loadUserCars();
    loadRentals();
    loadPayments();
    loadRecommendationQuestions(); 
    
    // 5. ATTACH LISTENERS (Using helper to prevent duplicates)
    attachFormListener('create-car-form', handleListCar);
    attachFormListener('rent-form', handleRentSubmit);
    attachFormListener('payment-form', handlePaymentSubmit);
    
    // 6. Asset Forms
    setupAssetForm('form-maintenance', 'maintenance');
    setupAssetForm('form-fuel', 'fuel');
    setupAssetForm('form-expense', 'expenses');
    setupAssetForm('form-doc', 'documents');
});

// Helper for Forms
function attachFormListener(id, handler) {
    const form = document.getElementById(id);
    if (form) {
        const newForm = form.cloneNode(true);
        form.parentNode.replaceChild(newForm, form);
        newForm.addEventListener('submit', handler);
        console.log(`Listener attached to #${id}`);
    }
}

// Global Logout
window.logout = function() {
    localStorage.clear();
    window.location.href = 'index.html';
};

/* =========================================
   FEATURE: CREATE CAR
   ========================================= */
async function handleListCar(e) {
    e.preventDefault(); 
    const formData = new FormData(e.target);
    const rawData = Object.fromEntries(formData.entries());

    const year = parseInt(rawData.year);
    const price = parseFloat(rawData.price);
    const ownerId = parseInt(localStorage.getItem('user_id'));

    if (isNaN(year) || isNaN(price) || isNaN(ownerId)) {
        alert("Error: Invalid Input. Please check numbers.");
        return;
    }

    const payload = {
        brand: rawData.brand,
        model: rawData.model,
        year: year,
        price: price,
        status: rawData.status,
        owner_id: ownerId
    };

    try {
        await api.post(API.CARS, payload);
        alert('Car listed successfully!');
        e.target.reset(); 
        // Switch to Garage tab
        document.querySelector('[data-tab="cars"]').click(); 
        setTimeout(loadUserCars, 300);
    } catch (err) {
        alert('Failed to list car: ' + err.message);
    }
}

/* =========================================
   FEATURE: MARKETPLACE & BUYING
   ========================================= */
async function loadMarketplace() {
    const container = document.getElementById('market-grid');
    if (!container) return;
    const userId = parseInt(localStorage.getItem('user_id'));
    
    try {
        const cars = await api.get(API.CARS);
        const marketCars = cars.filter(c => c.owner_id !== userId && (c.status === 'for_sale' || c.status === 'for_rent'));
        container.innerHTML = marketCars.length ? marketCars.map(c => renderCarCard(c, 'market')).join('') : '<p>No cars available.</p>';
    } catch (e) { container.innerHTML = '<p>Error loading market.</p>'; }
}

async function loadUserCars() {
    const container = document.getElementById('cars-grid');
    if (!container) return;
    const userId = parseInt(localStorage.getItem('user_id'));
    
    try {
        const cars = await api.get(API.CARS);
        const myCars = cars.filter(c => c.owner_id === userId);
        container.innerHTML = myCars.length ? myCars.map(c => renderCarCard(c, 'mine')).join('') : '<p>No cars in your garage.</p>';
    } catch (e) { console.error(e); }
}

function renderCarCard(car, type) {
    let btns = '';
    if (type === 'market') {
        if (car.status === 'for_rent') btns = `<button class="btn-sm" style="background:blue; color:white;" onclick="openRentModal(${car.car_id}, ${car.price})">Rent</button>`;
        // UPDATED: Calls buyCar correctly
        if (car.status === 'for_sale') btns = `<button class="btn-sm" style="background:green; color:white;" onclick="buyCar(${car.car_id}, ${car.price})">Buy</button>`;
    } else {
        btns = `<button class="btn-sm" style="background:#444; color:white;" onclick="openAssetModal(${car.car_id})">Manage Assets</button>`;
    }

    return `
    <div class="card">
        <h3>${car.brand} ${car.model}</h3>
        <p>Year: ${car.year} | Price: $${car.price}</p>
        <span style="background:#eee; padding:2px 5px; border-radius:3px; font-size:12px;">${car.status.toUpperCase()}</span>
        <div style="margin-top:10px;">${btns}</div>
    </div>`;
}

// *** CRITICAL FIX: BUY CAR LOGIC ***
window.buyCar = async (carId, price) => {
    if(!confirm(`Buy for $${price}?`)) return;
    try {
        // Send JSON { "car_id": 123 } to /api/car/buy
        await api.post(API.BUY_CAR, { car_id: parseInt(carId) });
        
        alert('Purchase Successful!'); 
        loadMarketplace(); 
        loadUserCars();
        loadPayments(); // Refresh balance
    } catch (err) { 
        console.error("Buy failed", err);
        alert('Transaction Failed: ' + (err.message || "Unknown error")); 
    }
};

/* =========================================
   FEATURE: ASSETS & MODALS
   ========================================= */
const setupAssetForm = (id, type) => {
    attachFormListener(id, async (e) => {
        e.preventDefault();
        const carId = document.getElementById('modal-car-id').value;
        const formData = new FormData(e.target);
        const data = Object.fromEntries(formData.entries());
        
        ['cost', 'amount', 'price', 'liters', 'mileage'].forEach(k => { if(data[k]) data[k] = parseFloat(data[k]); });
        if (type === 'maintenance' && data.service_date && !data.service_date.includes('T')) data.service_date += 'T00:00:00Z';

        try {
            await api.post(`/api/cars/${carId}/${type}`, data);
            alert('Record added!');
            e.target.reset();
            loadAssetHistory(carId, type, type === 'expenses' ? 'list-expense' : `list-${type}`);
        } catch(err) { alert('Error: ' + err.message); }
    });
};

window.openAssetModal = (carId) => {
    document.getElementById('modal-car-id').value = carId;
    document.getElementById('car-modal').classList.remove('hidden');
    ['maintenance', 'fuel', 'expenses', 'documents'].forEach(type => {
        loadAssetHistory(carId, type, type === 'expenses' ? 'list-expense' : `list-${type}`);
    });
};

async function loadAssetHistory(carId, type, listId) {
    const list = document.getElementById(listId);
    if (!list) return;
    list.innerHTML = '<li>Loading...</li>';
    try {
        const data = await api.get(`/api/cars/${carId}/${type}`);
        list.innerHTML = (data && data.length) ? data.map(item => {
            const d = (item.service_date || item.fill_date || item.expense_date || item.expiry_date || "").split('T')[0];
            return `<li>${d}: ${item.service_type || item.expense_type || "Record"} - $${item.cost || item.amount || item.price || 0}</li>`;
        }).join('') : '<li>No records.</li>';
    } catch (e) { list.innerHTML = '<li>Error loading.</li>'; }
}

/* =========================================
   FEATURE: PAYMENTS & RENTALS
   ========================================= */
// *** FIX: Open Payment Modal ***
window.openPaymentModal = () => {
    document.getElementById('payment-modal').classList.remove('hidden');
};

async function handlePaymentSubmit(e) {
    e.preventDefault();
    const data = { 
        amount: parseFloat(e.target.amount.value), 
        payment_type: e.target.payment_type.value,
        related_id: parseInt(document.getElementById('payment-ref-id').value) || 0
    };
    try { 
        await api.post(API.PAYMENTS, data); 
        alert('Payment Recorded!'); 
        closeModal('payment-modal'); 
        loadPayments(); 
    } catch (err) { alert(err.message); }
}

async function loadPayments() {
    try {
        const payments = await api.get(API.PAYMENTS);
        const tbody = document.getElementById('payments-table-body');
        if(tbody) tbody.innerHTML = payments.map(p => `<tr><td>${p.payment_id}</td><td>$${p.amount}</td><td>${p.payment_type}</td><td>${new Date(p.created_at).toLocaleDateString()}</td></tr>`).join('');
        
        // Mock Balance Update (Since we assume user starts with 1M)
        // If your backend API returns user balance in a dedicated endpoint, use that here.
        document.getElementById('nav-balance').innerText = "1,000,000 (DB)"; 
        document.getElementById('wallet-balance-display').innerText = "$ Check DB";
    } catch (e) {}
}

// Rental Logic
window.openRentModal = (id, price) => { 
    document.getElementById('rent-car-id').value = id; 
    currentRentPrice = parseFloat(price); 
    document.getElementById('rent-modal').classList.remove('hidden'); 
};

async function handleRentSubmit(e) {
    e.preventDefault();
    const data = { 
        car_id: parseInt(document.getElementById('rent-car-id').value), 
        start_date: e.target.start_date.value, 
        end_date: e.target.end_date.value,
        daily_price: currentRentPrice || 50 
    };
    try { 
        await api.post(API.RENTALS, data); 
        alert('Rental Request Sent!'); 
        closeModal('rent-modal'); 
        loadRentals(); 
    } catch (err) { alert(err.message); }
}

async function loadRentals() {
    try {
        const rentals = await api.get(API.RENTALS);
        const tbody = document.getElementById('rentals-table-body');
        if(tbody) tbody.innerHTML = rentals.map(r => `<tr><td>${r.rental_id}</td><td>${r.car_id}</td><td>${r.start_date.split('T')[0]}</td><td>${r.end_date.split('T')[0]}</td><td>$${r.total_price}</td></tr>`).join('');
    } catch (e) {}
}

/* =========================================
   UTILITIES
   ========================================= */
window.closeModal = (id) => document.getElementById(id).classList.add('hidden');

window.showMiniTab = (id) => {
    document.querySelectorAll('.mini-content').forEach(e => e.classList.add('hidden'));
    document.getElementById(id).classList.remove('hidden');
};

function setupTabs() {
    document.querySelectorAll('.tab-link').forEach(btn => {
        btn.addEventListener('click', () => {
            document.querySelectorAll('.tab-link').forEach(b => b.classList.remove('active'));
            document.querySelectorAll('.tab-content').forEach(c => c.classList.add('hidden'));
            btn.classList.add('active');
            document.getElementById(btn.getAttribute('data-tab')).classList.remove('hidden');
        });
    });
}

// Recommendations
async function loadRecommendationQuestions() {
    try {
        recQuestions = await api.get(`${API.RECOMMEND}/questions`);
        const container = document.getElementById('questions-container');
        if (container) container.innerHTML = recQuestions.map(q => `<div><label>${q.question}</label><select onchange="saveAnswer('${q.id}', this.value)"><option value="">Select...</option>${q.options.map(o => `<option value="${o}">${o}</option>`).join('')}</select></div>`).join('');
    } catch (e) {}
}
window.saveAnswer = (qId, val) => { userAnswers[qId] = val; };
window.getRecommendations = async () => {
    try {
        const res = await api.post(`${API.RECOMMEND}/result`, { answers: userAnswers });
        document.getElementById('recommendation-results').innerHTML = `<p>AI Suggests: ${res.recommended_type}</p>` + (res.cars || []).map(c => `<div>${c.brand} ${c.model}</div>`).join('');
    } catch (e) { alert(e.message); }
};