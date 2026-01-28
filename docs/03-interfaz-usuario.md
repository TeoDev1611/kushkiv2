---
id: interfaz
title: Manual de Interfaz
---

# Guía de Uso de la Interfaz

Kushki Facturador utiliza una interfaz intuitiva dividida en paneles lógicos.

## 1. Dashboard (Resumen)
Tu centro de mando con analítica avanzada.
*   **Filtros de Fecha:** En la parte superior derecha, selecciona "Desde" y "Hasta" para analizar un periodo específico.
*   **Tarjetas KPI:**
    *   *Ventas Totales:* Suma de facturas **Autorizadas** en el periodo.
    *   *Facturas:* Cantidad de documentos emitidos.
    *   *Pendientes:* Facturas que requieren tu atención.
    *   *Estado SRI:* Semáforo de conexión con el servicio de rentas.
*   **Gráficos Interactivos:**
    *   **Evolución de Ingresos:** Gráfico de línea suavizada que muestra tus ventas mes a mes.
    *   **Top 5 Clientes:** Gráfico de pastel para identificar tus cuentas clave.

## 2. Cotizaciones (Nuevo Panel)
Gestiona tus propuestas comerciales.
*   **Crear:** Selecciona cliente y productos igual que en una factura. Se genera un número secuencial único.
*   **PDF:** El botón de documento (📄) abre una proforma en PDF lista para enviar o imprimir.
*   **Convertir a Factura:** Usa el botón de cohete (🚀) para transformar esa cotización en una factura real. El sistema te llevará a la pantalla de emisión con todos los datos pre-cargados.

## 3. Emitir Factura
El proceso de venta simplificado.
1.  **Ambiente:** Un switch en la cabecera te permite cambiar entre `PRUEBAS` (Ámbar) y `PRODUCCIÓN` (Verde) al instante.
2.  **Cliente (Búsqueda Inteligente):** Escribe el nombre, RUC o email. El buscador tolerante a fallos encontrará al cliente aunque cometas errores tipográficos.
3.  **Productos:** Busca por nombre, código SKU o incluso por precio.
    *   El cálculo de impuestos (15%, 5%, 0%) es automático.
4.  **Emitir:** Al hacer clic en "Firmar y Emitir", el sistema:
    *   Genera el XML y lo firma digitalmente.
    *   Lo envía al SRI y genera el PDF.
    *   Envía el correo al cliente.

## 4. Historial
Tu archivo digital con **Búsqueda Global**.
*   **Barra de Búsqueda Inteligente:** Encuentra transacciones escribiendo el nombre del cliente, el número de factura, el RUC o incluso el monto total (ej. "50.00").
*   **Acciones Rápidas:**
    *   📄 **PDF:** Abre el RIDE visualmente.
    *   ✉️ **Email:** Reenvía la factura al cliente con un clic.
    *   🌐 **XML:** Abre el archivo fuente.
    *   📂 **Carpeta:** Te lleva a la ubicación física del archivo en tu disco.

## 5. Inventario y Clientes
Gestión básica de tus bases de datos (CRUD).
*   **Búsqueda Fuzzy:** Encuentra productos o clientes rápidamente usando términos aproximados.
*   Puedes Crear, Editar y Eliminar productos o clientes.
*   Los cambios se reflejan inmediatamente en todo el sistema.