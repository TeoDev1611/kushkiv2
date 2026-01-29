<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { fade } from 'svelte/transition';
    import { Backend } from '$lib/services/api';
    import { notifications } from '$lib/stores/notifications';
    import { withLoading } from '$lib/stores/app';
    import type { db } from 'wailsjs/go/models';
    import * as WailsApp from 'wailsjs/go/main/App';
    import { EventsOn } from '../../../../wailsjs/runtime/runtime';

    // Estado local
    let products: db.ProductDTO[] = [];
    let searchTerm = "";
    let sortCol: keyof db.ProductDTO = "Name";
    let sortAsc = true;
    let isEditing = false;
    let editingProduct: any = {
        SKU: "",
        Name: "",
        Price: 0,
        Stock: 0,
        TaxCode: "4", // 15% IVA
        TaxPercentage: 15
    };

    // Handler para evento global de guardado (Ctrl+S)
    const handleGlobalSave = () => {
        if (editingProduct.SKU && editingProduct.Name) {
            handleSave();
        }
    };

    onMount(() => {
        window.addEventListener('app-save', handleGlobalSave);
        loadProducts();

        // Real-time Sync Listener
        EventsOn("inventory-updated", (updatedProd: any) => {
            const index = products.findIndex(p => p.SKU === updatedProd.SKU);
            if (index !== -1) {
                products[index].Stock = updatedProd.Stock;
                // Forzar reactividad
                products = [...products]; 
                
                // Si estamos editando este producto, actualizar también el form
                if (isEditing && editingProduct.SKU === updatedProd.SKU) {
                    editingProduct.Stock = updatedProd.Stock;
                }
                
                notifications.show(`Stock actualizado: ${updatedProd.Name}`, "info");
            }
        });
    });

    onDestroy(() => {
        window.removeEventListener('app-save', handleGlobalSave);
    });

    async function loadProducts() {
        try {
            if (searchTerm.length > 2) {
                products = await withLoading(WailsApp.SearchProducts(searchTerm)) || [];
            } else {
                products = await withLoading(Backend.getProducts()) || [];
            }
        } catch (e) {
            notifications.show("Error cargando inventario: " + e, "error");
        }
    }

    let searchTimeout: any;
    function handleSearch() {
        clearTimeout(searchTimeout);
        searchTimeout = setTimeout(() => {
            loadProducts();
        }, 300);
    }

    function sort(col: keyof db.ProductDTO) {
        if (sortCol === col) {
            sortAsc = !sortAsc;
        } else {
            sortCol = col;
            sortAsc = true;
        }
    }

    $: sortedProducts = [...products].sort((a: any, b: any) => {
        let valA = a[sortCol];
        let valB = b[sortCol];
        
        // Manejo numérico
        if (typeof valA === 'number' && typeof valB === 'number') {
            return sortAsc ? valA - valB : valB - valA;
        }
        
        // Manejo string
        valA = String(valA).toLowerCase();
        valB = String(valB).toLowerCase();
        if (valA < valB) return sortAsc ? -1 : 1;
        if (valA > valB) return sortAsc ? 1 : -1;
        return 0;
    });

    function resetForm() {
        editingProduct = {
            SKU: "",
            Name: "",
            Price: 0,
            Stock: 0,
            TaxCode: "4",
            TaxPercentage: 15
        };
        isEditing = false;
    }

    function selectProduct(p: db.ProductDTO) {
        editingProduct = { ...p };
        isEditing = true;
    }

    async function handleSave() {
        if (!editingProduct.SKU || !editingProduct.Name) {
            notifications.show("Código y Nombre son obligatorios", "warning");
            return;
        }

        try {
            // Asegurar tipos numéricos
            editingProduct.Price = parseFloat(String(editingProduct.Price));
            editingProduct.Stock = parseInt(String(editingProduct.Stock));
            
            const res = await withLoading(Backend.saveProduct(editingProduct));
            
            if (res.startsWith("Error")) {
                notifications.show(res, "error");
            } else {
                notifications.show("Producto guardado", "success");
                await loadProducts();
                if (!isEditing) resetForm();
            }
        } catch (e) {
            notifications.show("Error guardando: " + e, "error");
        }
    }

    async function handleDelete(sku: string) {
        if (!confirm("¿Eliminar este producto permanentemente?")) return;

        try {
            const res = await withLoading(Backend.deleteProduct(sku));
            if (res.startsWith("Error")) {
                notifications.show(res, "error");
            } else {
                notifications.show("Producto eliminado", "success");
                if (editingProduct.SKU === sku) resetForm();
                await loadProducts();
            }
        } catch (e) {
            notifications.show("Error eliminando: " + e, "error");
        }
    }

    async function handleImportCSV() {
        try {
            const res = await WailsApp.ImportProductsCSV();
            if (res.startsWith("Éxito")) {
                notifications.show(res, "success");
                await loadProducts();
            } else if (res === "Cancelado") {
                // No hacer nada
            } else {
                notifications.show(res, "error");
            }
        } catch (e) {
            notifications.show("Error importando: " + e, "error");
        }
    }

    function updateTaxPercentage() {
        const map: Record<string, number> = { "0": 0, "2": 12, "4": 15, "5": 5 };
        editingProduct.TaxPercentage = map[editingProduct.TaxCode] || 0;
    }
</script>

<div class="panel full-height" in:fade={{ duration: 200 }}>
    <div class="header-row flex-row space-between align-center">
        <h1>Inventario de Productos</h1>
        <button class="btn-secondary" on:click={handleImportCSV}>
            📥 Importar CSV
        </button>
    </div>

    <div class="master-detail-layout">
        <!-- FORMULARIO -->
        <div class="sidebar-form card">
            <div class="form-header flex-row space-between">
                <h3>{isEditing ? "Editar Producto" : "Nuevo Producto"}</h3>
                <button class="btn-icon-mini" on:click={resetForm} title="Nuevo">✨</button>
            </div>
            
            <div class="sidebar-content">
                <div class="field">
                    <label for="p-sku">Código / SKU</label>
                    <input 
                        id="p-sku" 
                        bind:value={editingProduct.SKU} 
                        placeholder="PROD-001" 
                        class="mono" 
                        disabled={isEditing} 
                    />
                </div>

                <div class="field">
                    <label for="p-name">Descripción</label>
                    <input id="p-name" bind:value={editingProduct.Name} placeholder="Nombre del producto" />
                </div>

                <div class="grid col-2-tight">
                    <div class="field">
                        <label for="p-price">Precio Unit.</label>
                        <input id="p-price" type="number" step="0.01" bind:value={editingProduct.Price} />
                    </div>
                    <div class="field">
                        <label for="p-tax">IVA</label>
                        <select id="p-tax" bind:value={editingProduct.TaxCode} on:change={updateTaxPercentage}>
                            <option value="4">15%</option>
                            <option value="5">5%</option>
                            <option value="2">12%</option>
                            <option value="0">0%</option>
                        </select>
                    </div>
                </div>

                <div class="field">
                    <label for="p-stock">Stock Inicial</label>
                    <input id="p-stock" type="number" bind:value={editingProduct.Stock} />
                </div>
            </div>

            <div class="sidebar-footer mt-4">
                <button class="btn-primary full-width" on:click={handleSave}>
                    {isEditing ? "Actualizar" : "Guardar Producto"}
                </button>
                {#if isEditing}
                    <button class="btn-secondary full-width mt-2" on:click={resetForm}>Cancelar</button>
                {/if}
            </div>
        </div>

        <!-- LISTA -->
        <div class="card no-padding flex-col">
            <!-- Search Bar -->
            <div class="p-3 border-bottom">
                <input 
                    bind:value={searchTerm} 
                    on:input={handleSearch}
                    placeholder="🔍 Buscar producto por nombre o SKU..." 
                    style="width: 100%;"
                />
            </div>

            <div class="linear-grid flex-1 overflow-hidden flex-col">
                <div class="linear-header grid-columns-products">
                    <button class="cell header-btn sortable" on:click={() => sort('SKU')}>
                        SKU {sortCol === 'SKU' ? (sortAsc ? '↑' : '↓') : ''}
                    </button>
                    <button class="cell header-btn sortable" on:click={() => sort('Name')}>
                        Descripción {sortCol === 'Name' ? (sortAsc ? '↑' : '↓') : ''}
                    </button>
                    <div class="cell text-center">IVA %</div>
                    <button class="cell header-btn sortable text-right" on:click={() => sort('Price')}>
                        Precio {sortCol === 'Price' ? (sortAsc ? '↑' : '↓') : ''}
                    </button>
                    <div class="cell text-center">Acciones</div>
                </div>
                
                <div class="rows-container overflow-auto flex-1">
                    {#each sortedProducts as p}
                        <div 
                            class="linear-row grid-columns-products" 
                            role="button"
                            tabindex="0"
                            on:click={() => selectProduct(p)} 
                            on:keydown={(e) => (e.key === 'Enter' || e.key === ' ') && selectProduct(p)}
                            style="cursor: pointer;"
                        >
                            <div class="cell mono text-secondary">{p.SKU}</div>
                            <div class="cell">{p.Name}</div>
                            <div class="cell text-center text-secondary">{p.TaxPercentage}%</div>
                            <div class="cell text-right font-medium text-mint">${p.Price.toFixed(2)}</div>
                            <div class="cell text-center actions-cell">
                                <button class="btn-icon-mini" 
                                    on:click|stopPropagation={() => selectProduct(p)}
                                    title="Editar"
                                >✏️</button>
                                <button class="btn-icon-mini danger" 
                                    on:click|stopPropagation={() => handleDelete(p.SKU)}
                                    title="Eliminar"
                                >🗑️</button>
                            </div>
                        </div>
                    {/each}
                    
                    {#if sortedProducts.length === 0}
                        <div class="empty-state">
                            <div class="empty-state-icon">📦</div>
                            <p class="empty-state-text">No se encontraron productos.</p>
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    </div>
</div>

<style>
    .master-detail-layout {
        display: grid;
        grid-template-columns: 320px 1fr;
        gap: 24px;
        flex: 1;
        min-height: 0;
    }
    
    .flex-col { display: flex; flex-direction: column; }
    .flex-1 { flex: 1; }
    .overflow-hidden { overflow: hidden; }
    .overflow-auto { overflow-y: auto; }
    .text-mint { color: var(--accent-mint); }
    
    .grid-columns-products {
        display: grid;
        grid-template-columns: 120px 1fr 80px 100px 100px;
        align-items: center;
        padding: 0 16px;
    }

    .header-btn {
        background: none;
        border: none;
        color: inherit;
        font: inherit;
        text-align: left;
        padding: 0;
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
    }

    .sortable { cursor: pointer; user-select: none; }
    .sortable:hover { color: var(--text-primary); }

    .p-3 { padding: 16px; }
    .border-bottom { border-bottom: 1px solid var(--border-subtle); }
    
    @media (max-width: 900px) {
        .master-detail-layout {
            grid-template-columns: 1fr;
            grid-template-rows: auto 1fr;
        }
        .sidebar-form { max-height: 350px; overflow-y: auto; }
    }
</style>