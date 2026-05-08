# Contract: CLI splash flow

## Purpose

Definir el contrato funcional entre `internal/app`, `internal/presentation` e `internal/cli` para resolver un splash interactivo sin filtrar reglas del encounter dentro del adaptador de terminal.

## Inputs to UI

La UI debe recibir una vista de splash que, como minimo, exponga:

- etiqueta o titulo del evento
- salto actual
- total de saltos
- tiempo disponible del salto actual
- indicador visible de progreso
- mensaje de recompensa por salto, si existe

## UI responsibilities

- Mostrar el estado del splash antes de iniciar cada salto.
- Ejecutar la interaccion de timing del salto actual.
- Devolver un resultado binario por salto: exito o fallo por timeout/error.
- Repetir el flujo hasta que app indique que la secuencia termino o el pez escapo.

## App responsibilities

- Detectar que el engine dejo un splash pendiente.
- Pedir a presentation una vista de splash legible para CLI.
- Solicitar a UI la resolucion del salto o secuencia.
- Aplicar el resultado al engine/runtime.
- Continuar el flujo normal de ronda cuando el splash termina sin escape.

## Engine/runtime responsibilities

- Exponer el splash pendiente como estado del encounter o snapshot.
- Aplicar recompensa de distancia por salto ganado cuando corresponda.
- Marcar escape de splash como motivo terminal si hay fallo.
- Mantener snapshots finales consistentes para round summary y game over.

## CLI behavior guarantees

- El flujo sigue siendo jugable solo con terminal e input de teclado.
- Un timeout se interpreta como fallo del salto.
- Un splash de varios saltos comunica claramente progreso y resultado parcial.
- La salida de CLI conserva compatibilidad con los resumenes actuales de round y encounter.
