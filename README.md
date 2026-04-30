# Pesca

Motor de juego en Go para un duelo de pesca por rondas. El proyecto ya no es solo una base tecnica: hoy incluye un loop CLI jugable con presets de pez y jugador, cartas con efectos por fase y una arquitectura separada por dominio, runtime, contenido y presentacion.

## Estado actual del juego

- El jugador elige un preset de barajas propias, una `rod` base, un preset de aditamentos y una situacion de agua antes de empezar el encuentro.
- Cada pesca se abre con una lectura breve del agua y un minijuego de cast por timing que define la distancia inicial del duelo.
- La `rod` separa limites de apertura del lance frente a limites de track que gobiernan el tablero y los escapes.
- El pez aparece ahora segun el agua base, la apertura efectiva del lance y los habitats habilitados por el loadout.
- Cada ronda el jugador sigue eligiendo entre tres acciones base: `Tirar`, `Recoger` y `Soltar`.
- Cada accion del jugador consume la carta superior visible de una mini-baraja por color.
- El pez roba su siguiente carta desde una baraja configurable, con orden fijo o barajado segun el perfil.
- El duelo muestra el historial visible de descarte del pez para ayudar a leer patrones, con espacio para cartas o perfiles que oculten parte de esa informacion.
- El motor resuelve efectos `on_draw` y efectos post-outcome antes de aplicar progresion y chequeos finales.
- El encuentro usa dos ejes de presion: `distancia` y `profundidad`, mas cierres por captura, escape y agotamiento de mazo.

## Quick start

Requisitos:

- Go instalado.

Ejecuta el juego desde la raiz del repo:

```bash
go run ./cmd/fishing-duel
```

Flujo actual del CLI:

1. Elige un preset del jugador.
2. Elige una `rod` base.
3. Elige un preset de aditamentos.
4. Elige una situacion de agua y resuelve el cast inicial.
5. El juego resuelve que pez aparece segun esa ventana de lanzamiento.
6. Juega el duelo ronda a ronda desde la terminal.

Tambien puedes validar el proyecto con:

```bash
go test ./...
$(go env GOPATH)/bin/golangci-lint run
```

## Loop jugable actual

1. Se elige una situacion de agua que aporta lectura visible y una base de apertura para el encuentro.
2. La `rod` del jugador aporta la base estructural del loadout.
3. Los aditamentos modifican limites de apertura, limites de track y preparan tags de habitat futuras.
4. El jugador resuelve un cast por timing que fija la distancia inicial del duelo dentro de esos limites efectivos.
5. El juego selecciona un perfil de pez compatible con el agua base, la ventana de distancia/profundidad y los habitats del loadout.
6. Se inicializa el estado del encuentro, el loadout del jugador, las barajas del jugador y la baraja del pez.
7. El jugador ve sus tres acciones disponibles junto a la carta superior visible de cada color cuando aplica.
8. El pez revela su carta activa para la ronda al robar del mazo, y el juego conserva un historial visible de lo ya descartado en el ciclo actual.
9. El motor aplica efectos de `draw`, resuelve el outcome base del combate y luego aplica efectos condicionados por victoria, derrota o empate.
10. La progresion modifica distancia y profundidad del pez, y puede disparar eventos derivados como `chapotea en la superficie`.
11. El encuentro termina por captura, escape horizontal, escape por profundidad, escape por chapoteo o resolucion al agotarse la baraja del pez.

## Presets actuales

### Presets del jugador

- `Clasico`: tres cartas lisas por color, sin efectos.
- `Apertura preparada`: ventajas tacticas al revelar la carta superior.
- `Respuesta vertical`: respuestas segun el outcome para mover profundidad.
- `Corriente mixta`: mezcla efectos de apertura con efectos post-outcome.

### Presets de `rod`

- `Control costero`: abre cerca y poco profundo, pero tolera algo mas de track que de apertura.
- `Versatil estandar`: referencia equilibrada para probar el loop base.
- `Presion de fondo`: prioriza profundidad de apertura y de track a cambio de menos margen horizontal.

### Presets de aditamentos

- `Sin aditamentos`: referencia limpia para comparar rods sin modificadores.
- `Kit de fondo`: empuja la pesca hacia profundidad y prepara habitats de fondo/canal.
- `Kit de lance largo`: gana margen horizontal a costa de perder algo de acceso vertical temprano.
- `Kit de estabilidad`: mejora el track defensivo en ambos ejes y prepara habitats de maleza/roca.

### Presets del pez

- El juego ya no pide elegirlos manualmente en el flujo principal; aparecen segun agua, apertura y loadout.
- `Clasico`: referencia base sin efectos, barajada y con reciclado que retira cartas.
- `Apertura con anzuelo`: perfil de `draw_tempo` con orden fijo.
- `Presion horizontal`: empuja el encuentro hacia mar abierto.
- `Presion vertical`: hunde o hace respirar al pez segun resultado.
- `Control de superficie`: gira en torno a superficie y eventos legibles.
- `Agotamiento de mazo`: concentra su plan en ventanas de cierre por agotamiento.
- `Corriente mixta`: mezcla varias familias de efectos para probar el pipeline completo.

Los perfiles del pez viven en `internal/content/fishprofiles/` y los del jugador en `internal/content/playerprofiles/`.

## Arquitectura actual

### Bootstrap y aplicacion

- `cmd/fishing-duel/`: composition root del ejecutable CLI.
- `internal/app/`: coordinacion de la sesion y flujo general desacoplado de la UI.
- `internal/cli/`: adaptador de terminal, seleccion de presets, render e input.
- `internal/presentation/`: view models y catalogos de texto para la UI.

### Dominio y runtime del combate

- `internal/domain/`: tipos base del juego.
- `internal/cards/`: `FishCard`, `PlayerCard` y `CardEffect` compartidos.
- `internal/game/`: motor del round, fases y orquestacion.
- `internal/rules/`: resolucion base del combate `Blue/Red/Yellow`.
- `internal/progression/`: impactos del round sobre distancia, profundidad y eventos derivados.
- `internal/encounter/`: estado espacial del pez, thresholds base y transiciones como `splash`.
- `internal/endings/`: condiciones terminales del encounter.
- `internal/match/`: estado compartido acumulado de la partida.

### Mazo del pez, jugador y contenido configurable

- `internal/deck/`: mazo del pez, descarte, reciclado y politicas de retiro de cartas.
- `internal/player/`: runtime del jugador, incluyendo loadout, `rod` y barajas por color.
- `internal/content/`: perfiles, presets y contenido configurable reusable para pez, jugador, loadout y contextos de agua.

En terminos de arquitectura, el proyecto hoy funciona como un monolito modular pequeno: `cmd/` compone dependencias, `app` coordina el flujo, `presentation/cli` son adaptadores de borde y el dominio del duelo vive separado del contenido configurable.

## Como extender el juego

### Cambiar reglas de combate

1. Crea otro evaluador en `internal/rules/`.
2. Haz que implemente la interfaz usada por `internal/game`.
3. Conectalo en `cmd/fishing-duel/main.go` o en otro bootstrap.

### Cambiar como progresa el encuentro

1. Crea otra politica en `internal/progression/`.
2. Haz que implemente el contrato de progresion del motor.
3. Inyectala al crear `game.Engine`.

### Cambiar el criterio de fin

1. Crea otra condicion en `internal/endings/`.
2. Implementa el contrato de fin del motor.
3. Inyectala al construir la partida.

### Anadir contenido del pez o del jugador

1. Extiende `internal/content/fishprofiles/` para nuevos perfiles o presets del pez.
2. Extiende `internal/content/playerprofiles/` para nuevas barajas del jugador.
3. Reutiliza `cards.CardEffect` y los triggers existentes para introducir nuevas cartas sin cambiar el loop principal.

### Crear otra interfaz

1. Reutiliza `internal/app.Session`.
2. Implementa una UI propia con los tipos de `internal/presentation/`.
3. Usa un `presentation.Presenter` propio si quieres otro idioma, otro tono o otra lectura del estado.

## Documentacion util

- `docs/backlog/001-roadmap-roguelike.md`: backlog y direccion actual del proyecto.
- `docs/features/`: planes de trabajo e iteraciones ya planteadas.
- `docs/discoveries/`: discovery docs para decisiones de sistema mas amplias.

## Convencion recomendada

- `cards`, `domain`, `deck`, `encounter`, `match`, `rules`, `progression`, `endings`, `game`, `content` y `player` no deberian depender de ninguna UI.
- `presentation` convierte estado tecnico a contenido mostrable.
- `app` coordina el flujo.
- `cmd/...` solo compone dependencias.
