# Plan de feature: effective-go-naming-y-estructura-paquetes

## Objetivo
Mejorar la calidad del codigo del proyecto alineandolo con practicas de Effective Go mediante nombres mas descriptivos a lo largo de codigo productivo y tests, eliminacion de identificadores ambiguos o excesivamente genericos, y una refactorizacion estructural fuerte de paquetes para reforzar cohesion, direccion de dependencias y legibilidad general.

## Criterios de aceptacion
- No quedan variables, parametros, campos, helpers o tipos con nombres ambiguos o excesivamente genericos cuando exista una alternativa mas descriptiva y razonable.
- Las excepciones idiomaticas de Go (`err`, `ok`, `t`, indices cortos y receivers breves) se mantienen solo donde sigan siendo claras y estandar.
- Los nombres refactorizados mejoran la semantica del codigo sin alterar comportamiento funcional.
- La estructura de paquetes reduce acoplamientos innecesarios y mejora la direccion de dependencias entre modelo, orquestacion, reglas, progreso, finales, presentacion y app.
- La refactorizacion estructural extrae el modelo compartido de runtime fuera de `internal/game` para que `internal/game` quede enfocado en la orquestacion del motor.
- Se eliminan fugas de presentacion desde capas base cuando formen parte del refactor aprobado.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope
### Incluye
- Renombrado descriptivo de identificadores productivos y de testing.
- Revision y mejora de nombres de tipos, contratos e interfaces demasiado genericos.
- Refactorizacion estructural fuerte para separar el modelo compartido de runtime del paquete `internal/game`.
- Ajuste de dependencias, imports, tests y documentacion afectados por la nueva estructura.
- Revision de fugas de presentacion en capas no destinadas a UI.

### No incluye
- Cambios funcionales de gameplay.
- Introduccion de nuevas features de dominio.
- Reescritura total del proyecto fuera de los cambios necesarios para mejorar naming y estructura.
- Cambios esteticos sin mejora clara de legibilidad, mantenibilidad o direccion arquitectonica.

## Propuesta de implementacion
- Hacer un barrido de identificadores ambiguos o genericos y renombrarlos segun su semantica real.
- Priorizar contratos y tipos visibles entre paquetes como `Manager`, `Condition`, `EncounterCondition`, parametros `player` y `fish`, campos de contexto y helpers de test.
- Crear un paquete de modelo compartido de runtime para mover tipos como el estado del juego, resultado de ronda y estructuras relacionadas fuera de `internal/game`.
- Actualizar `internal/game` para que quede centrado en `Engine` y sus contratos de orquestacion, consumiendo el nuevo modelo compartido en lugar de definirlo internamente.
- Adaptar `progression`, `endings`, `presentation`, `app`, `cli` y bootstrap para depender del nuevo paquete de modelo segun corresponda.
- Revisar `internal/domain` y `internal/cli` para evitar dependencias de strings de presentacion donde deban usarse tipos semanticos.
- Mantener el refactor incremental y verificable, cuidando que cada renombrado o movimiento siga siendo mecanico y validable por tests y lint.

## Validacion
- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Verificar que los nombres finales mejoren semantica y lectura, no solo longitud.
- Verificar que las dependencias entre paquetes queden mas claras tras mover el modelo compartido fuera de `internal/game`.
- Verificar que no aparezcan comentarios `//nolint` como salida rapida para cerrar la iteracion.

## Riesgos o decisiones abiertas
- Mezclar un barrido amplio de naming con una refactorizacion estructural fuerte puede crecer rapidamente el alcance si no se mantiene disciplina de corte.
- Mover tipos compartidos fuera de `internal/game` tocara varios paquetes al mismo tiempo y puede generar churn amplio aunque sea mecanico.
- Corregir fugas de presentacion en capas bajas puede requerir pequenos ajustes coordinados en `presentation`, `cli` y tests.
- Habra que vigilar que los nuevos nombres sean realmente mejores y no solo mas largos.
