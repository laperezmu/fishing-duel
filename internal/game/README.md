# internal/game

Motor central del juego.

## Responsabilidad

- Orquestar una ronda completa.
- Pedir una carta al mazo.
- Resolver fases de efectos de carta como draw y post-outcome.
- Consultar el evaluador de reglas.
- Aplicar la politica de progresion.
- Refrescar el snapshot del mazo.
- Aplicar la condicion de fin.

## Contratos principales

- `RoundEvaluator`: decide el resultado de la ronda.
- `FishDeck`: contrato del mazo que entrega cartas y expone su estado.
- `PlayerMoveController`: valida la decision del jugador, expone su carta visible y consume la carta usada.
- `MatchProgressionPolicy`: modifica el estado acumulado.
- `MatchEndCondition`: determina si el encuentro termina.

## Estructuras clave

- `Engine`: coordinador principal.
- El estado compartido y el resultado de ronda viven en `internal/match/`.
- La carta robada del pez vive en `internal/cards/` y las capacidades del jugador viven bajo `internal/player/`.

## Si quieres anadir una mecanica nueva

- Preguntate primero si es regla de ronda, progresion o condicion de fin.
- Intenta resolverla mediante inyeccion de politicas antes de modificar `Engine`.

## Limite importante

- Este paquete no conoce textos de UI ni nombres tematicos.
