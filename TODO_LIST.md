# Lista de Tareas (TODO)

## ✅ Completado
- [x] **Fase 1:** Estructura base Wails + Svelte.
- [x] **Fase 2:** Base de datos SQLite y Modelos GORM.
- [x] **Fase 3:** Generación y Firma de XML (XAdES-BES nativo).
- [x] **Fase 4:** Conexión SOAP con SRI (Recepción/Autorización).
- [x] **Fase 5:** Generación de RIDE (PDF) con Maroto.
- [x] **Fase 6:** Dashboard y Reportería Básica.
- [x] **Fase 7:** Seguridad y Cloud.
    - [x] Implementar `CloudService` para API Deno.
    - [x] Sistema de Licenciamiento Node-Locked.
    - [x] UI de Bloqueo por Licencia.
    - [x] Wizard de Configuración Inicial Obligatorio.
    - [x] Eliminación de SMTP Legacy.

## 🚀 Pendiente (Roadmap)

### Refinamiento Técnico
- [ ] **Validación JWT:** Verificar firma del token de licencia en el cliente Go para evitar spoofing simple.
- [ ] **Offline Mode Mejorado:** Cola de reintento para envío de correos cuando vuelva internet (actualmente solo se intenta una vez al emitir).

### Funcionalidades Usuario
- [ ] **Importador Masivo:** Carga de productos/clientes desde Excel.
- [ ] **Personalización:** Permitir subir logo de empresa para el RIDE.
- [ ] **Multi-usuario:** (Futuro) Roles básicos (Admin/Vendedor).

### Mantenimiento
- [ ] **Tests Unitarios:** Aumentar cobertura en `CloudService` y `InvoiceService`.
- [ ] **CI/CD:** Configurar GitHub Actions para builds automáticos.
