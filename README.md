# Kushki Facturador v2.5 - Professional Edition 🚀

![Status](https://img.shields.io/badge/Estado-Producción_Estable-success)
![Version](https://img.shields.io/badge/Versión-2.5.0-blue)
![Licencia](https://img.shields.io/badge/Licencia-Node_Locked-orange)

Sistema de facturación electrónica de escritorio para Ecuador, diseñado con una arquitectura híbrida (Go + Svelte) que prioriza la autonomía del usuario, la seguridad de datos y una experiencia visual moderna.

## 🌟 Características Destacadas

### 🎨 Experiencia de Usuario (UX)
*   **Interfaz "Obsidian & Mint":** Diseño oscuro moderno con acentos visuales claros para estados (Éxito, Error, Pendiente).
*   **Dashboard Dinámico:** Gráficos de tendencias y KPIs filtrables por rangos de fecha personalizados.
*   **Feedback Visual:** Validaciones de formulario en tiempo real e indicadores de estado.

### 📧 Autonomía Total (SMTP Local)
*   **Sin Dependencias:** Envío de correos utilizando el servidor SMTP del propio usuario (Gmail, Outlook, Corporativo).
*   **Plantillas HTML:** Correos electrónicos profesionales con resumen de factura y adjunto PDF.
*   **Verificación:** Herramienta integrada para probar credenciales de correo.

### 🛡️ Seguridad y Auditoría
*   **Licenciamiento por Hardware:** Protección contra copias no autorizadas.
*   **Centro de Actividad:** Registro inmutable de cada correo enviado y evento de sincronización con el SRI.
*   **Respaldos Locales:** Generación de copias de seguridad (.zip) de base de datos y archivos XML/PDF con un clic.

### 🇪🇨 Cumplimiento SRI (2025-2026)
*   **Validaciones Estrictas:** Control de decimales, plazos de pago y montos máximos para consumidor final.
*   **Soporte Completo:** Manejo de RIMPE, Agentes de Retención y XML v1.1.0.
*   **Firma Electrónica:** Motor de firma XAdES-BES nativo (sin Java).

## 📚 Documentación

La documentación detallada se encuentra en la carpeta `docs/`:

1.  [Introducción y Alcance](docs/01-introduccion.md)
2.  [Instalación y Configuración](docs/02-instalacion-configuracion.md)
3.  [Manual de Interfaz](docs/03-interfaz-usuario.md)
4.  [Herramientas Avanzadas (Backups/Logs)](docs/04-herramientas-avanzadas.md)
5.  [Arquitectura Técnica](docs/05-arquitectura-tecnica.md)
6.  [Soporte Técnico y FAQ](docs/06-soporte-faq.md)

## 🛠️ Stack Tecnológico

*   **Core:** Go 1.24 (Backend), Wails v2 (Bridge).
*   **UI:** Svelte + Vite (Frontend).
*   **Datos:** SQLite + GORM (ORM).
*   **Reportes:** Maroto (PDF Engine).

## 🚀 Inicio Rápido (Desarrollo)

```bash
# Instalar dependencias
go mod tidy
cd frontend && npm install && cd ..

# Ejecutar en modo dev
wails dev
```

## 📦 Compilación (Producción)

```bash
wails build
```

---
**Desarrollado con ❤️ para Ecuador.**