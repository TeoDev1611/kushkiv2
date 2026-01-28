# GEMINI.md - Contexto del Proyecto "Kushki Facturador"

## Estado del Proyecto: FASE 10 (Pulido Final y Documentación)

Este proyecto ha evolucionado a un sistema de facturación electrónica profesional de alto nivel. Se han completado los módulos de **Cotizaciones**, **Analítica Avanzada** y **Búsqueda Inteligente**, consolidando la oferta de valor.

## 1. Arquitectura Implementada

*   **Frontend (`frontend/`):** Svelte. Diseño "Obsidian & Mint". Componentes reactivos (`ChartFrame`, `QuotationPanel`).
*   **Backend (`app.go`):** Controlador principal optimizado. Gestión de Licencia, Backups, SMTP Local y Bridge Wails.
*   **Servicios (`internal/service/`):**
    *   `InvoiceService`: Núcleo de facturación SRI con validación estricta.
    *   `QuotationService`: (NUEVO) Gestión de cotizaciones y conversión a facturas.
    *   `SearchService`: (NUEVO) Motor de búsqueda Fuzzy en memoria para Clientes, Productos y Facturas.
    *   `ChartService`: (NUEVO) Generador de gráficos interactivos con `go-echarts`.
    *   `MailService`: Motor SMTP Local.
    *   `SyncService` & `ReportService`: Utilidades de soporte.
*   **Base de Datos:** SQLite + GORM. Nuevas tablas: `Quotation`, `QuotationItem`.

## 2. Hitos Recientes (Q1 2026 - Feature Complete)

*   **Módulo de Cotizaciones:**
    *   Interfaz dedicada para crear y gestionar proformas.
    *   Generación automática de PDF.
    *   Conversión de "Cotización a Factura" con un clic.
*   **Inteligencia de Negocios:**
    *   **Gráficos:** Implementación de `go-echarts` para visualizar "Ingresos Mensuales" y "Top Clientes" directamente en el Dashboard.
    *   **Búsqueda Fuzzy:** Implementación de `sahilm/fuzzy` para búsquedas tolerantes a errores en todo el sistema.
*   **Autonomía de Correo (SMTP Local):**
    *   Soporte nativo para Gmail/Outlook.
    *   Validación de conexión SMTP.
*   **UI/UX Refinada:**
    *   Corrección masiva de accesibilidad (A11y).
    *   Dashboard con grid de gráficos interactivos.

## 3. Comandos Útiles

*   **Desarrollo:** `wails dev`
*   **Modo Debug:** `./build/bin/kushkiv2 --kushki-debug`
*   **Compilación:** `wails build`

## 4. Estado Final

El sistema es una solución "todo en uno" (All-in-One) para la facturación electrónica en Ecuador. Cumple con la normativa SRI, ofrece herramientas de venta (Cotizaciones) y análisis (Gráficos), todo bajo una arquitectura segura y local.
