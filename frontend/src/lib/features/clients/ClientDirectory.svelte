<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { fade } from 'svelte/transition';
    import { Backend } from '$lib/services/api';
    import { notifications } from '$lib/stores/notifications';
    import { withLoading } from '$lib/stores/app';
    import type { db } from 'wailsjs/go/models';

    // Tipos locales si no están completos en el modelo
    interface ClientUI extends db.ClientDTO {
        // extensiones de UI si fueran necesarias
    }

    let clients: ClientUI[] = [];
    let search = "";
    
    // Estado del formulario
    let isEditing = false;
    let editingClient: any = {
        ID: "",
        TipoID: "05", // Cédula por defecto
        Nombre: "",
        Direccion: "",
        Email: "",
        Telefono: ""
    };

    // Handler para evento global de guardado (Ctrl+S)
    const handleGlobalSave = () => {
        if (editingClient.ID && editingClient.Nombre) {
            handleSave();
        }
    };

    onMount(() => {
        window.addEventListener('app-save', handleGlobalSave);
        loadClients();
    });

    onDestroy(() => {
        window.removeEventListener('app-save', handleGlobalSave);
    });

    async function loadClients() {
        try {
            clients = await withLoading(Backend.getClients()) || [];
        } catch (e) {
            notifications.show("Error cargando clientes: " + e, "error");
        }
    }

    function resetForm() {
        editingClient = {
            ID: "",
            TipoID: "05",
            Nombre: "",
            Direccion: "",
            Email: "",
            Telefono: ""
        };
        isEditing = false; // "Nuevo" mode
    }

    function selectClient(client: ClientUI) {
        // Clonar para evitar mutación directa en la lista antes de guardar
        editingClient = JSON.parse(JSON.stringify(client));
        isEditing = true;
    }

    async function handleSave() {
        if (!editingClient.ID || !editingClient.Nombre) {
            notifications.show("Cédula/RUC y Nombre son obligatorios", "warning");
            return;
        }

        try {
            const res = await withLoading(Backend.saveClient(editingClient));
            if (res.startsWith("Error")) {
                notifications.show(res, "error");
            } else {
                notifications.show("Cliente guardado exitosamente", "success");
                await loadClients();
                if (!isEditing) resetForm(); // Limpiar solo si era creación nueva
            }
        } catch (e) {
            notifications.show("Error guardando: " + e, "error");
        }
    }

    async function handleDelete(id: string) {
        if (!confirm("¿Estás seguro de eliminar este cliente?")) return;
        
        try {
            const res = await withLoading(Backend.deleteClient(id));
            if (res.startsWith("Error")) {
                notifications.show(res, "error");
            } else {
                notifications.show("Cliente eliminado", "success");
                if (editingClient.ID === id) resetForm();
                await loadClients();
            }
        } catch (e) {
            notifications.show("Error eliminando: " + e, "error");
        }
    }

    // Filtrado reactivo
    $: filteredClients = clients.filter(c => 
        c.Nombre.toLowerCase().includes(search.toLowerCase()) || 
        c.ID.includes(search)
    );
</script>

<div class="panel full-height" in:fade={{ duration: 200 }}>
    <div class="header-row">
        <h1>Directorio de Clientes</h1>
        <div style="flex: 1"></div> 
        <!-- Espacio flexible -->
    </div>

    <div class="master-detail-layout">
        <!-- FORMULARIO (Sidebar Izquierda en Desktop) -->
        <div class="sidebar-form card">
            <div class="form-header flex-row space-between">
                <h3>{isEditing ? "Editar Cliente" : "Nuevo Cliente"}</h3>
                <button class="btn-icon-mini" title="Limpiar" on:click={resetForm}>✨</button>
            </div>
            
            <div class="sidebar-content">
                <div class="field">
                    <label for="c-id">Identificación (RUC/Cédula)</label>
                    <input id="c-id" bind:value={editingClient.ID} placeholder="099..." disabled={isEditing && false} />
                </div>

                <div class="field">
                    <label for="c-name">Razón Social / Nombre</label>
                    <input id="c-name" bind:value={editingClient.Nombre} placeholder="Nombre completo" />
                </div>

                <div class="field">
                    <label for="c-email">Correo Electrónico</label>
                    <input id="c-email" type="email" bind:value={editingClient.Email} placeholder="cliente@ejemplo.com" />
                </div>

                <div class="field">
                    <label for="c-addr">Dirección</label>
                    <input id="c-addr" bind:value={editingClient.Direccion} placeholder="Dirección completa" />
                </div>

                <div class="field">
                    <label for="c-tel">Teléfono</label>
                    <input id="c-tel" bind:value={editingClient.Telefono} placeholder="099..." />
                </div>
            </div>

            <div class="sidebar-footer mt-4">
                <button class="btn-primary full-width" on:click={handleSave}>
                    {isEditing ? "Actualizar Datos" : "Registrar Cliente"}
                </button>
                {#if isEditing}
                    <button class="btn-secondary full-width mt-2" on:click={resetForm}>Cancelar Edición</button>
                {/if}
            </div>
        </div>

        <!-- LISTA (Area Principal) -->
        <div class="card no-padding flex-col">
            <!-- Barra de búsqueda -->
            <div class="search-bar p-3 border-bottom">
                <input 
                    bind:value={search} 
                    placeholder="🔍 Buscar por nombre o RUC..." 
                    style="width: 100%; background: var(--bg-surface); border: 1px solid var(--border-subtle);"
                />
            </div>

            <!-- Tabla -->
            <div class="linear-grid flex-1 overflow-hidden flex-col">
                <div class="linear-header grid-columns-clients">
                    <div class="cell">ID</div>
                    <div class="cell">Nombre</div>
                    <div class="cell">Contacto</div>
                    <div class="cell text-center">Acciones</div>
                </div>
                
                <div class="rows-container overflow-auto flex-1">
                    {#each filteredClients as c}
                        <div 
                            class="linear-row grid-columns-clients" 
                            role="button"
                            tabindex="0"
                            on:click={() => selectClient(c)} 
                            on:keydown={(e) => (e.key === 'Enter' || e.key === ' ') && selectClient(c)}
                            style="cursor: pointer;"
                        >
                            <div class="cell mono text-secondary">{c.ID}</div>
                            <div class="cell font-medium">{c.Nombre}</div>
                            <div class="cell text-muted text-small">
                                {c.Email || c.Telefono || "-"}
                            </div>
                            <div class="cell text-center actions-cell">
                                <button class="btn-icon-mini danger" 
                                    on:click|stopPropagation={() => handleDelete(c.ID)}
                                    title="Eliminar Cliente"
                                    aria-label="Eliminar Cliente"
                                >🗑️</button>
                            </div>
                        </div>
                    {/each}
                    
                    {#if filteredClients.length === 0}
                        <div class="empty-state">
                            <div class="empty-state-icon">👥</div>
                            <p class="empty-state-text">
                                {search ? "No hay coincidencias." : "No hay clientes registrados."}
                            </p>
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    </div>
</div>

<style>
    /* Layout específico */
    .master-detail-layout {
        display: grid;
        grid-template-columns: 350px 1fr;
        gap: 24px;
        height: calc(100vh - 140px); /* Ajuste aproximado para el header */
    }

    .flex-col { display: flex; flex-direction: column; }
    .flex-1 { flex: 1; }
    .overflow-hidden { overflow: hidden; }
    .overflow-auto { overflow-y: auto; }
    .p-3 { padding: 16px; }
    .border-bottom { border-bottom: 1px solid var(--border-subtle); }

    .grid-columns-clients {
        display: grid;
        grid-template-columns: 140px 1fr 180px 80px;
        align-items: center;
        padding: 0 16px;
    }

    .text-small { font-size: 0.85rem; }

    /* Responsive: Stack en pantallas pequeñas */
    @media (max-width: 1000px) {
        .master-detail-layout {
            grid-template-columns: 1fr;
            grid-template-rows: auto 1fr;
        }
        .sidebar-form {
            height: auto;
            max-height: 400px;
            overflow-y: auto;
        }
    }
</style>
