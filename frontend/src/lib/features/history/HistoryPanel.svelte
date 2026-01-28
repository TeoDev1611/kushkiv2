<script lang="ts">
    import { onMount } from 'svelte';
    import { fade } from 'svelte/transition';
    import { Backend } from '$lib/services/api';
    import * as WailsApp from 'wailsjs/go/main/App';
    import { notifications } from '$lib/stores/notifications';
    import { withLoading } from '$lib/stores/app';

    // Estado
    let history: any[] = [];
    let currentPage = 1;
    let pageSize = 10;
    let totalItems = 0;
    let totalPages = 1;
    
    // Filtros
    let dateRange = {
        start: new Date(new Date().getFullYear(), new Date().getMonth(), 1).toISOString().split('T')[0],
        end: new Date().toISOString().split('T')[0]
    };
    let searchTerm = "";

    onMount(() => {
        loadHistory();
    });

    async function loadHistory() {
        try {
            // Nota: El backend original usaba paginación simple. 
            // Para búsqueda avanzada, podríamos necesitar implementar un endpoint de búsqueda combinado.
            // Por ahora usamos la paginación básica.
            const res = await withLoading(Backend.getFacturasPaginated(currentPage, pageSize));
            history = res.data || [];
            totalItems = res.total;
            totalPages = Math.ceil(totalItems / pageSize);
        } catch (e) {
            notifications.show("Error cargando historial: " + e, "error");
        }
    }
    
    function changePage(delta: number) {
        const nextPage = currentPage + delta;
        if (nextPage >= 1 && nextPage <= totalPages) {
            currentPage = nextPage;
            loadHistory();
        }
    }

    async function handleExportExcel() {
        try {
            const res = await withLoading(WailsApp.ExportSalesExcel(dateRange.start, dateRange.end));
            notifications.show(res, res.includes("Error") ? "error" : "success");
        } catch (e) {
            notifications.show("Error exportando: " + e, "error");
        }
    }

    // Acciones de Documento
    async function openPDF(clave: string) {
        try {
            const res = await WailsApp.OpenFacturaPDF(clave);
            notifications.show(res, "success");
        } catch (e) { notifications.show(String(e), "error"); }
    }

    async function openXML(clave: string) {
        try {
            const res = await WailsApp.OpenInvoiceXML(clave);
            notifications.show(res, res.includes("Error") ? "error" : "success");
        } catch (e) { notifications.show(String(e), "error"); }
    }

    async function openFolder(clave: string) {
        try {
            const res = await WailsApp.OpenInvoiceFolder(clave);
            notifications.show(res, res.includes("Error") ? "error" : "success");
        } catch (e) { notifications.show(String(e), "error"); }
    }

    async function resendEmail(clave: string) {
        try {
            const res = await withLoading(WailsApp.ResendInvoiceEmail(clave));
            notifications.show(res, res.includes("Error") ? "error" : "success");
        } catch (e) { notifications.show(String(e), "error"); }
    }

</script>

<div class="panel full-height" in:fade={{ duration: 200 }}>
    <div class="header-row">
        <h1>Historial de Transacciones</h1>
        <div class="header-actions">
            <button class="btn-secondary" on:click={handleExportExcel}>📊 Exportar Excel</button>
            <button class="btn-secondary" on:click={loadHistory}>🔄 Refrescar</button>
        </div>
    </div>

    <div class="card no-padding flex-col flex-1">
        <!-- Filtros -->
        <div class="filters-bar p-3 border-bottom flex-row space-between">
            <div class="input-group">
                <input type="date" bind:value={dateRange.start} />
                <input type="date" bind:value={dateRange.end} />
            </div>
            <!-- TODO: Conectar búsqueda al backend cuando soporte filtrado server-side o filtrar en memoria si son pocos datos -->
            <div style="width: 300px; opacity: 0.7;" title="Búsqueda global disponible próximamente">
                 <input bind:value={searchTerm} placeholder="🔍 Buscar..." disabled />
            </div>
        </div>

        <!-- Tabla -->
        <div class="linear-grid flex-1 overflow-hidden flex-col">
            <div class="linear-header grid-columns-history">
                <div class="cell">Estado</div>
                <div class="cell">Secuencial</div>
                <div class="cell">Fecha</div>
                <div class="cell">Cliente</div>
                <div class="cell text-right">Total</div>
                <div class="cell text-center">Acciones</div>
            </div>
            
            <div class="rows-container overflow-auto flex-1">
                {#each history as f}
                    <div class="linear-row grid-columns-history">
                        <div class="cell">
                            <span class="badge {f.estado}">{f.estado}</span>
                        </div>
                        <div class="cell mono text-secondary">{f.secuencial}</div>
                        <div class="cell mono">{f.fecha}</div>
                        <div class="cell text-truncate" title={f.cliente}>{f.cliente}</div>
                        <div class="cell text-right font-medium">${f.total.toFixed(2)}</div>
                        <div class="cell text-center actions-row">
                            {#if f.tienePDF}
                                <button class="btn-icon-mini" title="Ver PDF" aria-label="Ver PDF" on:click={() => openPDF(f.claveAcceso)}>📄</button>
                                <button class="btn-icon-mini" title="Reenviar Email" aria-label="Reenviar Email" on:click={() => resendEmail(f.claveAcceso)}>✉️</button>
                            {/if}
                            <button class="btn-icon-mini" title="Ver XML" aria-label="Ver XML" on:click={() => openXML(f.claveAcceso)}>🌐</button>
                            <button class="btn-icon-mini" title="Abrir Carpeta" aria-label="Abrir Carpeta" on:click={() => openFolder(f.claveAcceso)}>📂</button>
                        </div>
                    </div>
                {/each}

                {#if history.length === 0}
                    <div class="empty-state">
                        <div class="empty-state-icon">📭</div>
                        <p class="empty-state-text">No se encontraron facturas.</p>
                    </div>
                {/if}
            </div>
        </div>

        <!-- Footer Paginación -->
        <div class="pagination-footer p-3 border-top flex-row space-between">
            <button class="btn-secondary" disabled={currentPage === 1} on:click={() => changePage(-1)}>Anterior</button>
            <span class="text-muted">Página {currentPage} de {totalPages || 1} (Total: {totalItems})</span>
            <button class="btn-secondary" disabled={currentPage >= totalPages} on:click={() => changePage(1)}>Siguiente</button>
        </div>
    </div>
</div>

<style>
    .flex-col { display: flex; flex-direction: column; }
    .flex-1 { flex: 1; }
    .overflow-hidden { overflow: hidden; }
    .overflow-auto { overflow-y: auto; }
    .p-3 { padding: 16px; }
    .border-bottom { border-bottom: 1px solid var(--border-subtle); }
    .border-top { border-top: 1px solid var(--border-subtle); }
    .text-truncate { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

    .grid-columns-history {
        display: grid;
        grid-template-columns: 120px 140px 120px 1fr 120px 160px;
        align-items: center;
        padding: 0 16px;
    }
    
    .actions-row {
        display: flex;
        justify-content: center;
        gap: 4px;
    }
</style>
