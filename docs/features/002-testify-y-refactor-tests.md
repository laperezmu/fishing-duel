# Plan de feature: testify-y-refactor-tests

## Objetivo
Integrar la libreria `testify` al proyecto y refactorizar los tests unitarios existentes para que sigan un estilo mas consistente basado en table tests y uso de mocks para dependencias inyectadas, mejorando legibilidad, mantenibilidad y facilidad de extension de la suite.

## Criterios de aceptacion
- El proyecto incorpora `github.com/stretchr/testify` en `go.mod` y la suite compila sin dependencias de prueba rotas.
- Los tests unitarios relevantes usan `require` y/o `assert` de `testify` en lugar de comparaciones manuales repetitivas cuando aporta claridad.
- Los tests que validan multiples escenarios equivalentes quedan orientados a table tests.
- Las dependencias inyectadas que hoy se resuelven con stubs ad hoc se revisan y, donde aporte valor, se reemplazan o encapsulan mediante mocks consistentes con el flujo del proyecto.
- La suite completa de tests pasa con `go test ./...`.
- La refactorizacion no cambia el comportamiento funcional cubierto por los tests existentes.

## Scope
### Incluye
- Agregar `testify` como dependencia de testing.
- Refactorizar los tests unitarios actuales para adoptar un estilo consistente de table tests donde tenga sentido.
- Introducir mocks o test doubles mas estructurados para dependencias inyectadas en tests de unidad.
- Limpiar duplicacion evidente y helpers de test que puedan centralizarse sin perder claridad.

### No incluye
- Cambiar la logica de negocio de produccion salvo ajustes minimos necesarios para mejorar testabilidad.
- Incorporar frameworks externos adicionales de mocking o generacion automatica de mocks.
- Reescribir tests de integracion o cambiar el enfoque global de validacion mas alla de esta suite unitaria.
- Aumentar cobertura con nuevas features funcionales no relacionadas con la refactorizacion.

## Propuesta de implementacion
- Incorporar `testify` usando principalmente `require` y `assert`, y evaluar `testify/mock` solo donde el mock aporte mas valor que un stub manual.
- Revisar paquete por paquete la suite actual para convertir tests repetitivos en table tests con casos nombrados.
- Identificar puntos de inyeccion de dependencias en pruebas como UI, presenter, evaluadores o collaborators del motor, y normalizar su representacion con mocks o dobles reutilizables.
- Mantener el refactor enfocado en ergonomia y estructura de prueba, evitando mezclarlo con cambios funcionales no relacionados.
- Actualizar documentacion puntual si algun paquete necesita dejar mas claro como probar contratos inyectables.

## Validacion
- Ejecutar `go test ./...` al finalizar la refactorizacion.
- Revisar que los tests sigan cubriendo los mismos escenarios que antes del cambio y que los nuevos tables mantengan nombres descriptivos por caso.
- Verificar que los mocks no introduzcan falsos positivos por expectativas demasiado debiles o demasiado acopladas a implementacion interna.

## Riesgos o decisiones abiertas
- No todos los tests necesitan mocks; en algunos casos un stub simple puede seguir siendo la opcion mas clara y estable.
- Una refactorizacion demasiado agresiva puede hacer que ciertos tests queden mas abstractos pero menos legibles si se fuerza table tests donde no aportan.
- Si algun paquete no expone puntos de inyeccion claros, podria requerir pequenos ajustes de diseño para mejorar testabilidad.
