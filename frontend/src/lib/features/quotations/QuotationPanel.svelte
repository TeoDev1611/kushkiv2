<script>
    import {
        GetNextQuotationSecuencial,
        CreateQuotation,
        GetQuotations,
        OpenQuotationPDF,
        ConvertQuotationToInvoice,
        GetProducts,
        GetClients,
    } from "wailsjs/go/main/App.js";
    import { onMount, createEventDispatcher } from "svelte";
    import { fade } from "svelte/transition";

    export let clients = [];
    export let products = [];

    const dispatch = createEventDispatcher();

    let mode = "list"; // 'list' | 'create'
    let quotations = [];
    let totalQuotations = 0;
    let page = 1;
    let pageSize = 10;
    let loading = false;

    // Formulario Creation
    let newQuote = {
        secuencial: "",
        clienteID: "",
        clienteNombre: "",
        clienteDireccion: "",
        clienteEmail: "",
        clienteTelefono: "",
        observacion: "",
        items: [],
    };

    let clientSearch = "";
    let productSearch = "";
    let newItem = {
        codigo: "",
        nombre: "",
        cantidad: 1,
        precio: 0,
        codigoIVA: "4",
        porcentajeIVA: 15,
    };

    // Handler para evento global de guardado (Ctrl+S)
    const handleGlobalSave = () => {
        if (mode === 'create') {
            saveQuotation();
        }
    };

    import { onDestroy } from 'svelte';

    onMount(async () => {
        window.addEventListener('app-save', handleGlobalSave);
        loadQuotations();
    });

    onDestroy(() => {
        window.removeEventListener('app-save', handleGlobalSave);
    });

    async function loadQuotations() {
        loading = true;
        try {
            const res = await GetQuotations(page, pageSize);
            quotations = res.data || [];
            totalQuotations = res.total;
        } catch (e) {
            console.error(e);
        } finally {
            loading = false;
        }
    }

    async function initCreate() {
        try {
            const seq = await GetNextQuotationSecuencial();
            newQuote = {
                secuencial: seq,
                clienteID: "",
                clienteNombre: "",
                clienteDireccion: "",
                clienteEmail: "",
                clienteTelefono: "",
                observacion: "",
                items: [],
            };
            mode = "create";
        } catch (e) {
            console.error(e);
        }
    }

    function selectClient(c) {
        newQuote.clienteID = c.ID;
        newQuote.clienteNombre = c.Nombre;
        newQuote.clienteDireccion = c.Direccion;
        newQuote.clienteEmail = c.Email || "";
        newQuote.clienteTelefono = c.Telefono || "";
        clientSearch = "";
    }

    function selectProduct(p) {
        const taxCode = p.TaxCode || (p.TaxPercentage > 0 ? "4" : "0");
        const taxPerc =
            p.TaxPercentage !== undefined
                ? p.TaxPercentage
                : p.TaxCode == "4"
                  ? 15
                  : 0;
        newItem = {
            codigo: p.SKU,
            nombre: p.Name,
            cantidad: 1,
            precio: p.Price,
            codigoIVA: taxCode,
            porcentajeIVA: taxPerc,
        };
        productSearch = "";
    }

    function addItem() {
        if (!newItem.nombre || newItem.precio <= 0) return;
        newQuote.items = [...newQuote.items, { ...newItem }];
        newItem = {
            codigo: "",
            nombre: "",
            cantidad: 1,
            precio: 0,
            codigoIVA: "4",
            porcentajeIVA: 15,
        };
    }

    function removeItem(index) {
        newQuote.items = newQuote.items.filter((_, i) => i !== index);
    }

    async function saveQuotation() {
        if (!newQuote.clienteID || newQuote.items.length === 0) {
            // Simple validation feedback
            alert("Complete cliente y productos");
            return;
        }
        loading = true;
        try {
            const dto = {
                secuencial: newQuote.secuencial,
                clienteID: newQuote.clienteID,
                clienteNombre: newQuote.clienteNombre,
                clienteDireccion: newQuote.clienteDireccion,
                clienteEmail: newQuote.clienteEmail,
                clienteTelefono: newQuote.clienteTelefono,
                observacion: newQuote.observacion,
                items: newQuote.items.map((i) => ({
                    codigo: i.codigo,
                    nombre: i.nombre,
                    cantidad: parseFloat(i.cantidad),
                    precio: parseFloat(i.precio),
                    codigoIVA: i.codigoIVA.toString(),
                    porcentajeIVA: parseFloat(i.porcentajeIVA),
                })),
            };
            const res = await CreateQuotation(dto);
            if (res.includes("Error")) {
                alert(res);
            } else {
                mode = "list";
                loadQuotations();
                dispatch("toast", {
                    message: "Cotización creada",
                    type: "success",
                });
            }
        } catch (e) {
            alert(e);
        } finally {
            loading = false;
        }
    }

    async function openPDF(id) {
        const res = await OpenQuotationPDF(id);
        if (res.includes("Error")) alert(res);
    }

    async function convertToInvoice(id) {
        loading = true;
        try {
            const dto = await ConvertQuotationToInvoice(id);
            if (dto) {
                // Dispatch event to parent to switch tab and populate invoice
                dispatch("convertToInvoice", dto);
            }
        } catch (e) {
            console.error(e);
        } finally {
            loading = false;
        }
    }

    // Computed
    $: subtotal = newQuote.items.reduce(
        (sum, item) => sum + item.cantidad * item.precio,
        0,
    );
    $: total = newQuote.items.reduce(
        (sum, item) =>
            sum + item.cantidad * item.precio * (1 + item.porcentajeIVA / 100),
        0,
    );
</script>

<div class="panel" in:fade={{ duration: 100 }}>
    {#if mode === "list"}
        <div class="header-row">
            <h1>Cotizaciones</h1>
            <div class="header-actions">
                <button class="btn-secondary" on:click={loadQuotations}
                    >🔄 Refrescar</button
                >
                <button class="btn-primary" on:click={initCreate}
                    >➕ Nueva Cotización</button
                >
            </div>
        </div>

        <div class="card no-padding">
            <div class="linear-grid">
                <div class="linear-header grid-columns-quote">
                    <div class="cell">Secuencial</div>
                    <div class="cell">Fecha</div>
                    <div class="cell">Cliente</div>
                    <div class="cell text-right">Total</div>
                    <div class="cell text-center">Acciones</div>
                </div>
                <div class="rows-container">
                    {#each quotations as q}
                        <div class="linear-row grid-columns-quote">
                            <div class="cell mono text-secondary">
                                {q.secuencial}
                            </div>
                            <div class="cell mono">{q.fechaEmision}</div>
                            <div class="cell">{q.clienteNombre}</div>
                            <div class="cell text-right font-medium">
                                ${q.total.toFixed(2)}
                            </div>
                            <div
                                class="cell text-center"
                                style="gap: 4px; justify-content: center;"
                            >
                                <button
                                    class="btn-icon-mini"
                                    title="PDF"
                                    on:click={() => openPDF(q.id)}>📄</button
                                >
                                <button
                                    class="btn-icon-mini action"
                                    title="Facturar"
                                    on:click={() => convertToInvoice(q.id)}
                                    >🚀</button
                                >
                            </div>
                        </div>
                    {/each}
                </div>
                {#if quotations.length === 0}
                    <div class="empty-state">
                        No hay cotizaciones registradas
                    </div>
                {/if}
            </div>

            <!-- Simple Pagination -->
            <div class="flex-row space-between" style="padding: 16px;">
                <button
                    class="btn-secondary icon-only"
                    disabled={page === 1}
                    on:click={() => {
                        page--;
                        loadQuotations();
                    }}>«</button
                >
                <span class="text-secondary text-sm">Página {page}</span>
                <button
                    class="btn-secondary icon-only"
                    disabled={quotations.length < pageSize}
                    on:click={() => {
                        page++;
                        loadQuotations();
                    }}>»</button
                >
            </div>
        </div>
    {:else}
        <!-- CREATE MODE -->
        <div class="header-row">
            <h1>
                Nueva Cotización <span
                    class="mono text-mint"
                    style="font-size: 18px; margin-left: 8px;"
                    >#{newQuote.secuencial}</span
                >
            </h1>
            <button class="btn-secondary" on:click={() => (mode = "list")}
                >Cancelar</button
            >
        </div>

        <div class="split-view">
            <div class="left-panel">
                <div class="card">
                    <h3>Cliente</h3>
                    <div class="search-container mb-4">
                        <input
                            aria-label="Buscar cliente"
                            bind:value={clientSearch}
                            placeholder="Buscar cliente..."
                            class="full-search-input"
                        />
                        {#if clientSearch.length > 0}
                            <div class="dropdown">
                                {#each clients.filter( (c) => c.Nombre.toLowerCase().includes(clientSearch.toLowerCase()), ) as c}
                                    <button on:click={() => selectClient(c)}
                                        >{c.Nombre}</button
                                    >
                                {/each}
                            </div>
                        {/if}
                    </div>
                    <div class="form-stack">
                        <input
                            aria-label="Nombre cliente"
                            bind:value={newQuote.clienteNombre}
                            placeholder="Nombre Cliente"
                        />
                        <input
                            aria-label="RUC/CI cliente"
                            bind:value={newQuote.clienteID}
                            placeholder="RUC/CI"
                        />
                        <input
                            aria-label="Email cliente"
                            bind:value={newQuote.clienteEmail}
                            placeholder="Email"
                        />
                        <input
                            aria-label="Teléfono cliente"
                            bind:value={newQuote.clienteTelefono}
                            placeholder="Teléfono"
                        />
                        <input
                            aria-label="Dirección cliente"
                            bind:value={newQuote.clienteDireccion}
                            placeholder="Dirección"
                        />
                    </div>
                </div>

                <div class="card mt-4">
                    <h3>Totales</h3>
                    <div class="flex-row space-between text-secondary mb-2">
                        <span>Subtotal:</span>
                        <span>${subtotal.toFixed(2)}</span>
                    </div>
                    <div class="flex-row space-between grand-total">
                        <span>TOTAL:</span> <span>${total.toFixed(2)}</span>
                    </div>
                    <button
                        class="btn-primary full-width large mt-4"
                        on:click={saveQuotation}
                        disabled={loading}
                    >
                        {loading ? "Guardando..." : "💾 Guardar Cotización"}
                    </button>
                </div>
            </div>

            <div class="right-panel card">
                <h3>Productos</h3>
                <div class="add-product-bar">
                    <div class="search-box-product">
                        <input
                            aria-label="Buscar producto"
                            bind:value={productSearch}
                            placeholder="Buscar producto..."
                        />
                        {#if productSearch.length > 0}
                            <div class="dropdown">
                                {#each products.filter( (p) => p.Name.toLowerCase().includes(productSearch.toLowerCase()), ) as p}
                                    <button on:click={() => selectProduct(p)}
                                        >{p.Name} - ${p.Price}</button
                                    >
                                {/each}
                            </div>
                        {/if}
                    </div>
                    <input
                        aria-label="Cantidad"
                        class="input-qty"
                        type="number"
                        bind:value={newItem.cantidad}
                        min="1"
                        placeholder="#"
                    />
                    <input
                        aria-label="Precio"
                        class="input-price"
                        type="number"
                        bind:value={newItem.precio}
                        placeholder="$"
                    />
                    <button
                        aria-label="Añadir ítem"
                        class="btn-add"
                        on:click={addItem}>+</button
                    >
                </div>

                <div class="linear-grid mt-4">
                    <div class="linear-header grid-columns-items">
                        <div class="cell">Cant</div>
                        <div class="cell">Desc</div>
                        <div class="cell">P.Unit</div>
                        <div class="cell">Total</div>
                        <div class="cell"></div>
                    </div>
                    <div class="rows-container">
                        {#each newQuote.items as item, i}
                            <div class="linear-row grid-columns-items">
                                <div class="cell">{item.cantidad}</div>
                                <div class="cell">{item.nombre}</div>
                                <div class="cell">
                                    ${item.precio.toFixed(2)}
                                </div>
                                <div class="cell font-medium">
                                    ${(item.cantidad * item.precio).toFixed(2)}
                                </div>
                                <div class="cell">
                                    <button
                                        aria-label="Eliminar ítem"
                                        class="btn-icon-mini danger"
                                        on:click={() => removeItem(i)}>×</button
                                    >
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>

                <textarea
                    aria-label="Observaciones"
                    class="mt-4 full-width"
                    placeholder="Observaciones y notas adicionales..."
                    bind:value={newQuote.observacion}
                    style="height: 100px; resize: none;"
                ></textarea>
            </div>
        </div>
    {/if}
</div>

<style>
    /* Local Overrides matching Global System */
    .grid-columns-quote {
        grid-template-columns: 120px 140px 1fr 120px 100px;
    }
    .grid-columns-items {
        grid-template-columns: 60px 1fr 80px 80px 40px;
    }

    .panel {
        display: flex;
        flex-direction: column;
        gap: 24px;
        height: 100%;
    }

    .text-mint {
        color: var(--accent-mint);
    }
    .font-medium {
        font-weight: 500;
    }
    .text-secondary {
        color: var(--text-secondary);
    }

    /* Layouts */
    .split-view {
        display: grid;
        grid-template-columns: 350px 1fr;
        gap: 24px;
    }
    .left-panel {
        display: flex;
        flex-direction: column;
        gap: 24px;
    }
    .form-stack {
        display: flex;
        flex-direction: column;
        gap: 12px;
    }

    /* Inputs specific overrides */
    .full-search-input {
        width: 100%;
    }
    .search-container {
        position: relative;
    }
    .dropdown {
        position: absolute;
        top: 100%;
        left: 0;
        right: 0;
        background: var(--bg-panel);
        border: 1px solid var(--accent-mint);
        z-index: 100;
        max-height: 240px;
        overflow-y: auto;
        border-radius: 0 0 6px 6px;
        box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
    }
    .dropdown button {
        width: 100%;
        padding: 12px;
        background: transparent;
        border: none;
        color: var(--text-primary);
        text-align: left;
        cursor: pointer;
        border-bottom: 1px solid var(--border-subtle);
    }
    .dropdown button:hover {
        background: rgba(0, 255, 148, 0.1);
        color: var(--accent-mint);
    }

    /* Product Bar */
    .add-product-bar {
        display: flex;
        gap: 12px;
        background: var(--bg-hover);
        padding: 16px;
        border-radius: var(--radius-sm);
        border: 1px solid var(--border-subtle);
        align-items: center;
    }
    .search-box-product {
        flex: 1;
        position: relative;
    }
    .input-qty {
        width: 70px;
        text-align: center;
    }
    .input-price {
        width: 90px;
        text-align: right;
    }
    .btn-add {
        background: var(--accent-mint);
        color: black;
        border: none;
        width: 36px;
        height: 36px;
        border-radius: 6px;
        font-weight: bold;
        cursor: pointer;
        font-size: 18px;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.2s;
    }
    .btn-add:hover {
        transform: scale(1.05);
        box-shadow: 0 0 10px rgba(0, 255, 148, 0.4);
    }

    /* Totals */
    .grand-total {
        font-size: 18px;
        color: var(--text-primary);
        font-weight: 700;
        border-top: 1px solid var(--border-subtle);
        padding-top: 16px;
        margin-top: 16px;
    }

    .full-width {
        width: 100%;
    }
    .large {
        height: 48px;
        font-size: 14px;
        letter-spacing: 0.05em;
        text-transform: uppercase;
    }
    .empty-state {
        padding: 40px;
        text-align: center;
        color: var(--text-secondary);
        font-style: italic;
    }
    .btn-icon-mini.danger:hover {
        color: var(--status-error);
        background: rgba(255, 51, 51, 0.1);
    }
</style>
