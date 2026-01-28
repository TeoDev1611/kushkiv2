# 🎨 Kushki Design System (Obsidian & Mint)

Este documento define la guía de estilo visual para el ecosistema **Kushki**. El objetivo es mantener una identidad coherente entre la Landing Page y el Software de Escritorio/SaaS.

## 1. Filosofía de Diseño
* **Nombre del Tema:** Obsidian & Mint.
* **Enfoque:** "Dark Mode First" (Modo Oscuro Nativo).
* **Sensación:** Tecnológica, Financiera, Moderna, Fluida.
* **Ergonomía:** Alto contraste y transiciones suaves para una experiencia de usuario premium.

---

## 2. Paleta de Colores (Color Palette)

### Fondos y Superficies (Backgrounds)
Utilizados para la estructura principal de la interfaz.

| Nombre Variable | Código HEX | Uso |
| :--- | :--- | :--- |
| `bg-obsidian` | **#0B0F19** | Fondo principal de la ventana/pantalla. |
| `bg-surface` | **#161e31** | Tarjetas, Paneles laterales (Sidebar), Modales. |
| `bg-glass` | `rgba(11, 15, 25, 0.7)` | Overlays de carga (Loading). |
| `border-subtle` | `rgba(255, 255, 255, 0.05)` | Bordes sutiles para separar secciones. |

### Acentos de Marca (Brand Accents)
Utilizados para acciones principales, botones y estados de éxito.

| Nombre Variable | Código HEX | Uso |
| :--- | :--- | :--- |
| `primary-mint` | **#34d399** | Botones principales (CTA), Toasts de Éxito, Switch Activo. |
| `mint-dark` | **#059669** | Estado *Hover* de botones. |
| `accent-indigo` | **#6366f1** | Detalles secundarios, Toasts de Info. |
| `status-error` | **#EF3340** | Errores, Toasts de Error, Botones de borrar. |

---

## 3. Componentes UI (Especificaciones)

### Notificaciones (Toasts)
Alertas no intrusivas que aparecen en la esquina superior derecha.
*   **Dimensiones:** 340px x 70px.
*   **Fondo:** `bg-card`.
*   **Borde Izquierdo:** 4px sólido (Color según tipo: Mint, Rojo, Índigo).
*   **Sombra:** Drop Shadow difusa para efecto de flotación.
*   **Interacción:** Botón de cierre "X" discreto.

### Interruptores (Switches)
Evolución moderna del Checkbox tradicional.
*   **Estado Inactivo:** Borde gris sutil, indicador a la izquierda.
*   **Estado Activo:** Relleno `primary-mint`, indicador a la derecha.
*   **Animación:** Desplazamiento suave del indicador.

### Pantallas de Carga (Loading Overlay)
*   **Fondo:** Negro semitransparente (Obsidian con alpha).
*   **Indicador:** Spinner dibujado vectorialmente (arco giratorio) en color Mint.
*   **Texto:** Blanco, centrado, tipografía Outfit.

### Navegación (Sidebar)
*   **Item Normal:** Texto gris, fondo transparente.
*   **Item Activo:** Texto Mint, fondo con opacidad baja, **Borde izquierdo** sólido de 3px (Indicador de posición).
*   **Transiciones:** Cambio de página con efecto *Cross-Fade* (Fundido cruzado) de 300ms.

---

## 4. Tipografía (Typography)

La fuente oficial es **Outfit**.

* **Familia:** `Outfit`, sans-serif.
* **Pesos:**
    * `Regular (400)`: Cuerpo de texto, inputs.
    * `Medium (500)`: Etiquetas, Toasts.
    * `Bold (700)`: Botones, Encabezados de tabla.
    * `ExtraBold (800)`: Títulos principales, Valores monetarios.

---

*Documento actualizado para reflejar la versión v1.2.0 del sistema.*