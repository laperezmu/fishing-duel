# Data Model: Formalizar fishing-duel como sandbox de encounters

## Sandbox Configuration

- Purpose: Representa la configuracion completa que arranca un encounter de sandbox.
- Key fields:
  - `mode`: guiado, manual o escenario reutilizable.
  - `player_deck_preset_id`: preset base del jugador.
  - `rod_preset_id`: preset de cana.
  - `attachment_preset_id`: preset de aditamentos.
  - `fish_preset_id`: preset manual de pez cuando aplica.
  - `water_context_id`: contexto de agua seleccionado.
  - `seed`: valor reproducible de aleatoriedad.
  - `opening_override`: override de apertura o cast.
  - `state_overrides`: overrides controlados del estado inicial.
  - `card_selection`: selecciones concretas de cartas del jugador y del pez.
  - `scenario_id`: referencia a escenario guardado si se usa.
- Validation rules:
  - Debe poder resolverse a un encounter valido.
  - Debe distinguir claramente valores heredados del preset frente a valores sobrescritos.
  - No debe requerir conocimiento de presentation para validarse.

## Sandbox Card Selection

- Purpose: Describe una seleccion manual de cartas concretas usada para construir una prueba puntual.
- Key fields:
  - `owner_scope`: jugador o pez.
  - `source_preset_id`: preset o perfil base del que se parte.
  - `selected_cards`: lista de referencias estables de cartas concretas.
  - `selection_mode`: mantener preset, reemplazar parcialmente o definir escenario completo.
- Card identity mechanism: Cada carta se identifica por una referencia estable compuesta por `owner_scope`, `source_preset_id` y un `card_ref` catalogado dentro del preset base. `card_ref` debe ser unico dentro del preset y permanecer estable para soportar replay de escenarios. Si el contenido actual aun no expone un ID explicito por carta, el contrato del sandbox debe introducir uno derivado del catalogo cargado en memoria y persistirlo en escenarios.
- Tie-breaking rule: Si dos cartas comparten nombre, trigger o forma de efecto, el sandbox no debe resolver por coincidencia difusa. Solo acepta una `card_ref` unica; si la referencia no existe o es ambigua, devuelve un error explicito.
- Source tracking: Cada carta debe indicar su origen: preset_base, manual_replacement, o scenario_defined.
- Validation rules:
  - Debe conservar consistencia con el dominio del owner correspondiente.
  - Debe indicar si cada carta viene del preset base o de un reemplazo de escenario.

## Sandbox Override

- Purpose: Ajusta parte del estado derivado para una prueba de QA o debugging.
- Key fields:
  - `opening_override`: cast band, distancia inicial, profundidad inicial.
  - `round_threshold_overrides`: reservado para futuras iteraciones, fuera del MVP actual.
  - `deck_state_overrides`: recycle count, exhaustion state, visibilidad de descarte.
  - `encounter_state_overrides`: estado previo al splash u otras condiciones iniciales soportadas.
- MVP scope (FR-007-MVP): En la primera iteracion, solo se soportan overrides de distancia inicial, profundidad inicial, `CaptureDistance` base del `encounter.Config` y recycle count. `ExhaustionCaptureDistance` se hereda del config resuelto y no es configurable en el MVP. Los overrides de thresholds de round, exhaustion state, visibilidad de descarte y estado previo a splash se diferiran a iteraciones futuras.
- Validation rules:
  - Debe respetar limites del sandbox para evitar estados imposibles o incoherentes.
  - Debe generar mensajes claros cuando contradiga una combinacion base invalida.

## Sandbox Scenario

- Purpose: Configuracion reutilizable para reproducir una prueba de sandbox.
- Key fields:
  - `id`: identificador estable del escenario.
  - `name`: nombre legible.
  - `description`: proposito del escenario.
  - `base_configuration`: referencia a configuracion base.
  - `card_selection`: cartas concretas fijadas si aplica.
  - `seed`: semilla fija opcional.
  - `expected_observables`: resultados o trazas esperadas para comparacion.
- Validation rules:
  - Debe poder ejecutarse desde una ruta de replay de escenarios sin reconfigurar manualmente cada parametro del setup.
  - Debe ser compartible y reproducible entre usuarios.

## Resolution Trace

- Purpose: Evidencia observable de una ronda resuelta dentro del sandbox.
- Key fields:
  - `eligible_triggers`: triggers evaluados.
  - `activated_triggers`: triggers que dispararon.
  - `resolved_effects`: secuencia ordenada de efectos con prioridad.
  - `tie_break_notes`: evidencia de desempates aplicados.
  - `before_state`: snapshot estructurado previo a la resolucion.
  - `after_state`: snapshot estructurado posterior a la resolucion.
- Validation rules:
  - Debe derivarse del runtime real.
  - Debe ser suficiente para explicar el resultado sin abrir el codigo.

## State Transitions

- Configuration Draft -> Resolved Setup: la configuracion del sandbox se valida y se traduce a estado inicial del encounter.
- Resolved Setup -> Running Encounter: el sandbox arranca el engine con presets, cartas y overrides resueltos.
- Running Encounter -> Round Trace: cada ronda produce evidencia estructurada de la resolucion.
- Scenario Defined -> Scenario Replayed: una configuracion reusable se ejecuta de forma repetible con la misma semilla y observables esperados.
