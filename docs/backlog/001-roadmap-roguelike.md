# Roadmap roguelike

Este documento concentra el backlog activo del proyecto, con estado visible para distinguir que esta pendiente, que ya esta planificado y que ya quedo integrado en `main`.

## Convencion de estado

- `pending`: item identificado pero aun sin plan activo en `docs/features/`.
- `planned`: item ya convertido en plan de feature, pero todavia no mergeado en `main`.
- `done`: item ya integrado en `main`.
- `cancelled`: item descartado, absorbido por otro o fuera de foco.

## Foto actual

- `done`: `BL-005`, `BL-006`, `BL-018`, `BL-019`, `BL-020`, `BL-021`, `BL-022`, `BL-029`, `BL-030`, `BL-033`, `BL-041`, `BL-050`
- `planned`: `BL-034`, `BL-045`
- `pending`: resto del roadmap
- Foco recomendado inmediato: cerrar primero la base MVP de una run sin persistencia meta y conectar el runtime actual a los nuevos catalogos cerrados, empezando por `BL-001`, `BL-034`, `BL-035` y `BL-045`.

## Foco sugerido actual

- `BL-001`: fijar el contrato base y el MVP secuencial de una run sin persistencia meta.
- `BL-034`: desacoplar setup, opening y bootstrap antes del primer delivery real de run.
- `BL-035`: definir el flujo minimo navegable de zona y nodos para el MVP de expedicion.
- `BL-045`: conectar el spawn actual a fish pools cerradas para que el runtime deje de depender del catalogo global.
- `BL-050`: ya integrado; usarlo como base para que futuras fish pools no hereden un spawn totalmente deducible.

## Core Loop

### BL-001 Definir contrato base y MVP secuencial de una run
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: fijar la arquitectura minima de una expedicion jugable en una sola sesion, separando claramente encounter, run y futura meta, para que la implementacion pueda avanzar por etapas UI-agnostic y sin depender todavia de persistencia entre runs.
- **Resultado esperado**: contrato marco del MVP de run con estados base, handoff `encounter -> run`, recursos globales de expedicion, flujo minimo de zona y lista secuencial de subtareas necesarias para llegar a una primera run jugable sin meta-persistencia.
- **Direccion actual acordada**:
  - La run es una expedicion de pesca deportiva que empieza en aguas mas accesibles y avanza hacia mar abierto y peces de leyenda.
  - El objetivo de la run es llegar vivo y suficientemente fuerte al final de la expedicion para capturar un pez legendario final.
  - La derrota de la run ocurre cuando el jugador agota todo el hilo de la cana.
  - El retiro voluntario puede ocurrir solo en checkpoints o servicios y cierra la expedicion sin acceder al legendario final.
  - Cada zona tiene tiempo limitado, pero el tiempo funciona como presion local de zona y no como fail state global de la expedicion.
  - Entre zonas persisten la build, la economia, los patrocinadores, el desgaste del hilo y cualquier otro recurso de run.
  - Las zonas aumentan la dificultad y la cercania al final de la expedicion, pero no resetean la progresion acumulada en la run.
  - El MVP actual no incluye persistencia entre runs ni meta-progresion; esas capas deben quedar separadas, pero fuera de scope de la primera entrega.
- **Capas de estado a distinguir desde el analisis**:
  - `por encuentro`: estado tactico del duelo, cartas vistas, outcome y dano puntual del enfrentamiento.
  - `por run`: hilo restante, build del jugador, dinero, patrocinadores activos, fotos reservadas y progreso interno de la expedicion actual.
  - `entre runs`: reservado para backlog futuro, pero fuera de este MVP.
- **Recorte actual del MVP**:
  - run de una sola sesion, sin guardado ni profile permanente.
  - 1 o 2 zonas fijas como maximo.
  - flujo de nodos simple y navegable antes de abordar mapa rico o ocultacion compleja.
  - economia minima `captura -> fotos -> dinero -> servicio -> mejora`.
  - hilo como condicion global de derrota.
  - patrocinadores en forma minima y tardia, solo si no bloquean el primer vertical slice jugable.
- **Subtareas secuenciales necesarias para implementacion**:
  - `BL-034`: slice minimo de arquitectura UI-agnostic para setup, opening y bootstrap antes de crecer la run.
  - `BL-035`: definir flujo minimo navegable de zona y nodos del MVP.
  - `BL-036`: definir contrato de recursos globales de run y dano al hilo.
  - `BL-037`: definir economia minima de fotos, dinero, liquidacion y gasto en servicio.
  - `BL-038`: definir build minima y acciones de servicio del MVP.
  - `BL-039`: definir patrocinadores MVP o decidir explicitamente postergarlos para el primer vertical slice.
  - `BL-040`: implementar primer vertical slice jugable de run MVP sobre esas bases.
- **Criterios de cierre**:
  - queda definido el contrato de `RunState` y el handoff `encounter -> run`
  - queda fijado el MVP de run sin persistencia meta ni guardado
  - queda decidido el orden secuencial de implementacion hasta una primera run jugable
  - queda claro que partes se difieren a backlog posterior: meta, coleccion, guardado y mapa rico
- **Prioridad**: Alta

### BL-034 Desacoplar setup, opening y bootstrap para habilitar la run MVP
- **Estado**: `planned`
- **Tipo**: Discovery + Delivery
- **Objetivo**: ejecutar el slice minimo de arquitectura UI-agnostic necesario para que la run no nazca acoplada al CLI, dejando fuera del borde de terminal la orquestacion de setup, opening y entrada al encounter.
- **Resultado esperado**: `cmd/` reducido a composicion, opening reusable fuera de `internal/cli/`, y contratos de app/presentation suficientemente estables para que una run orqueste encounters sin depender de tipos renderizados por el terminal.
- **Dependencias**: `BL-001`, `BL-018`, `BL-020`, `BL-021`
- **Plan relacionado**: `docs/features/022-desacoplar-setup-opening-y-bootstrap-para-run-mvp.md`
- **Direccion actual acordada**:
  - No hace falta completar `BL-023` entero antes de la run.
  - Si conviene ejecutar primero el slice que saque del CLI el setup, el cast y el bootstrap mas acoplado.
  - Este item es prerrequisito tecnico del primer delivery real de run.
- **Prioridad**: Alta

### BL-035 Definir flujo minimo navegable de zona y nodos para el MVP
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: bajar la expedicion a una estructura minima jugable y secuencial, con el menor numero de nodos y transiciones necesario para validar una run completa sin entrar todavia en mapa rico.
- **Resultado esperado**: flujo concreto tipo `inicio -> pesca -> pesca -> servicio -> boss -> cierre`, con taxonomia minima de nodos, reglas de avance y puntos donde se inyectan encounters, tienda y cierre de zona.
- **Dependencias**: `BL-001`
- **Prioridad**: Alta

### BL-036 Definir recursos globales de run y reglas del hilo del MVP
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: concretar que recursos globales arrastra una expedicion y como el hilo actua como condicion principal de supervivencia a lo largo de toda la sesion.
- **Resultado esperado**: contrato de recursos de run con hilo, dinero y otros minimos necesarios; reglas de perdida, reparacion y daño provenientes de encounters o nodos; y criterios para no mezclar esos recursos con el runtime tactico.
- **Dependencias**: `BL-001`, `BL-035`
- **Prioridad**: Alta

### BL-037 Definir economia minima de fotos, dinero y liquidacion de servicio
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: recortar la economia editorial a su loop minimo jugable, suficiente para validar captura, conversion a dinero y decisiones de gasto sin entrar todavia en meta o reserva compleja.
- **Resultado esperado**: flujo economico `captura -> fotos -> liquidacion en servicio -> dinero -> gasto`, con una decision explicita sobre si la reserva entra ya en el MVP o se difiere al backlog posterior.
- **Dependencias**: `BL-001`, `BL-035`, `BL-036`
- **Prioridad**: Alta

### BL-038 Definir build minima y acciones de servicio del MVP
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: decidir el subconjunto minimo de mejoras, reparaciones y ajustes de build que el jugador puede realizar en servicios durante la primera run jugable.
- **Resultado esperado**: taxonomia reducida de acciones de servicio para el MVP, cubriendo al menos reparacion del hilo y una pequena capa de mejora o ajuste sobre `rod`, aditamentos o recursos del jugador.
- **Dependencias**: `BL-008`, `BL-036`, `BL-037`
- **Prioridad**: Alta

### BL-039 Definir patrocinadores MVP o postergarlos explicitamente
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: decidir si el primer vertical slice jugable ya necesita patrocinadores y, en caso afirmativo, cual es su forma minima sin arrastrar aun el sistema completo de afinidades y oferta ponderada.
- **Resultado esperado**: decision cerrada entre `sin patrocinadores en el primer slice` o `patrocinadores MVP minimos`, con su impacto en la identidad de run y su lugar en el flujo de zona.
- **Dependencias**: `BL-001`, `BL-035`, `BL-037`
- **Prioridad**: Media

### BL-040 Implementar primer vertical slice de run MVP
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: entregar una expedicion minima jugable en una sola sesion, con flujo de zona reducido, encounters integrados, hilo global, economia minima y decisiones de servicio suficientes para validar el bucle de run.
- **Resultado esperado**: primer vertical slice de run funcional, UI-agnostic en su capa de aplicacion, sin persistencia entre runs y sin requerir aun mapa rico, bestiario ni guardado.
- **Dependencias**: `BL-034`, `BL-035`, `BL-036`, `BL-037`, `BL-038`, `BL-039`
- **Prioridad**: Alta

### BL-002 Definir estructura de mapa y tipos de nodo
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: decidir como evoluciona el MVP de flujo secuencial hacia un mapa real por zonas, con bifurcaciones, informacion parcial y convergencias legibles.
- **Resultado esperado**: modelo de mapa por zona posterior al MVP base, con longitud objetivo, puntos de convergencia, familias de nodos, vocabulario de pistas tematicas y reglas de ocultacion del siguiente nodo.
- **Dependencias**: `BL-035`, `BL-040`
- **Prioridad**: Alta

### BL-003 Disenar progresion de dificultad entre zonas
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: decidir como escala el reto desde la costa inicial hasta mar abierto y los peces legendarios finales sin resetear la build ni la economia de run.
- **Resultado esperado**: reglas de progresion por zona para pools de peces, riesgo de nodos, presion de tiempo, desgaste del hilo, frecuencia de servicios y preparacion del boss final.
- **Dependencias**: `BL-035`, `BL-040`, `BL-002`
- **Prioridad**: Alta

## Fish y Encounters

### BL-007 Expandir perfiles data-driven de pez
- **Estado**: `cancelled`
- **Tipo**: Delivery
- **Objetivo**: item absorbido por tareas data-driven mas pequenas y modulares para no mezclar catalogo de peces, metadata de encounters, aparicion por rol y contenido editorial en una sola pieza demasiado amplia.
- **Resultado esperado**: el trabajo queda redistribuido en `BL-041`, `BL-042`, `BL-043`, `BL-044` y `BL-045`, cada una con ownership y dependencias mas concretas.
- **Dependencias**: `BL-006`
- **Prioridad**: Media

### BL-022 Definir aparicion de peces por aguas base y ventana de lanzamiento
- **Estado**: `done`
- **Tipo**: Discovery
- **Objetivo**: decidir como la combinacion entre aguas base del nodo, franja horizontal del cast y sesgo vertical del setup del jugador selecciona subconjuntos de peces y define la aparicion de encuentros compatibles con esa ventana de lanzamiento.
- **Resultado esperado**: modelo de pool base por nodo o zona, reglas de particion en subconjuntos por distancia y profundidad, metadata minima de aparicion en perfiles de pez y criterio para conectar cast, `rod`, aditamentos y habitats sin abrir todavia la capa de economia fotografica.
- **Dependencias**: `BL-003`, `BL-007`, `BL-020`, `BL-021`
- **Plan relacionado**: `docs/features/018-aparicion-de-peces-por-aguas-y-ventana-de-lanzamiento.md`
- **Notas de cierre**:
  - La eleccion manual del preset de pez ya no es el flujo principal; el juego resuelve `agua -> apertura -> spawn -> mazo del pez`.
  - Los perfiles de pez ya exponen metadata minima de aparicion por pool de agua, distancia, profundidad y habitats.
  - El spawn actual ya usa catalogos tipados para `water pools`, `habitats` y arquetipos, dejando menos superficie de strings crudos en runtime y UI.
- **Direccion actual acordada**:
  - Cada nodo de pesca parte de unas aguas base y luego se secciona en subconjuntos de pool segun la ventana horizontal del cast y la ventana vertical habilitada por la `rod` y sus aditamentos.
  - La aparicion de peces debe resolverse a partir de la apertura ya cerrada del encounter (`InitialDistance`, `InitialDepth`) y no directamente desde los limites de escape del tablero.
  - La aparicion de peces debe seguir siendo compatible con zonas, elites, bosses y roles especiales de encounter.
  - El sistema debe permitir que a veces el jugador tenga lectura explicita al entrar al nodo y a veces no, sin exigir que el mapa revele antes la estructura del subpool.
  - La capa de economia o tags editoriales de fotografia queda fuera de este item y se aborda dentro de la extension data-driven del pez.
- **Prioridad**: Alta

### BL-041 Externalizar catalogo base de peces a formato data-driven
- **Estado**: `done`
- **Tipo**: Delivery
- **Objetivo**: mover el catalogo actual de perfiles de pez fuera del codigo a un formato data-driven validable y definir pools cerradas de peces por encounter a partir de listas de `profile_ids`, para acelerar configuracion y ajuste de encuentros sin depender de recompilar ni de avanzar todavia en el sistema de nodos.
- **Resultado esperado**: loader y validacion de catalogo global de peces, schema o contrato estable para perfiles, mazos, arquetipos y metadata minima de aparicion, mas una capa de `fish pools` cerradas y reutilizables que referencien perfiles por id sin duplicar su definicion.
- **Dependencias**: `BL-006`, `BL-022`
- **Plan relacionado**: `docs/features/024-externalizar-catalogo-base-de-peces-y-fish-pools.md`
- **Notas de cierre**:
  - El catalogo base de peces ya vive en JSON embebido con carga y validacion explicita.
  - Las `fish pools` ya soportan subsets cerrados por id y entradas ponderadas por peso.
  - El spawn ya puede trabajar sobre subcatalogos resueltos desde pools sin cambiar su contrato base de `[]Profile`.
- **Prioridad**: Alta

### BL-042 Externalizar arquetipos y patrones de cartas de pez
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: separar del codigo la definicion reusable de arquetipos, patrones y bloques de cartas del pez para que nuevos perfiles puedan componerse desde datos sin duplicar wiring manual.
- **Resultado esperado**: capa data-driven para arquetipos o plantillas de mazo del pez, reutilizable por el catalogo de perfiles y preparada para crecer en complejidad de contenido.
- **Dependencias**: `BL-041`
- **Prioridad**: Media

### BL-043 Externalizar metadata editorial y de encounter de especies
- **Estado**: `pending`
- **Tipo**: Discovery + Delivery
- **Objetivo**: sacar a datos la metadata propia de cada especie o encuentro que despues alimentara fotografia, descubrimiento, roles especiales o lectura tematica, sin mezclarla todavia con economia o coleccion meta.
- **Resultado esperado**: contrato data-driven para nombre, descripcion, tags, rareza, notas editoriales y otros campos semanticos del pez que no pertenecen al runtime tactico puro.
- **Dependencias**: `BL-041`
- **Prioridad**: Media

### BL-044 Externalizar tablas de aparicion de peces por contexto de encounter
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: preparar las reglas de aparicion para que puedan resolverse desde tablas data-driven por agua, distancia, profundidad, habitats y futuros roles de encounter, siempre operando sobre pools cerradas de peces y no sobre el catalogo global completo.
- **Resultado esperado**: tablas o reglas configurables de spawn consumibles por el resolver actual, capaces de seleccionar una `fish pool` concreta por contexto de encounter y extensibles luego a zonas, elites, bosses y nodos de pesca.
- **Dependencias**: `BL-022`, `BL-041`
- **Prioridad**: Alta

### BL-045 Resolver spawns de encounter desde catalogos data-driven
- **Estado**: `planned`
- **Tipo**: Delivery
- **Objetivo**: conectar el flujo actual `agua -> apertura -> spawn -> mazo del pez` a los nuevos catalogos data-driven para que el runtime deje de depender de listas hardcodeadas de perfiles y pueda limitar cada encounter a una `fish pool` cerrada.
- **Resultado esperado**: pipeline de spawn consumiendo catalogos externos de peces, pools cerradas y reglas de aparicion, manteniendo tests, comportamiento de encuentro y contratos de app actuales.
- **Dependencias**: `BL-041`, `BL-044`
- **Plan relacionado**: `docs/features/025-resolver-spawns-desde-catalogos-y-fish-pools.md`
- **Prioridad**: Alta

### BL-050 Reducir el determinismo total del spawn de peces
- **Estado**: `done`
- **Tipo**: Discovery + Delivery
- **Objetivo**: evitar que el sistema actual de spawn sea completamente deducible para el jugador cuando conoce agua, distancia, profundidad y habitats, introduciendo variedad controlada sin perder legibilidad, control de balance ni capacidad de testeo.
- **Resultado esperado**: estrategia clara para desempates, pesos, variacion por pool o sampling controlado dentro de candidatos validos, con contratos reproducibles para tests y configuracion futura desde data-driven.
- **Dependencias**: `BL-022`, `BL-041`
- **Plan relacionado**: `docs/features/023-reducir-determinismo-del-spawn-de-peces.md`
- **Notas de cierre**:
  - El matching contextual del spawn se mantiene estable y separado de la seleccion final del pez.
  - El runtime ya permite inyectar randomizer para variar reproduciblemente entre candidatos empatados en score.
  - El camino sin randomizer sigue siendo estable para tests y para diagnostico del dominio.
- **Direccion actual acordada**:
  - El spawn no debe elegir sobre todo el catalogo global si el encounter ya opera sobre una `fish pool` cerrada.
  - Dentro de una pool valida, conviene reducir el determinismo absoluto del `best match` actual.
  - La variacion debe seguir siendo controlable y reproducible para debug, seeds y tests.
  - La solucion inicial deberia priorizar pesos o sampling simple antes que sistemas opacos o muy dificiles de balancear.
- **Prioridad**: Alta

## Items y Build

### BL-008 Definir categorias de objetos del jugador
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: clasificar las piezas de build que el jugador puede comprar, recibir o mejorar durante la expedicion sin solaparlas con los patrocinadores globales, partiendo de la build minima ya decidida para el MVP.
- **Resultado esperado**: taxonomia completa de acciones y objetos de tienda que expanda la build minima del MVP hacia mejoras de la `rod`, aditamentos, reparacion del hilo, intervenciones sobre el mazo, consumibles y otros ajustes de servicio.
- **Dependencias**: `BL-038`, `BL-011`
- **Prioridad**: Alta

### BL-021 Redefinir la cana como `rod` y sus aditamentos como base del setup del jugador
- **Estado**: `done`
- **Tipo**: Discovery
- **Objetivo**: reemplazar el vocabulario ambiguo de `rig` por un modelo explicito de `rod` para la cana base, separando sus limites estructurales de track frente a sus limites de apertura, y distinguiendo esa pieza base de los aditamentos que modifican la apertura vertical del encounter y el acceso a distintos habitats de pez.
- **Resultado esperado**: modelo de `rod` con limites estructurales separados para apertura y track, taxonomia inicial de aditamentos, nomenclatura clara para distinguir `rod` frente a setup completo del jugador, y reglas de tradeoff para que el equipo no se lea como una pila de mejoras universales.
- **Dependencias**: `BL-008`, `BL-020`
- **Plan relacionado**: `docs/features/017-rod-y-limites-de-apertura-y-track.md`
- **Notas de cierre**:
  - El vocabulario de `rig` ya fue sustituido por `rod` y el estado runtime usa `loadout` como composicion de `rod` y aditamentos.
  - La apertura del encounter ya valida `InitialDistance` e `InitialDepth` contra limites efectivos de apertura, mientras el tablero y los escapes usan limites de track.
  - El CLI ya permite elegir presets de `rod` y presets de aditamentos antes de abrir el encounter.
  - Los aditamentos ya modifican limites de apertura y track, y transportan `HabitatTags` listos para alimentar `BL-022`.
- **Direccion actual acordada**:
  - `rig` debe desaparecer como termino de dominio y pasar a `rod` cuando se hable de la cana base del jugador.
  - `rod` debe referirse a la pieza base de equipo; el conjunto `rod + aditamentos` debe nombrarse como setup o loadout para no volver a mezclar capas.
  - Los limites de apertura deben separarse de los limites de track: una cosa es hasta donde puede empezar la pesca tras el cast y otra hasta donde puede sostenerse el duelo antes del escape.
  - La `rod` debe exponer al menos una pareja de limites de apertura y otra de track, en horizontal y vertical, en vez de reutilizar un unico `MaxDistance` o `MaxDepth` para ambas fases.
  - Los limites de track siguen definiendo tablero, render y condiciones de escape; los limites de apertura solo validan hasta donde puede resolverse `InitialDistance` y `InitialDepth`.
  - Los aditamentos no reemplazan la `rod`; sesgan sobre todo la apertura vertical, la compatibilidad con ciertos habitats y la forma de acceder a subpools de peces.
  - La `rod` conserva el rol de limite estructural del encuentro y no debe mutarse por cartas durante el duelo.
  - Los setups deben tener costes y ventajas reales; por ejemplo, mas profundidad de apertura o mas tolerancia de track no deberian ser mejoras universales sin contrapartida.
- **Prioridad**: Alta

### BL-009 Disenar sistema de sinergias
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: definir como interactuan cartas del jugador, mejoras de `rod`, aditamentos, hilo, patrocinadores, fotos reservadas y economia editorial sin generar combinaciones opacas o imposibles de leer.
- **Resultado esperado**: reglas de stacking, limites, superficies de activacion y familias de build que mantengan el juego expresivo y entendible.
- **Dependencias**: `BL-005`, `BL-008`, `BL-011`
- **Prioridad**: Media

### BL-010 Implementar primer vertical slice de build
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: expandir el primer slice de run MVP hacia una build mas expresiva y una expedicion mas rica una vez validada la base jugable.
- **Resultado esperado**: prototipo ampliado con mapa de zona mas completo, build mas profunda, recompensas mas variadas y progresion mas rica que la del primer vertical slice base.
- **Dependencias**: `BL-040`, `BL-002`, `BL-003`, `BL-007`, `BL-008`, `BL-009`, `BL-011`, `BL-012`
- **Prioridad**: Media

## Economy y Meta

### BL-011 Definir economia ampliada de run y frontera futura de meta-progresion
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: ampliar la economia minima del MVP hacia una economia editorial mas rica y dejar formalizada la frontera entre recursos de run y futura meta-progresion, sin obligar a implementarla todavia.
- **Resultado esperado**: modelo economico ampliado con fotos, dinero, reserva, venta automatica, sinks de tienda, valor editorial y frontera explicita entre progreso de run y futura meta.
- **Dependencias**: `BL-037`, `BL-040`
- **Prioridad**: Alta

### BL-012 Disenar sistema de recompensas entre encuentros
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: decidir como se distribuyen fotos, dinero, ofertas de patrocinador, acceso a servicio, reparaciones y otras recompensas entre nodos y entre zonas.
- **Resultado esperado**: tabla de recompensas por tipo de nodo, outcome, zona y rol del pez para sostener el ritmo de la expedicion sin inflacion de recursos.
- **Dependencias**: `BL-040`, `BL-002`, `BL-008`, `BL-011`
- **Prioridad**: Alta

## Collection

### BL-013 Definir bestiario y coleccion
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: concretar como se registra entre runs la informacion descubierta de cada pez, que muestran los trofeos fotografiados y como se presenta el bestiario al jugador.
- **Resultado esperado**: definicion de registro de especies, composicion descubierta de sus barajas, categorias de legendarios y presentacion de trofeos sin convertir la coleccion en una fuente de poder.
- **Dependencias**: `BL-040`, `BL-011`
- **Prioridad**: Media

### BL-014 Disenar recompensas por coleccion
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: decidir si la coleccion permanece solo como trofeo, codex e informacion, o si desbloquea contenido no-poderoso como cosmeticos, lore o pistas adicionales.
- **Resultado esperado**: politica de recompensas de coleccion coherente con la direccion actual de mantener los trofeos como meta no mecanica.
- **Dependencias**: `BL-013`
- **Prioridad**: Baja

## Tech y UX

### BL-015 Disenar sistema de contenido data-driven
- **Estado**: `cancelled`
- **Tipo**: Discovery
- **Objetivo**: item descartado por ser demasiado amplio; su scope se redistribuye en tareas data-driven especificas por modulo para poder ejecutarlas de forma independiente conforme el flujo del desarrollo lo vaya habilitando.
- **Resultado esperado**: el plan global se reemplaza por items concretos para peces, apariciones, nodos, servicios, recompensas, patrocinadores y otros catalogos, evitando una mega-tarea transversal dificil de cerrar.
- **Dependencias**: ninguna
- **Prioridad**: Media

### BL-046 Externalizar catalogo de nodos y rutas del MVP
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: mover a datos la estructura minima de nodos y rutas del MVP de run una vez exista el flujo base de expedicion, para evitar que la primera navegacion de zona quede cableada a mano.
- **Resultado esperado**: catalogo data-driven para nodos, secuencias y variantes de ruta del MVP, desacoplado del render y del runtime tactico.
- **Dependencias**: `BL-035`, `BL-040`
- **Prioridad**: Media

### BL-047 Externalizar catalogo de servicios y acciones de build
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: llevar a datos la oferta minima de servicios, reparaciones y mejoras de run para que el MVP pueda crecer sin ampliar wiring duro en codigo.
- **Resultado esperado**: catalogo data-driven de acciones de servicio y build consumible por la capa de run y alineado con `BL-038`.
- **Dependencias**: `BL-038`, `BL-040`
- **Prioridad**: Media

### BL-048 Externalizar recompensas y economia contextual
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: convertir a catalogos configurables la distribucion de recompensas, conversion economica y variaciones contextuales que hoy nazcan primero como reglas de MVP.
- **Resultado esperado**: tablas data-driven de recompensas, payout editorial y variaciones por nodo u outcome, listas para crecer despues del primer vertical slice.
- **Dependencias**: `BL-037`, `BL-012`, `BL-040`
- **Prioridad**: Media

### BL-049 Externalizar catalogo de patrocinadores y ofertas
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: mover a datos las ofertas, efectos y compatibilidades basicas de patrocinadores cuando el sistema ya este validado como parte de run.
- **Resultado esperado**: catalogo data-driven de patrocinadores y reglas de oferta, separado del runtime de encounter y de la UI.
- **Dependencias**: `BL-039`, `BL-040`
- **Prioridad**: Baja

### BL-016 Disenar guardado de run y progreso meta
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: decidir como persistir expediciones en curso y datos meta como bestiario, trofeos legendarios y conocimiento ya descubierto del pez.
- **Resultado esperado**: contrato de persistencia para run activa y perfil permanente, con limites de versionado y reglas de compatibilidad.
- **Dependencias**: `BL-040`, `BL-011`, `BL-013`
- **Prioridad**: Media

### BL-017 Mejorar UX de lectura de build y estado de run
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: asegurar que el mapa, el hilo, los patrocinadores, la reserva de fotos, la economia y la build del jugador sigan siendo legibles durante la expedicion.
- **Resultado esperado**: requerimientos de HUD, feedback de ruta, lectura de servicios, valuacion editorial de fotos, resumen de build y transparencia de efectos globales.
- **Dependencias**: `BL-040`, `BL-008`, `BL-011`
- **Prioridad**: Media

### BL-023 Desacoplar flujos de interfaz y preparar arquitectura UI-agnostic
- **Estado**: `pending`
- **Tipo**: Discovery + Delivery
- **Objetivo**: refactorizar los flujos de setup, opening, presentacion y sesion para que la aplicacion deje de depender estructuralmente del CLI como unica forma de presentar e interactuar con el juego, dejando a `internal/cli/` como un adaptador de borde y preparando una futura UI grafica para consumir los mismos casos de uso.
- **Resultado esperado**: arquitectura donde el combate, el setup previo al encounter y la apertura del lance se expresen como casos de uso UI-agnostic; presenter y view models mas semanticos para setup/opening/combate; controlador reusable del cast fuera del CLI; y `cmd/` reducido a composicion de dependencias, de modo que una GUI futura no tenga que duplicar wiring ni parsear strings finales para reconstruir la experiencia.
- **Dependencias**: `BL-034`, `BL-018`, `BL-020`, `BL-021`
- **Contexto actual detectado**:
  - El loop de combate ya tiene una separacion sana entre motor, `app`, `presentation` y `cli`, especialmente en `internal/app/session.go`, donde la sesion consume interfaces de `Engine`, `UI` y `Presenter`.
  - El mayor acoplamiento actual vive en el setup previo a la partida y en la apertura del encounter: `cmd/fishing-duel/main.go` orquesta seleccion de deck, `rod`, aditamentos, pez y opening contra la UI concreta.
  - La logica del cast por timing esta incrustada en el adaptador CLI (`internal/cli/opening.go`), usando `scanner`, `time.Sleep` y redibujos de pantalla completos; eso no traslada bien a una UI grafica basada en eventos o ticks.
  - Parte de la presentacion del opening y de los presets aun se renderiza directamente desde tipos de dominio o contenido (`encounter.Opening`, presets de `rod`, aditamentos y agua) en vez de pasar por view models estables.
  - Los view models actuales de `presentation` ayudan para el combate, pero todavia cargan bastante texto final o labels compactados, lo que obligaria a una GUI a reutilizar strings ya formateados en vez de consumir datos estructurados.
  - El adaptador CLI mantiene estado visual propio (`intro`, `lastRound`, `opening`) y asume un flujo de repintado full-screen, senal de que parte de la UX aun vive dentro del borde de terminal.
- **Direccion actual acordada**:
  - El objetivo no es implementar una GUI todavia, sino dejar la arquitectura lista para que `CLI` y futura `GUI` cuelguen de los mismos flujos de aplicacion.
  - `cmd/` debe quedar como composition root y no como lugar donde vive la orquestacion de setup, seleccion de presets o minijuegos previos al combate.
  - El flujo previo a la sesion debe subir a `internal/app/` o a un subpaquete estable de casos de uso, para que elegir deck, `rod`, aditamentos, pez, contexto de agua y opening no dependa del adaptador concreto.
  - La logica temporal del cast debe extraerse a un controlador o estado reusable fuera de `internal/cli/`, de modo que cada interfaz solo decida como dibujarlo y capturar input.
  - `internal/presentation/` debe ampliarse para cubrir setup y opening, y a la vez volverse menos text-first en puntos clave: la GUI futura debe poder consumir datos semanticos y no solo labels finales.
  - El adaptador CLI debe tender a renderizar exclusivamente views o modelos de pantalla, evitando conocer de primera mano demasiado contenido configurable o runtime del dominio.
  - La migracion debe preservar la CLI actual como primer adaptador y mantener el juego jugable durante todo el refactor.
- **Propuesta de refactor por fases**:
  - `Fase 1`: extraer el flujo de setup desde `cmd/fishing-duel/` a una capa de aplicacion (`internal/app/setup/` o equivalente) con interfaces UI-neutrales para elecciones y confirmaciones.
  - `Fase 2`: mover la apertura del encounter a un flujo de aplicacion mas completo, donde la seleccion de agua, preview de loadout y resumen de apertura no dependan del renderer CLI.
  - `Fase 3`: sacar el cast timing de `internal/cli/opening.go` a un controlador reusable, con estado/ticks/resultado legibles por cualquier frontend.
  - `Fase 4`: ampliar `internal/presentation/` con view models de setup, opening, loadout y cast, reduciendo el paso directo de presets o structs de dominio al render.
  - `Fase 5`: adelgazar `internal/cli/` para que quede como adaptador de input/output de terminal, no como lugar donde vive la UX global de la app.
  - `Fase 6`: dejar preparado un segundo entrypoint futuro (`cmd/...`) para GUI sin duplicar el wiring de juego.
- **Criterios de cierre**:
  - el bootstrap ya no depende de metodos tipados contra `*cli.UI`
  - existe un flujo de setup/opening reusable desde `app` y consumible por mas de un adaptador
  - la logica del cast ya no vive en el paquete CLI
  - setup y opening tienen view models propios en `presentation`
  - la CLI sigue funcionando como adaptador despues del refactor
  - queda listo el camino para crear una GUI sin reescribir combate, setup ni opening
- **Riesgos o decisiones abiertas**:
  - si se intenta resolver a la vez la arquitectura UI-agnostic y una GUI completa, el scope puede crecer demasiado.
  - habra que decidir cuanto de `presentation` sigue siendo textual y cuanto pasa a ser estructurado sin romper la legibilidad actual de la CLI.
  - el cast puede necesitar un pequeño modelo de estado orientado a ticks o eventos antes de encajar bien en una GUI.
  - algunos flujos de seleccion hoy basados en presets concretos pueden necesitar un contrato mas abstracto para no arrastrar contenido directo hasta el borde.
- **Prioridad**: Media

### BL-029 Descomponer `match.State` y fijar fronteras del runtime tactico
- **Estado**: `done`
- **Tipo**: Calidad + Discovery
- **Objetivo**: reducir el acoplamiento actual del runtime tactico separando encounter, estado de mazo, recursos del jugador, thresholds de ronda y vistas de presentacion, evitando que `match.State` siga creciendo como estado universal del juego.
- **Resultado esperado**: ownership mas claro del estado mutable del combate, contratos mas estrechos entre `game`, `progression`, `endings`, `presentation` y futuros estados de run, y una ruta concreta para convivir con un `RunState` sin convertir `match.State` en un god object.
- **Dependencias**: `BL-018`
- **Plan relacionado**: `docs/features/019-descomponer-match-state-y-fronteras-runtime-tactico.md`
- **Notas de cierre**:
  - El runtime tactico ya separa `Round`, `Player` y `Lifecycle` como subestados explicitos dentro de `match.State`.
  - `engine`, `progression`, `endings`, `presentation`, `session` y el runtime de recursos del jugador ya consumen esa frontera mas explicita.
  - Queda fijado que `match.State` representa solo un duelo aislado y no debe absorber estado futuro de run.
- **Contexto actual detectado**:
  - `match.State` concentra encounter, deck, loadout, recursos del jugador, stats y flag de fin del duelo en una sola estructura compartida.
  - `game`, `progression`, `endings`, `presentation` y `app/session` consumen o mutan directamente ese estado, lo que hace que cualquier ampliacion ripplee varias capas.
  - El crecimiento futuro de la run, servicios, rewards y persistencia amenaza con empujar mas datos hacia este estado si no se fijan fronteras antes.
- **Direccion actual acordada**:
  - `match.State` debe seguir siendo tactico y no absorber runtime de expedicion.
  - El combate necesita subestados o contratos mas finos para aislar responsabilidades sin reescribir el engine completo en una sola iteracion.
  - La refactorizacion debe ocurrir antes de expandir `BL-001` para evitar que la run nazca sobre un estado tactico ya sobredimensionado.
  - Debe priorizar claridad de ownership y compatibilidad con futuras capas de run antes que una estetica perfecta de tipos.
- **Prioridad**: Alta

### BL-030 Consolidar runtime de combate y fronteras de paquetes
- **Estado**: `done`
- **Tipo**: Calidad + Delivery
- **Objetivo**: reducir la fragmentacion funcional actual entre `internal/game/`, `internal/rules/`, `internal/progression/`, `internal/endings/` e `internal/encounter/`, dejando ownership mas claro de la logica del duelo y preparando una arquitectura de combate mas estable para seguir creciendo.
- **Resultado esperado**: mapa de responsabilidades y una refactorizacion acotada que reduzca imports cruzados, clarifique que parte del combate evalua, que parte progresa, que parte cierra el encounter y que parte modela runtime, alineandose con la direccion ya expresada en la feature de arquitectura de paquetes.
- **Dependencias**: `BL-018`, `BL-029`
- **Plan relacionado**: `docs/features/020-consolidar-runtime-de-combate-y-fronteras-de-paquetes.md`
- **Notas de cierre**:
  - `presentation` ya consume snapshots tacticos mas estrechos para status, ronda y resumen final en lugar de depender del ensamblado completo del duelo.
  - Los thresholds y helpers de captura quedaron consolidados bajo `encounter`, dejando a `game` mas centrado en orquestacion.
  - La consolidacion principal del runtime ya esta integrada y la deuda residual queda trazada en `BL-033`.
- **Contexto actual detectado**:
  - La capacidad de combate sigue repartida entre varios paquetes transicionales y no converge todavia en una frontera estable.
  - Las reglas de progresion, evaluacion y cierre del encounter ya estan separadas, pero la expansion futura de zonas, roles de encounter o run puede volver mas opaca la superficie entre esos paquetes.
  - El propio plan de arquitectura ya reconocio esta convergencia como direccion futura, pero no existe un item especifico para materializarla.
- **Direccion actual acordada**:
  - La meta no es un gran rewrite, sino completar una consolidacion incremental del runtime de combate.
  - Esta tarea debe apoyarse en fronteras mejores de estado y no al reves.
  - Debe incluir la revision de si `presentation` puede depender de snapshots tacticos mas estrechos que el ensamblado completo actual.
  - La reorganizacion debe preservar el comportamiento actual del duelo y su capacidad de testeo.
- **Prioridad**: Media

### BL-033 Adelgazar contratos residuales del runtime tactico post-`BL-030`
- **Estado**: `done`
- **Tipo**: Calidad + Delivery
- **Objetivo**: cerrar el acoplamiento residual que sigue concentrado en contratos transicionales del duelo, especialmente en `match.RoundResult`, el mapeo de estado visible del mazo y las mutaciones directas de `*match.State` desde colaboradores tacticos, para que el runtime quede listo para crecer sin volver a recentralizar conocimiento en `game`.
- **Resultado esperado**: superficies de lectura y escritura mas estrechas entre `game`, `match`, `progression`, `endings` y `encounter`; ownership mas claro del estado visible del mazo; y resultados de ronda que expongan solo la informacion tactica que realmente necesitan las capas consumidoras.
- **Dependencias**: `BL-030`
- **Plan relacionado**: `docs/features/021-adelgazar-contratos-residuales-del-runtime-tactico.md`
- **Notas de cierre**:
  - `match.RoundResult` ya devuelve snapshots tacticos estrechos en lugar de reexportar `match.State` completo.
  - El mapping del estado visible del mazo ya vive bajo `match`, dejando a `game` mas centrado en orquestacion.
  - `progression`, `endings` y `player/playermoves` ya consumen contratos de mutacion mas finos (`ProgressionState`, `EndingState`, `PlayerMoveRuntime`).
- **Justificacion**:
  - `BL-030` ya resolvio la consolidacion principal, pero el merge deja visible una segunda capa de deuda mas fina que no conviene arrastrar a `BL-023` ni al futuro runtime de run.
  - Si `RoundResult` y los colaboradores tacticos siguen apoyandose en `match.State` o en mapeos orquestados desde `game`, el proyecto corre el riesgo de reintroducir un nuevo centro de gravedad accidental en el engine.
  - Capturar este follow-up como backlog explicito permite cerrar la frontera del duelo con una tarea acotada, en vez de diluirla dentro de futuras features de UI o producto.
- **Contexto actual detectado**:
  - El duelo ya tiene snapshots tacticos utiles para `presentation`, pero el resultado de ronda sigue siendo mas ancho de lo necesario para algunas capas consumidoras.
  - Parte del estado visible del mazo sigue dependiendo de mapeos ensamblados desde el engine en lugar de vivir claramente junto a su owner natural.
  - `progression` y `endings` todavia pueden apoyarse en mutaciones directas del ensamblado tactico completo en puntos donde convendria evaluar contratos mas finos.
- **Direccion actual acordada**:
  - Revisar si `match.RoundResult` puede adelgazar su superficie sin perder legibilidad ni testabilidad.
  - Mover el mapeo del estado visible del mazo hacia un owner mas natural que `game` cuando el cambio reduzca coupling real.
  - Evaluar interfaces o subestados mas estrechos para `progression` y `endings` sin convertir la tarea en un rewrite total del combate.
  - Mantener fuera de scope una reorganizacion cosmetica de paquetes que no mejore ownership real.
- **Prioridad**: Alta

### BL-031 Centralizar constantes y politicas de balance del encounter
- **Estado**: `pending`
- **Tipo**: Calidad + Discovery
- **Objetivo**: reducir magic numbers, strings repetidos y heuristicas opacas dentro del loop tactico, moviendo defaults y reglas sensibles a puntos de configuracion o politicas nombradas y documentadas.
- **Resultado esperado**: ownership claro para constantes de reciclado, recovery, agotamiento, timing de cast y scoring de spawn; menos drift entre runtime, UI y contenido; y una base mas segura para balancear el juego sin perseguir valores escondidos en multiples archivos.
- **Dependencias**: `BL-020`, `BL-021`, `BL-022`
- **Contexto actual detectado**:
  - Hay reglas importantes repartidas entre valores numericos embedidos en runtime, presenter, deck y contenido por defecto.
  - El resolver de spawn usa heuristicas numericas funcionales pero poco explicitadas, lo que dificulta tuning y lectura de intencion.
  - Parte de los textos de flujo y titulos globales se repite en varias capas del borde CLI o de presentacion.
- **Direccion actual acordada**:
  - No todo numero duro debe desaparecer, pero los que expresan politica de balance o contrato de UX deben tener nombre y ownership.
  - Conviene distinguir constantes de gameplay, constantes de UI y defaults de contenido para no mezclar capas.
  - La tarea debe dejar preparada una base de tuning antes de que la run y el contenido crezcan demasiado.
- **Prioridad**: Media

### BL-032 Automatizar pipeline tecnico de calidad
- **Estado**: `pending`
- **Tipo**: Infra + Delivery
- **Objetivo**: asegurar que test, lint y checks basicos de integracion no dependan solo del entorno local del desarrollador, formalizando un pipeline minimo de calidad y entrypoints consistentes para el workflow tecnico del proyecto.
- **Resultado esperado**: automatizacion reproducible para `go test`, `golangci-lint` y checks clave del repositorio, con una ruta simple para correrlos localmente y una base minima para CI o equivalentes.
- **Dependencias**: `BL-018`
- **Contexto actual detectado**:
  - El repo ya usa buenas herramientas de calidad, pero no tiene todavia una capa visible de automatizacion de pipeline en el propio proyecto.
  - Al crecer el backlog y los refactors, depender solo del ritual manual aumenta el riesgo de drift en el workflow.
  - La ausencia de entrypoints o automatizacion minima tambien dificulta colaborar o escalar el proyecto con menos friccion.
- **Direccion actual acordada**:
  - La tarea debe ser ligera: no busca introducir infraestructura sobredimensionada, sino endurecer el workflow real ya usado por el proyecto.
  - Debe respetar la simplicidad actual del repo y evitar acoplar la solucion a una plataforma unica si no es necesario.
- **Prioridad**: Media

## Completados

### BL-019 Hacer visible el descarte del pez y modular la lectura del historial
- **Estado**: `done`
- **Tipo**: Discovery + Delivery
- **Objetivo**: convertir el descarte del pez en una herramienta estrategica legible durante el encounter, manteniendo espacio para peces, cartas o eventos que oculten parcial o temporalmente esa informacion.
- **Resultado esperado**: estado de runtime, presentacion y UX que permitan ver el descarte visible del pez por ciclo, entender cuando el mazo recicla o se baraja y soportar excepciones de visibilidad por carta, arquetipo o evento.
- **Dependencias**: `BL-005`, `BL-006`, `BL-017`
- **Plan relacionado**: `docs/features/015-visibilidad-descarte-del-pez.md`
- **Prioridad**: Media

### BL-020 Disenar apertura del encounter de pesca y minijuego de cast
- **Estado**: `done`
- **Tipo**: Discovery + Delivery
- **Objetivo**: definir la fase previa al combate de pesca como una apertura autocontenida del encounter, donde el jugador lea la situacion del agua, ejecute un cast por timing y abra una ventana horizontal real sin depender todavia del sistema final de nodos o mapa, dejando ademas preparado al encounter para recibir `InitialDistance` e `InitialDepth` como entradas resueltas.
- **Resultado esperado**: flujo claro `leer situacion -> resolver cast -> abrir ventana horizontal`, vocabulario de aguas base y senales de lectura, bandas de distancia del cast, contrato minimo para inyectar contexto de agua desde presets o configuracion temporal, separacion entre informacion estetica visible y metadata interna de pool, y politica de fallo con impacto controlado sobre la apertura del encounter.
- **Dependencias**: `BL-001`
- **Plan relacionado**: `docs/features/016-apertura-encounter-pesca-y-cast.md`
- **Notas de cierre**:
  - El encounter ya resuelve una apertura previa con contexto de agua, minijuego de cast y configuracion inicial derivada.
  - El contrato queda listo para transportar `InitialDistance` e `InitialDepth`, aunque este slice solo altera la distancia inicial.
  - La version integrada sigue desacoplada de nodos, `rod`/aditamentos y subpools, que pasan a revisarse en `BL-002`, `BL-021` y `BL-022`.
- **Prioridad**: Alta

### BL-005 Disenar sistema extensible de efectos de cartas
- **Estado**: `done`
- **Tipo**: Delivery
- **Objetivo**: consolidar una arquitectura preparada para cartas de pez y de jugador sin seguir agregando casos aislados.
- **Resultado esperado**: sistema tecnico para encadenar efectos de carta por fases, trigger y owner.
- **Dependencias**: `docs/discoveries/008-taxonomia-de-efectos-de-carta.md`, plan 009, plan 012, plan 013
- **Plan relacionado**: `docs/features/009-pipeline-de-efectos-de-carta.md`, `docs/features/012-barajas-de-decision-del-jugador.md`, `docs/features/013-primeras-player-cards-con-efectos.md`
- **Prioridad**: Alta

### BL-006 Definir arquetipos de peces
- **Estado**: `done`
- **Tipo**: Discovery + Delivery
- **Objetivo**: definir arquetipos de pez faciles de configurar y llevarlos a una primera implementacion data-driven usable desde los presets del juego.
- **Resultado esperado**: perfiles mecanicos configurables que construyen barajas de pez y reemplazan wiring manual en presets de inicio.
- **Dependencias**: `docs/discoveries/008-taxonomia-de-efectos-de-carta.md`, plan 011, plan 013
- **Plan relacionado**: `docs/features/011-arquetipos-de-peces.md`
- **Prioridad**: Alta

### BL-018 Mejorar arquitectura y gobierno de paquetes
- **Estado**: `done`
- **Tipo**: Calidad + Delivery
- **Objetivo**: contener el crecimiento de `internal/` con reglas claras de organizacion y una primera refactorizacion acotada de paquetes.
- **Resultado esperado**: estrategia de estructura de paquetes mas sostenible, mas una mejora concreta que reduzca acoplamiento o dispersion actual.
- **Dependencias**: plan 014
- **Plan relacionado**: `docs/features/014-arquitectura-y-gobierno-de-paquetes.md`
- **Prioridad**: Alta

## Orden sugerido del trabajo pendiente

1. `BL-001`
2. `BL-034`
3. `BL-035`
4. `BL-041`
5. `BL-050`
6. `BL-036`
7. `BL-037`
8. `BL-038`
9. `BL-039`
10. `BL-040`
11. `BL-045`
12. `BL-044`
13. `BL-002`
14. `BL-003`
15. `BL-011`
16. `BL-008`
17. `BL-012`
18. `BL-046`
19. `BL-047`
20. `BL-048`
21. `BL-042`
22. `BL-043`
23. `BL-017`
24. `BL-013`
25. `BL-016`
26. `BL-010`
27. `BL-009`
28. `BL-049`
29. `BL-014`
30. `BL-023`
31. `BL-031`
32. `BL-032`
