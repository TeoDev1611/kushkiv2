---
id: arquitectura
title: Arquitectura Técnica
---

# Arquitectura del Sistema

Documentación técnica para desarrolladores y mantenimiento. Versión 2.6.0.

## Visión General
Kushki Facturador utiliza una arquitectura monolítica híbrida (Desktop Application).

*   **Frontend:** Single Page Application (SPA) construida con **Svelte**.
*   **Backend:** Binario compilado en **Go**, que actúa como servidor local y controlador.
*   **Comunicación:** Utiliza el puente de **Wails** (basado en WebKit/WebView2) para comunicar Frontend y Backend mediante bindings de JS.

## Tecnologías Clave

### Backend (Go)
*   **ORM:** `GORM` con driver `go-sqlite3`.
*   **XML/Firma:** Implementación nativa de **XAdES-BES** y estructuras estrictas SRI.
*   **Reportes:** `Maroto` para generación de PDFs (Facturas y Cotizaciones).
*   **Gráficos:** `Go-Echarts` para generación de snippets HTML/JS de visualización de datos.
*   **Búsqueda:** `Sahilm/Fuzzy` para indexación en memoria y búsqueda aproximada.
*   **SMTP:** Paquete `net/smtp` estándar para envío de correos.

### Frontend (Svelte)
*   **Estado:** Reactividad nativa de Svelte.
*   **Componentes:** Nuevos módulos `QuotationPanel` y `ChartFrame`.
*   **Estilos:** CSS puro con variables globales (Theme Obsidian).
*   **Build:** Vite.

## Estructura de Base de Datos

El archivo `kushki.db` contiene las siguientes tablas principales:

1.  **EmisorConfig:** Configuración singleton (RUC, Firma, SMTP, Logo).
2.  **Facturas:** Cabeceras de documentos autorizados.
3.  **FacturaItems:** Detalle de productos por factura.
4.  **Quotations:** (NUEVA) Cabeceras de cotizaciones.
5.  **QuotationItems:** (NUEVA) Detalle de productos en cotizaciones.
6.  **Clients / Products:** Catálogos maestros.
7.  **MailLogs:** Auditoría de envíos de correo.

## Servicios Internos

El backend se organiza en servicios inyectados en la estructura principal `App`:

*   `InvoiceService`: Lógica de emisión SRI.
*   `QuotationService`: Gestión de ciclo de vida de cotizaciones.
*   `SearchService`: Motor de búsqueda fuzzy. Carga proyecciones ligeras de datos en memoria para búsquedas instantáneas.
*   `ChartService`: Ejecuta consultas SQL de agregación y renderiza gráficos.
*   `MailService`: Cliente SMTP.

## Flujo de Emisión

1.  **UI:** Usuario llena formulario -> `CreateInvoice(DTO)`.
2.  **Backend:**
    *   Valida datos y reglas de negocio.
    *   Genera XML (`pkg/xml`).
    *   Firma XML (`pkg/crypto`).
    *   Envía a WS Recepción SRI (`pkg/sri`).
    *   Si es RECIBIDA -> Consulta WS Autorización SRI.
    *   Genera RIDE PDF (`pkg/pdf`).
    *   Envía Correo (`service/mail_service`).
    *   Guarda en DB y Logs.
3.  **UI:** Recibe notificación de éxito/error.