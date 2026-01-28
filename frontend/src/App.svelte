<script>
    import {
        CreateInvoice,
        GetEmisorConfig,
        SaveEmisorConfig,
        SelectCertificate,
        SelectStoragePath,
        SelectAndSaveLogo,
        GetProducts,
        SearchProducts,
        SaveProduct,
        DeleteProduct,
        GetClients,
        SearchClients,
        SaveClient,
        DeleteClient,
        GetNextSecuencial,
        GetDashboardStats,
        GetFacturasPaginated,
        OpenFacturaPDF,
        OpenInvoiceFolder,
        OpenInvoiceXML,
        ExportSalesExcel,
        ResendInvoiceEmail,
        GetTopProducts,
        GetSyncLogs,
        GetMailLogs,
        TriggerSyncManual,
        CheckLicense,
        ActivateLicense,
        TestSMTPConnection,
        CreateBackup,
        GetBackups,
        SelectBackupPath,
        SearchInvoicesSmart,
        GetStatisticsCharts,
    } from "../wailsjs/go/main/App.js";
    import { EventsOn } from "../wailsjs/runtime/runtime.js";
    import { onMount } from "svelte";
    import { fade, fly, slide } from "svelte/transition";
    import { flip } from "svelte/animate";
    import { cubicOut } from "svelte/easing";
    import Sidebar from "./components/Sidebar.svelte";
    import logo from "./assets/images/logo-universal.png";
    import Wizard from "./components/Wizard.svelte";
    import QuotationPanel from "./components/QuotationPanel.svelte";
    import ChartFrame from "./components/ChartFrame.svelte";

    // --- UTILIDADES ---
    function debounce(func, wait) {
        let timeout;
        return function (...args) {
            const context = this;
            clearTimeout(timeout);
            timeout = setTimeout(() => func.apply(context, args), wait);
        };
    }

    // --- ESTADO GLOBAL ---
    let activeTab = "dashboard";
    let loading = false;
    let initialLoading = true; // Estado para la Splash Screen
    let isLicensed = false; // Estado de Licencia
    let licenseKeyInput = ""; // Input para la licencia
    let showWizard = false;
    let contentArea;

    // --- NOTIFICACIONES & TOASTS ---
    let toasts = []; // { id, message, type }
    let notifications = []; // { id, message, type, time }
    let showNotifications = false; // Toggle para el panel de notificaciones
    let confirmationModal = {
        show: false,
        title: "",
        message: "",
        onConfirm: null,
    };

    // --- CONFIGURACIÓN Y FACTURACIÓN ---
    let config = {
        RUC: "",
        RazonSocial: "",
        NombreComercial: "",
        Direccion: "",
        P12Path: "",
        P12Password: "",
        Ambiente: 1,
        Estab: "001",
        PtoEmi: "001",
        Obligado: false,
        StoragePath: "",
    };
    let invoice = {
        secuencial: "000000001",
        clienteID: "",
        clienteNombre: "",
        clienteDireccion: "",
        clienteEmail: "",
        clienteTelefono: "",
        formaPago: "01",
        items: [],
    };
    let invoiceErrors = {}; // Para rastrear campos inválidos: { clienteID: true, items: true }
    let newItem = {
        codigo: "",
        nombre: "",
        cantidad: 1,
        precio: 0,
        codigoIVA: "4",
        porcentajeIVA: 15,
    };

    let clientSearch = "";
    let productSearch = "";
    let showPassword = false;

    // --- GESTIÓN DE PRODUCTOS ---
    let products = [];
    let editingProduct = {
        SKU: "",
        Name: "",
        Price: 0,
        Stock: 0,
        TaxCode: "4",
        TaxPercentage: 15,
    };

    // --- GESTIÓN DE CLIENTES ---
    let clients = [];
    let editingClient = {
        ID: "",
        TipoID: "05",
        Nombre: "",
        Direccion: "",
        Email: "",
        Telefono: "",
    };

    // --- ESTADO SINCRONIZACIÓN Y AUDITORÍA ---
    let syncLogs = [];
    let mailLogs = [];
    let selectedLog = null;
    let backups = [];

    // --- DASHBOARD CHARTS ---
    let dashboardCharts = { revenueBar: "", clientsPie: "" };

    // --- DASHBOARD & HISTORIAL ---
    let stats = {
        totalVentas: 0,
        totalFacturas: 0,
        pendientes: 0,
        sriOnline: true,
        salesTrend: [],
    };
    let topProducts = [];
    let historial = [];
    let historialSearch = "";
    let dateRange = { start: "", end: "" };

    // --- COMPUTED ---
    $: subtotal = invoice.items.reduce(
        (sum, item) => sum + item.cantidad * item.precio,
        0,
    );
    $: iva = invoice.items.reduce(
        (sum, item) =>
            sum + item.cantidad * item.precio * (item.porcentajeIVA / 100),
        0,
    );
    $: total = subtotal + iva;

    // Smart Search logic
    const handleHistorySearch = debounce(async (term) => {
        if (term.length > 1) {
            loading = true;
            try {
                historial = await SearchInvoicesSmart(term);
            } catch (e) {
                console.error(e);
            } finally {
                loading = false;
            }
        } else if (term.length === 0) {
            loadFacturasPage();
        }
    }, 400);

    // Use the search result directly, no manual filtering needed anymore for smart search
    $: filteredHistorial = historial;

    // Resetear scroll al cambiar de pestaña
    $: if (activeTab && contentArea) contentArea.scrollTop = 0;

    onMount(() => {
        // Timeout de seguridad por si el backend tarda demasiado
        const safetyTimer = setTimeout(() => {
            initialLoading = false;
            loading = false;
        }, 8000);

        // Cargar datos y manejar Splash Screen
        loadData().finally(() => {
            clearTimeout(safetyTimer);
            setTimeout(() => {
                initialLoading = false;
            }, 2000);
        });

        EventsOn("toast-notification", (data) => {
            showToast(data.message, data.type);
            if (data.type === "success") refreshDashboard();
        });
    });

    // --- KEYBOARD SHORTCUTS HANDLER ---
    function handleGlobalKeydown(event) {
        // Ignorar si el usuario está escribiendo en un input
        if (
            ["INPUT", "TEXTAREA", "SELECT"].includes(
                document.activeElement.tagName,
            )
        ) {
            // Permitir Esc para salir de foco
            if (event.key === "Escape") document.activeElement.blur();
            return;
        }

        // Ctrl/Cmd + S: Guardar (Contextual)
        if ((event.ctrlKey || event.metaKey) && event.key === "s") {
            event.preventDefault();
            if (activeTab === "config") handleSaveConfig();
            if (activeTab === "products") handleSaveProduct();
            if (activeTab === "clients") handleSaveClient();
            if (activeTab === "invoice")
                showToast("Use el botón para emitir (Seguridad)", "info");
        }

        // Ctrl/Cmd + N: Nueva Factura
        if ((event.ctrlKey || event.metaKey) && event.key === "n") {
            event.preventDefault();
            activeTab = "invoice";
            // Opcional: Resetear factura aquí si se desea
        }

        // Navegación rápida por tabs (1-6)
        if ((event.ctrlKey || event.metaKey) && !isNaN(parseInt(event.key))) {
            const num = parseInt(event.key);
            const tabs = [
                "dashboard",
                "invoice",
                "products",
                "clients",
                "history",
                "quotations",
            ];
            if (num >= 1 && num <= tabs.length) {
                event.preventDefault();
                activeTab = tabs[num - 1];
            }
        }
    }

    async function refreshSyncLogs() {
        loading = true;
        try {
            const [sLogs, mLogs, bkps] = await Promise.all([
                GetSyncLogs(),
                GetMailLogs(),
                GetBackups(),
            ]);
            syncLogs = sLogs || [];
            mailLogs = mLogs || [];
            backups = bkps || [];
        } catch (e) {
            console.error(e);
        } finally {
            loading = false;
        }
    }

    async function handleCreateBackup() {
        loading = true;
        try {
            const err = await CreateBackup();
            if (err) showToast("Error: " + err, "error");
            else {
                showToast("Respaldo creado exitosamente", "success");
                backups = await GetBackups();
            }
        } catch (e) {
            showToast(e, "error");
        } finally {
            loading = false;
        }
    }

    async function handleSelectBackupFolder() {
        const path = await SelectBackupPath();
        if (path)
            showToast(
                "Ruta seleccionada: " +
                    path +
                    ". (Recuerde guardar en Configuración)",
                "info",
            );
    }

    async function loadData() {
        loading = true;
        try {
            const licensed = await CheckLicense();
            if (!licensed) {
                isLicensed = false;
                return;
            }
            isLicensed = true;

            const cfg = await GetEmisorConfig();

            const isDefaultRUC =
                cfg &&
                (cfg.RUC === "1790011223001" || cfg.RUC === "9999999999999");
            const isMissingData =
                !cfg ||
                !cfg.RazonSocial ||
                cfg.RazonSocial === "EMISOR DE PRUEBA S.A." ||
                cfg.RazonSocial === "Nuevo Usuario";

            if (!cfg || isDefaultRUC || isMissingData) {
                showWizard = true;
                if (cfg) config = cfg;
            } else {
                config = cfg;
            }

            const [prods, cli, seq] = await Promise.all([
                GetProducts(),
                GetClients(),
                GetNextSecuencial(),
            ]);

            products = prods || [];
            clients = cli || [];
            invoice.secuencial = seq || "000000001";

            await refreshDashboard();
        } catch (e) {
            console.error("Error arrancando:", e);
            showToast("Error de conexión o arranque: " + e, "error");
        } finally {
            loading = false;
        }
    }

    async function handleActivation() {
        const licenseRegex = /^KSH-[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}$/;
        const key = licenseKeyInput.trim().toUpperCase();

        if (!licenseRegex.test(key)) {
            return showToast(
                "Formato inválido. Debe ser KSH-XXXX-XXXX-XXXX",
                "error",
            );
        }

        loading = true;
        try {
            const res = await ActivateLicense(key);
            if (res.startsWith("Éxito")) {
                showToast(res, "success");
                isLicensed = true;
                await loadData();
            } else {
                showToast(res, "error");
            }
        } catch (e) {
            showToast("Error crítico: " + e, "error");
        } finally {
            loading = false;
        }
    }

    function onWizardComplete() {
        showWizard = false;
        loadData();
    }

    const handleClientSearchInput = debounce(async (e) => {
        const term = e.target.value;
        if (term.length > 2) {
            await SearchClients(term);
        }
    }, 400);

    // Paginación
    let currentPage = 1;
    let pageSize = 10;
    let totalFacturas = 0;
    $: totalPages = Math.ceil(totalFacturas / pageSize);

    // Inicializar fechas
    const d = new Date();
    dateRange.start = new Date(d.getFullYear(), d.getMonth(), 1)
        .toISOString()
        .split("T")[0];
    dateRange.end = new Date().toISOString().split("T")[0];

    async function refreshDashboard() {
        loading = true;
        try {
            stats = await GetDashboardStats(dateRange.start, dateRange.end);
            await loadFacturasPage();
            topProducts = (await GetTopProducts()) || [];
            dashboardCharts = await GetStatisticsCharts();
            showToast("Datos actualizados", "info");
        } catch (e) {
            console.error("Error refreshing dashboard:", e);
            showToast("Error actualizando: " + e, "error");
        } finally {
            loading = false;
        }
    }

    // Auto-refrescar cuando cambian las fechas
    $: if (dateRange.start || dateRange.end) {
        if (isLicensed) refreshDashboard();
    }

    async function loadFacturasPage() {
        const res = await GetFacturasPaginated(currentPage, pageSize);
        historial = res.data || [];
        totalFacturas = res.total;
    }

    function changePage(newPage) {
        if (newPage >= 1 && newPage <= totalPages) {
            currentPage = newPage;
            loadFacturasPage();
        }
    }

    async function handleExportExcel() {
        loading = true;
        try {
            const res = await ExportSalesExcel(dateRange.start, dateRange.end);
            showToast(res, res.includes("Error") ? "error" : "success");
        } catch (err) {
            showToast("Error al exportar: " + err, "error");
        } finally {
            loading = false;
        }
    }

    function showToast(message, type = "success") {
        const id = Date.now() + Math.random();
        toasts = [...toasts, { id, message, type }];
        if (toasts.length > 3) toasts = toasts.slice(1);
        setTimeout(() => {
            toasts = toasts.filter((t) => t.id !== id);
        }, 6000);
        notifications = [
            { id, message, type, time: new Date() },
            ...notifications,
        ];
    }

    function showConfirmation(title, message, callback) {
        confirmationModal = { show: true, title, message, onConfirm: callback };
    }

    function handleConfirm() {
        if (confirmationModal.onConfirm) confirmationModal.onConfirm();
        confirmationModal.show = false;
    }

    async function handleSelectCert() {
        const path = await SelectCertificate();
        if (path) config.P12Path = path;
    }

    async function handleSelectStorage() {
        const path = await SelectStoragePath();
        if (path) config.StoragePath = path;
    }

    async function handleTestSMTP() {
        if (!config.SMTPHost || !config.SMTPUser || !config.SMTPPassword) {
            return showToast("Complete los datos de correo primero", "error");
        }
        loading = true;
        try {
            const res = await TestSMTPConnection(config);
            showToast(res, res.includes("Error") ? "error" : "success");
        } catch (e) {
            showToast("Fallo crítico: " + e, "error");
        } finally {
            loading = false;
        }
    }

    async function handleSaveConfig() {
        loading = true;
        try {
            config.Ambiente = parseInt(config.Ambiente);
            config.SMTPPort = parseInt(config.SMTPPort);
            const res = await SaveEmisorConfig(config);
            if (res.startsWith("Error")) showToast(res, "error");
            else showToast(res, "success");
        } catch (err) {
            showToast(err, "error");
        } finally {
            loading = false;
        }
    }

    async function handleOpenPDF(clave) {
        showToast(await OpenFacturaPDF(clave), "success");
    }

    async function handleOpenXML(clave) {
        const res = await OpenInvoiceXML(clave);
        showToast(res, res.includes("Error") ? "error" : "success");
    }

    async function handleOpenFolder(clave) {
        const res = await OpenInvoiceFolder(clave);
        showToast(res, res.includes("Error") ? "error" : "success");
    }

    async function handleResendEmail(clave) {
        const res = await ResendInvoiceEmail(clave);
        showToast(res, res.includes("Error") ? "error" : "success");
    }

    async function handleSaveProduct() {
        const res = await SaveProduct(editingProduct);
        showToast(res, res.includes("Error") ? "error" : "success");
        editingProduct = {
            SKU: "",
            Name: "",
            Price: 0,
            Stock: 0,
            TaxCode: "4",
            TaxPercentage: 15,
        };
        products = await GetProducts();
    }

    async function handleDeleteProduct(sku) {
        showConfirmation(
            "Eliminar Producto",
            "¿Estás seguro de que deseas eliminar este producto?",
            async () => {
                const res = await DeleteProduct(sku);
                showToast(res, res.includes("Error") ? "error" : "success");
                products = await GetProducts();
            },
        );
    }

    function selectProduct(p) {
        const taxCode = p.TaxCode || (p.TaxPercentage > 0 ? "4" : "0");
        const taxPerc =
            p.TaxPercentage !== undefined
                ? p.TaxPercentage
                : p.TaxCode == "4"
                  ? 15
                  : 0;
        newItem = {
            codigo: p.SKU,
            nombre: p.Name,
            cantidad: 1,
            precio: p.Price,
            codigoIVA: taxCode,
            porcentajeIVA: taxPerc,
        };
        productSearch = "";
    }

    async function handleSaveClient() {
        const res = await SaveClient(editingClient);
        showToast(res, res.includes("Error") ? "error" : "success");
        editingClient = {
            ID: "",
            TipoID: "05",
            Nombre: "",
            Direccion: "",
            Email: "",
            Telefono: "",
        };
        clients = await GetClients();
    }

    async function handleDeleteClient(id) {
        showConfirmation(
            "Eliminar Cliente",
            "¿Estás seguro de que deseas eliminar este cliente?",
            async () => {
                const res = await DeleteClient(id);
                showToast(res, res.includes("Error") ? "error" : "success");
                clients = await GetClients();
            },
        );
    }

    function selectClient(c) {
        invoice.clienteID = c.ID;
        invoice.clienteNombre = c.Nombre;
        invoice.clienteDireccion = c.Direccion;
        invoice.clienteEmail = c.Email || "";
        invoice.clienteTelefono = c.Telefono || "";
        clientSearch = "";
    }

    async function handleSyncNow() {
        const msg = await TriggerSyncManual();
        showToast(msg, "info");
        setTimeout(refreshSyncLogs, 1000);
    }

    function formatJSON(str) {
        if (!str) return "";
        try {
            const obj = JSON.parse(str);
            return JSON.stringify(obj, null, 2);
        } catch (e) {
            return str;
        }
    }

    function addItem() {
        if (!newItem.nombre || newItem.precio <= 0) return;
        invoice.items = [...invoice.items, { ...newItem }];
        newItem = {
            codigo: "",
            nombre: "",
            cantidad: 1,
            precio: 0,
            codigoIVA: "4",
            porcentajeIVA: 15,
        };
    }

    function removeItem(index) {
        invoice.items = invoice.items.filter((_, i) => i !== index);
    }

    async function handleEmit() {
        invoiceErrors = {};
        let isValid = true;

        if (!invoice.clienteID) {
            invoiceErrors.clienteID = true;
            isValid = false;
        }
        if (!invoice.clienteNombre) {
            invoiceErrors.clienteNombre = true;
            isValid = false;
        }
        if (!invoice.clienteDireccion) {
            invoiceErrors.clienteDireccion = true;
            isValid = false;
        }
        if (!invoice.clienteEmail || !invoice.clienteEmail.includes("@")) {
            invoiceErrors.clienteEmail = true;
            isValid = false;
        }

        if (invoice.items.length === 0) {
            showToast("Agregue al menos un producto", "error");
            return;
        }

        if (!isValid) {
            showToast(
                "Por favor complete los campos obligatorios marcados en rojo",
                "error",
            );
            return;
        }

        loading = true;
        try {
            const dto = {
                secuencial: invoice.secuencial,
                clienteID: invoice.clienteID,
                clienteNombre: invoice.clienteNombre,
                clienteDireccion: invoice.clienteDireccion,
                clienteEmail: invoice.clienteEmail,
                clienteTelefono: invoice.clienteTelefono,
                observacion: invoice.observacion,
                formaPago: invoice.formaPago,
                items: invoice.items.map((i) => ({
                    codigo: i.codigo,
                    nombre: i.nombre,
                    cantidad: parseFloat(i.cantidad),
                    precio: parseFloat(i.precio),
                    codigoIVA: i.codigoIVA.toString(),
                    porcentajeIVA: parseFloat(i.porcentajeIVA),
                })),
            };
            const result = await CreateInvoice(dto);
            if (result.startsWith("Éxito")) {
                showToast(result, "success");
                invoice.secuencial = (parseInt(invoice.secuencial) + 1)
                    .toString()
                    .padStart(9, "0");
                invoice.items = [];
                await refreshDashboard();
            } else {
                showToast(result, "error");
            }
        } catch (err) {
            showToast("Error crítico: " + err, "error");
        } finally {
            loading = false;
        }
    }

    function handleConvertToInvoice(event) {
        const dto = event.detail;
        invoice.clienteID = dto.clienteID;
        invoice.clienteNombre = dto.clienteNombre;
        invoice.clienteDireccion = dto.clienteDireccion;
        invoice.clienteEmail = dto.clienteEmail;
        invoice.clienteTelefono = dto.clienteTelefono;
        invoice.observacion = dto.observacion;

        invoice.items = dto.items.map((i) => ({
            codigo: i.codigo,
            nombre: i.nombre,
            cantidad: i.cantidad,
            precio: i.precio,
            codigoIVA: "2",
            porcentajeIVA: i.porcentajeIVA,
        }));

        activeTab = "invoice";
        showToast("Datos cargados desde Cotización", "info");
    }
</script>

<svelte:window on:keydown={handleGlobalKeydown} />

{#if initialLoading}
    <div class="splash-screen" out:fade={{ duration: 800, easing: cubicOut }}>
        <div class="splash-content">
             <div class="splash-brand" in:fly="{{ y: 20, duration: 800, delay: 200 }}">
                <img src={logo} alt="Kushki" style="width: 50px; height: auto;" />
                <span class="gradient-text">KUSHKI</span>
             </div>
            <div class="loader-premium"></div>
        </div>
    </div>
{/if}

<main>
    {#if !isLicensed}
        <div class="license-overlay" in:fade={{ duration: 300 }}>
            <div class="license-box card">
                <div class="lock-icon">🔐</div>
                <h2>Activación Requerida</h2>
                <p>
                    Bienvenido al Sistema de Facturación Kushki. <br />Ingresa
                    tu clave de producto para continuar.
                </p>

                <div class="license-form">
                    <input
                        aria-label="Clave de producto"
                        bind:value={licenseKeyInput}
                        placeholder="KSH-XXXX-XXXX-XXXX"
                        class="text-center"
                        on:input={(e) =>
                            (licenseKeyInput = e.target.value.toUpperCase())}
                    />
                    <button
                        class="btn-primary full-width"
                        on:click={handleActivation}
                        disabled={loading}
                    >
                        {loading
                            ? "Verificando con Servidor..."
                            : "Activar Licencia"}
                    </button>
                </div>
            </div>
        </div>
    {:else}
        <Sidebar {activeTab} on:change={(e) => (activeTab = e.detail)} />

        <section class="content" bind:this={contentArea}>
            {#if activeTab === "dashboard"}
                <div
                    class="panel"
                    in:fade={{ duration: 200, easing: cubicOut }}
                >
                    <div class="header-row">
                        <div>
                            <h1>Hola, <span class="gradient-text">{config.RazonSocial || "Emisor"}</span> 👋</h1>
                            <p class="subtitle">
                                Resumen general de facturación
                            </p>
                        </div>
                        <div class="header-actions">
                            <div class="notification-wrapper">
                                <button
                                    class="btn-icon-mini"
                                    on:click={() =>
                                        (showNotifications =
                                            !showNotifications)}
                                    title="Notificaciones"
                                    aria-label="Notificaciones"
                                >
                                    🔔
                                    {#if notifications.length > 0}<span
                                            class="notification-badge"
                                        ></span>{/if}
                                </button>

                                {#if showNotifications}
                                    <div
                                        class="notification-panel"
                                        transition:fly={{
                                            y: 10,
                                            duration: 200,
                                        }}
                                    >
                                        <div class="notif-header">
                                            <h4>Notificaciones</h4>
                                            <button
                                                class="link-btn-small"
                                                on:click={() =>
                                                    (notifications = [])}
                                                >Borrar todo</button
                                            >
                                        </div>
                                        <div class="notif-list">
                                            {#each notifications as n}
                                                <div
                                                    class="notif-item {n.type}"
                                                >
                                                    <div class="notif-msg">
                                                        {n.message}
                                                    </div>
                                                    <div class="notif-time">
                                                        {n.time.toLocaleTimeString()}
                                                    </div>
                                                </div>
                                            {/each}
                                            {#if notifications.length === 0}
                                                <div class="empty-notif">
                                                    No hay notificaciones
                                                    recientes
                                                </div>
                                            {/if}
                                        </div>
                                    </div>
                                {/if}
                            </div>

                            <div class="input-group flex-row" style="gap: 8px; background: var(--bg-surface); padding: 4px 8px; border-radius: 8px; border: 1px solid var(--border-subtle);">
                                <input
                                    type="date"
                                    bind:value={dateRange.start}
                                    aria-label="Fecha inicio"
                                    style="border: none; background: transparent; padding: 6px; width: auto;"
                                />
                                <span style="color: var(--text-tertiary);">-</span>
                                <input
                                    type="date"
                                    bind:value={dateRange.end}
                                    aria-label="Fecha fin"
                                    style="border: none; background: transparent; padding: 6px; width: auto;"
                                />
                            </div>
                            <button
                                class="btn-secondary"
                                on:click={refreshDashboard}
                                title="Actualizar">🔄</button
                            >
                        </div>
                    </div>

                    <div class="kpi-row">
                        <!-- KPI Cards remain the same -->
                        <div class="kpi-card">
                            <div class="kpi-icon mint">💰</div>
                            <div class="kpi-content">
                                <div class="title">Ventas Totales</div>
                                <div class="value gradient-text">
                                    ${stats.totalVentas.toFixed(2)}
                                </div>
                            </div>
                        </div>
                        <div class="kpi-card">
                            <div class="kpi-icon blue">📄</div>
                            <div class="kpi-content">
                                <div class="title">Facturas</div>
                                <div class="value">{stats.totalFacturas}</div>
                            </div>
                        </div>
                        <div class="kpi-card">
                            <div class="kpi-icon orange">⚠️</div>
                            <div class="kpi-content">
                                <div class="title">Pendientes</div>
                                <div class="value">{stats.pendientes}</div>
                            </div>
                        </div>
                        <div class="kpi-card">
                            <div
                                class="kpi-icon {stats.sriOnline
                                    ? 'green'
                                    : 'red'}"
                            >
                                {stats.sriOnline ? "🌐" : "🔌"}
                            </div>
                            <div class="kpi-content">
                                <div class="title">Estado SRI</div>
                                <div class="flex-row">
                                    <span
                                        class="status-dot {stats.sriOnline
                                            ? 'AUTORIZADO'
                                            : 'PENDIENTE'}"
                                    ></span>
                                    <span
                                        style="font-size: 14px; font-weight: 500;"
                                        >{stats.sriOnline
                                            ? "Online"
                                            : "Offline"}</span
                                    >
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="dashboard-layout mt-4" style="display: flex; flex-direction: column; gap: 24px;">
                        
                        <!-- Row 1: Charts & Top Products -->
                        <div class="charts-row" style="display: grid; grid-template-columns: 2fr 1fr; gap: 24px;">
                            <div class="section-chart card" style="min-height: 350px;">
                                <div class="chart-container" style="display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; height: 100%;">
                                    <div class="chart-box">
                                        <ChartFrame htmlContent={dashboardCharts.revenueBar} />
                                    </div>
                                    <div class="chart-box">
                                        <ChartFrame htmlContent={dashboardCharts.clientsPie} />
                                    </div>
                                </div>
                            </div>

                            <div class="card compact">
                                <h3>🏆 Top Productos</h3>
                                <div class="linear-list">
                                    {#each topProducts.slice(0, 5) as p}
                                        <div class="linear-item">
                                            <div class="item-info" style="flex:1">
                                                {p.name}
                                            </div>
                                            <div class="item-value">
                                                {p.quantity} un.
                                            </div>
                                        </div>
                                    {/each}
                                    {#if topProducts.length === 0}
                                        <div class="empty-state-text" style="padding: 20px;">Sin datos</div>
                                    {/if}
                                </div>
                            </div>
                        </div>

                        <!-- Row 2: Recent Activity (Full Width) -->
                        <div class="card">
                            <div class="flex-row space-between mb-4">
                                <h3>⚡ Actividad Reciente</h3>
                                <button
                                    class="btn-secondary"
                                    on:click={() => (activeTab = "history")}
                                    >Ver Historial Completo</button
                                >
                            </div>
                            <div class="linear-grid">
                                <div class="linear-header grid-columns-activity" style="grid-template-columns: 100px 140px 1fr 120px;">
                                    <div class="cell">Estado</div>
                                    <div class="cell">Secuencial</div>
                                    <div class="cell">Cliente</div>
                                    <div class="cell text-right">Total</div>
                                </div>
                                <div class="rows-container">
                                    {#each (historial || []).slice(0, 5) as f}
                                        <div class="linear-row grid-columns-activity" style="grid-template-columns: 100px 140px 1fr 120px;">
                                            <div class="cell">
                                                <span class="badge {f.estado}">{f.estado}</span>
                                            </div>
                                            <div class="cell mono text-muted">{f.secuencial}</div>
                                            <div class="cell">{f.cliente}</div>
                                            <div class="cell text-right font-medium">${f.total.toFixed(2)}</div>
                                        </div>
                                    {/each}
                                    {#if !historial || historial.length === 0}
                                        <div class="empty-state-text" style="padding: 40px; text-align: center;">No hay actividad reciente</div>
                                    {/if}
                                </div>
                            </div>
                        </div>
                    </div>

                    <style>
                        /* Local styles for dashboard lists */
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
                        .linear-item:last-child {
                            border-bottom: none;
                        }
                        .item-value {
                            font-weight: 600;
                            font-family: "SF Mono", monospace;
                        }
                        .item-meta {
                            color: var(--text-secondary);
                            margin-right: 16px;
                            font-size: 12px;
                        }
                        .status-dot {
                            width: 8px;
                            height: 8px;
                            border-radius: 50%;
                            display: block;
                        }
                        .status-dot.AUTORIZADO {
                            background: var(--status-success);
                            box-shadow: 0 0 8px rgba(0, 255, 148, 0.4);
                        }
                        .status-dot.PENDIENTE {
                            background: var(--status-warning);
                        }
                        .status-dot.ANULADO {
                            background: var(--status-error);
                        }
                    </style>
                </div>
            {:else if activeTab === "quotations"}
                <QuotationPanel
                    {clients}
                    {products}
                    on:convertToInvoice={handleConvertToInvoice}
                    on:toast={(e) => showToast(e.detail.message, e.detail.type)}
                />
            {:else if activeTab === "history"}
                <div class="panel" in:fade={{ duration: 200 }}>
                    <div class="header-row">
                        <h1>Historial de Transacciones</h1>
                        <div class="header-actions">
                            <button
                                class="btn-secondary"
                                on:click={handleExportExcel}
                                >📊 Exportar Excel</button
                            >
                            <button
                                class="btn-secondary"
                                on:click={loadFacturasPage}>🔄 Refrescar</button
                            >
                        </div>
                    </div>

                    <div class="card no-padding">
                        <div
                            class="flex-row space-between"
                            style="padding: 16px; border-bottom: 1px solid var(--border-subtle); background: var(--bg-panel);"
                        >
                            <div class="input-group">
                                <input
                                    type="date"
                                    bind:value={dateRange.start}
                                    aria-label="Desde"
                                />
                                <input
                                    type="date"
                                    bind:value={dateRange.end}
                                    aria-label="Hasta"
                                />
                            </div>
                            <div style="width: 320px;">
                                <input
                                    bind:value={historialSearch}
                                    on:input={(e) =>
                                        handleHistorySearch(e.target.value)}
                                    placeholder="🔍 Buscar (Nombre, RUC, Secuencial...)"
                                    aria-label="Buscar factura"
                                />
                            </div>
                        </div>

                        <div class="linear-grid">
                            <div class="linear-header grid-columns-history">
                                <div class="cell">Estado</div>
                                <div class="cell">Secuencial</div>
                                <div class="cell">Fecha</div>
                                <div class="cell">Cliente</div>
                                <div class="cell text-right">Total</div>
                                <div class="cell text-center">Acciones</div>
                            </div>
                            <div class="rows-container">
                                {#each filteredHistorial as f}
                                    <div
                                        class="linear-row grid-columns-history"
                                    >
                                        <div class="cell">
                                            <span class="badge {f.estado}"
                                                >{f.estado}</span
                                            >
                                        </div>
                                        <div class="cell mono text-secondary">
                                            {f.secuencial}
                                        </div>
                                        <div class="cell mono">{f.fecha}</div>
                                        <div class="cell" title={f.cliente}>
                                            {f.cliente}
                                        </div>
                                        <div
                                            class="cell text-right font-medium"
                                        >
                                            ${f.total.toFixed(2)}
                                        </div>
                                        <div
                                            class="cell text-center"
                                            style="gap: 4px; justify-content: center;"
                                        >
                                            {#if f.tienePDF}
                                                <button
                                                    class="btn-icon-mini"
                                                    title="Ver PDF"
                                                    on:click={() =>
                                                        handleOpenPDF(
                                                            f.claveAcceso,
                                                        )}>📄</button
                                                >
                                                <button
                                                    class="btn-icon-mini"
                                                    title="Email"
                                                    on:click={() =>
                                                        handleResendEmail(
                                                            f.claveAcceso,
                                                        )}>✉️</button
                                                >
                                            {/if}
                                            <button
                                                class="btn-icon-mini"
                                                title="XML"
                                                on:click={() =>
                                                    handleOpenXML(
                                                        f.claveAcceso,
                                                    )}>🌐</button
                                            >
                                            <button
                                                class="btn-icon-mini"
                                                title="Carpeta"
                                                on:click={() =>
                                                    handleOpenFolder(
                                                        f.claveAcceso,
                                                    )}>📂</button
                                            >
                                        </div>
                                    </div>
                                {/each}
                            </div>
                            {#if filteredHistorial.length === 0}
                                <div class="empty-state">
                                    <div class="empty-state-icon">📭</div>
                                    <div class="empty-state-text">
                                        No se encontraron transacciones en este rango de fechas.
                                    </div>
                                </div>
                            {/if}
                        </div>

                        <style>
                            .grid-columns-history {
                                grid-template-columns: 120px 140px 120px 1fr 120px 140px;
                            }
                            .font-medium {
                                font-weight: 500;
                            }
                            .text-secondary {
                                color: var(--text-secondary);
                            }
                        </style>
                        <div
                            class="flex-row space-between"
                            style="padding: 1rem;"
                        >
                            <button
                                class="btn-secondary"
                                disabled={currentPage === 1}
                                on:click={() => changePage(currentPage - 1)}
                                >Anterior</button
                            >
                            <span class="text-muted"
                                >Página {currentPage} de {totalPages || 1}</span
                            >
                            <button
                                class="btn-secondary"
                                disabled={currentPage === totalPages}
                                on:click={() => changePage(currentPage + 1)}
                                >Siguiente</button
                            >
                        </div>
                    </div>
                </div>
            {:else if activeTab === "invoice"}
                <div class="panel" in:fade={{ duration: 200 }}>
                    <div class="header-row">
                        <h1>Emitir Factura</h1>
                        <button
                            class="btn-secondary"
                            style="border-color: {config.Ambiente === 1
                                ? 'var(--status-warning)'
                                : 'var(--status-success)'}"
                            on:click={async () => {
                                config.Ambiente = config.Ambiente === 1 ? 2 : 1;
                                await handleSaveConfig();
                                showToast(
                                    `Ambiente: ${config.Ambiente === 1 ? "PRUEBAS" : "PRODUCCIÓN"}`,
                                    "info",
                                );
                            }}
                        >
                            {config.Ambiente === 1
                                ? "🧪 MODO PRUEBAS"
                                : "🚀 MODO PRODUCCIÓN"}
                        </button>
                    </div>

                    <!-- CLIENTE -->
                    <div class="card">
                        <div class="header-row" style="margin-bottom: 24px; border-bottom: none; padding-bottom: 0; align-items: flex-start;">
                            <div style="flex: 1;">
                                <div class="flex-row space-between" style="margin-bottom: 12px;">
                                    <h3>👤 Datos del Cliente</h3>
                                    <div
                                        class="badge mono"
                                        style="font-size: 13px; color: var(--accent-mint); border: 1px solid var(--mint-border); background: var(--mint-dim); padding: 6px 12px;"
                                    >
                                        DOC: {config.Estab}-{config.PtoEmi}-{invoice.secuencial}
                                    </div>
                                </div>
                                <div style="position: relative;">
                                    <input
                                        bind:value={clientSearch}
                                        on:input={(e) => handleClientSearchInput(e)}
                                        placeholder="🔍 Buscar o crear nuevo cliente..."
                                        style="background: var(--bg-hover); border-color: var(--border-medium); height: 44px; font-size: 15px;"
                                    />
                                    {#if clientSearch.length > 0}
                                        <div class="dropdown-menu">
                                            {#each clients.filter((c) => c.Nombre.toLowerCase().includes(clientSearch.toLowerCase()) || c.ID.includes(clientSearch)) as c}
                                                <button
                                                    class="dropdown-item"
                                                    on:click={() => selectClient(c)}
                                                    >{c.Nombre} ({c.ID})</button
                                                >
                                            {/each}
                                        </div>
                                    {/if}
                                </div>
                            </div>
                        </div>

                        <div class="invoice-client-grid">
                            <div class="field">
                                <label for="inv-id">Identificación</label>
                                <input
                                    id="inv-id"
                                    bind:value={invoice.clienteID}
                                    class={invoiceErrors.clienteID ? "invalid" : ""}
                                    placeholder="RUC / Cédula"
                                />
                            </div>
                            <div class="field">
                                <label for="inv-name">Razón Social</label>
                                <input
                                    id="inv-name"
                                    bind:value={invoice.clienteNombre}
                                    class={invoiceErrors.clienteNombre ? "invalid" : ""}
                                    placeholder="Nombre del Cliente"
                                />
                            </div>
                            <div class="field">
                                <label for="inv-email">Email</label>
                                <input
                                    id="inv-email"
                                    type="email"
                                    bind:value={invoice.clienteEmail}
                                    class={invoiceErrors.clienteEmail ? "invalid" : ""}
                                    placeholder="cliente@email.com"
                                />
                            </div>
                            <div class="field">
                                <label for="inv-tel">Teléfono</label>
                                <input
                                    id="inv-tel"
                                    bind:value={invoice.clienteTelefono}
                                    placeholder="0999999999"
                                />
                            </div>
                            <div class="field" style="grid-column: span 2;">
                                <label for="inv-addr">Dirección</label>
                                <input
                                    id="inv-addr"
                                    bind:value={invoice.clienteDireccion}
                                    class={invoiceErrors.clienteDireccion ? "invalid" : ""}
                                    placeholder="Dirección completa"
                                />
                            </div>
                        </div>
                    </div>

                    <!-- PRODUCTOS -->
                    <div class="card" style="min-height: 400px; display: flex; flex-direction: column;">
                        <div class="header-row" style="border-bottom: none; padding-bottom: 0; margin-bottom: 16px;">
                            <h3>📦 Items</h3>
                        </div>

                        <div class="invoice-product-bar">
                            <div style="position: relative; grid-column: span 1;">
                                <input
                                    bind:value={productSearch}
                                    placeholder="🔍 Buscar producto..."
                                    style="width: 100%;"
                                />
                                {#if productSearch.length > 0}
                                    <div class="dropdown-menu">
                                        {#each products.filter((p) => p.Name.toLowerCase().includes(productSearch.toLowerCase()) || p.SKU.includes(productSearch)) as p}
                                            <button
                                                class="dropdown-item"
                                                on:click={() => selectProduct(p)}
                                                >{p.Name} - ${p.Price}</button
                                            >
                                        {/each}
                                    </div>
                                {/if}
                            </div>

                            <input bind:value={newItem.codigo} placeholder="Código" />
                            <input bind:value={newItem.nombre} placeholder="Descripción" />
                            <input type="number" bind:value={newItem.cantidad} min="1" placeholder="Cant" class="text-center" />
                            <input type="number" step="0.01" bind:value={newItem.precio} placeholder="Precio" class="text-right" />
                            
                            <select bind:value={newItem.codigoIVA} on:change={() => {
                                const map = { "0": 0, "2": 12, "4": 15, "5": 5 };
                                newItem.porcentajeIVA = map[newItem.codigoIVA];
                            }}>
                                <option value="4">15%</option>
                                <option value="5">5%</option>
                                <option value="0">0%</option>
                                <option value="2">12%</option>
                            </select>

                            <button class="btn-neon icon-only" on:click={addItem} title="Agregar">
                                <span style="font-size: 20px; line-height: 1;">+</span>
                            </button>
                        </div>

                        <div class="linear-grid mt-4" style="flex: 1;">
                            <div class="linear-header grid-columns-invoice">
                                <div class="cell text-center">Cant</div>
                                <div class="cell">Descripción</div>
                                <div class="cell text-right">P. Unit</div>
                                <div class="cell text-center">IVA</div>
                                <div class="cell text-right">Total</div>
                                <div class="cell"></div>
                            </div>
                            <div class="rows-container">
                                {#each invoice.items as item, i (item)}
                                    <div
                                        class="linear-row grid-columns-invoice"
                                        animate:flip={{ duration: 300 }}
                                    >
                                        <div class="cell text-center">
                                            {item.cantidad}
                                        </div>
                                        <div class="cell">{item.nombre}</div>
                                        <div class="cell text-right">
                                            ${item.precio.toFixed(2)}
                                        </div>
                                        <div class="cell text-center">
                                            <span class="badge"
                                                >{item.porcentajeIVA}%</span
                                            >
                                        </div>
                                        <div
                                            class="cell text-right font-medium"
                                        >
                                            ${(
                                                item.cantidad *
                                                item.precio *
                                                (1 + item.porcentajeIVA / 100)
                                            ).toFixed(2)}
                                        </div>
                                        <div class="cell text-center">
                                            <button
                                                class="btn-icon-mini danger"
                                                on:click={() => removeItem(i)}
                                                >×</button
                                            >
                                        </div>
                                    </div>
                                {/each}
                            </div>
                            {#if invoice.items.length === 0}
                                <div
                                    style="padding: 40px; text-align: center; color: var(--text-secondary); font-style: italic;"
                                >
                                    Agrega productos para comenzar la factura
                                </div>
                            {/if}
                        </div>
                    </div>

                    <!-- FOOTER -->
                    <div class="footer-panel">
                        <div class="notes-area" style="flex:1">
                            <label for="inv-notes">Notas Adicionales</label>
                            <textarea
                                id="inv-notes"
                                bind:value={invoice.observacion}
                                placeholder="Observaciones..."
                                style="height: 120px; resize: none;"
                            ></textarea>
                        </div>
                        <div class="totals-area card">
                            <div class="total-row text-secondary">
                                <span>Subtotal</span><span
                                    >${subtotal.toFixed(2)}</span
                                >
                            </div>
                            <div class="total-row text-secondary">
                                <span>IVA</span><span>${iva.toFixed(2)}</span>
                            </div>
                            <div
                                class="total-row grand-total flex-row space-between"
                                style="border-top: 1px solid var(--border-subtle); margin-top: 16px; padding-top: 16px;"
                            >
                                <span style="font-size: 18px; font-weight: 700;"
                                    >TOTAL</span
                                >
                                <span
                                    style="font-size: 18px; font-weight: 700; color: var(--accent-mint);"
                                    >${total.toFixed(2)}</span
                                >
                            </div>
                            <button
                                class="btn-primary full-width"
                                style="margin-top: 24px; width: 100%; height: 48px; font-size: 14px;"
                                on:click={handleEmit}
                                disabled={loading}
                            >
                                {loading
                                    ? "Procesando..."
                                    : "🖋️ Firmar y Emitir"}
                            </button>
                        </div>
                    </div>

                    <style>
                        .grid-columns-invoice {
                            grid-template-columns: 80px 1fr 100px 80px 120px 60px;
                        }
                        .add-product-bar-linear {
                            display: flex;
                            gap: 12px;
                            align-items: center;
                            background: var(--bg-hover);
                            padding: 16px;
                            border-radius: var(--radius-sm);
                            border: 1px solid var(--border-subtle);
                            margin-top: 16px;
                        }
                        .dropdown-menu {
                            position: absolute;
                            top: 100%;
                            left: 0;
                            right: 0;
                            background: var(--bg-panel);
                            border: 1px solid var(--accent-mint);
                            z-index: 50;
                            max-height: 200px;
                            overflow-y: auto;
                            border-radius: 0 0 6px 6px;
                            box-shadow: 0 10px 40px rgba(0, 0, 0, 0.5);
                        }
                        .dropdown-item {
                            width: 100%;
                            padding: 12px;
                            background: transparent;
                            border: none;
                            border-bottom: 1px solid var(--border-subtle);
                            color: var(--text-primary);
                            text-align: left;
                            cursor: pointer;
                        }
                        .dropdown-item:hover {
                            background: rgba(0, 255, 148, 0.1);
                            color: var(--accent-mint);
                        }
                        .total-row {
                            display: flex;
                            justify-content: space-between;
                            margin-bottom: 8px;
                            font-size: 14px;
                        }
                    </style>
                </div>
            {:else if activeTab === "products"}
                <div class="panel full-height" in:fade={{ duration: 200 }}>
                    <div class="header-row">
                        <h1>Inventario de Productos</h1>
                    </div>

                    <div class="master-detail-layout">
                        <!-- FORMULARIO PRODUCTO -->
                        <div class="sidebar-form card">
                            <h3>
                                {editingProduct.SKU ? "Editar" : "Nuevo"} Producto
                            </h3>
                            <div class="field">
                                <label for="p-sku">SKU / Código</label>
                                <input
                                    id="p-sku"
                                    bind:value={editingProduct.SKU}
                                    placeholder="COD-001"
                                    style="font-family: monospace;"
                                />
                            </div>
                            <div class="field">
                                <label for="p-name">Nombre</label>
                                <input
                                    id="p-name"
                                    bind:value={editingProduct.Name}
                                    placeholder="Nombre del producto"
                                />
                            </div>
                            <div class="grid col-2-tight" style="gap: 12px; display: grid; grid-template-columns: 1fr 1fr;">
                                <div class="field">
                                    <label for="p-price">Precio</label>
                                    <input
                                        id="p-price"
                                        type="number"
                                        step="0.01"
                                        bind:value={editingProduct.Price}
                                    />
                                </div>
                                <div class="field">
                                    <label for="p-tax">Impuesto</label>
                                    <select
                                        id="p-tax"
                                        bind:value={editingProduct.TaxCode}
                                        on:change={() => {
                                            const map = { "0": 0, "2": 12, "4": 15, "5": 5 };
                                            editingProduct.TaxPercentage = map[editingProduct.TaxCode];
                                        }}
                                    >
                                        <option value="4">15%</option>
                                        <option value="5">5%</option>
                                        <option value="0">0%</option>
                                        <option value="2">12%</option>
                                    </select>
                                </div>
                            </div>
                            <div class="flex-row mt-4">
                                <button
                                    class="btn-primary full-width"
                                    style="flex:1"
                                    on:click={handleSaveProduct}>Guardar</button
                                >
                                {#if editingProduct.SKU}
                                    <button
                                        class="btn-secondary"
                                        on:click={() =>
                                            (editingProduct = {
                                                SKU: "",
                                                Name: "",
                                                Price: 0,
                                                Stock: 0,
                                                TaxCode: "4",
                                                TaxPercentage: 15,
                                            })}>Cancelar</button
                                    >
                                {/if}
                            </div>
                        </div>

                        <!-- LISTA PRODUCTOS -->
                        <div class="card no-padding" style="display: flex; flex-direction: column;">
                            <div class="linear-grid" style="flex: 1; overflow: hidden; display: flex; flex-direction: column;">
                                <div class="linear-header grid-columns-products">
                                    <div class="cell">SKU</div>
                                    <div class="cell">Nombre</div>
                                    <div class="cell text-right">Precio</div>
                                    <div class="cell text-center">Acciones</div>
                                </div>
                                <div class="rows-container" style="overflow-y: auto; flex: 1;">
                                    {#each products as p}
                                        <div class="linear-row grid-columns-products">
                                            <div class="cell mono text-secondary">{p.SKU}</div>
                                            <div class="cell">{p.Name}</div>
                                            <div class="cell text-right font-medium" style="color: var(--accent-mint);">
                                                ${p.Price.toFixed(2)}
                                            </div>
                                            <div class="cell text-center" style="gap: 4px; justify-content: center;">
                                                <button class="btn-icon-mini" on:click={() => (editingProduct = { ...p })}>✏️</button>
                                                <button class="btn-icon-mini danger" on:click={() => handleDeleteProduct(p.SKU)}>🗑️</button>
                                            </div>
                                        </div>
                                    {/each}
                                </div>
                                {#if products.length === 0}
                                    <div class="empty-state">
                                        <div class="empty-state-icon">📦</div>
                                        <div class="empty-state-text">No hay productos en inventario.</div>
                                    </div>
                                {/if}
                            </div>
                        </div>
                    </div>
                </div>
            {:else if activeTab === "clients"}
                <div class="panel full-height" in:fade={{ duration: 200 }}>
                    <div class="header-row">
                        <h1>Directorio de Clientes</h1>
                    </div>

                    <div class="master-detail-layout">
                        <!-- FORMULARIO CLIENTE -->
                        <div class="sidebar-form card" style="display: flex; flex-direction: column; height: 100%; padding: 24px;">
                            <div class="sidebar-content">
                                <h3>
                                    {editingClient.ID ? "Editar" : "Nuevo"} Cliente
                                </h3>
                                
                                <div class="field">
                                    <label for="c-id">Identificación</label>
                                    <input
                                        id="c-id"
                                        bind:value={editingClient.ID}
                                        placeholder="RUC / Cédula"
                                    />
                                </div>

                                <div class="field">
                                    <label for="c-name">Razón Social</label>
                                    <input
                                        id="c-name"
                                        bind:value={editingClient.Nombre}
                                        placeholder="Nombre completo"
                                    />
                                </div>

                                <div class="grid col-2-tight" style="gap: 12px; display: grid; grid-template-columns: 1fr 1fr;">
                                    <div class="field">
                                        <label for="c-email">Email</label>
                                        <input
                                            id="c-email"
                                            bind:value={editingClient.Email}
                                            type="email"
                                            placeholder="correo@ejemplo.com"
                                        />
                                    </div>
                                    <div class="field">
                                        <label for="c-tel">Teléfono</label>
                                        <input
                                            id="c-tel"
                                            bind:value={editingClient.Telefono}
                                            placeholder="Teléfono"
                                        />
                                    </div>
                                </div>

                                <div class="field">
                                    <label for="c-addr">Dirección</label>
                                    <input
                                        id="c-addr"
                                        bind:value={editingClient.Direccion}
                                        placeholder="Dirección"
                                    />
                                </div>
                            </div>

                            <div class="sidebar-footer" style="margin-top: 16px; border-top: 1px solid var(--border-subtle); padding-top: 16px;">
                                <div class="flex-row">
                                    <button
                                        class="btn-primary full-width"
                                        style="flex:1"
                                        on:click={handleSaveClient}>Guardar</button
                                    >
                                    {#if editingClient.ID}
                                        <button
                                            class="btn-secondary"
                                            on:click={() =>
                                                (editingClient = {
                                                    ID: "",
                                                    TipoID: "05",
                                                    Nombre: "",
                                                    Direccion: "",
                                                    Email: "",
                                                    Telefono: "",
                                                })}>Cancelar</button
                                        >
                                    {/if}
                                </div>
                            </div>
                        </div>

                        <!-- LISTA CLIENTES -->
                        <div class="card no-padding" style="display: flex; flex-direction: column;">
                            <div class="linear-grid" style="flex: 1; overflow: hidden; display: flex; flex-direction: column;">
                                <div class="linear-header grid-columns-clients">
                                    <div class="cell">ID</div>
                                    <div class="cell">Nombre</div>
                                    <div class="cell">Email</div>
                                    <div class="cell text-center">Acciones</div>
                                </div>
                                <div class="rows-container" style="overflow-y: auto; flex: 1;">
                                    {#each clients as c}
                                        <div class="linear-row grid-columns-clients">
                                            <div class="cell mono text-secondary">{c.ID}</div>
                                            <div class="cell">{c.Nombre}</div>
                                            <div class="cell text-secondary" style="font-size: 12px;">{c.Email}</div>
                                            <div class="cell text-center" style="gap: 4px; justify-content: center;">
                                                <button class="btn-icon-mini" on:click={() => (editingClient = { ...c })}>✏️</button>
                                                <button class="btn-icon-mini danger" on:click={() => handleDeleteClient(c.ID)}>🗑️</button>
                                            </div>
                                        </div>
                                    {/each}
                                </div>
                                {#if clients.length === 0}
                                    <div class="empty-state">
                                        <div class="empty-state-icon">👥</div>
                                        <div class="empty-state-text">
                                            No hay clientes registrados. Agrega uno nuevo desde el formulario.
                                        </div>
                                    </div>
                                {/if}
                            </div>
                        </div>
                    </div>
                </div>
            {:else if activeTab === "config"}
                <div class="panel" in:fade={{ duration: 200 }}>
                    <h1>Configuración</h1>
                    <div class="config-grid">
                        <div class="card">
                            <h3>🏢 Datos de Empresa</h3>
                            <div class="field">
                                <label for="cfg-ruc">RUC</label><input id="cfg-ruc"
                                    bind:value={config.RUC}
                                />
                            </div>
                            <div class="field">
                                <label for="cfg-razon">Razón Social</label><input id="cfg-razon"
                                    bind:value={config.RazonSocial}
                                />
                            </div>
                            <div class="field">
                                <label for="cfg-nombre">Nombre Comercial</label><input id="cfg-nombre"
                                    bind:value={config.NombreComercial}
                                />
                            </div>
                            <div class="field">
                                <label for="cfg-dir">Dirección</label><input id="cfg-dir"
                                    bind:value={config.Direccion}
                                />
                            </div>
                            <div class="grid col-2-tight">
                                <div class="field">
                                    <label for="cfg-estab">Estab (001)</label><input id="cfg-estab"
                                        bind:value={config.Estab}
                                        maxlength="3"
                                        class="text-center"
                                    />
                                </div>
                                <div class="field">
                                    <label for="cfg-ptoemi">PtoEmi (001)</label><input id="cfg-ptoemi"
                                        bind:value={config.PtoEmi}
                                        maxlength="3"
                                        class="text-center"
                                    />
                                </div>
                            </div>
                            <div class="field checkbox-field mt-1">
                                <label for="cfg-obligado"
                                    style="display: flex; align-items: center; gap: 8px; cursor: pointer;"
                                >
                                    <input id="cfg-obligado"
                                        type="checkbox"
                                        bind:checked={config.Obligado}
                                        style="width: auto;"
                                    />
                                    Obligado a Contabilidad
                                </label>
                            </div>
                        </div>

                        <div class="card">
                            <h3>🔐 Firma y Marca</h3>
                            <div class="field">
                                <label for="cfg-logo">Logo</label>
                                <div id="cfg-logo"
                                    class="flex-row"
                                    style="background: var(--bg-surface-hover); padding: 10px; border-radius: 8px;"
                                >
                                    {#if config.LogoPath}
                                        <img
                                            src={config.LogoPath.startsWith("/")
                                                ? "file://" + config.LogoPath
                                                : config.LogoPath}
                                            alt="Logo"
                                            style="width: 50px; height: 50px; object-fit: contain; background: white; border-radius: 4px;"
                                        />
                                    {/if}
                                    <button
                                        class="btn-secondary"
                                        on:click={async () => {
                                            const path =
                                                await SelectAndSaveLogo();
                                            if (
                                                path &&
                                                !path.startsWith("Error")
                                            )
                                                config.LogoPath = path;
                                        }}>📷 Subir Logo</button
                                    >
                                </div>
                            </div>

                            <div class="field">
                                <label for="cfg-p12">Archivo .p12</label>
                                <div class="input-group">
                                    <input id="cfg-p12"
                                        bind:value={config.P12Path}
                                        readonly
                                    />
                                    <button
                                        class="btn-secondary"
                                        on:click={handleSelectCert}>📂</button
                                    >
                                </div>
                            </div>
                            <div class="field">
                                <label for="cfg-pass">Contraseña Firma</label>
                                <div class="input-group">
                                    {#if showPassword}
                                        <input id="cfg-pass"
                                            type="text"
                                            bind:value={config.P12Password}
                                        />
                                    {:else}
                                        <input id="cfg-pass"
                                            type="password"
                                            bind:value={config.P12Password}
                                        />
                                    {/if}
                                    <button
                                        class="btn-secondary"
                                        on:click={() =>
                                            (showPassword = !showPassword)}
                                        >👁</button
                                    >
                                </div>
                            </div>
                        </div>

                        <div class="card">
                            <h3>📧 Correo (SMTP)</h3>
                            <div class="flex-row mb-1" style="gap: 12px;">
                                <button
                                    class="btn-secondary"
                                    on:click={() => {
                                        config.SMTPHost = "smtp.gmail.com";
                                        config.SMTPPort = 587;
                                    }}>Gmail</button
                                >
                                <button
                                    class="btn-secondary"
                                    on:click={() => {
                                        config.SMTPHost = "smtp.office365.com";
                                        config.SMTPPort = 587;
                                    }}>Outlook</button
                                >
                            </div>
                            <div class="field">
                                <label for="smtp-host">Host</label><input id="smtp-host"
                                    bind:value={config.SMTPHost}
                                />
                            </div>
                            <div class="grid col-2-tight">
                                <div class="field">
                                    <label for="smtp-port">Puerto</label><input id="smtp-port"
                                        type="number"
                                        bind:value={config.SMTPPort}
                                    />
                                </div>
                                <div class="field">
                                    <label for="smtp-user">Usuario</label><input id="smtp-user"
                                        bind:value={config.SMTPUser}
                                    />
                                </div>
                            </div>
                            <div class="field">
                                <label for="smtp-pass">Contraseña</label><input id="smtp-pass"
                                    type="password"
                                    bind:value={config.SMTPPassword}
                                />
                            </div>
                            <button
                                class="btn-secondary mt-1 full-width"
                                on:click={handleTestSMTP}>Test Conexión</button
                            >
                        </div>
                    </div>

                    <div
                        class="flex-row"
                        style="justify-content: flex-end; padding-top: 1rem; border-top: 1px solid var(--border-subtle);"
                    >
                        <button
                            class="btn-primary"
                            on:click={handleSaveConfig}
                            disabled={loading}
                        >
                            Guardar Configuración
                        </button>
                    </div>
                </div>
            {:else if activeTab === "sync"}
                <div class="panel full-height" in:fade={{ duration: 200 }}>
                    <div class="header-row">
                        <h1>Auditoría y Sync</h1>
                        <div class="header-actions">
                            <button
                                class="btn-secondary"
                                on:click={refreshSyncLogs}>Refrescar</button
                            >
                            <button class="btn-primary" on:click={handleSyncNow}
                                >Sincronizar SRI</button
                            >
                        </div>
                    </div>

                    <div
                        class="activity-grid"
                        style="display: grid; grid-template-rows: auto 1fr; gap: 24px; height: calc(100vh - 200px);"
                    >
                        <!-- MAIL LOGS -->
                        <div
                            class="card no-padding"
                            style="height: 300px; display: flex; flex-direction: column;"
                        >
                            <div
                                class="linear-grid"
                                style="flex: 1; overflow: hidden; display: flex; flex-direction: column;"
                            >
                                <div class="linear-header grid-columns-mail">
                                    <div class="cell">Estado</div>
                                    <div class="cell">Fecha</div>
                                    <div class="cell">Destinatario</div>
                                    <div class="cell">Mensaje</div>
                                </div>
                                <div
                                    class="rows-container"
                                    style="overflow-y: auto; flex: 1;"
                                >
                                    {#each mailLogs as ml}
                                        <div
                                            class="linear-row grid-columns-mail"
                                        >
                                            <div class="cell">
                                                <span class="badge {ml.estado}"
                                                    >{ml.estado}</span
                                                >
                                            </div>
                                            <div
                                                class="cell mono text-secondary"
                                            >
                                                {ml.fecha}
                                            </div>
                                            <div class="cell">{ml.email}</div>
                                            <div
                                                class="cell text-muted"
                                                style="white-space: nowrap; overflow: hidden; text-overflow: ellipsis;"
                                            >
                                                {ml.mensaje}
                                            </div>
                                        </div>
                                    {/each}
                                </div>
                            </div>
                        </div>

                        <!-- SYNC LOGS -->
                        <div
                            class="sync-container"
                            style="display: grid; grid-template-columns: 350px 1fr; gap: 24px; height: 100%; overflow: hidden;"
                        >
                            <div
                                class="card no-padding scrollable"
                                style="display: flex; flex-direction: column;"
                            >
                                <div
                                    class="linear-header"
                                    style="padding: 12px 16px; border-bottom: 1px solid var(--border-subtle); background: var(--bg-surface);"
                                >
                                    <h3 style="margin: 0; font-size: 14px;">
                                        Eventos de Sincronización
                                    </h3>
                                </div>
                                <div style="flex: 1; overflow-y: auto;">
                                    {#each syncLogs as log}
                                        <div
                                            role="button"
                                            tabindex="0"
                                            class="log-item {selectedLog === log
                                                ? 'active'
                                                : ''}"
                                            on:click={() => (selectedLog = log)}
                                            on:keydown={(e) => (e.key === 'Enter' || e.key === ' ') && (selectedLog = log)}
                                        >
                                            <div class="flex-row space-between">
                                                <span
                                                    class="mono text-muted"
                                                    style="font-size: 11px;"
                                                    >{log.timestamp}</span
                                                >
                                                <span
                                                    class="badge {log.status}"
                                                    style="transform: scale(0.9);"
                                                    >{log.status}</span
                                                >
                                            </div>
                                            <div
                                                style="font-weight: 500; font-size: 13px; margin-top: 4px; color: var(--text-primary);"
                                            >
                                                {log.action}
                                            </div>
                                        </div>
                                    {/each}
                                </div>
                            </div>

                            <div
                                class="card"
                                style="overflow: hidden; display: flex; flex-direction: column;"
                            >
                                {#if selectedLog}
                                    <h3 style="margin-bottom: 16px;">
                                        Detalle del Evento
                                    </h3>
                                    <div style="flex: 1; overflow-y: auto;">
                                        <div
                                            style="margin-bottom: 8px; font-size: 12px; color: var(--text-secondary);"
                                        >
                                            Request Payload
                                        </div>
                                        <pre class="code-block">{formatJSON(
                                                selectedLog.request,
                                            )}</pre>

                                        <div
                                            style="margin-bottom: 8px; margin-top: 16px; font-size: 12px; color: var(--text-secondary);"
                                        >
                                            Response Payload
                                        </div>
                                        <pre class="code-block">{formatJSON(
                                                selectedLog.response,
                                            )}</pre>
                                    </div>
                                {:else}
                                    <div class="empty-state">
                                        Selecciona un log para ver detalles
                                    </div>
                                {/if}
                            </div>
                        </div>
                    </div>

                    <style>
                        .grid-columns-mail {
                            grid-template-columns: 100px 160px 200px 1fr;
                        }
                        .log-item {
                            padding: 12px 16px;
                            border-bottom: 1px solid var(--border-subtle);
                            cursor: pointer;
                            transition: background 0.2s;
                        }
                        .log-item:hover {
                            background: var(--bg-hover);
                        }
                        .log-item.active {
                            background: rgba(0, 255, 148, 0.05);
                            border-left: 3px solid var(--accent-mint);
                        }
                        .code-block {
                            font-family: "SF Mono", monospace;
                            font-size: 11px;
                            background: var(--bg-hover);
                            padding: 12px;
                            border-radius: 6px;
                            color: var(--text-primary);
                            border: 1px solid var(--border-subtle);
                            white-space: pre-wrap;
                        }
                    </style>
                </div>
            {/if}

            {#if activeTab === "backups"}
                <div class="panel" in:fade={{ duration: 200 }}>
                    <div class="header-row">
                        <h1>Copias de Seguridad</h1>
                        <button
                            class="btn-primary"
                            on:click={handleCreateBackup}
                            disabled={loading}>Crear Respaldo</button
                        >
                    </div>

                    <div class="card mb-4">
                        <h3>Ruta de Almacenamiento</h3>
                        <div class="flex-row" style="gap: 12px;">
                            <input
                                value={config.StoragePath
                                    ? config.StoragePath + "/Backups"
                                    : "Default"}
                                readonly
                                style="flex: 1; font-family: monospace; background: var(--bg-hover);"
                            />
                            <button
                                class="btn-secondary"
                                on:click={handleSelectStorage}
                                >📂 Cambiar</button
                            >
                        </div>
                    </div>

                    <div class="card no-padding">
                        <div class="linear-grid" style="min-height: 300px;">
                            <div class="linear-header grid-columns-backups">
                                <div class="cell">Archivo</div>
                                <div class="cell">Fecha</div>
                                <div class="cell">Tamaño</div>
                                <div class="cell">Ruta Completa</div>
                            </div>
                            <div class="rows-container">
                                {#each backups as b}
                                    <div
                                        class="linear-row grid-columns-backups"
                                    >
                                        <div
                                            class="cell font-medium"
                                            style="color: var(--accent-mint);"
                                        >
                                            {b.name}
                                        </div>
                                        <div class="cell mono text-secondary">
                                            {b.date}
                                        </div>
                                        <div class="cell">{b.size}</div>
                                        <div
                                            class="cell text-muted"
                                            style="font-size: 11px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;"
                                            title={b.path}
                                        >
                                            {b.path}
                                        </div>
                                    </div>
                                {/each}
                                {#if backups.length === 0}
                                    <div class="empty-state">
                                        No hay respaldos disponibles
                                    </div>
                                {/if}
                            </div>
                        </div>
                    </div>

                    <style>
                        .grid-columns-backups {
                            grid-template-columns: 200px 150px 100px 1fr;
                        }
                    </style>
                </div>
            {/if}
        </section>

        {#if confirmationModal.show}
            <div class="modal-overlay" transition:fade>
                <div class="modal-card" in:fly={{ y: 20 }}>
                    <h3>{confirmationModal.title}</h3>
                    <p style="color: var(--text-muted); margin: 1rem 0;">
                        {confirmationModal.message}
                    </p>
                    <div class="modal-actions">
                        <button
                            class="btn-secondary"
                            on:click={() => (confirmationModal.show = false)}
                            >Cancelar</button
                        ><button class="btn-danger" on:click={handleConfirm}
                            >Confirmar</button
                        >
                    </div>
                </div>
            </div>
        {/if}

        <Wizard show={showWizard} on:complete={onWizardComplete} />
    {/if}

    <div class="toast-container">
        {#each toasts as t (t.id)}
            <div
                class="toast-card {t.type}"
                transition:fly={{ x: 20, duration: 300 }}
                animate:flip={{ duration: 300 }}
            >
                {t.message}
            </div>
        {/each}
    </div>
</main>

<style>
    /* All styles delegated to style.css */
</style>
