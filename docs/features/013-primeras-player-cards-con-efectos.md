# Plan de feature: primeras-player-cards-con-efectos

## Objetivo

Validar de forma jugable y concreta que el sistema de efectos de cartas ya funciona de manera simetrica para ambos lados del combate, introduciendo una primera camada pequena de `playerCards` con efectos reales.

La meta de esta feature no es llenar el juego de contenido todavia, sino convertir la nueva infraestructura del jugador en una mecanica viva: las 3 barajas de decision por color deben dejar de ser solo una migracion estructural y pasar a demostrar que pueden soportar cartas con identidad, efectos por fase y presion tactica real dentro del flujo actual del duelo.

## Criterios de aceptacion

- Existe al menos un conjunto pequeno de `playerCards` con efectos reales integrados al juego.
- Se validan en gameplay efectos del jugador en mas de una fase, incluyendo al menos `on_draw` y un trigger post-outcome.
- El pipeline de efectos resuelve cartas del jugador y del pez en la misma ronda sin romper el comportamiento existente.
- La configuracion base puede seguir reproduciendo la experiencia clasica sin efectos, pero existen presets o configuraciones de prueba para ejercer las nuevas `playerCards`.
- La UI CLI sigue conservando el flujo actual de 3 decisiones, con la minima informacion adicional necesaria para probar la carta visible de cada color si hiciera falta.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope

### Incluye
- Definir una primera tanda de `playerCards` con efectos simples, legibles y utiles para validar el sistema.
- Integrar esos efectos en las barajas del jugador por color sin cambiar el flujo principal de seleccion.
- Habilitar presets o configuraciones de prueba para que desde CLI se puedan jugar encuentros con `playerCards` reales.
- Asegurar que el engine trate a `playerCards` y `FishCard` de forma consistente dentro del pipeline por fases.
- Ajustar la presentacion solo si es necesario para que el jugador pueda entender que carta visible esta usando cada color.
- Cubrir con tests automatizados la resolucion de efectos del jugador y su convivencia con efectos del pez.

### No incluye
- Implementar todavia construccion de mazo del jugador, recompensas o edicion libre de barajas.
- Introducir una gran libreria de cartas del jugador o balance profundo de contenido.
- Redisenar por completo la UI del duelo para mostrar cartas grandes o manos complejas.
- Entrar aun en arquetipos de pez o perfiles data-driven mas alla de lo necesario para probar esta feature.

## Propuesta de implementacion

- Definir un pequeno repertorio inicial de `playerCards` para los tres colores, priorizando efectos faciles de verificar y que aprovechen las herramientas ya construidas.
- Como punto de partida, conviene incluir al menos:
  - una carta del jugador con `on_draw` que altere thresholds temporales del round
  - una carta del jugador con `on_owner_win` que modifique distancia o profundidad
  - una carta del jugador con `on_owner_lose` o `on_round_draw` para validar respuestas defensivas o de tempo
- Mantener una configuracion clasica por defecto con 3 cartas iguales sin efectos, pero anadir presets de prueba controlados para el jugador, de forma similar a lo que ya existe para el pez.
- Extender `internal/playermoves/` para aceptar configuraciones concretas de barajas por color y no solo la baraja base uniforme.
- Ajustar `cmd/fishing-duel/main.go` y el flujo CLI para poder arrancar partidas de prueba donde tanto el pez como el jugador usen configuraciones preparadas para validar efectos.
- Si la legibilidad lo pide, enriquecer `internal/presentation/` e `internal/cli/` con una pista compacta sobre la carta visible del jugador por color, sin convertir la interfaz en una vista de mano completa.
- Documentar los ejemplos iniciales de `playerCards` de forma que luego puedan servir de semilla para sistemas mayores de build, items o sinergias.

## Validacion

- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Cubrir con tests al menos:
  - activacion de `on_draw` del jugador
  - activacion de un efecto post-outcome del jugador
  - resolucion conjunta de efectos de jugador y pez en una misma ronda
  - preservacion del comportamiento clasico cuando las barajas del jugador no tienen efectos
  - disponibilidad correcta de la carta visible tras consumir y revelar la siguiente carta del color
- Validar manualmente desde CLI al menos un preset donde las `playerCards` alteren de forma visible el round.

## Riesgos o decisiones abiertas

- Si metemos demasiadas cartas distintas en esta iteracion, la feature puede dispersarse entre validacion estructural y balance de contenido.
- Si no mostramos ninguna pista en UI sobre la carta visible del jugador, puede ser dificil verificar manualmente algunos efectos desde CLI.
- Habra que elegir efectos iniciales que sean potentes para probar el sistema pero lo bastante simples para no confundir la lectura del duelo.
