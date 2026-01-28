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
        sriOnline: false, // Default false para evitar falsos positivos
        salesTrend: [],
    };
    let topProducts: any[] = [];
    let dashboardCharts = { revenueBar: "", clientsPie: "" };
    let recentActivity: any[] = [];
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
            // Cargar en paralelo pero manejar fallos individuales si es necesario
            // Usamos Promise.allSettled para mayor robustez si un servicio falla
            const [kpiRes, chartsRes, topProdRes, facturasRes] = await Promise.allSettled([
                Backend.getDashboardStats(dateRange.start, dateRange.end),
                Backend.getCharts(),
                Backend.getTopProducts(),
                Backend.getFacturasPaginated(1, 5)
            ]);

            // Procesar KPIs
            if (kpiRes.status === 'fulfilled') {
                stats = kpiRes.value;
            } else {
                console.error("Error loading KPIs", kpiRes.reason);
                notifications.show("Error cargando estadísticas", "error");
            }

            // Procesar Gráficos
            if (chartsRes.status === 'fulfilled') {
                dashboardCharts = chartsRes.value || { revenueBar: "", clientsPie: "" };
            }

            // Procesar Top Productos
            if (topProdRes.status === 'fulfilled') {
                topProducts = topProdRes.value || [];
            }

            // Procesar Actividad
            if (facturasRes.status === 'fulfilled') {
                recentActivity = facturasRes.value?.data || [];
            }

        } catch (e) {
            console.error(e);
            notifications.show("Error general en dashboard: " + e, "error");
        } finally {
            loading = false;
        }
    }

    function navigateToHistory() {
        activeTab.set('history');
    }

    onMount(() => {
        loadDashboardData();
    });
</script>

<div class="panel" in:fade={{ duration: 200 }}>
    <!-- Header -->
    <div class="header-row">
        <div>
            <h1>Resumen General</h1>
            <p class="subtitle">Métricas clave del periodo</p>
        </div>
        <div class="header-actions">
            <div class="input-group flex-row" style="gap: 8px; background: var(--bg-surface); padding: 4px 8px; border-radius: 8px; border: 1px solid var(--border-subtle);">
                <input
                    type="date"
                    bind:value={dateRange.start}
                    style="border: none; background: transparent; padding: 6px; width: auto;"
                    title="Fecha Inicio"
                />
                <span style="color: var(--text-tertiary);">-</span>
                <input
                    type="date"
                    bind:value={dateRange.end}
                    style="border: none; background: transparent; padding: 6px; width: auto;"
                    title="Fecha Fin"
                />
            </div>
            <button class="btn-secondary" on:click={loadDashboardData} title="Actualizar Datos">🔄</button>
        </div>
    </div>

    {#if loading}
        <div class="flex-row text-center" style="justify-content: center; padding: 40px;">
            <div class="loader-premium"></div>
        </div>
    {:else}
        <!-- KPIs -->
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
                    <div class="title">Facturas</div>
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
                    <div class="title">Estado SRI</div>
                    <div class="flex-row">
                        <span class="status-dot {stats.sriOnline ? 'AUTORIZADO' : 'PENDIENTE'}"></span>
                        <span style="font-size: 14px; font-weight: 500;">{stats.sriOnline ? "Online" : "Offline"}</span>
                    </div>
                </div>
            </div>
        </div>

        <div class="dashboard-layout mt-4" style="display: flex; flex-direction: column; gap: 24px;">
            <!-- Gráficos & Top Productos -->
            <div class="charts-row" style="display: grid; grid-template-columns: 2fr 1fr; gap: 24px;">
                <div class="section-chart card" style="min-height: 350px;">
                    {#if dashboardCharts.revenueBar || dashboardCharts.clientsPie}
                        <div class="chart-container" style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; height: 100%;">
                            <div class="chart-box">
                                <ChartFrame htmlContent={dashboardCharts.revenueBar} />
                            </div>
                            <div class="chart-box">
                                <ChartFrame htmlContent={dashboardCharts.clientsPie} />
                            </div>
                        </div>
                    {:else}
                        <div class="empty-state" style="border:none; background:none;">
                            <p>No hay datos suficientes para gráficos</p>
                        </div>
                    {/if}
                </div>

                <div class="card compact">
                    <h3>🏆 Top Productos</h3>
                    <div class="linear-list">
                        {#each topProducts.slice(0, 5) as p}
                            <div class="linear-item">
                                <div class="item-info" style="flex:1">{p.name}</div>
                                <div class="item-value">{p.quantity} un.</div>
                            </div>
                        {/each}
                        {#if topProducts.length === 0}
                            <div class="empty-state-text" style="padding: 20px;">Sin datos</div>
                        {/if}
                    </div>
                </div>
            </div>

            <!-- Actividad Reciente -->
            <div class="card">
                <div class="flex-row space-between mb-4">
                    <h3>⚡ Actividad Reciente</h3>
                    <button class="btn-secondary" on:click={navigateToHistory}>Ver Historial Completo</button>
                </div>
                <div class="linear-grid">
                    <div class="linear-header grid-columns-activity" style="grid-template-columns: 100px 140px 1fr 120px;">
                        <div class="cell">Estado</div>
                        <div class="cell">Secuencial</div>
                        <div class="cell">Cliente</div>
                        <div class="cell text-right">Total</div>
                    </div>
                    <div class="rows-container">
                        {#each recentActivity as f}
                            <div class="linear-row grid-columns-activity" style="grid-template-columns: 100px 140px 1fr 120px;">
                                <div class="cell">
                                    <span class="badge {f.estado}">{f.estado}</span>
                                </div>
                                <div class="cell mono text-muted">{f.secuencial}</div>
                                <div class="cell">{f.cliente}</div>
                                <div class="cell text-right font-medium">${(f.total || 0).toFixed(2)}</div>
                            </div>
                        {/each}
                        {#if recentActivity.length === 0}
                            <div class="empty-state-text" style="padding: 40px; text-align: center;">No hay actividad reciente</div>
                        {/if}
                    </div>
                </div>
            </div>
        </div>
    {/if}
</div>

<style>
    /* Estilos específicos que no estén en style.css global */
    .linear-list {
        display: flex;
        flex-direction: column;
    }
    .linear-item {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 12px 0;
        border-bottom: 1px solid var(--border-subtle);
    }
    .linear-item:last-child { border-bottom: none; }
    
    .status-dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        display: block;
        margin-right: 6px;
    }
    .status-dot.AUTORIZADO {
        background: var(--status-success);
        box-shadow: 0 0 8px rgba(0, 255, 148, 0.4);
    }
    .status-dot.PENDIENTE { background: var(--status-warning); }

    /* Responsive Charts */
    @media (max-width: 1100px) {
        .charts-row { grid-template-columns: 1fr !important; }
        .chart-container { grid-template-columns: 1fr !important; }
    }
</style>