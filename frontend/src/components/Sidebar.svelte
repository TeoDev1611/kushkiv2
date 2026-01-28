<script>
  export let activeTab;
  import { createEventDispatcher } from "svelte";
  import { fly, fade } from "svelte/transition";
  import logo from "../assets/images/logo-universal.png"; 
  
  const dispatch = createEventDispatcher();
  let collapsed = false;

  function setTab(tab) {
    dispatch("change", tab);
  }

  function toggleSidebar() {
    collapsed = !collapsed;
  }
</script>

<aside class:collapsed>
  <div class="header-group">
    <div class="logo-area">
       {#if !collapsed}
        <div class="brand-container" in:fade={{ duration: 200 }}>
             <img src={logo} alt="Kushki" class="brand-logo" />
             <div class="brand-text">
                 <span class="brand-name">KUSHKI</span>
                 <span class="brand-badge">PRO</span>
             </div>
        </div>
       {:else}
        <img src={logo} alt="K" class="brand-logo-mini" in:fade />
       {/if}
    </div>
    <button
      class="toggle-btn"
      on:click={toggleSidebar}
      title={collapsed ? "Expandir" : "Contraer"}
    >
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        {#if collapsed}
            <polyline points="13 17 18 12 13 7"></polyline>
            <polyline points="6 17 11 12 6 7"></polyline>
        {:else}
            <polyline points="11 17 6 12 11 7"></polyline>
            <polyline points="18 17 13 12 18 7"></polyline>
        {/if}
      </svg>
    </button>
  </div>

  <nav>
    <div class="nav-section">PRINCIPAL</div>
    
    <button class:active={activeTab === "dashboard"} on:click={() => setTab("dashboard")} title="Resumen">
      <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="9" rx="1"></rect><rect x="14" y="3" width="7" height="5" rx="1"></rect><rect x="14" y="12" width="7" height="9" rx="1"></rect><rect x="3" y="16" width="7" height="5" rx="1"></rect></svg>
      {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 150 }}>Dashboard</span>{/if}
    </button>

    <button class:active={activeTab === "invoice"} on:click={() => setTab("invoice")} title="Nueva Factura">
       <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="12" y1="18" x2="12" y2="12"></line><line x1="9" y1="15" x2="15" y2="15"></line></svg>
      {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 150 }}>Emitir Factura</span>{/if}
    </button>
    
    <button class:active={activeTab === "quotations"} on:click={() => setTab("quotations")} title="Cotizaciones">
      <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
      {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 150 }}>Cotizaciones</span>{/if}
    </button>

    <div class="nav-section">GESTIÓN</div>

    <button class:active={activeTab === "history"} on:click={() => setTab("history")} title="Historial">
       <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
      {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 150 }}>Historial</span>{/if}
    </button>

    <button class:active={activeTab === "products"} on:click={() => setTab("products")} title="Productos">
      <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line></svg>
      {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 150 }}>Productos</span>{/if}
    </button>

    <button class:active={activeTab === "clients"} on:click={() => setTab("clients")} title="Clientes">
      <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path><circle cx="9" cy="7" r="4"></circle><path d="M23 21v-2a4 4 0 0 0-3-3.87"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path></svg>
      {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 150 }}>Clientes</span>{/if}
    </button>

    <div class="nav-section mt-spacer">SISTEMA</div>

    <button class:active={activeTab === "sync"} on:click={() => setTab("sync")} title="Actividad">
      <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg>
      {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 150 }}>Actividad</span>{/if}
    </button>

    <button class:active={activeTab === "backups"} on:click={() => setTab("backups")} title="Respaldos">
      <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path><polyline points="17 21 17 13 7 13 7 21"></polyline><polyline points="7 3 7 8 15 8"></polyline></svg>
      {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 150 }}>Respaldos</span>{/if}
    </button>

    <button class:active={activeTab === "config"} on:click={() => setTab("config")} title="Configuración">
      <svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path></svg>
      {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 150 }}>Ajustes</span>{/if}
    </button>
  </nav>
</aside>

<style>
  aside {
    width: 280px;
    background-color: var(--bg-panel); /* Uses global var */
    border-right: 1px solid var(--border-subtle);
    display: flex;
    flex-direction: column;
    padding: 24px 16px;
    transition: width 0.4s cubic-bezier(0.16, 1, 0.3, 1);
    overflow: hidden;
    height: 100vh;
    z-index: 50;
  }
  aside.collapsed {
    width: 80px;
    padding: 24px 12px;
  }

  /* BRANDING */
  .header-group {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 40px;
    padding: 0 4px;
    height: 48px;
  }
  .brand-container {
      display: flex;
      align-items: center;
      gap: 12px;
  }
  .brand-logo {
      height: 32px;
      width: auto;
  }
  .brand-logo-mini {
      height: 28px;
      width: auto;
      margin: 0 auto;
      display: block;
  }
  .brand-text {
      display: flex;
      flex-direction: column;
      line-height: 1;
  }
  .brand-name {
      font-weight: 800;
      font-size: 18px;
      letter-spacing: -0.02em;
      color: var(--text-primary);
  }
  .brand-badge {
      font-size: 9px;
      font-weight: 700;
      color: var(--accent-mint);
      letter-spacing: 0.15em;
      margin-top: 2px;
  }

  .toggle-btn {
    background: transparent;
    border: none;
    color: var(--text-secondary);
    cursor: pointer;
    width: 32px; height: 32px;
    border-radius: 8px;
    display: flex; align-items: center; justify-content: center;
    transition: all 0.2s;
  }
  .toggle-btn:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }

  aside.collapsed .header-group {
      justify-content: center;
      padding: 0;
  }
  aside.collapsed .toggle-btn {
      display: none; 
  }
  /* Better toggle for collapsed: Click logo */
  aside.collapsed .brand-logo-mini {
      cursor: pointer;
  }
  aside.collapsed .header-group {
      cursor: pointer;
  }

  /* NAVIGATION */
  nav {
    display: flex;
    flex-direction: column;
    gap: 6px;
    flex: 1;
  }
  .nav-section {
      font-size: 10px;
      font-weight: 700;
      color: var(--text-tertiary);
      text-transform: uppercase;
      letter-spacing: 0.1em;
      margin: 16px 0 8px 12px;
  }
  aside.collapsed .nav-section { display: none; }

  nav button {
    display: flex;
    align-items: center;
    gap: 14px;
    width: 100%;
    padding: 12px 14px;
    background: transparent;
    border: none;
    color: var(--text-secondary);
    cursor: pointer;
    border-radius: 12px; /* Smooth */
    transition: all 0.2s cubic-bezier(0.2, 0, 0, 1);
    white-space: nowrap;
    font-size: 14px;
    font-weight: 500;
    position: relative;
  }

  nav button:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
    transform: translateX(4px);
  }

  nav button.active {
    background: linear-gradient(90deg, rgba(0, 230, 137, 0.1) 0%, transparent 100%);
    color: var(--accent-mint);
    font-weight: 600;
  }
  
  nav button.active::before {
      content: '';
      position: absolute;
      left: 0; top: 10px; bottom: 10px;
      width: 3px;
      background: var(--accent-mint);
      border-radius: 0 4px 4px 0;
      box-shadow: 0 0 12px var(--accent-mint);
  }

  aside.collapsed nav button {
    justify-content: center;
    padding: 14px;
    gap: 0;
  }
  aside.collapsed nav button:hover {
      transform: none;
      background: var(--bg-hover);
  }
  aside.collapsed nav button.active::before {
      left: 0;
  }

  .icon {
    width: 20px;
    height: 20px;
    min-width: 20px;
  }

  .mt-spacer { margin-top: 32px; }
</style>
