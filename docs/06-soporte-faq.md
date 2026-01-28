---
title: Soporte y FAQ
description: Resolución de problemas comunes.
sidebar:
  order: 6
---

# Soporte Técnico

Aunque el sistema es robusto, entendemos que la facturación electrónica puede ser compleja debido a factores externos (SRI, Internet, Correo).

## Preguntas Frecuentes (FAQ)

### ¿Por qué no veo mis facturas antiguas?
El sistema filtra por defecto las facturas del **mes actual** para mantener la velocidad. Usa los selectores de fecha en la parte superior del Dashboard o Historial para ver meses anteriores.

### ¿Puedo instalarlo en varias computadoras?
La licencia es **Node-Locked** (por equipo). Si necesitas facturar desde múltiples puntos, necesitarás licencias adicionales o contactar a soporte para una licencia flotante.

### El sistema dice "SRI Offline"
Esto verifica la conexión con los servidores del gobierno. Si el SRI está caído (común a fin de mes), el sistema guardará tus facturas como `PENDIENTE` o `FIRMADO`. Usa el botón **"Sincronizar"** en el panel de control más tarde para enviarlas.

## Resolución de Problemas

| Síntoma | Causa Probable | Solución |
| :--- | :--- | :--- |
| **Error SMTP** | Contraseña incorrecta o bloqueo de seguridad. | Si usas Gmail, asegúrate de usar una **Contraseña de Aplicación**, no tu clave normal. Revisa el puerto (587). |
| **Pantalla Blanca** | Error de renderizado gráfico. | Actualiza los drivers de video o Webview2. El sistema usa aceleración por hardware. |
| **No puedo editar una factura** | Factura ya autorizada. | Una vez autorizada por el SRI, un documento es legalmente inmutable. Debes emitir una Nota de Crédito. |

## Reportar un Bug

Si encuentras un error técnico:
1.  Ve a la pestaña **Sincronización/Logs**.
2.  Revisa si hay errores en rojo.
3.  Ejecuta la aplicación desde la terminal con `./kushkiv2 --kushki-debug` para ver trazas detalladas.
4.  Contacta a soporte con esa información.
