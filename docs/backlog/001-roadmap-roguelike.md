# Roadmap roguelike

Este documento concentra el backlog activo del proyecto, con estado visible para distinguir que esta pendiente, que ya esta planificado y que ya quedo integrado en `main`.

## Convencion de estado

- `pending`: item identificado pero aun sin plan activo en `docs/features/`.
- `planned`: item ya convertido en plan de feature, pero todavia no mergeado en `main`.
- `done`: item ya integrado en `main`.
- `cancelled`: item descartado, absorbido por otro o fuera de foco.

## Foto actual

- `done`: `BL-005`, `BL-006`, `BL-018`
- `planned`: ninguno
- `pending`: resto del roadmap
- Foco recomendado inmediato: cerrar la fundacion de la expedicion con `BL-001`, `BL-002` y `BL-011` antes de abrir mas delivery transversal.

## Foco sugerido actual

- `BL-001`: fijar el loop completo de la expedicion y sus capas de persistencia.
- `BL-002`: traducir ese loop a un mapa de zonas y nodos con informacion parcial.
- `BL-011`: cerrar la economia editorial, la reserva de fotos y la frontera entre run y meta.

## Core Loop

### BL-001 Definir loop completo de una run
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: cerrar como inicia, progresa y termina una expedicion de pesca deportiva roguelike apoyada en encuentros de pesca, ruta por nodos, economia editorial y build persistente durante la run.
- **Resultado esperado**: documento de flujo end-to-end de partida con capas de persistencia, estructura de zona, taxonomia inicial de patrocinadores y reglas base de economia de run.
- **Direccion actual acordada**:
  - La run es una expedicion de pesca deportiva que empieza en aguas mas accesibles y avanza hacia mar abierto y peces de leyenda.
  - El objetivo de la run es llegar vivo y suficientemente fuerte al final de la expedicion para capturar un pez legendario final.
  - La derrota de la run ocurre cuando el jugador agota todo el hilo de la cana.
  - El retiro voluntario puede ocurrir solo en checkpoints o servicios y cierra la expedicion sin acceder al legendario final.
  - Cada zona tiene tiempo limitado, pero el tiempo funciona como presion local de zona y no como fail state global de la expedicion.
  - Entre zonas persisten la build, la economia, los patrocinadores, el desgaste del hilo y cualquier otro recurso de run.
  - Las zonas aumentan la dificultad y la cercania al final de la expedicion, pero no resetean la progresion acumulada en la run.
  - La composicion de la baraja del pez es desconocida la primera vez que se enfrenta; despues de capturarlo, esa informacion queda registrada y persiste entre runs.
  - Las fotos de peces legendarios se conservan entre runs solo como trofeos, no como moneda ni mejora mecanica.
- **Capas de persistencia**:
  - `por encuentro`: estado tactico del duelo, cartas vistas, outcome y dano puntual del enfrentamiento.
  - `por run`: hilo restante, build del jugador, dinero, patrocinadores activos, fotos reservadas y conocimiento descubierto durante esa expedicion.
  - `entre runs`: registro de peces capturados, informacion desbloqueada de sus barajas y trofeos fotograficos de peces legendarios.
- **Loop base de la run**:
  - Preparar la expedicion con un setup inicial del jugador.
  - Entrar en una zona con tiempo limitado y avanzar por nodos con pistas tematicas y traduccion mecanica consistente.
  - Resolver nodos de pesca, servicio, evento o riesgo hasta alcanzar el cierre de zona.
  - Obtener fotos al capturar peces, venderlas automaticamente al llegar a servicios y convertir ese dinero en mejoras de tienda.
  - Elegir mejoras de patrocinador en hitos concretos para definir la identidad global de la run.
  - Mantener el hilo durante toda la expedicion, reparar o mejorar el equipo en servicios y llegar al legendario final.
  - Ganar la run al capturar el pez legendario final; perderla al agotar el hilo antes de conseguirlo.
- **Economia editorial de la run**:
  - Cada pesca exitosa genera fotos con valor economico y tags editoriales.
  - Las fotos se venden automaticamente al llegar a muelles, talleres o checkpoints de servicio.
  - El jugador puede reservar un numero pequeno de fotos para esperar sinergias de patrocinador o eventos que las paguen mejor mas adelante.
  - Las fotos reservadas no generan interes fijo; su revalorizacion depende de patrocinadores y eventos editoriales concretos.
  - Al cerrar una zona, toda foto reservada se liquida automaticamente antes de avanzar.
  - El dinero resultante se gasta en deckbuilding fino, reparaciones y mejoras del rig.
- **Reglas del hilo**:
  - El hilo representa la supervivencia global de la expedicion.
  - Su desgaste se arrastra durante toda la run y no se restaura automaticamente al cambiar de zona.
  - Las pescas fallidas, escapes violentos y ciertos eventos pueden romper tramos del hilo.
  - La tienda y algunas recompensas pueden reparar, reforzar o ampliar la tolerancia del hilo.
- **Taxonomia inicial de patrocinadores**:
  - `Revista de Trofeos`: premia capturas raras, elites, bosses y fotos de portada.
  - `Marca Tecnica`: mejora consistencia tactica, precision, thresholds y control del duelo.
  - `Taller Profesional`: favorece reparacion, mantenimiento, aguante del hilo y mejoras del rig.
  - `Patrocinio Offshore`: premia pesca profunda, mar abierto, encuentros duros y riesgo alto.
  - `Editorial de Aventura`: favorece reserva de fotos, eventos especiales y rutas arriesgadas.
  - `Canal Popular`: premia volumen, fotos frecuentes y rentabilidad estable de la expedicion.
- **Regla de oferta de patrocinadores**:
  - La run usa un numero acotado y tematico de patrocinadores para que la identidad sea legible.
  - Escoger una mejora de un patrocinador no garantiza volver a verlo, pero sesga la siguiente oferta hacia ese patrocinador y patrocinadores afines.
  - La pool siguiente mantiene elegible al patrocinador escogido, reduce temporalmente patrocinadores incompatibles y favorece que la run se enfoque en una o dos lineas principales.
  - Los patrocinadores actuan como efectos globales de run y no como cartas concretas del mazo.
- **Flujo exacto de una zona en el MVP**:
  - `Entrada de zona`: presenta tono, clima, dificultad y primeras pistas de ruta; no resetea la build.
  - `Tramo inicial`: 2 nodos resolubles con informacion parcial, normalmente pesca, evento o riesgo moderado.
  - `Bifurcacion`: al menos una decision entre dos rutas con pistas tematicas, sin revelar el tipo exacto del siguiente nodo.
  - `Servicio garantizado`: 1 nodo de muelle, taller o prensa donde se venden fotos no reservadas, se permite reservar fotos y se gasta dinero en tienda.
  - `Tramo final`: 1 o 2 nodos de mayor presion, con posibilidad de elite, evento fuerte o pesca exigente.
  - `Cierre de zona`: boss o objetivo fuerte de zona, liquidacion forzada de fotos reservadas y transicion a la siguiente zona.
  - `Hito de progresion`: al cierre de zona se ofrece una mejora de patrocinador antes de entrar en la siguiente etapa de la expedicion.
- **Taxonomia minima de nodos para el MVP**:
  - `Pesca`: encounter normal que genera fotos si se captura al pez.
  - `Elite`: encounter mas dificil con mejor recompensa economica o editorial.
  - `Evento`: decision o situacion especial que altera economia, ruta o build.
  - `Servicio`: muelle, taller o punto editorial que concentra venta automatica y tienda.
  - `Boss de zona`: cierre mecanico y tematico de cada tramo de la expedicion.
  - `Boss final legendario`: cierre de run y condicion principal de victoria.
- **Criterios de cierre**:
  - queda definido el inicio, avance por zonas, derrota, victoria y retiro
  - queda claro que recompensas son de run y cuales son meta
  - queda definido que persiste por encuentro, por run y entre runs
  - queda definida la economia `foto -> dinero -> tienda` y el rol de la reserva
  - queda definida la taxonomia inicial de patrocinadores y su regla de oferta ponderada
  - queda definido el flujo exacto de una zona del MVP y su taxonomia minima de nodos
- **Prioridad**: Alta

### BL-002 Definir estructura de mapa y tipos de nodo
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: decidir como se representa la expedicion por zonas, cuantas bifurcaciones reales existen y como funciona la informacion parcial de la ruta sin perder legibilidad.
- **Resultado esperado**: modelo de mapa por zona con longitud objetivo, puntos de convergencia, familias de nodos, vocabulario de pistas tematicas y reglas de ocultacion del siguiente nodo.
- **Dependencias**: `BL-001`
- **Prioridad**: Alta

### BL-003 Disenar progresion de dificultad entre zonas
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: decidir como escala el reto desde la costa inicial hasta mar abierto y los peces legendarios finales sin resetear la build ni la economia de run.
- **Resultado esperado**: reglas de progresion por zona para pools de peces, riesgo de nodos, presion de tiempo, desgaste del hilo, frecuencia de servicios y preparacion del boss final.
- **Dependencias**: `BL-001`, `BL-002`
- **Prioridad**: Alta

## Fish y Encounters

### BL-007 Expandir perfiles data-driven de pez
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: extender la primera base configurable de arquetipos hacia perfiles de pez que soporten encuentros normales, elites, bosses de zona y legendarios finales dentro de la expedicion.
- **Resultado esperado**: estructura ampliable para mazo, efectos, metadata de descubrimiento, tags editoriales de foto y condiciones de aparicion por zona o rol de encounter.
- **Dependencias**: `BL-006`, `BL-003`
- **Prioridad**: Media

## Items y Build

### BL-008 Definir categorias de objetos del jugador
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: clasificar las piezas de build que el jugador puede comprar, recibir o mejorar durante la expedicion sin solaparlas con los patrocinadores globales.
- **Resultado esperado**: taxonomia inicial de acciones y objetos de tienda que cubra mejoras del rig, reparacion del hilo, intervenciones sobre el mazo, consumibles y otros ajustes de servicio.
- **Dependencias**: `BL-001`, `BL-011`
- **Prioridad**: Alta

### BL-009 Disenar sistema de sinergias
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: definir como interactuan cartas del jugador, mejoras de rig, hilo, patrocinadores, fotos reservadas y economia editorial sin generar combinaciones opacas o imposibles de leer.
- **Resultado esperado**: reglas de stacking, limites, superficies de activacion y familias de build que mantengan el juego expresivo y entendible.
- **Dependencias**: `BL-005`, `BL-008`, `BL-011`
- **Prioridad**: Media

### BL-010 Implementar primer vertical slice de build
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: probar una expedicion completa con una version reducida pero funcional del loop de build, economia y progresion hacia un legendario final.
- **Resultado esperado**: prototipo jugable con mapa de zona, servicios, patrocinadores, economia de fotos, acciones de tienda y un conjunto pequeno de mejoras y bosses finales.
- **Dependencias**: `BL-001`, `BL-002`, `BL-003`, `BL-007`, `BL-008`, `BL-009`, `BL-011`, `BL-012`
- **Prioridad**: Media

## Economy y Meta

### BL-011 Definir economia de run y meta-progresion
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: formalizar la economia editorial de la expedicion y separar con claridad que recursos pertenecen a la run y cuales quedan entre partidas como coleccion o conocimiento.
- **Resultado esperado**: modelo economico con fotos, dinero, reserva, venta automatica, sinks de tienda, valor editorial y frontera explicita entre progreso de run y meta.
- **Dependencias**: `BL-001`
- **Prioridad**: Alta

### BL-012 Disenar sistema de recompensas entre encuentros
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: decidir como se distribuyen fotos, dinero, ofertas de patrocinador, acceso a servicio, reparaciones y otras recompensas entre nodos y entre zonas.
- **Resultado esperado**: tabla de recompensas por tipo de nodo, outcome, zona y rol del pez para sostener el ritmo de la expedicion sin inflacion de recursos.
- **Dependencias**: `BL-002`, `BL-008`, `BL-011`
- **Prioridad**: Alta

## Collection

### BL-013 Definir bestiario y coleccion
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: concretar como se registra entre runs la informacion descubierta de cada pez, que muestran los trofeos fotografiados y como se presenta el bestiario al jugador.
- **Resultado esperado**: definicion de registro de especies, composicion descubierta de sus barajas, categorias de legendarios y presentacion de trofeos sin convertir la coleccion en una fuente de poder.
- **Dependencias**: `BL-001`
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
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: preparar peces, patrocinadores, zonas, nodos, eventos y catalogos de tienda para crecer sin depender de wiring duro en codigo.
- **Resultado esperado**: estrategia de datos y versionado de contenido para perfiles de pez, tablas de ruta, ofertas de patrocinador, recompensas y servicios.
- **Dependencias**: `BL-002`, `BL-007`, `BL-008`, `BL-011`
- **Prioridad**: Media

### BL-016 Disenar guardado de run y progreso meta
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: decidir como persistir expediciones en curso y datos meta como bestiario, trofeos legendarios y conocimiento ya descubierto del pez.
- **Resultado esperado**: contrato de persistencia para run activa y perfil permanente, con limites de versionado y reglas de compatibilidad.
- **Dependencias**: `BL-001`, `BL-011`, `BL-013`
- **Prioridad**: Media

### BL-017 Mejorar UX de lectura de build y estado de run
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: asegurar que el mapa, el hilo, los patrocinadores, la reserva de fotos, la economia y la build del jugador sigan siendo legibles durante la expedicion.
- **Resultado esperado**: requerimientos de HUD, feedback de ruta, lectura de servicios, valuacion editorial de fotos, resumen de build y transparencia de efectos globales.
- **Dependencias**: `BL-001`, `BL-008`, `BL-011`
- **Prioridad**: Media

### BL-019 Hacer visible el descarte del pez y modular la lectura del historial
- **Estado**: `pending`
- **Tipo**: Discovery + Delivery
- **Objetivo**: convertir el descarte del pez en una herramienta estrategica legible durante el encounter, manteniendo espacio para peces, cartas o eventos que oculten parcial o temporalmente esa informacion.
- **Resultado esperado**: estado de runtime, presentacion y UX que permitan ver el descarte visible del pez por ciclo, entender cuando el mazo recicla o se baraja y soportar excepciones de visibilidad por carta, arquetipo o evento.
- **Dependencias**: `BL-005`, `BL-006`, `BL-017`
- **Direccion actual acordada**:
  - La regla general del combate pasa a ser que el jugador puede consultar el historial de descarte visible del pez durante el encuentro.
  - La lectura prioritaria debe centrarse en el ciclo actual del mazo; los ciclos anteriores pueden quedar resumidos, atenuados o separados si eso mejora claridad.
  - El reciclado del mazo debe comunicarse de forma explicita, incluyendo cuando el pez rebaraja y cuando parte del descarte ya no puede volver por `CardsToRemove`.
  - Algunas cartas, arquetipos o efectos concretos pueden ocultar una entrada individual del historial, mostrar solo el movimiento o volver opaco parte del descarte, pero no conviene que la excepcion por defecto oculte todo el panel.
  - La informacion visible del descarte debe nacer en el runtime del mazo y llegar a `match`, `presentation` y `cli` como un snapshot legible, en lugar de reconstruirse ad hoc en la UI.
  - La feature debe dejar preparada una superficie clara para futuras mecanicas de intel del pez, niebla informativa, bosses con fases opacas y desbloqueos de bestiario entre runs.
- **Prioridad**: Media

## Completados

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
2. `BL-002`
3. `BL-011`
4. `BL-003`
5. `BL-008`
6. `BL-012`
7. `BL-007`
8. `BL-009`
9. `BL-019`
10. `BL-017`
11. `BL-013`
12. `BL-015`
13. `BL-016`
14. `BL-010`
15. `BL-014`
