# Plan de implementacion - simulacion de combate en Go

## Objetivo
Construir un modulo de simulacion del combate de pesca en Go para ejecutar corridas Monte Carlo y responder preguntas de balance como:

- probabilidad de captura por tipo de pez
- probabilidad de escape por slot inicial 3, 4 y 5
- impacto de carnadas y mejoras
- duracion esperada del combate
- distribucion de fatigas
- valor esperado de cada accion del jugador en cada estado

La primera version debe aislar el subsistema de combate. No necesita modelar movimiento, mapa, almacen ni puntuacion global de la partida.

---

## Alcance funcional de la v1
La simulacion debe cubrir solo lo que afecta directamente el combate:

- posicion inicial del track: slot 3, 4 o 5
- tipo de pez: negro, monocolor o bicolor
- colores/arquetipos del pez
- mazo de comportamiento del pez como pila real
- regla de arquetipo
- fatiga del pez
- fatiga total
- carnada compatible que evita el primer escape
- mejoras de cana ofensivas
- politica de decision del jugador
- motor de simulacion repetible por semilla

Queda fuera de la v1:

- eleccion de zona
- sistema de lanzamiento con cartas de zona
- senuelo explorador
- economia de inventario
- eventos y puntos finales

---

## Analisis del problema

### 1. El combate es un proceso secuencial con informacion imperfecta
Cada ronda de combate ocurre sobre un estado observable parcial:

- el jugador conoce el slot actual del pez en el track
- conoce el tipo y color del pez
- conoce su propio equipamiento y si la carnada sigue activa
- no conoce la siguiente carta exacta del pez
- si puede conocer la composicion restante del mazo y del descarte

Por lo tanto, el problema del jugador no es de reaccion a una carta puntual sino de decision bajo distribucion de probabilidad.

### 2. El pez no necesita una IA activa
El pez no toma decisiones estrategicas. Su comportamiento emerge de una baraja fija que se roba en orden. La forma correcta de modelarlo es:

- un mazo barajado al inicio del combate
- una pila `draw` de la que se roba la carta superior
- una pila `discard` a la que van las cartas resueltas
- al agotarse `draw`, se baraja `discard`, se retiran 3 cartas del tope y se forma un nuevo `draw`
- si no se puede formar un nuevo mazo tras la fatiga, el pez es capturado automaticamente

Esto hace que la incertidumbre del pez sea muestreal y no heuristica.

### 3. La estrategia optima del jugador no deberia ser una tabla fija manual
Si se define a mano un "perfil optimo" se corre el riesgo de codificar sesgos falsos. Conviene separar:

- `PlayerPolicy`: interfaz generica para elegir accion
- `OptimalPolicy`: politica calculada por el sistema a partir del estado
- `HeuristicPolicy`: reglas simples para comparar rendimiento
- `RandomPolicy`: baseline de control

La politica optima debe elegir la accion que maximiza una funcion objetivo, por ejemplo:

- maximizar probabilidad de captura
- o maximizar valor esperado de resultado final

Para la v1 recomendamos optimizar `probabilidad de captura`, y dejar el objetivo configurable para futuras iteraciones.

### 4. El espacio de estados es finito y pequeno
Aunque el combate tiene incertidumbre, el numero de estados es manejable si se define bien. Un estado puede describirse con:

- posicion del track
- colores del pez
- composicion exacta de `draw`
- composicion exacta de `discard`
- numero de fatigas ocurridas
- si la proteccion de carnada sigue disponible
- mejoras activas del jugador

Con esto se puede resolver la accion optima por recursion con memoizacion sobre valor esperado, sin depender solo de Monte Carlo.

### 5. Monte Carlo sirve para validacion y agregados
Aunque la politica optima puede calcularse por esperanza exacta dado un estado, Monte Carlo sigue siendo util para:

- promediar resultados sobre miles de barajados iniciales
- comparar loadouts
- estimar sensibilidad de balance
- validar que la politica optima se comporta como esperamos
- analizar escenarios completos de combate sin enumerar manualmente todos los casos

---

## Decision clave: dos capas de simulacion

### Capa A - Evaluador exacto del estado
Componente que recibe un estado de combate y responde:

- valor esperado de captura desde ese estado
- mejor accion del jugador
- desglose por accion posible

Esta capa sirve para construir el perfil de estrategia optima.

### Capa B - Runner Monte Carlo
Componente que ejecuta muchos combates completos muestreando el mazo inicial del pez y usando una politica para el jugador.

Esta capa sirve para producir metricas agregadas.

La recomendacion es implementar primero la capa A y luego la capa B. Asi la simulacion Monte Carlo ya puede apoyarse en una politica fuerte y no en reglas improvisadas.

---

## Modelo de dominio propuesto

### Enumeraciones base
- `Action`: `Forzar`, `Tensar`, `Soltar`
- `FishFamily`: `Embiste`, `Aguante`, `Quiebre`
- `Color`: `None`, `Red`, `Blue`, `Yellow`
- `FishKind`: `Black`, `Mono`, `BiColor`
- `Outcome`: `Capture`, `Escape`, `Ongoing`

Mapeos de color del reglamento:

- `Forzar` <-> `Embiste` <-> `Red`
- `Tensar` <-> `Aguante` <-> `Blue`
- `Soltar` <-> `Quiebre` <-> `Yellow`

### Estructuras principales

```go
type CombatState struct {
    TrackPos            int
    Fish                 FishProfile
    Deck                 FishDeckState
    FatigueCount         int
    BaitGuardAvailable   bool
    RodMods              RodModifiers
}

type FishProfile struct {
    Kind   FishKind
    Colors [2]Color
}

type FishDeckState struct {
    Draw    []FishFamily
    Discard []FishFamily
}

type RodModifiers struct {
    Offensive map[Color]bool
}
```

### Estado derivado util
Tambien conviene exponer helpers:

- `PlayerColor(action)`
- `BaseResult(fishCard, action)`
- `ApplyArchetype(...)`
- `ApplyFatigue(...)`
- `ApplyRodModifier(...)`
- `ResolveTrackMovement(...)`

---

## Modelado de la baraja del pez

### Regla de negocio
La baraja del pez no es una bolsa de probabilidades abstracta. Debe comportarse como una pila real.

### Requisitos de implementacion
- crear el mazo base con `3 Embiste`, `3 Aguante`, `3 Quiebre`
- barajar con una fuente `rand.Rand` inyectable
- robar siempre desde el tope
- enviar cartas usadas a descarte
- cuando se intente robar una nueva carta y `draw` este vacio:
  - incrementar fatiga
  - barajar `discard`
  - retirar 3 cartas del tope del nuevo mazo
  - si el remanente es 0, disparar `fatiga total`

### API sugerida

```go
type FishDeck interface {
    Draw() (FishFamily, bool)
    ReshuffleForFatigue(rng *rand.Rand) FatigueResult
    Clone() FishDeck
    RemainingCounts() map[FishFamily]int
}
```

`Clone()` es importante porque el evaluador exacto necesitara ramificar estados sin mutar el original.

---

## Modelado de la decision del jugador

### Interfaz de politica

```go
type PlayerPolicy interface {
    ChooseAction(state CombatState) Action
    Name() string
}
```

### Politicas a implementar

#### 1. `RandomPolicy`
Elige una de las 3 acciones al azar. Sirve como baseline.

#### 2. `HeuristicPolicy`
Reglas simples, por ejemplo:

- evitar empates peligrosos contra arquetipos del pez
- priorizar la accion con mayor esperanza inmediata
- cambiar de sesgo segun track y fatiga

Sirve para comparar contra la politica optima.

#### 3. `OptimalPolicy`
Politica calculada por busqueda sobre el estado.

La accion elegida sera la que maximice:

`P(captura final | estado actual, accion)`

Si luego queremos optimizar EV de puntos o tiempo, solo cambiamos la funcion objetivo.

---

## Como calcular el perfil de estrategia optima

### Enfoque recomendado
Resolver el combate como un problema de valor esperado con recursion y memoizacion.

Para cada estado:

1. enumerar las acciones posibles del jugador
2. para cada accion, enumerar todas las posibles cartas superiores del pez compatibles con `draw`
3. ponderar por su probabilidad dentro del mazo restante
4. aplicar reglas del combate y obtener el siguiente estado
5. continuar recursivamente hasta `Capture` o `Escape`
6. guardar el mejor valor en cache

### Funcion de valor inicial

```text
V(estado terminal capturado) = 1.0
V(estado terminal escapado)  = 0.0
V(estado no terminal)        = max_accion sum_carta P(carta) * V(siguiente_estado)
```

### Ventajas
- produce una politica realmente optima para la especificacion cargada
- genera tablas reutilizables por estado
- evita depender de Monte Carlo para decidir cada accion

### Consideracion tecnica
Para que esto funcione bien, el estado debe ser serializable a una `cache key` estable. La representacion del mazo debe incluir el orden exacto de `draw` y `discard`, no solo sus conteos, si queremos modelar fielmente la pila real en ramificaciones deterministas.

### Optimizacion practica
Se puede arrancar con dos modos:

- `exact_ordered`: usa orden exacto de mazo, mas fiel, mas costoso
- `count_based`: usa solo conteos restantes, mas rapido, aproximado

Para la v1 recomendamos `exact_ordered` en el motor base y dejar `count_based` como optimizacion futura si hace falta.

---

## Resolucion de una ronda de combate
Orden de resolucion recomendado en codigo:

1. robar carta del pez
2. obtener resultado base de matriz
3. aplicar regla de arquetipo si el resultado base fue empate
4. aplicar fatiga acelerada a resultados no nulos
5. aplicar mejora ofensiva si corresponde
6. mover el track
7. evaluar captura o escape
8. si hubo intento de escape y `BaitGuardAvailable == true`, anular ese primer escape y fijar track en 5
9. mandar la carta del pez a descarte
10. si el combate sigue, comprobar la fatiga al intentar robar la siguiente carta

Nota: la proteccion de carnada debe consumirse solo cuando realmente evita un escape.

---

## Reglas congeladas para codificar
Estas decisiones ya quedan fijadas en el spec tecnico:

### 1. Momento exacto de la fatiga
La fatiga se evalua al intentar robar la siguiente carta de comportamiento y solo si el combate sigue.

### 2. Alcance de la mejora ofensiva
La mejora ofensiva aplica cuando:

- el jugador eligio la accion de ese color
- el resultado despues de arquetipo y fatiga queda en `0`

### 3. Multiple modifiers
En el reglamento actual solo existen mejoras ofensivas. En una ronda solo puede aplicar, como maximo, la mejora ofensiva del color de la accion elegida.

Estas reglas deben quedar escritas en el paquete como comentarios de especificacion y cubiertas por tests.

---

## Arquitectura recomendada del codigo

```text
internal/
  combat/
    types.go          // enums y modelos base
    matrix.go         // matriz base y helpers de color
    resolve.go        // resuelve una ronda de combate
    deck.go           // pila draw/discard y fatiga
    state.go          // estado, clonacion, cache key
    policy.go         // interfaz PlayerPolicy
    policy_random.go
    policy_heuristic.go
    policy_optimal.go // recursion + memoizacion
    simulator.go      // corre un combate completo con una politica
    stats.go          // metricas agregadas

internal/
  montecarlo/
    runner.go         // N simulaciones, semillas, agregacion
    scenarios.go      // escenarios predefinidos

cmd/
  combat-sim/
    main.go           // CLI para correr escenarios
```

Si se quiere mantener el repo muy pequeno al inicio, se puede arrancar con un solo paquete `internal/combat` y separar despues.

---

## Escenarios de entrada para Monte Carlo
Conviene definir un `Scenario` declarativo:

```go
type Scenario struct {
    InitialTrackPos      int
    Fish                 FishProfile
    BaitGuardAvailable   bool
    RodMods              RodModifiers
    Policy               string
    Simulations          int
    Seed                 int64
}
```

Esto permitira lanzar corridas como:

- pez negro desde slot 3 sin equipo
- pez rojo desde slot 5 con mejora ofensiva roja
- pez bicolor azul-rojo desde slot 4 con carnada activa
- comparativa de `RandomPolicy` vs `HeuristicPolicy` vs `OptimalPolicy`

---

## Salidas esperadas del runner
Cada corrida debe devolver al menos:

- `capture_rate`
- `escape_rate`
- `avg_rounds`
- `avg_fatigues`
- `capture_by_slot`
- `action_usage`
- `action_ev`
- `terminal_reason_breakdown`

Y opcionalmente:

- distribucion de longitud del combate
- distribucion de posicion terminal
- porcentaje de veces que la carnada salvo el combate
- porcentaje de veces que una mejora cambio un resultado

---

## Plan de implementacion por fases

### Fase 1 - Motor determinista del combate
Entregables:

- enums y modelos base
- matriz de combate
- aplicacion de arquetipo
- aplicacion de fatiga y fatiga total
- aplicacion de carnada y mejoras
- tests unitarios de reglas

Objetivo: poder correr un combate paso a paso con un mazo del pez ya fijado.

### Fase 2 - Baraja del pez como pila real
Entregables:

- constructor de mazo base
- shuffle reproducible por semilla
- draw/discard
- reshuffle por fatiga con trim de 3 cartas
- clonacion de estado
- tests de consistencia del mazo

Objetivo: que el pez ya tenga comportamiento completo y reproducible.

### Fase 3 - Politicas del jugador
Entregables:

- `RandomPolicy`
- `HeuristicPolicy`
- `OptimalPolicy` con memoizacion
- benchmark simple de costo por evaluacion

Objetivo: disponer de un perfil de decision reusable y comparable.

### Fase 4 - Simulador de combate completo
Entregables:

- `RunCombat(state, policy, rng)`
- eventos terminales normalizados
- recoleccion de estadisticas por combate

Objetivo: ejecutar un combate completo de inicio a fin.

### Fase 5 - Runner Monte Carlo
Entregables:

- `RunScenario(s Scenario)`
- agregacion de metricas
- export simple a JSON o tabla texto
- comparativa entre politicas

Objetivo: correr miles de simulaciones con una sola orden.

### Fase 6 - CLI y escenarios predefinidos
Entregables:

- binario `combat-sim`
- flags para slot, pez, seed, repeticiones, politica y loadout
- escenarios preconfigurados para balance

Objetivo: dejar la herramienta lista para iterar diseno.

---

## Estrategia de testing

### Unit tests obligatorios
- matriz base completa de 3x3
- regla de arquetipo para mono y bicolor
- fatiga de `+-1` a `+-2`
- cap de fatiga en `+-2`
- fatiga total cuando no queda mazo posible
- carnada evita exactamente el primer escape
- mejora ofensiva convierte `0` en `+1` cuando corresponde

### Tests de integracion
- combate completo con mazo fijo conocido
- combate completo con semilla fija
- simulaciones repetidas con la misma semilla devuelven el mismo agregado

### Tests de validacion de politica
- en estados triviales la politica optima coincide con la accion obvia
- la politica optima nunca rinde peor que la aleatoria en bateria de escenarios

---

## Riesgos tecnicos

### 1. Explosion del espacio de estados
Si `OptimalPolicy` usa orden exacto del mazo, la cache puede crecer. Mitigaciones:

- arrancar solo con combate aislado
- usar cache keys compactas
- medir antes de optimizar

### 2. Ambiguedades de reglas contaminando resultados
Si una sola regla queda ambigua, el Monte Carlo deja de ser confiable. Mitigacion:

- congelar reglas en tests y en este spec antes de correr balance serio

### 3. Mezclar evaluacion exacta y simulacion estocastica
Hay que evitar que `OptimalPolicy` use RNG. La politica optima debe ser pura y determinista sobre el estado recibido.

---

## Recomendacion final de implementacion
Orden recomendado de trabajo:

1. cerrar un spec tecnico minimo del combate con las ambiguedades resueltas
2. implementar el motor determinista y sus tests
3. implementar la baraja del pez como pila clonable
4. implementar `RandomPolicy` y `HeuristicPolicy`
5. implementar `OptimalPolicy` por recursion con memoizacion
6. montar el runner Monte Carlo
7. exponer una CLI para explorar balance

Este orden reduce riesgo porque primero garantiza que las reglas estan bien, luego modela al pez fielmente, y solo despues construye la capa de decision optima del jugador.

---

## Criterio de exito de la primera entrega
Consideraremos la primera entrega exitosa cuando podamos ejecutar algo equivalente a:

```bash
go run ./cmd/combat-sim --slot 4 --fish blue-red --policy optimal --simulations 100000 --seed 42
```

Y obtener un reporte con:

- tasa de captura
- tasa de escape
- rondas promedio
- fatigas promedio
- uso de acciones
- comparacion opcional contra una politica baseline
