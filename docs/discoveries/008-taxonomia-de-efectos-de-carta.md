# Discovery: taxonomia de efectos de carta

## Objetivo

Definir una taxonomia de efectos de carta para combate que sea agnostica al lado que juega la carta y que pueda reutilizarse tanto para `FishCard` como para futuras `playerCards`.

Este documento busca servir como puente entre la implementacion actual y la siguiente feature tecnica de arquitectura. No propone aun un motor completo en codigo, pero si una estructura de conceptos y reglas para evitar seguir agregando efectos de forma ad hoc.

## Scope de este discovery

### Incluye
- Efectos que nacen de cartas del pez o del jugador.
- Triggers, fases de resolucion y familias de impacto producidas por cartas.
- Restricciones de composicion, stacking y ownership de efectos.
- Mapeo entre la implementacion actual y la taxonomia objetivo.

### No incluye
- Efectos que no provengan de cartas.
- Redisenar reglas base de victoria o escape que hoy dependen del `encounter`, del `rig` o del agotamiento del mazo.
- Decidir aun el formato final de persistencia o contenido data-driven.

## Inventario del estado actual

### Lo que ya existe
- `internal/cards/cards.go` define `FishCard` y `PlayerCard` como tipos distintos que comparten `CardEffect`.
- `CardEffect` hoy puede desplazar `DistanceShift` y `DepthShift`.
- Los triggers existentes en datos ya son relativos al owner de la carta: `TriggerOnOwnerWin`, `TriggerOnOwnerLose` y `TriggerOnRoundDraw`.
- `internal/game/engine.go` resuelve la ronda base y luego entrega `match.ResolvedRound` a progresion.
- `internal/progression/track.go` aplica primero la progresion base del outcome y despues suma los efectos de carta ya filtrados por el motor.
- `internal/encounter/transition.go` transforma el resultado espacial en eventos derivados como `splash`.
- `internal/endings/encounter.go` resuelve captura o escape despues de la progresion; hoy el cierre no lo decide directamente ninguna carta.
- `internal/endings/encounter.go` ya expone un `RoundState` temporal para thresholds que deban durar solo un round.
- `internal/rules/` ya tiene hooks optativos para futuros efectos de carta que alteren outcome base sin sacar esa responsabilidad de `rules`.

### Lo que aun no existe
- No existen `playerCards` integradas en el loop jugable actual.
- No existe un sistema de prioridades o resolucion de conflictos entre efectos.
- No existe aun una fase real de `on_draw` conectada al motor.

## Problemas detectados en el modelo previo

### Triggers acoplados al lado equivocado

`TriggerOnPlayerWin` y `TriggerOnFishWin` describian el resultado desde el punto de vista del encuentro completo, no desde el punto de vista de la carta que emite el efecto. Ese problema ya se corrigio migrando a triggers relativos al owner.

### Efectos mezclados con resultados absolutos

La forma actual del modificador esta pensada para mutar el `encounter`, pero no diferencia con claridad:
- que ocurre por la regla base del combate
- que ocurre por una carta
- que ocurre como evento derivado del `encounter`

La separacion existe en la practica, pero todavia no esta explicitada como taxonomia formal.

### Modelo aun demasiado centrado en el pez

La estructura concreta nacio como `FishCard`, aunque `internal/cards/README.md` ya anticipaba reutilizacion futura. Esta feature ya bajo esa intencion a contratos compartidos, pero todavia falta integrarla en el loop con `playerCards` reales.

## Principios de la taxonomia objetivo

### 1. El sistema es agnostico al owner de la carta
- Un mismo tipo de efecto debe poder existir en una carta del pez o del jugador.
- La resolucion de un efecto no debe depender de si la carta pertenece a `fish` o `player`, sino de datos del contexto y del owner del emisor.

### 2. La carta emite efectos; el encounter resuelve consecuencias
- La carta no deberia decidir por si misma eventos derivados como `splash`.
- La carta emite una mutacion o una intencion; capas como `progression`, `encounter` y `endings` resuelven los efectos secundarios y los cierres terminales.

### 3. Los triggers deben ser relativos al owner
- En lugar de pensar en `OnPlayerWin` o `OnFishWin`, el sistema futuro deberia poder expresar `OnOwnerWin`, `OnOwnerLose` y `OnRoundDraw`.
- Esto permite reutilizar el mismo efecto desde ambos lados sin duplicar semantica.

### 4. Las fases deben ser explicitas y estables
- Cada efecto de carta debe declarar en que fase vive.
- El motor debe poder ordenar efectos por fase antes de ordenar por prioridad o stacking.

### 5. Las reglas base siguen separadas de los efectos de carta
- El resultado base del piedra-papel-tijera actual no pasa a ser un efecto de carta.
- Los limites del `rig`, la captura, el escape y el agotamiento del mazo siguen siendo reglas del encounter o del estado global, aunque una carta pueda mutar variables que luego impacten en esos chequeos.

## Vocabulario recomendado

### Owner de carta
- `fish`
- `player`

### Trigger de carta
- `on_draw`: la carta entra en juego o se revela.
- `on_owner_win`: el owner de la carta gana la ronda.
- `on_owner_lose`: el owner de la carta pierde la ronda.
- `on_round_draw`: la ronda termina en empate.
- `on_round_end`: fase de cierre posterior a outcome y efectos principales.

Nota: estos nombres son conceptuales. La implementacion puede escoger otros nombres, pero debe conservar la semantica relativa al owner.

## Fases de resolucion propuestas

### Fase 1: draw y revelacion
- La carta entra al contexto del round.
- Aqui viven efectos tipo `on_draw`.
- Esta fase debe existir de verdad si se mantiene ese trigger en el modelo futuro.

### Fase 2: resolucion base del combate
- Se determina el outcome base a partir de jugadas y reglas de combate.
- Por defecto, esta fase no aplica efectos de carta de progresion del encounter.

### Fase 3: emision de efectos post-outcome
- Se activan efectos condicionados por `on_owner_win`, `on_owner_lose` o `on_round_draw`.
- En esta fase solo se decide que efectos quedaron habilitados.

### Fase 4: aplicacion de impactos de carta
- Se aplican las mutaciones sobre el estado objetivo.
- Para el sistema actual, aqui vive la suma sobre `DistanceShift` y `DepthShift`.

### Fase 5: eventos derivados del encounter
- El estado resultante puede disparar eventos derivados como `splash`.
- Estos eventos no son, por defecto, un payload primario de carta.

### Fase 6: chequeos terminales
- Se resuelven captura, escape y cierres por agotamiento de mazo.
- Estas reglas siguen fuera de la taxonomia de carta, aunque usan el estado ya afectado por cartas.

## Familias de efectos de carta

### A. Efectos de progresion espacial del encounter

Impactan sobre la posicion del pez dentro del encounter.

### Subtipos
- desplazamiento horizontal
- desplazamiento vertical
- modificadores combinados horizontal + vertical

### Estado actual
- Implementado parcialmente.
- Hoy se expresa con `CardEffect.DistanceShift` y `CardEffect.DepthShift`.

### Observaciones
- Esta es la familia fundacional ya presente en el proyecto.
- `splash` no debe modelarse como subtipo de esta familia; es una consecuencia derivada cuando la mutacion vertical cruza superficie.

### B. Efectos sobre recursos del jugador

Impactan recursos consumibles o enfriamientos del lado jugador.

### Ejemplos futuros
- recuperar un uso de movimiento
- bloquear temporalmente una accion
- reducir o extender una recarga

### Estado actual
- No implementado como efecto de carta.
- Hoy los recursos viven en `internal/playermoves/` y son ajenos a cartas.

### C. Efectos sobre estado del mazo o flujo de cartas

Impactan draw, descarte, reciclado o visibilidad de cartas.

### Ejemplos futuros
- robar otra carta
- descartar la carta activa
- retrasar reciclado
- mirar la siguiente carta

### Estado actual
- No implementado como efecto de carta.
- El mazo existe, pero su comportamiento no esta gobernado por efectos de carta.

### D. Efectos sobre resolucion del combate

Impactan la forma en que se decide el outcome base.

### Ejemplos futuros
- ganar empates si la carta activa cumple una condicion
- convertir un empate en victoria del owner
- invertir el resultado base

### Estado actual
- No implementado como efecto de carta.
- La resolucion base sigue viviendo en `internal/rules/`.

### Observacion importante
- Esta familia requiere una fase propia anterior a la progresion espacial.
- No debe mezclarse con efectos que solo mutan `encounter` despues del outcome.

### E. Efectos sobre thresholds o chequeos terminales

Impactan variables que luego son leidas por los cierres del encounter.

### Ejemplos futuros
- aumentar temporalmente la distancia maxima alcanzable
- modificar la profundidad maxima alcanzable
- alterar la tolerancia de captura por agotamiento de mazo

### Estado actual
- No implementado como efecto de carta.
- Hoy esos valores viven en `encounter.Config` o `playerRig` y no se alteran por carta.

### Restriccion recomendada
- Las cartas deberian mutar thresholds o estado, pero no forzar directamente `captured` o `escaped` como resultado primario.
- Si una carta altera thresholds, esos cambios deben vivir como estado temporal del round.
- El `rig` no debe mutarse por cartas; define condiciones iniciales y limites estructurales del encounter.

### F. Efectos que emiten marcas o estados temporales

Agregan contexto para fases posteriores del mismo round o rounds futuros.

### Ejemplos futuros
- carta anclada para duplicar el siguiente efecto vertical
- estado de pez cansado
- bonificacion al siguiente `on_owner_win`

### Estado actual
- No implementado.

## Matriz resumida

| Familia | Trigger tipico | Fase principal | Target | Estado actual |
| --- | --- | --- | --- | --- |
| Progresion espacial | `on_owner_win`, `on_owner_lose`, `on_round_draw` | aplicacion de impactos | `encounter` | parcial |
| Recursos del jugador | `on_draw`, `on_owner_win` | aplicacion de impactos | `player resources` | no implementado |
| Flujo de mazo | `on_draw`, `on_round_end` | draw o cierre | `deck state` | no implementado |
| Resolucion del combate | `on_draw`, previo a outcome | resolucion base | `combat result` | no implementado |
| Thresholds terminales | `on_draw`, `on_owner_win` | aplicacion de impactos | `round context` | no implementado |
| Estados temporales | cualquier trigger | varias | `round context` | no implementado |

## Mapeo de mecanicas actuales a la taxonomia

| Mecanica actual | Es efecto de carta | Familia | Observacion |
| --- | --- | --- | --- |
| `FishCard.Move` | no | n/a | es jugada base de combate, no efecto |
| `DistanceShift` en `CardEffect` | si | progresion espacial | ya implementado |
| `DepthShift` en `CardEffect` | si | progresion espacial | ya implementado |
| `splash` | no directo | evento derivado | nace del `encounter` tras aplicar profundidad |
| limite de distancia del `rig` | no | n/a | regla global, no efecto de carta |
| limite de profundidad del `rig` | no | n/a | regla global, no efecto de carta |
| captura por distancia + superficie | no | n/a | cierre terminal posterior |
| captura por agotamiento de mazo | no | n/a | cierre terminal posterior |
| subir 1 nivel al ganar en distancia de captura | no | n/a | regla base de progresion actual |

## Reglas de composicion recomendadas

### Regla 1: primero fase, luego prioridad, luego stacking
- El sistema debe ordenar primero por fase.
- Dentro de la misma fase puede existir una prioridad opcional.
- Si no existe prioridad explicita, el comportamiento por defecto debe ser aditivo para payloads compatibles.

### Regla 2: los payloads aditivos deben agruparse
- Desplazamientos horizontales y verticales pueden agregarse antes de tocar `encounter`.
- Esto preserva bien el modelo actual de `Delta`.

### Regla 3: los payloads no conmutativos requieren fase separada
- Alterar outcome base y desplazar encounter no deberian vivir en la misma bolsa de efectos.
- Si dos familias no conmutan, se resuelven en fases distintas.

### Regla 4: clamp y eventos derivados pertenecen al dominio destino
- El efecto de carta no deberia decidir como se clampa profundidad o distancia.
- `encounter` conserva la responsabilidad de transformar mutaciones en eventos derivados o estados validos.

### Regla 5: los chequeos terminales quedan al final
- Ningun efecto de carta deberia saltarse el pipeline y cerrar la partida antes de que el `encounter` termine de resolver consecuencias.

### Regla 6: ningun efecto temporal de carta dura mas de un round
- Cualquier alteracion sobre thresholds o contexto debe modelarse como estado temporal del round.
- El `rig` conserva su rol de condicion inicial del encounter y no puede ser modificado por cartas.

## Recomendaciones para la siguiente feature tecnica

### 1. Separar contrato de carta y contrato de efecto
- `FishCard` y futuras `playerCards` pueden seguir existiendo como tipos distintos.
- Pero ambos deberian emitir una misma estructura de `CardEffect` o equivalente.
- No conviene forzar una carta comun unica.
- La simetria debe vivir en contratos compartidos, mientras que las cartas pueden conservar diferencias de uso, presentacion e identidad visual.

### 2. Reemplazar triggers absolutos por triggers relativos al owner
- `TriggerOnPlayerWin` y `TriggerOnFishWin` eran un limite del modelo previo.
- Esta migracion ya se implemento con triggers compatibles con cualquier owner de carta.

### 3. Hacer real la fase `on_draw` o eliminarla
- No conviene mantener un trigger sin fase real de ejecucion.
- La implementacion debe decidir entre soportarlo de verdad o sacarlo hasta necesitarlo.

### 4. Introducir un contexto de round orientado a efectos
- `match.ResolvedRound` ya es un buen punto de partida.
- La siguiente iteracion probablemente necesite un contexto mas rico que permita registrar owner, carta emisora, outcome relativo y efectos habilitados.

### 5. Mantener `encounter`, `rules` y `endings` desacoplados
- `rules` sigue resolviendo combate base.
- `progression` o una capa de efectos aplica payloads de carta.
- `encounter` resuelve clamps y eventos derivados.
- `endings` conserva los chequeos terminales.

## Decisiones cerradas por este discovery

- El scope de esta linea de trabajo son efectos de cartas, no efectos genericos de cualquier sistema.
- El sistema futuro debe ser agnostico al owner de la carta.
- Los efectos se clasifican por trigger, fase, familia de impacto y target.
- `splash` se considera un evento derivado del `encounter`, no un efecto primario de carta.
- La siguiente feature tecnica debe partir de un contrato comun reutilizable por `FishCard` y `playerCards`.
- No se modelara una carta comun unica; se compartiran contratos comunes entre tipos de carta distintos.
- Las diferencias de uso y presentacion entre cartas del pez y del jugador son validas y esperables.
- Los thresholds alterados por carta viven como estado temporal del round.
- El `rig` no puede ser modificado por cartas.
- Ningun efecto temporal de carta dura mas de un round.
- Los efectos futuros sobre outcome base deben seguir viviendo en `rules` mediante hooks mientras esa integracion escale de forma razonable.
- Solo si ese enfoque deja de escalar se evaluara una capa nueva anterior a `rules`.
- Esta feature ya implementa una primera version de contratos compartidos para efectos de carta y hooks optativos en `rules`.
