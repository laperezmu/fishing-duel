# Plan de feature: barajas-de-decision-del-jugador

## Objetivo

Reemplazar el modelo actual de usos planos por color por un sistema de 3 pequenas barajas de decision del jugador, una por cada color, donde el jugador siempre juega la carta superior disponible de la baraja elegida.

La intencion es mantener la experiencia de usuario actual casi intacta en el CLI mientras cambiamos la base mecanica a un modelo mucho mas expresivo: cada decision del jugador sigue siendo `Tirar`, `Recoger` o `Soltar`, pero ahora esa decision consume una carta concreta, la descarta y revela la siguiente, abriendo la puerta a efectos de carta especificos por color sin cambiar el flujo principal de combate.

Esta feature es una pieza clave para completar el sistema de efectos de cartas de ambos lados: despues de haber consolidado `FishCard`, `on_draw` y el pipeline por fases del pez, toca llevar esa misma riqueza al lado jugador con un modelo compatible con barajas por color.

## Criterios de aceptacion

- El jugador deja de usar un contador abstracto de usos por color y pasa a usar 3 barajas pequenas, una por cada color.
- Cada color expone siempre la carta superior actual de su baraja como la decision disponible para ese movimiento.
- Al jugar una decision, la carta superior correspondiente se descarta y se revela la siguiente carta de esa baraja.
- Cuando una baraja de color se vacia, ese color queda bloqueado durante 1 ronda completa y luego recupera su baraja barajada.
- La configuracion base sigue siendo equivalente a la experiencia actual: 3 cartas iguales por color, sin efectos.
- El sistema queda preparado para que futuras `playerCards` de un color puedan introducir efectos propios sin redisenar otra vez el flujo de seleccion.
- La UI actual sigue mostrando 3 decisiones, pero reflejando correctamente disponibilidad, agotamiento y recuperacion segun la nueva mecanica.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope

### Incluye
- Introducir un modelo de barajas del jugador por color dentro de `internal/playermoves/` o una capa cercana de dominio.
- Definir `playerCards` base sin efectos para los tres colores, con 3 cartas por color en la configuracion inicial.
- Actualizar el estado compartido de partida para reflejar la carta visible, el descarte y la recuperacion de cada baraja de decision.
- Adaptar la validacion y consumo de decisiones del jugador al nuevo sistema de cartas.
- Integrar el pipeline actual de efectos para que una `playerCard` pueda, en el futuro, emitir efectos del mismo modo que una `FishCard`.
- Mantener el flujo actual del CLI: el jugador sigue eligiendo entre 3 opciones numeradas, sin una mano adicional ni un paso extra de seleccion de carta.
- Cubrir con tests el agotamiento, la ronda de recuperacion, el rebarajado y la equivalencia del comportamiento base actual.

### No incluye
- Introducir aun recompensas, construccion de mazo o edicion de barajas del jugador fuera de la partida.
- Cambiar todavia la UI hacia una representacion visual completa de cartas del jugador.
- Anadir una gran variedad de `playerCards` con efectos complejos en la misma iteracion.
- Redisenar el sistema de fish deck o el loop roguelike completo.

## Propuesta de implementacion

- Crear un modelo explicito de `playerCards` por color que comparta contratos de efecto con `FishCard`, pero conserve su identidad como cartas del jugador.
- Evolucionar `internal/playermoves/` desde el sistema de `RemainingUses` hacia un sistema de barajas por color con:
  - mazo activo por color
  - descarte por color
  - carta visible actual por color
  - ronda de recuperacion al vaciarse
  - rebarajado al restaurarse
- Mantener la configuracion inicial equivalente a la actual: 3 cartas iguales sin efectos para azul, rojo y amarillo.
- Ajustar `internal/match/` para representar el nuevo estado del jugador sin perder la informacion que hoy consume la presentacion.
- Adaptar `internal/game/engine.go` para que el round pueda conocer tambien la carta de jugador utilizada, no solo el movimiento elegido y la `FishCard` robada.
- Preparar el pipeline para que una `playerCard` pueda activar efectos por fase de forma simetrica al pez, aunque en esta feature solo necesitemos validar bien el modelo base sin efectos.
- Mantener `internal/presentation/` e `internal/cli/` tan cercanos como sea posible al flujo actual, pero haciendo visible cuando una baraja esta agotada y cuando vuelve a estar disponible.
- Dejar una base clara para que la siguiente iteracion anada cartas del jugador con efectos reales sin cambiar el contrato central otra vez.

## Validacion

- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Cubrir con tests al menos:
  - inicializacion de 3 barajas base por color
  - revelado de la siguiente carta tras consumir la superior
  - bloqueo de una baraja vacia durante 1 ronda
  - recuperacion con rebarajado al cumplirse la ronda de espera
  - preservacion del flujo actual de seleccion de 3 decisiones en CLI
  - compatibilidad del comportamiento base con la experiencia actual cuando todas las cartas son equivalentes y no tienen efectos

## Riesgos o decisiones abiertas

- Habra que decidir cuanta informacion de carta visible mostramos ya en CLI sin sobrecargar una interfaz que hoy es muy compacta.
- Si la migracion reutiliza demasiado del modelo de usos actual, podemos terminar con una abstraccion hibrida dificil de extender.
- Si la feature intenta meter demasiados efectos reales de `playerCards` a la vez, puede mezclar la migracion estructural con balance y dificultar la validacion.
