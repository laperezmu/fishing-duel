# Docs

La documentacion del proyecto se divide en dos capas activas y una capa legacy.

## Estructura

- `docs/features/`: planes de features trabajadas bajo el flujo de implementacion. Cada feature debe resolverse en un solo documento.
- `docs/backlog/`: backlog de investigacion, priorizacion y desarrollo futuro del juego.
- `docs/discoveries/`: documentos legacy de discovery creados antes de la convencion actual. No deben crearse nuevos discoveries separados para una feature.

## Regla practica

- Si una tarea ya se va a implementar, se crea como plan en `docs/features/`.
- Si una idea aun necesita discovery, definicion o priorizacion, vive en `docs/backlog/`.
- El analisis, la propuesta y la implementacion de una feature deben convivir en un unico documento de plan. No se abre un discovery separado para descubrir la propuesta despues.

## Referencias

- La plantilla de plan de feature sigue definida en `PROJECT_CONTEXT.md`.
- El backlog tiene su propia guia y template en `docs/backlog/`.
