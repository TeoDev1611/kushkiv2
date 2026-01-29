# 📱 Plan de Trabajo: Soluciones y Nuevas Funcionalidades (Móvil)

Este documento detalla la solución técnica para el error de cámara y la especificación para la creación de productos desde el celular.

## 1. 🚨 Solución Crítica: Error "Camera Streaming Not Supported"

**El Problema:**
Al intentar abrir la cámara en el celular (Brave/Chrome/Safari), aparece el error `camera streaming not supported by browser`.

**La Causa Raíz:**
Los navegadores modernos bloquean el acceso a la API de la cámara (`navigator.mediaDevices.getUserMedia`) si la página se sirve a través de **HTTP** inseguro. Solo permiten la cámara en:
1.  `localhost` (El PC sí funciona).
2.  `https://` (Contexto Seguro).

Como el celular accede vía `http://192.168.x.x:8085`, el navegador bloquea la cámara por seguridad. **Cambiar de librería JS (Node package) no solucionará esto**; es una restricción del navegador.

**La Solución (Para implementar mañana):**
Convertir el servidor local de Go en un servidor **HTTPS** con certificados autofirmados generados al vuelo.

### Pasos de Implementación:
1.  **Generación de Certificados en Go:**
    *   Usar `crypto/x509` y `math/big` en `app.go` para generar un certificado SSL temporal en memoria al iniciar la app.
2.  **Activar TLS en Echo:**
    *   Cambiar `e.Start(...)` por `e.StartTLS(...)` usando el certificado generado.
3.  **Experiencia de Usuario:**
    *   El usuario escaneará el QR.
    *   El navegador mostrará una advertencia: *"La conexión no es privada"* (porque el certificado es local).
    *   El usuario deberá dar clic en "Avanzado" -> "Continuar a sitio no seguro".
    *   **Resultado:** La cámara funcionará perfectamente.

---

## 2. ✨ Nueva Funcionalidad: Crear Productos desde el Móvil

**Objetivo:**
Permitir que el bodeguero registre productos nuevos directamente en la estantería si encuentra uno que no existe en el sistema.

### Interfaz de Usuario (UI)
*   **Ubicación:** Botón flotante "Nuevo" (+) en la pantalla de Inventario (junto a la barra de búsqueda).
*   **Formulario Modal:**
    1.  **Código de Barras:** Campo con botón de escáner pequeño. (Si no tiene código, se genera un SKU interno automáticamente).
    2.  **Nombre:** Input de texto.
    3.  **Precio (PVP):** Input numérico.
    4.  **Stock Inicial:** Input numérico.
    5.  **Ubicación:** Input de texto.

### Backend (Go)
*   **Endpoint:** `POST /api/product/create`
*   **Lógica:**
    *   Validar que el SKU/Barcode no exista.
    *   Crear registro en DB con impuestos por defecto (IVA 15%).
    *   Emitir evento `inventory-updated` para que aparezca en el PC.

### Flujo de Trabajo
1.  Usuario está acomodando mercadería.
2.  Encuentra una caja nueva sin registrar.
3.  Abre la App Móvil -> Toca "+".
4.  Escanea el código de la caja.
5.  Ingresa "Pack Galletas x12", Precio "$2.50", Stock "20".
6.  Guardar.
7.  El producto queda disponible inmediatamente para la venta en el POS.
