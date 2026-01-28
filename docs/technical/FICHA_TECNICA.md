# 📋 Guía de Validación y Auditoría de Comprobantes Electrónicos (SRI - Offline)

Esta guía sirve como lista de verificación (checklist) para asegurar que los comprobantes electrónicos (Facturas, Retenciones, Notas de Crédito/Débito, Guías de Remisión) cumplan con los estándares técnicos y legales del SRI antes de su envío.

## 1. Requisitos Previos y Generales (Pre-Validación)

Antes de analizar la estructura del XML, verifique que el emisor cumpla con lo siguiente:

* [ ] **Certificado Digital:** El emisor posee una firma electrónica válida y vigente (archivo `.p12` o token).


* [ ] **Estado Tributario:** El RUC debe estar **ACTIVO** y no tener obligaciones pendientes.


* [ ] **Ambiente Correcto:**
* **Pruebas (Código 1):** Para desarrollo. No tienen validez tributaria.


* **Producción (Código 2):** Para emisión real. Tienen validez tributaria y legal.




* [ ] **Software de Firma:** La firma debe generarse bajo el estándar **XAdES-BES**.


* 
**Algoritmo:** RSA-SHA1.


* 
**Codificación:** UTF-8.





---

## 2. Validación de la Clave de Acceso (El "ADN" del Comprobante)

Cada comprobante debe tener una clave de 49 dígitos única. Si esta falla, el SRI rechazará todo el lote.

**Estructura Obligatoria (49 dígitos):**

1. [ ] **Fecha de Emisión:** 8 dígitos (ddmmaaaa).


2. [ ] **Tipo de Comprobante:** 2 dígitos (ver Tabla A).


3. [ ] **RUC:** 13 dígitos (del emisor).


4. [ ] **Tipo de Ambiente:** 1 dígito (1=Pruebas, 2=Producción).


5. [ ] **Serie:** 6 dígitos (Establecimiento + Punto de emisión).


6. [ ] **Secuencial:** 9 dígitos (Número único del comprobante).


7. [ ] **Código Numérico:** 8 dígitos (Generado internamente para seguridad).


8. [ ] **Tipo de Emisión:** 1 dígito (Siempre "1" para emisión normal offline).


9. [ ] **Dígito Verificador:** 1 dígito (Calculado con algoritmo **Módulo 11**).



> **🚨 Punto Crítico:** Verificar que el *Dígito Verificador* corresponda matemáticamente a los 48 dígitos anteriores.

---

## 3. Checklist de Estructura XML (Por Secciones)

### A. Información Tributaria (`<infoTributaria>`) - *Común para todos*

* [ ] **RUC y Razón Social:** Deben coincidir exactamente con el RUC del emisor.


* [ ] **Dirección Matriz:** Campo obligatorio (Máx. 300 caracteres).


* [ ] **Establecimiento y Punto de Emisión:** Deben estar activos en el SRI (3 dígitos c/u).



### B. Información de la Factura (`<infoFactura>`)

* [ ] **Fecha de Emisión:** Formato `dd/mm/aaaa`.


* [ ] **Identificación del Comprador:**
* *RUC:* 13 dígitos.
* *Cédula:* 10 dígitos.
* *Pasaporte:* 13 caracteres (aprox).
* 
*Consumidor Final:* Usar `9999999999999`.




* [ ] **Totales Obligatorios:**
* `totalSinImpuestos`: Suma de bases imponibles.
* `totalDescuento`: Suma de descuentos.
* 
`importeTotal`: Suma final incluyendo impuestos y propina.




* [ ] **Pagos:** Se debe especificar la forma de pago (Tabla 24), plazo y unidad de tiempo.



### C. Detalle de Productos (`<detalles>`)

* [ ] **Validación de Precios:**
* 
`precioUnitario` * `cantidad` = `precioTotalSinImpuesto` (Revisar redondeos a 2 decimales, excepto unitario que puede tener hasta 6).




* [ ] **Impuestos por Ítem:** Cada ítem debe tener su bloque `<impuestos>` detallando Código, Tarifa y Valor.



---

## 4. Validaciones Específicas por Tipo de Documento

### 🔹 Guía de Remisión (`<guiaRemision>`)

* [ ] **Fechas de Transporte:** Obligatorio `fechaIniTransporte` y `fechaFinTransporte`.


* [ ] **Transportista:** RUC/Cédula y Razón Social del transportista obligatorios.


* [ ] **Placa:** Número de placa del vehículo obligatorio.


* [ ] **Destinatarios:** Dirección de destino y motivo del traslado obligatorios.



### 🔹 Comprobante de Retención (`<comprobanteRetencion>`)

* [ ] **Periodo Fiscal:** Formato `mm/aaaa` obligatorio.


* [ ] **Impuestos a Retener:** Debe usar los códigos correctos (Renta=1, IVA=2, ISD=6).


* [ ] **Documento Sustento:** Número y tipo de documento sobre el cual se aplica la retención.



### 🔹 Nota de Crédito (`<notaCredito>`)

* [ ] **Motivo:** Campo obligatorio explicando la razón de la modificación.


* [ ] **Documento Modificado:** Debe referenciar el tipo (`codDocModificado`) y número (`numDocModificado`) de la factura original.



---

## 5. Tablas Maestras de Referencia (Catálogos)

El agente debe cruzar la información del XML con estas tablas para asegurar que los códigos existen:

| Tabla | Descripción | Valores Clave | Fuente |
| --- | --- | --- | --- |
| **Tabla 3** | Tipo Comprobante | 01=Factura, 04=Nota Crédito, 05=Nota Débito, 06=Guía Remisión, 07=Retención

 |  |
| **Tabla 6** | Identificación | 04=RUC, 05=Cédula, 06=Pasaporte, 07=Consumidor Final

 |  |
| **Tabla 16** | Códigos Impuesto | 2=IVA, 3=ICE, 5=IRBPNR

 |  |
| **Tabla 17** | Tarifa IVA | 0=0%, 2=12%, 3=14%, 6=No Objeto, 7=Exento

 |  |
| **Tabla 19** | Impuesto Retención | 1=Renta, 2=IVA, 6=ISD

 |  |
| **Tabla 24** | Formas de Pago | 01=Sin sistema financiero, 19=Tarjeta Crédito, 20=Otros con sist. financiero

 |  |

---

## 6. Errores Frecuentes a Detectar (Troubleshooting)

Si el SRI devuelve error, revise primero estos códigos comunes en la respuesta:

* 
**Error 35 (Documento Inválido):** Estructura XML mal formada o no cumple el esquema XSD.


* **Error 43 (Clave Registrada):** Ya se envió un comprobante con esa misma clave de acceso. Regenerar secuencial o verificar si ya se autorizó.


* 
**Error 46 (RUC no existe):** El RUC del emisor o comprador está mal escrito.


* **Error 70 (Clave en procesamiento):** El comprobante ya fue enviado y está siendo procesado. **NO reenviar** inmediatamente, esperar respuesta.



---

## 7. Proceso de Transmisión (Flujo del Agente)

1. **Generar XML:** Crear archivo según estándares anteriores.
2. **Firmar:** Aplicar firma XAdES-BES.
3. **Validar (Offline):** Correr validación de esquema XSD localmente.
4. **Enviar a Recepción:** Consumir Web Service de Recepción.
* *Si respuesta es "RECIBIDA":* Pasar al paso 5.
* *Si respuesta es "DEVUELTA":* Corregir errores reportados.


5. **Solicitar Autorización:** Consumir Web Service de Autorización.
* Verificar estado: **AUTORIZADO**.

---
**Nota Final:** Recuerde que los campos alfanuméricos no deben contener caracteres de escape mal formados (como `&` sin codificar como `&amp;`) ni saltos de línea innecesarios dentro de los tags.
