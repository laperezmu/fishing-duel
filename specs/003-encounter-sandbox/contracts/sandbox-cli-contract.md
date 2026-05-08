# Contract: Encounter Sandbox CLI

## Purpose

Definir el contrato funcional del sandbox CLI para configurar, ejecutar e inspeccionar encounters sin afectar el flujo de `fishing-run`.

## Supported Interaction Modes

1. Guided mode
   - Recorre un flujo amigable similar al actual.
   - Permite elegir presets base y avanzar con friccion baja.
2. Manual mode
   - Expone seleccion explicita de preset de jugador, cana, aditamentos, pez, contexto, cartas concretas y seed.
   - Permite aplicar overrides controlados del encounter.
3. Scenario mode
   - Ejecuta una configuracion reutilizable de QA o regresion.
4. Non-interactive or semi-reproducible mode
   - Permite ejecutar configuraciones compartibles sin recorrer todos los prompts manuales.

## Setup Contract

- El sandbox debe distinguir claramente entre:
  - valores heredados del preset base
  - valores seleccionados manualmente
  - valores fijados por escenario
  - valores sobrescritos por override
- El preset de pez puede seleccionarse manualmente en modo sandbox.
- Las cartas concretas del jugador y del pez pueden fijarse manualmente cuando la prueba lo requiera.
- La seed reproducible debe poder asociarse a la configuracion completa del sandbox.

## Output Contract

- El modo por defecto muestra el estado necesario para jugar o probar manualmente.
- Los modos mas detallados muestran:
  - triggers activados
  - efectos resueltos
  - prioridades
  - desempates
  - evidencia estructurada antes y despues de la resolucion
- La salida visible no define reglas; solo representa el resultado del runtime.

## Compatibility Expectations

- `fishing-run` conserva su flujo y su framing jugable.
- El sandbox puede reutilizar piezas de encounter, presentation y contenido, pero no debe imponer prompts ni conceptos nuevos sobre `fishing-run`.
- Los escenarios o exports del sandbox pueden ampliar la experiencia CLI sin alterar la fuente de verdad del engine.
