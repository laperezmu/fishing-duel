# internal/content/attachmentpresets

Define presets de aditamentos para componer el loadout del jugador sobre una `rod` base.

## Responsabilidad

- Ofrecer presets simples y jugables para testing manual del loop.
- Traducir esos presets a listas de `loadout.Attachment`.
- Mantener separados contenido equipable y runtime del loadout.

## Limites

- No administra inventario, rarezas ni economia.
- No resuelve aun slots complejos ni restricciones por tienda.
- Sus `HabitatTags` son preparacion para `BL-022`, no reglas de spawn activas todavia.
