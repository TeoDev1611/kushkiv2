<script lang="ts">
    import { onMount } from 'svelte';
    import { fade } from 'svelte/transition';
    
    // Components
    import Sidebar from '$lib/components/layout/Sidebar.svelte';
    import ToastContainer from '$lib/components/ui/ToastContainer.svelte';
    
    // Features
    import Dashboard from '$lib/features/dashboard/Dashboard.svelte';
    import ClientDirectory from '$lib/features/clients/ClientDirectory.svelte';
    import ConfigPanel from '$lib/features/settings/ConfigPanel.svelte';
    import ProductList from '$lib/features/inventory/ProductList.svelte';
    import HistoryPanel from '$lib/features/history/HistoryPanel.svelte';
    import InvoiceEmitter from '$lib/features/invoice/InvoiceEmitter.svelte';
    import QuotationPanel from '$lib/features/quotations/QuotationPanel.svelte';
    import SyncPanel from '$lib/features/sync/SyncPanel.svelte';

    // Stores & Services
    import { activeTab, isLicensed } from '$lib/stores/app';
    import { invoiceStore } from '$lib/stores/invoice';
    import { Backend } from '$lib/services/api';
    import { notifications } from '$lib/stores/notifications';

    let licenseKeyInput = "";
    let verificationLoading = false;

    // --- KEYBOARD SHORTCUTS ---
    function handleGlobalKeydown(event: KeyboardEvent) {
        const target = event.target as HTMLElement;
        const tagName = target.tagName;
        
        // Permitir ESC para salir de foco
        if (event.key === 'Escape') {
            target.blur();
            return;
        }

        // Ignorar atajos si el usuario está escribiendo (excepto Ctrl+S que suele ser deseado)
        const isInput = ['INPUT', 'TEXTAREA', 'SELECT'].includes(tagName);
        if (isInput && !(event.ctrlKey || event.metaKey)) return;

        // Ctrl/Cmd + S: Guardar Contextual
        if ((event.ctrlKey || event.metaKey) && event.key === 's') {
            event.preventDefault();
            
            // Excepción: Facturación (Por seguridad, se prefiere clic explícito)
            let currentTab;
            activeTab.subscribe(v => currentTab = v)();
            
            if (currentTab === 'invoice') {
                notifications.show("Por seguridad, use el botón 'Firmar y Emitir' para facturar.", "info");
            } else {
                // Emitir evento global para que el componente activo lo capture
                window.dispatchEvent(new CustomEvent('app-save'));
            }
        }

        // Ctrl/Cmd + N: Nueva Factura
        if ((event.ctrlKey || event.metaKey) && event.key === 'n') {
            event.preventDefault();
            activeTab.set('invoice');
            invoiceStore.reset();
            notifications.show("Nueva Factura iniciada", "info");
        }

        // Navegación por pestañas (Ctrl + 1..8)
        if ((event.ctrlKey || event.metaKey) && !isNaN(parseInt(event.key))) {
            const num = parseInt(event.key);
            // Orden coincidente con Sidebar.svelte
            const tabs = ['dashboard', 'invoice', 'quotations', 'products', 'clients', 'history', 'sync', 'config'];
            if (num >= 1 && num <= tabs.length) {
                event.preventDefault();
                activeTab.set(tabs[num - 1]);
            }
        }
    }

    onMount(async () => {
        try {
            const licensed = await Backend.checkLicense();
            isLicensed.set(licensed);
            
            if (licensed) {
               await Backend.getConfig();
            }
        } catch (e) {
            console.error("Error checking license", e);
        }
    });

    async function handleActivation() {
        if (!licenseKeyInput.trim()) return;
        
        verificationLoading = true;
        try {
            const res = await Backend.activateLicense(licenseKeyInput.toUpperCase());
            if (res.startsWith("Éxito")) {
                isLicensed.set(true);
                notifications.show("Licencia activada correctamente", "success");
            } else {
                notifications.show(res, "error");
            }
        } catch (e) {
            notifications.show("Error de conexión: " + e, "error");
        } finally {
            verificationLoading = false;
        }
    }

    function handleQuotationConversion(event: CustomEvent) {
        const dto = event.detail;
        invoiceStore.reset(); 
        invoiceStore.update(s => ({
            ...s,
            clienteID: dto.clienteID,
            clienteNombre: dto.clienteNombre,
            clienteDireccion: dto.clienteDireccion,
            clienteEmail: dto.clienteEmail,
            clienteTelefono: dto.clienteTelefono,
            observacion: dto.observacion,
            items: dto.items.map((i: any) => ({
                codigo: i.codigo,
                nombre: i.nombre,
                cantidad: i.cantidad,
                precio: i.precio,
                codigoIVA: "2",
                porcentajeIVA: i.porcentajeIVA
            }))
        }));
        activeTab.set('invoice');
        notifications.show("Datos cargados desde Cotización", "info");
    }
</script>

<svelte:window on:keydown={handleGlobalKeydown} />

<!-- Global Notifications -->
<ToastContainer />

<main class="app-root">
    {#if !$isLicensed}
        <div class="license-overlay" in:fade>
            <div class="license-box card">
                <div class="lock-icon">🔐</div>
                <h2>Activación Requerida</h2>
                <p>Ingresa tu clave de producto para continuar.</p>
                <div class="license-form">
                    <input
                        bind:value={licenseKeyInput}
                        placeholder="KSH-XXXX-XXXX-XXXX"
                        class="text-center"
                    />
                    <button
                        class="btn-primary full-width"
                        on:click={handleActivation}
                        disabled={verificationLoading}
                    >
                        {verificationLoading ? "Verificando..." : "Activar Licencia"}
                    </button>
                </div>
            </div>
        </div>
    {:else}
        <div class="layout-container">
            <Sidebar />
            
            <section class="content-area">
                {#if $activeTab === 'dashboard'}
                    <Dashboard />
                {:else if $activeTab === 'invoice'}
                    <InvoiceEmitter />
                {:else if $activeTab === 'quotations'}
                    <QuotationPanel on:convertToInvoice={handleQuotationConversion}/>
                {:else if $activeTab === 'clients'}
                    <ClientDirectory />
                {:else if $activeTab === 'products'}
                    <ProductList />
                {:else if $activeTab === 'history'}
                    <HistoryPanel />
                {:else if $activeTab === 'sync'}
                    <SyncPanel />
                {:else if $activeTab === 'config'}
                    <ConfigPanel />
                {:else}
                    <div class="placeholder-panel">
                        <h2>Módulo {$activeTab} en construcción 🚧</h2>
                    </div>
                {/if}
            </section>
        </div>
    {/if}
</main>

<style>
    /* Estructura Raíz - Solución al bug de fondo blanco */
    .app-root {
        width: 100vw;
        height: 100vh;
        background-color: var(--bg-app); /* Fondo base */
        /* Patrón sutil fijo */
        background-image: 
            radial-gradient(circle at 15% 50%, rgba(0, 230, 137, 0.03) 0%, transparent 25%),
            radial-gradient(circle at 85% 30%, rgba(41, 152, 255, 0.02) 0%, transparent 25%);
        background-attachment: fixed; /* CLAVE: El fondo no se mueve al hacer scroll */
        background-size: cover;
        color: var(--text-primary);
        overflow: hidden; /* El contenedor raíz no scrollea, scrollean los hijos */
    }

    .layout-container {
        display: flex;
        width: 100%;
        height: 100%;
    }

    .content-area {
        flex: 1;
        min-width: 0; /* CRÍTICO: Permite que el flex item se encoja correctamente */
        overflow-y: auto; 
        padding: 24px;
        position: relative;
    }

    /* License Overlay */
    .license-overlay {
        position: fixed;
        inset: 0;
        background: rgba(11, 15, 25, 0.95); /* Color sólido oscuro */
        z-index: 2000;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    .license-box {
        width: 400px;
        text-align: center;
        padding: 40px;
        border: 1px solid var(--border-subtle);
        background: var(--bg-panel);
    }
    .lock-icon { font-size: 40px; margin-bottom: 20px; }
    .license-form { display: flex; flex-direction: column; gap: 16px; margin-top: 24px; }
    
    .placeholder-panel {
        display: flex; 
        flex-direction: column; 
        align-items: center; 
        justify-content: center; 
        height: 100%;
        color: var(--text-secondary);
    }
</style>
