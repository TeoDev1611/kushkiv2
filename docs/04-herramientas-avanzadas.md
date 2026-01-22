---
id: herramientas
title: Herramientas Avanzadas
---

# Herramientas de Auditoría y Seguridad

Para el usuario avanzado o administrador, el sistema ofrece herramientas de control total.

## 1. Panel de Actividad (Auditoría)
Ubicado en la pestaña **"Actividad"** (Icono de gráfico 📈). Este panel se divide en dos secciones críticas:

### A. Historial de Correos ✉️
Aquí verás un registro detallado de cada intento de envío de correo.
*   **Estado:** `SUCCESS` (Verde) o `FAILED` (Rojo).
*   **Detalle:** Si falla, te dirá exactamente por qué (ej. "Contraseña incorrecta", "Host no encontrado").
*   **Fecha:** Hora exacta del envío.

### B. Logs del Sistema ⚙️
Registra la "conversación" técnica con el SRI.
*   Ideal para depuración.
*   Muestra el JSON exacto de la petición y la respuesta del servidor del SRI.
*   Útil si una factura es rechazada por motivos tributarios complejos.

## 2. Centro de Respaldos (Backups)
Ubicado en la pestaña **"Respaldos"** (Icono de disquete 💾).

Tus datos son lo más importante. Este módulo te permite:
*   **Ver Historial:** Lista de todos los respaldos generados anteriormente con su peso y fecha.
*   **Generar Respaldo:** El botón "Crear Respaldo Ahora" comprime:
    1.  Tu base de datos (`kushki.db`).
    2.  Todas tus carpetas de facturas (XMLs y PDFs).
*   El resultado es un archivo `.zip` listo para guardar en una nube externa o USB.

## 3. Notificaciones del Sistema
En la cabecera superior derecha (icono 🔔), encontrarás el centro de notificaciones de la sesión.
*   Guarda un historial temporal de lo que ha sucedido mientras usabas la app (ej. "Factura enviada", "Error de conexión").
*   Te permite revisar mensajes que desaparecieron de la pantalla (Toasts) si te los perdiste.
