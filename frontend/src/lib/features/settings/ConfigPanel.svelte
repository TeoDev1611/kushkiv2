    <script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { fade } from 'svelte/transition';
    import { Backend } from '$lib/services/api';
    import { notifications } from '$lib/stores/notifications';
    import { withLoading } from '$lib/stores/app';
    import * as WailsApp from 'wailsjs/go/main/App'; 
    import QRCode from 'qrcode';

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
        LogoPath: "",
        PDFTheme: "modern",
        SMTPHost: "",
        SMTPPort: 587,
        SMTPUser: "",
        SMTPPassword: ""
    };

    let showPassword = false;
    let satelliteInfo: any = null;

    // Handler para evento global de guardado (Ctrl+S)
    const handleGlobalSave = () => handleSaveConfig();

    onMount(async () => {
        window.addEventListener('app-save', handleGlobalSave);
        try {
            const cfg = await withLoading(Backend.getConfig());
            if (cfg) {
                config = { ...config, ...cfg };
            }
            loadSatelliteInfo();
        } catch (e) {
            notifications.show("Error cargando configuración: " + e, "error");
        }
    });

    onDestroy(() => {
        window.removeEventListener('app-save', handleGlobalSave);
    });

    async function loadSatelliteInfo() {
        try {
            const info = await WailsApp.GetSatelliteConnectionInfo();
            if (info) {
                // Si ya teníamos una IP manual editada, la conservamos. Si no, usamos la detectada.
                if (!satelliteInfo) {
                    satelliteInfo = info;
                }
            }
            updateQR();
        } catch (e) {
            console.error("Error loading satellite info:", e);
        }
    }

    function updateQR() {
        const canvas = document.getElementById('qr-canvas');
        if (canvas && satelliteInfo) {
            // Reconstruir URL con la IP (posiblemente editada)
            const url = `http://${satelliteInfo.ip}:${satelliteInfo.port}/?token=${satelliteInfo.token}`;
            
            QRCode.toCanvas(canvas, url, { 
                width: 180,
                margin: 1,
                color: {
                    dark: '#000000',
                    light: '#ffffff'
                }
            }, (error) => {
                if (error) console.error("QR Error:", error);
            });
        }
    }

    async function handleSaveConfig() {        try {
            // Asegurar tipos numéricos
            config.Ambiente = parseInt(String(config.Ambiente));
            config.SMTPPort = parseInt(String(config.SMTPPort));
            
            const res = await withLoading(Backend.saveConfig(config));
            
            if (res.startsWith("Error")) {
                notifications.show(res, "error");
            } else {
                notifications.show("Configuración guardada exitosamente", "success");
            }
        } catch (err) {
            notifications.show("Error guardando: " + err, "error");
        }
    }

    async function handleSelectCert() {
        try {
            const path = await WailsApp.SelectCertificate();
            if (path) config.P12Path = path;
        } catch (e) {
            notifications.show("Error seleccionando certificado: " + e, "error");
        }
    }

    async function handleSelectLogo() {
        try {
            const path = await WailsApp.SelectAndSaveLogo();
            if (path && !path.startsWith("Error")) {
                config.LogoPath = path;
            } else if (path) {
                notifications.show(path, "error");
            }
        } catch (e) {
            notifications.show("Error subiendo logo: " + e, "error");
        }
    }

    async function handleTestSMTP() {
        if (!config.SMTPHost || !config.SMTPUser || !config.SMTPPassword) {
            return notifications.show("Complete los datos de correo primero", "warning");
        }
        try {
            const res = await withLoading(WailsApp.TestSMTPConnection(config));
            notifications.show(res, res.includes("Error") ? "error" : "success");
        } catch (e) {
            notifications.show("Fallo crítico: " + e, "error");
        }
    }

    function setSMTPProvider(provider: 'gmail' | 'outlook') {
        if (provider === 'gmail') {
            config.SMTPHost = "smtp.gmail.com";
            config.SMTPPort = 587;
        } else {
            config.SMTPHost = "smtp.office365.com";
            config.SMTPPort = 587;
        }
    }
</script>

<div class="panel" in:fade={{ duration: 200 }}>
    <div class="header-row">
        <h1>Configuración del Sistema</h1>
        <button class="btn-primary" on:click={handleSaveConfig}>💾 Guardar Cambios</button>
    </div>

    <div class="config-grid">
        <!-- Datos de Empresa -->
        <div class="card">
            <h3>🏢 Datos de Empresa</h3>
            <div class="field">
                <label for="cfg-ruc">RUC</label>
                <input id="cfg-ruc" bind:value={config.RUC} placeholder="179..." />
            </div>
            <div class="field">
                <label for="cfg-razon">Razón Social</label>
                <input id="cfg-razon" bind:value={config.RazonSocial} />
            </div>
            <div class="field">
                <label for="cfg-nombre">Nombre Comercial</label>
                <input id="cfg-nombre" bind:value={config.NombreComercial} />
            </div>
            <div class="field">
                <label for="cfg-dir">Dirección Matriz</label>
                <input id="cfg-dir" bind:value={config.Direccion} />
            </div>
            <div class="grid col-2-tight">
                <div class="field">
                    <label for="cfg-estab">Estab (001)</label>
                    <input id="cfg-estab" bind:value={config.Estab} maxlength="3" class="text-center" />
                </div>
                <div class="field">
                    <label for="cfg-ptoemi">PtoEmi (001)</label>
                    <input id="cfg-ptoemi" bind:value={config.PtoEmi} maxlength="3" class="text-center" />
                </div>
            </div>
            <div class="field checkbox-field mt-1">
                <label for="cfg-obligado" style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
                    <input id="cfg-obligado" type="checkbox" bind:checked={config.Obligado} style="width: auto;" />
                    Obligado a llevar Contabilidad
                </label>
            </div>
        </div>

        <!-- Firma y Marca -->
        <div class="card">
            <h3>🔐 Firma Electrónica y Marca</h3>
            
            <!-- Logo -->
            <div class="field">
                <label for="cfg-logo">Logo Institucional</label>
                <div id="cfg-logo" class="flex-row" style="background: var(--bg-surface-hover); padding: 10px; border-radius: 8px; gap: 12px; align-items: center;">
                    {#if config.LogoPath}
                        <img 
                            src={config.LogoPath.startsWith("/") ? "file://" + config.LogoPath : config.LogoPath} 
                            alt="Logo" 
                            style="width: 50px; height: 50px; object-fit: contain; background: white; border-radius: 4px;"
                        />
                    {:else}
                        <div style="width: 50px; height: 50px; background: var(--bg-surface); border-radius: 4px; display: flex; align-items: center; justify-content: center; color: var(--text-tertiary);">?</div>
                    {/if}
                    <button class="btn-secondary" on:click={handleSelectLogo}>📷 Cambiar Logo</button>
                </div>
            </div>

            <!-- Tema PDF -->
            <div class="field">
                <label for="cfg-theme">Diseño de Factura PDF</label>
                <select id="cfg-theme" bind:value={config.PDFTheme}>
                    <option value="modern">✨ Moderno (Color Esmeralda)</option>
                    <option value="minimal">🖨️ Minimalista (Ahorro Tinta)</option>
                    <option value="corporate">🏢 Corporativo (Formal B/N)</option>
                </select>
                <p class="text-caption text-secondary mt-1">
                    Seleccione el estilo visual para los documentos generados.
                </p>
            </div>

            <!-- Certificado -->
            <div class="field">
                <label for="cfg-p12">Certificado (.p12)</label>
                <div class="input-group">
                    <input id="cfg-p12" bind:value={config.P12Path} readonly placeholder="Seleccione archivo..." />
                    <button class="btn-secondary" on:click={handleSelectCert}>📂</button>
                </div>
            </div>
            
            <div class="field">
                <label for="cfg-pass">Contraseña del Certificado</label>
                <div class="input-group">
                    {#if showPassword}
                        <input id="cfg-pass" type="text" bind:value={config.P12Password} />
                    {:else}
                        <input id="cfg-pass" type="password" bind:value={config.P12Password} />
                    {/if}
                    <button class="btn-secondary" on:click={() => showPassword = !showPassword}>👁</button>
                </div>
            </div>

            <!-- Ambiente -->
            <div class="field mt-2">
                <label for="cfg-env">Ambiente SRI</label>
                <select id="cfg-env" bind:value={config.Ambiente}>
                    <option value={1}>Pruebas</option>
                    <option value={2}>Producción</option>
                </select>
                <p class="text-caption text-secondary mt-1">
                    {#if config.Ambiente === 1}
                        ⚠️ En modo PRUEBAS las facturas no tienen validez tributaria.
                    {:else}
                        🚀 En modo PRODUCCIÓN las facturas son oficiales.
                    {/if}
                </p>
            </div>
        </div>

        <!-- Servidor de Correo -->
        <div class="card">
            <h3>📧 Configuración SMTP</h3>
            <p class="text-secondary text-caption mb-2">Para el envío automático de facturas a clientes.</p>
            
            <div class="flex-row mb-2" style="gap: 12px;">
                <button class="btn-secondary small" on:click={() => setSMTPProvider('gmail')}>Usar Gmail</button>
                <button class="btn-secondary small" on:click={() => setSMTPProvider('outlook')}>Usar Outlook</button>
            </div>
            
            <div class="field">
                <label for="smtp-host">Host SMTP</label>
                <input id="smtp-host" bind:value={config.SMTPHost} placeholder="smtp.gmail.com" />
            </div>
            
            <div class="grid col-2-tight">
                <div class="field">
                    <label for="smtp-port">Puerto</label>
                    <input id="smtp-port" type="number" bind:value={config.SMTPPort} />
                </div>
                <div class="field">
                    <label for="smtp-user">Usuario / Correo</label>
                    <input id="smtp-user" bind:value={config.SMTPUser} placeholder="tu@email.com" />
                </div>
            </div>
            
            <div class="field">
                <label for="smtp-pass">Contraseña (App Password)</label>
                <input id="smtp-pass" type="password" bind:value={config.SMTPPassword} placeholder="••••••••" />
            </div>
            
            <button class="btn-secondary mt-2 full-width" on:click={handleTestSMTP}>📡 Probar Conexión</button>
        </div>

        <!-- Satélite Móvil -->
        <div class="card">
            <h3>📱 Satélite Móvil</h3>
            <p class="text-secondary text-caption mb-2">Conecte su celular para gestión de inventario y escáner.</p>

            <div class="satellite-container" style="display: flex; flex-direction: column; align-items: center; gap: 15px; background: rgba(0,0,0,0.2); padding: 20px; border-radius: 8px;">
                <canvas id="qr-canvas" style="border-radius: 8px;"></canvas>
                
                {#if satelliteInfo}
                    <div class="pin-display" style="text-align: center;">
                        <span style="display: block; font-size: 12px; color: var(--text-tertiary);">PIN DE CONEXIÓN</span>
                        <span style="font-size: 32px; font-weight: 900; letter-spacing: 4px; color: var(--accent-mint);">{satelliteInfo.token}</span>
                    </div>
                    
                    <div class="ip-control" style="width: 100%;">
                        <label for="sat-ip" style="font-size: 12px; color: var(--text-secondary); display: block; margin-bottom: 4px;">Dirección IP del PC:</label>
                        <div style="display: flex; gap: 8px;">
                            <input 
                                id="sat-ip" 
                                type="text" 
                                bind:value={satelliteInfo.ip} 
                                on:input={updateQR}
                                style="text-align: center; font-family: monospace; letter-spacing: 1px;"
                            />
                        </div>
                        <p class="text-caption text-secondary mt-1 text-center">
                            Puerto: {satelliteInfo.port}
                        </p>
                    </div>

                    <div style="background: rgba(255, 165, 0, 0.15); border: 1px solid orange; padding: 10px; border-radius: 6px; font-size: 12px; color: #ffcc80;">
                        <strong>⚠️ ¿No conecta?</strong>
                        <ul style="padding-left: 20px; margin: 5px 0 0 0; display: flex; flex-direction: column; gap: 4px;">
                            <li>
                                <strong>Firewall:</strong> Asegúrese de permitir el tráfico en el puerto <strong>8085</strong> (TCP).
                                <br><span style="opacity: 0.8; font-size: 11px;">(Windows Defender, UFW, Firewalld, etc.)</span>
                            </li>
                            <li>Verifique que ambos dispositivos estén en la <strong>misma red Wi-Fi</strong>.</li>
                            <li><a href="#" style="color: var(--accent-mint); text-decoration: underline;" on:click|preventDefault={() => notifications.show("Revise la documentación en la carpeta 'docs'", "info")}>Ver Documentación de Ayuda</a></li>
                        </ul>
                    </div>
                {:else}
                    <p>Cargando información del servidor...</p>
                {/if}
            </div>
            <p class="text-caption text-secondary mt-2 text-center">
                Asegúrese de que el dispositivo móvil esté conectado a la misma red Wi-Fi.
            </p>
        </div>
    </div>
</div>

<style>
    .config-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
        gap: 24px;
        padding-bottom: 40px;
    }
    
    .text-caption { font-size: 0.85rem; }
    .btn-secondary.small { padding: 4px 12px; font-size: 0.85rem; }
</style>
