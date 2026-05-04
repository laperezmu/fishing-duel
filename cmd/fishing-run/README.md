# cmd/fishing-run

Punto de entrada del ejecutable principal de la run MVP.

## Responsabilidad

- Componer la sesion de run.
- Resolver una seleccion inicial de `AnglerProfile` una sola vez al comienzo de la expedicion.
- Recorrer la secuencia fija de nodos del MVP.
- Disparar encounters desde nodos de pesca y boss.
- Mostrar el resumen final de la expedicion.

## Que no debe hacer

- No debe contener reglas del encounter.
- No debe decidir logica de nodos fuera de la capa `internal/app/` y `internal/run/`.
- No debe reemplazar al sandbox manual de `cmd/fishing-duel/`.
