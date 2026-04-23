# internal/deck

Administracion de la baraja del pez.

## Responsabilidad

- Mantener mazo activo y descarte.
- Robar cartas del pez con su movimiento y modificadores.
- Rebarajar cuando el mazo se vacia.
- Aplicar una politica de reciclado.

## Puntos de extension

- `RecyclePolicy`: cambia como vuelve el descarte al mazo.
- `Shuffler`: cambia la estrategia de mezcla.
- `NewStandardFishCards()`: base para crear mazos predefinidos.
- `CustomFishDeck`: preset reusable para construir barajas de prueba o configuraciones especificas.

## Ejemplos de cambios futuros

- Peces con mazos de distinto tamano.
- Reciclados que retiren otra cantidad de cartas.
- Barajas con cartas raras o efectos especiales.

## Limites

- Este paquete no decide ganadores de ronda.
- Este paquete no decide captura ni escape.
