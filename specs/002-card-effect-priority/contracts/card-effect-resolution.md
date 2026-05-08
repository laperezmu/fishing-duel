# Contract: Card Effect Resolution

## Purpose

Definir el contrato interno entre contenido, engine, progression, encounter y presentation para la resolucion de efectos de cartas con prioridad determinista.

## Inputs

- Cartas de jugador y pez con multiples efectos independientes.
- Contexto del round o evento que determina que triggers son elegibles.
- Estado del encounter, thresholds, splash y recursos del jugador.

## Resolution Contract

1. El contenido declara efectos usando el catalogo cerrado de triggers y efectos soportados.
2. Cada efecto declara un trigger unico y una prioridad resoluble.
3. El runtime selecciona los efectos aplicables sin consultar UI.
4. El engine construye una secuencia ordenada por prioridad.
5. Si dos efectos comparten prioridad, el efecto del pez resuelve antes que el del jugador.
6. Encounter y progression aplican los efectos ya ordenados sobre el estado correspondiente.
7. Match snapshots exponen suficiente informacion para inspeccionar orden y resultado sin duplicar reglas.
8. Presentation y CLI solo renderizan el resultado resuelto.

## Invariants

- Una carta puede tener N efectos, pero cada efecto pertenece a un solo trigger.
- No puede haber ambiguedad entre trigger general y trigger especifico en la misma resolucion.
- Los efectos con variantes por entidad o color deben mantener esa distincion durante parseo, runtime y snapshots.
- La resolucion debe ser determinista bajo el mismo estado inicial.

## Compatibility Expectations

- `cmd/fishing-duel` y `cmd/fishing-run` siguen funcionando sin asumir conocimiento directo de prioridades.
- El contenido hardcoded del jugador y el catalogo JSON de pez deben poder migrarse al nuevo contrato dentro del mismo slice.
- Los snapshots y hints CLI pueden enriquecerse, pero no deben convertirse en la fuente de verdad de las reglas.
- La salida CLI puede mostrar la traza ordenada de efectos resueltos como evidencia de runtime, pero la secuencia sigue siendo generada exclusivamente por engine y snapshots.

## Migration Coverage Matrix

| Legacy Trigger or Effect | Status | Target Mapping | Notes |
|---|---|---|---|
| `TriggerOnDraw` | mantener | `TriggerOnDraw` / `TriggerOnCardUsed` | Se mantiene para activaciones al revelar carta; puede convivir con `on_card_used` cuando la semantica futura lo requiera. |
| `TriggerOnOwnerWin` | mantener | `TriggerOnOwnerWin` / `TriggerOnOwnerColorWin` | La variante por color se modela con `TargetMove`. |
| `TriggerOnOwnerLose` | mantener | `TriggerOnOwnerLose` / `TriggerOnOwnerColorLose` | La variante por color se modela con `TargetMove`. |
| `TriggerOnRoundDraw` | ajustar | `TriggerOnRoundDraw` / `TriggerOnColorDraw` | El empate general sigue existiendo y el empate por color se separa como trigger especifico. |
| `CaptureDistanceBonus` | reemplazar | `EffectTypeLegacyCaptureWindow` | Se conserva como capa de compatibilidad mientras migra al catalogo definitivo. |
| `SurfaceDepthBonus` | reemplazar | `EffectTypeLegacySurfaceWindow` | Se conserva como capa de compatibilidad mientras migra al catalogo definitivo. |
| `ExhaustionCaptureDistanceBonus` | ajustar | `EffectTypeLegacyExhaustionWindow` | Se mantiene como compatibilidad observable para captura por agotamiento. |
| `DistanceShift` | mantener | `EffectTypeAdvanceHorizontal` | Se normaliza con prioridad explicita. |
| `DepthShift` | mantener | `EffectTypeAdvanceVertical` | Se normaliza con prioridad explicita. |
