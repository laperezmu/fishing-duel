# Feature Specification: Formalizar fishing-duel como sandbox de encounters

**Feature Branch**: `003-encounter-sandbox`  
**Created**: 2026-05-08  
**Status**: Draft  
**Input**: User description: "Revisa el Issue 64 de github, ahi se encuentra todo el contexto necesario para desarrollar la tarea."

## User Scenarios & Testing *(mandatory)*

### User Story 0 - Framing de comando (Priority: P0)

Como usuario existente, quiero que el comando `fishing-duel` mantenga su nombre para no romper scripts o integraciones, pero la UX muestre un framing actualizado que indique que es un sandbox de encounters para QA y testing.

**Rationale**: Mantener el comando `fishing-duel` evita cambios en scripts, integraciones CI y habitos existentes. El cambio es solo de framing (texto, ayuda, prompts).

**Acceptance Scenarios**:

1. **Given** un usuario que ejecuta `fishing-duel`, **When** ve la ayuda o los prompts, **Then** observa framing de "encounter sandbox" o similar en lugar de "fishing duel".
2. **Given** un script existente que llama a `fishing-duel`, **When** se ejecuta, **Then** funciona sin cambios.

### User Story 1 - Configurar encounters de forma explicita (Priority: P1)

Como disenador o tester, quiero abrir un sandbox de encounters con seleccion directa de presets, cartas concretas y parametros clave para reproducir combinaciones concretas sin depender de filtros implicitos del flujo actual.

**Why this priority**: Es el valor base del sandbox. Sin setup explicito, el binario sigue siendo una experiencia parcial y no una herramienta confiable de QA o exploracion.

**Independent Test**: Puede probarse iniciando el sandbox y configurando manualmente una partida con preset de jugador, preset de pez, cartas concretas del jugador y del pez, cana, aditamentos y contexto de agua definidos por el usuario, verificando que el encounter arranca con esos valores.

**Acceptance Scenarios**:

1. **Given** un usuario dentro del sandbox, **When** selecciona preset de jugador, preset de pez, cartas concretas, cana, aditamentos y contexto de agua, **Then** el encounter se inicializa con esa combinacion exacta.
2. **Given** un usuario que quiere repetir una prueba, **When** fija una seed reproducible, **Then** el sandbox vuelve a producir la misma configuracion aleatoria derivada bajo las mismas elecciones.
3. **Given** un usuario que no necesita control total, **When** elige el modo guiado, **Then** puede recorrer un flujo simplificado similar al actual sin perder compatibilidad con el sandbox nuevo.

---

### User Story 2 - Inspeccionar la resolucion real del runtime (Priority: P2)

Como tester o balanceador, quiero ver con claridad que triggers se activaron, que efectos se resolvieron y en que orden para entender por que un encounter produjo un resultado concreto.

**Why this priority**: El rediseño reciente del sistema de efectos hace que la inspeccion de runtime sea clave para debugging, balance y regresion. Sin esto, el sandbox pierde gran parte de su valor operativo.

**Independent Test**: Puede probarse ejecutando una ronda con efectos simultaneos y verificando que el sandbox muestra la traza ordenada de resolucion, prioridades, desempates y estado final de forma consistente con el runtime.

**Acceptance Scenarios**:

1. **Given** una ronda con varios efectos aplicables, **When** termina la resolucion, **Then** el sandbox muestra triggers activados, efectos resueltos y su orden de prioridad.
2. **Given** dos efectos con la misma prioridad, **When** se muestra la traza de resolucion, **Then** queda visible que el efecto del pez se resolvio antes que el del jugador.
3. **Given** un usuario que necesita mas detalle, **When** activa una vista mas verbosa o de depuracion, **Then** el sandbox muestra informacion adicional del estado antes y despues de resolver sin convertir la UI en fuente de verdad de las reglas.

---

### User Story 3 - Reproducir escenarios y ejecutar pruebas de QA con baja friccion (Priority: P3)

Como QA o desarrollador, quiero lanzar escenarios guardados con seed fija para repetir bugs, smoke tests y pruebas de regresion de manera rapida y compartible.

**Why this priority**: El sandbox gana mucho valor cuando puede reutilizarse fuera del flujo manual. Esto habilita pruebas repetibles, colaboracion y validacion sistematica de casos complejos.

**Independent Test**: Puede probarse ejecutando un escenario guardado con seed fija y verificando que el sandbox levanta el encounter esperado y produce una salida reproducible y util para inspeccion.

**Acceptance Scenarios**:

1. **Given** un escenario predefinido de QA, **When** el usuario lo ejecuta, **Then** el sandbox aplica la configuracion completa sin requerir seleccion manual adicional.
2. **Given** un escenario compartido entre miembros del equipo, **When** dos usuarios lo ejecutan con la misma seed, **Then** ambos obtienen el mismo comportamiento observable del encounter.
3. **Given** una prueba de regresion sobre triggers, splash o agotamiento, **When** el usuario exporta o consulta la salida del sandbox, **Then** dispone de evidencia suficiente para comparar resultados entre ejecuciones.

---

### Edge Cases

- Que ocurre cuando el usuario elige una combinacion manual incompatible entre contexto de agua, habitats y preset de pez.
- Como responde el sandbox cuando una override manual contradice un valor derivado del preset base.
- Como se comporta la vista de depuracion cuando una ronda no activa ningun efecto.
- Que salida se genera cuando un escenario guardado referencia un preset inexistente o retirado.
- Como se evita que la ejecucion de escenarios reutilizables dependa de pasos de UI ocultos o selecciones intermedias.
- Como se mantiene legibilidad cuando el encounter produce multiples efectos, reshuffle, splash y cierre de partida en una misma secuencia.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST reemplazar el framing actual de `fishing-duel` por un sandbox formal de encounters orientado a exploracion, QA, debugging, balance y regresion.
- **FR-002**: El sandbox MUST ofrecer un modo guiado compatible con el flujo actual para usuarios que no necesiten control total del setup.
- **FR-003**: El sandbox MUST ofrecer un modo manual con seleccion explicita de preset de baraja del jugador, preset de cana, preset de aditamentos, preset de pez y contexto de agua.
- **FR-003a**: El modo manual MUST permitir seleccionar cartas concretas del jugador y del pez para construir encuentros de prueba sin limitarse al orden o composicion completa de los presets base. La identificacion de cartas se realiza por efecto y tipo de trigger (ej. "avance horizontal 2", "avance vertical 1 con reshuffle"), buscando en el catalogo del preset base. Se debe indicar claramente el origen de cada carta: preset_base, manual_replacement, o scenario_defined.
- **FR-004**: El sandbox MUST permitir fijar una seed reproducible para repetir la misma configuracion derivada y el mismo comportamiento observable bajo las mismas entradas.
- **FR-005**: El sandbox MUST permitir definir o sobrescribir al menos el resultado de apertura o cast band, la distancia inicial y la profundidad inicial del encounter.
- **FR-006**: El sandbox MUST permitir controlar si las barajas relevantes usan orden fijo o barajado cuando el escenario de prueba lo requiera.
- **FR-007**: El sandbox MUST definir un contrato de overrides de estado de prueba para encounter y deck que permita crecer de forma controlada sin mezclar reglas de runtime con presentation o CLI.
- **FR-007-MVP**: Como subconjunto MVP de FR-007, el sandbox MUST permitir sobrescribir solo: (a) distancia inicial del encounter, (b) profundidad inicial del encounter, (c) `CaptureDistance` base del `encounter.Config`, y (d) recycle count del deck. `ExhaustionCaptureDistance` se hereda del `encounter.Config` resuelto y no se sobrescribe en esta iteracion. Overrides de exhaustion state, visibilidad de descarte y estado previo a splash quedan diferidos.
- **FR-007a**: El sandbox MUST permitir que la seleccion manual de cartas concretas conviva con presets base y overrides, indicando claramente cuando una carta fue tomada del preset, reemplazada manualmente o definida por escenario.
- **FR-008**: El sandbox MUST validar las combinaciones seleccionadas y comunicar de forma clara cuando una configuracion manual sea incompatible o invalida.
- **FR-009**: El sandbox MUST mostrar la traza de resolucion de cada ronda con triggers activados, efectos resueltos, prioridad y desempates aplicados.
- **FR-010**: El sandbox MUST ofrecer al menos un modo de salida adicional al modo por defecto para inspeccion mas detallada del runtime.
- **FR-011**: La salida mas detallada MUST mostrar suficiente contexto para comparar estado antes y despues de la resolucion sin trasladar reglas de negocio a la UI.
- **FR-012**: El sandbox MUST permitir ejecutar escenarios guardados o configuraciones reutilizables para QA y regresion.
- **FR-012a**: Los escenarios guardados MUST poder fijar cartas concretas del jugador y del pez cuando la prueba requiera validar triggers, efectos o combinaciones puntuales.
- **FR-013**: El sandbox MUST ofrecer un camino reproducible basado en escenarios reutilizables con seed fija para lanzar una configuracion completa sin reconfigurar manualmente cada ejecucion.
- **FR-013-MVP**: Como subconjunto MVP de FR-013, el sandbox MUST soportar escenarios reutilizables con seed fija como baseline semi-reproducible, permitiendo que un QA ejecute la misma configuracion varias veces y obtenga el mismo resultado observable. El modo totalmente no interactivo (sin prompts) puede diferirse a una iteracion futura.
- **FR-014**: El sandbox MUST permitir exportar o presentar una evidencia util de ejecucion para debugging o comparacion, como resumen textual estructurado o snapshot estructurado.
- **FR-015**: El sistema MUST mantener la separacion actual de responsabilidades entre runtime, contenido, presentacion y CLI.
- **FR-016**: El sistema MUST mantener `fishing-run` como flujo jugable diferenciado del sandbox y evitar mezclar sus objetivos de uso.
- **FR-017**: El sistema MUST permitir seleccionar directamente el preset de pez sin depender solo del filtro de spawn por contexto cuando se use el modo manual.
- **FR-018**: El sistema MUST incluir una biblioteca inicial de escenarios de QA representativos, incluyendo al menos casos de empate por prioridad, splash de superficie, captura por agotamiento, regresion de migracion legacy y visibilidad de descarte.

### Key Entities *(include if feature involves data)*

- **Sandbox Configuration**: Configuracion completa elegida por el usuario para arrancar un encounter, incluyendo presets, cartas concretas, seed, modo y overrides.
- **Sandbox Card Selection**: Seleccion manual de cartas concretas del jugador y del pez para construir una prueba puntual a partir de presets base o escenarios guardados.
- **Sandbox Override**: Valor manual que reemplaza o ajusta parte del estado derivado de un preset o escenario base.
- **Sandbox Scenario**: Configuracion guardada o reutilizable que describe una prueba reproducible del sandbox.
- **Resolution Trace**: Evidencia observable de una ronda que lista triggers activados, efectos resueltos, prioridades, desempates y resultados.
- **Sandbox Mode**: Forma de operar el sandbox, como guiado, manual o replay de escenario reutilizable.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Un usuario puede configurar manualmente un encounter completo con jugador, pez, contexto y seed en menos de 3 minutos sin editar archivos ni codigo.
- **SC-002**: El 100% de los escenarios iniciales de QA definidos para la feature pueden ejecutarse de forma reproducible al menos dos veces con el mismo resultado observable.
- **SC-003**: Un tester puede identificar el orden de resolucion y el desempate aplicado en el 100% de 5 escenarios representativos sin revisar el codigo fuente.
- **SC-004**: El equipo puede lanzar un escenario reutilizable con seed fija compartido entre dos personas y obtener el mismo comportamiento observable en ambos casos.

## Assumptions

- El sandbox seguira siendo una experiencia CLI y no requiere una interfaz grafica nueva en esta fase.
- El sistema actual de presets y catalogos existentes se reutiliza como fuente de configuracion base.
- La funcionalidad nueva se implementa sobre el runtime actual de encounters sin redefinir las reglas del juego ya formalizadas.
- El flujo guiado del sandbox debe seguir siendo suficiente para demos manuales simples, mientras el modo manual cubre testing fino.
- La exportacion de evidencia puede comenzar con formatos simples y legibles antes de ampliar opciones futuras.

## Out of Scope

- Convertir `fishing-run` en sandbox o mezclar su loop progresivo con el flujo de prueba de encounters.
- Crear una interfaz visual fuera de CLI para el sandbox en esta version.
- Redisenar el sistema de reglas del juego, el contrato de cartas o el engine mas alla de lo necesario para exponer configuracion y trazas del sandbox.
- Construir automatizacion de CI completa alrededor del sandbox en esta tarea; el alcance se centra en la herramienta y sus escenarios reutilizables.
