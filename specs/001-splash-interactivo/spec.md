# Feature Specification: Splash interactivo con saltos y mejoras de cana

**Feature Branch**: `001-splash-interactivo-cana`  
**Created**: 2026-05-07  
**Status**: Draft  
**Input**: User description: "Disenos la tarea del Github Issue RMF-06, consulta toda la informacion del issue para el diseno del spec, realiza las preguntas necesarias para aclarar el desarrollo."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Resolver un splash con habilidad (Priority: P1)

Como jugador en un encounter, quiero responder al chapoteo del pez con una prueba breve de habilidad para que el resultado del splash dependa de mi reaccion y no solo de una tirada oculta.

**Why this priority**: Es el cambio central del item. Sin esta historia, el splash sigue siendo un evento aleatorio y no entrega la expresividad ni la tension buscadas para el loop base.

**Independent Test**: Puede probarse activando un splash en un encounter aislado y verificando que el jugador recibe una interaccion con tiempo limite cuyo resultado determina si el pez sigue sujeto o escapa.

**Acceptance Scenarios**:

1. **Given** un encounter en el que el pez dispara un splash, **When** comienza la prueba interactiva, **Then** el jugador ve una ventana de reaccion con limite de tiempo antes de que se resuelva el evento.
2. **Given** un splash activo, **When** el jugador acierta la accion requerida dentro del tiempo, **Then** el pez no escapa por ese salto y el encounter continua.
3. **Given** un splash activo, **When** el jugador falla la accion requerida o deja vencer el tiempo, **Then** el splash termina en escape del pez y el encounter se cierra con ese motivo.

---

### User Story 2 - Resolver peces con multiples saltos (Priority: P2)

Como disenador de encounters, quiero configurar peces que pidan varios saltos consecutivos para que el splash escale en dificultad y personalidad segun la especie o situacion.

**Why this priority**: La configuracion por especie forma parte del valor diferencial del refactor, pero depende de que el splash interactivo base ya exista.

**Independent Test**: Puede probarse con dos peces configurados con distinto numero de saltos y verificando que el jugador deba superar la cantidad exacta pedida para evitar el escape.

**Acceptance Scenarios**:

1. **Given** un pez configurado con tres saltos, **When** entra en splash, **Then** el jugador debe superar tres pruebas consecutivas para completar el evento sin escape.
2. **Given** un pez configurado con varios saltos, **When** el jugador falla cualquiera de ellos, **Then** el pez escapa inmediatamente sin resolver los saltos restantes.

---

### User Story 3 - Aprovechar mejoras de cana durante el splash (Priority: P3)

Como jugador con una mejora de cana orientada al splash, quiero recibir un beneficio tangible por cada salto ganado para que la build afecte este momento de forma visible y recompensante.

**Why this priority**: Amplia el valor del sistema hacia build y progresion, pero solo despues de que el loop interactivo y la configuracion de saltos esten claros.

**Independent Test**: Puede probarse activando un splash con una mejora de cana que acerque al pez tras cada salto ganado y comparando el resultado frente al mismo splash sin esa mejora.

**Acceptance Scenarios**:

1. **Given** una mejora de cana que recompensa cada salto ganado, **When** el jugador supera un salto, **Then** el pez se acerca la cantidad configurada antes de pasar al siguiente salto o cerrar el evento.
2. **Given** un splash de varios saltos con mejora activa, **When** el jugador supera todos los saltos, **Then** el acercamiento acumulado refleja todas las recompensas aplicables ganadas durante la secuencia.

---

### Edge Cases

- Que ocurre si el splash se dispara cuando el pez ya esta en una posicion que tambien cumpliria una condicion de captura o fin del encounter.
- Que ocurre si el tiempo expira exactamente al mismo tiempo que la confirmacion del jugador.
- Como se resuelve un pez configurado con el minimo o maximo de saltos permitidos.
- Que ocurre si una mejora de cana otorga acercamiento pero el pez ya esta en el limite minimo de distancia permitido.
- Como se presenta y resuelve un splash cuando el jugador encadena exitos parciales y luego falla un salto intermedio.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: El sistema MUST reemplazar la resolucion aleatoria actual del splash por una interaccion de habilidad del jugador cada vez que el encounter dispare ese evento.
- **FR-002**: El sistema MUST mantener los mismos detonantes de splash ya existentes en el encounter, salvo que la spec posterior explicite nuevas fuentes del evento.
- **FR-003**: El sistema MUST iniciar cada splash con una ventana de reaccion limitada en el tiempo y visible para el jugador.
- **FR-004**: El sistema MUST permitir que la duracion inicial de la ventana de reaccion sea configurable por contenido, partiendo de un valor base de 1 segundo para el primer slice.
- **FR-005**: El sistema MUST resolver cada salto con un resultado binario claro: exito si el jugador acierta dentro del tiempo, o fallo si se equivoca o deja vencer la ventana.
- **FR-006**: El sistema MUST cerrar el encounter por escape de splash si el jugador falla cualquier salto de la secuencia.
- **FR-007**: El sistema MUST permitir configurar entre 1 y 5 saltos consecutivos por splash para cada pez o contexto que use esta mecanica.
- **FR-008**: El sistema MUST exigir que el jugador complete todos los saltos configurados para considerar el splash superado sin escape.
- **FR-009**: El sistema MUST mostrar el progreso del splash de forma que el jugador pueda distinguir cuantos saltos faltan por resolver en una secuencia de varios intentos.
- **FR-010**: El sistema MUST permitir mejoras de cana que otorguen un acercamiento adicional del pez por cada salto ganado.
- **FR-011**: El sistema MUST aplicar la recompensa de acercamiento inmediatamente despues de cada salto ganado antes de resolver el siguiente salto, si existe uno.
- **FR-012**: El sistema MUST impedir que una recompensa de acercamiento mueva al pez mas alla del limite minimo valido del encounter.
- **FR-013**: El sistema MUST dejar el comportamiento de splash parametrizable por contenido para que distintas especies o presets cambien cantidad de saltos y recompensas aplicables sin redefinir el flujo base.
- **FR-014**: El sistema MUST registrar el resultado final del splash con suficiente claridad para que presentacion, resumen de ronda y condiciones terminales distingan entre exito del jugador y escape del pez.
- **FR-015**: El sistema MUST seguir permitiendo ejecutar y validar encounters manuales donde el splash pueda forzarse o reproducirse de manera consistente para pruebas de contenido y balance.

### Key Entities *(include if feature involves data)*

- **Splash Sequence**: Evento interactivo del encounter que define una serie ordenada de saltos, su limite de tiempo, su estado actual y su resultado final.
- **Splash Jump**: Unidad individual dentro de un splash que puede terminar en exito o fallo y que puede activar recompensas por salto ganado.
- **Splash Profile**: Configuracion reusable que describe cuantos saltos requiere un pez, que ventana base usa y que modificadores o recompensas pueden intervenir.
- **Rod Splash Bonus**: Mejora ligada a la cana que define beneficios aplicables cuando el jugador gana uno o mas saltos del splash.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: En el 100% de los encounters donde se active splash, el jugador recibe una interaccion jugable antes de que el pez pueda escapar por este motivo.
- **SC-002**: Un tester puede completar sin ayuda externa escenarios de splash con 1, 3 y 5 saltos y verificar en menos de 2 minutos el resultado esperado de exito total o escape inmediato por fallo.
- **SC-003**: En pruebas de contenido, un disenador puede diferenciar al menos dos perfiles de pez por dificultad de splash cambiando solo su configuracion de saltos y tiempo, sin redefinir las reglas base del encounter.
- **SC-004**: En pruebas comparativas, una mejora de cana orientada al splash produce un cambio observable en la distancia final del pez en el 100% de los saltos ganados donde la recompensa aplique.

## Assumptions

- El slice cubre el flujo actual de CLI y encounter runtime del MVP, sin exigir aun una interfaz grafica dedicada.
- Los detonantes actuales del splash se conservan en esta iteracion para limitar el alcance a la resolucion interactiva del evento.
- La accion exacta del jugador durante cada salto puede mantener un esquema simple de confirmacion y timing en esta primera version, siempre que siga siendo una prueba de habilidad con tiempo limite.
- La recompensa de cana por splash es acumulativa por salto ganado, no solo por completar la secuencia entera.
- El disenador necesita configurar la dificultad del splash desde contenido del pez o del encounter, no desde ajustes globales fijos.
- La spec no introduce nuevas economias, patrocinadores ni metaprogresion; solo deja preparado el punto de integracion para mejoras de cana relacionadas con splash.
