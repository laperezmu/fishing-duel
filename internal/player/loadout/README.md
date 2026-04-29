# internal/player/loadout

Define el setup activo del jugador como composicion de `rod` y futuros aditamentos.

## Responsabilidad

- Reunir la `rod` base con la lista de aditamentos equipados.
- Exponer los limites efectivos de apertura y track tras aplicar modificadores de aditamentos.
- Reservar el contrato de setup sin mezclar aun economia, tienda o inventario completo.

## Regla de arquitectura

- Este paquete no resuelve compra, slots ni progresion de build por si mismo.
- En el slice actual, los aditamentos ya pueden modificar limites y aportar `HabitatTags`, pero todavia no existen slots complejos ni economia de build.
