---
title: Herramientas Avanzadas
description: Control total sobre tus datos: Auditoría, Sincronización y Backups.
sidebar:
  order: 4
---

# Centro de Control y Seguridad

Ubicado en el módulo de **Sincronización** (`Ctrl + 7`), este panel te da control total sobre lo que ocurre "bajo el capó" del sistema.

## 1. Logs y Auditoría

La transparencia es clave. Aquí puedes ver:

### 📧 Logs de Correo
¿Un cliente dice que no recibió la factura?
*   Revisa este log.
*   Verás el **Estado Exacto** (Enviado, Fallido, Rebotado) y la fecha precisa.
*   Si falló (ej. "Password incorrect"), el sistema te lo dirá aquí.

### ☁️ Logs de Sincronización SRI
Historial técnico de la comunicación con el Servicio de Rentas Internas. Útil para contadores que necesitan saber por qué una factura específica fue "DEVUELTA" (ej. errores de validación XML).

## 2. Gestión de Respaldos (Backups)

Tu información es tu activo más valioso. Kushki v2 facilita su protección.

*   **Crear Respaldo Ahora:** Con un solo clic, el sistema:
    1.  Cierra temporalmente la base de datos para asegurar integridad.
    2.  Comprime la base de datos `kushki.db`.
    3.  Empaqueta todos los XMLs y PDFs generados.
    4.  Genera un archivo `.zip` con fecha y hora.
*   **Restauración:** Simplemente descomprime ese archivo en tu carpeta de instalación en caso de cambiar de computadora.

## 3. Sincronización Manual

Aunque el sistema sincroniza automáticamente cada vez que emites una factura, a veces necesitas forzar una actualización (ej. si trabajaste offline todo el día).

*   Botón **"Sincronizar SRI"**: Fuerza el reenvío de todos los comprobantes que estén en estado `PENDIENTE` o `FIRMADO` pero no `AUTORIZADO`.

:::note[Local-First]
Recuerda que **tú eres el dueño de tus datos**. No están en nuestra nube. Hacer respaldos periódicos es tu responsabilidad y tu mejor seguro.
:::
