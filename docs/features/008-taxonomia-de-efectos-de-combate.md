# Plan de feature: taxonomia-de-efectos-de-combate

## Objetivo

Definir una taxonomia clara de efectos de cartas de combate para dejar de crecer a partir de reglas aisladas y aplicar esa taxonomia en una primera implementacion tecnica de contratos compartidos para cartas del pez y del jugador.

El resultado esperado de esta feature es un documento de discovery utilizable como base para implementar un sistema de efectos de cartas agnostico al lado que juega la carta, alineado con lo que ya existe en `FishCard`, preparado para futuras `playerCards` y ya aterrizado en una primera capa de contratos comunes dentro del codigo.

## Criterios de aceptacion

- Existe un documento de referencia que clasifica los efectos de cartas por origen, trigger, fase de resolucion y tipo de impacto sobre el encounter.
- El documento mapea explicitamente las mecanicas ya implementadas para evitar contradicciones con la realidad actual del codigo.
- Queda explicito que el sistema futuro de efectos no depende de si la carta pertenece al pez o al jugador.
- El codigo migra los efectos actuales del pez a contratos de efecto de carta compartidos y relativos al owner de la carta.
- El sistema deja preparado un punto de extension para futuros efectos de carta que alteren outcome base desde `rules` mediante hooks.
- Quedan definidos limites y reglas de composicion suficientes para planear la siguiente feature de arquitectura sin ambiguedades importantes.

## Scope

### Incluye
- Documentar familias de efectos de cartas como movimiento horizontal, profundidad, recursos, mazo, estados, eventos y condiciones terminales.
- Definir una propuesta de fases de resolucion del round y en cuales fases puede actuar cada tipo de efecto de carta.
- Identificar origenes de efectos de carta actuales y futuros, al menos `FishCard` y futuras `playerCards`.
- Definir principios para que un mismo contrato de efecto pueda ser emitido y resuelto desde ambos lados del combate.
- Registrar restricciones de stacking, prioridades, conflictos y reglas de resolucion entre efectos.
- Implementar una primera version de contratos comunes para efectos de carta reutilizables por cartas del pez y del jugador.
- Migrar los efectos actuales del pez a triggers relativos al owner de la carta.
- Preparar en `rules` hooks para futuros efectos de carta que alteren el outcome base.
- Dejar recomendaciones concretas para convertir esta taxonomia en la proxima feature tecnica.

### No incluye
- Introducir nuevos efectos jugables o nuevas cartas en esta iteracion.
- Disenar en esta feature efectos que no provengan de cartas.
- Implementar aun un sistema completo de `playerCards` jugables dentro del loop actual.
- Redisenar la UI o la presentacion del combate mas alla de lo necesario para documentar ejemplos.

## Propuesta de implementacion

- Crear un documento de discovery principal que sirva como referencia para efectos de cartas y que use ejemplos sacados del estado actual del proyecto.
- Materializar ese documento en `docs/discoveries/008-taxonomia-de-efectos-de-carta.md` como salida principal de la feature.
- Partir desde la arquitectura existente en `internal/cards/`, `internal/progression/`, `internal/endings/`, `internal/encounter/` y `internal/playerrig/` para inventariar lo ya resuelto.
- Proponer una matriz simple con ejes de `emisor de carta -> trigger -> fase -> impacto -> restricciones` para evaluar nuevos efectos sin abrir decisiones desde cero cada vez.
- Cerrar una primera definicion de orden de resolucion que preserve el flujo actual del engine y deje un punto claro para futuros efectos emitidos tanto por `FishCard` como por `playerCards`.
- Implementar en `internal/cards/` contratos compartidos de efecto y tipos de carta separados por owner, sin forzar una carta comun unica.
- Ajustar `internal/game/`, `internal/match/` y `internal/progression/` para consumir los nuevos contratos sin cambiar el comportamiento jugable actual.
- Preparar `internal/rules/` con hooks optativos para futuros efectos de carta que alteren el resultado base de combate.
- Actualizar el backlog para que la siguiente feature natural sea la implementacion del sistema extensible derivado de esta taxonomia.

## Validacion

- Revisar manualmente que el documento cubra al menos los efectos de carta ya existentes: desplazamiento horizontal, desplazamiento vertical y splash, y que ubique correctamente su relacion con limites del rig y condiciones terminales.
- Verificar que `docs/discoveries/008-taxonomia-de-efectos-de-carta.md` deje una propuesta clara para la siguiente implementacion de `BL-005`.
- Verificar que desde el documento se pueda derivar sin dudas mayores una proxima feature de implementacion orientada a `BL-005`.
- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.

## Riesgos o decisiones abiertas

- Si la taxonomia intenta resolver demasiados casos futuros, puede volverse abstracta y perder utilidad practica para la siguiente feature.
- Si se documenta demasiado pegada a `FishCard`, puede dificultar la futura incorporacion simetrica de `playerCards`.
- Sigue abierta la decision de cuanto de la composicion futura debe modelarse como datos y cuanto como codigo especializado por tipo de efecto.
