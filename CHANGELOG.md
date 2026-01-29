# Historial de Cambios (Changelog)

Todas las modificaciones notables al proyecto "Kushki Facturador" se documentarán en este archivo.

## [2.6.0] - 2026-01-28

### ✨ Características Nuevas (POS & Satélite)
- **POS - Selección Avanzada de Clientes:**
    - Implementado modal de búsqueda de clientes en el punto de venta.
    - Añadida funcionalidad de "Registro Rápido" (+ Nuevo) para crear clientes sin salir del flujo de venta.
- **POS - Acceso Rápido Satélite:**
    - Añadido botón QR (📱) en la cabecera del POS para vinculación rápida de dispositivos móviles.
- **Mejora UX - Código QR:**
    - El código QR ahora codifica la URL directa de conexión (ej: `http://192.168.x.x:8085/?token=...`) para apertura inmediata en móviles.

### 🛠️ Refactorización
- **Componente ClientForm:** Extraída la lógica de formulario de clientes a `src/lib/components/ClientForm.svelte` para reutilización.

### 🐛 Correcciones
- **Reset de Formulario:** Corregido bug donde el formulario de "Nuevo Cliente" en el POS mantenía datos antiguos.

## [2.5.0] - 2026-01-21

### Añadido
- **Módulo de Respaldos:** Nueva pestaña para gestionar copias de seguridad locales (.zip).
- **Auditoría de Correos:** Registro en base de datos (`MailLog`) de todos los envíos SMTP con estado y fecha.
- **Centro de Notificaciones:** Panel lateral (Campana) que guarda el historial de eventos de la sesión.
- **Soporte de Logo:** Capacidad de subir imágenes (JPG/PNG), redimensionarlas automáticamente y mostrarlas en PDF y Correos.
- **Plantillas HTML:** Correos electrónicos con diseño corporativo y tabla de resumen.
- **Validación Visual:** Formularios con feedback de error en tiempo real (bordes rojos, mensajes claros).

### Cambiado
- **Motor de Correo:** Se eliminó la dependencia de la API Cloud. Ahora se utiliza exclusivamente SMTP Local configurado por el usuario.
- **Dashboard:** Las estadísticas ahora responden a un selector de rango de fechas (Desde/Hasta).
- **PDF (RIDE):** Rediseño completo con colores corporativos ("Mint"), tarjetas de información y soporte de logo.
- **XML SRI:** Ajustes estrictos para cumplir con esquema XSD 1.1.0 (Corrección de `tarifa`, `detallesAdicionales`).
- **Navegación:** Reordenamiento del menú lateral para seguir el flujo de trabajo lógico.

### Corregido
- **Error SQL:** Solucionado el error "incomplete input" en consultas de estadísticas mediante separación de cláusulas `Where`.
- **UI:** Arreglado el desbordamiento del buscador de clientes (ahora tiene scroll interno).
- **Carga Inicial:** Implementado timeout de seguridad para evitar "Pantalla Negra" si el backend tarda en responder.
- **Wails Bindings:** Solucionado conflicto de tipos `time.Time` mediante DTOs.

## [2.0.0] - 2026-01-15
- Lanzamiento inicial de la versión de escritorio con arquitectura híbrida.
- Firma electrónica XAdES-BES nativa.
- Base de datos SQLite local.
