# 📱 Plan de Modernización: Satélite Móvil v2.0

## 🎯 Objetivo
Transformar la actual "herramienta de ajuste de inventario" en un **Asistente Integral de Bodega y POS**. La aplicación móvil dejará de ser pasiva para convertirse en un controlador activo del sistema, respetando la identidad visual "Obsidian & Mint" del escritorio.

---

## 1. Rediseño de Interfaz (UI/UX)
**Meta:** Unificar la experiencia visual con la aplicación de escritorio.

### Estilo Visual "Obsidian & Mint"
*   **Paleta de Colores:**
    *   Fondo: `#1a1a1a` (Dark Grey) y `#000000` (Pure Black).
    *   Acento: `#34D399` (Mint) para acciones principales (Escanear, Guardar).
    *   Texto: Blanco (`#ffffff`) y Gris Terciario (`#9ca3af`).
    *   Bordes: Sutiles (`rgba(255,255,255,0.1)`).
*   **Tipografía:** *Nunito* (misma que Desktop).
*   **Layout Móvil:**
    *   **Header:** Estado de conexión y selector de modo (Inventario / POS).
    *   **Body:** Lista reactiva o visor de cámara.
    *   **Bottom Navigation Bar:** Navegación rápida entre:
        1.  📦 **Inventario** (Lista, Buscador).
        2.  📷 **Escanear** (Botón central flotante).
        3.  ⚙️ **Ajustes** (Desconectar, Info).

---

## 2. Nueva Funcionalidad: Escáner de Cámara
**Meta:** Eliminar la dependencia de pistolas lectoras bluetooth; usar la cámara del celular.

*   **Tecnología:** Librería `html5-qrcode` (Ligera, JS puro, funciona offline si se descarga).
*   **Modos de Uso:**
    1.  **Modo Consulta (Bodega):** Escanear un producto muestra su ficha técnica para editar stock, ubicación o precio.
    2.  **Modo POS (Caja):** Escanear un producto lo envía **directamente al carrito de ventas en la PC**.

---

## 3. Características Faltantes (Solicitadas)

### A. "Control Remoto" para el POS (El Factor "WOW")
El celular actuará como un periférico de entrada inteligente para la PC.
*   **Sincronización Instantánea:**
    *   **Escaneo:** Al leer un código con la cámara, el producto aparece **al instante** en la pantalla de ventas de la PC.
    *   **Cantidad (+1):** Si el producto ya está en la venta, pulsar `+1` en el celular aumentará la cantidad en la PC en tiempo real.
    *   **Feedback:** El celular vibrará/sonará confirmando que la PC recibió la orden.
*   **Flujo:**
    1.  Usuario selecciona "Modo POS" en el celular.
    2.  Escanea o pulsa botones de acción.
    3.  Celular envía `POST /api/pos/scan` o `/api/pos/update-qty`.
    4.  PC recibe el evento y actualiza la tabla de ventas ante los ojos del cliente.

### B. Gestión de Ubicación y Códigos
Solucionar el mensaje "Sin ubicación" y permitir asignar códigos nuevos.
*   **Editor de Producto Móvil:**
    *   Campo **Ubicación**: Editable (ej: "Pasillo 4, Estante B").
    *   Campo **Código de Barras**: Si un producto no tiene código, se puede escanear uno nuevo para asignárselo.

---

## 4. Arquitectura Técnica (Backend Go & Echo)

### Nuevos Endpoints en `app.go`
Se requiere ampliar la API REST del servidor incrustado:

| Verbo | Endpoint | Acción |
| :--- | :--- | :--- |
| `POST` | `/api/pos/scan` | Envía un código escaneado a la vista POS de escritorio. |
| `POST` | `/api/product/update` | Actualiza campos extendidos (Ubicación, Barcode, Precio). |
| `GET` | `/static/lib/html5-qrcode.min.js` | Servir la librería de escaneo localmente (Offline). |

### Comunicación Backend -> Frontend Desktop
*   Nuevo evento Wails: `remote-scan`.
*   Listener en `PosView.svelte`: Al recibir `remote-scan`, ejecuta la función `addProductByIdentifier`.

---

## 5. Plan de Ejecución (Paso a Paso)

1.  **Backend (Go):**
    *   Implementar endpoint `/api/pos/scan`.
    *   Implementar endpoint `/api/product/update`.
    *   Descargar y embeber `html5-qrcode.min.js`.

2.  **Frontend Móvil (HTML/CSS):**
    *   Reescribir `style.css` con variables CSS del tema Desktop.
    *   Reestructurar `index.html` con la barra de navegación inferior.

3.  **Lógica Móvil (JS):**
    *   Implementar el lector de cámara.
    *   Crear la lógica de "Modo POS" vs "Modo Inventario".
    *   Mejorar el formulario de edición (Ubicación).

4.  **Integración Desktop (Svelte):**
    *   Conectar `PosView.svelte` al evento de escaneo remoto.

---

## 6. Resultado Final Esperado
Un sistema donde el usuario puede estar en la bodega con su celular, escanear un producto, ver que le falta ubicación, corregirla ahí mismo, y luego ir a la caja, activar "Modo POS" en el celular y usarlo para despachar clientes rápidamente sin comprar lectores extra.
