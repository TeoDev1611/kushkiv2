---
title: Punto de Venta (POS) y Satélite Móvil
description: Guía de uso del módulo de facturación rápida y gestión de inventario remoto.
---

# 🛒 Punto de Venta (POS)

El módulo POS está diseñado para **ventas rápidas en mostrador** (Retail), optimizado para el uso de lectores de códigos de barras y teclado, minimizando el uso del mouse.

## Interfaz y Atajos

| Tecla | Acción | Descripción |
| :--- | :--- | :--- |
| **F5** | 🔍 Buscar | Abre el buscador manual de productos si no tiene el código a mano. |
| **F12** | 💳 Cobrar | Procesa la venta inmediatamente y genera la factura/nota de venta. |
| **ESC** | ❌ Cancelar | Limpia la pantalla actual o cierra ventanas modales. |
| **+ / -**| 🔢 Cantidad | Aumenta o disminuye la cantidad del último ítem añadido. |

## Flujo de Venta Rápida

1.  **Escanear:** Use su lector de códigos de barras. El producto se agregará automáticamente al carrito.
2.  **Ajustar:** Si el cliente lleva varios del mismo, escanee varias veces o use los botones `+` / `-`.
3.  **Cliente:** Por defecto es "CONSUMIDOR FINAL".
    *   Haga clic en el icono 👤 para **buscar** un cliente registrado o **crear uno nuevo** (+ Nuevo) sin salir de la venta.
4.  **Cobrar:** Presione `F12`. El sistema emitirá el documento electrónico y limpiará la pantalla para el siguiente cliente.

> **Tip:** ¿Necesita vincular un celular rápido? Use el botón 📱 en la parte superior del POS para ver el código QR de conexión sin ir a Configuración.

---

# 📱 Satélite Móvil (Inventario Remoto)

Convierta su teléfono celular en una extensión de Kushki. Ideal para realizar inventarios en bodega, ajustar stock en percha o verificar precios sin ir al computador.

## ¿Cómo conectar mi celular?

1.  Vaya a **Configuración** > **📱 Satélite Móvil**.
2.  Asegúrese de que su PC y su celular estén conectados a la **misma red Wi-Fi**.
3.  Escanee el **Código QR** que aparece en pantalla con la cámara de su celular.
4.  ¡Listo! Su celular se conectará automáticamente.

## Funciones del Satélite

*   **Buscador en tiempo real:** Escriba el nombre o código de un producto para ver su stock actual.
*   **Ajuste de Stock:** Toque cualquier producto para abrir el editor rápido.
    *   Use `+1` / `-1` para ajustes finos.
    *   Use `+10` / `-10` para ingresos masivos.
*   **Sincronización:** Cualquier cambio que haga en el celular se reflejará **instantáneamente** en la pantalla del computador (POS y Lista de Productos).

## Solución de Problemas

### 1. El celular no conecta (Timeout / Cargando infinito)
La causa más común es el **Firewall** del computador bloqueando el puerto `8085`.

**Solución para Linux:**
Abra una terminal y ejecute el comando según su distribución:
*   **Ubuntu / Linux Mint / Debian:**
    ```bash
    sudo ufw allow 8085/tcp
    ```
*   **Fedora / CentOS:**
    ```bash
    sudo firewall-cmd --zone=public --add-port=8085/tcp --permanent
    sudo firewall-cmd --reload
    ```

**Solución para Windows:**
1.  Abra "Seguridad de Windows" > "Firewall y protección de red".
2.  Seleccione "Permitir una aplicación a través del firewall".
3.  Busque `kushki.exe` y marque las casillas "Privada" y "Pública".

### 2. QR Incorrecto
Si tiene Docker o VPNs instalados, el código QR podría generar una IP interna (ej: `172.17.x.x`).
*   Vaya a **Configuración > Satélite Móvil**.
*   Edite manualmente el campo de IP con la dirección real de su PC (ej: `192.168.1.50`).
*   El QR se actualizará automáticamente.
