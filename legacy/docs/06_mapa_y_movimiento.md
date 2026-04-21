# Mapa y movimiento

## Principio general
El juego usa un mapa prototipo de **20 x 20** casillas.

El mapa distingue cuatro tipos de casilla:
- **Costa**: no navegable
- **Muelle**: punto inicial de jugador
- **Zona de pesca**: navegable y pescable
- **Mar abierto**: navegable y pescable

---

## Movimiento

### Regla base
- cada bote tiene **2 puntos de movimiento** por turno, más los bonos permanentes de sus mejoras de bote
- el movimiento es **ortogonal**
- un punto de movimiento permite avanzar **1 casilla**

### Restricción de proximidad entre botes
Un bote no puede **terminar su movimiento**:
- en una casilla ocupada por otro bote
- ni en cualquiera de las **8 casillas adyacentes** a otro bote

Esta restricción se aplica a la **casilla final**, no al trayecto completo.

### Casillas navegables
Un bote puede terminar su movimiento en:
- mar abierto
- una zona de pesca
- un muelle libre

Un bote no puede entrar en:
- costa
- isla

---

## Pesca según posición
- si el bote termina su turno en una **zona de pesca**, pesca en esa zona
- si el bote termina su turno en una casilla de **mar abierto**, pesca en mar abierto
- si el bote permanece en un **muelle**, también pesca en mar abierto

---

## Muelles iniciales
El mapa tiene **5 muelles** fijos.

### Regla de selección de muelle
1. El jugador inicial es quien haya ido a pescar más recientemente.
2. A partir de ese jugador, el orden de turno queda fijado de izquierda a derecha para toda la partida.
3. En **orden inverso** a ese orden de turno, cada jugador elige **1 muelle libre**.
4. Coloca su bote en ese muelle.

Reglas adicionales:
- cada muelle solo puede ser elegido por **1 jugador**
- en partidas de menos de 5 jugadores, los muelles no elegidos quedan vacíos

Esta elección forma parte de la estrategia inicial del mapa.

---

## Mapa prototipo v1

### Leyenda
- `#` = costa no navegable
- `1..5` = muelles
- `~` = mar abierto
- `I` = isla no navegable
- `E` = Estero / Manglar
- `A` = Arrecife somero
- `R` = Playa / Rompiente
- `P` = Punta / Cardumen pelágico
- `V` = Bajo profundo / Veril

```text
    01 02 03 04 05 06 07 08 09 10 11 12 13 14 15 16 17 18 19 20
01  #  #  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~
02  #  #  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~
03  #  1  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~
04  #  #  ~  ~  E  E  E  E  ~  ~  ~  ~  ~  ~  ~  ~  A  A  A  ~
05  #  #  ~  ~  E  E  E  E  ~  ~  ~  ~  ~  ~  ~  ~  A  I  A  ~
06  #  #  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  A  A  A  ~
07  #  2  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~
08  #  #  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~
09  #  #  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~
10  #  #  ~  ~  ~  ~  ~  ~  R  R  R  R  R  ~  ~  ~  ~  ~  ~  ~
11  #  #  3  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~
12  #  #  #  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~
13  #  #  #  ~  ~  ~  ~  ~  P  P  P  P  ~  ~  ~  ~  ~  ~  ~  ~
14  #  #  #  #  ~  ~  ~  ~  P  P  P  P  ~  ~  ~  ~  V  V  V  ~
15  #  #  #  #  4  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  V  V  V  ~
16  #  #  #  #  #  #  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  V  V  V  ~
17  #  #  #  #  #  #  #  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~
18  #  #  #  #  #  #  #  #  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~
19  #  #  #  #  #  #  #  #  5  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~  ~
20  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #  #
```

---

## Distancias mínimas desde muelles a zonas
Las siguientes distancias miden pasos ortogonales mínimos desde cada muelle hasta la casilla más cercana de cada zona.

| Muelle | E | A | R | P | V |
| --- | ---: | ---: | ---: | ---: | ---: |
| `1` | `4` | `16` | `14` | `17` | `26` |
| `2` | `5` | `16` | `10` | `13` | `22` |
| `3` | `8` | `19` | `7` | `8` | `17` |
| `4` | `10` | `21` | `9` | `5` | `12` |
| `5` | `15` | `21` | `9` | `5` | `11` |
