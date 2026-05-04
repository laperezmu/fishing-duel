# Plan de feature: definir-flujo-minimo-de-zona-y-nodos-para-el-mvp

## Objetivo

Definir e implementar el flujo minimo navegable de una run para el MVP, bajando la expedicion a una secuencia corta de nodos y transiciones que permita jugar una run completa sin mapa rico. Esta feature debe apoyarse en el contrato de `run` ya fijado por `BL-001`, pero ademas resolver una necesidad estructural nueva: separar el ejecutable CLI actual de encounter prototipo del futuro ejecutable CLI de run, para que el wiring del primer experimento no siga condicionando la composicion del juego real.

La meta no es reemplazar inmediatamente el prototipo actual, sino introducir un punto de entrada nuevo y limpio para la run MVP, mientras el binario existente sigue sirviendo como sandbox de encounter aislado, setup manual y pruebas de contenido tactico.

## Criterios de aceptacion

- Existe un ejecutable CLI dedicado a jugar la run MVP, separado del ejecutable actual de encounter prototipo.
- El ejecutable actual sigue funcionando como modo de prueba del encounter aislado.
- Queda definida una taxonomia minima de nodos para el MVP, al menos cubriendo `start`, `fishing`, `service/checkpoint`, `boss` y `end`.
- Queda definido un flujo secuencial inicial equivalente a `inicio -> pesca -> pesca -> servicio -> boss -> cierre`, aunque luego pueda ajustarse en detalle.
- La capa de aplicacion expone un wiring claro para recorrer nodos de run sin empujar esa logica a `cmd/`.
- Los puntos donde la run entra a encounter, consume su resultado y avanza al siguiente nodo quedan explicitados con contratos pequenos y estables.
- La documentacion y README dejan claro que binario corresponde al prototipo de encounter y cual corresponde a la run MVP.

## Scope

### Incluye

- Analizar el estado actual del binario `cmd/fishing-duel/` y recortarlo como ejecutable de prototipo/encounter.
- Introducir un nuevo `cmd/` para la run MVP con composition root separada.
- Definir el flujo secuencial minimo de nodos de la run MVP.
- Introducir contratos minimos de navegacion/orquestacion de run sobre la base de `internal/run/`.
- Definir como se inyecta un `encounter` dentro de un nodo de pesca y como su resultado vuelve a la run.
- Ajustar documentacion y naming para evitar confusion entre ejecutable de prototipo y ejecutable real de run.

### No incluye

- Mapa rico con bifurcaciones, ocultacion o generacion procedural.
- Economia ampliada, patrocinadores completos o meta-progresion.
- Eliminar el ejecutable actual de prototipo.
- Resolver toda la UX final de la run; primero importa la estructura y el wiring.

## Analisis de la situacion actual

### 1. El `main` actual ya no representa bien la direccion del producto

Hoy `cmd/fishing-duel/main.go` arranca un encounter aislado con setup manual completo:

- elige preset del jugador
- elige `rod`
- elige aditamentos
- elige agua
- resuelve cast
- resuelve spawn
- juega un duelo unico

Ese flujo fue util para iterar el sistema tactico y sigue siendo valioso como herramienta de prototipo. Pero para `BL-035` deja de ser un buen candidato a evolucion natural del juego, porque la run necesita otra composicion:

- arranque de expedicion
- secuencia de nodos
- persistencia de recursos entre encounters
- servicios o checkpoints
- cierre de run

Si intentamos meter ese flujo dentro del mismo `main`, el binario mezclaria dos responsabilidades distintas:

- sandbox de encounter
- juego real de run

### 2. El riesgo de seguir con un solo binario

Mantener un solo `cmd/fishing-duel/` para ambas experiencias tiene varios costos:

- obliga a llenar el `main` de flags, ramas o wiring condicional
- difumina que es prototipo y que es producto jugable de run
- mete presion innecesaria sobre `internal/cli/UI`, que hoy esta muy orientada al setup manual del encounter
- vuelve mas dificil evolucionar la UX de run sin romper el flujo de pruebas del duelo aislado

La separacion de mains no es solo una preferencia organizativa; es una forma de preservar una frontera util entre dos modos de uso legitimos del repo.

### 3. Que valor tiene conservar el binario actual

El ejecutable actual sigue siendo valioso incluso despues de abrir la run MVP:

- permite probar encounter, cast, spawn y contenido de pez sin recorrer una run completa
- sirve como harness manual de debugging para iterar presets, pools y balance tactico
- reduce friccion al validar cambios en sistemas del duelo que todavia evolucionan por separado

Por eso la recomendacion no es reemplazarlo, sino renombrar su rol conceptual: pasar de ser "el juego principal" a ser "el prototipo/manual runner del encounter".

### 4. Que necesita realmente el nuevo binario de run

El nuevo `main` de run no deberia pedir al usuario todo el setup tactico actual en cada inicio. Su responsabilidad deberia ser otra:

- crear o cargar una definicion minima de ruta secuencial
- inicializar un `run.State`
- orquestar el avance nodo a nodo
- delegar al encounter solo cuando el nodo lo requiera
- aplicar el `EncounterResult` a la run
- cerrar la expedicion con victoria, derrota o retiro

Eso hace que el nuevo ejecutable necesite contratos de app distintos del binario actual, aunque ambos compartan presenter, UI base o partes del adaptador CLI.

## Direccion propuesta

### 1. Separar ejecutables por intencion

Direccion recomendada:

- conservar `cmd/fishing-duel/` como ejecutable de encounter prototipo
- introducir un nuevo `cmd/` para la run MVP, por ejemplo `cmd/fishing-run/` o nombre equivalente mas claro

Recomendacion concreta:

- `cmd/fishing-duel/`: duelo aislado, setup manual, sandbox tactico
- `cmd/fishing-run/`: run MVP, secuencia de nodos, flujo principal del producto

Esto deja un lenguaje claro tambien para la documentacion y para futuros tests manuales.

### 2. Mantener `cmd/` como composition root minima

Igual que en `BL-034`, el nuevo `main` no debe resolver logica de negocio de run. Debe limitarse a:

- construir dependencias
- elegir el runner de run
- conectar adaptadores CLI/presentation
- ejecutar el flujo principal

La orquestacion real de nodos y transiciones debe vivir en `internal/app/` y/o `internal/run/`.

### 3. Introducir una primera capa de orquestacion de run

Para que el nuevo binario no quede vacio ni demasiado listo-para-todo, conviene introducir un caso de uso acotado, por ejemplo:

- `RunSession`
- `RunController`
- `PlayRun`

Responsabilidades minimas:

- inicializar `run.State`
- exponer el nodo actual
- decidir el siguiente paso valido
- disparar un encounter cuando el nodo es de pesca o boss
- aplicar `run.EncounterResult`
- marcar el cierre de la expedicion

### 4. Definir una ruta secuencial fija antes del mapa real

Para `BL-035` no hace falta un sistema de grafo general. Alcanza con una definicion fija y legible de nodos, por ejemplo:

- `start`
- `fishing-1`
- `fishing-2`
- `service-1`
- `boss-1`
- `end`

Eso permite validar de inmediato:

- persistencia de estado entre nodos
- handoff `run -> encounter -> run`
- puntos seguros de retiro
- cierre completo de la expedicion

### 5. No acoplar el nodo al setup manual del encounter

El nodo de pesca no deberia pedir el mismo setup interactivo completo que hoy usa el prototipo, salvo como medida temporal muy controlada. La direccion correcta para la run es:

- la run decide el contexto del nodo
- el bootstrap del encounter consume ese contexto
- el adaptador CLI solo muestra y recoge las decisiones que realmente sigan siendo del jugador en ese punto

Si el nuevo binario reutiliza al principio partes del setup actual, deberia hacerlo detras de contratos que luego permitan recortarlo, no incrustando otra vez el flujo completo en el `main`.

## Propuesta de implementacion

### 1. Reposicionar el binario actual como prototipo

- actualizar `README.md` y `cmd/fishing-duel/README.md`
- dejar claro que `cmd/fishing-duel/` es el runner de encounter aislado
- mantener su wiring actual mientras siga siendo util para pruebas de dominio

### 2. Crear un nuevo ejecutable para run

- agregar `cmd/fishing-run/main.go`
- componer ahi la UI, presenter y runner de run
- evitar que este nuevo binario dependa de decisiones manuales que pertenecen al prototipo salvo donde sean realmente necesarias en esta etapa

### 3. Crear un runner minimo de run en `internal/app/`

- introducir una sesion/controlador de run reutilizable
- recibir una ruta secuencial fija y un estado inicial de run
- devolver puntos claros de extension para `BL-036`, `BL-037` y `BL-038`

### 4. Modelar la secuencia minima de nodos

- definir catalogo o slice fijo del MVP en `internal/run/` o modulo equivalente
- representar nodo actual, nodo siguiente y tipos de nodo
- soportar al menos avance lineal y nodos terminales

### 5. Integrar encounters como efecto de nodos

- nodos `fishing` y `boss` disparan un encounter
- el resultado se transforma con `ResolveEncounterResult(...)`
- la run decide como avanzar o cerrar a partir de ese resultado

### 6. Dejar preparado el borde para la UX de run

- si hace falta, crear contratos CLI separados para run en vez de seguir cargando todo sobre `internal/cli.UI`
- evitar que los metodos pensados para el sandbox de encounter definan por accidente la UX completa de la run

## Archivos o zonas probables

- `cmd/fishing-duel/main.go`
- `cmd/fishing-duel/README.md`
- `cmd/fishing-run/main.go`
- `internal/app/`
- `internal/run/`
- `internal/cli/`
- `README.md`

## Riesgos o decisiones abiertas

- Hay que decidir si el nuevo binario de run reutiliza temporalmente el setup actual del jugador o si arranca con un loadout fijo del MVP.
- Hay que decidir si el nombre final del ejecutable de run debe enfatizar `run`, `roguelike` o `expedicion`; mi recomendacion es algo explicito y corto como `cmd/fishing-run/`.
- Hay que vigilar que `internal/cli.UI` no termine absorbiendo responsabilidades de dos experiencias distintas sin una frontera clara.
- Conviene no sobredisenar un sistema de nodos generico antes de validar la secuencia lineal base.

## Accionables inmediatos

- Pasar `BL-035` a `planned`.
- Crear el nuevo composition root de run.
- Definir la ruta fija inicial del MVP.
- Introducir el runner minimo de run en `internal/app/`.
- Ajustar la documentacion de ejecucion para explicar ambos binarios.

## Extension acordada despues del primer slice

- En la run, el agua ya no debe elegirse manualmente.
- Los nodos `fishing` y `boss` deben proveer o resolver su propio preset de agua.
- El jugador solo recibe los hints del agua del nodo actual y luego resuelve el cast.
- El sandbox `cmd/fishing-duel/` puede seguir conservando la eleccion manual para debugging.
