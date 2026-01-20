# GEMINI.md - Contexto del Proyecto "Kushki Facturador"

## Estado del Proyecto: IMPLEMENTADO (Fases 1-6)

Este proyecto es un sistema de facturación electrónica para Ecuador, construido con Go, Wails y Svelte.

## 1. Arquitectura Implementada

*   **Frontend (`frontend/`):** Interfaz en Svelte con formulario de emisión y reactividad básica.
*   **Backend Go (`main.go`, `app.go`):** Puente de comunicación y orquestación.
*   **Base de Datos (`internal/db/`):** SQLite con GORM. Migraciones automáticas activas para `EmisorConfig`, `Factura` y `Product`.
*   **Servicios (`internal/service/`):** `InvoiceService` que maneja el flujo de Clave de Acceso -> XML -> Firma -> SRI -> RIDE -> DB.
*   **Core (`pkg/`):**
    *   `crypto`: Firma XAdES-BES (Estructura base).
    *   `sri`: Cliente SOAP para Recepción y Autorización.
    *   `xml`: Constructor UBL 2.1.
    *   `pdf`: Generador de RIDE usando Maroto.
    *   `util`: Algoritmo Módulo 11.

## 2. Comandos Útiles

*   **Desarrollo:** `wails dev` (Levanta el entorno con Hot Reload).
*   **Compilación:** `wails build` (Genera el ejecutable nativo).

## 3. Próximos Pasos (Opcional)

1.  Completar el Marshalling detallado de XAdES-BES en `pkg/crypto/signer.go` (actualmente es un esqueleto funcional).
2.  Implementar la vista de Configuración en el Frontend para cargar el certificado `.p12` y datos del emisor.
3.  Añadir validaciones de esquema XML antes del envío.