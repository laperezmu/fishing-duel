# Plan de feature: reorganizacion-docs-y-backlog

## Objetivo
Reorganizar la documentacion del proyecto para separar con claridad los planes de features ya trabajados de un backlog de investigacion y desarrollo futuro, dejando una estructura mas escalable para seguir evolucionando la vision roguelike del juego.

## Criterios de aceptacion
- La carpeta `docs/` queda reorganizada con al menos `docs/features/` y `docs/backlog/`.
- Los planes de features existentes se almacenan bajo `docs/features/`.
- Existe una estructura inicial de backlog en `docs/backlog/` con la lista de pendientes propuesta y una plantilla para agregar nuevas tareas.
- `PROJECT_CONTEXT.md` y las guias de `docs/` reflejan la nueva organizacion documental.

## Scope
### Incluye
- Reorganizar la jerarquia de `docs/`.
- Mover los planes existentes a `docs/features/`.
- Crear la documentacion base del backlog y su template.
- Cargar un backlog inicial alineado con la vision roguelike del proyecto.

### No incluye
- Implementar nuevas features de gameplay.
- Convertir aun cada backlog item en una feature plan individual.
- Cambiar el flujo de ramas o PR mas alla de actualizar rutas documentales.

## Propuesta de implementacion
- Crear `docs/features/` como ubicacion oficial de planes de feature y mover ahi los documentos existentes.
- Crear `docs/backlog/` con una guia, una plantilla y un backlog inicial del producto.
- Actualizar `PROJECT_CONTEXT.md` y `docs/README.md` para reflejar la nueva convencion.

## Validacion
- Verificar que la estructura de `docs/` quede consistente y navegable.
- Verificar que las referencias a planes de feature apunten a `docs/features/`.

## Riesgos o decisiones abiertas
- Hay que decidir si el backlog vive como un solo documento maestro o como entradas separadas; esta primera version debe dejar ambas opciones abiertas.
