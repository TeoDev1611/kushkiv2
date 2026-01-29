# Ficha Técnica: Módulo POS y Satélite Móvil

## 1. Arquitectura del Servidor Satélite

El sistema implementa un servidor HTTP híbrido dentro del mismo proceso de la aplicación de escritorio (Wails/Go). Esto elimina la necesidad de instalar software adicional en el servidor o en los clientes móviles.

### Componentes:
*   **Servidor HTTP:** **Echo Framework (v4)** escuchando en `0.0.0.0:8085`.
*   **Routing:** `Echo Router` para mejor performance y manejo de errores.
*   **Middleware:**
    *   `Logger`: Registro de accesos en consola para debugging.
    *   `CORS`: Habilitado para compatibilidad móvil.
    *   `Recover`: Prevención de crashes por pánicos.
    *   `Auth`: Middleware personalizado para validación de `X-Kushki-Token`.
*   **Assets:** La Web App móvil se sirve mediante `e.StaticFS` sobre `embed.FS`.

### Diagrama de Flujo:
```mermaid
sequenceDiagram
    participant Mobile as WebApp (Mobile)
    participant Server as Go Server (PC)
    participant Desktop as Wails Frontend (PC)

    Note over Mobile, Server: Conexión Inicial (QR)
    Mobile->>Server: GET / (con Token en URL)
    Server-->>Mobile: Servir index.html (Embed)

    Note over Mobile, Server: Operativa
    Mobile->>Server: POST /api/stock {sku, qty}
    Server->>DB: Actualizar SQLite
    Server->>Desktop: runtime.EventsEmit("inventory-updated")
    Desktop->>Desktop: Actualizar UI (Svelte Store)
```

## 2. API Satélite

### Endpoints
| Método | Endpoint | Descripción | Auth |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/inventory` | Retorna lista completa de productos (DTO). | Token |
| `POST` | `/api/stock` | Actualiza stock. Body: `{sku, quantity, type}`. | Token |
| `GET` | `/api/status` | Health check. | Pública |

### Seguridad
*   **Token:** PIN numérico de 6 dígitos generado aleatoriamente en cada `App.startup`.
*   **Alcance:** La API solo es accesible dentro de la red local (LAN).
*   **CORS:** No configurado explícitamente (Same-origin policy aplica al servir los assets desde el mismo puerto).

## 3. Web App Móvil
*   **Stack:** Vanilla JS, HTML5, CSS3. (Sin frameworks ni build steps complejos).
*   **Ubicación:** `internal/mobile/static/`.
*   **Estado:** `localStorage` se usa para persistir el Token entre recargas.

## 4. Sincronización (Real-time)
Se utiliza el bus de eventos de Wails (`runtime.EventsEmit` y `EventsOn`) para notificar cambios desde el backend (Go) hacia el frontend de escritorio (Svelte).

*   **Evento:** `inventory-updated`
*   **Payload:** Objeto `Product` actualizado.
*   **Listeners:**
    *   `ProductList.svelte`: Actualiza la celda de stock.
    *   `PosView.svelte`: Actualiza resultados de búsqueda y valida carrito.

## 5. Interfaz POS (Punto de Venta)

El módulo de Punto de Venta ha sido mejorado para facilitar la operación rápida y la integración con el ecosistema móvil.

### Selección de Clientes
*   **Búsqueda Rápida:** Modal integrado para buscar clientes por Nombre o RUC/Cédula.
*   **Creación Rápida:** Botón "+ Nuevo" que despliega un formulario simplificado para registrar clientes sin abandonar la pantalla de ventas.
*   **Persistencia:** El cliente seleccionado se mantiene en el contexto de la venta actual hasta que se finaliza o se limpia.

### Integración Móvil
*   **Acceso Directo:** Botón "📱" en la cabecera del POS.
*   **Vinculación:** Despliega el código QR de conexión (el mismo de la configuración) directamente en el POS, permitiendo vincular dispositivos auxiliares (como verificadores de precios o inventariadores) al instante.

