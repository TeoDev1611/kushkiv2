<script lang="ts">
    import { onMount } from 'svelte';
    import { fade } from 'svelte/transition';
    import { Backend } from '$lib/services/api';
    import { withLoading, activeTab } from '$lib/stores/app';
    import { notifications } from '$lib/stores/notifications';
    import ChartFrame from '$lib/components/ui/ChartFrame.svelte';

    // Estado local
    let stats = {
        totalVentas: 0,
        totalFacturas: 0,
        pendientes: 0,
        sriOnline: false,
        salesTrend: [],
    };
    let topProducts: any[] = [];
    let dashboardCharts = { revenueBar: "", clientsPie: "" };
    let taxSummary = { ventas15: 0, ventas0: 0, ivaGenerado: 0, retencionesIva: 0, factorProporcion: 1, impuestoSugerido: 0 };
    let recentActivity: any[] = [];
    let emisorConfig: any = null;
    let loading = true;
    
    // Configuración de fechas (Mes actual por defecto)
    const d = new Date();
    let dateRange = {
        start: new Date(d.getFullYear(), d.getMonth(), 1).toISOString().split("T")[0],
        end: new Date().toISOString().split("T")[0]
    };

    // Carga de datos
    async function loadDashboardData() {
        loading = true;
        try {
            const [kpiRes, chartsRes, topProdRes, facturasRes, taxRes, confRes] = await Promise.allSettled([
                Backend.getDashboardStats(dateRange.start, dateRange.end),
                Backend.getCharts(),
                Backend.getTopProducts(),
                Backend.getFacturasPaginated(1, 8), // Pedimos un poco más para llenar la tabla
                Backend.getVATSummary(dateRange.start, dateRange.end),
                Backend.getConfig()
            ]);

            if (kpiRes.status === 'fulfilled') stats = kpiRes.value;
            if (chartsRes.status === 'fulfilled') dashboardCharts = chartsRes.value || { revenueBar: "", clientsPie: "" };
            if (topProdRes.status === 'fulfilled') topProducts = topProdRes.value || [];
            if (facturasRes.status === 'fulfilled') recentActivity = facturasRes.value?.data || [];
            if (taxRes.status === 'fulfilled') taxSummary = taxRes.value;
            if (confRes.status === 'fulfilled') emisorConfig = confRes.value;

        } catch (e) {
            console.error(e);
            notifications.show("Error cargando dashboard: " + e, "error");
        } finally {
            loading = false;
        }
    }

    function getDeadline(ruc: string) {
        if (!ruc || ruc.length < 9) return "--";
        const ninth = parseInt(ruc[8]);
        const day = ninth === 0 ? 28 : (ninth * 2) + 8;
        return `${day} del próximo mes`;
    }

    function navigateToHistory() {
        activeTab.set('history');
    }

    onMount(() => {
        loadDashboardData();
    });
</script>

<div class="panel" in:fade={{ duration: 200 }}>
    <!-- Header Minimalista -->
    <div class="header-row">
        <div>
            <h1>Dashboard</h1>
            <p class="subtitle text-secondary">Control de ventas y obligaciones SRI</p>
        </div>
        <div class="header-actions">
            <div class="date-selector-group flex-row">
                <input type="date" bind:value={dateRange.start} on:change={loadDashboardData} />
                <span class="text-muted">al</span>
                <input type="date" bind:value={dateRange.end} on:change={loadDashboardData} />
            </div>
            <button class="btn-secondary" on:click={loadDashboardData} title="Refrescar">🔄</button>
        </div>
    </div>

    {#if loading}
        <div class="flex-row text-center" style="justify-content: center; padding: 100px;">
            <div class="loader-premium"></div>
        </div>
    {:else}
        <!-- Fila 1: KPIs Esenciales -->
        <div class="kpi-row">
            <div class="kpi-card">
                <div class="kpi-icon mint">💰</div>
                <div class="kpi-content">
                    <div class="title">Ventas Totales</div>
                    <div class="value gradient-text">${(stats.totalVentas || 0).toFixed(2)}</div>
                </div>
            </div>
            <div class="kpi-card">
                <div class="kpi-icon blue">📄</div>
                <div class="kpi-content">
                    <div class="title">Facturas Emitidas</div>
                    <div class="value">{stats.totalFacturas || 0}</div>
                </div>
            </div>
            <div class="kpi-card">
                <div class="kpi-icon orange">⚠️</div>
                <div class="kpi-content">
                    <div class="title">Pendientes</div>
                    <div class="value">{stats.pendientes || 0}</div>
                </div>
            </div>
            <div class="kpi-card">
                <div class="kpi-icon {stats.sriOnline ? 'green' : 'red'}">
                    {stats.sriOnline ? "🌐" : "🔌"}
                </div>
                <div class="kpi-content">
                    <div class="title">Servidor SRI</div>
                    <div class="flex-row">
                        <span class="status-dot {stats.sriOnline ? 'AUTORIZADO' : 'PENDIENTE'}"></span>
                        <span style="font-weight: 600;">{stats.sriOnline ? "En Línea" : "Offline"}</span>
                    </div>
                </div>
            </div>
        </div>

        <div class="dashboard-grid-main mt-4">
            
            <!-- COLUMNA IZQUIERDA: Guía Fiscal y Gráfico -->
            <div class="left-stack">
                <!-- GUÍA FISCAL -->
                <div class="tax-card card">
                    <div class="flex-row space-between mb-4">
                        <h3 class="m-0">📋 Guía Formulario 104</h3>
                        {#if emisorConfig}
                            <div class="deadline-badge">
                                📅 Declara hasta: <strong>{getDeadline(emisorConfig.RUC)}</strong>
                            </div>
                        {/if}
                    </div>

                    <div class="tax-grid">
                        <div class="tax-item">
                            <span class="tax-label">Ventas 15% (Cas. 401)</span>
                            <span class="tax-value">${taxSummary.ventas15.toFixed(2)}</span>
                        </div>
                        <div class="tax-item">
                            <span class="tax-label">IVA Cobrado (Cas. 411)</span>
                            <span class="tax-value text-mint">${taxSummary.ivaGenerado.toFixed(2)}</span>
                        </div>
                        <div class="tax-item">
                            <span class="tax-label">Retenciones (Cas. 609)</span>
                            <span class="tax-value text-orange">${taxSummary.retencionesIva.toFixed(2)}</span>
                        </div>
                        <div class="tax-item highlight">
                            <span class="tax-label">Sugerido a Pagar</span>
                            <span class="tax-value">${taxSummary.impuestoSugerido.toFixed(2)}</span>
                        </div>
                    </div>
                </div>

                <!-- GRÁFICO TOP CLIENTES -->
                <div class="card" style="flex: 1; min-height: 380px; padding: 10px;">
                    <ChartFrame htmlContent={dashboardCharts.clientsPie} />
                </div>
            </div>

            <!-- COLUMNA DERECHA: Productos y Actividad -->
            <div class="right-stack">
                <!-- TOP PRODUCTOS -->
                <div class="card compact mb-4">
                    <h3 class="mb-3">🏆 Productos Estrella</h3>
                    <div class="linear-list">
                        {#each topProducts.slice(0, 4) as p}
                            <div class="linear-item">
                                <div style="flex:1">
                                    <div class="font-bold">{p.name}</div>
                                    <div class="text-xs text-muted">{p.sku}</div>
                                </div>
                                <div class="text-right">
                                    <div class="text-mint font-bold">{p.quantity} un.</div>
                                    <div class="text-xs text-secondary">${(p.total || 0).toFixed(2)}</div>
                                </div>
                            </div>
                        {/each}
                        {#if topProducts.length === 0}
                            <p class="empty-text">Sin ventas registradas</p>
                        {/if}
                    </div>
                </div>

                <!-- ACTIVIDAD RECIENTE -->
                <div class="card flex-1 no-padding overflow-hidden flex-col">
                    <div class="p-4 flex-row space-between border-bottom">
                        <h3 class="m-0">⚡ Últimos Movimientos</h3>
                        <button class="btn-icon-mini" on:click={navigateToHistory} title="Ver todo">➡</button>
                    </div>
                    <div class="linear-grid">
                        {#each recentActivity as f}
                            <div class="linear-row compact-row">
                                <div class="cell">
                                    <span class="badge {f.estado}" style="zoom: 0.8;">{f.estado}</span>
                                </div>
                                <div class="cell" style="flex: 1;">
                                    <div class="font-medium text-truncate" style="max-width: 120px;">{f.cliente}</div>
                                    <div class="text-xs mono text-muted">{f.secuencial}</div>
                                </div>
                                <div class="cell text-right font-bold text-mint">${(f.total || 0).toFixed(2)}</div>
                            </div>
                        {/each}
                        {#if recentActivity.length === 0}
                            <p class="p-4 text-center text-muted">No hay facturas este mes</p>
                        {/if}
                    </div>
                </div>
            </div>

        </div>
    {/if}
</div>

<style>
    .dashboard-grid-main {
        display: grid;
        grid-template-columns: 1fr 380px;
        gap: 24px;
        align-items: start;
    }

    .left-stack, .right-stack {
        display: flex;
        flex-direction: column;
        gap: 24px;
    }

    .right-stack {
        height: 100%;
    }

    .date-selector-group {
        background: var(--bg-surface);
        padding: 4px 12px;
        border-radius: 10px;
        border: 1px solid var(--border-subtle);
    }
    
    .date-selector-group input {
        border: none;
        background: transparent;
        width: auto;
        padding: 4px;
        font-size: 13px;
    }

    .tax-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 16px;
    }

    .tax-item {
        padding: 14px;
        background: rgba(255,255,255,0.03);
        border-radius: 12px;
        border: 1px solid var(--border-subtle);
    }
    
    .tax-label { font-size: 10px; color: var(--text-secondary); text-transform: uppercase; margin-bottom: 6px; }
    .tax-value { font-size: 18px; font-weight: 800; }
    .tax-item.highlight { background: rgba(52, 211, 153, 0.05); border-color: var(--mint-border); }

    .deadline-badge {
        font-size: 11px;
        padding: 4px 12px;
        background: var(--bg-hover);
        border-radius: 99px;
        border: 1px solid var(--border-medium);
    }

    .compact-row {
        grid-template-columns: auto 1fr auto;
        padding: 8px 16px;
    }

    .text-truncate {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .m-0 { margin: 0; }
    .p-4 { padding: 16px; }
    .border-bottom { border-bottom: 1px solid var(--border-subtle); }
    .linear-item {
        display: flex;
        align-items: center;
        padding: 10px 0;
        border-bottom: 1px solid var(--border-subtle);
    }
    .linear-item:last-child { border-bottom: none; }
    .empty-text { padding: 20px; text-align: center; color: var(--text-tertiary); font-style: italic; }

    @media (max-width: 1100px) {
        .dashboard-grid-main {
            grid-template-columns: 1fr;
        }
    }
</style>