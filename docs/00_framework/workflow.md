# Workflow de iteracion

## Objetivo

Desarrollar el MVP por capas sin perder pendientes ni mezclar decisiones de diseno con implementacion.

## Ciclo de trabajo recomendado

1. Escoger un tema del backlog.
2. Verificar dependencias.
3. Abrir su ficha o crearla desde `docs/_templates/`.
4. Resolver primero la intencion de diseno:
   - que fantasia cumple
   - que decision toma el jugador
   - que tension agrega
5. Resolver luego las reglas del sistema.
6. Identificar catalogos impactados.
7. Marcar preguntas abiertas.
8. Definir criterio de cierre.
9. Actualizar el backlog maestro.
10. Registrar la sesion.

## Orden recomendado para este proyecto

1. `01_foundations/03_run_loop.md`
2. `01_foundations/04_fail_success_conditions.md`
3. `03_content/01_zones.md`
4. `02_systems/01_casting.md`
5. `02_systems/02_encounter_flow.md`
6. `02_systems/03_distance_and_depth.md`
7. `02_systems/04_fish_decks_and_affinities.md`
8. `02_systems/05_reshuffle_and_exhaustion.md`
9. `03_content/02_fish_catalog.md`
10. `03_content/03_fish_abilities.md`
11. `02_systems/06_shop_and_economy.md`
12. `03_content/04_items_and_sponsorships.md`
13. `04_progression/01_run_progression.md`
14. `04_progression/02_meta_progression.md`
15. `04_progression/03_build_archetypes.md`
16. `05_presentation/01_functional_narrative.md`
17. `05_presentation/02_glossary.md`
18. `05_presentation/03_ui_information_model.md`
19. `05_presentation/04_mvp_exclusions.md`

## Regla de escalonado

No trabajar dos temas criticos a la vez. Si un tema depende de otro:

- queda `pendiente` si aun no se discute
- queda `en_discusion` si hay avances pero esta bloqueado por otra definicion

## Regla Balatro

Para cada tema preguntarse siempre:

- cual es la decision principal del jugador
- como se comunica el riesgo
- como alimenta la siguiente decision de build o progreso
- como escala sin romper claridad
