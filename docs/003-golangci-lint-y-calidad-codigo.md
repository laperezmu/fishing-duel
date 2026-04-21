# Plan de feature: golangci-lint-y-calidad-codigo

## Objetivo
Integrar `golangci-lint` como herramienta central de analisis estatico para reforzar el cumplimiento de practicas alineadas con Effective Go, mejorar la legibilidad del codigo y establecer una base consistente de calidad para desarrollo local y validacion futura en CI.

## Criterios de aceptacion
- El proyecto incorpora `golangci-lint` con una configuracion versionada en el repositorio.
- La configuracion activa como fase principal el conjunto de linters acordado para estilo, legibilidad, errores comunes, complejidad y seguridad.
- El conjunto principal incluye al menos: `gofmt`, `goimports`, `gofumpt`, `govet`, `staticcheck`, `errcheck`, `ineffassign`, `unused`, `revive`, `misspell`, `gocritic`, `gosec`, `nestif` y `cyclop`.
- Existe una forma documentada de ejecutar los linters localmente.
- El codigo del proyecto queda ajustado para que el analisis pase con la configuracion introducida, o bien cualquier exclusion necesaria queda explicitamente justificada y minimizada.
- La configuracion prioriza reglas que mejoren lectura, idiomaticidad y mantenibilidad, evitando reglas arbitrarias o ruidosas sin justificacion.

## Scope
### Incluye
- Agregar la configuracion de `golangci-lint` al repositorio.
- Seleccionar y parametrizar el set principal de linters orientado a Effective Go y calidad de lectura.
- Ajustar el codigo existente para cumplir con la nueva linea base de lint.
- Documentar el uso local del comando y el objetivo de la configuracion.
- Definir umbrales o reglas de complejidad cuando algun linter lo requiera, buscando un equilibrio entre legibilidad y ruido.

### No incluye
- Integracion completa con GitHub Actions u otro pipeline de CI en esta iteracion, salvo que resulte imprescindible dejar preparado el comando.
- Incorporacion de herramientas externas fuera de `golangci-lint` como flujo principal de validacion.
- Refactors funcionales no relacionados con cumplir las reglas de lint.
- Reglas de estilo personalizadas que contradigan el estilo idiomatico de Go o agreguen friccion sin beneficio claro.

## Propuesta de implementacion
- Crear un archivo `.golangci.yml` o equivalente como fuente unica de verdad para el stack de lint.
- Activar dentro de `golangci-lint` el set principal completo acordado: `gofmt`, `goimports`, `gofumpt`, `govet`, `staticcheck`, `errcheck`, `ineffassign`, `unused`, `revive`, `misspell`, `gocritic`, `gosec`, `nestif` y `cyclop`.
- Revisar si conviene fijar limites iniciales para `cyclop` y `nestif` que mejoren lectura sin disparar falsos positivos excesivos.
- Corregir el codigo y tests afectados por las nuevas reglas, manteniendo el comportamiento actual del proyecto.
- Documentar en `README.md` o en la documentacion mas apropiada como ejecutar `golangci-lint` localmente y cual es la intencion de esta capa de calidad.
- Si algun linter requiere exclusion puntual, dejarla localizada, comentada y justificada para evitar ocultar problemas reales.

## Validacion
- Ejecutar `golangci-lint run` sobre todo el repositorio con la configuracion nueva.
- Ejecutar `go test ./...` despues de los cambios de cumplimiento para verificar que los ajustes de lint no alteraron comportamiento.
- Verificar manualmente que la configuracion final mejora consistencia de formato, nombres, manejo de errores y legibilidad general del codigo.
- Verificar que no queden comentarios `//nolint` ignorando linters como salida rapida para cerrar la iteracion.

## Riesgos o decisiones abiertas
- `gosec`, `gocritic`, `nestif` y `cyclop` pueden introducir ruido inicial; la calibracion de severidad y umbrales sera clave para no volver el flujo punitivo.
- Algunas reglas pueden sugerir refactors amplios; conviene distinguir entre correcciones necesarias para la linea base y mejoras futuras que merezcan una feature propia.
- Habra que decidir si la documentacion de uso local vive en `README.md`, `PROJECT_CONTEXT.md` o ambas, dependiendo de si se quiere orientar mas a contribuidores o al flujo operativo del proyecto.
