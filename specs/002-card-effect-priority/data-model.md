# Data Model: Rediseno de triggers y efectos de cartas con prioridad de resolucion

## Trigger Catalog Entry

- Purpose: Representa un trigger soportado por el sistema de cartas.
- Key fields:
  - `id`: identificador canonico del trigger.
  - `owner_scope`: indica si aplica a pez, jugador o ambos.
  - `timing_window`: momento funcional en el que puede activarse.
  - `color_scope`: indica si el trigger depende de un color concreto o de un resultado general.
  - `semantic_notes`: documenta si mantiene semantica actual, cambia o es nuevo.
- Validation rules:
  - Debe pertenecer al catalogo cerrado definido por la spec.
  - Debe poder evaluarse sin depender de presentation o CLI.
- Relationships:
  - Un `Card Effect Binding` referencia exactamente un `Trigger Catalog Entry`.
  - Puede tener relacion de especializacion con otro trigger mas general.

## Effect Catalog Entry

- Purpose: Representa un tipo de efecto soportado y su comportamiento observable.
- Key fields:
  - `id`: identificador canonico del efecto.
  - `priority`: prioridad numerica o equivalente ordenable.
  - `target_scope`: pez, jugador, mazo actual, descarte u otra superficie de aplicacion.
  - `color_scope`: general o color concreto.
  - `effect_family`: movimiento, descarte, reshuffle, fatiga, splash u ocultamiento.
  - `deprecated_replacement_of`: referencia opcional a efectos legacy sustituidos.
- Validation rules:
  - Debe pertenecer al catalogo objetivo cerrado.
  - Debe tener prioridad definida.
  - Debe declarar restricciones cuando solo aplique a una entidad.
- Relationships:
  - Un `Card Effect Binding` referencia exactamente un `Effect Catalog Entry`.

## Entity-Specific Effect Variant

- Purpose: Diferencia variantes de un mismo efecto visible que cambian segun entidad o color.
- Key fields:
  - `base_effect_id`: efecto base compartido.
  - `entity_scope`: pez o jugador.
  - `color_requirement`: color requerido si aplica.
  - `application_rule`: describe sobre que runtime impacta.
- Validation rules:
  - Debe evitar ambiguedad entre variantes visualmente parecidas.
  - No puede eliminar la trazabilidad del efecto base.

## Card Effect Binding

- Purpose: Une un efecto concreto dentro de una carta con el trigger que lo activa.
- Key fields:
  - `trigger_id`: trigger unico del binding.
  - `effect_id`: efecto asociado.
  - `priority`: prioridad efectiva del binding si no se deriva del catalogo base.
  - `parameters`: datos del efecto, como magnitud, color o modo de reshuffle.
  - `source_owner`: pez o jugador.
- Validation rules:
  - Cada binding debe tener exactamente un trigger.
  - Una carta puede contener multiples bindings.
  - Cada binding debe ser resoluble en forma determinista.

## Resolution Priority Rule

- Purpose: Define el orden estable de todos los efectos aplicables en una misma resolucion.
- Key fields:
  - `phase`: draw, outcome u otro momento funcional equivalente.
  - `ordered_entries`: secuencia de bindings aplicables.
  - `tie_breaker`: regla de desempate.
- Validation rules:
  - Debe producir el mismo orden para el mismo estado inicial.
  - En empate de prioridad, el pez resuelve antes que el jugador.
  - Debe evitar doble activacion de triggers generales y especificos.

## Migration Coverage Map

- Purpose: Inventaria la relacion entre contrato actual y contrato objetivo.
- Key fields:
  - `legacy_trigger_or_effect`: identificador legacy.
  - `status`: mantener, ajustar, reemplazar o retirar.
  - `target_mapping`: trigger o efecto objetivo.
  - `notes`: impacto observable o diferencias de semantica.
- Validation rules:
  - Debe cubrir el 100% de triggers y efectos activos en el repo.
  - Cada item deprecado debe tener una ruta de tratamiento explicita.

## State Transitions

- Eligible -> Ordered: los efectos aplicables del round o evento se filtran y pasan a una secuencia determinista.
- Ordered -> Applied: cada efecto se aplica segun prioridad y desempate.
- Applied -> Snapshot: el resultado resuelto se expone para progreso, finales, tests y presentation.
- Legacy -> Migrated: cada trigger o efecto viejo se mapea al contrato nuevo o se marca como retirado.
