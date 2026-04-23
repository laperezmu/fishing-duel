# Plan de feature: arquetipos-de-peces

## Objetivo

Definir una primera taxonomia de arquetipos de peces que convierta la base mecanica actual del duelo en identidades jugables reconocibles, expresivas y reutilizables.

La meta de esta feature es dejar claro que tipos de pez queremos construir con las herramientas ya disponibles del sistema: distancia, profundidad, efectos de carta por fase, `on_draw`, thresholds temporales, condiciones de cierre del encounter y presets CLI para validacion manual.

No buscamos una lista tematica superficial, sino un marco de diseño accionable que permita pasar despues a perfiles data-driven y contenido concreto sin volver a discutir desde cero como debe sentirse cada pez.

## Criterios de aceptacion

- Existe un documento de discovery que define una primera familia de arquetipos de pez con identidad de gameplay clara.
- Cada arquetipo describe su fantasia, su patron de juego, su relacion con distancia y profundidad, y los tipos de efectos de carta que mejor lo expresan.
- El documento reutiliza explicitamente las herramientas ya implementadas en el combate actual en lugar de proponer sistemas paralelos.
- Quedan identificados arquetipos inmediatos para prototipar con el pipeline actual y arquetipos que requieren pasos tecnicos posteriores.
- El documento deja una base directa para la siguiente feature de perfiles data-driven de pez.

## Scope

### Incluye
- Definir una taxonomia inicial de arquetipos de peces para el duelo actual.
- Describir como se manifiesta cada arquetipo en terminos de cartas, fases, profundidad, distancia y condiciones de cierre.
- Identificar que herramientas del sistema actual ya sirven para cada arquetipo: `on_draw`, efectos post-outcome, thresholds temporales, splash, reciclado y composicion de baraja.
- Proponer un conjunto pequeno de arquetipos priorizados para prototipos tempranos y testing manual mediante presets CLI.
- Dejar recomendaciones concretas para traducir cada arquetipo a perfiles de datos en la siguiente iteracion.

### No incluye
- Implementar aun nuevos peces, nuevas cartas o nuevos presets jugables.
- Redefinir el loop roguelike completo, mapa, economia o meta-progresion.
- Cerrar todavia el formato tecnico final de serializacion data-driven.

## Propuesta de implementacion

- Crear un discovery en `docs/discoveries/` como salida principal de la feature.
- Partir del sistema real que hoy existe en el codigo y no de una fantasia abstracta: usar como base `FishCard`, efectos por fase, profundidad, captura dual, thresholds temporales y presets CLI.
- Organizar los arquetipos en una estructura consistente, por ejemplo:
  - fantasia del pez
  - patron de presion sobre el jugador
  - comportamiento espacial principal
  - firmas de cartas y efectos representativos
  - riesgos de diseno y solapamientos con otros arquetipos
  - viabilidad con el sistema actual
- Proponer una primera camada de arquetipos que aproveche bien lo ya construido. Como punto de partida, conviene explorar al menos:
  - perseguidor horizontal
  - buceador o escapista vertical
  - controlador de superficie
  - pez de tempo con `on_draw`
  - pez agotador orientado al cierre por mazo
  - pez mixto o elite que combine varias presiones
- Cruzar cada arquetipo con ejemplos de presets CLI que nos permitan probarlos rapidamente antes de formalizarlos como contenido definitivo.
- Cerrar el documento con una traduccion hacia perfiles data-driven: que campos, tags o bloques de configuracion necesitara cada arquetipo.

## Validacion

- Revisar manualmente que cada arquetipo se apoye en mecanicas ya existentes o deje explicita su dependencia futura.
- Verificar que la taxonomia cubra estilos de presion realmente distintos y no simples variaciones cosmeticas.
- Confirmar que del documento se pueda derivar con claridad la siguiente feature de perfiles data-driven de pez.

## Riesgos o decisiones abiertas

- Si definimos demasiados arquetipos demasiado pronto, la taxonomia puede perder foco y volverse decorativa.
- Si los arquetipos quedan demasiado pegados a presets de prueba, pueden limitar la expresividad del contenido futuro.
- Habra que decidir cuanto espacio dejamos para peces elite o boss sin sobredisenar el juego antes del loop roguelike completo.
