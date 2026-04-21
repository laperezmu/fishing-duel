# Reporte batch de simulacion de combate

- Escenarios: 24
- Simulaciones por escenario: 3000
- Seed base: 42

## Resumen global

| Metrica | Forzar | Tensar | Soltar |
| --- | ---: | ---: | ---: |
| Capture rate promedio forzado | 0.1493 | 0.1901 | 0.1726 |
| Escenarios donde la accion forzada fue la mejor | 9 | 7 | 8 |

| Politica | Capture rate promedio |
| --- | ---: |
| optimal | 0.6888 |
| heuristic | 0.6668 |
| random | 0.1840 |

## Escenarios baseline

| Escenario | Pez | Slot | Loadout | EV exacto F/T/S | Forced cap F/T/S | Mejor forzada | Optimal | Heuristic | Random |
| --- | --- | ---: | --- | --- | --- | --- | ---: | ---: | ---: |
| baseline-slot-3-black | black | 3 | none | 0.9664 / 0.9664 / 0.9664 | 0.6050 / 0.6003 / 0.6040 | forzar (0.6050) | 0.9647 | 0.9623 | 0.5093 |
| baseline-slot-3-blue | blue | 3 | none | 0.9367 / 0.8715 / 0.9550 | 0.6037 / 0.0100 / 0.5920 | forzar (0.6037) | 0.9603 | 0.9443 | 0.3310 |
| baseline-slot-3-blue-red | blue-red | 3 | none | 0.7230 / 0.7783 / 0.8911 | 0.0080 / 0.0130 / 0.6043 | soltar (0.6043) | 0.8920 | 0.8787 | 0.2127 |
| baseline-slot-3-blue-yellow | blue-yellow | 3 | none | 0.8911 / 0.7230 / 0.7783 | 0.6040 / 0.0120 / 0.0117 | forzar (0.6040) | 0.8983 | 0.8697 | 0.2167 |
| baseline-slot-3-red | red | 3 | none | 0.8715 / 0.9550 / 0.9367 | 0.0130 / 0.5967 / 0.6050 | soltar (0.6050) | 0.9650 | 0.9390 | 0.3360 |
| baseline-slot-3-red-yellow | red-yellow | 3 | none | 0.7783 / 0.8911 / 0.7230 | 0.0097 / 0.6017 / 0.0120 | tensar (0.6017) | 0.8883 | 0.8797 | 0.1967 |
| baseline-slot-3-yellow | yellow | 3 | none | 0.9550 / 0.9367 / 0.8715 | 0.5977 / 0.6010 / 0.0137 | tensar (0.6010) | 0.9533 | 0.9310 | 0.3360 |
| baseline-slot-4-black | black | 4 | none | 0.8188 / 0.8188 / 0.8188 | 0.2230 / 0.2247 / 0.2173 | tensar (0.2247) | 0.8127 | 0.8113 | 0.3530 |
| baseline-slot-4-blue | blue | 4 | none | 0.7466 / 0.6316 / 0.7730 | 0.2210 / 0.0000 / 0.2123 | forzar (0.2210) | 0.7670 | 0.7470 | 0.1973 |
| baseline-slot-4-blue-red | blue-red | 4 | none | 0.5089 / 0.5260 / 0.6201 | 0.0000 / 0.0000 / 0.2140 | soltar (0.2140) | 0.6133 | 0.5890 | 0.1127 |
| baseline-slot-4-blue-yellow | blue-yellow | 4 | none | 0.6201 / 0.5089 / 0.5260 | 0.2083 / 0.0000 / 0.0000 | forzar (0.2083) | 0.6333 | 0.5950 | 0.1033 |
| baseline-slot-4-red | red | 4 | none | 0.6316 / 0.7730 / 0.7466 | 0.0000 / 0.2110 / 0.2290 | soltar (0.2290) | 0.7670 | 0.7547 | 0.1907 |
| baseline-slot-4-red-yellow | red-yellow | 4 | none | 0.5260 / 0.6201 / 0.5089 | 0.0000 / 0.2020 / 0.0000 | tensar (0.2020) | 0.6227 | 0.6037 | 0.1070 |
| baseline-slot-4-yellow | yellow | 4 | none | 0.7730 / 0.7466 / 0.6316 | 0.2167 / 0.2123 / 0.0000 | forzar (0.2167) | 0.7697 | 0.7293 | 0.1927 |
| baseline-slot-5-black | black | 5 | none | 0.4870 / 0.4870 / 0.4870 | 0.0617 / 0.0673 / 0.0683 | soltar (0.0683) | 0.4943 | 0.4783 | 0.1843 |
| baseline-slot-5-blue | blue | 5 | none | 0.4282 / 0.2928 / 0.4177 | 0.0707 / 0.0000 / 0.0667 | forzar (0.0707) | 0.4153 | 0.4217 | 0.0870 |
| baseline-slot-5-blue-red | blue-red | 5 | none | 0.2761 / 0.2097 / 0.2842 | 0.0000 / 0.0000 / 0.0647 | soltar (0.0647) | 0.2773 | 0.2637 | 0.0377 |
| baseline-slot-5-blue-yellow | blue-yellow | 5 | none | 0.2842 / 0.2761 / 0.2097 | 0.0717 / 0.0000 / 0.0000 | forzar (0.0717) | 0.2900 | 0.2683 | 0.0417 |
| baseline-slot-5-red | red | 5 | none | 0.2928 / 0.4177 / 0.4282 | 0.0000 / 0.0677 / 0.0680 | soltar (0.0680) | 0.4260 | 0.4153 | 0.0933 |
| baseline-slot-5-red-yellow | red-yellow | 5 | none | 0.2097 / 0.2842 / 0.2761 | 0.0000 / 0.0700 / 0.0000 | tensar (0.0700) | 0.2733 | 0.2623 | 0.0480 |
| baseline-slot-5-yellow | yellow | 5 | none | 0.4177 / 0.4282 / 0.2928 | 0.0683 / 0.0660 / 0.0000 | forzar (0.0683) | 0.4417 | 0.4087 | 0.0777 |

## Escenarios loadout

| Escenario | Pez | Slot | Loadout | EV exacto F/T/S | Forced cap F/T/S | Mejor forzada | Optimal | Heuristic | Random |
| --- | --- | ---: | --- | --- | --- | --- | ---: | ---: | ---: |
| bait-slot-5-blue-red | blue-red | 5 | bait | 0.5198 / 0.5314 / 0.6298 | 0.0000 / 0.0000 / 0.2813 | soltar (0.2813) | 0.6317 | 0.5220 | 0.1140 |
| def-blue-slot-5-blue-red | blue-red | 5 | def=blue | 0.3327 / 0.9875 / 0.6526 | 0.0000 / 0.7903 / 0.0650 | tensar (0.7903) | 0.9883 | 0.9753 | 0.1430 |
| off-red-slot-4-red | red | 4 | off=red | 0.6316 / 0.7730 / 0.7466 | 0.0000 / 0.2173 / 0.2133 | tensar (0.2173) | 0.7857 | 0.7527 | 0.1950 |

