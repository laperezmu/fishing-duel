# Plan de feature: definir-contrato-base-y-mvp-secuencial-de-una-run

## Objetivo

Fijar el contrato base de una run roguelike de pesca como una expedicion jugable en una sola sesion, separando con claridad tres capas que ya no conviene mezclar: el runtime tactico de cada `encounter`, el estado persistente dentro de la `run` actual y la futura capa meta entre runs. La meta de esta feature no es implementar todavia el mapa, la economia completa ni el guardado, sino cerrar la frontera conceptual y tecnica minima que permita construir el primer vertical slice de expedicion sin rehacer despues los contratos principales.

Esta feature debe convertir la direccion ya acordada del backlog en un marco operativo: que vive en `RunState`, que entra y sale de un encounter, que recursos persisten entre nodos y zonas, cuales son los fail states del MVP y en que orden deben atacarse las subtareas siguientes hasta llegar a una primera run jugable.

## Criterios de aceptacion

- Queda definida la frontera entre estado `por encuentro`, estado `por run` y estado `entre runs`.
- Queda definido el handoff minimo `encounter -> run`, incluyendo que resultados del duelo impactan recursos globales de expedicion.
- Queda definido el contrato conceptual de `RunState` del MVP, aunque su implementacion llegue despues.
- Queda fijado el recorte del primer MVP de run: una sola sesion, sin guardado, sin meta-progresion y con flujo simple de nodos.
- Quedan explicitados los fail states, condiciones de retiro y condiciones de cierre del MVP.
- Queda decidido el orden de trabajo recomendado para `BL-035`, `BL-036`, `BL-037`, `BL-038`, `BL-039` y `BL-040`.
- Queda claro que partes se difieren para backlog posterior: mapa rico, coleccion, guardado, meta y sistemas expandidos.

## Scope

### Incluye

- Definir el proposito del sistema de run dentro de la vision roguelike actual.
- Delimitar responsabilidades entre `encounter`, `run` y futura `meta`.
- Identificar que datos deben vivir durante toda la expedicion del MVP.
- Definir el flujo macro minimo de una run y sus puntos de integracion con encounters y servicios.
- Bajar el alcance a un primer vertical slice implementable en pasos pequenos.
- Dejar accionables concretos para los items dependientes del backlog.

### No incluye

- Implementar todavia el runtime completo de run.
- Resolver el detalle del mapa por zonas posterior al MVP base.
- Disenar o implementar guardado, perfil permanente o meta-progresion.
- Cerrar todavia la economia ampliada, coleccion, bestiario o recompensas entre runs.
- Reabrir el slice tecnico ya cubierto por `BL-034`.

## Proposito de la tarea

El proposito real de `BL-001` no es producir una abstraccion elegante en el vacio, sino evitar que el primer MVP de expedicion nazca con responsabilidades mezcladas. Hoy el proyecto ya tiene una base de `encounter` bastante saneada: opening, cast, spawn, bootstrap y runtime tactico ya quedaron mas modulares y mas UI-agnostic. El cuello de botella ya no esta en el duelo aislado, sino en la falta de una capa intermedia que diga:

- cuando arranca y termina una expedicion
- que persiste entre encuentros
- como un resultado tactico impacta el progreso global
- donde viven las decisiones de nodo, servicio, retiro y cierre de zona

Sin esa frontera, cualquier avance sobre economia, nodos, servicios o patrocinadores corre riesgo de pegarse al runtime tactico actual o de empujar estado de run dentro de `match.State`, que justamente ya se venia limpiando para evitar eso.

## Analisis del alcance

### 1. Que problema debe resolver `BL-001`

El juego ya sabe resolver un duelo de pesca, pero todavia no sabe resolver una expedicion completa. El vacio no esta en una mecanica puntual, sino en el contrato que organiza la sesion completa.

El problema a resolver es doble:

- de producto: convertir encounters sueltos en una run con principio, progreso, tension y cierre
- de arquitectura: evitar que run, encounter y futura meta queden acoplados entre si

### 2. Que decisiones necesita cerrar ahora

Para habilitar el MVP, `BL-001` debe dejar fijadas al menos estas decisiones:

- la run es una sola sesion sin persistencia entre partidas
- el hilo de la cana es la condicion principal de derrota global
- el tiempo de zona es presion local y no derrota global automatica
- el retiro voluntario existe, pero solo en puntos seguros
- la build, la economia y el desgaste persisten durante toda la expedicion
- el primer MVP usa un flujo secuencial de nodos antes de tener mapa rico
- patrocinadores no deben bloquear el primer vertical slice

Estas decisiones ya aparecen insinuadas en backlog, pero `BL-001` debe consolidarlas como base de implementacion y no como notas sueltas.

### 3. Fronteras de estado que conviene fijar

#### Estado por encounter

Debe seguir siendo tactico, efimero y autocontenido. Incluye por ejemplo:

- apertura del encounter
- spawn del pez
- estado del duelo
- cartas vistas o jugadas
- outcome del encuentro
- dano puntual producido durante ese duelo

Este estado no deberia convertirse en almacenamiento natural de dinero, progreso de zona, historial de nodos o patrocinadores.

#### Estado por run

Debe representar todo lo que vive desde que empieza la expedicion hasta que termina. Para el MVP, el contrato minimo deberia contemplar:

- identidad o etapa actual de la expedicion
- zona actual y progreso dentro de la zona
- nodo actual o siguiente opcion navegable
- hilo restante y su maximo o capacidad operativa
- build activa del jugador
- dinero disponible
- fotos capturadas pendientes de liquidacion o reservadas, segun el recorte final de `BL-037`
- flags de retiro, derrota o acceso al boss final
- eventuales modificadores globales de run que luego podran alimentar patrocinadores o servicios

#### Estado entre runs

Debe quedar explicitamente fuera del MVP. Incluye:

- perfil permanente
- bestiario
- trofeos
- progreso de coleccion
- desbloqueos meta
- guardado/carga de carrera

La clave no es implementarlo, sino evitar dejar dependencias tempranas hacia esta capa.

### 4. Handoff `encounter -> run`

El punto mas importante del contrato es definir que sale del duelo y como se traduce a progreso de expedicion. El handoff minimo del MVP deberia contemplar un `EncounterResult` o snapshot equivalente con informacion como:

- resultado del encuentro: captura, escape, derrota, retiro forzado u otros outcomes relevantes
- dano neto al hilo o cambios permanentes de recursos globales
- pez capturado o metadata suficiente para convertirlo luego en valor fotografico/economico
- efectos persistentes ganados o perdidos al terminar el duelo
- si el nodo queda resuelto, fallado o reintenable segun la regla del MVP

La run no deberia inspeccionar `match.State` completo para decidir esto. El outcome del duelo debe salir resumido en un contrato pequeno y semantico.

### 5. Flujo macro minimo del MVP

Antes de entrar en mapas ricos, el MVP necesita una forma muy simple de expedicion completa. La direccion mas consistente hoy es algo cercano a:

`inicio de run -> nodo de pesca -> nodo de pesca -> servicio/checkpoint -> encounter final o boss -> cierre`

Ese flujo alcanza para validar:

- encadenado de multiples encounters dentro de una misma sesion
- persistencia de recursos globales entre encuentros
- gasto o reparacion en un servicio
- tension de llegar vivo al final de la expedicion
- cierre por victoria, derrota o retiro

No hace falta aun:

- ramas complejas
- ocultacion de informacion
- mapa procedural
- multiples tipos exoticos de nodo

### 6. Riesgo principal de alcance

El mayor riesgo de `BL-001` es intentar resolver demasiado temprano el juego completo y no el contrato base. Las alertas claras de sobrealcance son:

- disenar ya el mapa completo post-MVP
- resolver ya la meta-progresion
- meter patrocinadores completos antes de validar el loop base
- convertir `RunState` en una bolsa gigante de detalles futuros
- acoplar nodos, economia y servicios a structs tacticos del encounter

El recorte sano es definir solo lo necesario para que `BL-035` a `BL-040` puedan avanzar sin ambiguedad fuerte.

## Propuesta de contrato base

### 1. Entidades conceptuales minimas

Direccion propuesta:

- `RunState`: fotografia viva de la expedicion actual
- `RunNode`: unidad minima de avance del flujo de run
- `RunOutcome`: estado de cierre de la expedicion
- `EncounterResult`: resultado reducido que la run consume al salir de un duelo
- `RunAction` o comando equivalente: avanzar, entrar a encounter, usar servicio, retirarse

No hace falta fijar aun nombres definitivos de paquetes o structs, pero si su responsabilidad.

### 2. Contrato conceptual minimo de `RunState`

Direccion propuesta:

- `status`: en progreso, retirada, derrota, victoria
- `currentZone`: referencia a la zona actual del MVP
- `progress`: posicion actual dentro del flujo secuencial de nodos
- `thread`: estado global del hilo de la expedicion
- `loadout`: build persistente durante la run
- `currency`: dinero o recurso equivalente de gasto
- `captures`: resultados capturados pendientes de conversion o liquidacion
- `modifiers`: espacio pequeno para efectos globales de run si hacen falta

La idea no es fijar todos los campos finales, sino impedir que recursos de run terminen desperdigados entre `app`, `match`, `player` o adaptadores.

### 3. Estados terminales del MVP

Para el MVP alcanzan tres cierres principales:

- `victory`: se completa la expedicion y se captura o resuelve el objetivo final
- `defeat`: el jugador agota el hilo global y no puede continuar
- `retired`: el jugador decide cerrar la expedicion en un punto permitido

Esto deja fuera por ahora cierres meta, scoreboards persistentes o progresion entre runs.

## Orden recomendado de trabajo despues de `BL-001`

### `BL-035`

Definir el flujo minimo navegable de zona y nodos. Debe venir primero porque fija el esqueleto sobre el que se enchufan encounter, servicio y cierre.

### `BL-036`

Definir recursos globales y reglas del hilo. Debe venir enseguida porque la tension principal de la run vive ahi.

### `BL-037`

Definir economia minima. Conviene cerrarla despues del flujo y los recursos, porque depende de saber cuando se captura, cuando se liquida y donde se gasta.

### `BL-038`

Definir servicios y build minima. Necesita ya saber que recursos existen y en que nodos se consumen.

### `BL-039`

Tomar la decision final sobre patrocinadores MVP. Conviene dejarlo despues del loop base para que no bloquee el primer slice si termina siendo optativo.

### `BL-040`

Implementar el primer vertical slice solo cuando los contratos anteriores ya no esten moviendose fuerte.

## Accionables concretos

### Para cerrar `BL-001`

- Crear un documento de referencia del contrato de run MVP con vocabulario estable.
- Decidir el shape minimo del handoff `encounter -> run`.
- Decidir el shape minimo conceptual de `RunState`.
- Confirmar el flujo macro del MVP y sus estados terminales.
- Confirmar que queda explicitamente fuera de scope hasta despues del primer slice.

### Para abrir implementacion despues

- Abrir `BL-035` como plan de feature con el flujo secuencial de nodos.
- Abrir `BL-036` con contrato de recursos y reglas del hilo.
- Abrir `BL-037` con la economia minima del loop.
- Abrir `BL-038` con servicios y mejoras minimas.
- Resolver `BL-039` como decision corta de producto antes de `BL-040`.

## Riesgos o decisiones abiertas

- Hay que evitar que el concepto de zona entre demasiado detallado antes de cerrar el flujo secuencial base.
- Hay que decidir en `BL-037` si las fotos reservadas ya existen en el MVP o si se recortan a una sola capa de captura/liquidacion.
- Hay que decidir en `BL-036` si el hilo global se modela como recurso unico o como combinacion de capacidad maxima mas desgaste acumulado.
- Hay que decidir en `BL-035` cuanto control tiene el jugador sobre el siguiente nodo antes de introducir mapa real.
- Aunque `BL-034` ya este cerrado tecnicamente, conviene reflejarlo tambien en backlog para que no siga apareciendo como trabajo pendiente.
