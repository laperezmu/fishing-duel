# Plan de feature: arquetipos-de-peces

## Objetivo

Definir e implementar una primera taxonomia de arquetipos de peces que convierta la base mecanica actual del duelo en configuraciones reutilizables, faciles de ajustar y listas para crecer sin wiring manual por cada pez.

La meta de esta feature es aprovechar las herramientas ya disponibles del sistema - distancia, profundidad, efectos de carta por fase, `on_draw`, thresholds temporales, condiciones de cierre del encounter y presets CLI - para construir arquetipos mecanicos sencillos, poco tematicos y claramente configurables.

No buscamos peces con una personalidad demasiado cerrada ni una fantasia dificil de mantener. Buscamos patrones de comportamiento faciles de mover, renombrar, recombinar y convertir en datos, de forma que si la mecanica o la tematica cambian mas adelante no tengamos que rehacer la arquitectura ni los presets de prueba.

## Criterios de aceptacion

- Existe un documento de discovery que define una primera familia de arquetipos de pez con foco en comportamiento mecanico, no en fantasia rigida.
- Cada arquetipo describe su patron de juego, su relacion con distancia y profundidad, sus fases de efecto tipicas y los parametros minimos necesarios para configurarlo.
- La feature implementa una primera estructura data-driven o data-friendly para describir esos arquetipos y derivar barajas del pez desde configuracion.
- Los presets de inicio de partida del pez dejan de ser listas manuales ad hoc y pasan a apoyarse en los arquetipos configurados por esta feature.
- La taxonomia y la implementacion reutilizan explicitamente las herramientas ya existentes del combate actual en lugar de proponer sistemas paralelos.
- El resultado deja una base clara para ampliar despues perfiles de pez mas ricos sin reabrir el contrato central.

## Scope

### Incluye
- Definir una taxonomia inicial de arquetipos de peces para el duelo actual.
- Describir como se manifiesta cada arquetipo en terminos de cartas, fases, profundidad, distancia y condiciones de cierre.
- Identificar que herramientas del sistema actual ya sirven para cada arquetipo: `on_draw`, efectos post-outcome, thresholds temporales, splash, reciclado y composicion de baraja.
- Implementar una primera configuracion data-driven o data-friendly para arquetipos y perfiles basicos de pez.
- Actualizar los presets de inicio de partida del pez para que usen esos arquetipos configurados.
- Dejar una base simple y flexible para seguir ampliando perfiles de pez sin acoplarlos a nombres demasiado tematicos.

### No incluye
- Construir aun un sistema completo de contenido externo versionado o editable por usuario final.
- Redefinir el loop roguelike completo, mapa, economia o meta-progresion.
- Cerrar todavia la version final del formato data-driven para todo el juego.

## Propuesta de implementacion

- Crear un discovery en `docs/discoveries/011-arquetipos-de-peces.md` como salida principal de la feature, pero acompanar ese discovery con una primera implementacion real en codigo.
- Partir del sistema que hoy existe en el codigo y no de una fantasia abstracta: usar como base `FishCard`, efectos por fase, profundidad, captura dual, thresholds temporales y presets CLI.
- Redefinir los arquetipos con una estructura mecanica y configurable, por ejemplo:
  - id estable del arquetipo
  - objetivo mecanico principal
  - ejes de presion usados: horizontal, vertical, superficie, agotamiento, tempo, mixto
  - tipos de trigger prioritarios
  - repertorio de cartas o patrones de efecto permitidos
  - parametros minimos configurables
  - compatibilidad con el sistema actual
- Proponer una primera camada de arquetipos faciles de reconfigurar y de renombrar si la tematica cambia. Como punto de partida, conviene explorar al menos:
  - `horizontal_pressure`
  - `vertical_escape`
  - `surface_control`
  - `draw_tempo`
  - `deck_exhaustion`
  - `hybrid_pressure`
- Implementar una estructura de configuracion para perfiles de pez que permita derivar:
  - metadata del arquetipo
  - composicion de la baraja
  - politica de reciclado
  - orden fijo o barajado
  - etiquetas o detalles legibles para presets CLI
- Actualizar los presets de inicio de partida para que se construyan desde esos perfiles configurados, en vez de seguir cableando cartas manualmente dentro del preset.
- Mantener el resultado deliberadamente simple: mas importante que cubrir muchos casos es dejar una estructura facil de tocar, probar y cambiar.
- Cerrar el documento y el codigo con una base directa para ampliar luego perfiles mas ricos o mover la configuracion a formatos externos si hace falta.

## Validacion

- Revisar manualmente que cada arquetipo se apoye en mecanicas ya existentes o deje explicita su dependencia futura.
- Verificar que la taxonomia cubra estilos de presion realmente distintos y no simples variaciones cosmeticas.
- Validar que la implementacion permita construir presets de pez desde configuracion y no desde wiring manual repetido.
- Ejecutar `go test ./...`.
- Ejecutar `golangci-lint run`.

## Riesgos o decisiones abiertas

- Si definimos demasiados arquetipos demasiado pronto, la taxonomia puede perder foco y volverse decorativa.
- Si los arquetipos quedan demasiado pegados a presets de prueba, pueden limitar la expresividad del contenido futuro.
- Si la estructura de configuracion intenta resolver todo desde el inicio, puede volverse mas rigida de lo que conviene para una fase todavia cambiante del proyecto.
- Habra que decidir cuanto espacio dejamos para peces elite o boss sin sobredisenar el juego antes del loop roguelike completo.
