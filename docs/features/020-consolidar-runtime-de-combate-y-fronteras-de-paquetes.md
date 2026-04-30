# Plan de feature: consolidar-runtime-de-combate-y-fronteras-de-paquetes

## Objetivo

Reducir la fragmentacion funcional actual del runtime de combate para que la logica del duelo deje de estar repartida de forma difusa entre `internal/game/`, `internal/rules/`, `internal/progression/`, `internal/endings/`, `internal/encounter/` y `internal/presentation/`, dejando ownership mas claro de evaluacion, progresion, cierre, snapshots tacticos y presentacion.

La meta de esta feature no es un gran rewrite ni una mudanza completa a un paquete `combat/` definitivo en una sola iteracion. La meta es consolidar las fronteras principales del duelo sobre la base ya saneada por `BL-029`, revisar si `presentation` puede depender de snapshots mas estrechos que el ensamblado tactico completo, y dejar una arquitectura mas estable para el crecimiento de la run, roles de encounter y futuras interfaces.

## Criterios de aceptacion

- Queda explicitado que parte del duelo pertenece a `encounter`, que parte a `rules`, que parte a `progression`, que parte a `endings` y que parte a `game`.
- La implementacion reduce imports cruzados o conocimiento accidental entre esas capas.
- `presentation` deja documentado y, cuando sea razonable, aplicado un consumo de snapshots mas estrechos que `match.State` completo.
- El engine sigue orquestando el round completo sin absorber mas logica de negocio de la necesaria.
- La reorganizacion preserva el comportamiento actual del combate y su capacidad de testeo.
- La direccion resultante deja mas claro como introducir despues runtime de run, servicios, bosses y otras capas sin remezclar el duelo.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope

### Incluye

- Revisar y fijar ownership de las piezas principales del runtime de combate.
- Reducir fragmentacion funcional o dependencias opacas entre paquetes del duelo.
- Extraer o introducir snapshots tacticos mas estrechos cuando aporten claridad real, especialmente para `presentation`.
- Ajustar `engine`, `presentation` y las capas tacticas afectadas por esa consolidacion.
- Documentar la direccion resultante para el runtime de combate.

### No incluye

- Mover todo el duelo a un unico paquete definitivo por fuerza.
- Rehacer la UI o resolver aun la migracion UI-agnostic completa; eso sigue en `BL-023`.
- Implementar runtime de run, mapa o economia.
- Rebalancear reglas de gameplay salvo donde la consolidacion obligue a aclarar ownership o politicas existentes.

## Propuesta de implementacion

### 1. Fijar el mapa de ownership del duelo

La primera tarea de la feature es dejar claro el reparto esperado:

- `encounter/`
  - runtime propio del pez y del tablero tactico
  - configuracion y reglas de apertura o delta propias del encounter
- `rules/`
  - evaluacion del outcome del round
- `progression/`
  - traduccion de outcome y efectos a cambios de distancia/profundidad/eventos del encounter
- `endings/`
  - cierre del duelo y condiciones terminales
- `game/`
  - orquestacion del round, pipeline de fases y coordinacion de colaboradores
- `presentation/`
  - traduccion de snapshots tacticos a vistas, no acceso indiscriminado al runtime completo

El resultado buscado no es solo documental; debe guiar los cambios concretos de la feature.

### 2. Revisar la frontera del engine

`game.Engine` ya coordina bien el round, pero sigue siendo un punto sensible de acoplamiento porque refresca estado, aplica efectos, invoca progression y cierre, y devuelve resultados enriquecidos.

Direccion propuesta:

- mantener a `game` como orquestador del pipeline del round
- evitar que `game` siga absorbiendo detalles que deberian vivir en otras capas tacticas
- revisar si algunos helpers o snapshots hoy privados del engine deben moverse cerca de sus owners reales

### 3. Revisar snapshots mas estrechos para `presentation`

Hoy `presentation` consume `match.State` entero para construir `StatusView`, `RoundView` y `SummaryView`.

La feature debe revisar explicitamente si conviene estrechar esa dependencia. Direccion recomendada:

- identificar la minima informacion tactica que `presentation` necesita de verdad
- introducir snapshots o structs de lectura mas pequenos si reducen coupling real sin duplicar ruido innecesario
- priorizar especialmente:
  - status del tablero
  - descarte visible del pez
  - recursos y cartas del jugador visibles
  - stats y resumen final

No es obligatorio que toda `presentation` deje de leer el ensamblado tactico en esta iteracion, pero si debe quedar una direccion aplicada al menos en los puntos mas claros.

### 4. Consolidar dependencias entre paquetes tacticos

Sobre la base de `BL-029`, esta feature debe reducir la sensacion de capacidad repartida sin ownership claro.

Posibles movimientos aceptables:

- introducir snapshots tacticos compartidos en un lugar mas estable
- mover helpers que hoy viven donde no son dueños naturales
- simplificar superficies entre progression/endings/engine
- documentar por que una pieza sigue en un paquete transicional si aun no conviene moverla

La meta es dejar una arquitectura mas coherente, no mover archivos por estetica.

### 5. Slice minimo recomendado

El primer slice deberia cubrir:

- mapa de ownership explicitado en codigo y plan
- revision aplicada del consumo de `presentation`
- uno o dos movimientos concretos de consolidacion entre `game`, `progression`, `endings` y `encounter`
- actualizacion de tests y documentacion tactica

Con eso quedara una base mas estable para que `BL-023` y `BL-001` no sigan creciendo sobre una frontera de combate ambigua.

## Validacion

- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Verificar manualmente que el flujo actual `setup -> opening -> spawn -> combate` sigue funcionando.
- Verificar que `presentation` puede consumir informacion mas estrecha al menos en parte del flujo sin perder legibilidad.

## Riesgos o decisiones abiertas

- Si intentamos converger toda la arquitectura final del combate en una sola feature, el scope crecera demasiado.
- Si se introducen demasiados snapshots finos, puede aparecer duplicacion o ruido innecesario.
- Habra que decidir que piezas siguen en paquetes transicionales y cuales ya merecen migrar a una frontera mas estable.
- `BL-023` podria querer reutilizar algunos snapshots o contratos nuevos; conviene dejarlos suficientemente generales sin caer en abstraccion especulativa.
