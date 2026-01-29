# 🚀 MASTER PLAN: Módulo "Punto de Venta (POS) & Satélite Móvil"

**Estado Actual:** Proyecto Kushki V2 en FASE 10 (Feature Complete).
**Stack:** Go (Wails) + Svelte + SQLite + GORM.
**Objetivo:** Extender el sistema de facturación profesional con un modo de venta rápida (Retail POS) y conectividad móvil para gestión de bodega.

---

## FASE 1: Enriquecimiento del Modelo de Datos (Backend Go)

**Contexto:** El modelo de `Product` actual es básico. Necesitamos adaptarlo para retail y mejorar la búsqueda inteligente ya existente.

**Acciones:**

1. **Modificar `internal/db/models.go`:**
   * Actualizar struct `Product`:
     * `Barcode` (string, index único).
     * `AuxiliaryCode` (string).
     * `MinStock` (int).
     * `ExpiryDate` (time.Time).
     * `Location` (string, ej: "Estante A1").
   * Actualizar `ProductDTO` para reflejar estos cambios.

2. **Actualizar `internal/service/search_service.go`:**
   * Modificar `FuzzySearchProducts` para que el `SearchContent` incluya el `Barcode` y `AuxiliaryCode`.

3. **Migración DB:** Wails/GORM manejará la migración automática al añadir los campos al struct, pero se debe asegurar que `SKU` siga siendo la PK o evaluar si `Barcode` es mejor PK para retail. (Se mantendrá `SKU` como PK por compatibilidad).

---

## FASE 2: Interfaz "POS Mode" (Frontend Svelte)

**Contexto:** Crear una experiencia de usuario optimizada para teclado y escáner de barras, integrada en el sidebar actual.

**Acciones:**

1. **Crear `frontend/src/lib/features/pos/PosView.svelte`:**
   * **Layout:**
     * Área de escaneo (Input invisible con auto-focus).
     * Tabla de ítems con fuentes grandes (Mint/Obsidian style).
     * Panel lateral con Total Gigante y selección rápida de cliente.
   * **Atajos de Teclado:**
     * `F12` / `+`: Procesar Venta.
     * `F5`: Búsqueda manual de productos.
     * `ESC`: Limpiar/Cancelar.

2. **Integración con `App.CreateInvoice`:**
   * El POS debe construir el objeto `FacturaDTO` y llamar al método existente en `app.go`.

---

## FASE 3: Servidor Local para Satélite (Backend Go)

**Contexto:** Convertir la instancia de escritorio en un servidor para que los dispositivos móviles de la bodega se conecten vía Wi-Fi local.

**Acciones:**

1. **Implementar Servidor HTTP en `app.go`:**
   * Iniciar un servidor (ej: `chi` o `net/http`) en una goroutine durante el `startup`.
   * Endpoint `GET /api/inventory`: Lista de productos para el móvil.
   * Endpoint `POST /api/inventory/update`: Actualizar stock desde el móvil.

2. **Seguridad y Conectividad:**
   * Generar un PIN/Token de acceso temporal.
   * Función `GetLocalIP()` para mostrar el QR de conexión en la configuración.

---

## FASE 4: App Satélite Móvil (Web App Ligera)

**Contexto:** Una interfaz web optimizada para móviles que se sirve desde el servidor local de Go.

**Acciones:**

1. **Desarrollar `frontend/mobile/`:**
   * SPA ultra-ligera (Svelte o Vanilla JS).
   * Funcionalidades:
     * Escaneo de códigos (Cámara).
     * Ajuste de stock rápido (Ingreso/Egreso).
     * Consulta de precios en percha.

---

## FASE 5: Sincronización en Tiempo Real

**Contexto:** Mantener la UI de escritorio actualizada cuando se realicen cambios desde el móvil.

**Acciones:**

1. **Eventos de Wails:**
   * Al recibir un `POST` en el servidor local de Go, emitir un evento `runtime.EventsEmit(a.ctx, "inventory-updated", product)`.
2. **Listeners:**
   * `ProductList.svelte` y `PosView.svelte` deben escuchar este evento para actualizar stocks sin recargar.

---
*Plan actualizado el 28 de enero de 2026 para alinearse con la arquitectura Kushki V2.*