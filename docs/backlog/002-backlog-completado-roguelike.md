# Backlog completado roguelike

Este documento concentra las tareas ya integradas en `main`.

## Convencion de estado

- `done`: item ya completado e integrado en `main`.

## Core Loop

### BL-029 Descomponer `match.State` y fijar fronteras del runtime tactico
- **Estado**: `done`
- **Tipo**: Calidad + Discovery
- **Objetivo**: reducir el acoplamiento del runtime tactico separando encounter, estado de mazo, recursos del jugador, thresholds de ronda y vistas de presentacion.
- **Resultado esperado**: ownership mas claro del estado mutable del combate y contratos mas estrechos entre `game`, `progression`, `endings`, `presentation` y futuros estados de run.
- **Plan relacionado**: `docs/features/019-descomponer-match-state-y-fronteras-runtime-tactico.md`
- **Notas de cierre**:
  - El runtime tactico ya separa `Round`, `Player` y `Lifecycle` como subestados explicitos dentro de `match.State`.
  - `engine`, `progression`, `endings`, `presentation`, `session` y el runtime de recursos del jugador ya consumen esa frontera mas explicita.
  - Queda fijado que `match.State` representa solo un duelo aislado y no debe absorber estado futuro de run.
- **Prioridad**: Alta

### BL-030 Consolidar runtime de combate y fronteras de paquetes
- **Estado**: `done`
- **Tipo**: Calidad + Delivery
- **Objetivo**: reducir la fragmentacion funcional actual entre `internal/game/`, `internal/rules/`, `internal/progression/`, `internal/endings/` e `internal/encounter/`.
- **Resultado esperado**: mapa de responsabilidades y una refactorizacion acotada que reduzca imports cruzados y clarifique ownership del duelo.
- **Plan relacionado**: `docs/features/020-consolidar-runtime-de-combate-y-fronteras-de-paquetes.md`
- **Notas de cierre**:
  - `presentation` ya consume snapshots tacticos mas estrechos para status, ronda y resumen final.
  - Los thresholds y helpers de captura quedaron consolidados bajo `encounter`.
  - La consolidacion principal del runtime ya esta integrada y la deuda residual quedo trazada y cerrada en `BL-033`.
- **Prioridad**: Media

### BL-033 Adelgazar contratos residuales del runtime tactico post-`BL-030`
- **Estado**: `done`
- **Tipo**: Calidad + Delivery
- **Objetivo**: cerrar el acoplamiento residual concentrado en contratos transicionales del duelo.
- **Resultado esperado**: superficies de lectura y escritura mas estrechas entre `game`, `match`, `progression`, `endings` y `encounter`.
- **Plan relacionado**: `docs/features/021-adelgazar-contratos-residuales-del-runtime-tactico.md`
- **Notas de cierre**:
  - `match.RoundResult` ya devuelve snapshots tacticos estrechos en lugar de reexportar `match.State` completo.
  - El mapping del estado visible del mazo ya vive bajo `match`.
  - `progression`, `endings` y `player/playermoves` ya consumen contratos de mutacion mas finos.
- **Prioridad**: Alta

### BL-034 Desacoplar setup, opening y bootstrap para habilitar la run MVP
- **Estado**: `done`
- **Tipo**: Discovery + Delivery
- **Objetivo**: ejecutar el slice minimo de arquitectura UI-agnostic necesario para que la run no nazca acoplada al CLI, dejando fuera del borde de terminal la orquestacion de setup, opening y entrada al encounter.
- **Resultado esperado**: `cmd/` reducido a composicion, opening reusable fuera de `internal/cli/`, y contratos de app/presentation suficientemente estables para que una run orqueste encounters sin depender de tipos renderizados por el terminal.
- **Dependencias**: `BL-001`, `BL-018`, `BL-020`, `BL-021`
- **Plan relacionado**: `docs/features/022-desacoplar-setup-opening-y-bootstrap-para-run-mvp.md`
- **Notas de cierre**:
  - `cmd/fishing-duel/main.go` ya quedo reducido a composition root del flujo principal.
  - El setup, opening, cast y bootstrap del encounter ya pueden orquestarse desde `internal/app/` sin depender directamente del adaptador CLI.
  - La CLI actual sigue funcionando como adaptador sobre contratos mas UI-agnostic de app y presentation.
- **Prioridad**: Alta

## Fish y Encounters

### BL-005 Disenar sistema extensible de efectos de cartas
- **Estado**: `done`
- **Tipo**: Delivery
- **Objetivo**: consolidar una arquitectura preparada para cartas de pez y de jugador sin seguir agregando casos aislados.
- **Resultado esperado**: sistema tecnico para encadenar efectos de carta por fases, trigger y owner.
- **Plan relacionado**: `docs/features/009-pipeline-de-efectos-de-carta.md`, `docs/features/012-barajas-de-decision-del-jugador.md`, `docs/features/013-primeras-player-cards-con-efectos.md`
- **Prioridad**: Alta

### BL-006 Definir arquetipos de peces
- **Estado**: `done`
- **Tipo**: Discovery + Delivery
- **Objetivo**: definir arquetipos de pez faciles de configurar y llevarlos a una primera implementacion reusable.
- **Resultado esperado**: perfiles mecanicos configurables que construyen barajas de pez y reemplazan wiring manual en presets de inicio.
- **Plan relacionado**: `docs/features/011-arquetipos-de-peces.md`
- **Prioridad**: Alta

### BL-019 Hacer visible el descarte del pez y modular la lectura del historial
- **Estado**: `done`
- **Tipo**: Discovery + Delivery
- **Objetivo**: convertir el descarte del pez en una herramienta estrategica legible durante el encounter.
- **Resultado esperado**: estado de runtime, presentacion y UX que permitan ver el descarte visible del pez por ciclo.
- **Plan relacionado**: `docs/features/015-visibilidad-descarte-del-pez.md`
- **Prioridad**: Media

### BL-020 Disenar apertura del encounter de pesca y minijuego de cast
- **Estado**: `done`
- **Tipo**: Discovery + Delivery
- **Objetivo**: definir la fase previa al combate como una apertura autocontenida del encounter.
- **Resultado esperado**: flujo claro `leer situacion -> resolver cast -> abrir ventana horizontal` y contrato minimo para inyectar contexto de agua.
- **Plan relacionado**: `docs/features/016-apertura-encounter-pesca-y-cast.md`
- **Notas de cierre**:
  - El encounter ya resuelve una apertura previa con contexto de agua, minijuego de cast y configuracion inicial derivada.
  - El contrato ya transporta `InitialDistance` e `InitialDepth`.
- **Prioridad**: Alta

### BL-022 Definir aparicion de peces por aguas base y ventana de lanzamiento
- **Estado**: `done`
- **Tipo**: Discovery
- **Objetivo**: decidir como agua, cast y sesgo vertical del setup seleccionan subconjuntos de peces y encuentros compatibles.
- **Resultado esperado**: metadata minima de aparicion en perfiles de pez y criterio para conectar cast, `rod`, aditamentos y habitats.
- **Plan relacionado**: `docs/features/018-aparicion-de-peces-por-aguas-y-ventana-de-lanzamiento.md`
- **Notas de cierre**:
  - El juego ya resuelve `agua -> apertura -> spawn -> mazo del pez`.
  - Los perfiles ya exponen metadata minima de aparicion por pool de agua, distancia, profundidad y habitats.
  - El spawn actual ya usa catalogos tipados para `water pools`, `habitats` y arquetipos.
- **Prioridad**: Alta

### BL-041 Externalizar catalogo base de peces a formato data-driven
- **Estado**: `done`
- **Tipo**: Delivery
- **Objetivo**: mover el catalogo actual de perfiles de pez fuera del codigo a un formato data-driven validable y definir pools cerradas por encounter.
- **Resultado esperado**: catalogo global de peces y `fish pools` cerradas reutilizables que referencien perfiles por id sin duplicar su definicion.
- **Plan relacionado**: `docs/features/024-externalizar-catalogo-base-de-peces-y-fish-pools.md`
- **Notas de cierre**:
  - El catalogo base de peces ya vive en JSON embebido con carga y validacion explicita.
  - Las `fish pools` ya soportan subsets cerrados por id y entradas ponderadas por peso.
  - El spawn ya puede trabajar sobre subcatalogos resueltos desde pools sin cambiar su contrato base de `[]Profile`.
- **Prioridad**: Alta

### BL-045 Resolver spawns de encounter desde catalogos data-driven
- **Estado**: `done`
- **Tipo**: Delivery
- **Objetivo**: conectar el flujo `agua -> apertura -> spawn -> mazo del pez` a los catalogos y pools cerradas.
- **Resultado esperado**: pipeline de spawn consumiendo catalogos externos de peces y `fish_pool_id` concretas.
- **Plan relacionado**: `docs/features/025-resolver-spawns-desde-catalogos-y-fish-pools.md`
- **Notas de cierre**:
  - El bootstrap del encounter ya acepta catalogo y `fish_pool_id` como configuracion explicita.
  - El flujo actual mantiene fallback a catalogo y pool por defecto sin romper la CLI.
  - La capa de app ya separa la resolucion del subset de peces de la resolucion del spawn dentro de ese subset.
- **Prioridad**: Alta

### BL-050 Reducir el determinismo total del spawn de peces
- **Estado**: `done`
- **Tipo**: Discovery + Delivery
- **Objetivo**: evitar que el sistema actual de spawn sea completamente deducible.
- **Resultado esperado**: estrategia de variedad controlada con contratos reproducibles para tests y configuracion futura desde data-driven.
- **Plan relacionado**: `docs/features/023-reducir-determinismo-del-spawn-de-peces.md`
- **Notas de cierre**:
  - El matching contextual del spawn se mantiene estable y separado de la seleccion final del pez.
  - El runtime ya permite inyectar randomizer para variar reproduciblemente entre candidatos empatados en score.
  - El camino sin randomizer sigue siendo estable para tests y diagnostico del dominio.
- **Prioridad**: Alta

## Items y Build

### BL-021 Redefinir la cana como `rod` y sus aditamentos como base del setup del jugador
- **Estado**: `done`
- **Tipo**: Discovery
- **Objetivo**: reemplazar el vocabulario ambiguo de `rig` por un modelo explicito de `rod` y separar limites de apertura y track.
- **Resultado esperado**: modelo de `rod` con limites estructurales separados, taxonomia inicial de aditamentos y nomenclatura clara.
- **Plan relacionado**: `docs/features/017-rod-y-limites-de-apertura-y-track.md`
- **Notas de cierre**:
  - El vocabulario de `rig` ya fue sustituido por `rod`.
  - La apertura del encounter ya valida contra limites efectivos de apertura, mientras el tablero usa limites de track.
  - Los aditamentos ya modifican limites de apertura y track y transportan `HabitatTags`.
- **Prioridad**: Alta

## Tech y UX

### BL-018 Mejorar arquitectura y gobierno de paquetes
- **Estado**: `done`
- **Tipo**: Calidad + Delivery
- **Objetivo**: contener el crecimiento de `internal/` con reglas claras de organizacion y una primera refactorizacion acotada.
- **Resultado esperado**: estrategia de estructura de paquetes mas sostenible y una mejora concreta que reduzca acoplamiento o dispersion actual.
- **Plan relacionado**: `docs/features/014-arquitectura-y-gobierno-de-paquetes.md`
- **Prioridad**: Alta
