# Documentacion de diseno

Esta carpeta es la fuente de verdad del MVP actual.

- Usa `docs/90_backlog/master_backlog.md` para ver todo lo pendiente.
- Usa `docs/00_framework/workflow.md` para saber en que orden trabajar.
- Usa `docs/_templates/` para crear nuevas fichas sin perder consistencia.
- Considera `legacy/docs/` como referencia historica, no como version vigente.

## Estructura

- `00_framework/`: reglas de trabajo, criterio de cierre y proceso de iteracion.
- `01_foundations/`: vision, pilares, loop principal y reglas base de la run.
- `02_systems/`: sistemas jugables del MVP.
- `03_content/`: catalogos de zonas, peces, habilidades y objetos.
- `04_progression/`: progresion de run, meta y builds objetivo.
- `05_presentation/`: narrativa funcional, glosario, informacion visible y exclusiones.
- `90_backlog/`: pendientes, decisiones abiertas y parking lot.
- `_templates/`: plantillas para documentar temas nuevos.

## Regla de uso

1. Elige un tema del backlog con prioridad mas alta y dependencias resueltas.
2. Trabaja solo ese tema hasta dejarlo en `cerrado` o `en discusion`.
3. Actualiza su ficha y luego el backlog maestro.
4. Registra al final que quedo cerrado, que sigue abierto y cual es el siguiente paso.

## Estados permitidos

- `pendiente`
- `en_discusion`
- `cerrado`
- `aparcado`

## Convencion de archivos

- Un tema por archivo.
- Los nombres van con prefijo numerico para mantener el orden.
- Cada archivo debe incluir: objetivo, decision del jugador, reglas, preguntas abiertas y criterio de cierre.

## Referente principal

El referente estructural del MVP es `Balatro`:

- loop muy claro
- decisiones densas con informacion visible
- builds expresivas
- contenido escalable por capas
- expansion por catalogos, no por complejidad prematura
