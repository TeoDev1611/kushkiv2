---
title: Arquitectura Técnica
description: Detalles internos para desarrolladores y entusiastas de la tecnología.
sidebar:
  order: 5
---

# Arquitectura Moderna (v2.1.0)

Kushki Facturador v2 representa un salto cuántico respecto a las aplicaciones de escritorio tradicionales. Hemos adoptado una arquitectura **Híbrida de Alto Rendimiento**.

## Stack Tecnológico

| Capa | Tecnología | Versión | Ventaja |
| :--- | :--- | :--- | :--- |
| **Frontend** | **Svelte** | **4.2.9** | Reactividad quirúrgica sin Virtual DOM. UI ultra-rápida. |
| **Build Tool** | **Vite** | **5.0** | Compilación instantánea y módulos ES nativos. |
| **Lenguaje UI** | **TypeScript** | **5.3** | Tipado estricto para prevenir errores en tiempo de ejecución. |
| **Backend** | **Go (Golang)** | **1.21+** | Lógica de negocio compilada, binario nativo, concurrencia real. |
| **Bridge** | **Wails** | **v2** | Comunicación Frontend-Backend sin servidor HTTP, usando IPC nativo. |
| **Base de Datos** | **SQLite + GORM** | **Latest** | SQL robusto, local, en un solo archivo. |

## Estructura del Proyecto (Feature-Based)

Hemos refactorizado el código monolítico antiguo hacia una arquitectura modular basada en características (*Feature-based Architecture*). Esto facilita el mantenimiento y la escalabilidad.

```text
frontend/src/lib/
├── features/           # Módulos auto-contenidos (Vistas + Lógica)
│   ├── dashboard/      # Gráficos y KPIs
│   ├── invoice/        # Emisor, Lógica de Impuestos
│   ├── quotations/     # Gestión de Cotizaciones
│   ├── inventory/      # Productos
│   ├── clients/        # Directorio
│   ├── sync/           # Logs y Backups
│   └── settings/       # Configuración
├── stores/             # Gestión de Estado Global (Svelte Stores)
│   ├── app.ts          # Estado de la UI, Licencia
│   ├── invoice.ts      # Store complejo de la factura en curso
│   └── notifications.ts# Sistema de Toasts
└── services/           # Capa de Abstracción API
    └── api.ts          # Facade tipada para llamadas a Go (Wails)
```

## Ventajas para el Desarrollador/Usuario

1.  **Menor Consumo de RAM:** Al usar el motor de renderizado nativo del sistema operativo (WebView2 en Windows, WebKit en Linux/Mac) a través de Wails, la app pesa ~15MB en lugar de los ~150MB de una app Electron típica.
2.  **Type Safety:** Los modelos de datos de Go (`structs`) se convierten automáticamente a interfaces TypeScript (`models.ts`), garantizando que el frontend siempre sepa qué datos esperar.
3.  **Compilación Atómica:** El frontend es independiente. Podemos actualizar la interfaz visual sin tocar la lógica crítica de facturación en Go, y viceversa.

## Flujo de Datos

1.  **Interacción:** Usuario presiona `Ctrl + S` en *Clientes*.
2.  **Frontend:** `ClientDirectory.svelte` captura el evento, valida el formulario y llama a `Backend.saveClient()`.
3.  **Bridge:** Wails serializa los datos y los pasa a Go de forma segura.
4.  **Backend:** Go valida reglas de negocio, guarda en SQLite y retorna el nuevo ID.
5.  **Feedback:** Svelte recibe la respuesta y muestra una notificación *Toast* animada. Todo en milisegundos.
