# Base de données MongoDB

On utilise MongoDB dans Docker avec le dataset Steam (`steam.json/steam.json`, environ 60 Mo).
Au premier lancement, le fichier est importé automatiquement dans la collection `steam.games`.

## Lancer la base

```bash
cd database
cp .env.example .env
docker compose up -d
```

Il y a 3 conteneurs :
- `mongodb` : la base
- `mongo-seed` : importe le JSON et crée les index, puis il s'arrête tout seul
- `mongo-express` : une interface web pour voir les données

Pour voir où en est l'import :

```bash
docker compose logs -f mongo-seed
```

## Se connecter

- URI pour l'API : `mongodb://root:rootpassword@localhost:27017/steam?authSource=admin`
- mongo-express : http://localhost:8081 (admin / admin)
- en ligne de commande :
  ```bash
  docker compose exec mongodb mongosh -u root -p rootpassword --authenticationDatabase admin steam
  ```

## Exemple de document

```json
{
  "_id": 570,
  "title": "Dota 2",
  "required_age": 0,
  "is_free": true,
  "description": "...",
  "developers": ["Valve"],
  "publishers": ["Valve"],
  "platforms": ["windows", "mac", "linux"],
  "metacritic_score": 90,
  "genres": [{ "id": 1, "description": "Action" }],
  "release_date": "2013-07-09",
  "release_year": 2013
}
```

Le `_id` est l'id du jeu sur Steam.

## Réimporter les données

Si la collection n'est pas vide, l'import ne se refait pas. Pour le forcer :

```bash
SEED_FORCE=1 docker compose up mongo-seed
```

Pour tout supprimer et recommencer :

```bash
docker compose down -v
docker compose up -d
```
