---
title: Interfaz y Atajos
description: Domina la interfaz Obsidian y los atajos de teclado para máxima productividad.
sidebar:
  order: 3
---

# Guía de Interfaz y Productividad

Hemos diseñado la interfaz siguiendo los principios de **"Teclado Primero"**. Si eres un contador o facturador intensivo, podrás operar el sistema casi sin tocar el mouse.

## ⌨️ Atajos de Teclado (Power User)

¡Aprende estos comandos y duplica tu velocidad!

| Atajo | Acción | Descripción |
| :--- | :--- | :--- |
| **`Ctrl + N`** | **Nueva Factura** | Salta inmediatamente al módulo de facturación, limpia el formulario y pone el foco en el cliente. |
| **`Ctrl + S`** | **Guardar** | Guarda lo que estés haciendo (Cliente, Producto, Configuración, Cotización). *Nota: En facturación te pedirá confirmación por seguridad.* |
| **`Esc`** | **Blur** | Quita el foco de cualquier campo de texto para que puedas usar los atajos de navegación. |
| **`Ctrl + 1`** | Ir a Dashboard | Resumen general. |
| **`Ctrl + 2`** | Ir a Facturar | Módulo de emisión. |
| **`Ctrl + 3`** | Ir a Cotizaciones | Gestión de proformas. |
| **`Ctrl + 4`** | Ir a Productos | Inventario. |
| **`Ctrl + 5`** | Ir a Clientes | Directorio. |
| **`Ctrl + 6`** | Ir a Historial | Buscador de documentos. |
| **`Ctrl + 8`** | Ir a Ajustes | Configuración del sistema. |

## 🧭 Navegación Lateral (Sidebar)

La nueva barra lateral utiliza un diseño **"Rail"**:
*   **Colapsada (72px):** Muestra iconos SVG de alta definición. Ideal para tener más espacio en tablas y gráficos.
*   **Expandida:** Al pasar el mouse, se despliega suavemente para mostrar las etiquetas de texto.
*   **Inteligente:** En pantallas pequeñas (laptops), se mantiene compacta para evitar solapamientos.

## ⚡ Módulo de Facturación (Invoice Emitter)

El panel más importante (`Ctrl + 2`).

1.  **Búsqueda de Cliente:** Empieza escribiendo. El sistema busca por Nombre o RUC en tiempo real.
    *   *Tip:* Presiona `Enter` para seleccionar y saltar al siguiente campo.
2.  **Agregar Ítems:**
    *   Busca productos con el buscador inteligente.
    *   Los impuestos se calculan solos.
    *   Presiona el botón `+` o `Enter` en el precio para añadir.
3.  **Validación Visual:** Si falta un dato (ej. email del cliente), el campo se pondrá rojo y no te dejará emitir.

## 📊 Dashboard Interactivo

Tu centro de mando (`Ctrl + 1`).
*   **Sin Esperas:** Los datos se cargan en paralelo.
*   **Gráficos:** Renderizados con `Echarts`, son interactivos. Pasa el mouse para ver valores exactos.
*   **KPIs:** Indicadores de ventas y estado del SRI (Online/Offline).

:::tip[Accesibilidad]
Todas las tablas (Clientes, Productos, Historial) soportan navegación por teclado. Usa `Tab` para entrar en la lista y las flechas para moverte.
:::
