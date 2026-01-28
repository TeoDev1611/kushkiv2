<script>
    import { createEventDispatcher } from "svelte";
    import logo from "../../../assets/images/logo-universal.png";
    import { activeTab } from "$lib/stores/app";

    const dispatch = createEventDispatcher();
    
    const menuItems = [
        { id: 'dashboard', path: 'M3 13h8V3H3v10zm0 8h8v-6H3v6zm10 0h8V11h-8v10zm0-18v6h8V3h-8z', label: 'Dashboard' },
        { id: 'invoice', path: 'M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z', label: 'Facturar' },
        { id: 'quotations', path: 'M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6zM6 20V4h7v5h5v11H6z', label: 'Cotizaciones' },
        { id: 'products', path: 'M20 9V7c0-1.1-.9-2-2-2h-3L15 2H9L7.5 5H5c-1.1 0-2 .9-2 2v2c-1.11 0-2 .89-2 2v4h4v-2h14v2h4v-4c0-1.11-.89-2-2-2zm-6 0h-4V7h4v2z', label: 'Productos' },
        { id: 'clients', path: 'M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z', label: 'Clientes' },
        { id: 'history', path: 'M13 3c-4.97 0-9 4.03-9 9H1l3.89 3.89.07.14L9 12H6c0-3.87 3.13-7 7-7s7 3.13 7 7-3.13 7-7 7c-1.93 0-3.68-.79-4.94-2.06l-1.42 1.42C8.27 19.99 10.51 21 13 21c4.97 0 9-4.03 9-9s-4.03-9-9-9zm-1 5v5l4.28 2.54.72-1.21-3.5-2.08V8H12z', label: 'Historial' },
        { id: 'sync', path: 'M12 4V1L8 5l4 4V6c3.31 0 6 2.69 6 6 0 1.01-.25 1.97-.7 2.8l1.46 1.46C19.54 15.03 20 13.57 20 12c0-4.42-3.58-8-8-8zm0 14c-3.31 0-6-2.69-6-6 0-1.01.25-1.97.7-2.8L5.24 7.74C4.46 8.97 4 10.43 4 12c0 4.42 3.58 8 8 8v3l4-4-4-4v3z', label: 'Sync' },
        { id: 'config', path: 'M19.14 12.94c.04-.3.06-.61.06-.94 0-.32-.02-.64-.07-.94l2.03-1.58a.49.49 0 0 0 .12-.61l-1.92-3.32a.488.488 0 0 0-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54a.484.484 0 0 0-.48-.41h-3.84a.484.484 0 0 0-.48.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96a.488.488 0 0 0-.59.22L2.8 8.87a.49.49 0 0 0 .12.61l2.03 1.58c-.05.3-.09.63-.09.94s.02.64.07.94l-2.03 1.58a.49.49 0 0 0-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.48-.41l.36-2.54c.59-.24 1.13-.57 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32a.49.49 0 0 0-.12-.61l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z', label: 'Ajustes' },
    ];

    function selectTab(tab) {
        activeTab.set(tab);
        dispatch("change", tab);
    }
</script>

<aside class="sidebar">
    <div class="brand">
        <!-- Contenedor rígido para el logo -->
        <div class="logo-box">
            <img src={logo} alt="Kushki" class="logo" />
        </div>
        <div class="brand-text">KUSHKI</div>
    </div>

    <nav class="nav-menu">
        {#each menuItems as item}
            <button
                class="nav-item"
                class:active={$activeTab === item.id}
                on:click={() => selectTab(item.id)}
                title={item.label}
            >
                <div class="icon-box">
                    <svg viewBox="0 0 24 24" class="icon">
                        <path d={item.path} />
                    </svg>
                </div>
                <div class="label-box">
                    {item.label}
                </div>
                {#if $activeTab === item.id}
                    <div class="active-indicator" />
                {/if}
            </button>
        {/each}
    </nav>

    <div class="sidebar-footer">
        <span class="version-text">v2.1</span>
    </div>
</aside>

<style>
    /* COMPORTAMIENTO BASE (MOBILE/TABLET POR DEFECTO O FALLBACK) */
    .sidebar {
        width: 72px; /* Ancho fijo */
        min-width: 72px; /* Nunca encoger */
        flex-shrink: 0; /* Nunca ser aplastado por flexbox */
        height: 100%;
        background: var(--bg-panel);
        border-right: 1px solid var(--border-subtle);
        display: flex;
        flex-direction: column;
        z-index: 500;
        overflow: hidden; /* Corta todo lo que salga del ancho */
        transition: width 0.3s cubic-bezier(0.2, 0, 0, 1);
    }

    /* SOLO EN PANTALLAS GRANDES: Permitir expansión */
    @media (min-width: 1024px) {
        .sidebar:hover {
            width: 240px; /* Expansión suave */
            box-shadow: 4px 0 20px rgba(0,0,0,0.3);
        }
        
        .sidebar:hover .brand-text {
            opacity: 1;
            transform: translateX(0);
        }

        .sidebar:hover .label-box {
            opacity: 1;
            transform: translateX(0);
        }
        
        .sidebar:hover .version-text {
            opacity: 1;
        }
    }

    /* --- BRAND --- */
    .brand {
        height: 64px;
        display: flex;
        align-items: center;
        border-bottom: 1px solid rgba(255,255,255,0.03);
        margin-bottom: 8px;
        position: relative; /* Contexto para posicionamiento */
    }

    .logo-box {
        width: 72px; /* Ancho idéntico al sidebar base */
        min-width: 72px;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 2; /* Encima del texto si pasara algo raro */
        background: var(--bg-panel); /* Parche visual por si acaso */
    }

    .logo {
        width: 32px;
        height: 32px;
        object-fit: contain;
    }

    .brand-text {
        position: absolute;
        left: 72px; /* Empieza donde termina el icono */
        white-space: nowrap;
        font-weight: 800;
        font-size: 18px;
        letter-spacing: 0.05em;
        background: linear-gradient(90deg, #fff, var(--accent-mint));
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        opacity: 0; /* Oculto por defecto */
        transform: translateX(-10px);
        transition: all 0.2s;
    }

    /* --- NAVIGATION --- */
    .nav-menu {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 4px;
        padding: 8px 0;
        overflow-y: auto; /* Scroll si hay muchos items en vertical */
        overflow-x: hidden; /* Jamás scroll horizontal */
    }

    .nav-menu::-webkit-scrollbar { display: none; } /* Ocultar scrollbar visual */

    .nav-item {
        width: 100%;
        height: 48px;
        padding: 0;
        background: transparent;
        border: none;
        cursor: pointer;
        position: relative;
        color: var(--text-secondary);
        transition: all 0.2s;
        
        /* GRID LAYOUT: LA ESTRUCTURA MÁS SEGURA */
        display: grid;
        grid-template-columns: 72px 1fr; /* 72px Fijos | Resto Flexible */
        align-items: center;
    }

    .nav-item:hover {
        background: var(--bg-hover);
        color: var(--text-primary);
    }

    .nav-item.active {
        color: var(--accent-mint);
        background: rgba(52, 211, 153, 0.08);
    }

    /* ICONO: Celda 1 del Grid */
    .icon-box {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .icon {
        width: 22px;
        height: 22px;
        fill: currentColor;
    }

    /* TEXTO: Celda 2 del Grid */
    .label-box {
        text-align: left;
        font-size: 14px;
        font-weight: 500;
        white-space: nowrap;
        opacity: 0; /* Oculto por defecto */
        transform: translateX(-10px);
        transition: opacity 0.2s, transform 0.2s;
        padding-right: 16px;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    /* INDICADOR */
    .active-indicator {
        position: absolute;
        left: 0;
        top: 12px; bottom: 12px;
        width: 3px;
        background: var(--accent-mint);
        border-radius: 0 4px 4px 0;
    }

    /* FOOTER */
    .sidebar-footer {
        height: 48px;
        display: flex;
        align-items: center;
        justify-content: center;
        border-top: 1px solid rgba(255,255,255,0.03);
    }

    .version-text {
        font-size: 10px;
        color: var(--text-tertiary);
        font-family: monospace;
        opacity: 0;
        transition: opacity 0.2s;
        white-space: nowrap;
    }
</style>