# Plan de feature: apertura-encounter-pesca-y-cast

## Objetivo

Introducir una fase breve y autocontenida de apertura del encounter de pesca donde el jugador lea una situacion de agua, ejecute un minijuego de cast por timing y abra una ventana horizontal de inicio antes de entrar al duelo normal contra el pez.

La meta de esta feature no es resolver aun el sistema completo de nodos, pools de peces o economia fotografica. La meta es fijar el contrato base de la apertura tactica del encuentro y validarlo en gameplay real con el sistema actual, de forma que el cast deje de ser una idea abstracta y pase a ser una capa concreta del loop de pesca.

El resultado debe dejar una base que luego pueda conectarse con `BL-002`, `BL-021` y `BL-022`, pero sin quedar bloqueado por ellos. La primera iteracion debe poder vivir dentro del encounter actual, con contexto de agua inyectado desde presets o configuracion temporal y con impacto real sobre la distancia inicial del duelo.

Aunque el scope activo de esta feature solo modifica `InitialDistance`, el contrato de apertura debe dejar preparado al encounter para recibir tanto `InitialDistance` como `InitialDepth` como valores de entrada. La profundidad inicial seguira viniendo del valor base actual hasta que la cana y sus aditamentos entren en juego en `BL-021`, pero el encounter no debe quedar cerrado a ese futuro.

## Criterios de aceptacion

- Existe una fase de apertura del encounter antes del primer round del combate.
- La apertura recibe un contexto de agua autocontenido, sin depender aun del sistema de nodos o mapa.
- El jugador puede resolver un cast por timing que termina en una banda horizontal discreta, no en una coordenada libre continua.
- El contrato de apertura puede producir valores iniciales para `InitialDistance` y `InitialDepth`, aunque la primera iteracion solo altere `InitialDistance`.
- El resultado del cast impacta el estado inicial del encounter actual de una forma legible y testeable, al menos sobre `InitialDistance`.
- La CLI comunica con claridad la lectura del agua, el minijuego de cast y el resultado obtenido.
- La primera version valida que un cast largo no es automaticamente mejor en todos los contextos.
- La implementacion deja un contrato reutilizable para que futuras features conecten el cast con aditamentos de cana, subpools de peces y nodos desconocidos.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope

### Incluye

- Definir un modelo minimo de `contexto de agua` reusable por el encounter.
- Definir una fase de apertura que ocurre antes del combate normal.
- Implementar un minijuego de cast por timing en CLI con resultado en bandas discretas.
- Traducir el resultado del cast a una apertura horizontal concreta del encounter actual.
- Dejar preparado el contrato para que la apertura del encounter pueda recibir tambien `InitialDepth`, aunque el slice actual no lo modifique todavia.
- Crear al menos algunos contextos de agua o presets de prueba para validar comportamientos distintos.
- Ajustar presentacion y documentacion para que el nuevo flujo sea entendible.

### No incluye

- Implementar aun el sistema de nodos, rutas o mapa de zona.
- Seleccionar subpools reales de peces por cast o por profundidad; eso queda para `BL-022`.
- Redefinir aun la cana y sus aditamentos como sistema completo; eso queda para `BL-021`.
- Introducir economia fotografica o tags editoriales ligados al cast.
- Reemplazar por completo el flujo actual de seleccion de presets por un sistema final de expedicion.

## Propuesta de implementacion

### 1. Modelo minimo de apertura del encounter

La feature necesita un contrato pequeno y estable para no mezclar la idea del cast con el futuro sistema de nodos.

Propuesta base:

- Introducir un modelo equivalente a:
  - `WaterContext`
  - `CastBand`
  - `CastResult`
  - `EncounterOpening`
- `WaterContext` describe solo lo necesario para esta iteracion:
  - id estable
  - nombre corto legible
  - descripcion o señales visibles de lectura para el jugador
  - metadata interna del tipo de pool o del habitat que ese contexto activa
  - bandas horizontales validas o destacadas
  - modificadores base sobre distancia inicial
  - profundidad inicial base o valor derivado por defecto
- `CastResult` describe:
  - banda conseguida
  - calidad del cast si hace falta
  - modificador de distancia inicial derivado
- `EncounterOpening` empaqueta el contexto de agua y el resultado del cast para producir la configuracion inicial del encounter.

Punto importante:

- las señales visibles del agua y la metadata interna de pool no son lo mismo;
- la primera version puede almacenar ambas en el contexto, pero solo la capa visible debe llegar a la interfaz;
- la metadata de pool queda preparada para futuro y no forma parte de la lectura mostrada al jugador en este item.

La idea no es modelar aun habitats completos, sino fijar un contrato de apertura que despues pueda alimentar subpools reales.

### 2. Resolver primero el efecto sobre el encounter actual

Para reducir scope, la primera integracion no debe intentar resolver peces distintos ni nuevos tipos de encounter. Debe impactar el sistema actual que ya funciona.

Primera traduccion recomendada:

- El cast modifica `InitialDistance` antes de crear el estado del encounter.
- La distancia base del encounter deja de ser siempre fija y pasa a construirse como:
  - `distancia base del contexto de agua`
  - `+ resultado horizontal del cast`
- La profundidad inicial no cambia todavia por el cast, pero el encounter debe poder recibirla desde la apertura como valor explicito.
- La primera version puede dejar `InitialDepth` igual al valor base del contexto o al valor actual por defecto, hasta que `BL-021` defina el sesgo vertical de la cana.

Esto permite validar si el cast cambia de verdad la sensacion de apertura sin meter todavia el resto del sistema futuro.

### 3. Bandas horizontales discretas

El minijuego puede ser de precision, pero el sistema debe leerlo en bandas discretas para que siga siendo balanceable y legible.

Propuesta inicial:

- `muy corto`
- `corto`
- `medio`
- `largo`
- `muy largo`

No todos los contextos de agua tienen que valorar igual cada banda. La validacion minima debe cubrir al menos dos contextos donde la banda mas larga no sea automaticamente la mejor.

Ejemplos de contexto:

- `ensenada cercana`: favorece corto o medio
- `canal abierto`: favorece largo o muy largo
- `corriente irregular`: hace viable medio pero vuelve arriesgado el maximo

### 4. Minijuego CLI de cast

La primera version debe ser simple, repetible y suficientemente expresiva para probar el timing.

Propuesta:

- Mostrar una barra de carga que sube y baja en poco tiempo, en ciclo constante.
- El jugador hace un solo input para detener la barra en una seccion concreta.
- La seccion en la que termina la barra se traduce a una banda horizontal.
- La UI comunica con claridad:
  - contexto de agua actual
  - significado general de las secciones o bandas
  - resultado final del cast

No hace falta simular fisica de lance ni usar varios inputs. Un solo gesto de timing sobre una barra oscilante basta para validar la idea.

### 5. Inyeccion de contexto sin sistema de nodos

Para desacoplar la feature de `BL-002`, el contexto de agua debe entrar por una via simple y temporal.

Opciones validas para esta iteracion:

- preset fijo desde bootstrap
- seleccion previa en CLI junto a los presets de pez o jugador
- configuracion temporal dentro de `cmd/fishing-duel/`

Recomendacion por defecto:

- anadir una seleccion sencilla de `situacion de agua` en CLI antes de empezar el encounter.

Eso permite validar varios contextos sin tener que construir todavia un mapa o un sistema de nodos ocultos.

### 6. Integracion con arquitectura actual

La feature deberia respetar la separacion entre runtime, contenido y presentacion.

Direccion recomendada:

- `internal/content/`
  - presets o catalogos de contextos de agua
- `internal/encounter/` o un paquete hermano de runtime del encounter
  - tipos de apertura y aplicacion a la configuracion inicial
  - contrato para recibir `InitialDistance` y `InitialDepth` resueltos desde la apertura
- `internal/app/`
  - coordinacion del flujo `leer agua -> cast -> crear encounter`
- `internal/presentation/` y `internal/cli/`
  - lectura del contexto, render del minijuego y comunicacion del resultado

La regla principal es no meter toda la logica del cast dentro de `cmd/` ni convertir la CLI en el source of truth del sistema.

### 7. Slice minimo de implementacion

Para mantener la feature cerrable, conviene definir un primer slice pequeño:

- 3 contextos de agua de prueba
- 5 bandas horizontales discretas
- minijuego de barra oscilante en CLI con un solo input
- impacto visible sobre `InitialDistance`
- contrato listo para pasar tambien `InitialDepth`, aunque siga usando el valor base
- mensaje claro de apertura antes del primer round

Con eso ya se puede contestar si el cast aporta agencia real sin haber tocado aun pools de peces, nodos o cana.

### 8. Cobertura automatizada

La feature deberia cubrir al menos estos niveles:

- `runtime de apertura`
  - traduccion correcta de banda a modificador horizontal
  - aplicacion correcta del contexto de agua a la configuracion del encounter
  - generacion de una apertura que pueda transportar `InitialDistance` y `InitialDepth`
- `encounter`
  - construccion del estado inicial con distancia modificada por la apertura
  - aceptacion correcta de valores iniciales explicitos para distancia y profundidad
- `presentation`
  - textos o etiquetas legibles para contexto y resultado del cast
  - ausencia de fuga de metadata interna de pool a la UI
- `cli`
  - flujo de apertura y render del resultado sin romper el encounter actual

## Validacion

- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Validar manualmente desde CLI al menos dos contextos de agua donde la banda optima cambie.
- Validar manualmente que el encounter arranca con distinta distancia inicial segun el resultado del cast.
- Validar manualmente que la profundidad inicial sigue llegando como valor del encounter aunque el cast no la modifique todavia.
- Verificar que el juego sigue siendo jugable aunque el cast salga mal y que el fallo no vacia el encounter.

## Riesgos o decisiones abiertas

- Si el minijuego ocupa demasiado tiempo, puede frenar el ritmo de pesca en vez de enriquecerlo.
- Si todas las bandas se sienten iguales, el cast sera cosmetico y no tactico.
- Si el cast altera demasiado la distancia inicial, puede volverse mas determinante de lo deseado.
- Habra que decidir cuanto ruido o informacion exacta se muestra en la lectura del agua antes del input y cuanto queda solo como metadata interna de pool.
- La integracion futura con nodos, cana y subpools debe reutilizar este contrato en vez de duplicarlo desde otro sistema paralelo.
