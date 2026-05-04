# cmd/fishing-run

Punto de entrada del ejecutable principal de la run MVP.

## Responsabilidad

- Componer la sesion de run.
- Resolver una seleccion inicial de `AnglerProfile` una sola vez al comienzo de la expedicion.
- Recorrer la secuencia fija de nodos del MVP, hoy extendida a 8 fases de agua.
- Disparar encounters desde nodos de pesca y boss.
- Mostrar el resumen final de la expedicion.

## Estructura actual de la ruta

- 24 combates aproximados en total.
- 8 bosses en total, con 1 boss cada 3 combates.
- 8 fases encadenadas; cada fase cambia el preset de agua y prueba una combinacion distinta de pesca.

## Que no debe hacer

- No debe contener reglas del encounter.
- No debe decidir logica de nodos fuera de la capa `internal/app/` y `internal/run/`.
- No debe reemplazar al sandbox manual de `cmd/fishing-duel/`.
