# Plan de feature: descomponer-match-state-y-fronteras-runtime-tactico

## Objetivo

Reducir el acoplamiento actual del runtime tactico del duelo separando responsabilidades que hoy viven agrupadas en `match.State`, de forma que el combate siga siendo jugable pero deje de crecer como un estado universal compartido entre engine, progression, endings, presentation y futuras capas de run.

La meta de esta feature no es rehacer todo el sistema de combate ni introducir aun el runtime de expedicion. La meta es fijar fronteras mas sanas para el estado mutable del duelo, dejar ownership claro de encounter, mazo, recursos del jugador, thresholds de ronda y stats, y preparar el terreno para que `BL-001` y `BL-030` no se monten sobre un objeto central sobredimensionado.

## Criterios de aceptacion

- `match.State` deja de ser el contenedor principal de todas las piezas tacticas del duelo en una sola estructura plana.
- Quedan definidos subestados o contratos tacticos con ownership explicito para encounter, deck, recursos del jugador, estado de ronda y estadisticas.
- `game`, `progression`, `endings` y `presentation` dependen de superficies mas estrechas o de un ensamblado tactico mas claro.
- La presentacion deja de necesitar conocer mas datos tacticos de los que realmente muestra.
- La refactorizacion preserva el comportamiento actual del duelo y el contrato funcional de `app/session`.
- La nueva estructura deja claro por que `RunState` futuro no debe vivir dentro del estado tactico del encounter.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope

### Incluye

- Revisar y explicitar las responsabilidades que hoy se mezclan en `match.State`.
- Definir la estructura objetivo del runtime tactico y los subestados que lo componen.
- Refactorizar el estado compartido del duelo para reflejar esas fronteras.
- Ajustar engine, progression, endings, presenter y tests afectados por el nuevo ensamblado tactico.
- Documentar la frontera entre runtime tactico y futuro runtime de expedicion.

### No incluye

- Implementar todavia `RunState`, mapa, economia de run o servicios.
- Reorganizar por completo todos los paquetes de combate; esa convergencia se aborda despues en `BL-030`.
- Cambiar el comportamiento de gameplay del duelo salvo donde una aclaracion estructural lo exija de forma acotada.
- Redisenar aun la UX o hacer la migracion UI-agnostic completa de setup/opening.

## Propuesta de implementacion

### 1. Partir `match.State` por ownership real

Hoy `match.State` mezcla al menos cinco responsabilidades distintas:

- progreso tactico del encounter
- estado del mazo del pez y descarte visible
- recursos y barajas del jugador
- thresholds y efectos scoped a la ronda
- estadisticas y flag terminal del duelo

Direccion propuesta:

- mantener un ensamblado tactico de alto nivel para el engine
- pero dividir su interior en subestados con nombres y ownership estables

Ejemplo de direccion esperada:

- `EncounterState` o reutilizar `encounter.State` como owner exclusivo del pez y del tablero
- `DeckState` como snapshot tactico del mazo y descarte
- `PlayerResourcesState` para moves/loadout del jugador
- `RoundState` para thresholds y efectos efimeros de un turno
- `MatchStats` o equivalente para contadores de resultados

La meta no es cosmetica; es que cada capa muta solo lo que le corresponde.

### 2. Estrechar la superficie de mutacion entre paquetes

El engine hoy inicializa, refresca, progresa y cierra el encounter tocando directamente gran parte del estado.

Direccion propuesta:

- reducir los puntos donde varias capas mutan el mismo struct global
- dejar mas claro que parte del estado pertenece a:
  - `game`
  - `progression`
  - `endings`
  - `presentation`
- si hace falta, introducir helpers o wrappers tacticos pequenos antes que interfaces especulativas grandes

### 3. Fijar la frontera con la futura run

Esta feature debe dejar muy explicito que el runtime de run no pertenece a `match.State`.

Resultado buscado:

- el estado tactico representa solo un duelo aislado
- el futuro runtime de expedicion consumira ese duelo como un paso dentro de una run, no como su contenedor global
- la refactorizacion facilita que `BL-001` introduzca `RunState` sin volver a mezclar capas

### 4. Reducir dependencia de presentation sobre estado bruto

`presentation` hoy consume `match.State` entero para construir vistas.

Direccion propuesta:

- revisar si `presentation` debe leer el ensamblado tactico completo o snapshots mas pequenos
- quitar dependencias innecesarias del presenter sobre piezas que solo existen porque el estado actual esta demasiado centralizado
- mantener la salida actual de CLI, pero con una frontera mas clara entre runtime y vistas

### 5. Slice minimo recomendado

Para mantener la iteracion cerrable, el primer slice de implementacion deberia cubrir:

- extraer subestados tacticos estables dentro de `match/` o un hogar equivalente cercano
- actualizar el engine para operar sobre esa nueva estructura sin cambiar comportamiento
- ajustar progression/endings para mutar solo las piezas necesarias
- adaptar presentation y tests a los nuevos accesos
- dejar una nota o README corto si hace falta para fijar la frontera tactica resultante

Con eso `BL-030` ya podra consolidar mejor los paquetes del combate y `BL-001` podra introducir runtime de run sobre una base mas sana.

## Validacion

- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Verificar manualmente que el flujo actual `setup -> opening -> spawn -> combate` sigue funcionando.
- Verificar que un cambio en recursos del jugador, deck o thresholds de ronda no requiere volver a tocar todas las capas consumidoras del estado tactico.

## Riesgos o decisiones abiertas

- Si se intenta mover demasiadas fronteras de paquetes a la vez, `BL-029` puede solaparse en exceso con `BL-030`.
- Si la descomposicion introduce demasiada abstraccion temprana, puede empeorar la legibilidad en vez de mejorarla.
- Habra que decidir si los nuevos subestados viven aun bajo `match/` como capa transicional o si alguna pieza ya merece migrar a un paquete propio.
- El presenter puede seguir necesitando parte del estado combinado; el objetivo es estrechar la frontera, no duplicar snapshots sin necesidad.
