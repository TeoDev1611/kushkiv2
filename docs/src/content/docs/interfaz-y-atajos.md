---
title: Interfaz y Productividad
description: Guía de navegación, componentes visuales y atajos de teclado para Kushki Facturador.
---

Kushki Facturador utiliza un sistema de diseño denominado **"Obsidian & Mint"**, optimizado para entornos de alta productividad y fatiga visual reducida.

## Atajos de Teclado (Power User)

Para acelerar tu flujo de trabajo, hemos implementado atajos globales que funcionan en cualquier parte de la aplicación:

| Comando | Acción |
| :--- | :--- |
| <kbd>Ctrl</kbd> + <kbd>N</kbd> | **Nueva Factura:** Limpia el formulario y te lleva a la emisión. |
| <kbd>Ctrl</kbd> + <kbd>S</kbd> | **Guardar:** Guarda cambios en el panel activo (Productos, Clientes o Configuración). |
| <kbd>Ctrl</kbd> + <kbd>1-8</kbd> | **Navegación:** Cambia entre pestañas (1: Dashboard, 2: Factura, 3: Cotización...). |
| <kbd>Esc</kbd> | **Salir:** Quita el foco de campos de texto o cierra menús desplegables. |

## Elementos de la Interfaz

### Barra Lateral (Sidebar)
Ubicada a la izquierda, organiza el sistema en 4 grandes áreas:
1.  **Operación:** Dashboard y Emisión de Facturas/Cotizaciones.
2.  **Catálogos:** Gestión de Productos e Inventario, y Directorio de Clientes.
3.  **Auditoría:** Historial de Documentos y Logs de Sincronización.
4.  **Sistema:** Panel de Configuración avanzada.

### Notificaciones (Toasts)
En la esquina inferior derecha aparecerán avisos flotantes:
- <span style="color: #34d399">●</span> **Verde:** Operación exitosa (Factura autorizada, datos guardados).
- <span style="color: #6366f1">●</span> **Azul:** Información del sistema o procesos en segundo plano.
- <span style="color: #EF3340">●</span> **Rojo:** Errores de validación o fallos de conexión con el SRI.

### Búsqueda Inteligente (Fuzzy Search)
En los módulos de Clientes y Productos, el buscador es "tolerante a errores". Puedes escribir parte del nombre, el RUC o el SKU, y el motor encontrará coincidencias aproximadas instantáneamente.

## Estados de Documentos

A lo largo de la interfaz verás etiquetas de colores para el estado SRI:

- **AUTORIZADO:** El documento es legalmente válido y tiene firma electrónica.
- **PENDIENTE:** El documento está en cola para ser enviado al SRI.
- **DEVUELTA / RECHAZADA:** Hubo un error en la validación (RUC inválido, firma caducada, etc.). Revisa el historial para ver el motivo exacto.

:::tip[Pro-Tip]
Puedes usar <kbd>Ctrl</kbd> + <kbd>S</kbd> mientras editas un producto para guardarlo sin necesidad de mover el ratón al botón de "Guardar".
:::
