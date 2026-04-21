# Plan de feature: limite-usos-movimientos-jugador

## Objetivo
Introducir un sistema de usos limitados para los movimientos del jugador que anada friccion y estrategia a cada partida, obligando a variar las decisiones de combate sin cambiar la regla base de resolucion entre jugador y pez.

La implementacion debe ser escalable, parametrizable y configurable para permitir futuros modificadores que alteren los usos iniciales o el ciclo de recuperacion sin reescribir la mecanica base.

## Criterios de aceptacion
- El jugador inicia cada partida con 3 usos disponibles de `Blue`, `Red` y `Yellow`.
- El valor `3` no queda hardcodeado como unica opcion del sistema; existe una configuracion explicita para definir usos iniciales por movimiento.
- Cada vez que el jugador usa un movimiento, ese movimiento consume 1 uso antes de la siguiente ronda.
- Cuando un movimiento agota su tercer uso en la ronda `N`, queda bloqueado durante toda la ronda `N+1`.
- El movimiento bloqueado recupera automaticamente sus 3 usos al inicio de la ronda `N+2`.
- La configuracion permite evolucionar hacia peces, partidas o modificadores que cambien los usos iniciales o la recuperacion sin romper la arquitectura.
- El motor impide jugar un movimiento sin usos disponibles o en cooldown; la restriccion no depende solo de la UI.
- La presentacion y la CLI muestran con claridad los usos restantes y el estado de bloqueo o recarga de cada movimiento.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope
### Incluye
- Modelar en el estado compartido de la partida el recurso de usos del jugador por movimiento.
- Modelar una configuracion explicita del sistema de recursos del jugador, en vez de quemar valores fijos en el motor.
- Aplicar la logica de consumo, bloqueo de una ronda y recarga completa en el motor o en una politica dedicada de recursos del jugador.
- Ajustar presentacion y CLI para reflejar movimientos disponibles, agotados o temporalmente bloqueados.
- Agregar o actualizar tests unitarios y de integracion para cubrir el nuevo ciclo de usos y recuperacion.

### No incluye
- Cambiar la regla base de combate `Blue > Red > Yellow > Blue`.
- Cambiar la baraja del pez o su sistema de reciclado.
- Introducir nuevas cartas, efectos especiales o habilidades de pez adicionales en esta iteracion.
- Agregar persistencia externa o configuracion por archivo para los usos del jugador.
- Implementar aun modificadores concretos que alteren los usos iniciales, mas alla de dejar la mecanica preparada para soportarlos.

## Propuesta de implementacion
- Extender `internal/match` con un modelo explicito del estado de recursos del jugador por movimiento, incluyendo usos restantes y round de recarga cuando aplique.
- Introducir una configuracion dedicada para los recursos del jugador, con defaults equivalentes a 3 usos por movimiento y una ronda completa de espera antes de recargar.
- Mantener la validacion de disponibilidad en el motor para que una jugada invalida no pueda resolverse aunque la UI falle en filtrarla.
- Introducir una capa dedicada para administrar el ciclo de consumo y recuperacion del jugador, evitando mezclar esta responsabilidad con la progresion del track del pez.
- Actualizar `internal/presentation` para exponer el estado de cada accion y `internal/cli` para comunicar con claridad cuando una accion esta disponible, agotada o en espera de recarga.
- Ajustar el flujo de sesion y los tests necesarios para manejar jugadas rechazadas o no disponibles de forma consistente.

## Validacion
- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Cubrir con tests al menos: consumo normal, agotamiento en el tercer uso, bloqueo de una ronda completa, recarga total en la ronda siguiente y rechazo de jugadas no disponibles.
- Verificar manualmente en la CLI que los usos visibles se actualizan ronda a ronda y que la recuperacion ocurre con la semantica acordada.
- Verificar que no aparezcan comentarios `//nolint` como salida rapida para cerrar la iteracion.

## Riesgos o decisiones abiertas
- Habra que decidir si el rechazo de una jugada invalida devuelve error de motor, fuerza una nueva seleccion en la sesion o ambas cosas mediante capas separadas.
- Si la logica de recursos del jugador se acopla demasiado al motor, podria dificultar futuras mecanicas similares; conviene dejar un punto de extension claro.
- La UX de la CLI debe informar bien el estado de cada accion para que la nueva friccion se sienta estrategica y no confusa.
