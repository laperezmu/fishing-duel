# Backlog

Esta carpeta contiene backlog, priorizacion y estado visible del trabajo activo y completado del producto.

## Objetivo

- Capturar ideas futuras antes de convertirlas en una feature implementable.
- Separar backlog de la ejecucion activa de features.
- Dar visibilidad a la vision roguelike, sus dependencias y el estado real de cada item.

## Estado visible

- `pending`: item identificado pero aun sin plan activo de implementacion.
- `planned`: item ya convertido en plan dentro de `docs/features/`, pero todavia no integrado en `main`.
- `done`: item ya integrado en `main` y movido al backlog completado.
- `cancelled`: item descartado, absorbido por otro o ya no vigente.

## Convencion

- Usa `TEMPLATE.md` para nuevas entradas.
- Cada item del backlog debe tener un identificador propio de su linea, por ejemplo `BL-###` para el backlog roguelike principal o `GUI-###` para el backlog de interfaz grafica.
- Cada item del backlog debe declarar un `Estado` visible.
- Un backlog item puede vivir dentro de un documento maestro o pasar a un archivo propio si crece.
- Cuando un item este suficientemente claro para implementarse, se convierte en un plan dentro de `docs/features/`.
- `docs/discoveries/` queda como contexto legacy y solo debe referenciarse cuando aporte contexto historico util.

## Documentos iniciales

- `docs/backlog/001-backlog-activo-roguelike.md`: backlog activo del producto con estado visible.
- `docs/backlog/002-backlog-completado-roguelike.md`: backlog completado ya integrado en `main`.
- `docs/backlog/003-backlog-interfaz-grafica-y-prototipo-visual.md`: backlog separado para cliente grafico, investigacion de librerias y primer vertical slice visual.
- `docs/backlog/TEMPLATE.md`: plantilla para nuevas tareas.
