# Plan de feature: reducir-determinismo-del-spawn-de-peces

## Objetivo

Reducir el determinismo absoluto del spawn de peces sin depender todavia de una capa data-driven completa. La meta es que el sistema deje de ser totalmente deducible cuando el jugador conoce agua, distancia, profundidad y habitats, pero sin perder legibilidad, control de balance, reproducibilidad en tests ni compatibilidad futura con catalogos y pools externas.

Esta feature no busca resolver aun el catalogo data-driven de peces ni el sistema final de nodos. Busca mejorar la estrategia de seleccion dentro del resolver actual para que la variedad quede modelada como una capacidad del dominio y no como un efecto accidental del wiring futuro.

## Criterios de aceptacion

- El spawn deja de elegir siempre el mismo perfil valido bajo las mismas condiciones salvo que la configuracion lo fuerce explicitamente.
- La variacion sigue siendo reproducible cuando se proporciona seed o randomizer controlado.
- El algoritmo mantiene reglas legibles y tuneables para balancear variedad frente a coherencia contextual.
- La solucion no rompe el flujo actual `agua -> apertura -> spawn -> mazo del pez`.
- La implementacion deja preparada una superficie clara para que pesos o reglas futuras puedan venir luego desde data-driven.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope

### Incluye

- Revisar el algoritmo actual de `ResolveSpawn(...)` y detectar donde se introduce el determinismo absoluto.
- Definir una estrategia minima de variedad controlada dentro de candidatos validos.
- Introducir un mecanismo reproducible de seleccion o desempate que pueda usar randomizer/seed.
- Ajustar tests para cubrir variedad controlada sin volverlos flaky.
- Documentar la direccion acordada si el cambio altera expectativas del spawn.

### No incluye

- Externalizar aun catalogos a JSON.
- Implementar todavia pools por nodo o zonas data-driven completas.
- Cambiar el runtime tactico del encounter.
- Introducir sistemas opacos de ML, heuristicas emergentes complejas o tuning dificil de depurar.

## Contexto actual

Hoy el spawn:

- recibe un slice de perfiles validos y filtra por `water pool`, distancia, profundidad y habitats
- calcula un `MatchScore` por perfil
- ordena candidatos por score y luego por indice original
- elige siempre el primero

Esto hace que, bajo el mismo contexto, el resultado sea completamente deducible. La unica variedad actual depende de cambiar manualmente el contexto o el catalogo pasado al resolver.

## Propuesta de implementacion

### 1. Separar matching de seleccion final

El primer paso es distinguir mejor dos responsabilidades:

- determinar que candidatos son validos
- decidir cual se selecciona entre esos candidatos

Direccion propuesta:

- mantener el matching actual como capa de elegibilidad
- extraer la seleccion final a una funcion o politica explicita
- dejar esa politica lista para evolucionar despues a configuracion data-driven

### 2. Introducir variedad controlada entre candidatos validos

La primera version no necesita un sistema complejo. Debe ser simple, legible y facil de testear.

Opciones razonables para el primer slice:

- seleccionar aleatoriamente entre todos los candidatos con mejor score
- o seleccionar aleatoriamente entre una banda corta de candidatos cercanos en score
- o aplicar pesos derivados del score y samplear de forma reproducible

Recomendacion para esta feature:

- empezar por `random among top score ties` o una banda muy acotada de scores cercanos
- evitar todavia pesos muy sofisticados si no mejoran claramente la expresividad

### 3. Mantener reproducibilidad explicita

La variedad solo sirve si sigue siendo depurable.

Direccion propuesta:

- pasar un randomizer o dependencia equivalente al punto de seleccion
- permitir semillas controladas en tests
- evitar apoyarse en random global no inyectado

### 4. Preparar la transicion futura a data-driven

Aunque esta feature no implemente JSON, si debe dejar el dominio preparado para ello.

Direccion propuesta:

- nombrar conceptos como `selection policy`, `weight`, `candidate pool` o equivalente
- no hardcodear decisiones irreversibles dentro de `ResolveSpawn(...)`
- dejar facil que `BL-041`/`BL-044` puedan despues inyectar pesos o reglas por pool sin reescribir la seleccion

### 5. Ajustar la estrategia de tests

El cambio no debe volver fragil la suite.

Direccion propuesta:

- testear matching y scoring por separado de la seleccion final
- testear la seleccion con randomizer fake o seed fija
- cubrir casos como empate, score dominante unico y banda corta de candidatos validos

## Slice minimo recomendado

El primer slice deberia cubrir:

- extraccion de la politica de seleccion final del spawn
- seleccion reproducible con randomizer inyectado
- variedad minima entre candidatos empatados o muy cercanos
- tests nuevos para esa politica

Con eso ya se reduce el determinismo total sin bloquear el trabajo posterior de catalogos y pools data-driven.

## Implementacion aplicada

- `ResolveSpawn(...)` mantiene un camino estable y deterministicamente ordenado cuando no se inyecta randomizer.
- `ResolveSpawnWithRandomizer(...)` permite introducir variedad controlada entre candidatos empatados en score sin alterar el matching contextual.
- `app` ya usa ese camino reproducible con randomizer inyectado desde el bootstrap actual del encounter, preparando la futura integracion con pools y catalogos externos.

## Validacion

- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Verificar manualmente que encuentros repetidos sobre el mismo contexto ya no devuelven siempre el mismo pez cuando hay mas de un candidato razonable.
- Verificar que con seed fija el resultado siga siendo reproducible.

## Riesgos o decisiones abiertas

- Si la banda de candidatos validos es demasiado amplia, la fantasia contextual del agua puede perder fuerza.
- Si la banda es demasiado estrecha, la variedad puede seguir sintiendose irrelevante.
- Habra que decidir si la primera version varia solo empates exactos o tambien scores cercanos.
- Conviene evitar que la aleatoriedad vuelva opaco el balance o dificulte comprender por que un pez era elegible.
