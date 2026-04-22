# Plan de feature: combate-empates-por-color

## Objetivo
Permitir que ciertos peces ganen empates cuando la carta del jugador coincide con un color en el que ese pez tiene ventaja, sin acoplar la solucion a una unica regla fija y dejando preparada la arquitectura para integrar futuras condiciones que alteren la resolucion del combate.

## Criterios de aceptacion
- Se puede definir, por pez, un conjunto de 0, 1, 2 o 3 colores en los que gana los empates.
- Si jugador y pez juegan el mismo color y ese color pertenece al conjunto de ventaja del pez, el resultado de la ronda es `FishWin`.
- Si jugador y pez juegan el mismo color y ese color no pertenece al conjunto de ventaja del pez, el resultado sigue siendo `Draw`.
- La regla circular actual (`Blue > Red`, `Red > Yellow`, `Yellow > Blue`) se mantiene igual para rondas que no dependan de una condicion especial.
- La solucion introduce un punto de extension explicito para futuras condiciones de combate sin obligar a reescribir el motor o la regla base cada vez.
- La suite de tests cubre la ventaja de empate por color y protege el comportamiento actual cuando un pez no tiene ventaja de empate configurada.

## Scope
### Incluye
- Modelar la ventaja de empate como parte de la configuracion o perfil de combate del pez.
- Refactorizar la resolucion de combate para soportar condiciones adicionales sobre la regla base.
- Implementar la primera condicion adicional: ventaja de empate por color.
- Ajustar bootstrap, fixtures y tests para que el evaluador reciba la configuracion necesaria.

### No incluye
- Un catalogo completo de especies seleccionables desde la CLI.
- Persistencia externa de perfiles de pez en JSON, YAML o similar.
- Nuevas condiciones de combate aparte de la ventaja de empate por color.
- Cambios en la presentacion visual mas alla de lo necesario para mantener coherencia funcional.

## Propuesta de implementacion
- Introducir una configuracion de combate del pez separada de la logica base del motor, con un conjunto de colores con ventaja en empate; el caso actual se representa con conjunto vacio para preservar compatibilidad.
- Refactorizar `internal/rules` para que la evaluacion del combate tenga una regla base y una capa de condiciones adicionales aplicables en orden.
- Implementar una interfaz o contrato de condicion de combate que permita agregar nuevas reglas especializadas sin reescribir `game.Engine` ni duplicar la logica circular existente.
- Implementar la condicion `tie advantage by color` como la primera condicion enchufable sobre esa capa.
- Tocar principalmente `internal/rules/`, el modelado compartido necesario en `internal/domain/` o un paquete dedicado de perfiles, y la composicion en `cmd/fishing-duel/main.go` y tests relacionados.

## Validacion
- Tests unitarios del evaluador para validar empate sin ventaja, empate con ventaja y victorias normales fuera de empate.
- Tests de integracion del motor o session para asegurar que la configuracion sin ventaja mantiene el comportamiento actual.
- Verificacion manual basica en la CLI si se expone temporalmente un pez de prueba con ventaja de empate configurada.

## Riesgos o decisiones abiertas
- Hay que decidir donde vive el perfil de combate del pez para no mezclar responsabilidades con encounter, deck o presentacion.
- Conviene definir desde ahora el orden de aplicacion de condiciones de combate para evitar ambiguedades cuando existan varias reglas especiales simultaneas.
- Si en el futuro algunas condiciones dependen de estado de ronda o del historial, la interfaz de condicion podria necesitar un contexto mas rico que solo `player`, `fish` y la configuracion del pez.
