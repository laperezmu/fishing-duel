# Plan de feature: rod-y-limites-de-apertura-y-track

## Objetivo

Convertir la actual capa de `rig` del jugador en un modelo explicito de `rod` que separe dos familias de limites del encounter: los limites de apertura que validan hasta donde puede empezar la pesca despues del cast, y los limites de track que definen el tablero real del duelo y las condiciones de escape.

La meta de esta feature no es resolver aun subpools completos de peces ni la economia de build. La meta es fijar el vocabulario tecnico correcto para la cana base del jugador, eliminar la ambiguedad de `MaxDistance` y `MaxDepth`, y conectar la apertura vertical del encounter con una primera resolucion basada en `rod`, dejando el sistema listo para que `BL-022` use la apertura ya resuelta y no los limites de escape del tablero.

## Criterios de aceptacion

- El termino `rig` deja de usarse como nombre de dominio para la cana base del jugador y pasa a `rod`.
- El modelo del jugador separa limites de apertura y limites de track en horizontal y vertical.
- Los limites de track siguen gobernando render del tablero y condiciones de escape del encounter.
- Los limites de apertura se usan para validar o resolver `InitialDistance` y `InitialDepth` antes de crear el `encounter.State`.
- La apertura del encounter puede producir una profundidad inicial derivada de la `rod`, aunque el cast siga resolviendo solo la banda horizontal.
- La implementacion deja espacio para introducir aditamentos despues sin volver a mezclar `rod` con setup completo del jugador.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope

### Incluye

- Renombrar la capa actual de `playerrig` a un modelo de `rod`.
- Separar limites de apertura y track en el estado base del jugador.
- Integrar esos limites en apertura, endings, render y presentacion.
- Resolver una primera profundidad inicial desde la `rod` antes de arrancar el encounter.
- Ajustar documentacion y tests para reflejar el nuevo vocabulario.

### No incluye

- Implementar sistema de compra, slots complejos o inventario persistente de aditamentos.
- Resolver todavia subpools reales de pez por habitat o profundidad; eso queda para `BL-022`.
- Redisenar el mapa o los nodos de expedicion.
- Introducir cartas que muten la `rod` durante el duelo.
- Cerrar todavia una UX final de seleccion de equipo para run completa.

## Propuesta de implementacion

### 1. Sustituir `rig` por `rod` como pieza base del jugador

La capa actual de `playerrig` ya expresa limites estructurales del duelo, pero el nombre arrastra ambiguedad y dificulta separar la pieza base del setup completo.

Direccion propuesta:

- Renombrar el paquete actual a `internal/player/rod`.
- Renombrar `Config` y `State` para expresar una `rod` base del jugador.
- Reservar la idea de setup o loadout para el conjunto futuro `rod + aditamentos`, sin implementarlo aun como sistema completo.

### 2. Separar limites de apertura y limites de track

La misma variable no debe servir a la vez para decir hasta donde puede empezar la pesca y hasta donde puede mantenerse el pez sin escapar.

Propuesta base para la `rod`:

- `OpeningMaxDistance`
- `OpeningMaxDepth`
- `TrackMaxDistance`
- `TrackMaxDepth`

Semantica esperada:

- `Opening*` limita la apertura del encounter tras resolver contexto de agua, cast y sesgo vertical de la `rod`.
- `Track*` define el tablero visible del duelo y las condiciones de escape horizontal o vertical.

### 3. Resolver una apertura vertical minima desde la `rod`

La feature no necesita aun un minijuego vertical. Basta con que la `rod` pueda influir en la profundidad inicial de forma estructural y testeable.

Primera traduccion recomendada:

- `WaterContext` sigue aportando una profundidad base.
- La `rod` define hasta que profundidad puede abrirse la pesca.
- La apertura resuelta usa la profundidad base del contexto, acotada por `OpeningMaxDepth`.
- Si en el futuro hay aditamentos, esa resolucion podra desplazarse antes de validar el valor final.

En horizontal:

- El cast sigue resolviendo la banda.
- La `rod` valida que la distancia inicial final no supere `OpeningMaxDistance`.

### 4. Integracion con el encounter actual

La separacion debe entrar en los puntos donde hoy se usa `MaxDistance` y `MaxDepth` con doble significado.

Impactos previstos:

- `internal/app/`
  - coordinar la resolucion de apertura con la `rod` antes de crear el encounter.
- `internal/encounter/`
  - aceptar una apertura validada contra limites de `Opening*`.
- `internal/endings/`
  - usar solo `TrackMaxDistance` y `TrackMaxDepth` para escapes.
- `internal/presentation/` y `internal/cli/`
  - seguir mostrando el tablero segun limites de `Track*`.
- `cmd/fishing-duel/`
  - crear estado de `rod` en vez de `playerrig` y pasarlo al flujo de bootstrap.

### 5. Slice minimo de implementacion

Para mantener la feature cerrable, el primer slice deberia cubrir:

- renombre de `playerrig` a `rod`
- nuevos campos `Opening*` y `Track*`
- integracion de esos campos en summary, tablero y endings
- validacion de la apertura contra `OpeningMaxDistance`
- resolucion de `InitialDepth` acotada por `OpeningMaxDepth`
- actualizacion de tests y backlog

Con eso ya queda fijada la frontera entre apertura y tablero antes de introducir aditamentos o spawn contextual de peces.

### 6. Slices siguientes recomendados dentro de `BL-021`

Una vez cerrado el renombre y la separacion de limites, conviene encadenar pequeños pasos que dejen preparada la capa de setup sin mezclarla aun con economia o tienda.

#### a. Presets o seleccion de `rod` en CLI

- Exponer al menos 2 o 3 presets de `rod` para testing manual.
- Permitir validar rods con distinto alcance de apertura y distinto margen de track.
- Mantener la seleccion simple, igual que los presets actuales de jugador, pez y agua.

#### b. Contrato minimo de setup o loadout

- Introducir una estructura ligera que reserve el espacio de `rod + aditamentos`.
- Permitir ya una primera capa de aditamentos con modificadores de apertura y track, sin abrir aun la economia de build.
- Evitar que futuras features vuelvan a meter toda la semantica del equipo dentro de `rod`.

#### c. Puente directo hacia `BL-022`

- Dejar explicito que la aparicion del pez mirara la apertura ya resuelta (`InitialDistance`, `InitialDepth`).
- Evitar que `BL-022` dependa de los limites de track o del tablero visible para decidir spawn.
- Preparar metadata minima para que el resolver futuro pueda leer la apertura efectiva del lance y no solo los maximos del equipo.

## Validacion

- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Validar manualmente que el tablero sigue renderizando con los limites de track.
- Validar manualmente que una `rod` con menos alcance de apertura no permite empezar mas lejos o mas profundo aunque el contexto de agua lo proponga.

## Riesgos o decisiones abiertas

- Habra que decidir despues si `OpeningMaxDepth` solo recorta la profundidad inicial o si tambien habilita bandas discretas de apertura vertical.
- Si `rod` y setup no quedan bien diferenciados ahora, los aditamentos volveran a contaminar el modelo base.
- `BL-022` necesitara decidir si el spawn mira solo `InitialDistance` e `InitialDepth` o tambien la ventana disponible antes de aplicar clamps.
- Queda abierto si en una siguiente iteracion los aditamentos se eligen uno por uno con slots reales o si pasan antes por presets de loadout mas completos.
