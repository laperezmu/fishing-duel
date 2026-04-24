# Backlog

Esta carpeta contiene roadmap, priorizacion y estado visible de los pendientes del producto.

## Objetivo

- Capturar ideas futuras antes de convertirlas en una feature implementable.
- Separar backlog de la ejecucion activa de features.
- Dar visibilidad a la vision roguelike, sus dependencias y el estado real de cada item.

## Estado visible

- `pending`: item identificado pero aun sin plan activo de implementacion.
- `planned`: item ya convertido en plan dentro de `docs/features/`, pero todavia no integrado en `main`.
- `done`: item ya integrado en `main`.
- `cancelled`: item descartado, absorbido por otro o ya no vigente.

## Convencion

- Usa `TEMPLATE.md` para nuevas entradas.
- Cada item del backlog debe tener un identificador `BL-###`.
- Cada item del backlog debe declarar un `Estado` visible.
- Un backlog item puede vivir dentro de un documento maestro o pasar a un archivo propio si crece.
- Cuando un item este suficientemente claro para implementarse, se convierte en un plan dentro de `docs/features/`.
- `docs/discoveries/` queda como contexto legacy y solo debe referenciarse cuando aporte contexto historico util.

## Documentos iniciales

- `docs/backlog/001-roadmap-roguelike.md`: backlog inicial del producto con estado visible.
- `docs/backlog/TEMPLATE.md`: plantilla para nuevas tareas.
