# Plan de feature: externalizar-catalogo-base-de-peces-y-fish-pools

## Objetivo

Mover el catalogo actual de peces fuera del codigo a un formato data-driven validable e introducir una capa de `fish pools` cerradas que referencien perfiles por `id`. La meta es que el sistema de spawn y el runtime de encounter dejen de depender de `DefaultProfiles()` como lista hardcodeada y puedan consumir catalogos externos sin perder el comportamiento actual ni bloquear la futura integracion con nodos, zonas o rutas de run.

Esta feature debe resolver dos necesidades inmediatas:

- hacer configurable el catalogo base de peces sin recompilar para cada ajuste de contenido
- preparar subconjuntos cerrados y reutilizables de perfiles para que futuros encounters no trabajen sobre todo el catalogo global

No busca todavia externalizar toda la estrategia de aparicion por contexto ni introducir el sistema final de nodos. Busca dejar la base de datos de peces y sus pools en una frontera estable y extensible.

## Criterios de aceptacion

- El catalogo base de peces puede cargarse desde archivos de datos externos en lugar de depender solo de `DefaultProfiles()`.
- Existe un contrato claro y validable para perfiles de pez y para `fish pools` cerradas por `profile_ids`.
- El loader detecta ids duplicados, referencias rotas, valores invalidos y errores estructurales con mensajes legibles.
- El runtime puede resolver un subcatalogo cerrado de peces a partir de una `fish pool` sin duplicar definiciones de perfiles.
- El flujo actual de spawn y combate sigue funcionando aunque todavia exista un fallback temporal al catalogo embebido.
- La implementacion deja una superficie clara para que despues `BL-044` seleccione pools por contexto y `BL-045` conecte el runtime completo al catalogo externo.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope

### Incluye

- Definir el formato de datos del catalogo base de peces.
- Definir el formato de datos de `fish pools` cerradas por `profile_ids`.
- Implementar carga, parseo, validacion y resolucion de referencias entre perfiles y pools.
- Preparar una API de acceso al catalogo para que app/runtime pidan perfiles o pools resueltas sin conocer el formato del archivo.
- Mantener compatibilidad con el spawn actual y con el comportamiento actual de perfiles mientras se migra.
- Cubrir con tests la carga, validacion y resolucion de pools.

### No incluye

- Seleccionar aun que pool usar segun nodo, zona o rol de encounter; eso queda para `BL-044`.
- Reescribir todavia todo el bootstrap del encounter para leer siempre desde archivos; eso puede entrar por slices.
- Externalizar aun metadata editorial rica o sistemas de coleccion; eso queda para `BL-043`.
- Introducir aun plantillas sofisticadas de arquetipos data-driven si no son necesarias para el primer slice; eso queda para `BL-042`.

## Estado actual y limitaciones

Hoy el sistema ya tiene una frontera de dominio razonable:

- `fishprofiles.Profile` expresa la definicion semantica del pez
- `ResolveSpawn(...)` recibe un `[]Profile`, asi que tecnicamente ya puede trabajar sobre subcatalogos cerrados
- `DefaultProfiles()` solo es un origen hardcodeado de ese slice, no una dependencia estructural del algoritmo

Eso es una ventaja importante: no hace falta rehacer el motor de spawn para soportar pools cerradas. Hace falta desacoplar el origen de los perfiles y anadir una capa de catalogo/pools validables.

## Modelo recomendado

La direccion recomendada es separar tres conceptos:

### 1. `FishProfile`

Define que es un pez:

- id
- arquetipo actual
- nombre y descripcion
- metadata de aparicion minima
- mazo/cartas
- politica de reciclado y shuffle

Debe seguir siendo la unidad de dominio que luego consume el spawn y que puede convertirse a `FishDeckPreset`.

### 2. `FishCatalog`

Es la coleccion global de perfiles disponibles.

Responsabilidades:

- cargar perfiles desde archivos
- validar ids unicos
- exponer acceso por id
- resolver listas completas o subsets

### 3. `FishPool`

Es un subconjunto cerrado y nombrado de perfiles, definido por referencias a ids de perfil.

Responsabilidades:

- declarar que peces estan permitidos en un cierto contexto
- no duplicar el contenido del perfil
- poder reutilizarse luego desde encounters, nodos, zonas o tablas de aparicion

## Formato de datos recomendado

### Catalogo de perfiles

Recomendacion: usar JSON en una primera iteracion por simplicidad de tooling y legibilidad.

Estructura sugerida:

```json
{
  "profiles": [
    {
      "id": "classic",
      "archetype_id": "baseline_cycle",
      "name": "Clasico",
      "description": "Baraja base de referencia sin efectos.",
      "details": [
        "Arquetipo: ciclo base sin presion especializada."
      ],
      "appearance": {
        "water_pool_tags": ["shoreline", "mixed_current"],
        "min_initial_distance": 0,
        "max_initial_distance": 2,
        "min_initial_depth": 0,
        "max_initial_depth": 2,
        "required_habitat_tags": []
      },
      "cards": [
        { "move": "blue" },
        { "move": "red" }
      ],
      "cards_to_remove": 3,
      "shuffle": true
    }
  ]
}
```

### Catalogo de pools

Estructura sugerida:

```json
{
  "pools": [
    {
      "id": "shoreline_intro_pool",
      "name": "Shoreline Intro",
      "description": "Pool basica de encuentros costeros iniciales.",
      "profile_ids": ["classic", "surface-control", "deck-exhaustion"]
    }
  ]
}
```

Principio clave:

- las pools solo referencian ids
- nunca duplican definiciones completas de perfiles

## Propuesta de implementacion

### 1. Introducir un paquete o capa de catalogo estable

Hace falta una frontera donde la app/runtime pida perfiles o pools ya resueltas sin conocer como se serializan.

Direccion propuesta:

- mantener `internal/content/fishprofiles/` como owner semantico del dominio pez
- introducir dentro de ese modulo una capa de carga/serializacion o un subpaquete enfocado a catalogos
- separar tipos de archivo (`fishProfileRecord`, `fishPoolRecord`) de tipos de dominio (`Profile`, `FishPool`)

Esto evita contaminar el dominio con detalles de JSON y deja abierta la puerta a otros formatos en el futuro.

### 2. Definir DTOs de archivo y conversores a dominio

No conviene serializar directamente todos los tipos de dominio actuales si eso mezcla demasiado la validacion con el parseo.

Direccion propuesta:

- crear records/DTOs pensados para parsear JSON
- convertirlos luego a `Profile` y `FishPool`
- reutilizar las validaciones ya existentes donde tenga sentido (`Profile.Validate()`, `Appearance.Validate()`, etc.)

Ventajas:

- errores de parseo y errores de dominio quedan mejor separados
- el dominio puede seguir evolucionando sin arrastrar tags de serializacion en todas partes

### 3. Implementar validacion fuerte del catalogo

La validacion es parte central de la feature, no un detalle accesorio.

Debe cubrir al menos:

- ids de perfil no vacios
- ids de perfil unicos
- ids de pool no vacios
- ids de pool unicos
- `profile_ids` no vacios dentro de una pool cuando la pool debe ser utilizable
- referencias de `profile_ids` que existan en el catalogo
- errores estructurales de cartas, moves, arquetipos, habitats y water pools

La salida de errores debe ser accionable, por ejemplo:

- `fish pool shoreline_intro_pool references unknown profile id surface-control-v2`
- `fish profile classic: appearance: ...`

### 4. Resolver pools a slices de dominio reutilizables

El runtime actual ya trabaja con `[]Profile`. Eso hay que aprovecharlo.

Direccion propuesta:

- exponer algo como `ResolvePool(id string) ([]Profile, error)` o equivalente
- mantener la resolucion como una operacion pura sobre catalogos ya cargados
- devolver clones o copias seguras si hace falta evitar mutaciones accidentales del contenido compartido

Esto habilita de forma natural los siguientes usos:

- encounters de prueba con subcatalogos cerrados
- futuros nodos de pesca que apunten a una pool por id
- tablas de aparicion que primero eligen pool y luego delegan al spawn normal

### 5. Mantener compatibilidad con el flujo actual durante la migracion

No hace falta forzar que todo el runtime lea archivos de inmediato.

Direccion propuesta:

- mantener temporalmente `DefaultProfiles()` como fallback o como fuente interna del catalogo por defecto
- anadir una ruta nueva de carga desde archivos sin romper el bootstrap actual
- permitir una adopcion incremental: primero se valida y testea el catalogo externo, luego se conecta el runtime completo

Esto reduce el riesgo y facilita que `BL-045` sea una integracion posterior mas pequena.

### 6. Preparar extensiones futuras sin sobredisenar

BL-041 debe dejar puntos de extension claros, pero no resolverlos todos ahora.

Extensiones previstas:

- `BL-042`: arquetipos o plantillas de cartas de pez reutilizables desde datos
- `BL-043`: metadata editorial/semantica mas rica de especies
- `BL-044`: tablas data-driven que eligen que pool usar segun contexto
- `BL-045`: bootstrap y runtime de spawn consumiendo catalogos externos como camino principal

La feature debe dejar esos caminos faciles, pero sin introducir un framework de contenido sobredimensionado.

## API y uso futuro recomendados

La capa resultante deberia permitir usos como estos:

### Uso actual o de pruebas

- cargar catalogo por defecto
- resolver una pool concreta para un encounter de prueba
- pasar ese subset a `ResolveFishSpawn(...)`

### Uso futuro desde nodos

- un nodo o tabla de aparicion decide `fish_pool_id`
- el catalogo resuelve la pool a `[]Profile`
- el spawn actual opera sobre ese subset, no sobre todo el catalogo

### Uso futuro desde tooling

- validar catalogos offline
- listar perfiles y pools
- detectar referencias rotas antes de ejecutar el juego

## Slice minimo recomendado

El primer slice de implementacion deberia cubrir:

- tipos de archivo para perfiles y pools
- parseo JSON y validacion basica
- resolucion de pools por `profile_ids`
- tests de carga, ids duplicados y referencias invalidas
- una ruta de uso minima desde app o tests que demuestre que un encounter puede usar una pool cerrada

Eso ya aporta valor inmediato sin esperar al sistema de nodos ni a una externalizacion mas rica de contenido.

## Validacion

- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Verificar que un subcatalogo resuelto desde una pool produce el mismo comportamiento de spawn que el slice equivalente construido a mano.
- Verificar que errores de ids duplicados y referencias rotas fallen con mensajes claros.

## Riesgos o decisiones abiertas

- Si el schema se acopla demasiado a los tipos actuales, futuras extensiones de contenido pueden ser mas rigidas de lo deseable.
- Si se intenta resolver arquetipos, metadata editorial y tablas de aparicion en la misma feature, el scope volvera a crecer demasiado.
- Habra que decidir si el runtime carga archivos siempre en arranque o si primero se usa solo en tests/herramientas hasta que `BL-045` lo convierta en camino principal.
- Conviene evitar que las pools se conviertan en un segundo lugar donde se duplique logica de aparicion; deben declarar subsets cerrados, no reglas completas de spawn.
