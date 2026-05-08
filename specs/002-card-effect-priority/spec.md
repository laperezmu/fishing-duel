# Feature Specification: Rediseno de triggers y efectos de cartas con prioridad de resolucion

**Feature Branch**: `002-card-effect-priority`  
**Created**: 2026-05-07  
**Status**: Draft  
**Input**: User description: "He actualizado la lista de triggers y efectos en el issue de Github, actualiza el spec acorde al requerimiento."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Definir un catalogo jugable de efectos y triggers (Priority: P1)

Como disenador del juego, quiero un catalogo claro de triggers y efectos soportados para poder crear cartas de pez y jugador sin ambiguedades sobre cuando se activa cada efecto.

**Why this priority**: El valor principal del item es cerrar la definicion funcional del sistema de cartas. Sin este catalogo, el contenido nuevo y la migracion del contenido actual siguen bloqueados.

**Independent Test**: Puede probarse revisando una carta de pez y una del jugador y verificando que cada efecto se puede clasificar en un trigger permitido y en un efecto permitido sin recurrir a interpretaciones ad hoc.

**Acceptance Scenarios**:

1. **Given** una carta con multiples efectos, **When** se revisa su definicion, **Then** cada efecto queda asociado a un unico trigger valido y mantiene independencia respecto de los demas efectos de la misma carta.
2. **Given** una propuesta de carta nueva, **When** el disenador selecciona sus reglas, **Then** solo puede usar triggers y efectos incluidos en el catalogo objetivo aprobado por esta spec, incluyendo variantes por color y por entidad cuando correspondan.
3. **Given** efectos historicos basados en umbrales de captura o superficie, **When** se evalua el catalogo nuevo, **Then** queda definido con que efectos o combinaciones deben reemplazarse y en que casos dejan de existir.

---

### User Story 2 - Resolver efectos simultaneos de forma determinista (Priority: P2)

Como jugador y tester, quiero que los efectos simultaneos se resuelvan siempre en un orden estable y comprensible para que el resultado de una ronda no dependa de interpretaciones inconsistentes.

**Why this priority**: El issue exige prioridad explicita de resolucion y desempate a favor del pez. Sin esta regla, el sistema sigue siendo dificil de balancear, probar y explicar.

**Independent Test**: Puede probarse construyendo una ronda en la que ambos lados disparen efectos al mismo tiempo y verificando que el orden y el resultado final coincidan siempre con la tabla de prioridad definida.

**Acceptance Scenarios**:

1. **Given** una ronda donde pez y jugador activan efectos del mismo nivel de prioridad, **When** se resuelven, **Then** el efecto del pez se aplica antes que el del jugador.
2. **Given** una ronda con efectos de distinta prioridad, **When** se resuelven, **Then** todos los efectos se aplican de mayor a menor prioridad segun la regla formal definida por la spec.
3. **Given** una ronda con multiples efectos validos dentro de la misma carta, **When** se resuelven, **Then** el orden depende de la prioridad de cada efecto y no de la posicion accidental de la carta o de la carga del contenido.

---

### User Story 3 - Migrar contenido y reglas sin perder cobertura funcional (Priority: P3)

Como responsable de la plataforma, quiero una estrategia de migracion explicita para cartas, runtime y progreso para poder reemplazar el contrato actual sin dejar huecos funcionales ni regresiones invisibles.

**Why this priority**: El item mezcla discovery y calidad. La migracion no debe preceder a la definicion del catalogo ni a la prioridad, pero es necesaria para convertir la decision en trabajo implementable por slices.

**Independent Test**: Puede probarse tomando el contenido actual del juego y verificando que cada trigger y efecto existente queda clasificado como soportado, reemplazado o retirado con una ruta de migracion definida.

**Acceptance Scenarios**:

1. **Given** el inventario actual de triggers y efectos usados por el contenido del repo, **When** se compara con el catalogo objetivo, **Then** cada caso queda marcado como mantener, cambiar de semantica, reemplazar o deprecar.
2. **Given** una carta existente que hoy usa `CaptureDistanceBonus` o `SurfaceDepthBonus`, **When** se planifica su migracion, **Then** la spec indica con que efecto objetivo se sustituye y que comportamiento observable debe preservarse o cambiar, distinguiendo si el resultado aplica al pez o al jugador.
3. **Given** el alcance aprobado de esta tarea, **When** el equipo pase a planificacion, **Then** existe una secuencia clara de slices para runtime, contenido, progresion y validacion sin ambiguedades sobre dependencias o alcance.

---

### Edge Cases

- Que ocurre cuando una carta contiene varios efectos que comparten trigger pero apuntan a prioridades distintas.
- Que ocurre cuando ambos lados activan un efecto del mismo tipo y misma prioridad en la misma ronda.
- Como se resuelve un empate de color concreto cuando tambien aplica la regla general de empate.
- Como se resuelve una victoria o derrota por color concreto cuando otro efecto del mismo trigger cambia antes la posicion del encounter.
- Como se interpreta un efecto de fatiga cuando dos colores del jugador podrian verse afectados por reglas distintas durante la misma ronda.
- Que ocurre con cartas legacy cuyo comportamiento actual combina en un solo efecto cambios de umbral y cambios de posicion.
- Como se trata contenido existente que usa efectos retirados pero no tiene un reemplazo uno a uno exacto.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST definir un catalogo objetivo cerrado de triggers soportados para efectos de cartas de pez y jugador.
- **FR-002**: El catalogo objetivo de triggers MUST incluir, como minimo, activacion cuando el dueno de la carta gana, cuando pierde, cuando el resultado es empate, cuando el resultado es empate de un color concreto, cuando el dueno de la carta gana con un color concreto, cuando el dueno de la carta pierde con un color concreto, cuando la carta es usada, cuando el pez ya hizo splash, cuando la carta esta en descarte, cuando el pez ya hizo reshuffle y cuando el pez ya se fatigo.
- **FR-003**: El sistema MUST definir para cada trigger objetivo si mantiene la semantica actual, cambia de semantica o se incorpora como capacidad nueva.
- **FR-004**: El sistema MUST establecer que cada efecto individual pertenece a un unico trigger.
- **FR-005**: El sistema MUST permitir que una misma carta contenga multiples efectos independientes, incluso si usan triggers distintos.
- **FR-006**: El sistema MUST definir un catalogo objetivo cerrado de efectos soportados.
- **FR-007**: El catalogo objetivo de efectos MUST incluir, como minimo, avance horizontal, avance vertical, reshuffle de la baraja actual, obligar a descartar de un color concreto, fatiga instantanea aplicada solo al pez, fatiga de un color concreto aplicada solo al jugador, alterar el orden del descarte del jugador, ocultar descarte temporalmente y acercamiento por splash exitoso.
- **FR-007a**: La semantica observable de cada efecto del catalogo MUST quedar definida de forma explicita. En particular, "ocultar descarte temporalmente" significa que las cartas afectadas permanecen en descarte pero su informacion visible deja de mostrarse hasta que termine la ventana definida por el efecto, y "acercamiento por splash exitoso" significa que un splash exitoso aplica un avance horizontal inmediato del pez hacia captura dentro de la misma resolucion.
- **FR-008**: El sistema MUST deprecar `CaptureDistanceBonus` y `SurfaceDepthBonus` como efectos objetivos del modelo nuevo.
- **FR-009**: El sistema MUST documentar para cada efecto deprecado si se reemplaza por un efecto nuevo, por una combinacion de efectos, o si su comportamiento queda fuera de alcance.
- **FR-010**: El sistema MUST diferenciar explicitamente los efectos que comparten representacion visual pero aplican a entidades distintas, para evitar ambiguedad entre fatiga del pez y fatiga del jugador por color.
- **FR-011**: El sistema MUST definir una regla formal de prioridad de resolucion aplicable a todos los efectos simultaneos de una ronda o evento.
- **FR-011a**: La prioridad MUST interpretarse con una convencion unica y estable: un valor numerico mayor se resuelve antes que un valor numerico menor.
- **FR-012**: La regla formal de prioridad MUST producir siempre el mismo orden de resolucion ante el mismo estado inicial y los mismos efectos activados.
- **FR-013**: Cuando dos efectos aplicables tengan la misma prioridad, el sistema MUST resolver primero el efecto del pez y luego el del jugador.
- **FR-014**: El sistema MUST especificar como se resuelven multiples efectos aplicables dentro de una misma carta cuando comparten trigger pero no comparten prioridad.
- **FR-015**: El sistema MUST definir la relacion entre triggers especificos y triggers generales para evitar dobles activaciones ambiguas, incluyendo empate general vs empate de color y victoria o derrota general vs victoria o derrota por color concreto.
- **FR-016**: El sistema MUST incluir un analisis de cobertura del contenido actual indicando que triggers y efectos ya existen en el repo, cuales faltan y cuales requieren cambios de semantica.
- **FR-017**: El sistema MUST incluir una estrategia de migracion del runtime de cartas, encounter y progresion hacia el nuevo contrato.
- **FR-018**: La estrategia de migracion MUST desglosar al menos los frentes de contrato de cartas, contenido existente, resolucion de rounds, resumenes o estados derivados y validacion de regresion.
- **FR-019**: El sistema MUST dejar el alcance listo para planificacion por slices posteriores sin introducir sistemas meta ni redisenos fuera del dominio de cartas y encounters.

### Key Entities *(include if feature involves data)*

- **Trigger Catalog Entry**: Definicion funcional de un trigger permitido, con su condicion de activacion, su semantica observable y su relacion con triggers mas generales o mas especificos.
- **Effect Catalog Entry**: Definicion funcional de un efecto permitido, con su resultado observable, su prioridad de resolucion y sus restricciones de uso.
- **Entity-Specific Effect Variant**: Variante de un efecto que comparte lectura visual con otra, pero cambia su destino, restricciones o consecuencias segun aplique al pez o al jugador.
- **Card Effect Binding**: Asociacion entre un efecto individual y su trigger unico dentro de una carta concreta.
- **Resolution Priority Rule**: Regla que ordena todos los efectos activados en una misma resolucion, incluyendo desempates entre pez y jugador.
- **Migration Coverage Map**: Inventario que vincula triggers y efectos actuales con su estado objetivo: mantener, ajustar, reemplazar o retirar.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: El 100% de los triggers y efectos usados por el contenido actual del repositorio quedan clasificados en una matriz de cobertura como soportados, ajustados, reemplazados o retirados.
- **SC-002**: Un tester puede revisar al menos 5 escenarios de resolucion simultanea representativos y obtener el mismo orden y resultado esperado en el 100% de las repeticiones.
- **SC-003**: Un disenador puede describir una carta con multiples efectos y determinar sin ayuda externa que trigger activa cada efecto y en que momento se resuelve.
- **SC-004**: El equipo puede derivar un plan de implementacion por slices a partir de esta spec sin abrir preguntas de alcance sobre deprecaciones, prioridad o cobertura del contenido actual.

## Assumptions

- La tarea se centra en definir el contrato funcional del sistema de cartas y no en implementar en esta etapa toda la migracion tecnica.
- La dependencia `RMF-06` no altera el objetivo central del item y se asume como contexto de backlog pendiente, no como bloqueo para redactar esta spec.
- El comportamiento actual basado en `DistanceShift`, `DepthShift` y otros efectos no deprecados sigue siendo valido salvo donde esta spec indique cambios de semantica o relaciones nuevas con otros triggers.
- La prioridad de resolucion se define a nivel de efecto individual y no a nivel de carta completa.
- Cuando un efecto tenga variantes por color o por entidad, esa distincion forma parte del catalogo funcional y no de una nota editorial opcional.
- El primer alcance no abre nuevas economias, nuevas interfaces de usuario ni cambios de meta-progresion; se limita a cartas, encounter y estados derivados asociados.

## Out of Scope

- Redisenar sistemas de economia, meta-progresion o progresion global fuera de cartas y encounters.
- Introducir nuevas interfaces de usuario o cambiar el flujo principal de `cmd/fishing-duel` y `cmd/fishing-run` mas alla de la compatibilidad necesaria con el contrato nuevo.
- Sustituir la arquitectura modular actual por un sistema nuevo de runtime, contenido o presentacion.
