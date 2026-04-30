# Roadmap roguelike

Este documento concentra el backlog activo del proyecto, con estado visible para distinguir que esta pendiente, que ya esta planificado y que ya quedo integrado en `main`.

## Convencion de estado

- `pending`: item identificado pero aun sin plan activo en `docs/features/`.
- `planned`: item ya convertido en plan de feature, pero todavia no mergeado en `main`.
- `done`: item ya integrado en `main`.
- `cancelled`: item descartado, absorbido por otro o fuera de foco.

## Foto actual

- `done`: `BL-005`, `BL-006`, `BL-018`, `BL-019`, `BL-020`, `BL-021`
- `planned`: `BL-022`
- `pending`: resto del roadmap
- Foco recomendado inmediato: cerrar la fundacion de la expedicion con `BL-001`, `BL-011` y `BL-002`; si seguimos en la capa tactica, el siguiente bloque natural es `BL-022`.

## Foco sugerido actual

- `BL-001`: fijar el loop completo de la expedicion y sus capas de persistencia.
- `BL-011`: formalizar la economia de run y la frontera clara entre recursos de expedicion y progreso meta.
- `BL-002`: traducir ese loop a un mapa de zonas y nodos con informacion parcial.

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
  - El dinero resultante se gasta en deckbuilding fino, reparaciones y mejoras de la cana y sus aditamentos.
- **Reglas del hilo**:
  - El hilo representa la supervivencia global de la expedicion.
  - Su desgaste se arrastra durante toda la run y no se restaura automaticamente al cambiar de zona.
  - Las pescas fallidas, escapes violentos y ciertos eventos pueden romper tramos del hilo.
  - La tienda y algunas recompensas pueden reparar, reforzar o ampliar la tolerancia del hilo.
- **Taxonomia inicial de patrocinadores**:
  - `Revista de Trofeos`: premia capturas raras, elites, bosses y fotos de portada.
  - `Marca Tecnica`: mejora consistencia tactica, precision, thresholds y control del duelo.
  - `Taller Profesional`: favorece reparacion, mantenimiento, aguante del hilo y mejoras de la cana.
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
- **Objetivo**: extender la primera base configurable de arquetipos hacia perfiles de pez que soporten encuentros normales, elites, bosses de zona y legendarios finales dentro de la expedicion, concentrando tambien la metadata editorial o fotografica propia de cada especie o encuentro.
- **Resultado esperado**: estructura ampliable para mazo, efectos, metadata de descubrimiento, tags editoriales de foto y condiciones de aparicion por zona o rol de encounter, dejando claro que la capa fotografica vive en el contenido data-driven del pez y no en la resolucion generica del cast.
- **Dependencias**: `BL-006`, `BL-003`
- **Prioridad**: Media

### BL-022 Definir aparicion de peces por aguas base y ventana de lanzamiento
- **Estado**: `planned`
- **Tipo**: Discovery
- **Objetivo**: decidir como la combinacion entre aguas base del nodo, franja horizontal del cast y sesgo vertical del setup del jugador selecciona subconjuntos de peces y define la aparicion de encuentros compatibles con esa ventana de lanzamiento.
- **Resultado esperado**: modelo de pool base por nodo o zona, reglas de particion en subconjuntos por distancia y profundidad, metadata minima de aparicion en perfiles de pez y criterio para conectar cast, `rod`, aditamentos y habitats sin abrir todavia la capa de economia fotografica.
- **Dependencias**: `BL-003`, `BL-007`, `BL-020`, `BL-021`
- **Plan relacionado**: `docs/features/018-aparicion-de-peces-por-aguas-y-ventana-de-lanzamiento.md`
- **Direccion actual acordada**:
  - Cada nodo de pesca parte de unas aguas base y luego se secciona en subconjuntos de pool segun la ventana horizontal del cast y la ventana vertical habilitada por la `rod` y sus aditamentos.
  - La aparicion de peces debe resolverse a partir de la apertura ya cerrada del encounter (`InitialDistance`, `InitialDepth`) y no directamente desde los limites de escape del tablero.
  - La aparicion de peces debe seguir siendo compatible con zonas, elites, bosses y roles especiales de encounter.
  - El sistema debe permitir que a veces el jugador tenga lectura explicita al entrar al nodo y a veces no, sin exigir que el mapa revele antes la estructura del subpool.
  - La capa de economia o tags editoriales de fotografia queda fuera de este item y se aborda dentro de la extension data-driven del pez.
- **Prioridad**: Alta

## Items y Build

### BL-008 Definir categorias de objetos del jugador
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: clasificar las piezas de build que el jugador puede comprar, recibir o mejorar durante la expedicion sin solaparlas con los patrocinadores globales.
- **Resultado esperado**: taxonomia inicial de acciones y objetos de tienda que cubra mejoras de la `rod`, aditamentos, reparacion del hilo, intervenciones sobre el mazo, consumibles y otros ajustes de servicio.
- **Dependencias**: `BL-001`, `BL-011`
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

### BL-023 Desacoplar flujos de interfaz y preparar arquitectura UI-agnostic
- **Estado**: `pending`
- **Tipo**: Discovery + Delivery
- **Objetivo**: refactorizar los flujos de setup, opening, presentacion y sesion para que la aplicacion deje de depender estructuralmente del CLI como unica forma de presentar e interactuar con el juego, dejando a `internal/cli/` como un adaptador de borde y preparando una futura UI grafica para consumir los mismos casos de uso.
- **Resultado esperado**: arquitectura donde el combate, el setup previo al encounter y la apertura del lance se expresen como casos de uso UI-agnostic; presenter y view models mas semanticos para setup/opening/combate; controlador reusable del cast fuera del CLI; y `cmd/` reducido a composicion de dependencias, de modo que una GUI futura no tenga que duplicar wiring ni parsear strings finales para reconstruir la experiencia.
- **Dependencias**: `BL-018`, `BL-020`, `BL-021`
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
2. `BL-011`
3. `BL-002`
4. `BL-008`
5. `BL-003`
6. `BL-022`
7. `BL-012`
8. `BL-007`
9. `BL-009`
10. `BL-017`
11. `BL-023`
12. `BL-013`
13. `BL-015`
14. `BL-016`
15. `BL-010`
16. `BL-014`
