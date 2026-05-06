# Mapa de migracion a GitHub Issues

Este documento reemplaza la nomenclatura legacy `BL-*` y `GUI-*` por una numeracion ordenada por milestone. El backlog vivo debe gestionarse en GitHub Issues; `docs/backlog/` queda como contexto historico.

## Milestones

### Run MVP Foundation

- `RMF-01` Definir recursos globales de run, reglas del hilo y `AnglerProfile` del MVP
- `RMF-02` Definir economia minima de fotos, dinero y liquidacion de servicio
- `RMF-03` Definir build minima y acciones de servicio del MVP
- `RMF-04` Definir patrocinadores MVP o postergarlos explicitamente
- `RMF-05` Implementar primer vertical slice de run MVP

### Zone Progression

- `ZP-01` Definir estructura de mapa y tipos de nodo
- `ZP-02` Disenar progresion de dificultad entre zonas
- `ZP-03` Disenar sistema de recompensas entre encuentros
- `ZP-04` Externalizar catalogo de nodos y rutas del MVP

### Fish Data-Driven Expansion

- `FDE-01` Externalizar arquetipos y patrones de cartas de pez
- `FDE-02` Externalizar metadata editorial y de encounter de especies
- `FDE-03` Externalizar tablas de aparicion de peces por contexto de encounter

### Build and Services

- `BS-01` Definir categorias de objetos del jugador
- `BS-02` Disenar sistema de sinergias
- `BS-03` Implementar primer vertical slice de build
- `BS-04` Externalizar catalogo de servicios y acciones de build

### Economy and Meta Boundary

- `EMB-01` Definir economia ampliada de run y frontera futura de meta-progresion
- `EMB-02` Definir bestiario y coleccion
- `EMB-03` Disenar recompensas por coleccion
- `EMB-04` Disenar guardado de run y progreso meta
- `EMB-05` Externalizar recompensas y economia contextual
- `EMB-06` Externalizar catalogo de patrocinadores y ofertas

### Architecture and Delivery

- `AD-01` Mejorar UX de lectura de build y estado de run
- `AD-02` Desacoplar flujos de interfaz y preparar arquitectura UI-agnostic
- `AD-03` Centralizar constantes y politicas de balance del encounter
- `AD-04` Automatizar pipeline tecnico de calidad

### Graphical Client Prototype

- `GCP-01` Evaluar `Ebitengine` como base del cliente grafico
- `GCP-02` Disenar arquitectura de adaptador grafico sobre la app actual
- `GCP-03` Implementar composition root de prototipo grafico
- `GCP-04` Construir sistema minimo de escenas y navegacion visual
- `GCP-05` Definir direccion artistica minima del primer slice visual
- `GCP-06` Implementar HUD y tablero visual del encounter
- `GCP-07` Implementar presentacion visual de run, nodos y servicios
- `GCP-08` Definir pipeline minimo de assets y recursos visuales
- `GCP-09` Implementar animacion y feedback visual minimo del loop
- `GCP-10` Entregar primer vertical slice grafico jugable

## Labels recomendadas

- `stream:roguelike`
- `stream:gui`
- `stream:platform`
- `area:core-loop`
- `area:fish-encounters`
- `area:items-build`
- `area:economy-meta`
- `area:collection`
- `area:tech-ux`
- `area:graphics`
- `type:discovery`
- `type:delivery`
- `type:discovery+delivery`
- `type:quality`
- `type:infra`
- `type:quality+discovery`
- `type:infra+delivery`
- `priority:high`
- `priority:medium`
- `priority:low`
- `status:ready-for-spec`
- `status:in-spec`
- `status:planned`
- `status:ready-for-delivery`
- `status:blocked`

## Mapeo legacy -> nueva numeracion

| Legacy | Nuevo ID | Milestone |
| --- | --- | --- |
| `BL-036` | `RMF-01` | Run MVP Foundation |
| `BL-037` | `RMF-02` | Run MVP Foundation |
| `BL-038` | `RMF-03` | Run MVP Foundation |
| `BL-039` | `RMF-04` | Run MVP Foundation |
| `BL-040` | `RMF-05` | Run MVP Foundation |
| `BL-002` | `ZP-01` | Zone Progression |
| `BL-003` | `ZP-02` | Zone Progression |
| `BL-012` | `ZP-03` | Zone Progression |
| `BL-046` | `ZP-04` | Zone Progression |
| `BL-042` | `FDE-01` | Fish Data-Driven Expansion |
| `BL-043` | `FDE-02` | Fish Data-Driven Expansion |
| `BL-044` | `FDE-03` | Fish Data-Driven Expansion |
| `BL-008` | `BS-01` | Build and Services |
| `BL-009` | `BS-02` | Build and Services |
| `BL-010` | `BS-03` | Build and Services |
| `BL-047` | `BS-04` | Build and Services |
| `BL-011` | `EMB-01` | Economy and Meta Boundary |
| `BL-013` | `EMB-02` | Economy and Meta Boundary |
| `BL-014` | `EMB-03` | Economy and Meta Boundary |
| `BL-016` | `EMB-04` | Economy and Meta Boundary |
| `BL-048` | `EMB-05` | Economy and Meta Boundary |
| `BL-049` | `EMB-06` | Economy and Meta Boundary |
| `BL-017` | `AD-01` | Architecture and Delivery |
| `BL-023` | `AD-02` | Architecture and Delivery |
| `BL-031` | `AD-03` | Architecture and Delivery |
| `BL-032` | `AD-04` | Architecture and Delivery |
| `GUI-001` | `GCP-01` | Graphical Client Prototype |
| `GUI-002` | `GCP-02` | Graphical Client Prototype |
| `GUI-003` | `GCP-03` | Graphical Client Prototype |
| `GUI-004` | `GCP-04` | Graphical Client Prototype |
| `GUI-005` | `GCP-05` | Graphical Client Prototype |
| `GUI-006` | `GCP-06` | Graphical Client Prototype |
| `GUI-007` | `GCP-07` | Graphical Client Prototype |
| `GUI-008` | `GCP-08` | Graphical Client Prototype |
| `GUI-009` | `GCP-09` | Graphical Client Prototype |
| `GUI-010` | `GCP-10` | Graphical Client Prototype |

## Criterio de migracion

- Crear una issue por cada item abierto del mapa anterior.
- Usar el nuevo ID ordenado en el titulo, por ejemplo: `[RMF-01] Definir recursos globales de run...`.
- Conservar el ID legacy solo dentro del campo `Source context` o en una nota de migracion.
- Los items `done` y `cancelled` no se migran como issues activas.
- Toda nueva iniciativa debe nacer ya en GitHub Issues + `specs/`, sin volver a crear IDs legacy en `docs/backlog/`.
