# Plan de feature: definir-recursos-globales-de-run-y-angler-profiles

## Objetivo

Extender el alcance de `BL-036` para que no solo fije los recursos globales de la run y las reglas del hilo del MVP, sino tambien la identidad inicial jugable con la que el jugador entra a la expedicion: los `AnglerProfile`.

La idea central es que el inicio de una run deje de componerse desde elecciones tecnicas separadas de mazo, `rod` y aditamentos, y pase a resolverse desde un perfil de pescador seleccionable, equivalente al rol que cumple la baraja inicial en otros roguelikes de runs. Ese pescador debe definir un paquete inicial coherente de estilo de juego y, al mismo tiempo, convertirse en el punto de apoyo correcto para los futuros desbloqueos entre runs.

Esta feature debe cerrar dos cosas a la vez:

- que recursos globales persisten durante toda la run, especialmente el `thread` como condicion principal de derrota
- como se encapsula el paquete inicial de la run dentro de un `AnglerProfile`, sin acoplar `rod`, barajas y meta-progresion en el lugar incorrecto

## Criterios de aceptacion

- Queda definido el contrato del recurso global `thread` del MVP: capacidad inicial, perdida, reparacion parcial y condicion de derrota.
- Queda definido que el `thread` es parte del estado global de run y no del runtime tactico del encounter.
- Existe un concepto explicito de `AnglerProfile` como identidad inicial jugable de la run.
- Un `AnglerProfile` define al menos:
  - baraja inicial del jugador
  - `rod` inicial
  - aditamentos iniciales o build base equivalente
  - `thread` inicial
- Queda definida la frontera entre `AnglerProfile` y `rod`:
  - `rod` sigue modelando herramienta y limites tacticos
  - `AnglerProfile` modela el paquete inicial de run
- El flujo de inicio de `cmd/fishing-run/` pasa a depender de la seleccion de `AnglerProfile`, no de elegir por separado deck, `rod` y attachments.
- Queda preparado el borde para futuros desbloqueos entre runs sin exigir persistencia real todavia.
- Se explicita que el sandbox de `cmd/fishing-duel/` puede seguir conservando su setup manual independiente.

## Scope

### Incluye

- Definir las reglas MVP del `thread` como recurso global de run.
- Definir el shape de `AnglerProfile` y sus dependencias sobre contenido existente.
- Reemplazar en la run MVP el setup manual fragmentado por una seleccion de pescador.
- Definir un catalogo inicial de `AnglerProfile` para el MVP.
- Introducir los contratos necesarios para que un `AnglerProfile` construya el estado inicial de run.
- Dejar modelados los metadatos minimos para futuros desbloqueos entre runs, aunque no se implementen aun.

### No incluye

- Persistencia real de progreso meta entre runs.
- Sistema final de desbloqueos permanentes.
- Economia ampliada del servicio mas alla de lo necesario para reparar o preservar `thread` en el MVP.
- Reemplazar el sandbox manual de `cmd/fishing-duel/`.
- Disenar todavia pasivas complejas, perks unicos o arboles de progresion para cada pescador.

## Problema que resuelve

Hoy la run MVP ya existe como estructura secuencial minima, pero arranca todavia desde una composicion demasiado tecnica para la experiencia que queremos:

- el jugador elige deck preset
- luego elige `rod`
- luego elige aditamentos

Eso sirve como harness de prototipo, pero no como identidad de run. Si el objetivo es que la seleccion inicial se parezca mas a elegir un personaje o un estilo de expedicion, hace falta una entidad superior que encapsule ese arranque.

Al mismo tiempo, `BL-036` ya iba a encargarse del `thread` como recurso global. Es mejor extender ese item en vez de abrir otra linea aislada, porque ambas decisiones estan conectadas:

- el `thread` inicial ya no deberia salir de defaults genricos, sino del pescador elegido
- el pescador es quien fija el paquete base de supervivencia y estilo

## Direccion conceptual

### 1. `AnglerProfile` como identidad jugable de la run

Un `AnglerProfile` debe representar un personaje o arquetipo inicial de expedicion. No es una `rod` expandida ni una mera etiqueta cosmetica. Es el paquete inicial con el que el jugador entra a la run.

Direccion recomendada:

- `AnglerProfile` = identidad seleccionable del jugador para una run
- `rod` = herramienta/equipo dentro de ese perfil
- `loadout` = estado tactico derivado de esa seleccion

Eso permite preservar responsabilidades limpias:

- `rod` sigue describiendo limites y configuracion del duelo
- `loadout` sigue siendo el estado usable por el encounter
- `AnglerProfile` describe desde que paquete arranca la run

### 2. `thread` como recurso global asociado al pescador y a la run

El `thread` debe quedar definido como una reserva global de supervivencia de la expedicion, no como una barra local del encuentro. La capa de run decide cuanto `thread` queda, cuando se pierde, cuando se repara y cuando se agota la run.

Direccion MVP propuesta:

- cada `AnglerProfile` define un `starting_thread`
- los encounters y, mas adelante, algunos nodos/servicios pueden producir dano o reparacion de `thread`
- si `thread` llega a 0, la run termina en derrota
- el `thread` no se resetea entre nodos

### 3. Desbloqueos futuros sin persistencia inmediata

Los `AnglerProfile` deben nacer con la idea de que algunos estaran disponibles al inicio y otros no. Pero en esta etapa no hace falta implementar guardado ni meta-progresion real.

La recomendacion es modelar desde ya metadatos como:

- `unlocked_by_default`
- `unlock_key` o `unlock_id`
- `unlock_description` opcional

La run MVP puede filtrar solo los perfiles desbloqueados por defecto mientras la capa meta aun no existe.

## Propuesta de modelo

### 1. Nuevo modulo de contenido

Crear un modulo dedicado, por ejemplo:

- `internal/content/anglerprofiles/`

No conviene mezclarlo con `playerprofiles/` porque ese paquete hoy representa presets de baraja, no identidades completas de run.

### 2. Shape recomendado de `AnglerProfile`

Direccion propuesta:

- `id`
- `name`
- `description`
- `details`
- `deck_preset_id`
- `rod_preset_id`
- `attachment_preset_id` o lista equivalente si luego hiciera falta mas de uno
- `starting_thread`
- `unlocked_by_default`
- `unlock_id` opcional para el futuro

Posibles extensiones posteriores, fuera del MVP:

- pasiva unica
- dinero inicial distinto
- modificadores de spawn
- afinidades con habitats o servicios

### 3. Contrato de build inicial

Cada `AnglerProfile` debe poder resolverse en algo como un `RunStartLoadout` o `ResolvedAnglerStart`, con:

- deck preset del jugador ya resuelto
- `loadout.State` inicial ya construido
- `thread` inicial

Eso evita que `RunSession` tenga que conocer demasiados detalles de como combinar presets internos.

## Flujo propuesto de inicio de run

### Estado actual

`cmd/fishing-run/` reutiliza hoy el setup manual separado del sandbox:

- deck
- `rod`
- attachments

### Estado deseado para esta feature

`cmd/fishing-run/` debe arrancar asi:

- elegir `AnglerProfile`
- resolver su paquete inicial
- inicializar `run.State` con ese paquete
- entrar a la ruta de nodos

El sandbox de `cmd/fishing-duel/` puede seguir funcionando con el setup manual actual porque su objetivo sigue siendo tactico y de debugging.

## Reglas MVP del `thread`

Esta feature debe cerrar al menos estas reglas:

- `thread` es entero y vive en `run.State.Thread`
- `starting_thread` depende del `AnglerProfile`
- el dano al `thread` se aplica al terminar el encounter a traves de `EncounterResult`
- el `thread` nunca puede exceder su maximo actual sin una regla explicita
- los nodos de servicio podran reparar `thread`, aunque el detalle economico completo quede para `BL-037` y `BL-038`
- si `thread == 0`, la run termina en `defeat`

Decision recomendada para MVP:

- tratar `thread` como recurso actual + maximo, donde el maximo inicial tambien viene del pescador

## Catalogo inicial recomendado

Para el MVP bastan 3 perfiles jugables y 1 o 2 bloqueados a nivel de datos:

- `coastal-specialist`
  - menos `thread`
  - mejor control cercano
- `deep-angler`
  - mejor apertura/profundidad
  - tension mayor en horizontal o menor `thread`
- `steady-handler`
  - mas `thread`
  - build mas sobria o menos explosiva

Perfiles bloqueados de ejemplo:

- `storm-reader`
- `channel-hunter`

Aunque no se puedan desbloquear todavia, ya sirven para validar la estructura del catalogo.

## Implementacion propuesta por fases

### Fase 1. Modelado y contenido base

- Crear `internal/content/anglerprofiles/`.
- Definir el tipo `AnglerProfile`.
- Crear catalogo por defecto con perfiles iniciales.
- Resolver referencias a `playerprofiles`, `rodpresets` y `attachmentpresets`.

### Fase 2. Resolucion del paquete inicial de run

- Crear un resolvedor o builder de inicio de run a partir de `AnglerProfile`.
- Generar desde ahi:
  - preset de baraja del jugador
  - `loadout.State`
  - `thread` inicial
- Validar que el contrato resultante sea apto para inicializar `run.State`.

### Fase 3. Rewire del inicio de `cmd/fishing-run/`

- Reemplazar `resolvePlayerSetup(...)` en la run por una seleccion de pescador.
- Mantener el setup manual actual solo para `cmd/fishing-duel/`.
- Extender UI y presentation con vistas de seleccion/confirmacion de `AnglerProfile`.

### Fase 4. Integracion formal con `thread`

- Hacer que `run.NewState(...)` reciba el `thread` inicial resuelto desde el pescador.
- Confirmar reglas de dano, derrota y limites de reparacion.
- Dejar tests que cubran perfiles con distintos valores de `thread`.

### Fase 5. Preparacion del borde meta

- Introducir flags de desbloqueo en datos.
- Filtrar en la UI de run solo perfiles disponibles por defecto.
- Dejar la interfaz lista para que una futura capa meta provea el conjunto desbloqueado.

## Archivos o zonas probables

- `internal/content/anglerprofiles/`
- `internal/app/run_session.go`
- `internal/run/`
- `internal/cli/`
- `internal/presentation/`
- `cmd/fishing-run/main.go`
- `README.md`

## Riesgos y decisiones abiertas

- Hay que decidir si el `AnglerProfile` puede apuntar a un solo preset de attachments o si conviene dejarlo como lista para no trabar crecimiento posterior.
- Hay que decidir si el `thread` inicial y su maximo siempre coinciden en el MVP; mi recomendacion es que si.
- Hay que evitar meter dentro de `AnglerProfile` demasiadas reglas futuras desde ahora; primero debe resolver identidad inicial y `thread`.
- Hay que mantener la distincion entre contenido de run (`AnglerProfile`) y contenido tactico (`playerprofiles`, `rodpresets`, `attachmentpresets`).

## Accionables inmediatos

- Pasar `BL-036` a `planned` con alcance extendido a `AnglerProfile`.
- Crear este nuevo modulo de contenido.
- Reemplazar el setup inicial de la run por la seleccion de pescador.
- Fijar el contrato de `thread` inicial derivado del perfil.
- Dejar listos tests de resolucion e inicializacion de run por perfil.
