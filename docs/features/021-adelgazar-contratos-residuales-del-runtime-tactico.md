# Plan de feature: adelgazar-contratos-residuales-del-runtime-tactico

## Objetivo

Cerrar la deuda residual que sigue concentrada en contratos transicionales del duelo despues de `BL-030`, especialmente en `match.RoundResult`, el mapeo del estado visible del mazo y las mutaciones demasiado anchas sobre `match.State`, para que el runtime tactico quede listo para crecer sin volver a recentralizar conocimiento en `internal/game/`.

La meta de esta feature no es rehacer otra vez la arquitectura de combate ni mover paquetes por estetica. La meta es completar el estrechamiento de superficies que `BL-029` y `BL-030` dejaron encaminado, de forma que `BL-023` y el futuro runtime de run puedan apoyarse en contratos tacticos mas claros y menos acoplados.

## Criterios de aceptacion

- `match.RoundResult` expone solo la informacion tactica que realmente necesitan sus consumidores.
- El estado visible del mazo deja de depender de mapeos ensamblados desde `game` cuando exista un owner mas natural.
- `progression` y `endings` reducen su dependencia del ensamblado tactico completo o dejan documentado por que aun no conviene estrecharla mas.
- La refactorizacion deja mas claro que parte del runtime lee snapshots, que parte muta subestados y que parte solo orquesta el round.
- La salida funcional del duelo no cambia y la cobertura relevante sigue pasando.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope

### Incluye

- Revisar el contrato real de `match.RoundResult` y adelgazarlo si su superficie actual es mas ancha de lo necesario.
- Reubicar helpers o mapping del estado visible del mazo si hoy viven en `game` sin ser su owner natural.
- Evaluar contratos mas finos para `progression` y `endings`, ya sea via subestados, helpers tacticos o interfaces pequenas con ownership claro.
- Ajustar `app/session`, `presentation` y tests afectados por esa frontera mas estrecha.
- Documentar la direccion resultante si hace falta fijar decisiones de ownership.

### No incluye

- Rehacer el flujo UI-agnostic completo; eso sigue en `BL-023`.
- Introducir runtime de run, servicios o economia.
- Reorganizar todos los paquetes del combate otra vez sin mejora real de ownership.
- Abrir una capa nueva de snapshots si solo duplica datos sin reducir coupling real.

## Propuesta de implementacion

### 1. Adelgazar `match.RoundResult`

`RoundResult` sigue siendo un punto sensible porque puede arrastrar mas estado del necesario hacia `app/session`, `presentation` u otros consumidores.

Direccion propuesta:

- identificar que campos leen realmente los consumidores del resultado de ronda
- mover el resto a snapshots o lecturas derivadas mas especificas cuando aporte claridad
- evitar que el resultado de ronda actue como un segundo ensamblado universal del duelo

La meta no es tener el tipo mas pequeno posible, sino uno que exprese el outcome del round sin reexportar accidentalmente todo el runtime tactico.

### 2. Mover el mapping del estado visible del mazo cerca de su owner

Si `game` sigue traduciendo visibilidad de descarte o estado de mazo desde tipos de otra capa, el engine mantiene conocimiento que no deberia poseer.

Direccion propuesta:

- localizar el mapping exacto que hoy hace `game`
- evaluar si debe vivir en `match`, `encounter`, una capa de snapshots tacticos o un helper de lectura claramente nombrado
- dejar al engine consumiendo una operacion lista para usar, no ensamblando representaciones visibles

### 3. Estrechar la superficie de `progression` y `endings`

`progression` y `endings` ya tienen ownership mas claro que antes, pero todavia pueden apoyarse en `*match.State` de forma mas ancha de lo deseable.

Direccion propuesta:

- detectar que partes de `match.State` necesitan mutar o leer de verdad
- introducir contratos pequenos o acceso por subestados cuando eso reduzca coupling de forma legible
- preferir soluciones directas y locales antes que interfaces abstractas demasiado generales

Si algun punto no conviene estrecharlo todavia, debe quedar explicitado para no reabrir la misma discusion en cada feature futura.

### 4. Ajustar consumidores y documentacion

El estrechamiento de contratos puede afectar a `app/session`, `presentation` y tests de integracion del combate.

Direccion propuesta:

- adaptar consumidores a los nuevos resultados o snapshots
- mantener la salida actual del presenter y del CLI
- dejar una nota corta en la documentacion tactica si la nueva frontera cambia ownership relevante

### 5. Slice minimo recomendado

Para mantener el scope controlado, el primer slice deberia cubrir:

- adelgazamiento aplicado de `match.RoundResult`
- un movimiento concreto del mapping de estado visible del mazo fuera de `game`
- un estrechamiento real en `progression` o `endings`
- ajuste de tests y documentacion minima

Con eso el runtime tactico quedara mas preparado para `BL-023` y para el futuro runtime de run sin volver a inflar el engine.

## Validacion

- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Verificar manualmente que el flujo actual `setup -> opening -> spawn -> combate` sigue funcionando.
- Verificar que el engine ya no necesita conocer detalles visibles del mazo o resultados de ronda mas anchos de lo necesario.

## Riesgos o decisiones abiertas

- Si se intenta adelgazar todos los contratos a la vez, el scope puede crecer y mezclarse con una nueva ronda de arquitectura global.
- Si `RoundResult` se reduce demasiado, algunos consumidores pueden acabar recomponiendo lecturas dispersas peores que el contrato actual.
- Habra que distinguir entre coupling real y simple acceso a un subestado ya razonablemente owned.
- `BL-023` puede querer reutilizar parte de estos contratos, asi que conviene mantenerlos semanticos y no overly tailored a una sola interfaz.
