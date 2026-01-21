# GEMINI.md - Contexto del Proyecto "Kushki Facturador"

## Estado del Proyecto: FASE 7 (Seguridad y Nube)

Este proyecto es un sistema de facturación electrónica híbrido (Desktop/Cloud) para Ecuador.

## 1. Arquitectura Implementada

*   **Frontend (`frontend/`):** Svelte. Incluye Panel de Licencia (Bloqueante) y Wizard de Configuración.
*   **Backend (`app.go`):** Controlador principal. Gestiona el ciclo de vida, seguridad y puente Wails.
*   **Servicios (`internal/service/`):**
    *   `CloudService` (NUEVO): Cliente HTTP para API Deno (Licencias/Email).
    *   `InvoiceService`: Núcleo de facturación SRI.
    *   `SyncService`: Sincronización de fondo.
    *   `ReportService`: Generador PDF/Excel.
*   **Persistencia:** SQLite/GORM. Tabla `EmisorConfig` actualizada con campos de licencia.

## 2. Cambios Recientes (Q1 2026)

*   **Seguridad:** Implementación de Licenciamiento Node-Locked. UI de bloqueo si no hay licencia activa.
*   **Cloud:** Integración con API Deno para delegar el envío de correos.
*   **Limpieza:** Eliminación del módulo `MailService` (SMTP Legacy) y su configuración en UI.
*   **UX:** Flujo forzado: Licencia -> Wizard -> Dashboard.

## 3. Comandos Útiles

*   **Desarrollo:** `wails dev`
*   **Compilación:** `wails build`

## 4. Próximos Pasos

1.  Refinar la validación del Token JWT de licencia en el lado del cliente (verificar firma).
2.  Implementar webhooks en el Backend Cloud para notificar estado de correos.
3.  Optimizar la carga de reportes históricos masivos.
