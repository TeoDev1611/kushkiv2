<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { fade } from 'svelte/transition';
    import { Backend } from '$lib/services/api';
    import { notifications } from '$lib/stores/notifications';
    import { withLoading } from '$lib/stores/app';
    import type { db } from 'wailsjs/go/models';
    import * as WailsApp from 'wailsjs/go/main/App';
    import ClientForm from '$lib/components/ClientForm.svelte';

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
        // Delegado al componente ClientForm si es necesario, pero aquí usamos el botón del form.
    };

    onMount(() => {
        loadClients();
    });

    async function loadClients() {
        try {
            if (search.length > 2) {
                // Usar búsqueda backend
                clients = await withLoading(WailsApp.SearchClients(search)) || [];
            } else {
                clients = await withLoading(Backend.getClients()) || [];
            }
        } catch (e) {
            notifications.show("Error cargando clientes: " + e, "error");
        }
    }

    let searchTimeout: any;
    function handleSearch() {
        clearTimeout(searchTimeout);
        searchTimeout = setTimeout(() => {
            loadClients();
        }, 300);
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

    function onClientSaved() {
        loadClients();
        if (!isEditing) resetForm();
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

    async function handleImportCSV() {
        try {
            const res = await WailsApp.ImportClientsCSV();
            if (res.startsWith("Error")) {
                notifications.show(res, "error");
            } else if (res !== "Cancelado") {
                notifications.show(res, "success");
                await loadClients();
            }
        } catch (e) {
            notifications.show("Error importando: " + e, "error");
        }
    }
</script>

<div class="panel full-height" in:fade={{ duration: 200 }}>
    <div class="header-row">
        <h1>Directorio de Clientes</h1>
        <div style="flex: 1"></div> 
        <button class="btn-secondary" on:click={handleImportCSV}>
            📥 Importar CSV
        </button>
    </div>

    <div class="master-detail-layout">
        <!-- FORMULARIO (Sidebar Izquierda en Desktop) -->
        <div class="sidebar-form card">
            <div class="form-header flex-row space-between">
                <h3>{isEditing ? "Editar Cliente" : "Nuevo Cliente"}</h3>
                <button class="btn-icon-mini" title="Limpiar" on:click={resetForm}>✨</button>
            </div>
            
            <div class="sidebar-content mt-4">
                <ClientForm 
                    bind:client={editingClient} 
                    isEditing={isEditing}
                    on:saved={onClientSaved}
                    on:cancel={resetForm}
                />
            </div>
        </div>

        <!-- LISTA (Area Principal) -->
        <div class="card no-padding flex-col">
            <!-- Barra de búsqueda -->
            <div class="search-bar p-3 border-bottom">
                <input 
                    bind:value={search} 
                    on:input={handleSearch}
                    placeholder="🔍 Buscar por nombre o RUC..." 
                    style="width: 100%; background: var(--bg-surface); border: 1px solid var(--border-subtle);"
                />
            </div>

            <!-- Tabla -->
            <div class="linear-grid flex-1 overflow-hidden flex-col">
                <div class="linear-header grid-columns-clients">
                    <div class="cell">ID</div>
                    <div class="cell">Nombre</div>
                    <div class="cell">Email</div>
                    <div class="cell">Dirección</div>
                    <div class="cell">Teléfono</div>
                    <div class="cell text-center">Acciones</div>
                </div>
                
                <div class="rows-container overflow-auto flex-1">
                    {#each clients as c}
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
                            <div class="cell text-small">{c.Email || "-"}</div>
                            <div class="cell text-small text-truncate">{c.Direccion || "-"}</div>
                            <div class="cell text-small">{c.Telefono || "-"}</div>
                            
                            <div class="cell text-center actions-cell">
                                <button class="btn-icon-mini" 
                                    on:click|stopPropagation={() => selectClient(c)}
                                    title="Editar"
                                >✏️</button>
                                <button class="btn-icon-mini danger" 
                                    on:click|stopPropagation={() => handleDelete(c.ID)}
                                    title="Eliminar Cliente"
                                    aria-label="Eliminar Cliente"
                                >🗑️</button>
                            </div>
                        </div>
                    {/each}
                    
                    {#if clients.length === 0}
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
        flex: 1;
        min-height: 0;
    }

    .flex-col { display: flex; flex-direction: column; }
    .flex-1 { flex: 1; }
    .overflow-hidden { overflow: hidden; }
    .overflow-auto { overflow-y: auto; }
    .p-3 { padding: 16px; }
    .border-bottom { border-bottom: 1px solid var(--border-subtle); }

    .grid-columns-clients {
        display: grid;
        grid-template-columns: 110px 1fr 180px 180px 100px 90px;
        align-items: center;
        padding: 0 16px;
    }

    .text-small { font-size: 0.85rem; color: var(--text-secondary); }
    .text-truncate {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    /* Responsive: Stack en pantallas pequeñas */
    @media (max-width: 1200px) {
        .grid-columns-clients {
            grid-template-columns: 110px 1fr 140px 90px;
        }
        /* Ocultar dirección y teléfono en pantallas medianas */
        .grid-columns-clients > :nth-child(4),
        .grid-columns-clients > :nth-child(5) {
            display: none;
        }
    }

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