---
id: instalacion
title: Instalación y Configuración
---

# Instalación y Configuración Inicial

Al iniciar **Kushki Facturador** por primera vez, el sistema te guiará a través de un proceso seguro para garantizar que tu entorno de facturación esté listo.

## 1. Activación de Licencia

El sistema utiliza protección **Node-Locked**.
1.  Al abrir la app, verás una pantalla de bloqueo.
2.  Ingresa tu clave de producto (Formato: `KSH-XXXX-XXXX-XXXX`).
3.  El sistema validará tu hardware contra el servidor de licencias.
4.  Si es exitoso, accederás al **Asistente de Configuración (Wizard)**.

## 2. Asistente de Configuración (Wizard)

Este asistente de 4 pasos es obligatorio la primera vez.

### Paso 1: Datos de Empresa
Ingresa los datos tributarios tal como constan en tu RUC.
*   **RUC:** 13 dígitos obligatorios.
*   **Dirección Matriz:** La dirección fiscal principal.
*   **Opcionales:** Si eres Contribuyente RIMPE o Agente de Retención, llena estos campos para que aparezcan en el XML/PDF.
*   **Logo:** Haz clic en el botón de cámara (📷) y selecciona tu logo (PNG/JPG). El sistema lo ajustará automáticamente.

### Paso 2: Firma Electrónica
El "pasaporte" de tus facturas.
*   **Archivo .p12:** Selecciona tu archivo de firma electrónica.
*   **Contraseña:** La clave de tu firma. El sistema valida inmediatamente si es correcta.

### Paso 3: Almacenamiento
*   Define dónde se guardarán tus facturas.
*   Por defecto, el sistema crea una estructura organizada por `Año/Mes` dentro de tu carpeta de usuario.

### Paso 4: Correo Electrónico (SMTP)
Configura cómo se enviarán las facturas a tus clientes.
*   **Botones Rápidos:** Usa "Gmail" o "Outlook" para pre-llenar los servidores.
*   **Contraseña:** Si usas Gmail, recuerda usar una **Contraseña de Aplicación**, no tu clave personal.

---

## Modificar Configuración

Si necesitas cambiar algo después (ej. actualizaste tu firma):
1.  Ve a la pestaña **Configuración** en el menú lateral.
2.  Modifica los datos necesarios.
3.  Usa el botón **"Probar Conexión"** en la sección de correo para verificar que todo funcione.
4.  Haz clic en "Guardar Toda la Configuración".
