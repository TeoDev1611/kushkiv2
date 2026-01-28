---
id: soporte-faq
title: Soporte Técnico y FAQ
---

# 🛠️ Soporte Técnico y Preguntas Frecuentes

Esta guía ayuda a resolver los problemas más comunes que pueden surgir durante el uso de **Kushki Facturador**.

## ❓ Preguntas Frecuentes (FAQ)

### 1. ¿Por qué mi factura dice "DEVUELTA" por el SRI?
Esto ocurre generalmente por errores de validación. Revisa la pestaña **Actividad**:
*   **Error 35:** Estructura XML inválida (revisa que los datos del cliente estén completos).
*   **Error 43:** Clave de acceso registrada (probablemente ya enviaste esa factura).
*   **RUC no existe:** Verifica que el RUC del emisor o comprador tenga 13 dígitos y sea válido.

### 2. ¿Cómo configuro mi correo de Gmail para enviar facturas?
Gmail requiere un paso de seguridad adicional:
1.  Activa la "Verificación en dos pasos" en tu cuenta de Google.
2.  Busca "Contraseñas de Aplicaciones".
3.  Genera una nueva contraseña para "Kushki App".
4.  Usa esa contraseña de 16 caracteres en la configuración SMTP, **no tu clave personal**.

### 3. El logo no aparece en el PDF, ¿qué hago?
Asegúrate de:
*   Haber guardado la configuración después de subir el logo.
*   Que el archivo original sea un formato válido (JPG o PNG).
*   Si cambiaste de computadora, recuerda volver a subir el logo para que el sistema genere la ruta local correcta.

### 4. ¿Puedo convertir una Cotización en Factura?
Sí. En la pestaña de Cotizaciones, busca la proforma que deseas facturar y haz clic en el botón de "Cohete" (🚀). Esto te llevará a la pantalla de emisión con todos los datos del cliente y productos ya cargados, listos para firmar.

### 5. ¿Dónde están mis archivos físicamente?
El sistema guarda todo en la ruta que elegiste en **Configuración -> Carpeta de Guardado**. Por defecto, se organizan así:
`Ruta/Año/Mes/FACTURA-000000XXX.pdf`

---

## 🚨 Resolución de Problemas (Troubleshooting)

### Error: "dial tcp: lookup smtp.gmail.com: no such host"
*   **Causa:** No tienes conexión a internet o el nombre del servidor está mal escrito.
*   **Solución:** Revisa tu Wi-Fi/Cable y verifica que no hayas escrito `smpt` en lugar de `smtp`. Usa los botones de pre-configuración.

### Error: "Archivo P12 inválido o contraseña incorrecta"
*   **Causa:** La contraseña de tu firma electrónica no coincide.
*   **Solución:** Vuelve a ingresar la contraseña en la pestaña **Configuración**. El sistema validará el archivo inmediatamente.

### La aplicación se queda en pantalla negra al iniciar
*   **Causa:** El backend está procesando una base de datos bloqueada o muy grande.
*   **Solución:** Espera 8 segundos; el sistema tiene un "Safety Timer" que forzará la carga de la interfaz automáticamente.

### Las estadísticas no se actualizan
*   **Causa:** Solo se suman las facturas con estado `AUTORIZADO`.
*   **Solución:** Si tus facturas están `PENDIENTES` o `DEVUELTAS`, no contarán para el total de ventas. Verifica el estado en el Historial.

---

## 📞 ¿Necesitas más ayuda?
Si el problema persiste, contacta al administrador del sistema adjuntando los archivos de log:
1.  Ejecuta la app en modo debug: `./kushkiv2 --kushki-debug`.
2.  Envía el archivo `kushki_app.log` generado en la raíz del programa.