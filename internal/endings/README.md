# internal/endings

Condiciones de finalizacion del encuentro.

## Responsabilidad

- Inspeccionar el `game.State` actualizado.
- Marcar si la partida termino.
- Registrar el estado terminal del encuentro.

## Implementacion actual

- `EncounterCondition`: captura por proximidad, escape por exceso de distancia y resolucion especial cuando el mazo se agota.

## Como extenderlo

- Crea nuevas condiciones con reglas especiales.
- Ejemplos: limite de rondas, captura instantanea por evento, escape por tormenta.

## Regla de arquitectura

- Debe trabajar sobre el estado ya progresado.
- No debe encargarse de renderizar mensajes finales.
