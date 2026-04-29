# Plan de feature: visibilidad-descarte-del-pez

## Objetivo

Convertir el descarte del pez en una fuente de informacion estrategica real dentro del duelo, de forma que el jugador pueda leer mejor el patron del pez sin perder tension ni espacio para excepciones.

La meta de esta feature no es mostrar "mas datos" por mostrar. Es mover el combate desde una opacidad excesiva hacia una lectura tactica mas interesante: el jugador debe poder ver que cartas ya uso el pez, entender cuando esas cartas aun no pueden volver, y detectar cuando un reciclado o una regla del pez cambia de nuevo el espacio de posibilidades.

Al mismo tiempo, la feature debe dejar preparada una superficie clara para que algunas cartas, arquetipos o eventos reduzcan esa informacion de forma parcial y legible. El objetivo no es introducir una IA mas tramposa, sino una capa de intel del pez mas rica, con visibilidad por defecto y excepciones controladas.

## Criterios de aceptacion

- El runtime del mazo del pez expone un snapshot legible de descarte e informacion de ciclo, sin obligar a reconstruir esa informacion desde la UI.
- `match.State`, `presentation` y `cli` pueden mostrar el descarte visible del pez durante el encounter.
- La lectura por defecto prioriza el ciclo actual del mazo y comunica de forma explicita cuando el pez recicla, si rebaraja y cuantas cartas retira por ciclo.
- El modelo soporta al menos cuatro niveles de lectura de una carta descartada: `full`, `move_only`, `masked` y `hidden`.
- La implementacion incluye al menos un camino real de validacion para lectura reducida, ya sea mediante un preset de prueba, una carta del pez o un arquetipo que no revele toda la informacion del historial.
- La CLI sigue siendo legible en terminal sin convertirse en una vista grande de mano o debug.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope

### Incluye

- Extender el mazo del pez para producir un snapshot de descarte visible, separado entre ciclo actual y resumen de ciclos previos.
- Modelar la visibilidad del historial como una capacidad de runtime reutilizable y no como logica hardcodeada en la CLI.
- Transportar esa informacion desde `deck` hasta `match`, `presentation` y `cli`.
- Mejorar la UI CLI para mostrar el historial visible del pez, el estado del reciclado y el caracter del mazo (`shuffle` y `CardsToRemove`).
- Introducir una primera forma concreta de lectura parcial u opaca para validar que el sistema no solo sirve para el caso completamente visible.
- Actualizar documentacion afectada para explicar la nueva lectura del duelo.

### No incluye

- Implementar aun bestiario persistente, memoria entre runs o desbloqueos meta ligados al descarte del pez.
- Redisenar toda la interfaz del combate o introducir una vista compleja de cartas detalladas.
- Crear todavia un sistema completo de niebla informativa para todos los sistemas del juego.
- Cambiar el balance general del combate mas alla de lo necesario para soportar y validar la nueva lectura.

## Propuesta de implementacion

### 1. Modelo de informacion visible del mazo del pez

La informacion del descarte debe nacer en el runtime del mazo, no en `presentation` ni en `cli`.

Propuesta base:

- Introducir un snapshot de solo lectura en `internal/deck/` con tipos equivalentes a:
  - `VisibleDiscardEntry`
  - `VisibleDiscardCycleSummary`
  - `DeckVisibilitySnapshot`
- El snapshot debe incluir al menos:
  - descarte visible del ciclo actual
  - resumen de ciclos anteriores
  - `RecycleCount`
  - si el mazo rebaraja al reciclar
  - cuantas cartas retira el reciclado
  - si el mazo esta agotado
- El snapshot no debe exponer slices mutables internos del mazo.

La idea es que `deck` siga siendo el source of truth tanto para conteos como para lectura visible del historial.

### 2. Niveles de visibilidad por carta descartada

Para que la feature no quede cerrada al caso "todo visible siempre", conviene modelar la lectura por entrada de historial.

Propuesta:

- Anadir a `cards.FishCard` una metadata de historial, por ejemplo `DiscardVisibility`.
- Niveles iniciales:
  - `full`: muestra nombre si existe y, en su defecto, movimiento legible
  - `move_only`: muestra solo el movimiento del pez
  - `masked`: muestra que hubo una carta pero no cual
  - `hidden`: no agrega entrada visible al historial
- Valor por defecto para cartas existentes: `full`.

Esto deja una superficie simple y extensible para cartas futuras sin acoplar la excepcion a la UI.

### 3. Limites de lectura recomendados

La informacion de mayor valor estrategico es el descarte del ciclo actual. Por eso la feature debe priorizar esa lectura y no intentar mostrar un historial largo sin jerarquia.

Regla propuesta:

- `ciclo actual`: mostrar entradas visibles en orden de uso.
- `ciclos previos`: mostrar resumen compacto, no lista completa detallada.
- Si el pez recicla con `shuffle: true`, la UI debe indicarlo claramente para que el jugador entienda que las cartas pueden volver.
- Si el pez recicla con `shuffle: false`, la UI debe comunicar que el patron sigue siendo mas predecible.
- Si `CardsToRemove > 0`, la UI debe indicar que parte del descarte queda fuera del siguiente ciclo.

Esta decision busca maximizar claridad sin inflar la interfaz.

### 4. Integracion con `match`, `presentation` y `cli`

#### 4.1 `internal/game/` y `internal/match/`

- Extender la interfaz `FishDeck` usada por `internal/game/engine.go` para poder pedir el snapshot visible del mazo.
- Ampliar `match.DeckState` con una representacion del historial visible y metadata de reciclado, en lugar de quedarse solo con conteos.
- Hacer que `engine.refreshState()` copie el snapshot del mazo al estado compartido.

#### 4.2 `internal/presentation/`

- Crear view models especificos para el historial del pez, por ejemplo:
  - `FishDiscardEntryView`
  - `FishDiscardView`
  - `FishDiscardCycleSummaryView`
- Traducir movimientos, nombres ocultos y estados de reciclado a etiquetas legibles y compactas.
- Conservar el tono actual de `Presenter`: informacion funcional, breve y coherente con la terminal.

#### 4.3 `internal/cli/`

- Anadir una nueva seccion de render tipo `Historial del pez` en la pantalla principal del duelo.
- Formato recomendado:
  - `Ciclo actual: Embestir | ? | Zafarse`
  - `Reciclado: baraja / orden fijo | retira 3 | ciclos 1`
  - `Ciclos previos: 1 ciclo resumido`
- Si no hay historial aun, mostrar una linea minima del estilo `Sin cartas usadas todavia.`
- Si una carta es `masked`, renderizar `?`.
- Si es `move_only`, renderizar solo el label del movimiento del pez.
- Si es `hidden`, no agregarla al listado visible, pero el resumen de ciclo debe seguir siendo coherente con los conteos globales.

### 5. Primera validacion real de lectura parcial

La feature no deberia cerrar dejando la modularidad solo "lista para mas adelante". Conviene validar al menos un caso real de informacion parcial.

Opcion recomendada:

- Anadir un preset de pez de prueba o extender uno existente con una carta cuya visibilidad sea `masked` o `move_only`.

Condicion importante:

- La excepcion debe ser legible y acotada. No conviene que el primer ejemplo oculte todo el historial del encuentro.

Esto permite probar que el sistema sirve tanto para el caso base como para variaciones futuras de arquetipo o evento.

### 6. Cobertura automatizada

La feature deberia introducir tests en cuatro niveles:

- `internal/deck/`
  - snapshot del descarte visible antes y despues de robar
  - separacion entre ciclo actual y ciclos previos tras reciclar
  - tratamiento correcto de `full`, `move_only`, `masked` y `hidden`
- `internal/game/`
  - refresco del estado compartido con el snapshot visible del mazo
- `internal/presentation/`
  - traduccion correcta del historial a views legibles
- `internal/cli/`
  - render compacto del historial visible y del estado de reciclado

### 7. Documentacion a actualizar

- `README.md` debe mencionar que el duelo muestra historial visible del pez como parte del loop actual.
- Si cambia la lectura de presets o de UX, ajustar tambien la documentacion del backlog o del plan segun haga falta.

## Validacion

- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Validar manualmente desde CLI un encounter con mazo barajado y otro con orden fijo para comprobar que el historial visible ayuda a leer patrones.
- Validar manualmente un caso con reciclado para verificar que la UI comunica bien que las cartas vuelven o dejan de volver.
- Validar manualmente al menos un caso con una carta `masked` o `move_only` para comprobar que la excepcion sigue siendo entendible.

## Riesgos o decisiones abiertas

- Si mostramos demasiado detalle de ciclos anteriores, la interfaz puede ganar ruido y perder valor tactico.
- Si usamos `hidden` con demasiada frecuencia, el sistema puede volver a sentirse opaco en vez de estrategico.
- Habra que decidir si el historial visible muestra nombre de carta, movimiento o ambas cosas cuando una carta tiene identidad fuerte y texto largo.
- Hay que cuidar la coherencia entre lo que el historial oculta y los conteos visibles del mazo, para que el jugador no sienta contradicciones.
