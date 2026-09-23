# ML-Elec

Système open source qui transforme des signaux techniques en décisions de maintenance fiables, prioritaires et actionnables avant la panne.

## Language

### Décision

**Décision**:
Une cause probable + une priorité métier + une action recommandée + un temps restant estimé + une gravité 0-100 + une confiance 0-1.
_Avoid_: alerte, notification, prédiction

**Gravité**:
Le max des scores par signal sur seuils bas/haut, avec paliers basse <40 / haute 40-75 / critique >75 et hystérésis 5 points.
_Avoid_: score ML, probabilité

**Signal**:
Une mesure technique brute issue d'un capteur (vibration, température, courant).
_Avoid_: donnée, métrique

**Alerte crédible**:
Un signal filtré dont la gravité et la menace ont été validées, avec un taux de faux positifs mesuré comme KPI principal. Faux positifs = rejetées / (retenues + rejetées) sur 30j glissants, cibles démo <20% et atelier <10%.
_Avoid_: alarme, notification brute

### Mémoire opérationnelle

**Mémoire opérationnelle**:
L'ensemble des cas validés par les techniciens et réutilisables pour décider (fiches de cas, règles métier, relations).
_Avoid_: base de données, historique, dataset ML

**Fiche de cas**:
Un épisode structuré : symptôme, contexte, cause, action, résultat, confiance, avec statut en attente puis retenue ou rejetée.
_Avoid_: ticket, rapport d'incident

**Règle métier**:
Un scénario connu encodé sous forme "si X alors Y + gravité + temps restant", lisible par un profil métier.
_Avoid_: règle ML, seuil, modèle
