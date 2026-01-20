# 🛠️ TECH.md - Documentación Técnica y Arquitectura

> **Solo para Desarrolladores.** Este documento detalla la estructura interna, flujos de datos y decisiones de diseño de Kushki Facturador v2.

---

## 1. Arquitectura de Alto Nivel

El sistema sigue una arquitectura **Híbrida Nativa** utilizando el patrón **Backend-as-a-Service (BaaS)** local.

*   **Frontend (UI):** Svelte (JavaScript/HTML/CSS). Se ejecuta en un WebView del sistema (WebKit/WebView2).
*   **Backend (Core):** Go (Golang). Maneja toda la lógica de negocio, criptografía, base de datos y comunicación con el SRI.
*   **Bridge:** [Wails v2](https://wails.io). Interconecta JS y Go. Las funciones de Go en `app.go` se exponen como Promesas en JS.

---

## 2. Mapa del Código (Directory Structure)

### 🟢 Backend (Go)
*   **`main.go`**: Entrypoint. Configura Wails, ciclo de vida (`OnStartup`, `OnShutdown`) y gestiona el cierre seguro de la DB.
*   **`app.go`**: **Controlador Principal**. Aquí están los métodos públicos expuestos al Frontend (e.g., `CreateInvoice`, `GetDashboardStats`). Actúa como orquestador.
*   **`internal/`**: Código privado de la aplicación.
    *   **`db/`**: Modelos GORM (`models.go`), conexión SQLite (`connection.go`) y migraciones (`migrations.go`).
    *   **`service/`**: Lógica de negocio pura.
        *   `invoice_service.go`: Orquesta XML, Firma, SRI y PDF.
        *   `sync_service.go`: Maneja la cola de reintentos y workers concurrentes.
        *   `mail_service.go`: Envío SMTP.
*   **`pkg/`**: Librerías reutilizables/bajas.
    *   **`crypto/`**: Implementación manual de **XAdES-BES** y manejo de certificados P12.
    *   **`sri/`**: Cliente SOAP para Recepción y Autorización.
    *   **`xml/`**: Estructuras UBL 2.1 (Factura Electrónica).
    *   **`pdf/`**: Generador de RIDE usando `maroto`.

### 🟠 Frontend (Svelte)
*   **`src/App.svelte`**: Single Page Application. Contiene el Router (Tabs), Estado Global y lógica de UI.
*   **`src/components/`**:
    *   `Sidebar.svelte`: Navegación lateral colapsable.
    *   `Wizard.svelte`: Asistente de configuración inicial.
*   **`wailsjs/`**: **NO TOCAR**. Código autogenerado por Wails que conecta JS con Go.

---

## 3. Flujos Críticos

### A. Emisión de Factura (El "Hot Path")
1.  **Frontend:** Recoge datos → Llama a `CreateInvoice(dto)`.
2.  **Go (App):** Recibe DTO → Pasa a `InvoiceService`.
3.  **InvoiceService:**
    *   Genera **Clave de Acceso** (Algoritmo Modulo 11).
    *   Construye XML (UBL 2.1) en memoria.
    *   **Firma:** Usa `pkg/crypto` para inyectar la firma XAdES-BES en el XML.
    *   **SRI Recepción:** Envía XML firmado al WebService del SRI.
    *   **SRI Autorización:** Consulta estado.
    *   **PDF:** Genera el RIDE con código QR.
    *   **DB:** Guarda la transacción en SQLite.
4.  **Go (App):** Guarda archivos físicos (`/Año/Mes/FACTURA-001...`) y encola email.
5.  **Frontend:** Recibe "Éxito" y actualiza Dashboard.

### B. Sincronización (Worker Pool)
Para no congelar la UI al procesar facturas pendientes:
1.  `SyncService` inicia un **Worker** en segundo plano (Goroutine).
2.  Usa un canal semáforo (`make(chan struct{}, 3)`) para limitar a **3 envíos simultáneos** al SRI.
3.  Si el SRI responde, actualiza el estado en DB y genera logs detallados en memoria para el panel "Sincronización".

### C. Dashboard (Concurrencia)
Al cargar `GetDashboardStats`, Go lanza **4 Goroutines** en paralelo usando `sync.WaitGroup`:
1.  Suma de ventas del mes.
2.  Conteo de facturas.
3.  Conteo de pendientes.
4.  Cálculo de tendencia (Gráfico) mediante SQL optimizado.
Esto reduce el tiempo de carga de ~500ms a ~50ms.

---

## 4. Base de Datos (SQLite)

### Configuración
*   **Archivo:** `kushki.db` en la raíz.
*   **Modo:** `WAL` (Write-Ahead Logging) habilitado en `internal/db/connection.go`. Permite lecturas mientras se escribe.
*   **Índices:** Se añaden índices manuales en `migrations.go` para:
    *   `fecha_emision` + `estado_sri` (Dashboard).
    *   `created_at` (Historial).

### Tablas Clave
*   `emisor_configs`: Configuración única (RUC, Firma, SMTP).
*   `facturas`: Cabecera de documentos. Contiene BLOBs para `xml_firmado` y `pdf_ride`.
*   `factura_items`: Detalle de productos por factura.
*   `products`: Inventario.
*   `clients`: Directorio.

---

## 5. Criptografía (XAdES-BES)

La firma **NO** usa librerías externas de Java ni OpenSSL. Es una implementación nativa en Go (`pkg/crypto/signer.go`).

*   **Proceso:**
    1.  Calcula Hash SHA1 del XML ("Canonicalizado").
    2.  Firma el Hash con la llave privada del `.p12`.
    3.  Construye la estructura `KeyInfo`, `SignedProperties` y `SignatureValue` según estándar SRI.
    4.  Inyecta el nodo `<ds:Signature>` en el XML original.

> **Nota:** Si cambias algo en la estructura del XML antes de firmar, la firma se romperá (Error "Firma Inválida"). La canonicalización es estricta.

---

## 6. Guía para Extender

### ¿Cómo añadir un nuevo reporte?
1.  Crea la función en `internal/service/report_service.go`.
2.  Exponla en `app.go` (struct `App`).
3.  Ejecuta `wails dev` para regenerar los bindings en `frontend/wailsjs/`.
4.  Consúmela en Svelte importándola desde `../wailsjs/go/main/App.js`.

### ¿Cómo añadir un nuevo tipo de documento (e.g., Retenciones)?
1.  Define la estructura XML en `pkg/xml/structures.go`.
2.  Crea un nuevo servicio o extiende `invoice_service.go`.
3.  Asegúrate de cambiar el "Tipo de Comprobante" en la generación de la Clave de Acceso (`pkg/util/mod11.go`).

---

## 7. Solución de Problemas Comunes

*   **Error "Database Locked":** Ocurre si dos goroutines intentan escribir sin usar el pool de conexiones correcto. GORM + WAL mode ya lo mitiga, pero asegura siempre usar `db.GetDB()`.
*   **Interfaz Lenta en Linux:** Verifica que no hayas reintroducido `backdrop-filter: blur` en CSS. Wails usa WebKitGTK que no optimiza bien ese filtro.
*   **Firma Inválida en SRI:** Revisa `signer.go`. El SRI exige que los namespaces XML (`xmlns`) estén exactamente en el orden correcto y sin espacios extra.

---

**Desarrollado con ❤️ y Go.**
