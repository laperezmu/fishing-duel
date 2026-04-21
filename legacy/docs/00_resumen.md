# Resumen general del juego

## Premisa
El juego representa una competencia de pesca de **12 horas** para **2 a 5 jugadores**.

Cada jugador navega por un mapa con su bote, se abastece en un **almacén** común, decide si pesca o se reposiciona, y cuando un pez pica, resuelve un combate de tira y afloja usando información imperfecta, lectura del patrón del pez y preparación previa.

La partida dura **12 rondas**. Cada ronda representa una hora de competencia.

---

## Ejes principales del juego
El sistema se apoya en cuatro pilares:

1. **Movimiento y posicionamiento en el mapa**
2. **Gestión de inventario**
3. **Elección de zona y lanzamiento**
4. **Combate contra el pez**

---

## Objetivo
Al final de la partida, gana el jugador con más puntos.

Si varios jugadores empatan con la puntuación más alta, comparten la victoria.

Los puntos provienen de:

- peces capturados
- calidad de captura según la distancia del lanzamiento
- el objetivo **Coleccionista de crustáceos**
- el objetivo **Limpiador de la zona**
- cartas de **Tesoro hundido**

---

## Estructura general de ronda
La partida dura **12 rondas**.

En cada ronda, los jugadores realizan turnos de pesca siguiendo el flujo definido del turno individual.

El jugador inicial se determina al comienzo de la partida: inicia quien haya ido a pescar más recientemente. A partir de ahí, los turnos se resuelven siempre en orden fijo de izquierda a derecha desde ese jugador inicial.

---

## Estructura del mundo de juego
El juego utiliza:

- **5 zonas de pesca**
- **mar abierto**
- un **track de combate** de 5 espacios
- un **almacén común**
- un sistema de **cartas de comportamiento del pez**
- un sistema de **objetos consumibles y equipables**

---

## Conceptos clave

### Mar abierto
Toda casilla de agua que no pertenezca a una zona específica se considera **mar abierto**.

Si un jugador no se encuentra en una zona específica, igualmente puede pescar usando la baraja de mar abierto.

El mapa también incluye **5 muelles** iniciales sobre la costa.

### Calidad
Cuando capturas un pez, la distancia de lanzamiento determina calidad adicional:

- Slot 3 → +0 PV
- Slot 4 → +1 PV
- Slot 5 → +2 PV

Esta bonificación se registra durante la partida con fichas de calidad.

### Carnada
La carnada ya no depende de cada pez individual.  
La afinidad de carnada está determinada por la **zona de pesca**.

Cada zona normal tendrá **2 carnadas compatibles** impresas en el mapa.

Asignación propuesta:
- Estero / Manglar → Camarón + Cangrejo
- Arrecife somero → Cangrejo + Calamar
- Playa / Rompiente → Gusano marino + Sardina
- Punta / Cardumen pelágico → Sardina + Calamar
- Bajo profundo / Veril → Camarón + Gusano marino

**Mar abierto no es compatible con ninguna carnada.**

### Fatiga del pez
Cuando un combate se alarga lo suficiente, el pez sufre fatiga.  
La fatiga reduce el tamaño futuro de su mazo y además acelera el combate.

Tras la **primera fatiga**, todos los resultados no nulos del combate pasan de mover **1** a mover **2** espacios.

---

## Estado actual del diseño
El juego tiene completamente definidos:

- flujo del turno individual
- sistema de lanzamiento
- sistema de combate
- fatiga del pez
- catálogo del almacén
- tipos de eventos
- sistema de puntuación
- lenguaje visual del combate
- simplificación del sistema de peces
- simplificación de afinidades de carnada por zona
- mapa prototipo y muelles iniciales

La versión actual ya fija también el orden inicial de turno, la colocación inversa en muelles, el uso de fichas de calidad y la condición de victoria compartida en caso de empate final.
