---
title: Configuración del Sistema
description: Pasos para configurar el emisor, la firma y el correo.
---

Antes de emitir tu primera factura, debes completar la configuración básica.

## 1. Datos del Emisor
Ingresa tu RUC (13 dígitos), Razón Social y Dirección Matriz. Asegúrate de seleccionar el **Ambiente** correcto:
- **Pruebas:** Para validar el sistema (documentos sin valor legal).
- **Producción:** Para emitir facturas oficiales autorizadas por el SRI.

## 2. Firma Electrónica
Carga tu archivo `.p12` y proporciona la contraseña. El sistema validará que la firma:
- No esté caducada.
- Corresponda al RUC ingresado.

## 3. Servidor de Correo (SMTP)
Configura tu cuenta de Gmail u Outlook para enviar automáticamente los PDFs a tus clientes. 
- **Gmail:** Requiere crear una "Contraseña de Aplicación".
- **Puerto:** Generalmente 587 (STARTTLS).

:::caution[Seguridad]
Nunca compartas tu archivo `.p12` ni tu contraseña con terceros. Kushki nunca te pedirá estos datos por internet.
:::
