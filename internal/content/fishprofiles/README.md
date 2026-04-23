# internal/content/fishprofiles

Define perfiles y arquetipos de pez configurables para construir barajas sin wiring manual por preset.

## Responsabilidad

- Declarar arquetipos mecanicos de pez con ids estables.
- Describir perfiles concretos que derivan cartas, detalles y comportamiento base de baraja.
- Exponer presets de pez construidos desde esos perfiles para bootstrap y testing manual.
- Servir como primera capa data-friendly antes de mover contenido a formatos externos.
- Mantener separado el contenido del pez respecto al runtime puro del mazo.

## Limites

- No roba cartas ni administra el mazo activo.
- No resuelve rounds ni condiciones terminales.
- No reemplaza aun un sistema general de contenido externo.
