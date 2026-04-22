# Plan de feature: verticalidad-y-efectos-de-carta

## Objetivo
Agregar un segundo eje espacial al encounter mediante profundidad o verticalidad del pez, de forma que ciertas cartas del pez puedan modificarla al ganar una ronda y abrir un primer sistema extensible de efectos de carta preparado para futuros efectos de items y cartas del jugador.

La arquitectura debe separar con claridad las condiciones del entorno, el comportamiento del pez y los limites operativos de las herramientas del jugador, de modo que la distancia maxima alcanzable y la profundidad maxima alcanzable dependan del loadout o rig del jugador y no del pez.

## Criterios de aceptacion
- El encounter modela de forma explicita un estado de profundidad independiente de la distancia horizontal.
- El sistema separa al menos tres fuentes de configuracion: encounter o entorno, perfil del pez y rig o herramientas del jugador.
- La superficie pertenece al encounter o entorno, mientras que la distancia maxima alcanzable y la profundidad maxima alcanzable dependen del rig o herramientas del jugador.
- El sistema deja de asumir que una carta del pez es solo un `Move`; existe una representacion preparada para adjuntar efectos configurables a cartas del pez.
- Se puede configurar al menos un primer efecto de carta del pez: al ganar la ronda, desplaza la profundidad del pez hacia arriba o hacia abajo por una magnitud configurable.
- Si la profundidad del pez supera la profundidad maxima alcanzable por el jugador, el encounter termina en escape con una razon terminal explicita.
- Si la distancia horizontal del pez supera la distancia maxima alcanzable por el jugador, el encounter sigue terminando en escape, pero ese limite queda asociado a las herramientas del jugador y no al pez.
- Si el pez ya esta en superficie y un efecto intenta hacerlo subir otra vez, se dispara un evento de chapoteo que resuelve un escape con probabilidad configurable.
- La arquitectura deja preparada la composicion futura de efectos provenientes no solo de cartas del pez, sino tambien de cartas o items del jugador.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope
### Incluye
- Extender el modelo de `encounter` para soportar profundidad como nuevo eje del combate.
- Introducir un modelo explicito para las capacidades del jugador o su rig de pesca, incluyendo al menos alcance horizontal y profundidad maxima alcanzable.
- Introducir una representacion mas rica de carta del pez que pueda describir efectos configurables ademas del movimiento base.
- Implementar el primer modificador de carta: desplazamiento vertical al ganar.
- Resolver nuevos finales de encounter asociados a las capacidades del jugador y al nuevo eje: escape por profundidad maxima alcanzable, escape por distancia horizontal maxima y posible escape por chapoteo en superficie.
- Ajustar el motor, progresion, finales, presentacion y CLI para reflejar el nuevo eje y sus eventos.
- Agregar tests unitarios e integracion para profundidad, chapoteo y efectos de carta.

### No incluye
- Implementar aun items del jugador con efectos activos durante el combate.
- Reemplazar completamente el sistema actual de reglas de combate.
- Agregar un framework generico y abierto para cualquier tipo de scripting de efectos.
- Introducir aun multiples tipos de efectos de carta fuera del primer desplazamiento vertical al ganar.

## Propuesta de implementacion
- Extender `internal/encounter` con profundidad como estado real del pez y con configuracion propia del entorno, manteniendo ahi conceptos globales como superficie y razones terminales del encounter.
- Introducir un modelo o config de rig del jugador que defina limites operativos como distancia maxima alcanzable, profundidad maxima alcanzable y la base modificable del evento de chapoteo.
- Evolucionar la representacion del mazo del pez para que robe una carta de encounter mas rica que contenga el `Move` base y una lista o descriptor de efectos configurables, en lugar de solo `domain.Move`.
- Mantener `internal/rules` centrado en resolver el resultado base de la ronda, y mover la aplicacion de efectos espaciales al flujo de progresion o a una capa de efectos de encounter preparada para crecer.
- Introducir una capa o contrato de efectos de carta por fases del round que, en esta primera iteracion, soporte modificar la progresion del encounter y disparar eventos derivados como chapoteo.
- Mantener `internal/game` como orquestador del flujo: validar jugada, robar carta del pez, resolver round, calcular progresion base, aplicar efectos de carta, refrescar estado y evaluar condiciones terminales en funcion del encounter y de las capacidades del jugador.
- Ajustar `internal/presentation` y `internal/cli` para mostrar profundidad actual, cambios verticales y resultado del evento de chapoteo sin mezclar la logica visual con la resolucion del efecto.
- Dejar la capa de efectos desacoplada de la fuente concreta del efecto para que en futuras iteraciones pueda recibir tambien modificadores provenientes de cartas o items del jugador.

## Validacion
- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Cubrir con tests al menos: configuracion valida del nuevo eje, efecto vertical al ganar, escape por profundidad maxima alcanzable, preservacion del escape por distancia maxima alcanzable, evento de chapoteo en superficie y preservacion del flujo actual cuando ninguna carta tenga efecto vertical.
- Verificar manualmente en la CLI que se muestran tanto distancia como profundidad y que el chapoteo se comunica de forma clara.
- Verificar que no aparezcan comentarios `//nolint` como salida rapida para cerrar la iteracion.

## Riesgos o decisiones abiertas
- Habra que decidir si la nueva representacion de carta del pez vive en `domain`, `deck` o en un paquete de combate mas especifico para no mezclar conceptos demasiado pronto.
- Si el primer sistema de efectos se vuelve demasiado generico desde el inicio, puede sobredisenarse; si se hace demasiado especifico al pez, luego costara reutilizarlo para items del jugador.
- El orden exacto entre progresion base, efecto vertical y chequeos terminales debe quedar bien definido para evitar ambiguedades en edge cases.
- El evento de chapoteo introduce azar en un sistema hasta ahora casi totalmente determinista; la inyeccion de su resolucion probabilistica debe quedar desacoplada para poder testearla con control.
