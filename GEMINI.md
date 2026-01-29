# GEMINI.md - Contexto del Proyecto "Kushki Facturador"

## Estado del Proyecto: FASE 12 (Mobile v2 & Optimización) - EN CURSO

Este proyecto ha evolucionado a un ecosistema de facturación omnicanal. Se han optimizado los módulos de **Punto de Venta (Retail)** y **Conectividad Móvil (Echo Server)**, permitiendo una experiencia de usuario fluida y profesional.

## 1. Arquitectura Implementada

*   **Frontend (`frontend/`):** Svelte. Diseño "Obsidian & Mint". Componentes reactivos (`PosView`, `ProductList`).
*   **Backend (`app.go`):** Controlador principal optimizado con **Echo Framework**.
*   **Móvil (`internal/mobile/`):** App Web modernizada con soporte para escaneo de cámara (HTML5-QRCode) y Modo POS remoto.
*   **Servicios (`internal/service/`):**
    *   `InvoiceService`: Facturación con re-generación de secuencial atómica.
    *   `SearchService`: Motor de búsqueda Fuzzy optimizado.
*   **Base de Datos:** SQLite + GORM. Optimizada para actualizaciones de stock concurrentes.

## 2. Hitos Recientes (Enero 2026)

*   **Módulo POS (Retail):**
    *   Búsqueda y creación rápida de clientes (+Nuevo).
    *   Sincronización instantánea con Escáner Móvil.
    *   Acceso QR directo desde la cabecera.
*   **Satélite Móvil v2:**
    *   Migración a **Echo Framework** con Gzip y CORS.
    *   Interfaz rediseñada (Obsidian & Mint).
    *   Lector de cámara integrado para inventario y ventas.
*   **Debug & Logs:**
    *   Logs condicionales (`--kushki-debug`).
    *   Formateo de respuestas SRI XML a JSON legible en consola.
*   **Robustez:**
    *   Corrección de conflictos de concurrencia en `ClaveAcceso` y `Unique Constraints`.
    *   Detección inteligente de IP LAN y override manual.

## 3. Comandos Útiles

*   **Desarrollo:** `wails dev`
*   **Modo Debug:** `./build/bin/kushkiv2 --kushki-debug`
*   **Compilación:** `wails build`

## 4. Estado Final

El sistema es una solución "todo en uno" para Ecuador. Integra Facturación Electrónica, Punto de Venta (POS) y Gestión de Bodega Móvil en un solo ejecutable, sin dependencias externas complejas.
