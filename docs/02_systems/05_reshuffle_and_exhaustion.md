# Reshuffle y agotamiento

## Estado

- Estado: `pendiente`
- Prioridad: `alta`
- Depende de: `04_fish_decks_and_affinities`

## Objetivo

Definir que pasa cuando el mazo del pez se agota y como esto genera un arco de combate.

## Base actual

- Al agotarse el mazo, se rehace desde el descarte.
- Despues se pierden cartas del tope.
- Si ya no quedan recursos, el pez entra en agotamiento.
- La captura es automatica si el pez esta suficientemente cerca.

## Criterio de cierre

- El proceso completo cabe en una secuencia corta.
- Se entiende por que el agotamiento existe y como recompensa buena lectura o control.
