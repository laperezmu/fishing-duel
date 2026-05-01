# Plan de feature: resolver-spawns-desde-catalogos-y-fish-pools

## Objetivo

Conectar el flujo actual `agua -> apertura -> spawn -> mazo del pez` a los catalogos data-driven ya existentes, de modo que el runtime deje de depender de `DefaultProfiles()` como catalogo global implicito y pueda arrancar encounters usando una `fish pool` cerrada concreta. La meta de esta feature no es decidir aun que pool corresponde a cada nodo o zona, sino dejar preparado el bootstrap y la capa de app para consumir subsets cerrados de peces desde ids o desde catalogos cargados.

Esta feature debe ser el puente entre la infraestructura ya creada en `BL-041` y el futuro sistema de aparicion por contexto de `BL-044`, permitiendo que el juego use pools cerradas reales antes de tener mapa o nodos completos.

## Criterios de aceptacion

- El bootstrap de encounter puede recibir un catalogo o una `fish pool` concreta en vez de depender siempre de `DefaultProfiles()`.
- La app puede resolver un spawn usando perfiles provenientes de `Catalog.ResolvePool(...)` sin adaptaciones manuales fuera del dominio.
- Sigue existiendo un camino por defecto compatible con el flujo actual cuando no se proporciona una pool concreta.
- El contrato nuevo deja claro donde vive la seleccion del subset de peces y donde vive la resolucion del spawn dentro de ese subset.
- La CLI actual sigue funcionando.
- Los tests cubren al menos un encounter inicializado contra una pool cerrada y el fallback por defecto.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope

### Incluye

- Introducir una forma explicita de pasar un `fishprofiles.Catalog`, una `fish_pool_id` o un subset de perfiles al bootstrap del encounter.
- Ajustar `internal/app/` para que la resolucion del spawn consuma pools cerradas cuando se le indiquen.
- Mantener compatibilidad con el flujo interactivo actual y con el bootstrap actual por defecto.
- Anadir tests de integracion en `app` para el nuevo wiring.
- Documentar la nueva frontera de uso si hace falta fijar el contrato.

### No incluye

- Elegir aun pools desde tablas por agua, zona o rol de encounter; eso sigue en `BL-044`.
- Cargar aun pools desde nodos de run o mapas; eso quedara para trabajo posterior del core loop.
- Cambiar el algoritmo tactico del spawn ni su scoring/variedad mas alla de consumir subsets cerrados.
- Rehacer la UX del setup de CLI.

## Estado actual y hueco a cerrar

Hoy ya tenemos:

- `fishprofiles.Catalog`
- `fishprofiles.ResolvePool(...)`
- `ResolveFishSpawn(...)` trabajando sobre `[]Profile`
- `DefaultProfiles()` leyendo desde el catalogo embebido

Pero el bootstrap actual sigue cableando esto de forma global:

- `ResolveFishSpawnWithRandomizer(..., fishprofiles.DefaultProfiles(), ...)`

Eso significa que el runtime sigue sabiendo como llegar al catalogo global, en vez de recibir un subset cerrado como dependencia explicita.

## Propuesta de implementacion

### 1. Introducir un contrato explicito de fuente de peces para el encounter

El bootstrap necesita dejar de asumir que siempre usara todo el catalogo por defecto.

Direccion propuesta:

- crear un pequeno contrato de entrada para el bootstrap, por ejemplo una `EncounterFishSource` o una configuracion equivalente
- permitir variantes como:
  - catalogo por defecto + pool por defecto
  - catalogo explicito + `fish_pool_id`
  - subset de perfiles ya resuelto

La meta es que la capa de app reciba una intencion explicita sobre el subset de peces y no lo derive internamente de un unico hardcode.

### 2. Mantener un camino por defecto para el flujo actual

No conviene romper el flujo actual de CLI mientras todavia no existe un nodo real de run.

Direccion propuesta:

- mantener `BootstrapEncounter(...)` como entrypoint simple para el juego actual
- introducir una variante mas explicita, algo tipo `BootstrapEncounterWithCatalog(...)` o una configuracion estructurada equivalente
- dejar el camino viejo delegando al nuevo con el catalogo por defecto y una pool por defecto razonable

Esto evita churn innecesario en `cmd/` y `cli/`.

### 3. Resolver la pool en la capa de app, no en el adaptador de UI

La seleccion del subset de peces debe vivir fuera del borde de interfaz.

Direccion propuesta:

- hacer que `internal/app/` resuelva `fish_pool_id -> []Profile`
- mantener `ResolveFishSpawn(...)` agnostico del origen de ese slice
- evitar que CLI o futuros nodos tengan que conocer detalles internos de `Catalog`

Asi la frontera queda limpia:

- una capa selecciona la pool
- otra resuelve el spawn dentro de esa pool

### 4. Exponer un punto de extension claro para nodos futuros

Aunque todavia no existan nodos data-driven, esta feature debe dejar el hook listo.

Direccion propuesta:

- definir una estructura de configuracion del encounter que pueda crecer luego con campos como:
  - `FishPoolID`
  - `WaterContextID`
  - `Seed`
  - overrides futuros de spawn o encounter
- no hace falta llenar todos esos campos ahora; basta con fijar la direccion correcta

### 5. Cubrir con tests el wiring nuevo

El valor de esta feature esta en el wiring, asi que los tests tienen que demostrarlo.

Direccion propuesta:

- test de bootstrap o de resolucion donde una pool cerrada limite realmente el subset disponible
- test del camino por defecto para garantizar compatibilidad
- test de error legible si la pool solicitada no existe

## Opciones de API razonables

### Opcion A: funcion nueva de bootstrap con config

Ventajas:

- mas escalable para nodos y run
- evita proliferar parametros posicionales

Ejemplo conceptual:

```go
type EncounterBootstrapConfig struct {
    Catalog    fishprofiles.Catalog
    FishPoolID string
}
```

### Opcion B: helper puntual para resolver subset antes del bootstrap

Ventajas:

- cambio pequeno
- facil de introducir rapido

Desventaja:

- menos expresivo para el futuro runtime de run

Recomendacion:

- prefiero la opcion A, aunque sea minima, porque encaja mejor con futuros nodos de pesca

## Slice minimo recomendado

El primer slice deberia cubrir:

- una configuracion minima de bootstrap con catalogo y `fish_pool_id`
- fallback al catalogo/pool por defecto
- resolucion de perfiles desde `Catalog.ResolvePool(...)`
- test de uso real con una pool cerrada

Con eso el runtime ya deja de estar pegado al catalogo global y queda listo para que `BL-044` solo decida que pool usar, en lugar de tener que cambiar otra vez el wiring del encounter.

## Extensiones futuras habilitadas

Esta feature deja preparado el camino para:

- nodos de pesca que referencien `fish_pool_id`
- encuentros elites o bosses con pools dedicadas
- tablas de aparicion por agua/contexto que primero elijan pool y luego deleguen al spawn normal
- tooling o debug scenarios que ejecuten encuentros concretos contra subsets cerrados

## Validacion

- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Verificar manualmente que el flujo actual sigue funcionando con el catalogo/pool por defecto.
- Verificar que una pool cerrada cambia realmente el subset elegible del spawn.

## Riesgos o decisiones abiertas

- Si el contrato de bootstrap se vuelve demasiado ambicioso ahora, puede mezclarse con decisiones futuras de nodos o run.
- Si la seleccion de pool se deja demasiado difusa, `BL-044` volvera a empujar logica al wiring de app en vez de reutilizar esta frontera.
- Habra que decidir cual es la pool por defecto del juego actual mientras no existan nodos reales; probablemente una pool global explicita y no una llamada directa al catalogo completo.
