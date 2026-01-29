<script lang="ts">
    import { onMount } from 'svelte';
    import { Backend } from '$lib/services/api';
    import { notifications } from '$lib/stores/notifications';
    import type { db } from 'wailsjs/go/models';
    import { EventsOn } from '../../../../wailsjs/runtime/runtime';
    import * as WailsApp from 'wailsjs/go/main/App'; 
    import QRCode from 'qrcode';
    import ClientForm from '$lib/components/ClientForm.svelte';

    // Estado del POS
    let items: any[] = [];
    let barcodeInput = "";
    let searchInput = "";
    let isSearching = false;
    let searchResults: db.ProductDTO[] = [];
    let selectedClient: db.ClientDTO = {
        ID: "9999999999999",
        Nombre: "CONSUMIDOR FINAL",
        TipoID: "07",
        Direccion: "S/N",
        Email: "",
        Telefono: ""
    };
    
    // Client Search State
    let showClientSearch = false;
    let showNewClientForm = false;
    let newClientData = { ID: "", TipoID: "05", Nombre: "", Direccion: "", Email: "", Telefono: "" };
    let clientSearchTerm = "";
    let clientSearchResults: db.ClientDTO[] = [];

    // QR State
    let showQR = false;
    let satelliteInfo: any = null;
    
    let total = 0;
    let subtotal0 = 0;
    let subtotal15 = 0;
    let iva = 0;
    let secuencial = "";
    let inputElement: HTMLInputElement;

    $: calculateTotals(items);

    function calculateTotals(currentItems: any[]) {
        subtotal0 = 0;
        subtotal15 = 0;
        iva = 0;
        currentItems.forEach(item => {
            const sub = item.cantidad * item.precio;
            if (item.porcentajeIVA > 0) {
                subtotal15 += sub;
                iva += sub * (item.porcentajeIVA / 100);
            } else {
                subtotal0 += sub;
            }
        });
        total = subtotal0 + subtotal15 + iva;
    }

    async function handleBarcode(e: KeyboardEvent) {
        if (e.key === 'Enter' && barcodeInput.trim()) {
            await addProductByIdentifier(barcodeInput.trim());
            barcodeInput = "";
        }
    }

    async function addProductByIdentifier(id: string) {
        const products = await Backend.searchProducts(id);
        if (products && products.length > 0) {
            // Si hay coincidencia exacta por Barcode o SKU, tomamos el primero
            const exact = products.find(p => p.Barcode === id || p.SKU === id) || products[0];
            addItem(exact);
        } else {
            notifications.show("Producto no encontrado", "error");
        }
    }

    function addItem(prod: db.ProductDTO) {
        const existing = items.find(i => i.codigo === prod.SKU);
        if (existing) {
            existing.cantidad += 1;
            items = [...items];
        } else {
            items = [...items, {
                codigo: prod.SKU,
                nombre: prod.Name,
                cantidad: 1,
                precio: prod.Price,
                codigoIVA: prod.TaxCode,
                porcentajeIVA: prod.TaxPercentage
            }];
        }
        notifications.show(`${prod.Name} añadido`, "success");
    }

    function removeItem(index: number) {
        items.splice(index, 1);
        items = [...items];
    }

    function updateQuantity(index: number, delta: number) {
        if (items[index].cantidad + delta > 0) {
            items[index].cantidad += delta;
            items = [...items];
        }
    }

    async function processVenta() {
        if (items.length === 0) {
            notifications.show("No hay productos en la venta", "error");
            return;
        }

        try {
            const currentSec = await Backend.getNextSecuencial();
            const invoiceData: any = {
                secuencial: currentSec,
                clienteID: selectedClient.ID,
                clienteNombre: selectedClient.Nombre,
                clienteDireccion: selectedClient.Direccion,
                clienteEmail: selectedClient.Email,
                clienteTelefono: selectedClient.Telefono,
                observacion: "Venta POS",
                formaPago: "01", // Sin utilización del sistema financiero
                plazo: "0",
                unidadTiempo: "dias",
                items: items
            };

            const res = await Backend.createInvoice(invoiceData);
            if (res.startsWith("Éxito")) {
                notifications.show(res, "success");
                items = [];
                barcodeInput = "";
                inputElement.focus();
            } else {
                notifications.show(res, "error");
            }
        } catch (e) {
            notifications.show("Error procesando venta: " + e, "error");
        }
    }

    function handleGlobalKeyDown(e: KeyboardEvent) {
        if (isSearching || showClientSearch || showQR || showNewClientForm) return;

        if (e.key === 'F12' || e.key === '+') {
            e.preventDefault();
            processVenta();
        } else if (e.key === 'F5') {
            e.preventDefault();
            isSearching = true;
        } else if (e.key === 'Escape') {
            items = [];
            barcodeInput = "";
            notifications.show("Venta cancelada", "info");
        }
    }

    async function searchManual() {
        if (searchInput.length > 2) {
            searchResults = await Backend.searchProducts(searchInput);
        } else {
            searchResults = [];
        }
    }

    // --- Cliente Logic ---
    async function searchClients() {
        if (clientSearchTerm.length > 2) {
            clientSearchResults = await Backend.searchClients(clientSearchTerm);
        } else {
            clientSearchResults = [];
        }
    }

    function selectFoundClient(c: db.ClientDTO) {
        selectedClient = c;
        showClientSearch = false;
        clientSearchTerm = "";
        notifications.show("Cliente asignado", "info");
    }

    function onNewClientSaved(event: CustomEvent) {
        const c = event.detail;
        selectedClient = c;
        showNewClientForm = false;
        showClientSearch = false; // Close parent modal if open
        notifications.show("Nuevo cliente asignado", "info");
    }

    function openNewClientForm() {
        newClientData = { ID: "", TipoID: "05", Nombre: "", Direccion: "", Email: "", Telefono: "" };
        showNewClientForm = true;
    }

    // --- QR Logic ---
    async function openQR() {
        showQR = true;
        if (!satelliteInfo) {
             satelliteInfo = await WailsApp.GetSatelliteConnectionInfo();
        }
        // Wait for DOM update then render
        setTimeout(() => {
            const canvas = document.getElementById('pos-qr-canvas');
            if (canvas && satelliteInfo) {
                // Generar URL directa
                const payload = satelliteInfo.url;
                QRCode.toCanvas(canvas, payload, { width: 250, margin: 2 }, (e) => {
                   if(e) console.error(e);
                });
            }
        }, 100);
    }

    onMount(async () => {
        secuencial = await Backend.getNextSecuencial();
        if(inputElement) inputElement.focus();

        // Escuchar cambios de stock remotos
        EventsOn("inventory-updated", (updatedProd: any) => {
            // Actualizar resultados de búsqueda si están visibles
            if (isSearching && searchResults.length > 0) {
                const idx = searchResults.findIndex(p => p.SKU === updatedProd.SKU);
                if (idx !== -1) {
                    searchResults[idx].Stock = updatedProd.Stock;
                    searchResults = [...searchResults];
                }
            }
            
            // Notificar si está en el carrito (Opcional: Podríamos validar stock disponible)
            const inCart = items.find(i => i.codigo === updatedProd.SKU);
            if (inCart) {
                notifications.show(`Stock de ${updatedProd.Name} cambió a ${updatedProd.Stock}`, "info");
            }
        });

        // Escuchar escaneos remotos desde el celular
        EventsOn("pos-scan-event", (data: any) => {
            if (data.product) {
                // Convertir modelo DB a modelo POS si es necesario, o usar el que viene
                // El addItem espera un db.ProductDTO. Adaptamos lo que falte.
                const prod: db.ProductDTO = {
                    SKU: data.product.sku || data.product.SKU, // Tolerancia a mayus/minus
                    Name: data.product.name || data.product.Name,
                    Price: data.product.price || data.product.Price,
                    Stock: data.product.stock || data.product.Stock,
                    TaxCode: String(data.product.tax_code || data.product.TaxCode),
                    TaxPercentage: data.product.tax_percentage || data.product.TaxPercentage,
                    Barcode: data.product.barcode || data.product.Barcode,
                    AuxiliaryCode: data.product.auxiliary_code || data.product.AuxiliaryCode,
                    MinStock: 0,
                    ExpiryDate: "",
                    Location: ""
                };
                
                // Añadir N veces
                const qty = data.quantity || 1;
                for (let k = 0; k < qty; k++) {
                     addItem(prod);
                }
                notifications.show(`📱 Escaneo remoto: ${prod.Name}`, "success");
            }
        });
    });
</script>

<svelte:window on:keydown={handleGlobalKeyDown} />

<div class="pos-container">
    <div class="pos-main">
        <header class="pos-header">
            <div class="scanner-zone">
                <div class="icon-scan">🔍</div>
                <input 
                    bind:this={inputElement}
                    bind:value={barcodeInput}
                    on:keydown={handleBarcode}
                    placeholder="Escanee código o escriba SKU..."
                    class="barcode-input"
                />
            </div>
            <div style="display: flex; gap: 10px; align-items: center;">
                <div class="client-mini-card">
                    <span class="label">Cliente:</span>
                    <span class="value">{selectedClient.Nombre}</span>
                    <button class="btn-icon-small" title="Cambiar Cliente" on:click={() => showClientSearch = true}>👤</button>
                </div>
                <button class="btn-icon-small" title="Vincular App Móvil" on:click={openQR}>📱</button>
            </div>
        </header>

        <div class="pos-table-container">
            <table class="pos-table">
                <thead>
                    <tr>
                        <th>Producto</th>
                        <th class="text-center">Cant.</th>
                        <th class="text-right">P. Unit</th>
                        <th class="text-right">Total</th>
                        <th class="text-center">Acciones</th>
                    </tr>
                </thead>
                <tbody>
                    {#each items as item, i}
                        <tr class="item-row">
                            <td>
                                <div class="item-name">{item.nombre}</div>
                                <div class="item-sku">{item.codigo}</div>
                            </td>
                            <td class="text-center">
                                <div class="qty-control">
                                    <button on:click={() => updateQuantity(i, -1)}>-</button>
                                    <span class="qty-val">{item.cantidad}</span>
                                    <button on:click={() => updateQuantity(i, 1)}>+</button>
                                </div>
                            </td>
                            <td class="text-right">${item.precio.toFixed(2)}</td>
                            <td class="text-right font-bold">${(item.cantidad * item.precio).toFixed(2)}</td>
                            <td class="text-center">
                                <button class="btn-delete" on:click={() => removeItem(i)}>🗑️</button>
                            </td>
                        </tr>
                    {/each}
                    {#if items.length === 0}
                        <tr>
                            <td colspan="5" class="empty-state">
                                🛒 Esperando productos... (F5 para buscar manual)
                            </td>
                        </tr>
                    {/if}
                </tbody>
            </table>
        </div>
    </div>

    <aside class="pos-sidebar">
        <div class="total-display">
            <span class="total-label">TOTAL A PAGAR</span>
            <span class="total-value">${total.toFixed(2)}</span>
        </div>

        <div class="summary-card">
            <div class="summary-row">
                <span>Subtotal 0%</span>
                <span>${subtotal0.toFixed(2)}</span>
            </div>
            <div class="summary-row">
                <span>Subtotal 15%</span>
                <span>${subtotal15.toFixed(2)}</span>
            </div>
            <div class="summary-row">
                <span>IVA (15%)</span>
                <span>${iva.toFixed(2)}</span>
            </div>
        </div>

        <div class="action-buttons">
            <button class="btn-pay" on:click={processVenta}>
                <span class="btn-text">COBRAR (F12)</span>
                <span class="btn-icon">💳</span>
            </button>
            <button class="btn-cancel" on:click={() => items = []}>
                CANCELAR (ESC)
            </button>
        </div>

        <div class="shortcuts-guide">
            <div class="shortcut"><span>F5</span> Búsqueda Manual</div>
            <div class="shortcut"><span>F12</span> Cobrar</div>
            <div class="shortcut"><span>ESC</span> Limpiar</div>
            <div class="shortcut"><span>+ / -</span> Cantidades</div>
        </div>
    </aside>
</div>

{#if isSearching}
    <div class="modal-overlay" on:click|self={() => isSearching = false}>
        <div class="search-modal card">
            <h3>Búsqueda de Productos</h3>
            <input 
                bind:value={searchInput} 
                on:input={searchManual}
                placeholder="Nombre o descripción..."
                class="full-width"
                autofocus
            />
            <div class="search-results">
                {#each searchResults as p}
                    <button class="result-item" on:click={() => { addItem(p); isSearching = false; searchInput = ""; }}>
                        <div class="res-info">
                            <span class="res-name">{p.Name}</span>
                            <span class="res-sku">{p.SKU} | Stock: {p.Stock}</span>
                        </div>
                        <span class="res-price">${p.Price.toFixed(2)}</span>
                    </button>
                {/each}
            </div>
        </div>
    </div>
{/if}

{#if showClientSearch}
    <div class="modal-overlay" on:click|self={() => showClientSearch = false}>
        <div class="search-modal card">
            <div class="flex-row space-between">
                <h3>Seleccionar Cliente</h3>
                <button class="btn-secondary small" on:click={openNewClientForm}>+ Nuevo</button>
            </div>
            
            <input 
                bind:value={clientSearchTerm} 
                on:input={searchClients}
                placeholder="RUC, Cédula o Nombre..."
                class="full-width"
                autofocus
            />
            <div class="search-results">
                {#each clientSearchResults as c}
                    <button class="result-item" on:click={() => selectFoundClient(c)}>
                        <div class="res-info">
                            <span class="res-name">{c.Nombre}</span>
                            <span class="res-sku">{c.ID} | {c.Email || "Sin email"}</span>
                        </div>
                        <span class="res-price">➜</span>
                    </button>
                {/each}
                {#if clientSearchResults.length === 0 && clientSearchTerm.length > 2}
                     <div style="padding: 20px; text-align: center; color: #888;">No se encontraron clientes</div>
                {/if}
            </div>
        </div>
    </div>
{/if}

{#if showNewClientForm}
    <div class="modal-overlay" on:click|self={() => showNewClientForm = false}>
        <div class="card" style="width: 500px;">
            <h3>Registrar Nuevo Cliente</h3>
            <div class="mt-4">
                <ClientForm 
                    bind:client={newClientData} 
                    on:saved={onNewClientSaved}
                    on:cancel={() => showNewClientForm = false}
                />
            </div>
        </div>
    </div>
{/if}

{#if showQR}
    <div class="modal-overlay" on:click|self={() => showQR = false}>
        <div class="search-modal card" style="text-align: center; width: 400px;">
            <h3>Satélite Móvil</h3>
            <div style="background: white; padding: 10px; display: inline-block; border-radius: 8px; margin: 20px 0;">
                <canvas id="pos-qr-canvas"></canvas>
            </div>
            {#if satelliteInfo}
                <div style="font-size: 24px; font-weight: bold; color: var(--accent-mint); margin-bottom: 10px;">
                    PIN: {satelliteInfo.token}
                </div>
                <div style="font-size: 12px; color: #888;">
                    Escanee con la cámara de su celular para conectar
                </div>
            {:else}
                <p>Cargando...</p>
            {/if}
            <button class="btn-secondary full-width mt-4" on:click={() => showQR = false}>Cerrar</button>
        </div>
    </div>
{/if}

<style>
    .pos-container {
        display: flex;
        gap: 20px;
        height: 100%;
        color: var(--text-primary);
    }

    .pos-main {
        flex: 1;
        display: flex;
        flex-direction: column;
        background: var(--bg-panel);
        border: 1px solid var(--border-subtle);
        border-radius: 12px;
        overflow: hidden;
    }

    .pos-header {
        padding: 20px;
        background: rgba(0,0,0,0.2);
        display: flex;
        justify-content: space-between;
        align-items: center;
        border-bottom: 1px solid var(--border-subtle);
    }

    .scanner-zone {
        display: flex;
        align-items: center;
        gap: 12px;
        background: var(--bg-app);
        padding: 8px 16px;
        border-radius: 8px;
        border: 1px solid var(--accent-mint);
        box-shadow: 0 0 15px rgba(52, 211, 153, 0.1);
        width: 400px;
    }

    .barcode-input {
        background: transparent;
        border: none;
        color: var(--text-primary);
        font-size: 18px;
        width: 100%;
        outline: none;
    }

    .client-mini-card {
        display: flex;
        align-items: center;
        gap: 10px;
        background: rgba(255,255,255,0.05);
        padding: 8px 16px;
        border-radius: 8px;
    }

    .label { color: var(--text-tertiary); font-size: 12px; }
    .value { font-weight: 600; color: var(--accent-mint); }

    .pos-table-container {
        flex: 1;
        overflow-y: auto;
    }

    .pos-table {
        width: 100%;
        border-collapse: collapse;
    }

    .pos-table th {
        background: rgba(0,0,0,0.1);
        padding: 12px 20px;
        text-align: left;
        color: var(--text-tertiary);
        font-size: 13px;
        text-transform: uppercase;
    }

    .item-row {
        border-bottom: 1px solid rgba(255,255,255,0.03);
        transition: background 0.2s;
    }

    .item-row:hover { background: rgba(255,255,255,0.02); }

    .pos-table td { padding: 16px 20px; }

    .item-name { font-weight: 600; font-size: 16px; }
    .item-sku { font-size: 12px; color: var(--text-tertiary); }

    .qty-control {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 12px;
    }

    .qty-control button {
        width: 28px;
        height: 28px;
        border-radius: 6px;
        border: 1px solid var(--border-subtle);
        background: var(--bg-app);
        color: white;
        cursor: pointer;
    }

    .qty-val { font-size: 18px; font-weight: bold; min-width: 30px; }

    .btn-delete {
        background: transparent;
        border: none;
        cursor: pointer;
        font-size: 18px;
        opacity: 0.5;
    }
    .btn-delete:hover { opacity: 1; }

    .empty-state {
        text-align: center;
        padding: 100px;
        color: var(--text-tertiary);
        font-size: 18px;
    }

    /* SIDEBAR POS */
    .pos-sidebar {
        width: 350px;
        display: flex;
        flex-direction: column;
        gap: 20px;
    }

    .total-display {
        background: var(--accent-mint);
        color: #000;
        padding: 30px;
        border-radius: 12px;
        display: flex;
        flex-direction: column;
        align-items: center;
        box-shadow: 0 10px 30px rgba(52, 211, 153, 0.2);
    }

    .total-label { font-size: 14px; font-weight: 800; opacity: 0.7; }
    .total-value { font-size: 52px; font-weight: 900; }

    .summary-card {
        background: var(--bg-panel);
        border: 1px solid var(--border-subtle);
        border-radius: 12px;
        padding: 20px;
        display: flex;
        flex-direction: column;
        gap: 10px;
    }

    .summary-row {
        display: flex;
        justify-content: space-between;
        color: var(--text-secondary);
    }

    .action-buttons {
        display: flex;
        flex-direction: column;
        gap: 12px;
    }

    .btn-pay {
        background: var(--accent-mint);
        color: #000;
        border: none;
        padding: 20px;
        border-radius: 12px;
        font-size: 20px;
        font-weight: 800;
        cursor: pointer;
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 15px;
        transition: transform 0.2s, box-shadow 0.2s;
    }

    .btn-pay:hover {
        transform: translateY(-2px);
        box-shadow: 0 5px 20px rgba(52, 211, 153, 0.3);
    }

    .btn-cancel {
        background: rgba(255, 100, 100, 0.1);
        color: #ff6464;
        border: 1px solid #ff6464;
        padding: 12px;
        border-radius: 12px;
        cursor: pointer;
        font-weight: 600;
    }

    .shortcuts-guide {
        margin-top: auto;
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 8px;
    }

    .shortcut {
        font-size: 11px;
        color: var(--text-tertiary);
    }
    .shortcut span {
        background: var(--bg-panel);
        padding: 2px 6px;
        border-radius: 4px;
        border: 1px solid var(--border-subtle);
        margin-right: 4px;
    }

    /* MODAL */
    .modal-overlay {
        position: fixed;
        inset: 0;
        background: rgba(0,0,0,0.8);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
    }

    .search-modal {
        width: 600px;
        padding: 30px;
        background: var(--bg-panel);
    }

    .search-results {
        margin-top: 20px;
        max-height: 400px;
        overflow-y: auto;
        display: flex;
        flex-direction: column;
        gap: 8px;
    }

    .result-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 12px 16px;
        background: rgba(255,255,255,0.03);
        border: 1px solid var(--border-subtle);
        border-radius: 8px;
        color: white;
        cursor: pointer;
        text-align: left;
    }

    .result-item:hover {
        background: var(--accent-mint);
        color: #000;
    }

    .res-info { display: flex; flex-direction: column; }
    .res-name { font-weight: 600; }
    .res-sku { font-size: 11px; opacity: 0.7; }
    .res-price { font-weight: bold; font-size: 18px; }
</style>
