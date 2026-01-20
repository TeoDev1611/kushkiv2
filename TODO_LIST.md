# Lista de Tareas - Implementación Técnica (Basada en TECH.md v1.1.0)

Este documento desglosa los requerimientos arquitectónicos de `TECH.md` en tareas ejecutables, separadas por módulos y prioridades.

---

## 1. Validación del Núcleo (Core Verification)
*Asegurar que la implementación actual cumpla con las definiciones estrictas de los capítulos 3, 4 y 6.*

- [x] **Persistencia (GORM/SQLite)**:
    - [x] Verificar que `EmisorConfig` incluya los campos de correo (`SMTPHost`, `SMTPUser`, `SMTPPass`) y `StoragePath`.
    - [x] Confirmar que `P12Password` se esté guardando cifrada (AES-GCM) y no en texto plano.
- [x] **Lógica de Facturación**:
    - [x] Auditar `EmitirFactura` para confirmar que el proceso ocurre en una Goroutine (no bloqueante).
    - [x] Verificar que la generación de Clave de Acceso usa la implementación "Zero-allocation" descrita en `pkg/util/mod11.go`.
- [x] **Seguridad**:
    - [x] Confirmar implementación de `crypto/aes` para datos sensibles.

---

## 2. Módulo de Reportería y Contabilidad (Cap. 9)
*Implementación del servicio de exportación usando `excelize`.*

- [x] **Infraestructura**:
    - [x] Instalar dependencia: `go get github.com/xuri/excelize/v2`.
    - [x] Crear estructura `internal/service/report_service.go`.
- [x] **Reporte ATS (Excel)**:
    - [x] Implementar método `GenerateSalesExcel(startDate, endDate time.Time)`.
    - [x] Crear consulta GORM con filtro de fechas (`BETWEEN`).
    - [x] Mapear resultados a filas de Excel optimizadas (Stream writer si el volumen es alto).
- [ ] **Visualización (Dashboard)**:
    - [ ] Implementar método `GetTopProducts()` que retorne JSON para Svelte.
    - [ ] Integrar librería de gráficos (Chart.js o similar) en `frontend/src/pages/Dashboard.svelte`.

---

## 3. Gestión de Correos y Notificaciones (Cap. 10)
*Sistema SMTP asíncrono con cola de persistencia.*

- [ ] **Base de Datos**:
    - [ ] Crear modelo `EmailQueue` en `internal/db/models.go` (ID, To, Subject, Body, Status, RetryCount, CreatedAt).
    - [ ] Ejecutar migración automática.
- [ ] **Backend (Go)**:
    - [ ] Implementar `SMTPService` con `gopkg.in/gomail.v2` o `net/smtp`.
    - [ ] Crear **Worker** (Goroutine) que inicie al arrancar la app y revise la cola cada 60s.
    - [ ] Implementar lógica de reintentos (Backoff exponencial) para correos fallidos.
    - [ ] Crear templates HTML básicos para el cuerpo del correo.
- [ ] **Frontend**:
    - [ ] Añadir opción "Reenviar Correo" en el historial de facturas.

---

## 4. Sistema de Respaldo (Cap. 11)
*Protección de datos locales.*

- [x] **Lógica de Backup**:
    - [x] Crear función `CreateBackup()` en `app.go`.
    - [x] Implementar compresión `.zip` nativa (`archive/zip`).
    - [x] Incluir en el zip: Archivo `kushki.db` y carpetas de XML/PDFs generados.
- [x] **Automatización**:
    - [x] Vincular `CreateBackup` al evento `OnShutdown` (cierre de aplicación) en `main.go`.
    - [x] Añadir configuración en frontend para definir ruta de destino del backup.

---

## 5. Módulo de Contingencia Offline (Cap. 12)
*Resiliencia ante fallos de red.*

- [x] **Detección de Fallos**:
    - [x] Modificar `soap_client.go` para distinguir entre "Error de Validación" (Rechazo SRI) y "Error de Red" (Timeout/DNS).
- [x] **Manejo de Estado**:
    - [x] Si es error de red, guardar factura con estado `PENDIENTE_ENVIO` (no `ERROR`).
    - [x] Permitir generar XML firmado y RIDE incluso sin conexión.
- [x] **Worker de Sincronización**:
    - [x] Crear `SyncService` que corra en background.
    - [x] Implementar "Ping" periódico al SRI.
    - [x] Al recuperar conexión: Buscar facturas `PENDIENTE_ENVIO`, enviarlas en lote y notificar vía Wails Events.