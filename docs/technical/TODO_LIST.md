# Lista de Tareas (TODO)

## ✅ Completado (Sprint Final - Q1 2026)

### Core & SRI
- [x] **Validación XML:** Estructura `InfoAdicional`, `Detalles`, `Tarifa` corregida.
- [x] **Normativa:** Soporte RIMPE, Agente de Retención y Reglas de Pago (>$1000).
- [x] **Precisión:** Implementación de redondeo estricto a 2 decimales.

### Funcionalidad de Correo
- [x] **SMTP Local:** Implementación completa de envío nativo.
- [x] **Plantillas:** HTML profesional con resumen de factura.
- [x] **Pruebas:** Botón "Test Connection" y validación de credenciales.
- [x] **Independencia:** Eliminación de fallback a API Cloud para envíos.

### Interfaz y Experiencia (UI/UX)
- [x] **Personalización:** Subida y redimensionado de Logo.
- [x] **Dashboard:** Filtros de fecha dinámicos y gráficos reactivos.
- [x] **Notificaciones:** Toasts apilables y Centro de Historial (Campana).
- [x] **Navegación:** Reordenamiento lógico del Sidebar.
- [x] **Validación:** Feedback visual (bordes rojos) en formularios incompletos.

### Gestión de Datos
- [x] **Backups:** Módulo para listar y crear respaldos ZIP.
- [x] **Auditoría:** Tabla `MailLog` y visualización en pestaña "Actividad".
- [x] **Búsqueda:** Dropdown de clientes optimizado con scroll.

## 🚀 Mantenimiento Futuro (Post-Entrega)

### Testing & QA
- [x] **Cobertura:** Tests unitarios para servicios Core (`Quotation`, `Search`, `Chart`, `Product`).
- [x] **Integración:** Validar flujos completos (Cotización -> Factura).
- [x] **Importador Masivo:** Permitir carga de productos desde CSV.
- [ ] **Multi-empresa:** Soporte para gestionar múltiples RUCs en la misma instalación.
