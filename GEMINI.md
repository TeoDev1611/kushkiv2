# GEMINI.md - Contexto del Proyecto "Kushki Facturador"

## Estado del Proyecto: FASE 9 (Consolidación UI, SMTP Local y Auditoría)

Este proyecto ha evolucionado a un sistema de facturación electrónica profesional, priorizando la autonomía del usuario (SMTP Local), la seguridad de datos (Backups) y una experiencia de usuario pulida (Dashboard Dinámico, UI Moderna).

## 1. Arquitectura Implementada

*   **Frontend (`frontend/`):** Svelte. Diseño "Obsidian & Mint". Panel de Actividad con Auditoría dual (Correos/Sistema). Wizard de Configuración con selector de archivos.
*   **Backend (`app.go`):** Controlador principal optimizado. Gestión de Licencia, Backups, SMTP Local y Bridge Wails.
*   **Servicios (`internal/service/`):**
    *   `InvoiceService`: Núcleo de facturación SRI con validación estricta (Decimales, Plazos, XML 1.1.0).
    *   `MailService`: Motor SMTP Local con plantillas HTML y adjuntos PDF.
    *   `SyncService`: Worker de sincronización y registro de logs técnicos.
    *   `ReportService`: Generador PDF/Excel.
    *   `CloudService`: Cliente HTTP Seguro (Limitado solo a Licenciamiento).
*   **Base de Datos:** SQLite + GORM. Nuevas tablas: `MailLog` (Historial de correos). Nuevos campos en `EmisorConfig` (Logo, RIMPE, Retención).

## 2. Hitos Recientes (Q1 2026 - Sprint Final)

*   **Autonomía de Correo (SMTP Local):**
    *   Eliminación de dependencia de API Cloud para envíos.
    *   Soporte nativo para Gmail/Outlook con configuración guiada.
    *   Validación de conexión SMTP ("Test Button").
    *   Plantillas de correo HTML profesionales.
*   **UI/UX de Primer Nivel:**
    *   **Dashboard:** Estadísticas dinámicas con filtrado por rango de fechas real.
    *   **Notificaciones:** Sistema de Toasts apilables y Centro de Notificaciones (Historial de sesión).
    *   **Personalización:** Soporte para Logo de empresa (Redimensionado automático).
*   **Cumplimiento SRI:**
    *   Corrección de estructura XML (Detalles Adicionales, Tarifas).
    *   Manejo de RIMPE y Agente de Retención.
    *   Validación de montos >$1000 (Solo producción).
*   **Seguridad y Auditoría:**
    *   Módulo de **Backups**: Generación y gestión de copias de seguridad (.zip).
    *   Pestaña **Actividad**: Registro detallado de cada correo enviado (Éxito/Fallo).

## 3. Comandos Útiles

*   **Desarrollo:** `wails dev`
*   **Modo Debug:** `./build/bin/kushkiv2 --kushki-debug`
*   **Compilación:** `wails build`

## 4. Estado Final

El sistema es estable y funcional para entornos de Producción y Pruebas. Cumple con la ficha técnica del SRI y ofrece herramientas de gestión completas para el usuario final.