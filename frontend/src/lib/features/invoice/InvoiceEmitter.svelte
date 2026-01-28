<script lang="ts">
    import { onMount } from 'svelte';
    import { fade } from 'svelte/transition';
    import { flip } from 'svelte/animate';
    import { invoiceStore, invoiceTotals } from '$lib/stores/invoice';
    import { notifications } from '$lib/stores/notifications';
    import { Backend } from '$lib/services/api';
    import { withLoading } from '$lib/stores/app';

    // Estado local para búsquedas y UI
    let clientSearch = "";
    let productSearch = "";
    let clients: any[] = [];
    let products: any[] = [];
    let showClientDropdown = false;
    let showProductDropdown = false;
    
    // Item temporal para agregar
    let newItem = {
        codigo: "",
        nombre: "",
        cantidad: 1,
        precio: 0,
        codigoIVA: "4",
        porcentajeIVA: 15
    };

    // Errores de validación visual
    let errors: Record<string, boolean> = {};

    onMount(async () => {
        try {
            const [prods, cli, seq] = await Promise.all([
                Backend.getProducts(),
                Backend.getClients(),
                Backend.getNextSecuencial()
            ]);
            products = prods || [];
            clients = cli || [];
            
            // Actualizar secuencial si la tienda está "fresca" (vacía de items)
            if ($invoiceStore.items.length === 0) {
                invoiceStore.updateSecuencial(seq || "000000001");
            }
        } catch (e) {
            console.error(e);
        }
    });

    // --- Manejo de Clientes ---
    function selectClient(c: any) {
        invoiceStore.setClient(c);
        clientSearch = "";
        showClientDropdown = false;
    }

    // --- Manejo de Productos ---
    function selectProduct(p: any) {
        const taxCode = p.TaxCode || (p.TaxPercentage > 0 ? "4" : "0");
        const taxPerc = p.TaxPercentage !== undefined ? p.TaxPercentage : (p.TaxCode == "4" ? 15 : 0);
        
        newItem = {
            codigo: p.SKU,
            nombre: p.Name,
            cantidad: 1,
            precio: p.Price,
            codigoIVA: taxCode,
            porcentajeIVA: taxPerc
        };
        productSearch = "";
        showProductDropdown = false;
    }

    function addItem() {
        if (!newItem.nombre || newItem.precio <= 0) return;
        
        // Convertir tipos por seguridad
        const itemToAdd = {
            ...newItem,
            cantidad: parseFloat(String(newItem.cantidad)),
            precio: parseFloat(String(newItem.precio)),
            porcentajeIVA: parseFloat(String(newItem.porcentajeIVA))
        };

        invoiceStore.addItem(itemToAdd);
        
        // Reset item
        newItem = {
            codigo: "",
            nombre: "",
            cantidad: 1,
            precio: 0,
            codigoIVA: "4",
            porcentajeIVA: 15
        };
    }

    function updateItemTax() {
        const map: Record<string, number> = { "0": 0, "2": 12, "4": 15, "5": 5 };
        newItem.porcentajeIVA = map[newItem.codigoIVA] || 0;
    }

    // --- Emisión ---
    async function handleEmit() {
        errors = {};
        let isValid = true;
        const inv = $invoiceStore;

        if (!inv.clienteID) { errors.clienteID = true; isValid = false; }
        if (!inv.clienteNombre) { errors.clienteNombre = true; isValid = false; }
        if (!inv.clienteDireccion) { errors.clienteDireccion = true; isValid = false; }
        if (!inv.clienteEmail || !inv.clienteEmail.includes("@")) { errors.clienteEmail = true; isValid = false; }
        
        if (inv.items.length === 0) {
            notifications.show("Agregue al menos un producto", "error");
            return;
        }

        if (!isValid) {
            notifications.show("Complete los campos obligatorios marcados en rojo", "error");
            return;
        }

        try {
            const res = await withLoading(Backend.createInvoice(inv));
            if (res.startsWith("Éxito")) {
                notifications.show(res, "success");
                
                // Preparar para la siguiente
                const nextSeq = (parseInt(inv.secuencial) + 1).toString().padStart(9, "0");
                invoiceStore.reset(nextSeq);
            } else {
                notifications.show(res, "error");
            }
        } catch (e) {
            notifications.show("Error crítico: " + e, "error");
        }
    }
</script>

<div class="panel" in:fade={{ duration: 200 }}>
    <div class="header-row">
        <h1>Emitir Factura</h1>
        <!-- Aquí se podría poner el switcher de ambiente si fuera necesario -->
    </div>

    <!-- SECCIÓN CLIENTE -->
    <div class="card mb-4">
        <div class="flex-row space-between mb-2">
            <h3>👤 Datos del Cliente</h3>
            <div class="badge mono" style="font-size: 14px; background: var(--mint-dim); color: var(--accent-mint); border: 1px solid var(--mint-border);">
                SEC: {$invoiceStore.secuencial}
            </div>
        </div>

        <!-- Buscador de Cliente -->
        <div class="search-container mb-3" style="position: relative;">
             <input
                bind:value={clientSearch}
                on:focus={() => showClientDropdown = true}
                on:blur={() => setTimeout(() => showClientDropdown = false, 200)}
                placeholder="🔍 Buscar cliente registrado..."
                style="background: var(--bg-surface);"
            />
            {#if showClientDropdown && clientSearch.length > 0}
                <div class="dropdown-menu">
                    {#each clients.filter(c => c.Nombre.toLowerCase().includes(clientSearch.toLowerCase()) || c.ID.includes(clientSearch)) as c}
                        <button class="dropdown-item" on:click={() => selectClient(c)}>
                            {c.Nombre} <span class="text-secondary text-small">({c.ID})</span>
                        </button>
                    {/each}
                </div>
            {/if}
        </div>

        <div class="invoice-client-grid">
            <div class="field">
                <label for="inv-id">Identificación</label>
                <input id="inv-id" bind:value={$invoiceStore.clienteID} class:invalid={errors.clienteID} placeholder="RUC / Cédula" />
            </div>
            <div class="field">
                <label for="inv-name">Razón Social</label>
                <input id="inv-name" bind:value={$invoiceStore.clienteNombre} class:invalid={errors.clienteNombre} />
            </div>
            <div class="field">
                <label for="inv-email">Email</label>
                <input id="inv-email" type="email" bind:value={$invoiceStore.clienteEmail} class:invalid={errors.clienteEmail} />
            </div>
            <div class="field">
                <label for="inv-tel">Teléfono</label>
                <input id="inv-tel" bind:value={$invoiceStore.clienteTelefono} />
            </div>
            <div class="field span-2">
                <label for="inv-addr">Dirección</label>
                <input id="inv-addr" bind:value={$invoiceStore.clienteDireccion} class:invalid={errors.clienteDireccion} />
            </div>
        </div>
    </div>

    <!-- SECCIÓN ITEMS -->
    <div class="card flex-col flex-1" style="min-height: 400px;">
        <h3 class="mb-2">📦 Detalle de Factura</h3>

        <!-- Barra de agregar item -->
        <div class="invoice-product-bar mb-3">
             <div style="position: relative; flex: 2;">
                <input
                    bind:value={productSearch}
                    on:focus={() => showProductDropdown = true}
                    on:blur={() => setTimeout(() => showProductDropdown = false, 200)}
                    placeholder="🔍 Buscar producto..."
                />
                {#if showProductDropdown && productSearch.length > 0}
                    <div class="dropdown-menu">
                        {#each products.filter(p => p.Name.toLowerCase().includes(productSearch.toLowerCase())) as p}
                            <button class="dropdown-item" on:click={() => selectProduct(p)}>
                                {p.Name} - <span class="text-mint">${p.Price}</span>
                            </button>
                        {/each}
                    </div>
                {/if}
            </div>
            
            <input bind:value={newItem.codigo} placeholder="Cód." style="flex: 0.5;" />
            <input bind:value={newItem.nombre} placeholder="Descripción" style="flex: 2;" />
            <input type="number" bind:value={newItem.cantidad} placeholder="Cant" style="flex: 0.5;" class="text-center" />
            <input type="number" step="0.01" bind:value={newItem.precio} placeholder="$$" style="flex: 0.8;" class="text-right" />
            
            <select bind:value={newItem.codigoIVA} on:change={updateItemTax} style="flex: 0.8;">
                <option value="4">15%</option>
                <option value="5">5%</option>
                <option value="0">0%</option>
                <option value="2">12%</option>
            </select>

            <button class="btn-primary icon-only" on:click={addItem} title="Agregar">+</button>
        </div>

        <!-- Tabla de Items -->
        <div class="linear-grid flex-1 overflow-hidden flex-col">
            <div class="linear-header grid-columns-invoice">
                <div class="cell text-center">Cant</div>
                <div class="cell">Descripción</div>
                <div class="cell text-right">P. Unit</div>
                <div class="cell text-center">IVA</div>
                <div class="cell text-right">Total</div>
                <div class="cell"></div>
            </div>
            <div class="rows-container overflow-auto flex-1">
                {#each $invoiceStore.items as item, i (item)}
                    <div class="linear-row grid-columns-invoice" animate:flip={{ duration: 300 }}>
                        <div class="cell text-center">{item.cantidad}</div>
                        <div class="cell">{item.nombre}</div>
                        <div class="cell text-right">${item.precio.toFixed(2)}</div>
                        <div class="cell text-center"><span class="badge">{item.porcentajeIVA}%</span></div>
                        <div class="cell text-right font-medium">
                            ${(item.cantidad * item.precio * (1 + item.porcentajeIVA / 100)).toFixed(2)}
                        </div>
                        <div class="cell text-center">
                            <button class="btn-icon-mini danger" on:click={() => invoiceStore.removeItem(i)}>×</button>
                        </div>
                    </div>
                {/each}
                {#if $invoiceStore.items.length === 0}
                    <div class="empty-state">
                        <p class="empty-state-text text-small">Agregue productos para comenzar.</p>
                    </div>
                {/if}
            </div>
        </div>

        <!-- Footer Totales -->
        <div class="invoice-footer mt-3 pt-3 border-top flex-row" style="align-items: flex-start; gap: 24px;">
            <div class="notes-area flex-1">
                <label for="inv-obs">Observaciones / Notas</label>
                <textarea id="inv-obs" bind:value={$invoiceStore.observacion} rows="3" class="full-width"></textarea>
            </div>
            
            <div class="totals-area card compact bg-surface" style="width: 300px;">
                <div class="flex-row space-between mb-1">
                    <span class="text-secondary">Subtotal</span>
                    <span>${$invoiceTotals.subtotal.toFixed(2)}</span>
                </div>
                <div class="flex-row space-between mb-2">
                    <span class="text-secondary">IVA Total</span>
                    <span>${$invoiceTotals.iva.toFixed(2)}</span>
                </div>
                <div class="flex-row space-between pt-2 border-top">
                    <span class="font-bold text-lg">TOTAL</span>
                    <span class="font-bold text-lg text-mint">${$invoiceTotals.total.toFixed(2)}</span>
                </div>
                
                <button class="btn-primary full-width mt-3" on:click={handleEmit}>
                    🖋️ Firmar y Emitir
                </button>
            </div>
        </div>
    </div>
</div>

<style>
    .invoice-client-grid {
        display: grid;
        grid-template-columns: 1fr 2fr 1.5fr 1fr;
        gap: 16px;
    }
    .span-2 { grid-column: span 2; }

    .invoice-product-bar {
        display: flex;
        gap: 8px;
        align-items: center;
        background: var(--bg-hover);
        padding: 12px;
        border-radius: 8px;
    }

    .grid-columns-invoice {
        display: grid;
        grid-template-columns: 60px 1fr 100px 60px 100px 50px;
        padding: 0 12px;
        align-items: center;
    }

    .dropdown-menu {
        position: absolute;
        top: 100%; left: 0; right: 0;
        background: var(--bg-panel);
        border: 1px solid var(--accent-mint);
        z-index: 50;
        max-height: 200px;
        overflow-y: auto;
        box-shadow: 0 8px 16px rgba(0,0,0,0.5);
    }
    .dropdown-item {
        width: 100%; padding: 8px 12px;
        text-align: left; background: none; border: none;
        color: var(--text-primary); cursor: pointer;
        border-bottom: 1px solid var(--border-subtle);
    }
    .dropdown-item:hover { background: rgba(0,255,148,0.1); }
    
    .invalid { border-color: var(--status-error) !important; }
    
    .flex-col { display: flex; flex-direction: column; }
    .flex-1 { flex: 1; }
    .overflow-hidden { overflow: hidden; }
    .overflow-auto { overflow-y: auto; }
    .border-top { border-top: 1px solid var(--border-subtle); }
    .text-small { font-size: 0.85rem; }
    
    @media (max-width: 900px) {
        .invoice-client-grid { grid-template-columns: 1fr 1fr; }
        .invoice-product-bar { flex-wrap: wrap; }
    }
</style>
