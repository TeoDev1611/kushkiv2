<script lang="ts">
    import { onMount } from 'svelte';
    import { fade } from 'svelte/transition';
    import * as WailsApp from 'wailsjs/go/main/App';
    import { notifications } from '$lib/stores/notifications';
    import { withLoading } from '$lib/stores/app';

    let syncLogs: any[] = [];
    let mailLogs: any[] = [];
    let backups: any[] = [];
    let activeSubTab = 'logs'; // 'logs' | 'backups'

    onMount(() => {
        refreshData();
    });

    async function refreshData() {
        try {
            const [sLogs, mLogs, bkps] = await withLoading(Promise.all([
                WailsApp.GetSyncLogs(),
                WailsApp.GetMailLogs(),
                WailsApp.GetBackups()
            ]));
            syncLogs = sLogs || [];
            mailLogs = mLogs || [];
            backups = bkps || [];
        } catch (e) {
            console.error(e);
            notifications.show("Error actualizando logs", "error");
        }
    }

    async function handleSyncNow() {
        try {
            notifications.show("Iniciando sincronización...", "info");
            const msg = await WailsApp.TriggerSyncManual();
            notifications.show(msg, "success");
            setTimeout(refreshData, 2000);
        } catch (e) {
            notifications.show("Error en sync: " + e, "error");
        }
    }

    async function handleCreateBackup() {
        try {
            const err = await withLoading(WailsApp.CreateBackup());
            if (err) notifications.show("Error: " + err, "error");
            else {
                notifications.show("Respaldo creado exitosamente", "success");
                refreshData();
            }
        } catch (e) {
            notifications.show(String(e), "error");
        }
    }

    async function handleSelectBackupPath() {
        const path = await WailsApp.SelectBackupPath();
        if (path) notifications.show("Ruta seleccionada: " + path, "info");
    }
</script>

<div class="panel full-height" in:fade={{ duration: 200 }}>
    <div class="header-row">
        <div>
            <h1>Centro de Control</h1>
            <p class="subtitle">Sincronización, Auditoría y Respaldos</p>
        </div>
        <div class="header-actions">
            <button class="btn-secondary" on:click={refreshData}>🔄 Refrescar</button>
            <button class="btn-primary" on:click={handleSyncNow}>☁️ Sincronizar SRI</button>
        </div>
    </div>

    <!-- Sub-tabs locales -->
    <div class="tabs mb-4">
        <button class="tab-btn" class:active={activeSubTab === 'logs'} on:click={() => activeSubTab = 'logs'}>Logs & Auditoría</button>
        <button class="tab-btn" class:active={activeSubTab === 'backups'} on:click={() => activeSubTab = 'backups'}>Respaldos</button>
    </div>

    {#if activeSubTab === 'logs'}
        <div class="logs-grid">
            <!-- Logs de Correo -->
            <div class="card flex-col">
                <h3>📧 Envíos de Correo</h3>
                <div class="linear-grid flex-1 overflow-hidden flex-col">
                    <div class="linear-header grid-logs">
                        <div class="cell">Estado</div>
                        <div class="cell">Fecha</div>
                        <div class="cell">Destinatario</div>
                    </div>
                    <div class="rows-container overflow-auto flex-1">
                        {#each mailLogs as l}
                            <div class="linear-row grid-logs">
                                <div class="cell"><span class="badge {l.estado}">{l.estado}</span></div>
                                <div class="cell mono text-muted">{l.fecha}</div>
                                <div class="cell">{l.email}</div>
                            </div>
                        {/each}
                        {#if mailLogs.length === 0}
                            <div class="empty-state small">Sin registros recientes</div>
                        {/if}
                    </div>
                </div>
            </div>

            <!-- Logs de Sync -->
            <div class="card flex-col">
                <h3>🔄 Sincronización SRI</h3>
                <div class="linear-grid flex-1 overflow-hidden flex-col">
                    <div class="linear-header grid-logs">
                        <div class="cell">Estado</div>
                        <div class="cell">Fecha</div>
                        <div class="cell">Detalle</div>
                    </div>
                    <div class="rows-container overflow-auto flex-1">
                        {#each syncLogs as l}
                            <div class="linear-row grid-logs">
                                <div class="cell"><span class="badge {l.status === 'OK' ? 'success' : 'warning'}">{l.status}</span></div>
                                <div class="cell mono text-muted">{l.timestamp}</div>
                                <div class="cell text-small">{l.detail}</div>
                            </div>
                        {/each}
                        {#if syncLogs.length === 0}
                            <div class="empty-state small">Sin registros recientes</div>
                        {/if}
                    </div>
                </div>
            </div>
        </div>
    {:else}
        <!-- Backups -->
        <div class="card">
            <div class="flex-row space-between mb-4">
                <h3>🗄️ Copias de Seguridad</h3>
                <div class="flex-row">
                    <button class="btn-secondary" on:click={handleSelectBackupPath}>📂 Configurar Ruta</button>
                    <button class="btn-primary" on:click={handleCreateBackup}>➕ Crear Respaldo Ahora</button>
                </div>
            </div>
            
            <div class="linear-grid">
                <div class="linear-header grid-backups">
                    <div class="cell">Archivo</div>
                    <div class="cell">Fecha</div>
                    <div class="cell">Tamaño</div>
                    <div class="cell">Ruta</div>
                </div>
                <div class="rows-container">
                    {#each backups as b}
                        <div class="linear-row grid-backups">
                            <div class="cell font-medium">{b.name}</div>
                            <div class="cell mono text-muted">{b.date}</div>
                            <div class="cell mono">{b.size}</div>
                            <div class="cell text-muted text-small">{b.path}</div>
                        </div>
                    {/each}
                    {#if backups.length === 0}
                        <div class="empty-state">No hay respaldos generados</div>
                    {/if}
                </div>
            </div>
        </div>
    {/if}
</div>

<style>
    .logs-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 24px;
        height: calc(100vh - 220px);
    }

    .flex-col { display: flex; flex-direction: column; }
    .flex-1 { flex: 1; }
    .overflow-hidden { overflow: hidden; }
    .overflow-auto { overflow-y: auto; }

    .grid-logs {
        display: grid;
        grid-template-columns: 100px 140px 1fr;
        font-size: 12px;
    }

    .grid-backups {
        display: grid;
        grid-template-columns: 1fr 150px 100px 2fr;
    }

    .empty-state.small { padding: 20px; min-height: auto; }
    .text-small { font-size: 0.85em; }

    /* Tabs Simples */
    .tabs {
        display: flex;
        gap: 16px;
        border-bottom: 1px solid var(--border-subtle);
    }
    .tab-btn {
        background: none;
        border: none;
        padding: 12px 24px;
        color: var(--text-secondary);
        border-bottom: 2px solid transparent;
        cursor: pointer;
        font-size: 14px;
    }
    .tab-btn:hover { color: var(--text-primary); }
    .tab-btn.active {
        color: var(--accent-mint);
        border-bottom-color: var(--accent-mint);
    }

    @media (max-width: 1000px) {
        .logs-grid { grid-template-columns: 1fr; height: auto; }
    }
</style>
