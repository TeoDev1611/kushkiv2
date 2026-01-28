---
title: Guía de Cumplimiento SRI
description: Cómo utilizar el asistente fiscal para tu declaración de IVA.
---

El sistema Kushki incluye un motor de inteligencia fiscal diseñado para simplificar la declaración mensual o semestral del IVA en Ecuador.

## Asistente para Formulario 104

El Dashboard principal genera automáticamente los valores para los casilleros más importantes:

### Casilleros Soportados
- **401**: Ventas locales gravadas con tarifa diferente de cero (15%).
- **403**: Ventas locales gravadas con tarifa 0% que no dan derecho a crédito tributario.
- **411**: IVA generado en ventas locales tarifa 15%.
- **609**: Retenciones en la fuente de IVA efectuadas por los clientes.

## Factor de Proporcionalidad (Campo 702)

Si realizas ventas con tarifa 0% y 15%, el sistema calcula automáticamente tu factor de proporcionalidad basado en la fórmula:

$$FP = \frac{Ventas 15\%}{Ventas Totales}$$

## Calendario de Declaraciones

Kushki utiliza el noveno dígito de tu RUC para calcular tu fecha límite. 

| 9no Dígito | Fecha Límite |
| :--- | :--- |
| 1 | Día 10 |
| 2 | Día 12 |
| ... | ... |
| 0 | Día 28 |

:::tip[Nota Importante]
Si la fecha límite cae en fin de semana o feriado, se traslada al siguiente día hábil.
:::
