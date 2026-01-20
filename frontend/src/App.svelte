<script>
  import { 
    CreateInvoice, GetEmisorConfig, SaveEmisorConfig, SelectCertificate, SelectStoragePath,
    GetProducts, SearchProducts, SaveProduct, DeleteProduct,
    GetClients, SearchClients, SaveClient, DeleteClient,
    GetNextSecuencial, GetDashboardStats, GetFacturasPaginated, OpenFacturaPDF,
    OpenInvoiceFolder, OpenInvoiceXML, ExportSalesExcel, ResendInvoiceEmail,
    GetTopProducts, GetSyncLogs, TriggerSyncManual
  } from '../wailsjs/go/main/App.js'
  import { EventsOn } from '../wailsjs/runtime/runtime.js'
  import { onMount } from 'svelte'
  import { fade, fly, slide } from 'svelte/transition'
  import { cubicOut } from 'svelte/easing'
  import { flip } from 'svelte/animate'
  import Sidebar from './components/Sidebar.svelte'
  import Wizard from './components/Wizard.svelte'

  // --- UTILIDADES ---
  function debounce(func, wait) {
    let timeout;
    return function(...args) {
      const context = this;
      clearTimeout(timeout);
      timeout = setTimeout(() => func.apply(context, args), wait);
    };
  }

  // --- ESTADO GLOBAL ---
  let activeTab = 'dashboard' 
  let loading = false
  let initialLoading = true // Estado para la Splash Screen
  let showWizard = false
  let contentArea
  let toast = { show: false, message: '', type: 'success' } 
  let confirmationModal = { show: false, title: '', message: '', onConfirm: null }

  // --- CONFIGURACIÓN Y FACTURACIÓN ---
  let config = { RUC: '', RazonSocial: '', NombreComercial: '', Direccion: '', P12Path: '', P12Password: '', Ambiente: 1, Estab: '001', PtoEmi: '001', Obligado: false, SMTPHost: '', SMTPUser: '', SMTPPass: '', StoragePath: '' }
  let invoice = { secuencial: '000000001', clienteID: '', clienteNombre: '', clienteDireccion: '', clienteEmail: '', clienteTelefono: '', formaPago: '01', items: [] }
  let newItem = { codigo: "", nombre: "", cantidad: 1, precio: 0, codigoIVA: "4", porcentajeIVA: 15 }
  
  let clientSearch = ''
  let productSearch = ''
  let showPassword = false

  // --- GESTIÓN DE PRODUCTOS ---
  let products = []
  let editingProduct = { SKU: '', Name: '', Price: 0, Stock: 0, TaxCode: "4", TaxPercentage: 15 }

  // --- GESTIÓN DE CLIENTES ---
  let clients = []
  let editingClient = { ID: '', TipoID: '05', Nombre: '', Direccion: '', Email: '', Telefono: '' }

  // --- ESTADO SINCRONIZACIÓN ---
  let syncLogs = []
  let selectedLog = null

  // --- DASHBOARD & HISTORIAL ---
  let stats = { totalVentas: 0, totalFacturas: 0, pendientes: 0, sriOnline: true, salesTrend: [] }
  let topProducts = []
  let historial = []
  let historialSearch = ''
  let dateRange = { start: '', end: '' }

  // --- COMPUTED ---
  $: subtotal = invoice.items.reduce((sum, item) => sum + (item.cantidad * item.precio), 0)
  $: iva = invoice.items.reduce((sum, item) => sum + (item.cantidad * item.precio * (item.porcentajeIVA / 100)), 0)
  $: total = subtotal + iva
  $: filteredHistorial = (historial || []).filter(f => {
      const cliente = f.cliente || '';
      const secuencial = f.secuencial || '';
      return cliente.toLowerCase().includes(historialSearch.toLowerCase()) || 
             secuencial.includes(historialSearch)
  })

  // Resetear scroll al cambiar de pestaña
  $: if (activeTab && contentArea) contentArea.scrollTop = 0

  onMount(() => {
      // Cargar datos y manejar Splash Screen
      loadData().then(() => {
          setTimeout(() => {
              initialLoading = false;
          }, 2000); // 2 segundos de splash screen para efecto visual
      });
      
      EventsOn("toast-notification", (data) => {
          showToast(data.message, data.type)
          if (data.type === 'success') refreshDashboard()
      })
  })

  async function loadData() {
    loading = true
    try {
      const cfg = await GetEmisorConfig()
      if (!cfg || !cfg.RUC) {
          showWizard = true
          // No return here, allow loading other lists if possible, but wizard takes precedence visually
      } else {
          config = cfg
      }
      
      // Cargar el resto aunque no haya config completa (para que la UI no rompa)
      const [prods, cli, seq] = await Promise.all([
          GetProducts(),
          GetClients(),
          GetNextSecuencial()
      ])
      
      products = prods || []
      clients = cli || []
      invoice.secuencial = seq || "000000001"
      
      await refreshDashboard()
    } catch (e) { 
      console.error("Error arrancando:", e) 
      showToast("Error de conexión o arranque: " + e, 'error')
    } finally {
      loading = false 
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
  let currentPage = 1
  let pageSize = 10
  let totalFacturas = 0
  $: totalPages = Math.ceil(totalFacturas / pageSize)

  // Inicializar fechas
  const d = new Date()
  dateRange.start = new Date(d.getFullYear(), d.getMonth(), 1).toISOString().split('T')[0]
  dateRange.end = new Date().toISOString().split('T')[0]

  async function refreshDashboard() {
      stats = await GetDashboardStats()
      await loadFacturasPage()
      topProducts = await GetTopProducts() || []
  }

  async function loadFacturasPage() {
      const res = await GetFacturasPaginated(currentPage, pageSize)
      historial = res.data || []
      totalFacturas = res.total
  }

  function changePage(newPage) {
      if (newPage >= 1 && newPage <= totalPages) {
          currentPage = newPage
          loadFacturasPage()
      }
  }

  function getBarHeight(val) {
      if (!stats.salesTrend || stats.salesTrend.length === 0) return 0;
      const max = Math.max(...stats.salesTrend.map(d => d.total));
      if (max === 0) return 0;
      return (val / max) * 100;
  }

  async function handleExportExcel() {
      loading = true
      try {
          const res = await ExportSalesExcel(dateRange.start, dateRange.end)
          showToast(res, res.includes('Error') ? 'error' : 'success')
      } catch (err) {
          showToast("Error al exportar: " + err, 'error')
      } finally {
          loading = false
      }
  }

  function showToast(message, type = 'success') {
    toast = { show: true, message, type }
    setTimeout(() => toast.show = false, 4000)
  }

  function showConfirmation(title, message, callback) {
      confirmationModal = { show: true, title, message, onConfirm: callback }
  }

  function handleConfirm() {
      if (confirmationModal.onConfirm) confirmationModal.onConfirm()
      confirmationModal.show = false
  }

  async function handleSelectCert() {
    const path = await SelectCertificate()
    if (path) config.P12Path = path
  }

  async function handleSelectStorage() {
    const path = await SelectStoragePath()
    if (path) config.StoragePath = path
  }

  async function handleSaveConfig() {
    loading = true
    try {
      config.Ambiente = parseInt(config.Ambiente)
      const res = await SaveEmisorConfig(config)
      if (res.startsWith("Error")) showToast(res, 'error')
      else showToast(res, 'success')
    } catch (err) {
      showToast(err, 'error')
    } finally {
      loading = false
    }
  }

  async function handleOpenPDF(clave) {
    showToast(await OpenFacturaPDF(clave), 'success')
  }

  async function handleOpenXML(clave) {
    const res = await OpenInvoiceXML(clave)
    showToast(res, res.includes('Error') ? 'error' : 'success')
  }

  async function handleOpenFolder(clave) {
    const res = await OpenInvoiceFolder(clave)
    showToast(res, res.includes('Error') ? 'error' : 'success')
  }

  async function handleResendEmail(clave) {
      const res = await ResendInvoiceEmail(clave)
      showToast(res, res.includes('Error') ? 'error' : 'success')
  }

  async function handleSaveProduct() {
    const res = await SaveProduct(editingProduct)
    showToast(res, res.includes('Error') ? 'error' : 'success')
    editingProduct = { SKU: '', Name: '', Price: 0, Stock: 0, TaxCode: "4", TaxPercentage: 15 }
    products = await GetProducts()
  }

  async function handleDeleteProduct(sku) {
    showConfirmation('Eliminar Producto', '¿Estás seguro de que deseas eliminar este producto?', async () => {
        const res = await DeleteProduct(sku)
        showToast(res, res.includes('Error') ? 'error' : 'success')
        products = await GetProducts()
    })
  }

  function selectProduct(p) {
    const taxCode = p.TaxCode || (p.TaxPercentage > 0 ? "4" : "0");
    const taxPerc = p.TaxPercentage !== undefined ? p.TaxPercentage : (p.TaxCode == "4" ? 15 : 0);
    newItem = { codigo: p.SKU, nombre: p.Name, cantidad: 1, precio: p.Price, codigoIVA: taxCode, porcentajeIVA: taxPerc }
    productSearch = ""
  }

  async function handleSaveClient() {
    const res = await SaveClient(editingClient)
    showToast(res, res.includes('Error') ? 'error' : 'success')
    editingClient = { ID: '', TipoID: '05', Nombre: '', Direccion: '', Email: '', Telefono: '' }
    clients = await GetClients()
  }

  async function handleDeleteClient(id) {
    showConfirmation('Eliminar Cliente', '¿Estás seguro de que deseas eliminar este cliente?', async () => {
        const res = await DeleteClient(id)
        showToast(res, res.includes('Error') ? 'error' : 'success')
        clients = await GetClients()
    })
  }

  function selectClient(c) {
    invoice.clienteID = c.ID
    invoice.clienteNombre = c.Nombre
    invoice.clienteDireccion = c.Direccion
    invoice.clienteEmail = c.Email || ''
    invoice.clienteTelefono = c.Telefono || ''
    clientSearch = ""
  }

  async function refreshSyncLogs() {
      loading = true
      try { syncLogs = await GetSyncLogs() } catch (e) { console.error(e) } finally { loading = false }
  }

  async function handleSyncNow() {
      const msg = await TriggerSyncManual()
      showToast(msg, 'info')
      setTimeout(refreshSyncLogs, 1000)
  }

  function formatJSON(str) {
      if (!str) return '';
      try { const obj = JSON.parse(str); return JSON.stringify(obj, null, 2); } catch (e) { return str; }
  }

  function addItem() {
    if (!newItem.nombre || newItem.precio <= 0) return
    invoice.items = [...invoice.items, { ...newItem }]
    newItem = { codigo: "", nombre: "", cantidad: 1, precio: 0, codigoIVA: "4", porcentajeIVA: 15 }
  }

  function removeItem(index) {
    invoice.items = invoice.items.filter((_, i) => i !== index)
  }

  async function handleEmit() {
    if (!invoice.clienteID || invoice.items.length === 0) {
      showToast("Ingrese cliente y al menos un producto", 'error')
      return
    }
    loading = true
    try {
      const dto = {
        secuencial: invoice.secuencial,
        clienteID: invoice.clienteID,
        clienteNombre: invoice.clienteNombre,
        clienteDireccion: invoice.clienteDireccion,
        clienteEmail: invoice.clienteEmail,
        clienteTelefono: invoice.clienteTelefono,
        formaPago: invoice.formaPago,
        items: invoice.items.map(i => ({
            codigo: i.codigo, nombre: i.nombre, cantidad: parseFloat(i.cantidad),
            precio: parseFloat(i.precio), codigoIVA: i.codigoIVA.toString(), porcentajeIVA: parseFloat(i.porcentajeIVA)
        }))
      }
      const result = await CreateInvoice(dto)
      if (result.startsWith("Éxito")) {
        showToast(result, 'success')
        invoice.secuencial = (parseInt(invoice.secuencial) + 1).toString().padStart(9, '0')
        invoice.items = []
        await refreshDashboard()
      } else {
        showToast(result, 'error')
      }
    } catch (err) {
      showToast("Error crítico: " + err, 'error')
    } finally {
      loading = false
    }
  }
</script>

{#if initialLoading}
    <div class="splash-screen" out:fade={{ duration: 800 }}>
        <div class="splash-content">
            <div class="logo-container-splash">
                <h1 class="splash-title">KUSHKI</h1>
                <div class="splash-line"></div>
                <p class="splash-subtitle">FACTURACIÓN ELECTRÓNICA</p>
            </div>
            <div class="loader-bar">
                <div class="loader-progress"></div>
            </div>
        </div>
    </div>
{/if}

<main>
  <Sidebar {activeTab} on:change={(e) => activeTab = e.detail} />

  <section class="content" bind:this={contentArea}>
    {#if activeTab === 'dashboard'}
      <div class="panel dashboard-layout" in:fade={{ duration: 100 }}>
        <div class="header-row">
            <div>
                <h1>Hola, {config.RazonSocial || 'Emisor'} 👋</h1>
                <p class="subtitle">Aquí tienes el resumen de tu facturación.</p>
            </div>
            <button class="btn-icon-soft" on:click={refreshDashboard} title="Actualizar">🔄</button>
        </div>
        
        <div class="kpi-row">
            <div class="kpi-card glass">
                <div class="kpi-icon mint">💰</div>
                <div class="kpi-content">
                    <div class="title">Ventas Totales</div>
                    <div class="value gradient-text">${stats.totalVentas.toFixed(2)}</div>
                </div>
            </div>
            <div class="kpi-card glass">
                <div class="kpi-icon blue">📄</div>
                <div class="kpi-content">
                    <div class="title">Facturas</div>
                    <div class="value">{stats.totalFacturas}</div>
                </div>
            </div>
            <div class="kpi-card glass {stats.pendientes > 0 ? 'warning-border' : ''}">
                <div class="kpi-icon orange">⚠️</div>
                <div class="kpi-content">
                    <div class="title">Pendientes</div>
                    <div class="value">{stats.pendientes}</div>
                </div>
            </div>
            <div class="kpi-card glass">
                <div class="kpi-icon {stats.sriOnline ? 'green' : 'red'}">
                    {stats.sriOnline ? '🌐' : '🔌'}
                </div>
                <div class="kpi-content">
                    <div class="title">Estado SRI</div>
                    <div class="status-badge {stats.sriOnline ? 'online' : 'offline'}">
                        {stats.sriOnline ? 'ONLINE' : 'OFFLINE'}
                    </div>
                </div>
            </div>
        </div>

        <div class="section-chart glass">
            <div class="chart-header">
                <h3>Tendencia de Ventas (7 Días)</h3>
            </div>
            <div class="chart-container">
                {#if stats.salesTrend && stats.salesTrend.length > 0}
                    <div class="bar-chart">
                        {#each stats.salesTrend as day}
                            <div class="bar-group" title="{day.date}: ${day.total}">
                                <div class="bar-value" style="bottom: {getBarHeight(day.total) + 5}%; opacity: {day.total > 0 ? 1 : 0}">${day.total > 0 ? day.total : ''}</div>
                                <div class="bar" style="height: {getBarHeight(day.total)}%;"></div>
                                <div class="label">{day.date.slice(5)}</div>
                            </div>
                        {/each}
                    </div>
                {:else}
                    <div class="empty-state">Sin datos suficientes para graficar</div>
                {/if}
            </div>
        </div>

        <div class="grid col-2">
            <div class="card glass compact">
                <h3>🏆 Top Productos</h3>
                <div class="mini-table">
                    {#each topProducts as p}
                        <div class="mini-row">
                            <div class="prod-name">{p.name}</div>
                            <div class="prod-stat">{p.quantity} un.</div>
                            <div class="prod-price">${p.total.toFixed(2)}</div>
                        </div>
                    {/each}
                    {#if topProducts.length === 0}
                        <div class="empty-state-small">Sin ventas aún</div>
                    {/if}
                </div>
            </div>

            <div class="card glass compact">
                <div class="header-row">
                    <h3>⚡ Actividad Reciente</h3>
                    <button class="link-btn" on:click={() => activeTab = 'history'}>Ver todo →</button>
                </div>
                <div class="mini-table">
                    {#each (historial || []).slice(0, 5) as f}
                        <div class="mini-row">
                            <span class="status-dot {f.estado}"></span>
                            <div class="col-main">
                                <span class="secuencial">{f.secuencial}</span>
                                <span class="client">{f.cliente}</span>
                            </div>
                            <div class="amount">${f.total.toFixed(2)}</div>
                        </div>
                    {/each}
                </div>
            </div>
        </div>
      </div>

    {:else if activeTab === 'history'}
      <div class="panel" in:fade={{ duration: 100 }}>
        <div class="header-row">
            <h1>Historial de Transacciones</h1>
            <div class="header-actions">
                <button class="btn-secondary" on:click={handleExportExcel}>📊 Exportar Excel</button>
                <button class="btn-secondary" on:click={loadFacturasPage}>🔄 Refrescar</button>
            </div>
        </div>

        <div class="card no-padding">
             <!-- Toolbar Unificada -->
             <div class="toolbar">
                <div class="date-group">
                    <span class="label-tiny">Desde</span>
                    <input type="date" bind:value={dateRange.start} class="input-tiny" />
                    <span class="label-tiny">Hasta</span>
                    <input type="date" bind:value={dateRange.end} class="input-tiny" />
                </div>
                <div class="search-wrapper">
                    <input class="search-input-compact" bind:value={historialSearch} placeholder="🔍 Buscar por cliente, RUC o secuencial..." />
                </div>
            </div>

            <div class="table-container dense-table">
                <table>
                    <thead>
                        <tr>
                            <th>Estado</th>
                            <th>Secuencial</th>
                            <th>Fecha</th>
                            <th>Cliente</th>
                            <th class="text-right">Total</th>
                            <th class="text-center">Acciones</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each filteredHistorial as f}
                            <tr>
                                <td><span class="badge {f.estado}">{f.estado}</span></td>
                                <td class="mono">{f.secuencial}</td>
                                <td class="mono">{f.fecha}</td>
                                <td>{f.cliente}</td>
                                <td class="bold text-right">${f.total.toFixed(2)}</td>
                                <td class="text-center">
                                    {#if f.tienePDF}
                                        <button class="btn-icon-mini" title="Ver PDF" on:click={() => handleOpenPDF(f.claveAcceso)}>📄</button>
                                        <button class="btn-icon-mini" title="Email" on:click={() => handleResendEmail(f.claveAcceso)}>✉️</button>
                                    {/if}
                                    <button class="btn-icon-mini" title="XML" on:click={() => handleOpenXML(f.claveAcceso)}>🌐</button>
                                    <button class="btn-icon-mini" title="Carpeta" on:click={() => handleOpenFolder(f.claveAcceso)}>📂</button>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
                {#if filteredHistorial.length === 0}
                    <div class="empty-state">No se encontraron transacciones</div>
                {/if}
            </div>
            <div class="pagination-controls compact">
                <button class="btn-secondary small" disabled={currentPage === 1} on:click={() => changePage(currentPage - 1)}>Anterior</button>
                <span>Página {currentPage} de {totalPages || 1} ({totalFacturas} items)</span>
                <button class="btn-secondary small" disabled={currentPage === totalPages} on:click={() => changePage(currentPage + 1)}>Siguiente</button>
            </div>
        </div>
      </div>

    {:else if activeTab === 'invoice'}
      <div class="panel" in:fade={{ duration: 100 }}>
        <div class="header-row">
          <h1>Emitir Factura</h1>
          <div class="badge-pill">Ambiente: {config.Ambiente === 1 ? 'PRUEBAS' : 'PRODUCCIÓN'}</div>
        </div>

        <!-- SECCIÓN 1: DATOS DEL CLIENTE -->
        <div class="card">
          <div class="section-header">
             <h3>👤 Datos del Cliente</h3>
             <div class="secuencial-box">
                <span class="label">Nro. Factura</span>
                <span class="number">{config.Estab}-{config.PtoEmi}-{invoice.secuencial}</span>
             </div>
          </div>
          
          <div class="search-container mb-1">
            <input 
                bind:value={clientSearch} 
                on:input={(e) => handleClientSearchInput(e)}
                placeholder="🔍 Buscar cliente por nombre o RUC/CI..." 
                class="full-search-input" 
            />
            {#if clientSearch.length > 0}
                <div class="dropdown">
                    {#each clients.filter(c => c.Nombre.toLowerCase().includes(clientSearch.toLowerCase()) || c.ID.includes(clientSearch)) as c}
                        <button on:click={() => selectClient(c)}>{c.Nombre} ({c.ID})</button>
                    {/each}
                </div>
            {/if}
          </div>

          <div class="client-grid">
            <div class="field">
                <label>Identificación</label>
                <input bind:value={invoice.clienteID} placeholder="1712345678" />
            </div>
            <div class="field span-2">
                <label>Razón Social</label>
                <input bind:value={invoice.clienteNombre} placeholder="Nombre Cliente" />
            </div>
            <div class="field">
                <label>Email</label>
                <input bind:value={invoice.clienteEmail} placeholder="cliente@email.com" type="email" />
            </div>
            <div class="field span-2">
                <label>Dirección</label>
                <input bind:value={invoice.clienteDireccion} placeholder="Dirección completa" />
            </div>
            <div class="field">
                <label>Teléfono</label>
                <input bind:value={invoice.clienteTelefono} placeholder="0991234567" />
            </div>
            <div class="field">
                <label>Forma de Pago</label>
                <select bind:value={invoice.formaPago} class="dark-select">
                    <option value="01">Efectivo (01)</option>
                    <option value="16">Tarjeta de Débito (16)</option>
                    <option value="19">Tarjeta de Crédito (19)</option>
                    <option value="20">Otros con SF (20)</option>
                </select>
            </div>
          </div>
        </div>

        <!-- SECCIÓN 2: PRODUCTOS -->
        <div class="card">
          <div class="section-header">
            <h3>📦 Detalle de Productos</h3>
          </div>

          <div class="add-product-bar">
             <div class="search-box-product">
                <input bind:value={productSearch} placeholder="🔍 Buscar producto..." />
                {#if productSearch.length > 0}
                    <div class="dropdown">
                        {#each products.filter(p => p.Name.toLowerCase().includes(productSearch.toLowerCase()) || p.SKU.includes(productSearch)) as p}
                            <button on:click={() => selectProduct(p)}>{p.Name} - ${p.Price}</button>
                        {/each}
                    </div>
                {/if}
             </div>
             
             <!-- Campos de edición manual rápida -->
             <input class="input-code" bind:value={newItem.codigo} placeholder="Código" />
             <input class="input-desc" bind:value={newItem.nombre} placeholder="Descripción" />
             <input class="input-qty" type="number" bind:value={newItem.cantidad} min="1" placeholder="Cant" />
             <input class="input-price" type="number" step="0.01" bind:value={newItem.precio} placeholder="Precio" />
             <select bind:value={newItem.codigoIVA} on:change={() => {
                const map = { "0": 0, "2": 12, "4": 15, "5": 5 };
                newItem.porcentajeIVA = map[newItem.codigoIVA];
             }} class="dark-select input-tax">
                <option value="4">15%</option>
                <option value="5">5%</option>
                <option value="0">0%</option>
                <option value="2">12%</option>
             </select>
             <button class="btn-add" on:click={addItem}>Añadir</button>
          </div>

          <div class="table-container invoice-table">
            <table>
              <thead>
                <tr>
                  <th width="10%">Cant</th>
                  <th width="40%">Descripción</th>
                  <th width="15%" class="text-right">P. Unit</th>
                  <th width="10%" class="text-center">IVA</th>
                  <th width="15%" class="text-right">Total</th>
                  <th width="5%"></th>
                </tr>
              </thead>
              <tbody>
                {#each invoice.items as item, i (item)}
                  <tr>
                    <td class="text-center">{item.cantidad}</td>
                    <td>{item.nombre}</td>
                    <td class="text-right">${item.precio.toFixed(2)}</td>
                    <td class="text-center"><span class="tag-tax">{item.porcentajeIVA}%</span></td>
                    <td class="text-right bold">${(item.cantidad * item.precio * (1 + item.porcentajeIVA/100)).toFixed(2)}</td>
                    <td><button class="btn-remove" on:click={() => removeItem(i)}>×</button></td>
                  </tr>
                {/each}
                {#if invoice.items.length === 0}
                    <tr><td colspan="6" class="empty-row">Agrega productos para comenzar la factura</td></tr>
                {/if}
              </tbody>
            </table>
          </div>
        </div>

        <!-- SECCIÓN 3: TOTALES Y ACCIONES -->
        <div class="footer-panel">
            <div class="notes-area">
                <label>Notas Adicionales (Opcional)</label>
                <textarea placeholder="Ej: Entregar en recepción..."></textarea>
            </div>
            <div class="totals-area">
                <div class="total-row">
                    <span>Subtotal Sin Impuestos</span>
                    <span class="mono">${subtotal.toFixed(2)}</span>
                </div>
                <div class="total-row">
                    <span>IVA (15% / 5%)</span>
                    <span class="mono">${iva.toFixed(2)}</span>
                </div>
                <div class="total-row grand-total">
                    <span>TOTAL A PAGAR</span>
                    <span class="mono">${total.toFixed(2)}</span>
                </div>
                <button class="btn-primary large full-width" on:click={handleEmit} disabled={loading}>
                    {loading ? 'Procesando...' : '🖋️ Firmar y Emitir Factura'}
                </button>
            </div>
        </div>
      </div>

    {:else if activeTab === 'products'}
      <div class="panel" in:fade={{ duration: 100 }}>
        <h1>Inventario de Productos</h1>
        <div class="master-detail-layout">
            <!-- Formulario Lateral (Sticky) -->
            <div class="sidebar-form card">
                <h3>{editingProduct.SKU ? 'Editar' : 'Nuevo'} Producto</h3>
                <div class="form-stack">
                    <div class="field">
                        <label>SKU / Código</label>
                        <input bind:value={editingProduct.SKU} placeholder="COD-001" />
                    </div>
                    <div class="field">
                        <label>Nombre</label>
                        <input bind:value={editingProduct.Name} placeholder="Nombre del producto" />
                    </div>
                    <div class="grid col-2-tight">
                        <div class="field">
                            <label>Precio</label>
                            <input type="number" step="0.01" bind:value={editingProduct.Price} />
                        </div>
                        <div class="field">
                            <label>Impuesto</label>
                            <select bind:value={editingProduct.TaxCode} on:change={() => {
                                const map = { "0": 0, "2": 12, "4": 15, "5": 5 };
                                editingProduct.TaxPercentage = map[editingProduct.TaxCode];
                            }} class="dark-select">
                                <option value="4">15%</option>
                                <option value="5">5%</option>
                                <option value="0">0%</option>
                                <option value="2">12%</option>
                            </select>
                        </div>
                    </div>
                </div>
                <div class="form-actions">
                     <button class="btn-primary full-width" on:click={handleSaveProduct}>Guardar</button>
                     {#if editingProduct.SKU}
                        <button class="btn-secondary full-width" on:click={() => editingProduct = { SKU: '', Name: '', Price: 0, Stock: 0, TaxCode: "4", TaxPercentage: 15 }}>Cancelar</button>
                     {/if}
                </div>
            </div>

            <!-- Tabla Principal -->
            <div class="card no-padding scrollable-area">
                <div class="table-header">
                    <h3>Lista de Productos</h3>
                    <span class="badge-pill">{products.length} items</span>
                </div>
                <div class="table-body-scroll">
                    <table class="dense-table">
                        <thead><tr><th>SKU</th><th>Nombre</th><th class="text-right">Precio</th><th class="text-center">Acciones</th></tr></thead>
                        <tbody>
                            {#each products as p}
                                <tr>
                                    <td class="mono">{p.SKU}</td>
                                    <td>{p.Name}</td>
                                    <td class="text-right bold text-mint">${p.Price.toFixed(2)}</td>
                                    <td class="text-center">
                                        <button class="btn-icon-mini" on:click={() => editingProduct = {...p}}>✏️</button>
                                        <button class="btn-icon-mini danger" on:click={() => handleDeleteProduct(p.SKU)}>🗑️</button>
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
      </div>

    {:else if activeTab === 'clients'}
      <div class="panel" in:fade={{ duration: 100 }}>
        <h1>Directorio de Clientes</h1>
        <div class="master-detail-layout">
            <div class="sidebar-form card">
                <h3>{editingClient.ID ? 'Editar' : 'Nuevo'} Cliente</h3>
                <div class="form-stack">
                    <div class="field">
                        <label>Identificación (RUC/CI)</label>
                        <input bind:value={editingClient.ID} placeholder="1712345678" />
                    </div>
                    <div class="field">
                        <label>Nombre / Razón Social</label>
                        <input bind:value={editingClient.Nombre} />
                    </div>
                    <div class="field">
                        <label>Email</label>
                        <input bind:value={editingClient.Email} type="email" placeholder="Para facturación" />
                    </div>
                    <div class="field">
                        <label>Dirección</label>
                        <input bind:value={editingClient.Direccion} />
                    </div>
                    <div class="field">
                        <label>Teléfono</label>
                        <input bind:value={editingClient.Telefono} />
                    </div>
                </div>
                <div class="form-actions">
                    <button class="btn-primary full-width" on:click={handleSaveClient}>Guardar</button>
                    {#if editingClient.ID}
                        <button class="btn-secondary full-width" on:click={() => editingClient = { ID: '', TipoID: '05', Nombre: '', Direccion: '', Email: '', Telefono: '' }}>Cancelar</button>
                    {/if}
                </div>
            </div>

            <div class="card no-padding scrollable-area">
                 <div class="table-header">
                    <h3>Mis Clientes</h3>
                    <span class="badge-pill">{clients.length} registros</span>
                </div>
                <div class="table-body-scroll">
                    <table class="dense-table">
                        <thead><tr><th>ID</th><th>Nombre</th><th>Email</th><th class="text-center">Acciones</th></tr></thead>
                        <tbody>
                            {#each clients as c}
                                <tr>
                                    <td class="mono">{c.ID}</td>
                                    <td>{c.Nombre}</td>
                                    <td class="text-small">{c.Email}</td>
                                    <td class="text-center">
                                        <button class="btn-icon-mini" on:click={() => editingClient = {...c}}>✏️</button>
                                        <button class="btn-icon-mini danger" on:click={() => handleDeleteClient(c.ID)}>🗑️</button>
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
      </div>

    {:else if activeTab === 'config'}
      <div class="panel" in:fade={{ duration: 100 }}>
        <h1>Configuración de Emisor</h1>
        <div class="config-grid">
            <!-- SECCIÓN EMPRESA -->
            <div class="card config-card">
                <h3>🏢 Datos de Empresa</h3>
                <div class="form-stack">
                    <div class="field">
                        <label>RUC</label>
                        <input bind:value={config.RUC} />
                    </div>
                    <div class="field">
                        <label>Razón Social</label>
                        <input bind:value={config.RazonSocial} />
                    </div>
                    <div class="field">
                        <label>Nombre Comercial</label>
                        <input bind:value={config.NombreComercial} placeholder="Opcional" />
                    </div>
                    <div class="field">
                        <label>Dirección Matriz</label>
                        <input bind:value={config.Direccion} />
                    </div>
                    <div class="grid col-2-tight">
                        <div class="field">
                            <label>Estab (001)</label>
                            <input bind:value={config.Estab} maxlength="3" class="text-center" />
                        </div>
                        <div class="field">
                            <label>PtoEmi (001)</label>
                            <input bind:value={config.PtoEmi} maxlength="3" class="text-center" />
                        </div>
                    </div>
                    <div class="field checkbox-field">
                        <label><input type="checkbox" bind:checked={config.Obligado} /> Obligado a Contabilidad</label>
                    </div>
                </div>
            </div>

            <!-- SECCIÓN FIRMA Y RUTAS -->
            <div class="card config-card">
                <h3>🔐 Firma Electrónica</h3>
                <div class="form-stack">
                    <div class="field">
                        <label>Archivo .p12</label>
                        <div class="input-group">
                            <input bind:value={config.P12Path} readonly class="text-small" />
                            <button class="btn-secondary compact" on:click={handleSelectCert}>📂</button>
                        </div>
                    </div>
                    <div class="field">
                        <label>Contraseña Firma</label>
                        <div class="input-group">
                            {#if showPassword}
                                <input type="text" bind:value={config.P12Password} />
                            {:else}
                                <input type="password" bind:value={config.P12Password} />
                            {/if}
                            <button class="btn-secondary compact" on:click={() => showPassword = !showPassword}>
                                {showPassword ? '🙈' : '👁️'}
                            </button>
                        </div>
                    </div>
                    <div class="field mt-1">
                        <label>Carpeta de Guardado (XML/PDF)</label>
                        <div class="input-group">
                            <input bind:value={config.StoragePath} readonly class="text-small" placeholder="Por defecto" />
                            <button class="btn-secondary compact" on:click={handleSelectStorage}>📂</button>
                        </div>
                        <small class="hint">Se crearán carpetas por Año/Mes automáticamente.</small>
                    </div>
                </div>
            </div>

            <!-- SECCIÓN CORREO -->
            <div class="card config-card">
                <h3>📧 Servidor de Correo (SMTP)</h3>
                <div class="form-stack">
                    <div class="field">
                        <label>Host SMTP (Ej: smtp.gmail.com:587)</label>
                        <input bind:value={config.SMTPHost} />
                    </div>
                    <div class="field">
                        <label>Usuario / Email</label>
                        <input bind:value={config.SMTPUser} />
                    </div>
                    <div class="field">
                        <label>Contraseña de Aplicación</label>
                        <input type="password" bind:value={config.SMTPPass} />
                    </div>
                    <div class="alert-box">
                        <small>Nota: Para Gmail, usa una "Contraseña de Aplicación" si tienes 2FA activado.</small>
                    </div>
                </div>
            </div>
        </div>
        
        <div class="actions-right">
            <button class="btn-primary large" on:click={handleSaveConfig} disabled={loading}>
              Guardar Toda la Configuración
            </button>
        </div>
      </div>

    {:else if activeTab === 'sync'}
      <div class="panel full-height" in:fade={{ duration: 100 }}>
        <div class="header-row">
            <h1>Estado de Sincronización</h1>
            <div class="header-actions"><button class="btn-secondary" on:click={refreshSyncLogs}>🔄 Refrescar Logs</button><button class="btn-primary" on:click={handleSyncNow}>🚀 Sincronizar Ahora</button></div>
        </div>
        <div class="sync-container">
            <div class="log-list card scrollable">
                <h3>Historial</h3>
                {#if syncLogs.length === 0}<div class="empty-state">Sin registros</div>
                {:else}
                    {#each syncLogs as log (log.id)}
                        <!-- svelte-ignore a11y-click-events-have-key-events -->
                        <div class="log-item {selectedLog === log ? 'active' : ''}" on:click={() => selectedLog = log}>
                            <div class="log-header"><span class="timestamp">{log.timestamp}</span><span class="badge {log.status}">{log.status}</span></div>
                            <div class="log-title">{log.action}</div><div class="log-detail">{log.detail}</div>
                        </div>
                    {/each}
                {/if}
            </div>
            <div class="log-detail-panel card">
                {#if selectedLog}
                    <h3>Detalle</h3>
                    <div class="detail-grid">
                        <div class="detail-section"><label>Petición</label><div class="code-block"><pre>{formatJSON(selectedLog.request)}</pre></div></div>
                        <div class="detail-section"><label>Respuesta</label><div class="code-block"><pre>{formatJSON(selectedLog.response)}</pre></div></div>
                    </div>
                {:else}<div class="empty-selection"><p>Selecciona un evento</p></div>{/if}
            </div>
        </div>
      </div>
    {/if}
  </section>

  {#if confirmationModal.show}
    <div class="modal-overlay" transition:fade>
      <div class="modal-card" in:fly={{ y: 20 }}>
        <h3>{confirmationModal.title}</h3><p>{confirmationModal.message}</p>
        <div class="modal-actions"><button class="btn-secondary" on:click={() => confirmationModal.show = false}>Cancelar</button><button class="btn-danger" on:click={handleConfirm}>Confirmar</button></div>
      </div>
    </div>
  {/if}
  
  <Wizard show={showWizard} on:complete={onWizardComplete} />
</main>

<style>
  :global(:root) {
    --bg-obsidian: #0B0F19; --bg-surface: #161e31; --bg-glass: rgba(11, 15, 25, 0.7); --border-subtle: rgba(255, 255, 255, 0.08);
    --primary-mint: #34d399; --mint-dark: #059669; --accent-indigo: #6366f1; --status-error: #EF3340; --status-warning: #fbbf24; --status-info: #60a5fa;
    --radius-md: 12px; --radius-lg: 16px; --elevation-1: 0 1px 3px rgba(0,0,0,0.12), 0 1px 2px rgba(0,0,0,0.24); --elevation-2: 0 3px 6px rgba(0,0,0,0.16), 0 3px 6px rgba(0,0,0,0.23);
    --font-family: 'Outfit', sans-serif;
  }
  :global(body) { margin: 0; font-family: var(--font-family); background-color: var(--bg-obsidian); color: #ecf0f1; overflow: hidden; }
  main { display: flex; height: 100vh; }
  .content { flex: 1; padding: 2rem; overflow-y: auto; background-color: #0B0F19; }
  .panel { max-width: 1100px; margin: 0 auto; display: flex; flex-direction: column; gap: 1.5rem; }
  .card, .kpi-card { 
      background: var(--bg-surface); padding: 1.5rem; border-radius: var(--radius-lg); 
      border: 1px solid var(--border-subtle); 
      /* Optimización: Sin sombras pesadas por defecto */
      box-shadow: none;
      transition: transform 0.2s ease, border-color 0.2s ease; 
  }
  .card:hover, .kpi-card:hover { 
      transform: translateY(-2px); 
      border-color: rgba(255,255,255,0.15);
      /* Sombra ligera solo en hover */
      box-shadow: 0 4px 12px rgba(0,0,0,0.3);
  }
  input, select { 
      background: rgba(255, 255, 255, 0.03); 
      border: 1px solid transparent; 
      border-bottom: 1px solid #475569; 
      color: white; padding: 14px 16px; 
      border-radius: 8px 8px 0 0; 
      font-family: inherit; width: 100%; box-sizing: border-box; 
      transition: background-color 0.2s, border-color 0.2s; 
  }
  /* Custom Dark Select */
  .dark-select {
      appearance: none;
      -webkit-appearance: none;
      -moz-appearance: none;
      background-color: rgba(255, 255, 255, 0.03);
      background-image: url("data:image/svg+xml;charset=UTF-8,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='white' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
      background-repeat: no-repeat;
      background-position: right 12px center;
      background-size: 16px;
      color: white;
      padding-right: 40px; /* Espacio para la flecha */
      border-bottom: 1px solid #475569;
      border-radius: 8px 8px 0 0;
  }
  .dark-select option {
      background-color: #1e293b; /* Fondo oscuro para las opciones desplegadas */
      color: white;
  }
  .dark-select:focus {
      background-color: rgba(255, 255, 255, 0.06);
      border-bottom-color: var(--primary-mint);
  }
  .btn-primary { 
      background: var(--primary-mint); color: #064e3b; border: none; 
      padding: 12px 24px; border-radius: var(--radius-md); font-weight: 700; cursor: pointer; 
      transition: transform 0.1s, background-color 0.2s; 
      box-shadow: 0 2px 4px rgba(0,0,0,0.2); 
  }
  .btn-primary:hover { background: var(--mint-dark); color: #fff; transform: translateY(-1px); }
  .btn-primary:disabled { background: #334155; color: #94a3b8; cursor: not-allowed; transform: none; box-shadow: none; }
  .btn-secondary { 
      background: rgba(255, 255, 255, 0.05); color: #fff; border: 1px solid var(--border-subtle); 
      padding: 10px 20px; border-radius: var(--radius-md); cursor: pointer; 
      transition: background-color 0.2s, border-color 0.2s; 
  }
  .btn-danger { background: var(--status-error); color: #fff; border: none; padding: 10px 20px; border-radius: var(--radius-md); cursor: pointer; transition: background 0.2s; }
  .btn-icon { background: var(--primary-mint); color: #000; width: 40px; height: 40px; border: none; border-radius: 12px; cursor: pointer; font-size: 1.5rem; transition: background 0.2s; }
  .btn-icon-small { background: transparent; border: none; cursor: pointer; font-size: 1.1rem; transition: transform 0.1s; }
  .text-danger { color: #ef4444; background: transparent; border: none; cursor: pointer; font-size: 1.2rem; }
  .btn-primary:active, .btn-secondary:active, .btn-danger:active { transform: scale(0.96); }
  .loading-overlay { position: fixed; inset: 0; background: rgba(11, 15, 25, 0.8); backdrop-filter: blur(8px); display: flex; flex-direction: column; align-items: center; justify-content: center; z-index: 200; }
  .spinner { width: 50px; height: 50px; border: 3px solid rgba(52, 211, 153, 0.1); border-top-color: var(--primary-mint); border-radius: 50%; animation: spin 0.8s cubic-bezier(0.68, -0.55, 0.27, 1.55) infinite; filter: drop-shadow(0 0 8px rgba(52, 211, 153, 0.5)); }
  @keyframes spin { to { transform: rotate(360deg); } }
  ::-webkit-scrollbar { width: 8px; height: 8px; } ::-webkit-scrollbar-track { background: #0B0F19; } ::-webkit-scrollbar-thumb { background: #334155; border-radius: 4px; }
  .kpi-card .title { font-size: 0.85rem; color: #94a3b8; margin-bottom: 0.5rem; } .kpi-card .value { font-size: 1.8rem; font-weight: bold; color: white; }
  .status-indicator { display: flex; align-items: center; gap: 8px; font-weight: bold; margin-top: 5px; }
  .dot { width: 10px; height: 10px; border-radius: 50%; } .dot.green { background: #34d399; } .dot.red { background: #ef4444; }
  .grid { display: grid; gap: 1rem; } .col-2 { grid-template-columns: 1fr 1fr; } .col-3 { grid-template-columns: 1fr 2fr 1fr; }
  .header-row { display: flex; justify-content: space-between; align-items: center; } .header-actions { display: flex; gap: 10px; align-items: center; }
  .date-selector { display: flex; align-items: center; gap: 5px; background: #161e31; padding: 4px 10px; border-radius: 8px; border: 1px solid #334155; }
  .date-selector input { background: transparent; border: none; padding: 4px; color: white; }
  .field { display: flex; flex-direction: column; gap: 6px; } .input-group { display: flex; gap: 8px; }
  .search-input { background: #0f1523; border: none; padding: 8px 12px; color: white; border-radius: 6px; }
  .search-input.full-width { width: 100%; font-size: 1rem; padding: 12px; background: rgba(0,0,0,0.2); border: 1px solid rgba(255,255,255,0.1); }
  .dropdown { position: absolute; top: 100%; left: 0; right: 0; background: #1e293b; border: 1px solid #34d399; border-radius: 0 0 8px 8px; z-index: 10; max-height: 200px; overflow-y: auto; }
  .dropdown button { width: 100%; padding: 10px; background: transparent; border: none; color: white; text-align: left; cursor: pointer; }
  .chart-container { height: 250px; display: flex; align-items: flex-end; justify-content: center; padding-top: 20px; }
  .bar-chart { display: flex; align-items: flex-end; gap: 15px; height: 100%; width: 100%; }
  .bar-group { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: flex-end; position: relative; }
  .bar { width: 100%; border-radius: 4px 4px 0 0; background: linear-gradient(to top, #059669, #34d399); transition: height 0.5s ease; }
  .add-item-row { display: flex; gap: 10px; background: #0f1523; padding: 10px; border-radius: 8px; align-items: center; }
  .grow { flex: 1; } .small { width: 80px; } .medium { width: 120px; }
  table { width: 100%; border-collapse: collapse; } th { text-align: left; color: #64748b; font-size: 0.8rem; padding: 10px; border-bottom: 1px solid #334155; } td { padding: 10px; border-bottom: 1px solid rgba(255,255,255,0.05); }
  .badge { padding: 4px 10px; border-radius: 12px; font-size: 0.7rem; font-weight: 700; text-transform: uppercase; }
  .badge.AUTORIZADO { background: rgba(52, 211, 153, 0.1); color: var(--primary-mint); }
  .modal-overlay { position: fixed; inset: 0; background: rgba(11, 15, 25, 0.7); display: flex; align-items: center; justify-content: center; z-index: 150; }
  .modal-card { background: #161e31; padding: 2rem; border-radius: 12px; width: 400px; }
  .dashboard-layout { display: flex; flex-direction: column; gap: 1.5rem; }
  .kpi-row { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 1rem; }
  .glass { 
      background: #161e31; 
      border: 1px solid rgba(255,255,255,0.05); 
      border-radius: 16px; padding: 1.2rem;
      transition: transform 0.2s, box-shadow 0.2s;
  }
  .kpi-card { display: flex; align-items: center; gap: 1rem; }
  .kpi-icon { width: 48px; height: 48px; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: 1.5rem; }
  .kpi-icon.mint { background: rgba(52, 211, 153, 0.15); color: #34d399; }
  .kpi-icon.blue { background: rgba(96, 165, 250, 0.15); color: #60a5fa; }
  .kpi-icon.orange { background: rgba(251, 191, 36, 0.15); color: #fbbf24; }
  .section-chart { padding: 1.5rem; height: 320px; display: flex; flex-direction: column; }
  .mini-table { display: flex; flex-direction: column; gap: 0.8rem; }
  .mini-row { display: flex; align-items: center; justify-content: space-between; padding: 0.8rem; background: rgba(255,255,255,0.02); border-radius: 8px; }
  .status-badge { font-size: 0.75rem; font-weight: 700; padding: 2px 8px; border-radius: 4px; }
  .status-badge.online { background: rgba(34, 197, 94, 0.2); color: #4ade80; }
  .status-badge.offline { background: rgba(239, 68, 68, 0.2); color: #f87171; }
  .gradient-text { background: linear-gradient(45deg, #34d399, #60a5fa); -webkit-background-clip: text; -webkit-text-fill-color: transparent; }
  .status-dot { width: 8px; height: 8px; border-radius: 50%; margin-right: 10px; }
  .status-dot.AUTORIZADO { background: #34d399; }
  .col-main { display: flex; flex-direction: column; flex: 1; }
  .secuencial { font-family: monospace; font-size: 0.8rem; color: #94a3b8; }
  .client { font-size: 0.9rem; font-weight: 500; }
  .amount { font-weight: 700; color: white; }
  .link-btn { background: none; border: none; color: #60a5fa; cursor: pointer; font-size: 0.85rem; }
  .link-btn:hover { text-decoration: underline; }
  .btn-icon-soft { background: rgba(255,255,255,0.05); border: none; width: 36px; height: 36px; border-radius: 50%; color: white; cursor: pointer; }
  
  /* --- INVOICE UI STYLES --- */
  .section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; border-bottom: 1px solid rgba(255,255,255,0.05); padding-bottom: 0.5rem; }
  .section-header h3 { margin: 0; font-size: 1rem; color: #e2e8f0; display: flex; align-items: center; gap: 8px; }
  
  .secuencial-box { background: rgba(52, 211, 153, 0.1); padding: 4px 12px; border-radius: 6px; border: 1px solid rgba(52, 211, 153, 0.2); text-align: right; }
  .secuencial-box .label { display: block; font-size: 0.65rem; color: #34d399; font-weight: 700; text-transform: uppercase; }
  .secuencial-box .number { font-family: monospace; font-size: 1rem; font-weight: 700; color: white; letter-spacing: 1px; }

  .client-grid { display: grid; grid-template-columns: 1fr 2fr 1fr; gap: 1rem; margin-top: 1rem; }
  .span-2 { grid-column: span 2; }
  .full-search-input { width: 100%; font-size: 1rem; padding: 12px; background: rgba(0,0,0,0.2); border: 1px solid rgba(255,255,255,0.1); border-radius: 8px; color: white; }

  .add-product-bar { display: flex; gap: 10px; background: #0f1523; padding: 12px; border-radius: 12px; align-items: center; border: 1px solid rgba(255,255,255,0.05); margin-bottom: 1rem; flex-wrap: wrap; }
  .search-box-product { flex: 2; position: relative; min-width: 200px; }
  .search-box-product input { background: transparent; border: none; border-bottom: 1px solid #475569; padding: 8px; width: 100%; color: white; }
  
  .input-code { width: 80px; }
  .input-desc { flex: 2; min-width: 150px; }
  .input-qty { width: 60px; text-align: center; }
  .input-price { width: 90px; text-align: right; }
  .input-tax { width: 80px; }
  
  .btn-add { background: var(--primary-mint); color: #064e3b; border: none; padding: 0 20px; height: 44px; border-radius: 8px; font-weight: 700; cursor: pointer; transition: transform 0.1s; }
  .btn-add:hover { transform: scale(1.02); background: var(--mint-dark); color: white; }
  .btn-remove { background: transparent; color: #ef4444; border: 1px solid rgba(239, 68, 68, 0.3); width: 28px; height: 28px; border-radius: 50%; cursor: pointer; display: flex; align-items: center; justify-content: center; font-size: 1.2rem; transition: background 0.2s; }
  .btn-remove:hover { background: rgba(239, 68, 68, 0.1); }

  .invoice-table th { background: rgba(255,255,255,0.02); font-weight: 600; color: #94a3b8; font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.5px; padding: 10px; text-align: left; }
  .invoice-table td { padding: 12px 10px; font-size: 0.9rem; border-bottom: 1px solid rgba(255,255,255,0.03); }
  .tag-tax { background: #334155; padding: 2px 6px; border-radius: 4px; font-size: 0.7rem; color: #cbd5e1; }
  .text-right { text-align: right; }
  .text-center { text-align: center; }
  .empty-row { text-align: center; padding: 2rem; color: #64748b; font-style: italic; }

  .footer-panel { display: flex; gap: 2rem; margin-top: 1rem; }
  .notes-area { flex: 1; }
  .notes-area textarea { width: 100%; height: 100px; background: rgba(255,255,255,0.03); border: 1px solid #475569; border-radius: 8px; padding: 10px; color: white; resize: none; font-family: inherit; }
  .totals-area { width: 320px; background: #0f1523; padding: 1.5rem; border-radius: 12px; border: 1px solid rgba(255,255,255,0.05); }
  .total-row { display: flex; justify-content: space-between; margin-bottom: 8px; font-size: 0.9rem; color: #94a3b8; }
  .total-row.grand-total { border-top: 1px solid #334155; padding-top: 12px; margin-top: 12px; font-size: 1.2rem; color: white; font-weight: 700; }
  .total-row.grand-total .mono { color: var(--primary-mint); }
  .full-width { width: 100%; margin-top: 1rem; }

  /* --- SPECIALIZED LAYOUTS --- */
  
  /* Toolbar (History) */
  .toolbar { display: flex; gap: 1rem; margin-bottom: 1rem; align-items: center; justify-content: space-between; padding: 0.5rem 1rem; border-bottom: 1px solid rgba(255,255,255,0.05); }
  .date-group { display: flex; align-items: center; gap: 8px; background: rgba(0,0,0,0.2); padding: 4px 10px; border-radius: 6px; border: 1px solid rgba(255,255,255,0.05); }
  .label-tiny { font-size: 0.75rem; color: #94a3b8; font-weight: 600; text-transform: uppercase; }
  .input-tiny { border: none; background: transparent; color: white; padding: 2px; font-size: 0.85rem; width: auto; }
  .search-wrapper { flex: 1; max-width: 400px; }
  .search-input-compact { width: 100%; background: rgba(0,0,0,0.2); border: 1px solid rgba(255,255,255,0.1); padding: 6px 10px; border-radius: 6px; color: white; font-size: 0.9rem; }
  .search-input-compact:focus { border-color: var(--primary-mint); }
  
  /* Dense Table Optimized */
  .dense-table th { padding: 12px 15px; font-size: 0.75rem; background: rgba(0,0,0,0.15); text-transform: uppercase; color: #94a3b8; font-weight: 600; text-align: left; }
  .dense-table td { padding: 12px 15px; font-size: 0.95rem; border-bottom: 1px solid rgba(255,255,255,0.03); color: #e2e8f0; }
  
  /* Config Action Spacing */
  .actions-right { display: flex; justify-content: flex-end; margin-top: 3rem; padding-top: 2rem; border-top: 1px solid rgba(255,255,255,0.05); }

  /* Sync Spacing */
  .sync-container { display: grid; grid-template-columns: 350px 1fr; gap: 2rem; flex: 1; overflow: hidden; min-height: 0; }
  .log-item { padding: 1.2rem 1.5rem; border-bottom: 1px solid rgba(255,255,255,0.05); cursor: pointer; transition: background 0.2s; display: flex; flex-direction: column; gap: 4px; }

  /* Restored styles */
  .btn-icon-mini { background: transparent; border: 1px solid rgba(255,255,255,0.1); width: 26px; height: 26px; border-radius: 4px; display: inline-flex; align-items: center; justify-content: center; font-size: 0.9rem; cursor: pointer; margin-right: 4px; transition: all 0.1s; color: #cbd5e1; }
  .btn-icon-mini:hover { background: rgba(255,255,255,0.1); color: white; transform: scale(1.1); }
  .btn-icon-mini.danger:hover { background: rgba(239, 68, 68, 0.2); border-color: #ef4444; }

  .master-detail-layout { display: grid; grid-template-columns: 320px 1fr; gap: 1.5rem; height: calc(100vh - 180px); }
  .sidebar-form { height: fit-content; position: sticky; top: 0; display: flex; flex-direction: column; gap: 1rem; }
  .sidebar-form h3 { margin: 0 0 1rem 0; font-size: 1.1rem; border-bottom: 1px solid rgba(255,255,255,0.05); padding-bottom: 0.5rem; }
  .form-stack { display: flex; flex-direction: column; gap: 12px; }
  .col-2-tight { grid-template-columns: 1fr 1fr; gap: 10px; }
  .scrollable-area { display: flex; flex-direction: column; overflow: hidden; height: 100%; }
  .table-header { padding: 10px 15px; background: rgba(0,0,0,0.15); display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid rgba(255,255,255,0.05); }
  .table-header h3 { margin: 0; font-size: 0.95rem; }
  .table-body-scroll { flex: 1; overflow-y: auto; }
  .text-mint { color: #34d399; }
  .text-small { font-size: 0.8rem; color: #94a3b8; }
  .no-padding { padding: 0 !important; overflow: hidden; }
  
  .config-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1.5rem; }
  .config-card { height: 100%; }
  .config-card h3 { font-size: 1rem; color: var(--primary-mint); margin-bottom: 1rem; text-transform: uppercase; letter-spacing: 0.5px; }
  .hint { display: block; margin-top: 4px; font-size: 0.75rem; color: #64748b; }
  .alert-box { background: rgba(251, 191, 36, 0.1); border: 1px solid rgba(251, 191, 36, 0.2); color: #fbbf24; padding: 8px; border-radius: 6px; font-size: 0.75rem; margin-top: 1rem; }
    .compact { padding: 4px 8px; font-size: 0.8rem; }
  
      /* --- SPLASH SCREEN --- */
      .splash-screen {
          position: fixed; inset: 0; background: #0B0F19; z-index: 9999;
          display: flex; align-items: center; justify-content: center;
          flex-direction: column;
      }    .splash-content {
        text-align: center; display: flex; flex-direction: column; align-items: center; gap: 2rem;
    }
    .splash-title {
        font-size: 4rem; font-weight: 200; color: white; margin: 0;
        letter-spacing: 12px;
        background: linear-gradient(to right, #fff, #94a3b8);
        -webkit-background-clip: text; -webkit-text-fill-color: transparent;
        animation: tracking-in 1.5s cubic-bezier(0.215, 0.610, 0.355, 1.000) both;
    }
    .splash-subtitle {
        color: var(--primary-mint); font-size: 0.8rem; letter-spacing: 4px; font-weight: 600;
        margin-top: 0.5rem; opacity: 0;
        animation: fade-in-up 1s ease-out 0.5s forwards;
    }
    .splash-line {
        width: 0; height: 1px; background: rgba(255,255,255,0.1); margin: 10px auto;
        animation: expand-width 1s ease-out 0.2s forwards;
    }
    
    .loader-bar {
        width: 200px; height: 2px; background: rgba(255,255,255,0.1); border-radius: 2px; overflow: hidden;
        margin-top: 2rem;
    }
    .loader-progress {
        height: 100%; background: var(--primary-mint); width: 0;
        animation: load-progress 2s cubic-bezier(0.23, 1, 0.32, 1) forwards;
        box-shadow: 0 0 10px var(--primary-mint);
    }
  
    @keyframes tracking-in {
        0% { letter-spacing: -0.5em; opacity: 0; }
        40% { opacity: 0.6; }
        100% { opacity: 1; }
    }
    @keyframes expand-width { to { width: 100px; } }
    @keyframes fade-in-up { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
    @keyframes load-progress { 0% { width: 0; } 50% { width: 70%; } 100% { width: 100%; } }
  
  </style>
