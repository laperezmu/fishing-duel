# Backlog interfaz grafica y prototipo visual

Este documento concentra la nueva linea de trabajo para pasar del prototipo CLI actual a un primer videojuego jugable tambien a nivel visual. Su objetivo no es reemplazar de inmediato la arquitectura de run ya existente, sino abrir una pista separada de analisis, prototipado y delivery para una interfaz grafica real.

## Convencion de estado

- `pending`: item identificado pero aun sin plan activo en `docs/features/`.
- `planned`: item ya convertido en plan de feature, pero todavia no mergeado en `main`.
- `cancelled`: item descartado, absorbido por otro o fuera de foco.

## Convencion de identificadores

- Este backlog usa su propia notacion independiente: `GUI-###`.
- La numeracion del scope grafico no depende del backlog roguelike principal.
- Los futuros planes de feature de esta linea deberian reflejar tambien esta separacion conceptual.

## Foto actual

- `pending`: todos los items de esta linea
- Recomendacion inicial: usar `Ebitengine` como libreria objetivo del primer prototipo grafico.
- Alternativas consideradas: `raylib-go`, `Fyne`

## Analisis inicial de librerias

### Opcion recomendada: `Ebitengine`

Por que encaja mejor con el proyecto ahora:

- esta pensada especificamente para juegos 2D en Go, no para apps de escritorio genericas
- su modelo `update/draw` encaja bien con la capa de runtime y estados que el repo ya tiene
- soporta desktop, mobile y WebAssembly, lo que deja una salida razonable para prototipos distribuibles
- tiene utilidades nativas para input, audio, texto, sprites, shaders y rendering 2D sin meter una pila externa enorme
- es una opcion bastante adoptada en el ecosistema Go y tiene comunidad/documentacion solidas

Limitantes a asumir desde el principio:

- no es un toolkit de widgets tradicional; la UI del juego hay que construirla como UI de juego
- fuera de Windows, el cross-compiling sigue teniendo limites practicos por dependencias de plataforma
- si mas adelante quisieramos tooling visual tipo editor o pantallas muy de escritorio, haria falta otra capa o herramientas separadas
- obliga a pensar desde el principio en escenas, camara, loop de render y assets, no solo en layouts declarativos

### Opcion alternativa: `raylib-go`

Por que podria servir:

- API muy directa y expresiva para prototipos visuales rapidos
- buena ergonomia para rendering 2D, audio y primitives
- tambien esta orientada a videojuegos mas que a apps de escritorio

Principales limitantes frente a `Ebitengine`:

- depende mas claramente de toolchains y requisitos nativos, especialmente con `cgo`
- el primer build y la distribucion suelen tener algo mas de friccion operativa
- su ecosistema Go para arquitectura idiomatica de juegos parece menos alineado con el objetivo de mantener el proyecto como monolito modular pequeno y portable

### Opcion descartable para este objetivo: `Fyne`

Por que no deberia ser la opcion principal para este reto:

- esta pensada para apps de escritorio nativas, no para un videojuego 2D con loop de render protagonista
- resuelve muy bien widgets, formularios y layouts, pero menos naturalmente escenas de juego, tableros, animacion continua y direccion audiovisual de videojuego
- serviria mejor para tooling, launcher o utilidades futuras que para el juego principal

## Direccion recomendada

- Mantener el dominio, la run y los contratos de `internal/app/` como fuente de verdad.
- Introducir una nueva capa de adaptador grafico separada del CLI.
- Construir un primer vertical slice visual sobre `Ebitengine` sin romper `cmd/fishing-run/` ni `cmd/fishing-duel/`.
- Tratar el primer prototipo grafico como otro composition root, no como un reemplazo abrupto del proyecto actual.

## Tareas

### GUI-001 Evaluar `Ebitengine` como base del cliente grafico
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: confirmar que `Ebitengine` cubre suficientemente bien el primer prototipo visual del juego en desktop y definir sus restricciones reales para este repo.
- **Resultado esperado**: decision tecnica cerrada sobre la libreria grafica principal, con tradeoffs documentados frente a `raylib-go` y descarte explicito de toolkits tipo `Fyne` para el cliente principal.
- **Dependencias**: ninguna
- **Prioridad**: Alta

### GUI-002 Disenar arquitectura de adaptador grafico sobre la app actual
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: decidir como se conecta un cliente grafico al runtime existente sin mezclar dominio, run y rendering dentro de un solo paquete acoplado.
- **Resultado esperado**: propuesta de `cmd/` grafico y adaptadores visuales sobre `internal/app/`, `internal/presentation/` y futuros view models de escena.
- **Dependencias**: `GUI-001`
- **Prioridad**: Alta

### GUI-003 Implementar composition root de prototipo grafico
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: crear un nuevo ejecutable grafico separado del CLI actual para montar el primer prototipo visual del juego.
- **Resultado esperado**: nuevo `cmd/` para cliente grafico con loop basico, inicializacion de ventana, entrada a la run y borde listo para escenas.
- **Dependencias**: `GUI-001`, `GUI-002`
- **Prioridad**: Alta

### GUI-004 Construir sistema minimo de escenas y navegacion visual
- **Estado**: `pending`
- **Tipo**: Discovery + Delivery
- **Objetivo**: definir la estructura minima de escenas del prototipo grafico para pasar de menus/seleccion a run y encounter sin depender de pantallas CLI.
- **Resultado esperado**: escenas base para inicio, seleccion de pescador, nodo actual, encounter y cierre de run con transiciones simples.
- **Dependencias**: `GUI-002`, `GUI-003`
- **Prioridad**: Alta

### GUI-005 Definir direccion artistica minima del primer slice visual
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: fijar una identidad visual inicial suficientemente clara para que el prototipo no sea solo funcional sino tambien legible y con atmosfera.
- **Resultado esperado**: guia minima de color, tipografia, escala, framing, ritmo visual y tratamiento del agua, pez, tablero y HUD.
- **Dependencias**: `GUI-001`
- **Prioridad**: Alta

### GUI-006 Implementar HUD y tablero visual del encounter
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: traducir el estado tactico actual del duel a una presentacion grafica jugable, legible y mas evocadora que la version terminal.
- **Resultado esperado**: tablero visual con distancia, profundidad, tension, cartas visibles, historial clave y feedback de ronda.
- **Dependencias**: `GUI-003`, `GUI-004`, `GUI-005`
- **Prioridad**: Alta

### GUI-007 Implementar presentacion visual de run, nodos y servicios
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: hacer visible la estructura de run actual en una interfaz grafica con lectura clara de nodo actual, hilo, capturas y progreso de expedicion.
- **Resultado esperado**: HUD o scene de run capaz de mostrar el recorrido minimo, los puntos de pesca, el servicio y el cierre de expedicion.
- **Dependencias**: `GUI-003`, `GUI-004`, `GUI-005`
- **Prioridad**: Alta

### GUI-008 Definir pipeline minimo de assets y recursos visuales
- **Estado**: `pending`
- **Tipo**: Discovery
- **Objetivo**: decidir como se organizan, cargan y empaquetan imagenes, tipografias, iconos, audio y placeholders del primer prototipo visual.
- **Resultado esperado**: convencion minima para assets y carga reproducible dentro del repo sin bloquear iteracion rapida.
- **Dependencias**: `GUI-001`, `GUI-003`, `GUI-005`
- **Prioridad**: Media

### GUI-009 Implementar animacion y feedback visual minimo del loop
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: agregar movimiento, transiciones y feedback suficiente para que el prototipo se sienta como juego y no como maqueta estatica.
- **Resultado esperado**: animaciones simples de entrada, cast, resolucion de ronda, captura, escape y transicion entre nodos.
- **Dependencias**: `GUI-006`, `GUI-007`, `GUI-008`
- **Prioridad**: Media

### GUI-010 Entregar primer vertical slice grafico jugable
- **Estado**: `pending`
- **Tipo**: Delivery
- **Objetivo**: integrar run minima, encounter visual, seleccion de pescador y presentacion audiovisual basica en un prototipo jugable de punta a punta.
- **Resultado esperado**: primer build grafico capaz de ejecutar una run simple completa con identidad visual propia, aunque todavia use placeholders y contenido minimo.
- **Dependencias**: `GUI-003`, `GUI-004`, `GUI-005`, `GUI-006`, `GUI-007`, `GUI-008`, `GUI-009`
- **Prioridad**: Alta
