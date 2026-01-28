---
title: Arquitectura Técnica
description: Stack tecnológico y patrones de diseño utilizados en Kushki Facturador.
---

Kushki Facturador es una aplicación híbrida de escritorio construida sobre una arquitectura **Local-First**.

## Stack Tecnológico

- **Backend:** Go (Golang) 1.21+
- **Frontend:** Svelte 4 con TypeScript
- **Bridge:** Wails v2 (Comunicación bidireccional Go <-> JS)
- **Base de Datos:** SQLite 3 con GORM (Cifrado local opcional)
- **PDF:** Maroto v2 (Generación de RIDE basada en componentes)

## Patrones de Diseño

### Strategy Pattern (PDF Themes)
Utilizado para permitir múltiples plantillas de factura (Moderno, Minimalista, Corporativo) sin duplicar lógica de negocio.

### Singleton (Database)
Asegura que el motor SQLite mantenga la integridad de los datos mediante una única conexión compartida con modo **WAL (Write-Ahead Logging)** para alto rendimiento.

## Seguridad de la Firma

La firma electrónica `.p12` se almacena localmente y la contraseña se cifra utilizando **AES-256 GCM** antes de ser guardada en la base de datos.
