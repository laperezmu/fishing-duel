# Backlog activo roguelike

Este documento concentra el backlog activo y pendiente del proyecto. Las tareas ya integradas en `main` viven en `docs/backlog/002-backlog-completado-roguelike.md`.

## Convencion de estado

- `pending`: item identificado pero aun sin plan activo en `docs/features/`.
- `planned`: item ya convertido en plan de feature, pero todavia no mergeado en `main`.
- `cancelled`: item descartado, absorbido por otro o fuera de foco.

## Archivos del backlog

- `docs/backlog/001-backlog-activo-roguelike.md`: backlog activo y pendiente.
- `docs/backlog/002-backlog-completado-roguelike.md`: tareas ya completadas e integradas en `main`.

## Foto actual

- `planned`: `BL-001`
- `pending`: resto del backlog activo
- `done`: ver `docs/backlog/002-backlog-completado-roguelike.md`
- Foco recomendado inmediato: cerrar primero la base MVP de una run sin persistencia meta y usar la base ya saneada del encounter para empezar a modelar la expedicion, empezando por `BL-001` y `BL-035`.

## Foco sugerido actual

- `BL-001`: fijar el contrato base y el MVP secuencial de una run sin persistencia meta.
- `BL-035`: definir el flujo minimo navegable de zona y nodos para el MVP de expedicion.
- `BL-044`: externalizar tablas de aparicion por contexto de encounter cuando ya queramos elegir `fish_pool_id` desde reglas de contenido.
- `BL-023`: mantener en cola el refactor UI-agnostic amplio despues del slice tecnico minimo.

## Core Loop

### BL-001 Definir contrato base y MVP secuencial de una run
- **Estado**: `planned`
- **Tipo**: Discovery
- **Objetivo**: fijar la arquitectura minima de una expedicion jugable en una sola sesion, separando claramente encounter, run y futura meta, para que la implementacion pueda avanzar por etapas UI-agnostic y sin depender todavia de persistencia entre runs.
- **Resultado esperado**: contrato marco del MVP de run con estados base, handoff `encounter -> run`, recursos globales de expedicion, flujo minimo de zona y lista secuencial de subtareas necesarias para llegar a una primera run jugable sin meta-persistencia.
- **Plan relacionado**: `docs/features/026-definir-contrato-base-y-mvp-secuencial-de-una-run.md`
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

### BL-035 Definir flujo minimo navegable de zona y nodos para el MVP
- **Estado**: `planned`
- **Tipo**: Discovery
- **Objetivo**: bajar la expedicion a una estructura minima jugable y secuencial, con el menor numero de nodos y transiciones necesario para validar una run completa sin entrar todavia en mapa rico.
- **Resultado esperado**: flujo concreto tipo `inicio -> pesca -> pesca -> servicio -> boss -> cierre`, con taxonomia minima de nodos, reglas de avance y puntos donde se inyectan encounters, tienda y cierre de zona.
- **Plan relacionado**: `docs/features/027-definir-flujo-minimo-de-zona-y-nodos-para-el-mvp.md`
- **Dependencias**: `BL-001`
- **Prioridad**: Alta

### BL-036 Definir recursos globales de run, reglas del hilo y `AnglerProfile` del MVP
- **Estado**: `planned`
- **Tipo**: Discovery
- **Objetivo**: concretar que recursos globales arrastra una expedicion, como el hilo actua como condicion principal de supervivencia a lo largo de toda la sesion y como el arranque de la run pasa a encapsularse dentro de un `AnglerProfile` seleccionable.
- **Resultado esperado**: contrato de recursos de run con hilo, dinero y otros minimos necesarios; reglas de perdida, reparacion y dano provenientes de encounters o nodos; y definicion del paquete inicial de run via `AnglerProfile`, cubriendo mazo inicial, `rod`, aditamentos base y `thread` inicial sin mezclar esos recursos con el runtime tactico.
- **Plan relacionado**: `docs/features/028-definir-recursos-globales-de-run-y-angler-profiles.md`
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

## Items y Build

### BL-008 Definir categorias de objetos del jugador
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: clasificar las piezas de build que el jugador puede comprar, recibir o mejorar durante la expedicion sin solaparlas con los patrocinadores globales, partiendo de la build minima ya decidida para el MVP.
- **Resultado esperado**: taxonomia completa de acciones y objetos de tienda que expanda la build minima del MVP hacia mejoras de la `rod`, aditamentos, reparacion del hilo, intervenciones sobre el mazo, consumibles y otros ajustes de servicio.
- **Dependencias**: `BL-038`, `BL-011`
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
- **Dependencias**: `BL-040`, `BL-002`, `BL-003`, `BL-008`, `BL-009`, `BL-011`, `BL-012`
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
- **Objetivo**: refactorizar los flujos de setup, opening, presentacion y sesion para que la aplicacion deje de depender estructuralmente del CLI como unica forma de presentar e interactuar con el juego.
- **Resultado esperado**: arquitectura donde el combate, el setup previo al encounter y la apertura del lance se expresen como casos de uso UI-agnostic y `cmd/` quede reducido a composicion de dependencias.
- **Dependencias**: `BL-034`, `BL-018`, `BL-020`, `BL-021`
- **Prioridad**: Media

### BL-031 Centralizar constantes y politicas de balance del encounter
- **Estado**: `pending`
- **Tipo**: Calidad + Discovery
- **Objetivo**: reducir magic numbers, strings repetidos y heuristicas opacas dentro del loop tactico.
- **Resultado esperado**: ownership claro para constantes de reciclado, recovery, agotamiento, timing de cast y scoring de spawn.
- **Dependencias**: `BL-020`, `BL-021`, `BL-022`
- **Prioridad**: Media

### BL-032 Automatizar pipeline tecnico de calidad
- **Estado**: `pending`
- **Tipo**: Infra + Delivery
- **Objetivo**: asegurar que test, lint y checks basicos de integracion no dependan solo del entorno local del desarrollador.
- **Resultado esperado**: automatizacion reproducible para `go test`, `golangci-lint` y checks clave del repositorio.
- **Dependencias**: `BL-018`
- **Prioridad**: Media

### BL-046 Externalizar catalogo de nodos y rutas del MVP
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: mover a datos la estructura minima de nodos y rutas del MVP de run una vez exista el flujo base de expedicion.
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
- **Resultado esperado**: tablas data-driven de recompensas, payout editorial y variaciones por nodo u outcome.
- **Dependencias**: `BL-037`, `BL-012`, `BL-040`
- **Prioridad**: Media

### BL-049 Externalizar catalogo de patrocinadores y ofertas
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: mover a datos las ofertas, efectos y compatibilidades basicas de patrocinadores cuando el sistema ya este validado como parte de run.
- **Resultado esperado**: catalogo data-driven de patrocinadores y reglas de oferta, separado del runtime de encounter y de la UI.
- **Dependencias**: `BL-039`, `BL-040`
- **Prioridad**: Baja

## Orden sugerido del trabajo pendiente

1. `BL-001`
2. `BL-034`
3. `BL-035`
4. `BL-036`
5. `BL-037`
6. `BL-038`
7. `BL-039`
8. `BL-040`
9. `BL-044`
10. `BL-002`
11. `BL-003`
12. `BL-011`
13. `BL-008`
14. `BL-012`
15. `BL-046`
16. `BL-047`
17. `BL-048`
18. `BL-042`
19. `BL-043`
20. `BL-017`
21. `BL-013`
22. `BL-016`
23. `BL-010`
24. `BL-009`
25. `BL-049`
26. `BL-014`
27. `BL-023`
28. `BL-031`
29. `BL-032`
