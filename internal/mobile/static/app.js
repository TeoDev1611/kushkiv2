const state = {
    token: null,
    products: [],
    filtered: [],
    editingProduct: null,
    tempStock: 0,
    isPosMode: false,
    scanner: null
};

// DOM Elements
const views = {
    login: document.getElementById('login-view'),
    app: document.getElementById('app-view'),
    inventory: document.getElementById('inventory-view'),
    scanner: document.getElementById('scanner-view')
};

const dom = {
    pinInput: document.getElementById('pin-input'),
    btnLogin: document.getElementById('btn-login'),
    searchInput: document.getElementById('search-input'),
    list: document.getElementById('product-list'),
    
    // Header
    modeBadge: document.getElementById('mode-badge'),
    
    // Nav
    btnScan: document.getElementById('btn-scan'),
    btnTogglePos: document.getElementById('btn-toggle-pos'),
    posLabel: document.getElementById('pos-label'),
    
    // Modal
    modal: document.getElementById('edit-modal'),
    modalTitle: document.getElementById('modal-title'),
    modalSku: document.getElementById('modal-sku'),
    modalStock: document.getElementById('modal-stock'),
    modalLocation: document.getElementById('modal-location'),
    
    posControls: document.getElementById('pos-controls'),
    invControls: document.getElementById('inventory-controls'),
    posQtyDisplay: document.getElementById('pos-qty-display'), // New element
    
    btnCloseModal: document.getElementById('btn-close-modal'),
    btnSaveStock: document.getElementById('btn-save-stock'),
    btnSendPos: document.getElementById('btn-send-pos'),
    
    // Scanner
    btnStopScan: document.getElementById('btn-stop-scan'),
    
    toast: document.getElementById('toast')
};

// Init
window.addEventListener('DOMContentLoaded', () => {
    const urlParams = new URLSearchParams(window.location.search);
    const tokenParam = urlParams.get('token');

    if (tokenParam) {
        login(tokenParam);
    } else {
        const savedToken = localStorage.getItem('kushki_token');
        if (savedToken) login(savedToken);
    }
});

// Login
dom.btnLogin.addEventListener('click', () => {
    const pin = dom.pinInput.value.trim();
    if (pin.length === 6) login(pin);
    else showToast('PIN inválido');
});

// Search
dom.searchInput.addEventListener('input', (e) => filterProducts(e.target.value));

// Toggle POS Mode
dom.btnTogglePos.addEventListener('click', () => {
    state.isPosMode = !state.isPosMode;
    updateModeUI();
});

function updateModeUI() {
    if (state.isPosMode) {
        dom.modeBadge.innerText = "MODO CAJA";
        dom.modeBadge.style.color = "#34D399"; // Mint
        dom.posLabel.innerText = "Salir POS";
        dom.btnTogglePos.classList.add('active');
        showToast("Modo POS Activo: Escanee para vender");
    } else {
        dom.modeBadge.innerText = "INVENTARIO";
        dom.modeBadge.style.color = "#9ca3af"; // Gray
        dom.posLabel.innerText = "Modo POS";
        dom.btnTogglePos.classList.remove('active');
    }
}

// Scanner Logic
dom.btnScan.addEventListener('click', startScanner);
dom.btnStopScan.addEventListener('click', stopScanner);

async function startScanner() {
    views.inventory.classList.add('hidden');
    views.scanner.classList.remove('hidden');

    if (!state.scanner) {
        state.scanner = new Html5Qrcode("reader");
    }

    try {
        await state.scanner.start(
            { facingMode: "environment" }, 
            { fps: 10, qrbox: { width: 250, height: 250 } },
            onScanSuccess,
            (err) => { /* ignore failures */ }
        );
    } catch (err) {
        showToast("Error cámara: " + err);
        stopScanner();
    }
}

function stopScanner() {
    if (state.scanner) {
        state.scanner.stop().then(() => {
            views.scanner.classList.add('hidden');
            views.inventory.classList.remove('hidden');
        }).catch(err => console.error(err));
    } else {
        views.scanner.classList.add('hidden');
        views.inventory.classList.remove('hidden');
    }
}

function onScanSuccess(decodedText, decodedResult) {
    // Si estamos en modo POS, enviar directo
    if (state.isPosMode) {
        sendToPos(decodedText);
        showToast(`Escaneado: ${decodedText}`);
        // Pequeña pausa para no escanear el mismo mil veces
        state.scanner.pause();
        setTimeout(() => state.scanner.resume(), 1500);
    } else {
        // Modo Inventario: Buscar y abrir modal
        stopScanner();
        const product = state.products.find(p => p.Barcode === decodedText || p.SKU === decodedText);
        if (product) {
            openModal(product.SKU);
        } else {
            showToast("Producto no encontrado: " + decodedText);
        }
    }
}

// List Interactions
dom.list.addEventListener('click', (e) => {
    const card = e.target.closest('.product-card');
    if (card) {
        openModal(card.dataset.sku);
    }
});

// Modal Actions
function openModal(sku) {
    const p = state.products.find(x => x.SKU === sku);
    if (!p) return;

    state.editingProduct = p;
    state.tempStock = p.Stock;
    state.posQty = 1; // Reset POS Qty

    dom.modalTitle.innerText = p.Name;
    dom.modalSku.innerText = p.SKU;
    dom.modalStock.innerText = state.tempStock;
    dom.modalLocation.value = p.Location || "";
    if (dom.posQtyDisplay) dom.posQtyDisplay.innerText = "1";
    
    if (state.isPosMode) {
        dom.posControls.classList.remove('hidden');
        dom.invControls.classList.add('hidden');
    } else {
        dom.posControls.classList.add('hidden');
        dom.invControls.classList.remove('hidden');
    }

    dom.modal.classList.remove('hidden');
}

dom.btnCloseModal.addEventListener('click', () => dom.modal.classList.add('hidden'));

// Inventory Logic
document.querySelectorAll('.btn-qty').forEach(btn => {
    btn.addEventListener('click', () => {
        const delta = parseInt(btn.dataset.delta);
        state.tempStock += delta;
        if (state.tempStock < 0) state.tempStock = 0;
        dom.modalStock.innerText = state.tempStock;
    });
});

// POS Qty Logic
document.querySelectorAll('.btn-pos-qty').forEach(btn => {
    btn.addEventListener('click', () => {
        const delta = parseInt(btn.dataset.delta);
        state.posQty = (state.posQty || 1) + delta;
        if (state.posQty < 1) state.posQty = 1;
        if (dom.posQtyDisplay) dom.posQtyDisplay.innerText = state.posQty;
    });
});

dom.btnSaveStock.addEventListener('click', async () => {
    if (!state.editingProduct) return;
    
    // Aquí podríamos guardar también Location, pero el endpoint actual solo guarda stock.
    // TODO: Actualizar endpoint backend para aceptar Location.
    
    try {
        const res = await fetch('/api/stock', {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'X-Kushki-Token': state.token 
            },
            body: JSON.stringify({
                sku: state.editingProduct.SKU,
                quantity: state.tempStock,
                type: 'set'
            })
        });
        
        if (res.ok) {
            state.editingProduct.Stock = state.tempStock;
            renderList();
            dom.modal.classList.add('hidden');
            showToast('Guardado');
        }
    } catch (e) {
        showToast('Error de conexión');
    }
});

// POS Logic
dom.btnSendPos.addEventListener('click', () => {
    if (state.editingProduct) {
        sendToPos(state.editingProduct.SKU, state.posQty || 1);
        // Feedback visual
        dom.btnSendPos.innerText = "¡ENVIADO!";
        dom.btnSendPos.style.background = "#fff";
        dom.btnSendPos.style.color = "#000";
        setTimeout(() => {
            dom.btnSendPos.innerText = "🛒 ENVIAR A CAJA";
            dom.btnSendPos.style.background = "";
            dom.btnSendPos.style.color = "";
        }, 1000);
    }
});

async function sendToPos(sku, qty = 1) {
    try {
        await fetch('/api/pos/scan', {
            method: 'POST',
            headers: { 
                'Content-Type': 'application/json',
                'X-Kushki-Token': state.token 
            },
            body: JSON.stringify({ sku: sku, quantity: qty })
        });
    } catch (e) {
        showToast("Error enviando a POS");
    }
}

// Core Functions
async function login(token) {
    state.token = token;
    localStorage.setItem('kushki_token', token);
    try {
        await loadInventory();
        views.login.classList.add('hidden');
        views.app.classList.remove('hidden');
    } catch (e) {
        showToast('Error de conexión');
        views.login.classList.remove('hidden');
        views.app.classList.add('hidden');
    }
}

async function loadInventory() {
    const res = await fetch('/api/inventory', {
        headers: { 'X-Kushki-Token': state.token }
    });
    if (!res.ok) throw new Error();
    state.products = await res.json();
    state.filtered = state.products;
    renderList();
}

function renderList() {
    dom.list.innerHTML = '';
    state.filtered.forEach(p => {
        const el = document.createElement('div');
        el.className = `product-card ${p.Stock <= (p.MinStock || 0) ? 'low-stock' : ''}`;
        el.dataset.sku = p.SKU;
        el.innerHTML = `
            <div class="p-info">
                <h3>${p.Name}</h3>
                <p>${p.SKU} | $${p.Price.toFixed(2)}</p>
                <p style="font-size: 11px; color: #999;">${p.Location || 'Sin ubicación'}</p>
            </div>
            <div class="p-stock">${p.Stock}</div>
        `;
        dom.list.appendChild(el);
    });
}

function filterProducts(query) {
    const q = query.toLowerCase();
    state.filtered = state.products.filter(p => 
        p.Name.toLowerCase().includes(q) || 
        p.SKU.toLowerCase().includes(q) ||
        (p.Barcode && p.Barcode.includes(q))
    );
    renderList();
}

function showToast(msg) {
    dom.toast.innerText = msg;
    dom.toast.classList.remove('hidden');
    setTimeout(() => dom.toast.classList.add('hidden'), 3000);
}