---
id: introduccion
title: Introducción
sidebar_label: Introducción
---

# Introducción a Kushki Facturador

**Kushki Facturador** es una solución de escritorio moderna diseñada para simplificar la emisión, autorización y distribución de comprobantes electrónicos en Ecuador. A diferencia de las soluciones web tradicionales, Kushki ofrece una experiencia **Local-First**, garantizando que tus datos estén siempre en tu equipo.

## ¿Para quién es esto?

*   **PyMEs y Profesionales:** Que necesitan facturar rápidamente sin depender de suscripciones mensuales costosas o sistemas web lentos.
*   **Contadores:** Que requieren una herramienta fiable que genere los RIDE (PDF) y XMLs cumpliendo estrictamente la normativa del SRI.
*   **Negocios en Zonas con Internet Inestable:** Gracias a su arquitectura, puedes generar y firmar facturas sin internet, y sincronizarlas cuando recuperes la conexión.

## Capacidades Principales

### 1. Facturación Electrónica SRI
El núcleo del sistema. Permite generar facturas cumpliendo con la ficha técnica **XML 1.1.0** (versión vigente 2026).
*   Cálculo automático de impuestos (IVA 15%, 5%, 0%).
*   Manejo de formas de pago.
*   Validaciones de montos para Consumidor Final.

### 2. Gestión de Marca
Tu factura es tu imagen. El sistema permite:
*   Subir tu **Logo Corporativo**.
*   Redimensionamiento inteligente de imágenes para que se vean perfectas en el PDF.
*   Personalización de colores en el RIDE generado.

### 3. Centro de Control (Dashboard)
Olvídate de las hojas de cálculo. El dashboard te ofrece:
*   Ventas totales en tiempo real.
*   Estado de conexión con el SRI.
*   Gráfico de tendencias de ventas (filtrable por fecha).

### 4. Independencia Tecnológica
*   **Tu Correo, Tus Reglas:** Configura tu propio Gmail, Outlook o servidor corporativo para enviar las facturas. No dependes de un tercero para que tus correos lleguen.
*   **Tus Datos:** Base de datos SQLite encriptada en tu disco duro.
