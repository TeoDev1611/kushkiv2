# Documentación Técnica (TECH.md)

## Arquitectura del Sistema

### 1. Diagrama de Componentes

```mermaid
graph TD
    UI[Frontend Svelte] <-->|Wails Bridge| App[Go Backend Desktop]
    App <-->|GORM| DB[(SQLite Local)]
    App -->|HTTPS/JSON| API[Cloud API (Deno)]
    App -->|SOAP| SRI[SRI Web Services]
    
    subgraph "Go Internal Services"
        InvoiceSvc[Invoice Service]
        CloudSvc[Cloud Service]
        SyncSvc[Sync Service]
    end
    
    App --> InvoiceSvc
    App --> CloudSvc
    App --> SyncSvc
```

### 2. Módulo de Seguridad (Licenciamiento)

El sistema implementa un modelo de seguridad "Node-Locked" para prevenir la piratería.

*   **Machine ID:** Se genera un hash SHA-256 único basado en: `Hostname + OS + Arch`.
*   **Activación:**
    *   Endpoint: `POST /api/v1/license/activate`
    *   Payload: `{ license_key: "...", machine_id: "..." }`
    *   Respuesta: Token JWT que se almacena localmente (`EmisorConfig.LicenseToken`).
*   **Persistencia:** La clave y el token se guardan en la tabla `emisor_configs`.
*   **Bloqueo UI:** El Frontend verifica el estado de la licencia antes de montar el Dashboard.

### 3. Servicio de Nube (`CloudService`)

Reemplaza la antigua infraestructura SMTP local. Centraliza la comunicación externa.

*   **Responsabilidades:**
    1.  Validación de Licencias.
    2.  Envío de Correos Transaccionales (Facturas PDF).
*   **Envío de PDF:**
    *   Usa `multipart/form-data`.
    *   Envía el PDF generado en memoria (sin guardar en disco temporalmente para el envío) directamente a la API.
    *   Endpoint: `POST /enviar-pdf`.

### 4. Flujo de Emisión de Factura

1.  **UI:** Usuario llena formulario y da clic en "Emitir".
2.  **App:** `InvoiceService` construye XML, firma (XAdES-BES) y envía al SRI.
3.  **App:** Si SRI autoriza, `ReportService` genera el RIDE (PDF).
4.  **App:** Guarda XML y PDF en disco local (organizado por Año/Mes).
5.  **App (Async):** Invoca `CloudService.SendPDFReport()` en una goroutine para enviar el correo al cliente final vía API Cloud.

### 5. Base de Datos

*   **Motor:** SQLite 3.
*   **ORM:** GORM.
*   **Tablas Críticas:**
    *   `emisor_configs`: Configuración y Credenciales de Licencia.
    *   `facturas`: Historial transaccional completo (incluye blobs de XML/PDF).
    *   `products` / `clients`: Catálogos maestros.

### 6. Configuración de Usuario

Se ha eliminado la configuración técnica compleja (SMTP). El usuario solo configura:
1.  **Empresa:** RUC, Razón Social, Dirección.
2.  **Firma:** Archivo `.p12` y Contraseña.
3.  **Rutas:** Carpeta de almacenamiento (opcional).

Todo lo demás (servidor de correos, validación) es gestionado por la plataforma Cloud.