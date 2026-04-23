# Plan de feature: arquitectura-y-gobierno-de-paquetes

## Objetivo

Reducir la dispersion actual de paquetes bajo `internal/` y completar una refactorizacion inicial de arquitectura que deje separadas, de forma mas consistente, las capas de runtime, contenido configurable y bordes de aplicacion.

La meta de esta feature no es "ordenar carpetas" por estetica. Es dejar una propuesta arquitectonica concreta ya aplicada en las zonas mas criticas del proyecto, para que las siguientes features no sigan ampliando `internal/` con paquetes creados por necesidad local de cada iteracion.

## Criterios de aceptacion

- El plan define una propuesta explicita de organizacion de paquetes y no depende de un discovery separado para decidirla despues.
- Quedan definidas reglas practicas para crear, fusionar o extender paquetes.
- Se ejecuta una refactorizacion suficiente para separar con claridad contenido configurable y runtime al menos en las zonas de pez, jugador y presets.
- La nueva organizacion mejora la ubicacion de features futuras sin romper el comportamiento actual.
- `go test ./...` pasa.
- `golangci-lint run` pasa.

## Scope

### Incluye
- Revisar la estructura actual de `internal/` y clasificar sus paquetes por rol arquitectonico.
- Definir una propuesta objetivo de organizacion para runtime, contenido configurable, reglas de combate, presentacion, bootstrap y tooling.
- Definir reglas de gobierno para evitar seguir creando paquetes por feature puntual cuando corresponde extender una capacidad existente.
- Consolidar una primera reorganizacion transversal sobre las zonas de presets, perfiles y configuracion de pez y jugador.
- Reubicar paquetes concretos cuando hoy mezclen runtime y contenido configurable de forma innecesaria.
- Actualizar documentacion de arquitectura afectada.

### No incluye
- Reescribir todo el proyecto en una sola iteracion.
- Mezclar el refactor con cambios grandes de gameplay no relacionados.
- Congelar una arquitectura definitiva si el dominio aun necesita evolucionar.

## Propuesta de implementacion

### 1. Modelo de organizacion propuesto

La estructura futura debe tender a agruparse por capacidad estable y no por feature puntual.

Capas objetivo:

- `internal/domain/`
  - tipos puros, enums y conceptos base del juego sin orquestacion
- `internal/combat/`
  - logica del duelo y runtime del encounter
  - aqui deberian converger, progresivamente, piezas hoy repartidas entre `encounter`, `progression`, `endings`, `rules` y parte de `game`
- `internal/content/`
  - contenido configurable y perfiles derivados de datos
  - aqui encajan especialmente perfiles de pez, presets, cartas predefinidas y configuraciones reusables
  - esta feature debe dejar ya un primer slice real de esta idea y no solo mencionarla como direccion futura
- `internal/player/`
  - herramientas, barajas y recursos del jugador
  - aqui deberian converger, progresivamente, piezas hoy repartidas entre `playermoves` y `playerrig`
- `internal/presentation/`
  - transformacion de estado a vistas y catalogos de texto
- `internal/cli/`
  - render, input y flujo de consola
- `internal/app/`
  - composicion de la sesion y coordinacion alto nivel

Este es el objetivo final y debe de estar implementado al finalizar el plan

### 2. Reglas de gobierno de paquetes

- Crear un paquete nuevo solo si introduce una capacidad estable con vocabulario y ownership propios.
- No crear un paquete nuevo solo porque una feature nueva necesita 2 o 3 archivos relacionados.
- Si el codigo nuevo describe contenido configurable, debe vivir cerca de `content/` o de un paquete de perfiles/configuracion equivalente, no mezclado con runtime.
- Si el codigo nuevo resuelve comportamiento del duelo en tiempo de ejecucion, debe vivir cerca de `combat/` o de sus paquetes transicionales actuales.
- Si el codigo nuevo representa recursos o decisiones del jugador, debe vivir cerca de `player/` o de sus paquetes transicionales actuales.
- `cmd/` y bootstrap no deben acumular logica de configuracion de dominio que pueda vivir en paquetes reusables.
- Los presets son contenido de prueba o configuracion; no deben convertirse en el lugar donde se define el dominio de forma manual y repetida.

### 3. Estado actual ya alineado con el plan

Antes de extender el scope completo, ya se introdujo un primer movimiento valido:

- los arquetipos y perfiles del pez viven en `internal/content/fishprofiles/`
- los presets del pez ya dejaron de vivir en `internal/deck/`
- los presets del jugador viven en `internal/content/playerprofiles/` y ya no pertenecen a `internal/player/playermoves/`
- `internal/deck/` vuelve a quedar mas centrado en mecanica de mazo

Este avance no se considera el cierre completo del plan, sino el primer paso de la refactorizacion.

### 4. Refactorizacion a completar en esta feature

La refactorizacion completa de esta feature debe cubrir al menos estas zonas:

#### 4.1 Consolidar el contenido configurable del pez

- consolidar `internal/content/fishprofiles/` como origen de verdad para arquetipos y perfiles de pez
- dejar `internal/deck/` enfocado en mecanica de mazo: robar, reciclar, barajar y contar cartas

#### 4.2 Separar mejor la configuracion del jugador del runtime del jugador

- revisar `internal/player/playermoves/` para distinguir que parte es runtime de barajas y que parte es contenido configurable
- mover presets o configuraciones base del jugador fuera de los paquetes que ejecutan la mecanica de consumo y recuperacion si hoy estan mezclados
- dejar `internal/player/playermoves/` mas centrado en comportamiento runtime y validacion del sistema de decision

#### 4.3 Definir una primera ubicacion estable para contenido configurable

- consolidar el primer slice ya dentro de `internal/content/`, evitando una fase transicional adicional
- aplicar la misma regla a nuevos presets o perfiles para que no vuelvan a caer en paquetes de runtime o bootstrap

#### 4.4 Limpiar dependencias de borde

- reducir la logica de seleccion de presets y configuracion de contenido que hoy pueda estar demasiado pegada a `cmd/`
- mantener `cmd/` como composicion de alto nivel y no como catalogo de configuracion del juego

Resultado esperado de la refactorizacion completa:

- `internal/deck/` deja de mezclar runtime de mazo con catalogo de presets
- el contenido configurable del pez y del jugador queda mas cerca de perfiles y configuracion reusable
- futuras expansiones data-driven tienen una ubicacion mas obvia y menos acoplada a runtime
- el bootstrap consume configuracion, pero no la define de forma manual

### 5. Cambios concretos previstos

- Revisar paquetes actuales y clasificar cuales son de runtime, cuales son de contenido y cuales son de borde de aplicacion.
- Consolidar la reubicacion ya iniciada en la capa de pez.
- Reubicar la configuracion del jugador si hoy sigue demasiado mezclada con runtime.
- Ajustar imports, tests, bootstrap y CLI afectados por esa reorganizacion.
- Dejar documentadas reglas concretas sobre donde deben vivir presets, perfiles, runtime y composicion.
- Actualizar la documentacion con la convencion de organizacion resultante.

## Validacion

- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.
- Verificar que la refactorizacion reduce mezcla entre runtime y contenido configurable tanto para pez como para jugador.
- Verificar que una feature futura de contenido del pez o del jugador tiene una ubicacion mas clara que antes.
- Verificar que la propuesta evita seguir abriendo paquetes por feature puntual.

## Riesgos o decisiones abiertas

- Si intentamos completar toda la arquitectura objetivo en una sola iteracion, el refactor puede crecer demasiado.
- Extender el scope sobre jugador y bootstrap puede tocar varias rutas a la vez y producir churn relativamente amplio.
- Habra que balancear nombres arquitectonicos buenos con el costo real de mover codigo que hoy aun funciona.
