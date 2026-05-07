# Data Model: Splash interactivo con saltos y mejoras de cana

## SplashProfile

- Purpose: Configuracion reusable del splash para una especie, preset o encounter.
- Fields:
  - `jump_count`: cantidad total de saltos requeridos; rango valido `1..5`
  - `time_limit`: duracion base de cada salto interactivo
  - `success_reward_distance`: acercamiento adicional otorgado por salto ganado desde la build del jugador
- Relationships:
  - Pertenece al contexto de encounter o al perfil del pez que puede disparar splash.
- Validation rules:
  - `jump_count` no puede ser menor que 1 ni mayor que 5
  - `time_limit` debe ser mayor que cero
  - `success_reward_distance` no puede forzar distancia por debajo del minimo efectivo del encounter

## SplashSequenceState

- Purpose: Estado runtime de un splash activo o recien resuelto.
- Fields:
  - `status`: pendiente, resuelto con exito, resuelto con escape
  - `total_jumps`: total configurado de saltos
  - `resolved_jumps`: cantidad ya resuelta
  - `current_jump`: salto actualmente activo cuando exista uno
  - `trigger_source`: motivo tactico que disparo el splash actual
  - `last_result`: ultimo resultado aplicado de la secuencia
- Relationships:
  - Vive dentro del estado del encounter o del snapshot tactico de ronda.
  - Se consume desde `internal/app` para pedir resolucion a la UI.
- Validation rules:
  - `resolved_jumps` nunca puede exceder `total_jumps`
  - una secuencia terminada no puede volver a estado pendiente sin generar un nuevo splash

## SplashJumpState

- Purpose: Unidad individual de decision dentro de la secuencia.
- Fields:
  - `index`: posicion del salto dentro de la secuencia
  - `time_limit`: tiempo disponible para ese salto
  - `result`: exito, fallo o pendiente
  - `reward_applied`: indica si el bonus por salto ganado ya fue consumido
- Relationships:
  - Forma parte de `SplashSequenceState`
- Validation rules:
  - `index` debe caer dentro del rango de la secuencia
  - `reward_applied` solo puede ser verdadero en saltos con exito

## SplashResolution

- Purpose: Resultado entregado por la UI/app al runtime despues de resolver uno o varios saltos.
- Fields:
  - `completed`: indica si la secuencia ya termino
  - `escaped`: indica si el pez escapo por fallo de splash
  - `successful_jumps`: numero de saltos ganados en la resolucion actual
  - `distance_reward_applied`: total de acercamiento aplicado por recompensas de cana
- Relationships:
  - Actualiza `SplashSequenceState`
  - Puede modificar el `encounter.State`
- Validation rules:
  - `escaped` y `completed` deben ser coherentes con la cantidad de saltos restantes
  - `distance_reward_applied` debe coincidir con los saltos exitosos que otorguen recompensa

## RodSplashBonus

- Purpose: Capacidad del loadout del jugador para recompensar saltos ganados.
- Fields:
  - `distance_per_success`: cuanto se acerca el pez por salto ganado
  - `enabled`: indica si la build aplica recompensa de splash
- Relationships:
  - Vive en el loadout efectivo del jugador y puede originarse en la cana o en aditamentos futuros.
- Validation rules:
  - `distance_per_success` no puede ser negativo

## State Transitions

- `No splash` -> `Splash pending`: una ronda provoca un intento de subida por encima de superficie.
- `Splash pending` -> `Jump success`: el jugador gana el salto actual dentro del tiempo.
- `Jump success` -> `Splash pending`: quedan saltos por resolver.
- `Jump success` -> `Splash cleared`: era el ultimo salto de la secuencia.
- `Splash pending` -> `Splash escape`: el jugador falla o vence el tiempo.
- `Splash escape` -> `Encounter finished`: el encounter conserva `EndReasonSplashEscape`.
