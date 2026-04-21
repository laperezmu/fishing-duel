# internal/encounter

Modelo del enfrentamiento sobre el track de distancia.

## Responsabilidad

- Definir la configuracion del encuentro (`Config`).
- Validar reglas estructurales del track.
- Mantener el estado persistente del encuentro (`State`).

## Parametros importantes

- `InitialDistance`: posicion inicial del pez.
- `CaptureDistance`: umbral de captura por acercamiento.
- `EscapeDistance`: umbral de escape por alejamiento.
- `ExhaustionCaptureDistance`: umbral especial cuando el mazo se agota.
- `PlayerWinStep` y `FishWinStep`: cuanto se mueve el pez por ronda.

## Cuando extenderlo

- Para soportar peces con tracks distintos.
- Para anadir configuraciones de encuentro por especie, zona o modo de juego.

## Limites

- No modifica el estado por si mismo; eso lo hace `internal/progression/`.
- No decide el fin; eso lo hace `internal/endings/`.
