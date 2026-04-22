# internal/encounter

Modelo del enfrentamiento sobre los ejes espacial y vertical del pez.

## Responsabilidad

- Definir la configuracion del encuentro (`Config`).
- Validar reglas estructurales del encounter.
- Mantener el estado persistente del encuentro (`State`).
- Aplicar deltas espaciales y resolver eventos de superficie como chapoteo.

## Parametros importantes

- `InitialDistance`: posicion inicial del pez.
- `InitialDepth`: profundidad inicial del pez.
- `SurfaceDepth`: limite superior del eje vertical.
- `CaptureDistance`: umbral de captura por acercamiento.
- `ExhaustionCaptureDistance`: umbral especial cuando el mazo se agota.
- `SplashEscapeChance`: probabilidad base de escape al chapotear en superficie.
- `PlayerWinStep` y `FishWinStep`: cuanto se mueve el pez por ronda.

## Cuando extenderlo

- Para soportar peces con tracks distintos.
- Para anadir configuraciones de encuentro por especie, zona o modo de juego.
- Para introducir eventos o capas de espacio adicionales.

## Limites

- Expone helpers para aplicar deltas del encounter, pero no decide la progresion base del round; eso lo hace `internal/progression/`.
- No decide el fin; eso lo hace `internal/endings/`.
