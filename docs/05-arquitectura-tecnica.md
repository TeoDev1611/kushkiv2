---
id: arquitectura
title: Arquitectura Técnica
---

# Arquitectura del Sistema

Documentación técnica para desarrolladores y mantenimiento.

## Visión General
Kushki Facturador utiliza una arquitectura monolítica híbrida (Desktop Application).

*   **Frontend:** Single Page Application (SPA) construida con **Svelte**.
*   **Backend:** Binario compilado en **Go**, que actúa como servidor local y controlador.
*   **Comunicación:** Utiliza el puente de **Wails** (basado en WebKit/WebView2) para comunicar Frontend y Backend mediante bindings de JS.

## Tecnologías Clave

### Backend (Go)
*   **ORM:** `GORM` con driver `go-sqlite3`. Manejo de concurrencia activado (WAL mode).
*   **XML:** Paquete `encoding/xml` con estructuras estrictas (`pkg/xml/structures.go`) alineadas al XSD 1.1.0 del SRI.
*   **Firma:** Implementación nativa de **XAdES-BES** (`pkg/crypto/signer.go`). Realiza canonicalización, hashing SHA1 y firma RSA.
*   **PDF:** Librería `Maroto` para generación declarativa de PDFs.
*   **SMTP:** Paquete `net/smtp` estándar para envío de correos con soporte AUTH PLAIN/LOGIN.

### Frontend (Svelte)
*   **Estado:** Reactividad nativa de Svelte.
*   **Estilos:** CSS puro con variables globales (Theme Obsidian).
*   **Build:** Vite.

## Estructura de Base de Datos

El archivo `kushki.db` contiene las siguientes tablas principales:

1.  **EmisorConfig:** Configuración singleton (RUC, Firma, SMTP, Logo).
2.  **Facturas:** Cabeceras de documentos (Clave de acceso, Totales, Estado SRI).
3.  **FacturaItems:** Detalle de productos por factura (para reportes).
4.  **Clients / Products:** Catálogos maestros.
5.  **MailLogs:** Auditoría de envíos de correo.

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
