# Plan de feature: resolver-agua-desde-el-nodo-de-run

## Objetivo

Extender `BL-035` para que el agua del encounter deje de elegirse manualmente en la run y pase a resolverse desde el nodo actual. La run debe definir el contexto de agua de cada nodo de pesca o boss, mientras el jugador solo recibe las senales visibles del agua que le toco jugar y luego resuelve el cast dentro de esa situacion.

En esta etapa no hace falta introducir todavia una generacion rica de aguas por zona. Alcanza con fijar una resolucion deterministica o prefijada por nodo dentro del ejemplo actual, pero dejando la arquitectura lista para que luego cada nodo pueda seguir reglas con algo de aleatoriedad entre runs.

## Criterios de aceptacion

- `cmd/fishing-run` ya no pide elegir manualmente el tipo de agua.
- Los nodos `fishing` y `boss` resuelven su `water context` desde la capa de run.
- El jugador sigue viendo los hints/senales del agua actual cuando entra al encounter.
- `cmd/fishing-duel` conserva la eleccion manual de agua como sandbox de prototipo.
- La ruta de run puede asociar un preset de agua por nodo, aunque sea fijo por ahora.
- El bootstrap de encounter soporta recibir un `water context` ya resuelto sin depender de `ChooseWaterContext(...)`.
- Queda preparada la frontera para introducir reglas de agua con aleatoriedad controlada mas adelante.

## Scope

### Incluye

- Agregar metadatos de agua a nodos de run del MVP.
- Permitir resolver apertura de encounter desde un preset de agua ya fijado.
- Rewirear `fishing-run` para usar esa resolucion de agua por nodo.
- Mantener intacto el flujo manual de `fishing-duel`.

### No incluye

- Generacion procedural completa de agua por zona.
- Catalogo externo de reglas de agua por nodo.
- Ocultacion avanzada o capas de intel mas complejas.

## Direccion propuesta

- `RunNode` define o referencia el preset de agua que le corresponde.
- La run resuelve ese preset antes de bootstrapear el encounter.
- El encounter solo consume el `water context` y muestra sus hints.
- La decision del jugador sigue estando en el cast y en el duelo, no en elegir el agua.

## Implementacion propuesta

### 1. Extender nodos de run

- Agregar un campo opcional de agua a `NodeState`, por ejemplo `WaterPresetID`.
- Asignar presets fijos a `fishing-1`, `fishing-2` y `boss-1` en la ruta MVP.

### 2. Extender el bootstrap de apertura

- Agregar una variante de `ResolveEncounterOpening(...)` que reciba un preset de agua ya resuelto.
- Mantener la variante manual existente para el prototipo.

### 3. Rewire de la run

- Hacer que `RunSession` pase el agua del nodo actual al bootstrap del encounter.
- El encounter de run ya no debe invocar `ChooseWaterContext(...)`.

### 4. Preparar la siguiente etapa

- Dejar helpers pequenos para que luego el agua pueda resolverse por reglas/pools y no solo por ID fijo.

## Accionables inmediatos

- Pasar `BL-035` a incluir este ajuste documentalmente.
- Implementar preset de agua por nodo en la ruta MVP.
- Separar bootstrap manual vs bootstrap desde nodo.
- Actualizar README si hace falta para aclarar que la run no pregunta el agua.
