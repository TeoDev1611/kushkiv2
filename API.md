# Guía de Integración: Validación de Licencias

Este documento es interno y explica cómo integrar tu aplicación cliente con la API de validación de licencias de Kushki.

**⚠️ IMPORTANTE:** Este archivo está ignorado por git (`.gitignore`) para uso privado.

---

## Endpoint de Activación / Validación

Este endpoint sirve para dos propósitos:
1.  **Activar** una licencia nueva en una máquina específica.
2.  **Validar** una licencia existente cada vez que la aplicación arranca.

### Detalles HTTP

- **URL:** `https://<TU-DOMINIO-O-LOCALHOST>/api/v1/license/activate`
- **Método:** `POST`
- **Content-Type:** `application/json`

### Cuerpo de la Petición (JSON)

Debes enviar la clave de licencia ingresada por el usuario y un identificador único de su máquina (Machine ID).

```json
{
  "license_key": "KSH-ABCD-1234-EFGH",
  "machine_id": "MACHINE-UNIQUE-ID-999"
}
```

> **Nota sobre `machine_id`:** Es responsabilidad del cliente generar este ID. Puede ser un hash de la MAC address, número de serie del disco, o un UUID generado y guardado localmente en la primera instalación.

---

## Respuestas

### ✅ 1. Éxito (Licencia Válida)

Si la licencia es válida y corresponde a la máquina (o es la primera vez que se usa y se enlaza exitosamente):

**Código HTTP:** `200 OK`

```json
{
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "message": "Activado"
}
```

*   **`token`**: Es un JSON Web Token (JWT) firmado digitalmente por la API usando RSA-256.
*   **Uso del Token:** Tu aplicación cliente debe verificar la firma de este token usando la **Clave Pública** correspondiente para asegurarse de que la respuesta es legítima y no fue falsificada.

### ❌ 2. Error (Licencia Inválida)

Si la licencia no existe, está desactivada, o ya fue usada en otra máquina.

**Código HTTP:** `403 Forbidden`

```json
{
  "error": "Licencia inválida o desactivada"
}
```

o si pertenece a otra máquina:

```json
{
  "error": "Licencia en uso por otra máquina"
}
```

---

## Ejemplo de Implementación (Cliente JavaScript/TypeScript)

```typescript
async function validateLicense(key: string, machineId: string) {
  const response = await fetch('http://localhost:8000/api/v1/license/activate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      license_key: key,
      machine_id: machineId
    })
  });

  const data = await response.json();

  if (response.ok) {
    console.log("Licencia Válida. Token:", data.token);
    // TODO: Verificar firma del JWT aquí
    return true;
  } else {
    console.error("Error de validación:", data.error);
    return false;
  }
}
```
