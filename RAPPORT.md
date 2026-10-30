# Rapport de réalisation — Application de gestion des dépenses

## 1. Objectif
Le projet consiste à développer une application web en Go permettant d'enregistrer et de gérer des dépenses. Les données sont stockées dans une base SQLite afin de rester disponibles après l'arrêt du serveur.

## 2. Étape 1 — Serveur HTTP et notions Go
Le projet utilise un module Go déclaré dans `go.mod`. Le programme démarre dans `cmd/server/main.go`. Le paquet standard `net/http` fournit le serveur et le routeur `http.ServeMux`.
- `http.ListenAndServe` démarre le serveur sur le port 8080.
- `http.ResponseWriter` sert à écrire la réponse envoyée au navigateur.
- `*http.Request` contient la méthode, l'URL, les paramètres et les données du formulaire.
- `/ping` répond avec le statut 200 et le texte `pong`.
- Les routes sont déclarées avec les méthodes HTTP ; une méthode non autorisée reçoit une réponse 405 du routeur.

## 3. Étape 2 — Templates et organisation
Le code est séparé en couches : handlers, models, repository et fichiers web. Les pages sont générées avec `html/template`. Ce moteur échappe automatiquement les valeurs insérées dans le HTML, ce qui réduit le risque d'injection XSS.
Le formulaire envoie une requête POST. Après l'enregistrement, le serveur répond par une redirection HTTP 303 vers `/expenses` (Post/Redirect/Get), évitant une nouvelle soumission lors de l'actualisation.

## 4. Étape 3 — SQLite et repository
Le paquet `database/sql` fournit l'interface d'accès aux bases SQL et `modernc.org/sqlite` est le pilote SQLite. La fonction `Migrate` crée les tables `categories` et `expenses` et ajoute des catégories initiales.
Le repository regroupe les requêtes SQL : `Create`, `GetAll`, `GetByID`, `Update`, `Delete` et `Categories`. Les handlers ne contiennent donc pas directement les requêtes SQL.
Les valeurs des requêtes sont transmises avec des paramètres `?` : cela évite de concaténer les données utilisateur dans le SQL et protège contre les injections SQL.

## 5. Étape 4 — Catégories et filtres
Chaque dépense possède un `category_id`, clé étrangère vers la table `categories`. La liste utilise une jointure SQL `JOIN` pour afficher le nom de la catégorie en même temps que la dépense.
Les filtres sont transmis dans l'URL, par exemple `/expenses?category=2&from=2025-01-01&to=2025-01-31`. Le repository ajoute uniquement les conditions nécessaires et transmet leurs valeurs comme paramètres SQL. Sans filtre, toutes les dépenses sont affichées.

## 6. Cycle de traitement d'une requête
1. Le navigateur envoie une requête HTTP.
2. `ServeMux` choisit le handler selon le chemin et la méthode.
3. Le handler lit et valide les paramètres ou le formulaire.
4. Le handler appelle le repository.
5. Le repository exécute une requête SQLite paramétrée.
6. Le handler prépare les données et rend un template HTML, ou redirige après une modification.
7. `ResponseWriter` transmet le statut, les en-têtes et le contenu au navigateur.

## 7. Vérifications manuelles
Après lancement :
- Ouvrir `/ping` et vérifier le texte `pong`.
- Ouvrir `/expenses`, ajouter une dépense puis vérifier sa présence dans la liste.
- Modifier une dépense, puis la supprimer.
- Choisir une catégorie et/ou une période pour vérifier les filtres.
- Redémarrer le serveur et vérifier que les données sont toujours présentes dans `expenses.db`.
- Envoyer `POST /ping` avec curl : le serveur doit répondre 405.

Exemples de commandes :
```bash
curl -i http://localhost:8080/ping
curl -i -X POST http://localhost:8080/ping
```

## 8. Limite du dossier fourni
Les photos disponibles décrivent les étapes 1, 2, 3 et 4. L'énoncé de l'étape 5 n'a pas été transmis ; aucune exigence de cette étape n'est inventée dans cette version.
