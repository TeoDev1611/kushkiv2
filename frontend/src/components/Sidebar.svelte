<script>
  export let activeTab;
  import { createEventDispatcher } from 'svelte';
  import { fly } from 'svelte/transition';
  const dispatch = createEventDispatcher();

  let collapsed = false;

  function setTab(tab) {
      dispatch('change', tab);
  }

  function toggleSidebar() {
      collapsed = !collapsed;
  }
</script>

<aside class:collapsed={collapsed}>
    <div class="header-group">
        <div class="logo">
            {#if !collapsed}
                <div class="logo-text" in:fly={{ x: -10, duration: 200 }}>
                    <h2>Kushki</h2>
                    <span>Facturador</span>
                </div>
            {/if}
        </div>
        <button class="toggle-btn" on:click={toggleSidebar} title={collapsed ? "Expandir" : "Contraer"}>
            {collapsed ? '☰' : '«'}
        </button>
    </div>
    
    <nav>
      <button class:active={activeTab === 'dashboard'} on:click={() => setTab('dashboard')} title="Resumen">
        <span class="icon">📊</span> 
        {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 100 }}>Resumen</span>{/if}
      </button>
      <button class:active={activeTab === 'invoice'} on:click={() => setTab('invoice')} title="Nueva Factura">
        <span class="icon">📄</span> 
        {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 100 }}>Nueva Factura</span>{/if}
      </button>
      <button class:active={activeTab === 'history'} on:click={() => setTab('history')} title="Historial">
        <span class="icon">🕒</span> 
        {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 100 }}>Historial</span>{/if}
      </button>
      <button class:active={activeTab === 'products'} on:click={() => setTab('products')} title="Productos">
        <span class="icon">📦</span> 
        {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 100 }}>Productos</span>{/if}
      </button>
      <button class:active={activeTab === 'clients'} on:click={() => setTab('clients')} title="Clientes">
        <span class="icon">👥</span> 
        {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 100 }}>Clientes</span>{/if}
      </button>
      <button class:active={activeTab === 'sync'} on:click={() => setTab('sync')} title="Actividad">
        <span class="icon">📈</span> 
        {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 100 }}>Actividad</span>{/if}
      </button>
      <button class:active={activeTab === 'backups'} on:click={() => setTab('backups')} title="Respaldos">
        <span class="icon">💾</span> 
        {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 100 }}>Respaldos</span>{/if}
      </button>
      <button class:active={activeTab === 'config'} on:click={() => setTab('config')} title="Configuración">
        <span class="icon">⚙️</span> 
        {#if !collapsed}<span class="label" in:fly={{ x: -5, duration: 100 }}>Configuración</span>{/if}
      </button>
    </nav>
</aside>

<style>
  aside {
    width: 240px; background-color: #0B0F19;
    border-right: 1px solid rgba(255,255,255,0.05);
    display: flex; flex-direction: column; padding: 1rem;
    transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    overflow: hidden;
  }
  aside.collapsed {
      width: 70px;
      padding: 1rem 0.5rem;
  }

  .header-group { display: flex; align-items: center; justify-content: space-between; margin-bottom: 2rem; min-height: 40px; }
  .logo { display: flex; align-items: center; overflow: hidden; white-space: nowrap; }
  .logo h2 { margin: 0; color: #fff; font-size: 1.2rem; }
  .logo span { color: #34d399; font-size: 0.7rem; font-weight: bold; display: block; }
  
  .toggle-btn { background: transparent; border: none; color: #64748b; font-size: 1.2rem; cursor: pointer; padding: 4px; border-radius: 4px; }
  .toggle-btn:hover { color: white; background: rgba(255,255,255,0.1); }
  aside.collapsed .header-group { justify-content: center; }

  nav { display: flex; flex-direction: column; gap: 4px; }
  nav button {
    display: flex; align-items: center; gap: 12px; width: 100%;
    padding: 12px 16px; background: transparent; border: none;
    border-left: 3px solid transparent;
    color: #94a3b8; cursor: pointer; border-radius: 0 8px 8px 0; 
    transition: all 0.2s; white-space: nowrap;
    overflow: hidden;
  }
  nav button:hover { background: rgba(255,255,255,0.03); color: #fff; }
  nav button.active { 
      background: rgba(52, 211, 153, 0.1); 
      color: #34d399; 
      font-weight: 600; 
      border-left: 3px solid #34d399;
  }
  
  aside.collapsed nav button {
      justify-content: center;
      padding: 12px;
      gap: 0;
  }
  .icon { font-size: 1.2rem; min-width: 24px; text-align: center; }
</style>