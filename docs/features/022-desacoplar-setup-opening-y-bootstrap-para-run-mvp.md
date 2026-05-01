# Plan de feature: desacoplar-setup-opening-y-bootstrap-para-run-mvp

## Objetivo

Ejecutar el slice minimo de arquitectura UI-agnostic necesario para que el primer MVP de run no nazca acoplado al CLI. La meta de esta feature no es completar toda la vision de `BL-023`, sino sacar del borde de terminal la orquestacion mas sensible de setup, opening y entrada al encounter, dejando a `cmd/` como composition root y a `internal/cli/` como adaptador de input/output.

Esta feature debe preparar el terreno para que una futura run pueda encadenar setup, opening, spawn y combate como casos de uso reutilizables desde mas de una interfaz, sin obligar todavia a implementar GUI ni a rehacer el presenter completo del juego.

## Criterios de aceptacion

- `cmd/fishing-duel/` deja de orquestar directamente el flujo completo de setup y opening contra `*cli.UI`.
- Existe un flujo de aplicacion reusable para setup y entrada al encounter consumible por interfaces UI-neutrales.
- La logica temporal o de estado del cast deja de vivir incrustada en `internal/cli/`.
- `presentation` incorpora los view models minimos necesarios para setup/opening sin depender de render directo de structs de dominio o contenido en la CLI.
- La CLI actual sigue siendo jugable tras el refactor.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope

### Incluye

- Extraer el wiring principal de setup/opening fuera de `cmd/fishing-duel/main.go` hacia una capa de aplicacion reusable.
- Introducir interfaces UI-neutrales para elecciones de setup, visualizacion de apertura y resolucion del cast.
- Mover el estado/control del cast a un componente reusable fuera de `internal/cli/`.
- Agregar view models minimos en `internal/presentation/` para setup/opening/cast cuando reduzcan coupling real.
- Adaptar la CLI para consumir esos contratos como adaptador.

### No incluye

- Implementar una GUI.
- Resolver todo el alcance de `BL-023`.
- Redisenar toda la UX del juego.
- Implementar todavia el runtime de run o el mapa de expedicion.

## Propuesta de implementacion

### 1. Crear un flujo de app para entrada al encounter

La preparacion previa al combate hoy esta demasiado repartida entre `cmd/`, `app` y `cli`. El primer paso debe ser encapsular el flujo completo de entrada al encounter en una capa de aplicacion reutilizable.

Direccion propuesta:

- crear un caso de uso o controlador de aplicacion que cubra setup inicial, opening, spawn y arranque de sesion
- mantener `cmd/` limitado a construir dependencias y delegar la ejecucion
- evitar que este flujo dependa de `*cli.UI` o de tipos concretos del adaptador terminal

### 2. Separar el contrato UI del setup y del opening

La CLI no deberia recibir presets y structs de dominio sin una frontera estable de app/presentation.

Direccion propuesta:

- definir interfaces pequenas para elecciones de setup y opening
- decidir que datos del dominio deben convertirse en opciones o views antes de llegar al adaptador
- mantener el contrato lo bastante semantico para que luego pueda reutilizarlo una GUI

### 3. Extraer el cast timing del paquete CLI

El timing del cast es hoy el punto mas dificil de reutilizar desde otra interfaz.

Direccion propuesta:

- mover la maquina de estado o controlador del cast a `internal/app/`, `internal/encounter/` o un paquete estable equivalente
- dejar al adaptador CLI solo con la responsabilidad de dibujar estado y capturar input
- preferir un modelo basado en estado/eventos/ticks en vez de sleeps o repintado incrustado en el borde

### 4. Introducir view models minimos de setup/opening

No hace falta resolver todo `presentation`, pero si estabilizar lo suficiente los datos consumidos por la CLI.

Direccion propuesta:

- agregar view models pequenos para opciones de setup, resumen de loadout, contexto de opening y estado del cast
- evitar que la CLI renderice directamente presets o structs de dominio cuando eso filtre demasiada estructura interna
- mantener el enfoque pragmatco: solo extraer los modelos realmente necesarios para el slice

### 5. Mantener la CLI funcional durante la migracion

La feature debe dejar el juego operativo y no romper el flujo actual.

Direccion propuesta:

- migrar por pasos pequenos
- conservar `session` y el combate actual como piezas reutilizables
- cubrir con tests el wiring nuevo y los contratos de app/presentation donde haya riesgo de regresion

## Archivos o zonas probables

- `cmd/fishing-duel/main.go`
- `internal/app/`
- `internal/app/opening.go`
- `internal/app/session.go`
- `internal/cli/`
- `internal/cli/opening.go`
- `internal/cli/ui.go`
- `internal/presentation/`

## Validacion

- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Verificar manualmente que el flujo completo actual `setup -> opening -> spawn -> combate` sigue siendo jugable desde la CLI.
- Verificar que `cmd/` ya solo compone dependencias y que la logica del cast no vive en `internal/cli/`.

## Riesgos o decisiones abiertas

- Si el contrato de setup/opening se vuelve demasiado abstracto demasiado pronto, puede introducir sobreingenieria antes de que la run concrete sus necesidades reales.
- Si el cast se modela con demasiada semantica de terminal, la futura GUI seguira heredando ese coupling aunque cambie de paquete.
- Habra que elegir cuanto de `presentation` pasa a modelos estructurados sin perder la legibilidad actual de la CLI.
- Conviene mantener este slice acotado; el resto de `BL-023` puede venir despues del primer vertical slice de run.
