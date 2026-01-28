<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { fade } from 'svelte/transition';
    import { Backend } from '$lib/services/api';
    import { notifications } from '$lib/stores/notifications';
    import { withLoading } from '$lib/stores/app';
    import type { db } from 'wailsjs/go/models';

    // Estado local
    let products: db.ProductDTO[] = [];
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
    // Solo guarda si hay datos válidos en el formulario
    const handleGlobalSave = () => {
        if (editingProduct.SKU && editingProduct.Name) {
            handleSave();
        } else {
            // Opcional: Feedback si intenta guardar vacío
            // notifications.show("Complete el formulario para guardar", "info");
        }
    };

    onMount(() => {
        window.addEventListener('app-save', handleGlobalSave);
        loadProducts();
    });

    onDestroy(() => {
        window.removeEventListener('app-save', handleGlobalSave);
    });

    async function loadProducts() {
        try {
            products = await withLoading(Backend.getProducts()) || [];
        } catch (e) {
            notifications.show("Error cargando inventario: " + e, "error");
        }
    }

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

    function updateTaxPercentage() {
        const map: Record<string, number> = { "0": 0, "2": 12, "4": 15, "5": 5 };
        editingProduct.TaxPercentage = map[editingProduct.TaxCode] || 0;
    }
</script>

<div class="panel full-height" in:fade={{ duration: 200 }}>
    <div class="header-row">
        <h1>Inventario de Productos</h1>
        <div style="flex: 1"></div>
    </div>

    <div class="master-detail-layout">
        <!-- FORMULARIO -->
        <div class="sidebar-form card">
            <div class="form-header flex-row space-between">
                <h3>{isEditing ? "Editar Producto" : "Nuevo Producto"}</h3>
                <button class="btn-icon-mini" on:click={resetForm}>✨</button>
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
            <div class="linear-grid flex-1 overflow-hidden flex-col">
                <div class="linear-header grid-columns-products">
                    <div class="cell">SKU</div>
                    <div class="cell">Descripción</div>
                    <div class="cell text-right">Precio</div>
                    <div class="cell text-center">Acciones</div>
                </div>
                
                <div class="rows-container overflow-auto flex-1">
                    {#each products as p}
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
                            <div class="cell text-right font-medium text-mint">${p.Price.toFixed(2)}</div>
                            <div class="cell text-center actions-cell">
                                <button class="btn-icon-mini danger" 
                                    on:click|stopPropagation={() => handleDelete(p.SKU)}
                                    title="Eliminar"
                                >🗑️</button>
                            </div>
                        </div>
                    {/each}
                    
                    {#if products.length === 0}
                        <div class="empty-state">
                            <div class="empty-state-icon">📦</div>
                            <p class="empty-state-text">Tu inventario está vacío.</p>
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
        height: calc(100vh - 140px);
    }
    
    .flex-col { display: flex; flex-direction: column; }
    .flex-1 { flex: 1; }
    .overflow-hidden { overflow: hidden; }
    .overflow-auto { overflow-y: auto; }
    .text-mint { color: var(--accent-mint); }
    
    .grid-columns-products {
        display: grid;
        grid-template-columns: 120px 1fr 100px 80px;
        align-items: center;
        padding: 0 16px;
    }
    
    @media (max-width: 900px) {
        .master-detail-layout {
            grid-template-columns: 1fr;
            grid-template-rows: auto 1fr;
        }
        .sidebar-form { max-height: 350px; overflow-y: auto; }
    }
</style>
