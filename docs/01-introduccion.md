---
title: Introducción a Kushki v2
description: Descubre la potencia de la facturación local-first moderna.
sidebar:
  order: 1
---

# Bienvenido a la Nueva Era de Facturación

**Kushki Facturador v2** no es solo una actualización; es una reingeniería total. Hemos abandonado las arquitecturas web lentas para ofrecerte una aplicación de escritorio nativa, construida con **tecnologías de vanguardia (Go + Svelte 4)** que garantizan velocidad instantánea y seguridad total.

## ¿Por qué Kushki v2?

A diferencia de los sistemas tradicionales que dependen de la velocidad de tu internet para cada clic, Kushki opera bajo la filosofía **Local-First**:

1.  **Velocidad Instantánea:** La interfaz carga en milisegundos. No hay "spinners" de carga innecesarios.
2.  **Privacidad Absoluta:** Tus datos viven en un archivo `SQLite` encriptado en *tu* disco duro. Nosotros no vemos tus clientes ni tus montos.
3.  **Resiliencia:** ¿Se fue el internet? Sigue facturando. El sistema firmará y guardará los XMLs localmente y los sincronizará cuando vuelvas a estar online.

## Capacidades del Sistema

### 🚀 Núcleo de Alto Rendimiento
*   **Facturación SRI XML 1.1.0:** Cumplimiento estricto de la normativa 2026.
*   **Cálculo Reactivo:** Cambia cantidades o precios y verás el IVA y Totales recalcularse instantáneamente sin recargar la pantalla.

### 💼 Gestión Comercial "Todo en Uno"
*   **Cotizaciones Integradas:** Crea proformas y conviértelas en facturas con un solo clic.
*   **Inventario Inteligente:** Búsqueda *Fuzzy* (difusa) que encuentra productos aunque escribas mal su nombre.

### 🎨 Experiencia de Usuario (UX) Obsidian
*   **Interfaz Oscura:** Diseñada para descansar la vista durante largas jornadas de trabajo.
*   **Navegación por Teclado:** Usa `Ctrl + 1` para moverte, `Ctrl + N` para facturar. Despídete del mouse.
*   **Sidebar Inteligente:** Un menú lateral que respeta tu espacio, colapsándose a 72px para darte más área de trabajo.

### 📊 Inteligencia de Negocios
*   **Dashboard en Tiempo Real:** Gráficos de ventas y KPIs que se actualizan al instante.
*   **Auditoría Total:** Logs detallados de cada conexión con el SRI y cada correo enviado.

:::tip[Ventaja Técnica]
Gracias al uso de **Wails y Svelte 4**, esta aplicación consume hasta un **80% menos de memoria RAM** que las aplicaciones basadas en Electron convencionales.
:::
