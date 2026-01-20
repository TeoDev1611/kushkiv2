# ⚡ Kushki Facturador v2.0

> **Sistema de Facturación Electrónica para Ecuador (SRI) de Alto Rendimiento.**
> Construido con tecnología híbrida nativa para máxima velocidad, seguridad y una experiencia de usuario cinematográfica.

![Status](https://img.shields.io/badge/Estado-Producción-34d399?style=for-the-badge)
![Tech](https://img.shields.io/badge/Stack-Go_xp_Svelte-blue?style=for-the-badge)
![License](https://img.shields.io/badge/Licencia-Proprietaria-orange?style=for-the-badge)

## 📖 Descripción General

**Kushki Facturador** no es solo un emisor de facturas; es una suite de gestión fiscal diseñada para la velocidad. Elimina la complejidad de la facturación electrónica del SRI mediante una arquitectura de software moderna que combina la potencia bruta de **Go (Golang)** en el backend con la interactividad fluida de **Svelte** en el frontend, todo empaquetado en un binario nativo ligero usando **Wails**.

Diseñado bajo principios de **Material Design 3**, con un tema visual "Obsidian & Mint" optimizado para reducir la fatiga visual y maximizar la eficiencia operativa.

---

## 🚀 Tecnologías y Arquitectura

### 🧠 Backend (Core de Potencia)
*   **Lenguaje:** Go (Golang) 1.21+.
*   **Bridge Nativo:** [Wails v2](https://wails.io) - Comunicación bidireccional Go ↔ JS sin servidores HTTP latentes.
*   **Firma Electrónica (Crypto):** Implementación **nativa y manual** del estándar **XAdES-BES** (XML Advanced Electronic Signatures). *No depende de librerías externas opacas ni binarios de Java.*
*   **Concurrencia:** Uso extensivo de **Goroutines** y **WaitGroups** para cálculos de métricas y firmas en paralelo.
*   **Worker Pools:** Sistema de colas (Semáforos) para la sincronización masiva con el SRI, respetando límites de tasa y evitando bloqueos de UI.

### 🎨 Frontend (Interfaz de Usuario)
*   **Framework:** Svelte - Sin Virtual DOM para un renderizado instantáneo.
*   **Estilo:** CSS Artesanal (Sin frameworks pesados como Tailwind/Bootstrap) optimizado para **renderizado por GPU**.
*   **Animaciones:** Transiciones cinemáticas (`fade`, `fly`) y **Splash Screen** de carga inicial.
*   **Componentes:**
    *   **Sidebar Colapsable:** Maximización de espacio de trabajo.
    *   **Master-Detail Layouts:** Navegación fluida en inventarios y clientes.
    *   **Gráficos:** SVG dinámicos renderizados en tiempo real.

### 💾 Persistencia de Datos
*   **Motor:** SQLite 3.
*   **Modo:** **WAL (Write-Ahead Logging)** habilitado para permitir lecturas y escrituras simultáneas sin bloqueos.
*   **ORM:** GORM con optimización de consultas e **índices compuestos** manuales para búsquedas instantáneas en historiales masivos.
*   **Seguridad:** Encriptación AES para el almacenamiento de contraseñas de firmas digitales (.p12).

---

## 💎 Ecosistema de Módulos: Potencia en cada Pixel

Kushki Facturador no es solo software, es una suite completa dividida en paneles especializados para cubrir cada aspecto de tu negocio.

### 📊 1. Centro de Mando (Dashboard "Bento Grid")
Toma decisiones basadas en datos, no en intuiciones. Nuestro Dashboard de diseño moderno te ofrece una visión de 360° de tu negocio al instante.
*   **KPIs en Tiempo Real:** Visualiza Ventas Totales, Documentos Emitidos y Pendientes con indicadores de estado semafóricos.
*   **Monitor de Estado SRI:** Verificación constante de conexión con el Servicio de Rentas Internas. Si el SRI cae, tú lo sabes primero.
*   **Tendencias de Venta:** Gráficos interactivos de alto rendimiento que muestran tu evolución financiera de los últimos 7 días.
*   **Top Products:** Identifica tus "Best Sellers" automáticamente.

### ⚡ 2. Motor de Facturación Relámpago
Olvídate de los formularios lentos y complejos del SRI. Hemos diseñado el facturador más rápido del mercado.
*   **Diseño "Master-Detail":** Formulario inteligente a la izquierda, vista previa de items a la derecha. Todo en una sola pantalla.
*   **Cálculo Tributario Automático:** El sistema maneja complejidades como IVA 15%, 5%, 0% y Exento sin que tengas que usar la calculadora.
*   **Autocompletado Inteligente:** Busca clientes y productos por nombre, RUC o código mientras escribes (Debounce Search).
*   **Workflow "One-Click":** Un solo botón para Firmar, Autorizar, Generar PDF y Enviar por Email.

### 📂 3. Auditoría y Control Total (Historial)
Tu contabilidad, siempre inmaculada y accesible.
*   **Tabla de Alta Densidad:** Visualiza decenas de transacciones sin scroll innecesario.
*   **Acciones Rápidas:** Botones inmediatos para re-imprimir RIDE (PDF), descargar XML firmado, reenviar correos o abrir la carpeta contenedora.
*   **Exportación Ejecutiva:** Genera reportes en **Excel** compatibles con cualquier sistema contable con un solo clic.
*   **Búsqueda Global:** Encuentra cualquier factura por cliente, secuencial o fecha en milisegundos.

### 📦 4. Gestión de Activos (Clientes y Productos)
Mantén tu base de datos organizada sin esfuerzo.
*   **Inventario Persistente:** Guarda productos con sus códigos de impuestos predefinidos para no repetir datos nunca más.
*   **Directorio de Clientes:** Agenda ilimitada de clientes con validación de datos.
*   **Edición "In-Place":** Modifica precios o datos de clientes sobre la marcha desde los paneles laterales.

### 🔄 5. Sincronización Resiliente
¿El SRI está caído? No pares de vender.
*   **Cola de Procesamiento:** Si el SRI falla, el sistema guarda la factura y permite reintentar el envío cuando el servicio se restablezca.
*   **Logs Técnicos Detallados:** Acceso transparente a las respuestas XML/SOAP para auditoría técnica o depuración.

---
*   **Asistente de Inicio (Wizard):** Guía paso a paso para la configuración inicial (Carga de firma, RUC, Logo).
*   **Gestión de Certificados:** Soporte para archivos `.p12` y `.pfx`.
*   **Respaldo Automático:** Generación de backups `.zip` de la base de datos y repositorio de documentos.
*   **Multi-Ambiente:** Switch instantáneo entre SRI Pruebas y SRI Producción.

---

## 🔐 Seguridad y Criptografía de Grado Bancario

Kushki Facturador ha sido diseñado bajo la premisa de **"Privacidad por Diseño"**. A diferencia de los sistemas contables en la nube, tus datos sensibles nunca abandonan tu máquina sin encriptación.

*   **🔒 Encriptación AES-256 GCM:** Las contraseñas de tus firmas electrónicas y las credenciales de correo (SMTP) se almacenan utilizando el estándar de encriptación avanzada **AES-256**. Ni siquiera alguien con acceso físico a la base de datos puede leer tus secretos.
*   **🔑 Gestión de Secretos Local:** La llave de encriptación se genera de forma única, asegurando que tus datos estén protegidos contra ataques de fuerza bruta y accesos no autorizados.
*   **🛡️ Firma XAdES-BES Nativa:** El proceso de firmado electrónico ocurre íntegramente en la memoria volátil del sistema. Tu certificado digital `.p12` nunca se expone a servidores externos ni a APIs de terceros.
*   **🚫 Zero-Cloud Dependency:** No dependemos de servidores externos para procesar tus datos. Tú eres el único dueño de tu información fiscal y contable.

---

## 📸 Experiencia Visual (UI/UX)

La interfaz ha sido pulida pixel a pixel para ofrecer una experiencia "Premium":

*   **Dark Mode Nativo:** Paleta de colores `#0B0F19` (Obsidian) con acentos `#34d399` (Mint) y `#6366f1` (Indigo).
*   **Selectores Personalizados:** Dropdowns estilizados con SVGs para consistencia en cualquier SO.
*   **Optimización Low-Level:** Eliminación de efectos costosos (`backdrop-filter`) para garantizar 60FPS incluso en hardware modesto o drivers gráficos genéricos en Linux.
*   **Feedback Inmediato:** Sistema de notificaciones "Toast" para cada acción del sistema.

---

## 🛠 Instalación y Desarrollo

### Requisitos Previos
*   Go 1.21+
*   Node.js 16+
*   NPM

### Comandos
```bash
# 1. Instalar dependencias de Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 2. Clonar repositorio
git clone https://github.com/tu-usuario/kushki-facturador.git

# 3. Ejecutar en modo desarrollo (Hot Reload)
wails dev

# 4. Compilar para Producción (Binario Optimizado)
wails build
```

---

## 📂 Estructura del Proyecto

```
kushkiv2/
├── frontend/          # SPA en Svelte
│   ├── src/
│   │   ├── components/# Wizard, Sidebar, etc.
│   │   └── App.svelte # Lógica principal y Layout
│   └── wailsjs/       # Bindings automáticos Go -> JS
├── internal/
│   ├── db/            # Modelos GORM, Migraciones, Índices
│   └── service/       # Lógica de Negocio (Invoice, Sync, Mail)
├── pkg/
│   ├── crypto/        # Firma XAdES-BES (Core Crítico)
│   ├── pdf/           # Generador RIDE
│   ├── sri/           # Cliente SOAP SRI
│   └── xml/           # Constructor UBL 2.1
├── app.go             # Controlador principal (Puente Wails)
└── main.go            # Entrypoint y gestión de ciclo de vida
```

---

> **Kushki Facturador** - Potencia, Elegancia y Cumplimiento Tributario.