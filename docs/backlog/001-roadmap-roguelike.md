# Backlog inicial roguelike

Este documento recopila los primeros pendientes de discovery y delivery para llevar el proyecto desde el duelo de pesca actual hasta un roguelike completo.

## Core Loop

### BL-001 Definir loop completo de una run
- **Tipo**: Discovery
- **Objetivo**: cerrar como inicia, progresa y termina una run.
- **Resultado esperado**: documento de flujo end-to-end de partida.
- **Criterios de cierre**:
  - queda definido el inicio, avance por zonas, derrota, victoria y retiro
  - queda claro que recompensas son de run y cuales son meta
- **Prioridad**: Alta

### BL-002 Definir estructura de mapa y tipos de nodo
- **Tipo**: Discovery
- **Objetivo**: decidir como avanza el jugador por zonas de pesca.
- **Resultado esperado**: taxonomia de nodos como pesca, evento, tienda, descanso, elite y boss.
- **Dependencias**: `BL-001`
- **Prioridad**: Alta

### BL-003 Diseñar progresion de dificultad entre zonas
- **Tipo**: Discovery
- **Objetivo**: decidir como escala el reto de una zona a la siguiente.
- **Resultado esperado**: reglas de dificultad por zona, rareza y encounter.
- **Dependencias**: `BL-001`, `BL-002`
- **Prioridad**: Alta

## Combat

### BL-004 Catalogar efectos que alteran el combate base
- **Tipo**: Discovery
- **Objetivo**: construir una taxonomia de efectos antes de seguir agregando reglas aisladas.
- **Resultado esperado**: lista de categorias como alteracion de movimientos, recursos, distancia, mazo, orden de resolucion o estado.
- **Prioridad**: Alta

### BL-005 Diseñar sistema extensible de efectos de combate
- **Tipo**: Delivery
- **Objetivo**: dejar una arquitectura preparada para peces y objetos con efectos composables.
- **Resultado esperado**: sistema tecnico para encadenar modificadores del encounter.
- **Dependencias**: `BL-004`
- **Prioridad**: Alta

## Fish

### BL-006 Definir arquetipos de peces
- **Tipo**: Discovery
- **Objetivo**: decidir familias de pez por identidad de gameplay.
- **Resultado esperado**: categorias como agresivo, evasivo, controlador, agotador, combo o boss.
- **Dependencias**: `BL-004`
- **Prioridad**: Alta

### BL-007 Diseñar perfiles data-driven de pez
- **Tipo**: Delivery
- **Objetivo**: describir peces por datos y no por wiring manual en codigo.
- **Resultado esperado**: estructura configurable para mazo, efectos, recompensas y condiciones.
- **Dependencias**: `BL-006`
- **Prioridad**: Alta

## Items y Build

### BL-008 Definir categorias de objetos del jugador
- **Tipo**: Discovery
- **Objetivo**: decidir que tipos de objetos construyen la build del jugador.
- **Resultado esperado**: clasificacion inicial de reliquias, consumibles, mejoras pasivas y modificadores de movimientos.
- **Prioridad**: Alta

### BL-009 Diseñar sistema de sinergias
- **Tipo**: Discovery
- **Objetivo**: definir como interactuan objetos, recursos y efectos de combate.
- **Resultado esperado**: reglas de stacking, activacion y limites de combinacion.
- **Dependencias**: `BL-008`, `BL-004`
- **Prioridad**: Alta

### BL-010 Implementar primer vertical slice de build
- **Tipo**: Delivery
- **Objetivo**: probar una run con un set pequeno de objetos y sinergias.
- **Resultado esperado**: prototipo jugable con recompensas y elecciones de build.
- **Dependencias**: `BL-008`, `BL-009`, `BL-005`
- **Prioridad**: Media

## Economy y Meta

### BL-011 Definir economia de run y meta-progresion
- **Tipo**: Discovery
- **Objetivo**: separar recursos temporales de recursos permanentes.
- **Resultado esperado**: lista de monedas, usos y sinks por capa de progresion.
- **Dependencias**: `BL-001`
- **Prioridad**: Alta

### BL-012 Diseñar sistema de recompensas entre encuentros
- **Tipo**: Discovery
- **Objetivo**: decidir como el jugador obtiene recursos, objetos o elecciones tras cada reto.
- **Resultado esperado**: tabla de recompensas por nodo y por tipo de pez.
- **Dependencias**: `BL-011`, `BL-002`
- **Prioridad**: Alta

## Collection

### BL-013 Definir bestiario y coleccion
- **Tipo**: Discovery
- **Objetivo**: concretar que significa completar la coleccion del juego.
- **Resultado esperado**: definicion de capturas, variantes, rarezas y criterios de registro.
- **Prioridad**: Media

### BL-014 Diseñar recompensas por coleccion
- **Tipo**: Discovery
- **Objetivo**: decidir si la coleccion solo es cosmetica o tambien desbloquea contenido.
- **Resultado esperado**: reglas de desbloqueo por hitos del bestiario.
- **Dependencias**: `BL-013`
- **Prioridad**: Media

## Tech y UX

### BL-015 Diseñar sistema de contenido data-driven
- **Tipo**: Discovery
- **Objetivo**: preparar peces, objetos y zonas para crecer sin depender de codigo duro.
- **Resultado esperado**: estrategia de datos y versionado de contenido.
- **Dependencias**: `BL-007`, `BL-008`
- **Prioridad**: Media

### BL-016 Diseñar guardado de run y progreso meta
- **Tipo**: Discovery
- **Objetivo**: decidir como persistir estado de run, coleccion y desbloqueos.
- **Resultado esperado**: contrato de persistencia y limites de versionado.
- **Dependencias**: `BL-001`, `BL-011`, `BL-013`
- **Prioridad**: Media

### BL-017 Mejorar UX de lectura de build y estado de run
- **Tipo**: Discovery
- **Objetivo**: asegurar que la complejidad futura siga siendo legible para el jugador.
- **Resultado esperado**: requerimientos de HUD, feedback de sinergias, log de efectos y resumen de build.
- **Dependencias**: `BL-008`, `BL-009`
- **Prioridad**: Media

## Orden sugerido

1. `BL-001`
2. `BL-002`
3. `BL-011`
4. `BL-004`
5. `BL-006`
6. `BL-008`
7. `BL-009`
8. `BL-013`
9. `BL-015`
