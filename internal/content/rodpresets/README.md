# internal/content/rodpresets

Define presets de `rod` para bootstrap y testing manual del loop de pesca.

## Responsabilidad

- Declarar presets listos para CLI con ids, nombres y detalles legibles.
- Traducir esos presets a un `loadout.State` minimo basado en `rod`.
- Mantener el contenido de equipo separado del runtime del jugador.

## Limites

- No administra inventario ni tienda.
- No decide spawn de peces ni subpools.
- No reemplaza un futuro sistema data-driven mas general de build.
