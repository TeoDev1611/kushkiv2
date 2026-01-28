---
title: Instalación y Configuración
description: Configura tu entorno de facturación en menos de 2 minutos.
sidebar:
  order: 2
---

# Puesta en Marcha

El sistema está diseñado para ser "Instalar y Olvidar". Una vez configurado, rara vez tendrás que volver a esta sección.

## 1. Activación Segura (Node-Locked)

Para proteger tu inversión, el sistema se vincula a tu hardware.
1.  Al abrir la app, verás la pantalla de bloqueo.
2.  Ingresa tu licencia: `KSH-XXXX-XXXX-XXXX`.
3.  El sistema validará y desencriptará la base de datos local.

## 2. Asistente de Configuración (Wizard)

Si es tu primera vez, el sistema detectará que faltan datos y abrirá el **Wizard**.

### Datos Tributarios
Llena estos datos con precisión, ya que irán firmados en cada XML.
*   **RUC & Razón Social:** Tal como constan en tu ficha del SRI.
*   **Obligado a Contabilidad:** Marca la casilla si aplica.

### Firma Electrónica (.p12)
El corazón de la facturación.
*   Selecciona tu archivo `.p12` o `.pfx`.
*   Ingresa la contraseña. El sistema intentará abrir el archivo en segundo plano para verificar que la clave sea correcta.

### Configuración de Correo (SMTP)
Kushki actúa como tu propio servidor de correos. Esto garantiza que el remitente sea **TU empresa**, no un tercero.

| Proveedor | Host | Puerto | Requisito Especial |
| :--- | :--- | :--- | :--- |
| **Gmail** | `smtp.gmail.com` | `587` | Requiere "Contraseña de Aplicación" (2FA) |
| **Outlook** | `smtp.office365.com` | `587` | Usa tu contraseña normal o de aplicación |

:::note[¿Por qué configurar mi propio correo?]
Al usar tu propio SMTP, evitas caer en la carpeta de SPAM de tus clientes, ya que el correo sale legítimamente desde tu cuenta y no desde un servidor masivo de facturación.
:::

## 3. Personalización de Marca

Ve a la pestaña **Configuración** (`Ctrl + 8`) para:
*   **Subir Logo:** El sistema redimensiona automáticamente tu imagen para optimizar el peso del PDF (RIDE).
*   **Ambiente:** Cambia entre `PRUEBAS` y `PRODUCCIÓN` con un solo clic.

### Guardado Contextual
Recuerda que puedes usar el atajo `Ctrl + S` en el panel de configuración para guardar los cambios inmediatamente.