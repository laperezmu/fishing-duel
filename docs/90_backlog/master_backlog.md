# Backlog maestro de diseno

## Como usarlo

- No borres temas; cambia su estado.
- Trabaja solo un tema critico a la vez.
- Cuando cierres una ficha, actualiza esta tabla.
- Si surge una idea fuera de alcance, muevela a `parking_lot.md`.

## Prioridad critica

| Tema | Estado | Depende de | Archivo |
| --- | --- | --- | --- |
| Vision del juego | en_discusion | - | `docs/01_foundations/01_game_vision.md` |
| Pilares de diseno | en_discusion | - | `docs/01_foundations/02_design_pillars.md` |
| Loop principal de run | pendiente | Vision, Pilares | `docs/01_foundations/03_run_loop.md` |
| Condiciones de progreso y fracaso | pendiente | Loop principal | `docs/01_foundations/04_fail_success_conditions.md` |
| Catalogo de zonas | pendiente | Loop principal, Condiciones | `docs/03_content/01_zones.md` |
| Sistema de lances | pendiente | Loop principal | `docs/02_systems/01_casting.md` |
| Flujo del encuentro | pendiente | Lances | `docs/02_systems/02_encounter_flow.md` |

## Prioridad alta

| Tema | Estado | Depende de | Archivo |
| --- | --- | --- | --- |
| Tracks de distancia y profundidad | pendiente | Flujo del encuentro | `docs/02_systems/03_distance_and_depth.md` |
| Mazos del pez y afinidades | pendiente | Flujo del encuentro | `docs/02_systems/04_fish_decks_and_affinities.md` |
| Reshuffle y agotamiento | pendiente | Mazos del pez | `docs/02_systems/05_reshuffle_and_exhaustion.md` |
| Catalogo de peces | pendiente | Zonas, Mazos del pez | `docs/03_content/02_fish_catalog.md` |
| Catalogo de habilidades del pez | pendiente | Catalogo de peces | `docs/03_content/03_fish_abilities.md` |
| Progresion de run | pendiente | Condiciones, Zonas | `docs/04_progression/01_run_progression.md` |
| Tienda y economia de run | pendiente | Loop principal, Progresion de run | `docs/02_systems/06_shop_and_economy.md` |
| Catalogo de objetos y patrocinios | pendiente | Tienda y economia | `docs/03_content/04_items_and_sponsorships.md` |
| Glosario | pendiente | - | `docs/05_presentation/02_glossary.md` |
| Modelo de informacion visible | pendiente | - | `docs/05_presentation/03_ui_information_model.md` |

## Prioridad media

| Tema | Estado | Depende de | Archivo |
| --- | --- | --- | --- |
| Legendarios y trofeos | pendiente | Zonas, Peces, Meta | `docs/03_content/05_legendaries_and_trophies.md` |
| Meta progresion | pendiente | Condiciones, Legendarios | `docs/04_progression/02_meta_progression.md` |
| Arquetipos de build | pendiente | Objetos | `docs/04_progression/03_build_archetypes.md` |
| Narrativa funcional | pendiente | Vision del juego | `docs/05_presentation/01_functional_narrative.md` |
| Exclusiones del MVP | pendiente | - | `docs/05_presentation/04_mvp_exclusions.md` |

## Siguiente recomendacion de trabajo

1. Cerrar `docs/01_foundations/01_game_vision.md`
2. Cerrar `docs/01_foundations/02_design_pillars.md`
3. Pasar a `docs/01_foundations/03_run_loop.md`
