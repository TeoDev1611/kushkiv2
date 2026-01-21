# Kushki Facturador v2.0 🚀

![Status](https://img.shields.io/badge/Estado-Producción_Q1_2026-success)
![Tech](https://img.shields.io/badge/Stack-Go_Wails_Svelte-blue)
![Security](https://img.shields.io/badge/Licencia-Node_Locked-orange)

Sistema de facturación electrónica de escritorio para Ecuador, diseñado para alta eficiencia, seguridad robusta y experiencia de usuario moderna ("Obsidian & Mint").

## ✨ Características Principales

*   **Arquitectura Híbrida Segura:** Aplicación de escritorio (Go/Wails) con validación de licenciamiento y servicios en la nube (API Deno).
*   **Licenciamiento Node-Locked:** El sistema se vincula al hardware específico del usuario, impidiendo copias no autorizadas.
*   **Emisión "Zero-Config" de Correos:** Envío de comprobantes PDF vía API Cloud, eliminando la compleja configuración de SMTP local para el usuario.
*   **Firma Electrónica Nativa:** Implementación pura en Go (XAdES-BES) sin dependencias de Java o librerías externas pesadas.
*   **Dashboard en Tiempo Real:** Métricas de ventas, estado del SRI y tendencias gráficas.
*   **Base de Datos Local:** SQLite con GORM para persistencia rápida y segura de comprobantes.
*   **Modo Offline Resiliente:** Permite facturar y firmar localmente (la sincronización requiere internet).

## 🛠️ Stack Tecnológico

| Componente | Tecnología | Descripción |
| :--- | :--- | :--- |
| **Frontend** | Svelte + Vite | Interfaz reactiva, rápida y ligera. |
| **Backend Desktop** | Go 1.24 (Wails) | Lógica de negocio, firma XML, base de datos. |
| **Cloud API** | Deno (Oak) | Microservicio de Licenciamiento y Envío de Correos. |
| **Database** | SQLite + GORM | Almacenamiento local de facturas y configuración. |
| **Reportes** | Maroto (PDF) | Generación de RIDE vectorial de alta calidad. |

## 🚀 Instalación y Uso

### Prerrequisitos
*   Go 1.21+
*   Node.js 18+
*   Wails v2 (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Ejecución en Desarrollo
```bash
wails dev
```

### Compilación (Producción)
```bash
wails build
```

## 🔒 Flujo de Seguridad y Activación

1.  **Instalación:** Al abrir la app por primera vez, se mostrará el **Panel de Activación**.
2.  **Activación:** El usuario ingresa su Clave de Producto. El sistema genera un `MachineID` único y lo valida contra la Nube.
3.  **Configuración:** Si la activación es exitosa, se inicia el **Asistente de Configuración** (Wizard) obligatorio para cargar RUC y Firma Electrónica.
4.  **Uso:** El Dashboard se desbloquea solo con una licencia válida y configuración completa.

## 📂 Estructura del Proyecto

```
kushkiv2/
├── app.go                 # Controlador principal (Bridge Frontend-Backend)
├── frontend/              # UI Svelte
├── internal/
│   ├── db/                # Modelos GORM y Migraciones
│   └── service/
│       ├── cloud_service.go  # Cliente API Deno (Licencias/Email)
│       ├── invoice_service.go # Lógica de Facturación SRI
│       └── report_service.go  # Generación PDF/Excel
└── pkg/                   # Librerías Core (Firma XAdES, XML, SRI SOAP)
```

## 📝 Licencia

Este software es propietario y requiere una licencia comercial activa para su funcionamiento. Protegido por sistema de validación de hardware.