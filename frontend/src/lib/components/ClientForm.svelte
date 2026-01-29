<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { Backend } from '$lib/services/api';
    import { notifications } from '$lib/stores/notifications';
    import { withLoading } from '$lib/stores/app';

    export let client: any = {
        ID: "",
        TipoID: "05",
        Nombre: "",
        Direccion: "",
        Email: "",
        Telefono: ""
    };
    
    // Si es true, muestra botones de guardar/cancelar integrados. 
    // Si es false, el padre controla el guardado (útil para layouts custom).
    export let showActions = true; 
    export let isEditing = false;

    const dispatch = createEventDispatcher();

    async function handleSave() {
        if (!client.ID || !client.Nombre) {
            notifications.show("Cédula/RUC y Nombre son obligatorios", "warning");
            return;
        }

        try {
            const res = await withLoading(Backend.saveClient(client));
            if (res.startsWith("Error")) {
                notifications.show(res, "error");
            } else {
                notifications.show("Cliente guardado exitosamente", "success");
                dispatch('saved', client);
            }
        } catch (e) {
            notifications.show("Error guardando: " + e, "error");
        }
    }

    function handleCancel() {
        dispatch('cancel');
    }
</script>

<div class="client-form">
    <div class="field">
        <label for="cf-id">Identificación (RUC/Cédula)</label>
        <input id="cf-id" bind:value={client.ID} placeholder="099..." disabled={isEditing} />
    </div>

    <div class="field">
        <label for="cf-name">Razón Social / Nombre</label>
        <input id="cf-name" bind:value={client.Nombre} placeholder="Nombre completo" />
    </div>

    <div class="field">
        <label for="cf-email">Correo Electrónico</label>
        <input id="cf-email" type="email" bind:value={client.Email} placeholder="cliente@ejemplo.com" />
    </div>

    <div class="field">
        <label for="cf-addr">Dirección</label>
        <input id="cf-addr" bind:value={client.Direccion} placeholder="Dirección completa" />
    </div>

    <div class="field">
        <label for="cf-tel">Teléfono</label>
        <input id="cf-tel" bind:value={client.Telefono} placeholder="099..." />
    </div>

    {#if showActions}
        <div class="form-actions mt-4">
            <button class="btn-primary full-width" on:click={handleSave}>
                {isEditing ? "Actualizar Datos" : "Registrar Cliente"}
            </button>
            {#if isEditing}
                <button class="btn-secondary full-width mt-2" on:click={handleCancel}>Cancelar</button>
            {/if}
        </div>
    {/if}
</div>

<style>
    .client-form {
        display: flex;
        flex-direction: column;
        gap: 16px;
    }
    
    .field {
        display: flex;
        flex-direction: column;
        gap: 6px;
    }
    
    label {
        font-size: 0.85rem;
        color: var(--text-secondary);
        font-weight: 500;
    }
    
    input {
        background: rgba(0, 0, 0, 0.2);
        border: 1px solid var(--border-subtle);
        padding: 10px;
        border-radius: 6px;
        color: var(--text-primary);
        font-size: 0.95rem;
    }
    
    input:focus {
        border-color: var(--accent-mint);
        outline: none;
    }
    
    input:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .mt-4 { margin-top: 1.5rem; }
    .mt-2 { margin-top: 0.5rem; }
    .full-width { width: 100%; }
</style>