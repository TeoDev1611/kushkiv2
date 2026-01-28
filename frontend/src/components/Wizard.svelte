<script>
  import { createEventDispatcher } from 'svelte';
  import { SelectCertificate, SelectStoragePath, SaveEmisorConfig, SelectAndSaveLogo } from '../../wailsjs/go/main/App.js';
  
  export let show = false;
  const dispatch = createEventDispatcher();

  let step = 1;
  let config = {
      RUC: '',
      RazonSocial: '',
      Direccion: '',
      P12Path: '',
      P12Password: '',
      StoragePath: '',
      Ambiente: 1,
      Estab: '001',
      PtoEmi: '001',
      Obligado: false,
      ContribuyenteRimpe: '',
      AgenteRetencion: '',
      LogoPath: '',
      SMTPHost: '',
      SMTPPort: 587,
      SMTPUser: '',
      SMTPPassword: ''
  };
  let loading = false;
  let showPassword = false;
  let showSMTPPassword = false;

  async function handleSelectCert() {
      const path = await SelectCertificate();
      if (path) config.P12Path = path;
  }

  async function handleSelectStorage() {
      const path = await SelectStoragePath();
      if (path) config.StoragePath = path;
  }

  async function handleSelectLogo() {
      const path = await SelectAndSaveLogo();
      if (path && !path.startsWith("Error")) {
          config.LogoPath = path;
      } else if (path) {
          alert(path);
      }
  }

  async function finishSetup() {
      if (!config.RUC || !config.RazonSocial || !config.P12Path || !config.P12Password) {
          alert("Por favor complete los campos obligatorios (*)");
          return;
      }

      loading = true;
      try {
          config.Ambiente = parseInt(config.Ambiente);
          config.SMTPPort = parseInt(config.SMTPPort);
          const res = await SaveEmisorConfig(config);
          if (res.startsWith("Error")) {
              alert(res);
          } else {
              dispatch('complete');
          }
      } catch (err) {
          alert("Error crítico: " + err);
      } finally {
          loading = false;
      }
  }
</script>

{#if show}
<div class="wizard-overlay">
  <div class="wizard-card">
      <div class="wizard-header">
          <h2>🚀 Configuración Inicial</h2>
          <p>Bienvenido a Kushki Facturador. Configuremos su entorno.</p>
      </div>

      <div class="wizard-steps">
          <div class="step {step >= 1 ? 'active' : ''}">1. Empresa</div>
          <div class="step {step >= 2 ? 'active' : ''}">2. Firma</div>
          <div class="step {step >= 3 ? 'active' : ''}">3. Almacenamiento</div>
          <div class="step {step >= 4 ? 'active' : ''}">4. Correo</div>
      </div>

      <div class="wizard-content">
          {#if step === 1}
              <div class="field">
                  <label for="w-ruc">RUC *</label>
                  <input id="w-ruc" bind:value={config.RUC} placeholder="1790000000001" />
              </div>
              <div class="field">
                  <label for="w-razon">Razón Social *</label>
                  <input id="w-razon" bind:value={config.RazonSocial} placeholder="Mi Empresa S.A." />
              </div>
              <div class="field">
                  <label for="w-dir">Dirección Matriz</label>
                  <input id="w-dir" bind:value={config.Direccion} placeholder="Av. Principal 123" />
              </div>
              <div class="field">
                  <label for="w-agente">Agente de Retención</label>
                  <input id="w-agente" bind:value={config.AgenteRetencion} placeholder="Ej: Resolución Nro. 1" />
              </div>
              <div class="field">
                  <label for="w-logo">Logo de Empresa (Opcional)</label>
                  <div class="input-group">
                      <input id="w-logo" bind:value={config.LogoPath} readonly placeholder="Seleccione su logo..." />
                      <button class="btn-secondary" on:click={handleSelectLogo}>📷</button>
                  </div>
              </div>
          {:else if step === 2}
              <div class="field">
                  <label for="w-p12">Archivo de Firma (.p12) *</label>
                  <div class="input-group">
                      <input id="w-p12" bind:value={config.P12Path} readonly placeholder="Seleccione archivo..." />
                      <button class="btn-secondary" on:click={handleSelectCert}>📂</button>
                  </div>
              </div>
              <div class="field">
                  <label for="w-pass">Contraseña de Firma *</label>
                  <div class="input-group">
                      {#if showPassword}
                          <input id="w-pass" type="text" bind:value={config.P12Password} />
                      {:else}
                          <input id="w-pass" type="password" bind:value={config.P12Password} />
                      {/if}
                      <button class="btn-secondary" on:click={() => showPassword = !showPassword}>
                          {showPassword ? '🙈' : '👁️'}
                      </button>
                  </div>
              </div>
          {:else if step === 3}
              <div class="field">
                  <label for="w-store">Carpeta de Facturas</label>
                  <p class="hint">Donde se guardarán los XMLs y PDFs generados.</p>
                  <div class="input-group">
                      <input id="w-store" bind:value={config.StoragePath} readonly placeholder="Por defecto: Carpeta del usuario" />
                      <button class="btn-secondary" on:click={handleSelectStorage}>📂</button>
                  </div>
              </div>
              <div class="field checkbox">
                  <label>
                      <input type="checkbox" bind:checked={config.Obligado} /> Obligado a llevar contabilidad
                  </label>
              </div>
          {:else if step === 4}
              <div class="field">
                  <label for="w-smtp">Servidor SMTP (Host)</label>
                  <input id="w-smtp" bind:value={config.SMTPHost} placeholder="smtp.gmail.com" />
              </div>
              <div class="field">
                  <label for="w-port">Puerto SMTP</label>
                  <input id="w-port" type="number" bind:value={config.SMTPPort} placeholder="587" />
              </div>
              <div class="field">
                  <label for="w-user">Usuario / Correo</label>
                  <input id="w-user" bind:value={config.SMTPUser} placeholder="tu@empresa.com" />
              </div>
              <div class="field">
                  <label for="w-smtppass">Contraseña</label>
                  <div class="input-group">
                      {#if showSMTPPassword}
                          <input id="w-smtppass" type="text" bind:value={config.SMTPPassword} />
                      {:else}
                          <input id="w-smtppass" type="password" bind:value={config.SMTPPassword} />
                      {/if}
                      <button class="btn-secondary" on:click={() => showSMTPPassword = !showSMTPPassword}>
                          {showSMTPPassword ? '🙈' : '👁️'}
                      </button>
                  </div>
              </div>
          {/if}
      </div>

      <div class="wizard-footer">
          {#if step > 1}
              <button class="btn-secondary" on:click={() => step--}>Atrás</button>
          {:else}
              <div></div> <!-- Spacer -->
          {/if}

          {#if step < 4}
              <button class="btn-primary" on:click={() => {
                  config.SMTPPort = parseInt(config.SMTPPort);
                  step++;
              }}>Siguiente</button>
          {:else}
              <button class="btn-primary" on:click={finishSetup} disabled={loading}>
                  {loading ? 'Guardando...' : 'Finalizar y Empezar'}
              </button>
          {/if}
      </div>
  </div>
</div>
{/if}

<style>
  .wizard-overlay {
      position: fixed; inset: 0; background: rgba(11, 15, 25, 0.95);
      display: flex; align-items: center; justify-content: center; z-index: 9999;
      backdrop-filter: blur(5px);
  }
  .wizard-card {
      background: #161e31; width: 500px; padding: 2rem; border-radius: 16px;
      border: 1px solid rgba(255,255,255,0.1);
      box-shadow: 0 20px 50px rgba(0,0,0,0.5);
      display: flex; flex-direction: column; gap: 1.5rem;
  }
  .wizard-header h2 { margin: 0; color: white; font-size: 1.5rem; }
  .wizard-header p { margin: 5px 0 0; color: #94a3b8; font-size: 0.9rem; }
  
  .wizard-steps { display: flex; gap: 1rem; padding-bottom: 1rem; border-bottom: 1px solid rgba(255,255,255,0.05); }
  .step { font-size: 0.85rem; color: #64748b; font-weight: 600; }
  .step.active { color: #34d399; }

  .wizard-content { display: flex; flex-direction: column; gap: 1rem; min-height: 200px; }
  
  .wizard-footer { display: flex; justify-content: space-between; margin-top: 1rem; }

  /* Reutilizando estilos de inputs del sistema principal */
  .field { display: flex; flex-direction: column; gap: 6px; }
  label { font-size: 0.85rem; color: #cbd5e1; font-weight: 500; }
  input {
      background: #0B0F19; border: 1px solid #334155; color: white;
      padding: 12px 16px; border-radius: 8px; width: 100%; box-sizing: border-box;
  }
  input:focus { border-color: #34d399; outline: none; }
  .input-group { display: flex; gap: 8px; }
  .hint { font-size: 0.8rem; color: #64748b; margin: -4px 0 8px; }
  
  .checkbox label { display: flex; align-items: center; gap: 8px; cursor: pointer; }
  input[type="checkbox"] { width: 18px; height: 18px; }

  .btn-primary { 
      background: #34d399; color: #064e3b; border: none; padding: 10px 24px; 
      border-radius: 8px; font-weight: 700; cursor: pointer; transition: all 0.2s;
  }
  .btn-secondary {
      background: transparent; color: white; border: 1px solid #334155; 
      padding: 10px 20px; border-radius: 8px; cursor: pointer;
  }
</style>