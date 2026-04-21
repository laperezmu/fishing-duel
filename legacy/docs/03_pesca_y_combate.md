# Pesca y combate

## Track
El juego usa un track de **5 espacios**.

- Los espacios **3, 4 y 5** se usan en el lanzamiento.
- Los espacios **1 y 2** representan margen hacia la captura.

### Condiciones del track
- Si el pez se mueve **más allá del espacio 5**, escapa.
- Si el pez se mueve **por debajo del espacio 1**, es capturado.

---

## Calidad de captura
Cuando capturas un pez, obtienes puntos adicionales según el slot donde inició el combate:

- **slot 3** → +0 PV
- **slot 4** → +1 PV
- **slot 5** → +2 PV

La calidad se suma al valor del pez al final de la partida.

Si un pez revelado escapa durante el combate:
- no se conserva
- vuelve a la baraja de la zona de la que salió
- esa baraja se rebaraja inmediatamente

---

## Acciones del jugador
El jugador siempre tiene disponibles estas 3 acciones:

- **Forzar**
- **Tensar**
- **Soltar**

No se gastan.  
Cada ronda de combate simplemente elige una.

### Colores de las acciones
- **Forzar** = rojo
- **Tensar** = azul
- **Soltar** = amarillo

---

## Familias del pez
Las cartas de comportamiento del pez pertenecen a una de tres familias:

- **Embiste**
- **Aguante**
- **Quiebre**

Estas tres familias forman el núcleo del combate.

### Colores de las familias
- **Embiste** = rojo
- **Aguante** = azul
- **Quiebre** = amarillo

---

## Matriz base del combate

| Pez \ Jugador | Forzar | Tensar | Soltar |
|---|---:|---:|---:|
| **Embiste** | 0 | -1 | +1 |
| **Aguante** | +1 | 0 | -1 |
| **Quiebre** | -1 | +1 | 0 |

### Interpretación
- **+1** = el pez se mueve 1 espacio hacia **captura**
- **0** = empate, no se mueve
- **-1** = el pez se mueve 1 espacio hacia **escape**

---

## Lectura del ciclo
- **Forzar** vence a **Aguante**
- **Tensar** vence a **Quiebre**
- **Soltar** vence a **Embiste**
- misma familia = empate

---

## Secuencia de una ronda de combate
1. El controlador del pez toma la carta de comportamiento superior.
2. La mira y la coloca boca abajo.
3. El jugador activo elige en secreto una de sus 3 acciones.
4. Se revelan ambas a la vez.
5. Se determina el resultado del combate.
6. Se mueve el pez en el track.
7. La carta de comportamiento usada va al descarte.
8. Si no fue capturado ni escapó, empieza una nueva ronda.

---

## Valor del pez por tipo visual
La carta de pez no lleva PV impresos. Su valor se determina por su tipo:

- **pez negro** = 1 PV
- **pez monocolor** = 2 PV
- **pez bicolor** = 3 PV

---

## Arquetipos
El arquetipo ya no necesita texto en la carta del pez. Se determina por el color del pez.

### Colores
- **Rojo** = familia roja
- **Azul** = familia azul
- **Amarillo** = familia amarilla

### Pez normal
Un pez **negro** no tiene arquetipo.

### Pez monocolor
Un pez de un solo color tiene un arquetipo.

### Pez bicolor
Un pez bicolor tiene dos arquetipos.

---

## Regla de arquetipo
Si la carta de comportamiento revelada coincide en color con uno de los colores del pez **y** el jugador juega ese mismo color, el empate se convierte en **derrota del jugador**.

En otras palabras:
- en ese caso, un empate de color deja de ser `0`
- y pasa a ser `-1`

Esto aplica:
- a peces de un color
- y a peces bicolores en cualquiera de sus dos colores

---

## Mazo de comportamiento del pez
Cada combate usa una baraja base de **9 cartas**:
- 3 Embiste
- 3 Aguante
- 3 Quiebre

Cada combate empieza con una baraja nueva de esas 9 cartas, barajada por separado.

---

## Fatiga del pez
La fatiga no se comprueba al descartar la carta recién usada. Se comprueba cuando el combate continúa y se intenta robar la siguiente carta de comportamiento.

Si en ese momento el mazo del pez está agotado:

1. Se baraja su descarte.
2. Se retiran **3 cartas del tope** fuera de este combate.
3. El combate continúa con las cartas restantes.

Si el pez vuelve a agotar su mazo, se repite el proceso.

---

## Fatiga acelerada
Cuando un pez sufre su **primera fatiga**, todos los resultados de combate que normalmente moverían al pez **1 espacio** pasan a moverlo **2 espacios** en la misma dirección.

Esto afecta solo a resultados no nulos.

### Antes de fatiga
- victoria del jugador = `+1`
- empate = `0`
- victoria del pez = `-1`

### Después de la primera fatiga
- victoria del jugador = `+2`
- empate = `0`
- victoria del pez = `-2`

Las fatigas posteriores **no aumentan más** ese valor.  
El efecto queda **capado en ±2**.

---

## Fatiga total
Si al intentar formar un nuevo mazo por fatiga el pez ya no puede formar mazo:
- el pez sufre **fatiga total**
- y es **capturado automáticamente**

Esto aplica en cualquier slot.

---

## Carnadas

### Afinidad por zona
La carnada ya no depende del pez individual.

Cada zona normal tiene **2 afinidades de carnada** impresas en el mapa.

Si el jugador pesca en una zona y usa una carnada compatible con esa zona:
- la primera vez que el pez fuera a escaparse
- permanece en el espacio 5 en lugar de escapar

### Mar abierto
Mar abierto **no es compatible con ninguna carnada**.

Si pescas en mar abierto, la carnada no da beneficio.

---

## Señuelo
### Señuelo Explorador
Durante la preparación, después de colocar las 3 cartas de la zona en los espacios 5, 4 y 3:
- puedes mirar en privado 1 de esas cartas
- luego eliges normalmente tu lanzamiento

---

## Mejoras de caña

### Regla general
Las mejoras de caña modifican el resultado del combate en la casilla ofensiva de su color.

Las mejoras de caña son equipables de color:
- rojo
- azul
- amarillo

### Importante
Las mejoras ofensivas actúan como **sumatoria sobre el resultado del combate** solo cuando:

- el jugador eligió la acción de ese color
- el resultado tras matriz, arquetipo y fatiga quedó en `0`

### Orden de resolución
1. matriz base
2. arquetipo
3. fatiga
4. mejora de caña
5. mover el pez

Eso significa que una mejora ofensiva añade **+1 hacia captura** en su casilla correspondiente.

---

## Mejoras ofensivas
Estas mejoran el resultado del empate de su color.

### Roja — Carrete de Choque
Si el jugador eligió **Forzar** y el resultado quedó en `0`, añade +1 al resultado.

### Azul — Freno de Fondo
Si el jugador eligió **Tensar** y el resultado quedó en `0`, añade +1 al resultado.

### Amarilla — Bobina de Deriva
Si el jugador eligió **Soltar** y el resultado quedó en `0`, añade +1 al resultado.

---

## Restricción de mejoras de caña
No puedes tener en tu inventario otro equipable del mismo color.

Esto incluye:
- otra mejora de caña del mismo color
- o una mejora de bote del mismo color

---

## Ejemplos de resolución

### Ejemplo 1: empate sin fatiga + mejora ofensiva
- Resultado base: `0`
- Fatiga: no aplica
- Mejora ofensiva: `+1`

Resultado final: `+1`

### Ejemplo 2: empate con arquetipo
- Resultado base: `0`
- Arquetipo: `-1`
- Fatiga no activa
- Sin mejora

Resultado final: `-1`

### Ejemplo 3: empate con arquetipo y mejora ofensiva
- Resultado base: `0`
- Arquetipo: `-1`
- Fatiga no activa
- Mejora ofensiva: `0`

Resultado final: `0`
