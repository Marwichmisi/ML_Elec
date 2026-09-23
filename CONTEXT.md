# ML-Elec

Système open source qui transforme des signaux techniques en décisions de maintenance fiables, prioritaires et actionnables avant la panne.

## Language

### Décision

**Décision**:
Une cause probable + une priorité métier + une action recommandée + un temps restant estimé.
_Avoid_: alerte, notification, prédiction

**Signal**:
Une mesure technique brute issue d'un capteur (vibration, température, courant).
_Avoid_: donnée, métrique

**Alerte crédible**:
Un signal filtré dont la gravité et la menace ont été validées, avec un taux de faux positifs mesuré comme KPI principal.
_Avoid_: alarme, notification brute

### Mémoire opérationnelle

**Mémoire opérationnelle**:
L'ensemble des cas validés par les techniciens et réutilisables pour décider (fiches de cas, règles métier, relations).
_Avoid_: base de données, historique, dataset ML

**Fiche de cas**:
Un épisode validé structuré : symptôme, contexte, cause, action, résultat, confiance.
_Avoid_: ticket, rapport d'incident

**Règle métier**:
Un scénario connu encodé sous forme "si X alors Y", lisible par un profil métier.
_Avoid_: règle ML, seuil, modèle
