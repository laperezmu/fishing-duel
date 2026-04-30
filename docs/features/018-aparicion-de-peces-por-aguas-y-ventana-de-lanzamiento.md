# Plan de feature: aparicion-de-peces-por-aguas-y-ventana-de-lanzamiento

## Objetivo

Resolver la aparicion del pez a partir del contexto de agua elegido, la apertura ya resuelta del lance y el setup efectivo del jugador, de forma que el encounter deje de depender de una eleccion manual del preset de pez y empiece a comportarse como una seleccion contextual de subpool.

La meta de esta feature no es cerrar todavia el sistema final de nodos, zonas, elites o economia editorial. La meta es introducir un primer resolver jugable de aparicion que conecte `WaterContext`, `InitialDistance`, `InitialDepth`, `rod`, aditamentos y habitats dentro del flujo actual del duelo aislado.

## Criterios de aceptacion

- La seleccion manual del preset de pez deja de ser el flujo principal del juego.
- Los perfiles de pez exponen metadata minima de aparicion por agua, distancia, profundidad y habitats.
- Existe un resolver determinista de aparicion a partir de la apertura resuelta del encounter y del loadout del jugador.
- El resolver puede distinguir al menos entre peces de costa, mar abierto, corriente mixta y sesgos de habitat aportados por aditamentos.
- El juego sigue construyendo el mazo del pez desde perfiles data-driven, pero el perfil elegido ya no se inyecta manualmente desde CLI.
- La CLI muestra que pez o arquetipo aparecio antes del combate.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope

### Incluye

- Extender perfiles de pez con metadata minima de aparicion.
- Introducir un contexto de spawn derivado de `WaterContext`, `InitialDistance`, `InitialDepth` y `HabitatTags` del loadout.
- Resolver automaticamente el perfil de pez antes de crear el encounter.
- Mostrar un resumen de aparicion en la CLI actual.
- Actualizar tests, backlog y README afectados por el nuevo flujo.

### No incluye

- Integrar aun nodos o mapa de expedicion.
- Soportar pesos complejos, rarezas, elites, bosses o tablas probabilisticas por zona.
- Abrir la capa editorial o economica de fotografias.
- Cerrar todavia el modelo final de perfiles data-driven que luego ampliara `BL-007`.

## Propuesta de implementacion

### 1. Metadata minima de aparicion en perfiles

Cada perfil necesita una descripcion minima de cuando puede aparecer.

Primer contrato recomendado:

- `WaterPoolTags`
- `MinInitialDistance` / `MaxInitialDistance`
- `MinInitialDepth` / `MaxInitialDepth`
- `HabitatTags`

Semantica:

- `WaterPoolTags` conecta cada perfil con las aguas base visibles o internas del contexto.
- Distancia y profundidad se leen desde la apertura efectiva del lance, no desde los limites de escape del tablero.
- `HabitatTags` representa acceso o sesgo aportado por aditamentos del loadout.

### 2. Contexto de spawn derivado de la apertura

El resolver no debe mirar el tablero completo ni el track del encuentro. Solo necesita un contexto minimo de aparicion:

- `WaterPoolTag`
- `InitialDistance`
- `InitialDepth`
- `HabitatTags`

Ese contexto puede construirse directamente desde `encounter.Opening` y `loadout.State`.

### 3. Resolver de aparicion inicial

El primer slice puede ser determinista. No necesita aun pesos ni rarezas complejas.

Comportamiento recomendado:

- filtra perfiles compatibles con el contexto de spawn
- calcula una prioridad simple para favorecer perfiles mas especificos sobre perfiles genericos
- elige un resultado estable y reproducible dentro del conjunto compatible
- devuelve tambien informacion suficiente para presentar al jugador que pez o arquetipo aparecio

### 4. Integracion con bootstrap y app

El flujo actual debe cambiar de:

- elegir preset de pez manualmente

a:

- elegir deck del jugador
- elegir `rod`
- elegir aditamentos
- resolver agua y apertura
- resolver aparicion del pez
- construir la baraja del pez desde el perfil seleccionado

Para mantener la direccion arquitectonica, esta seleccion conviene vivir en `app`, no en el adaptador CLI.

### 5. Slice minimo de implementacion

El primer slice deberia cubrir:

- metadata de aparicion en `fishprofiles`
- resolver de spawn con tests
- flujo de app para resolver aparicion y notificarla a la UI
- integracion en `main` para dejar de pedir un preset manual del pez
- resumen de aparicion visible en la CLI

Con eso ya queda la base tactica lista para que `BL-007` profundice perfiles de pez y para que un futuro mapa solo inyecte mejores contextos de agua.

## Validacion

- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Verificar manualmente que distintas combinaciones de agua, apertura y aditamentos cambian el pez o arquetipo resultante.
- Verificar manualmente que la CLI sigue mostrando la apertura y ahora tambien el pez aparecido antes del combate.

## Riesgos o decisiones abiertas

- Habra que decidir despues si el spawn final usa pesos aleatorios o si algunos perfiles son completamente deterministas por contexto.
- Si la metadata de aparicion se vuelve demasiado rica antes de `BL-007`, puede duplicarse trabajo al expandir los perfiles data-driven.
- Queda abierta la decision de cuando introducir roles de encounter como normal, elite, boss de zona o legendario dentro del mismo resolver.
