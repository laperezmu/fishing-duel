# Pesca

Motor de juego en Go para un duelo de pesca por rondas. El proyecto esta separado en capas para que la logica del encuentro, la presentacion y las interfaces de usuario evolucionen por separado.

## Mapa del proyecto

- `cmd/fishing-duel/`: composicion del ejecutable CLI.
- `internal/domain/`: tipos base del juego.
- `internal/deck/`: mazo del pez, descarte y reciclado.
- `internal/encounter/`: configuracion y estado del track de distancia.
- `internal/match/`: estado compartido y resultado acumulado de la partida.
- `internal/playermoves/`: recursos y recarga de movimientos del jugador.
- `internal/rules/`: resolucion de rondas `Blue/Red/Yellow`.
- `internal/progression/`: efectos de una ronda sobre el estado del encuentro.
- `internal/endings/`: condiciones de fin de partida.
- `internal/game/`: motor orquestador del juego.
- `internal/presentation/`: traduccion de estado a textos y view models.
- `internal/app/`: flujo de sesion desacoplado de la UI concreta.
- `internal/cli/`: adaptador de terminal.

## Como extender el juego

### Cambiar reglas de combate

1. Crea otro evaluador en `internal/rules/` que implemente la interfaz usada por `internal/game`.
2. Conectalo en `cmd/fishing-duel/main.go` o en otro bootstrap.

### Cambiar como se mueve el pez en el track

1. Crea otra politica en `internal/progression/`.
2. Haz que implemente el contrato de progresion del motor.
3. Inyectala al crear `game.Engine`.

### Cambiar el criterio de fin

1. Crea otra condicion en `internal/endings/`.
2. Implementa el contrato de fin del motor.
3. Inyectala al construir la partida.

### Crear otra interfaz

1. Reutiliza `internal/app.Session`.
2. Implementa una UI propia con los tipos de `internal/presentation/`.
3. Usa un `presentation.Presenter` propio si quieres otro idioma, otro tono o incluso otra lectura visual del estado.

## Convencion recomendada

- `domain`, `deck`, `encounter`, `match`, `rules`, `progression`, `endings` y `game` no deberian depender de ninguna UI.
- `presentation` convierte estado tecnico a contenido mostrable.
- `app` coordina el flujo.
- `cmd/...` solo compone dependencias.

## Calidad de codigo

- Ejecuta `go test ./...` para validar la suite.
- Ejecuta `$(go env GOPATH)/bin/golangci-lint run` para revisar formato, estilo, errores comunes, complejidad y seguridad.
- La configuracion versionada vive en `.golangci.yml`.
