# MANUAL TÉCNICO COMPLETO DEL SISTEMA DE FACTURACIÓN ELECTRÓNICA

**Versión del Documento:** 1.1.0
**Fecha:** Enero 2026
**Proyecto:** Kushki Facturador (Go/Wails/Svelte)

Este documento detalla de forma exhaustiva la arquitectura, lógica de negocio, estructuras de datos, flujos de ejecución e interfaz de usuario del software. Está dirigido a desarrolladores y auditores técnicos que requieren una comprensión total del funcionamiento interno de esta aplicación compilada de alto rendimiento.

---

## 1. INTRODUCCIÓN Y ALCANCE

El sistema es una aplicación de escritorio híbrida (Nativa/Web) diseñada para la emisión, firma, envío y almacenamiento de comprobantes electrónicos bajo normativa del SRI. A diferencia de soluciones interpretadas, este software se compila a código máquina nativo, utilizando el motor de renderizado del sistema operativo (WebView2 en Windows / WebKit en Linux).

### Capacidades Técnicas

* **Generación XML UBL 2.1**: Serialización ultra-rápida utilizando `encoding/xml` y Structs nativos de Go.
* **Firma XAdES-BES**: Implementación de criptografía de bajo nivel (RSA-SHA1/SHA256) sin dependencias externas pesadas.
* **Comunicación SOAP**: Cliente HTTP nativo (`net/http`) con gestión de Envelopes optimizada para los Web Services del SRI.
* **Renderizado PDF**: Generación programática de vectores PDF mediante **Maroto** (sin motor de navegador intermedio).
* **Persistencia**: SQLite integrado con **GORM** (Go Object Relational Mapper).
* **Interfaz**: Frontend en **Svelte** (compilado a JS puro) comunicado mediante **Wails Bridge** (Cero latencia).

---

## 2. ESTRUCTURA DEL PROYECTO

El proyecto sigue el estándar de diseño de paquetes de Go ("Standard Go Project Layout"), separando el frontend (UI) del backend (Go).

```text
/
├── main.go                     # Punto de entrada. Inicializa Wails y el ciclo de vida.
├── app.go                      # BRIDGE: Métodos expuestos al Frontend (JS puede llamar a estas funciones).
├── wails.json                  # Configuración del compilador y metadatos del binario.
├── go.mod                      # Gestión de dependencias (Module Graph).
├── frontend/                   # INTERFAZ DE USUARIO (SVELTE)
│   ├── src/
│   │   ├── App.svelte          # Componente Raíz.
│   │   ├── main.js             # Inicialización de Svelte.
│   │   ├── stores/             # Gestión de estado reactivo (Svelte Stores).
│   │   ├── pages/              # Vistas (Dashboard, Factura, Config).
│   │   └── components/         # UI Kit (Botones, Inputs, Modales).
│   └── wailsjs/                # Bindings generados automáticamente (Go -> JS).
├── internal/                   # LÓGICA PRIVADA DE APLICACIÓN
│   ├── db/                     # Capa de Acceso a Datos.
│   │   ├── connection.go       # Singleton de conexión SQLite + GORM.
│   │   └── migrations.go       # Auto-migración de esquemas.
│   └── service/                # Orquestadores de negocio.
│       └── invoice_service.go  # Lógica principal de facturación.
└── pkg/                        # LÓGICA DE NEGOCIO REUTILIZABLE (CORE)
    ├── sri/                    # INTEGRACIÓN SRI
    │   ├── soap_client.go      # Cliente HTTP para SOAP (Recepcion/Autorizacion).
    │   └── types.go            # Structs de respuesta del SRI.
    ├── xml/                    # CONSTRUCTOR XML
    │   ├── builder.go          # Marshaler de Struct a XML UBL 2.1.
    │   └── structures.go       # Definición estricta de tags XML.
    ├── crypto/                 # CRIPTOGRAFÍA
    │   └── signer.go           # Algoritmo de firma XAdES-BES (ds:Signature).
    └── pdf/                    # REPORTES
        └── ride_generator.go   # Motor gráfico (Maroto) para PDF.

```

---

## 3. ESQUEMA DE BASE DE DATOS (GORM)

La persistencia se maneja con SQLite. GORM utiliza `Structs` de Go para definir el esquema.

### Struct: `EmisorConfig`

Almacena configuración tributaria.

* **`RUC`** (PK, string): Registro Único.
* `RazonSocial` (string): Nombre legal.
* `P12Path` (string): Ruta al certificado.
* `P12Password` (string): Contraseña (Cifrada con AES-GCM antes de guardar).
* `Ambiente` (int): 1 (Pruebas), 2 (Producción).
* `Estab` (string): '001'.
* `PtoEmi` (string): '001'.
* `Obligado` (bool): Contabilidad.
* `SMTPHost`, `SMTPUser`, `SMTPPass`: Configuración de correo.

### Struct: `Factura`

Registro transaccional.

* **`ClaveAcceso`** (PK, string, size:49): ID único.
* `Secuencial` (string): 9 dígitos.
* `FechaEmision` (time.Time): Fecha real.
* `ClienteID` (string): FK.
* `Total` (float64): Importe final.
* `EstadoSRI` (string): 'PENDIENTE', 'AUTORIZADO', etc.
* `XMLFirmado` ([]byte): BLOB con el XML final para evitar regeneración.
* `MensajeError` (string): Respuesta de fallo del SRI.
* `Subtotal15` (float64): Base imponible 15%.
* `Subtotal0` (float64): Base imponible 0%.
* `IVA` (float64): Valor impuesto.

### Struct: `Product`

Inventario optimizado.

* **`SKU`** (PK, string).
* `Name` (string, index).
* `Price` (float64).
* `Stock` (int).
* `TaxCode` (int): 2 (IVA).
* `TaxPercentage` (int): 0, 15, etc.

---

## 4. LÓGICA CORE: CICLO DE EMISIÓN (GO)

La función `EmitirFactura` en `internal/service/invoice_service.go` orquesta el proceso. Al ser Go, este proceso ocurre en un hilo del sistema operativo (Goroutine) y no bloquea la UI.

### Paso 1: Mapeo de Datos y Concurrencia

1. Recibe un JSON desde el Frontend (Svelte).
2. Go deserializa el JSON a un struct `InvoiceDTO`.
3. Consulta `EmisorConfig` y obtiene el siguiente secuencial de forma atómica.

### Paso 2: Generación de Clave de Acceso (Algoritmo Módulo 11)

Implementación nativa en `pkg/util/mod11.go`.

* Cálculo matemático puro sobre arreglos de bytes (Zero-allocation) para máxima velocidad.
* Generación de `ClaveAcceso` de 49 dígitos.

### Paso 3: Serialización XML (`pkg/xml`)

No se usa construcción de DOM (lento). Se usa **Marshaling directo**.

1. Se define una estructura jerárquica de Go con "Struct Tags":
```go
type FacturaXML struct {
    XMLName id `xml:"factura"`
    InfoTributaria InfoTrib `xml:"infoTributaria"`
    // ...
}

```


2. `xml.Marshal(factura)` convierte la estructura en bytes XML instantáneamente.

### Paso 4: Firma Electrónica (`pkg/crypto`)

Este es el componente más crítico y rápido comparado con Python.

1. Carga el `.p12` en memoria.
2. Extrae la llave privada RSA y certificado X.509 usando `software.ssl.pkcs12`.
3. **Canonicalización**: Transforma el XML a formato canónico (C14N) requerido por XAdES.
4. **Digest**: Calcula el SHA-1/SHA-256 del XML.
5. **Firma**: Genera la firma RSA y construye el nodo `<ds:Signature>`.
6. Inserta la firma en el XML sin romper la estructura.

### Paso 5: Envío SOAP (`pkg/sri`)

Uso de `net/http` con timeouts estrictos (30s).

* **Recepción**: Envía el XML en Base64 dentro de un Envelope SOAP manual.
* **Autorización**: Si recepción es OK, consulta el endpoint de autorización.
* **Parsing**: Usa `xml.Unmarshal` para leer la respuesta del SRI y extraer el estado y fecha.

### Paso 6: Generación PDF (`pkg/pdf`)

Librería: **Maroto**.

* No usa HTML ni CSS. Dibuja el PDF usando coordenadas y grillas.
* Es 10x más rápido que renderizar HTML.
* Genera el QR en memoria y lo estampa en el documento.

---

## 5. INTERFAZ DE USUARIO (SVELTE + WAILS)

La UI vive en `frontend/`. Durante el desarrollo es un servidor web local, en producción es un bundle de assets embebidos en el `.exe`.

### Comunicación Frontend-Backend

* **Sin REST API**: No hay llamadas HTTP locales.
* **Wails Bindings**: Go expone métodos como `App.CrearFactura()`.
* **Promesas JS**: En Svelte se llama:
```javascript
// Svelte
let resultado = await CreateInvoice(datos);

```


Esto ejecuta código binario Go directamente.

### Dashboard (`pages/Dashboard.svelte`)

* **Stores**: Usa `svelte/store` para manejar el estado global.
* **Rendimiento**: Renderiza listas largas de facturas usando virtualización (solo renderiza lo visible en pantalla).

### Facturación (`pages/Invoice.svelte`)

* **Reactividad**: Svelte recalcula totales (Subtotal, IVA, Total) automáticamente cuando cambia una variable, sin necesidad de listeners manuales complejos.
* **Validación**: Validación de inputs en tiempo real (lado cliente) antes de enviar a Go.

---

## 6. SEGURIDAD Y RENDIMIENTO

### Compilación Binaria

* **Ofuscación**: Al compilar a código máquina, la lógica de firma y las claves de API no son legibles en texto plano como en Python/JS.
* **Memory Safety**: Go maneja la memoria automáticamente, previniendo leaks comunes en aplicaciones de escritorio de larga duración.

### Cifrado de Datos Sensibles

* Usa **AES-GCM** (Galois/Counter Mode) implementado en `crypto/aes` de la librería estándar de Go.
* Las contraseñas de P12 y Correo se cifran antes de tocar el disco (SQLite).

---

## 7. MANEJO DE ERRORES Y LOGS

### Go Error Handling

No usa excepciones (`try/catch`). Usa retorno explícito de errores:

```go
if err := s.sriClient.Enviar(xml); err != nil {
    return AppError{Code: "SRI_OFFLINE", Msg: err.Error()}
}

```

Esto asegura que **ningún** error pase desapercibido o rompa la aplicación inesperadamente.

### Wails Events

Si ocurre un error crítico en el backend (ej. "Firma caducada"), Go emite un evento:
`runtime.EventsEmit(ctx, "error-toast", "Su firma ha caducado")`.
Svelte escucha este evento y muestra una notificación nativa.

---

## 8. REQUISITOS DEL SISTEMA

Debido a la arquitectura compilada, los requisitos son mínimos.

### Dependencias de Desarrollo

* **Go 1.21+**: Compilador del lenguaje.
* **Node.js 18+**: Para compilar el frontend Svelte.
* **Wails CLI**: Herramienta de construcción (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`).
* **GCC**: Compilador C (necesario para CGO/SQLite).

### Dependencias del Usuario Final (Runtime)

* **Windows**: WebView2 Runtime (Preinstalado en Win10/11).
* **Linux**: `libwebkit2gtk-4.0` (Estándar en la mayoría de distros, incluyendo CachyOS).
* **RAM**: ~30MB - 50MB (vs 200MB+ de Electron/Python).
* **Disco**: ~15MB (Tamaño del ejecutable).

---

## 9. MÓDULO DE REPORTERÍA Y CONTABILIDAD (NUEVO)

La exportación de datos es vital para la contabilidad y análisis de negocio.

### Arquitectura de Exportación (Go + Excelize)

Se implementa un servicio `ReportService` que utiliza `github.com/xuri/excelize/v2` para generación nativa de Excel (.xlsx).

* **Reporte de Ventas (ATS):**
    * **Input:** Rango de fechas.
    * **Query:** GORM `Where("fecha_emision BETWEEN ? AND ?", start, end)`.
    * **Procesamiento:** Itera el resultado y escribe celda por celda en el stream de Excelize (High Performance).
    * **Output:** Archivo `.xlsx` formateado.

* **Reporte de Productos Más Vendidos:**
    * **Lógica:** Query de agregación SQL `SUM(cantidad)` sobre los items de factura.
    * **Visualización:** Se envían los datos JSON al frontend (Svelte) para renderizar gráficos con **Chart.js** o similar.

---

## 10. GESTIÓN DE CORREOS Y NOTIFICACIONES (SMTP ASÍNCRONO)

El envío del comprobante no bloquea la UI. Se utiliza el modelo de concurrencia de Go.

### Cola de Envíos (Go Routines)

1. **Persistencia:** Tabla `email_queue` en SQLite para persistencia ante cierres inesperados.
2. **Worker Pool:** 
    * Una **Goroutine** en background verifica la cola cada 60 segundos o ante evento de nueva factura.
    * Utiliza `gopkg.in/gomail.v2` o `net/smtp` para el envío TLS/SSL.
    * Maneja reintentos automáticos (backoff exponencial).
3. **Templates:**
    * Uso de `html/template` nativo de Go para generar correos HTML dinámicos.

---

## 11. SISTEMA DE RESPALDO Y RECUPERACIÓN (BACKUP)

### Estrategia de Respaldo

1. **Backup Automático:**
    * Al evento `OnShutdown` de Wails.
    * Uso de `archive/zip` de la librería estándar de Go.
    * Comprime `kushki.db` y carpetas de XML/PDFs firmados.
2. **Ubicación:** Configurable por el usuario (ej. carpeta sincronizada con la nube).

---

## 12. MÓDULO DE CONTINGENCIA (OFFLINE)

Garantiza la operatividad sin internet.

### Lógica de "Emisión Offline"

1. **Detección:** El cliente HTTP del SRI retorna error de conexión/timeout.
2. **Estado:** La factura se guarda en DB con estado `PENDIENTE_ENVIO` y se genera el XML firmado (que es válido legalmente).
3. **Sync Worker:**
    * Goroutine en segundo plano que hace "ping" al SRI periódicamente.
    * Al recuperar conexión, envía lotes de facturas pendientes.
    * Actualiza estados y notifica al usuario vía eventos Wails (`EventsEmit`).