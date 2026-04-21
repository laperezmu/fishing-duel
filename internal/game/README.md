# internal/game

Motor central del juego.

## Responsabilidad

- Orquestar una ronda completa.
- Pedir una carta al mazo.
- Consultar el evaluador de reglas.
- Aplicar la politica de progresion.
- Refrescar el snapshot del mazo.
- Aplicar la condicion de fin.

## Contratos principales

- `Evaluator`: decide el resultado de la ronda.
- `ProgressionPolicy`: modifica el estado acumulado.
- `EndCondition`: determina si el encuentro termina.

## Estructuras clave

- `State`: estado completo del juego.
- `RoundResult`: snapshot util tras una ronda.
- `Engine`: coordinador principal.

## Si quieres anadir una mecanica nueva

- Preguntate primero si es regla de ronda, progresion o condicion de fin.
- Intenta resolverla mediante inyeccion de politicas antes de modificar `Engine`.

## Limite importante

- Este paquete no conoce textos de UI ni nombres tematicos.
