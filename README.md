# Gestion des dépenses — Go

Application pédagogique correspondant aux étapes 1 à 4 du document fourni.

## Prérequis
- Go 1.22 ou plus récent
- Connexion Internet au premier lancement pour télécharger `modernc.org/sqlite`

## Lancer le projet
Depuis le dossier du projet :

```bash
go mod tidy
go run ./cmd/server
```

Ouvrir ensuite http://localhost:8080

La base `expenses.db` est créée automatiquement dans le dossier courant. Pour choisir un autre chemin :
- Windows PowerShell : `$env:DB_PATH="data/expenses.db"`
- Linux/macOS : `DB_PATH=data/expenses.db go run ./cmd/server`

## Routes
- `GET /` : redirection vers la liste
- `GET /expenses` : liste et filtres (`category`, `from`, `to`)
- `GET /expenses/new` et `POST /expenses` : formulaire d'ajout
- `GET /expenses/{id}/edit` et `POST /expenses/{id}/update` : modification
- `POST /expenses/{id}/delete` : suppression
- `GET /about` : présentation
- `GET /ping` : réponse texte `pong`

## Organisation
- `cmd/server` : point d'entrée
- `internal/models` : structures de données
- `internal/repository` : SQLite, migration et requêtes SQL
- `internal/handlers` : routes et traitement HTTP
- `web/templates` : pages HTML
- `web/static` : feuille de style

## Remarque
L'étape 5 n'apparaît pas dans les photos transmises. Le projet livré couvre les étapes 1 à 4 visibles ; il faudra adapter ou compléter l'étape 5 lorsque son énoncé sera disponible.
