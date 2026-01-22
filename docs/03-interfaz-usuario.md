---
id: interfaz
title: Manual de Interfaz
---

# Guía de Uso de la Interfaz

Kushki Facturador utiliza una interfaz intuitiva dividida en paneles lógicos.

## 1. Dashboard (Resumen)
Tu centro de mando.
*   **Filtros de Fecha:** En la parte superior derecha, selecciona "Desde" y "Hasta" para analizar un periodo específico.
*   **Tarjetas KPI:**
    *   *Ventas Totales:* Suma de facturas **Autorizadas** en el periodo.
    *   *Facturas:* Cantidad de documentos emitidos.
    *   *Pendientes:* Facturas que requieren tu atención (rechazadas o error de red).
    *   *Estado SRI:* Semáforo de conexión con el servicio de rentas.
*   **Gráfico:** Evolución de ventas día a día.

## 2. Emitir Factura
El proceso de venta simplificado.
1.  **Ambiente:** Un switch en la cabecera te permite cambiar entre `PRUEBAS` (Ámbar) y `PRODUCCIÓN` (Verde) al instante.
2.  **Cliente:** Escribe el nombre o RUC en el buscador. El sistema autocompletará los datos. Si es nuevo, llena los campos y se guardará automáticamente al emitir.
    *   *Validación:* Si olvidas un dato clave (ej. email), el campo se pondrá rojo.
3.  **Productos:** Busca por nombre/código o ingrésalo manualmente.
    *   El cálculo de impuestos (15%, 5%, 0%) es automático.
4.  **Emitir:** Al hacer clic en "Firmar y Emitir", el sistema:
    *   Genera el XML.
    *   Lo firma digitalmente.
    *   Lo envía al SRI.
    *   Genera el PDF.
    *   Envía el correo al cliente.

## 3. Historial
Tu archivo digital.
*   Visualiza todas las facturas emitidas.
*   **Acciones Rápidas:**
    *   📄 **PDF:** Abre el RIDE visualmente.
    *   ✉️ **Email:** Reenvía la factura al cliente con un clic.
    *   🌐 **XML:** Abre el archivo fuente.
    *   📂 **Carpeta:** Te lleva a la ubicación física del archivo en tu disco.

## 4. Inventario y Clientes
Gestión básica de tus bases de datos (CRUD).
*   Puedes Crear, Editar y Eliminar productos o clientes.
*   Los cambios se reflejan inmediatamente en el buscador de la pantalla "Emitir".
