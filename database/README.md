# Base de donnees MongoDB (dataset Steam)

Stack Docker qui demarre un MongoDB et y importe automatiquement `steam.json/steam.json`
(~60 Mo, un tableau JSON de jeux Steam) dans la collection `steam.games`.

## Contenu

| Fichier | Role |
| --- | --- |
| `docker-compose.yml` | MongoDB 7, le conteneur d'import, et mongo-express |
| `docker/seed.sh` | Attend Mongo, lance `mongoimport`, puis les index |
| `docker/indexes.js` | Index crees apres l'import |
| `.env.example` | Variables (identifiants, ports, base, collection) |

## Demarrage

```bash
cd database
cp .env.example .env     # deja fait si .env existe
docker compose up -d
```

Le service `mongo-seed` est un conteneur one-shot : il attend que MongoDB soit sain,
importe le JSON, cree les index, puis s'arrete (statut `Exited (0)`). C'est normal.

Suivre l'import :

```bash
docker compose logs -f mongo-seed
```

## Acces

- **URI de connexion** (depuis la machine hote, pour l'API) :
  `mongodb://root:rootpassword@localhost:27017/steam?authSource=admin`
- **Depuis un autre conteneur du meme reseau** : remplacer `localhost` par `mongodb`
- **mongo-express** : http://localhost:8081 (login `admin` / `admin`)
- **Shell mongo** :
  ```bash
  docker compose exec mongodb mongosh -u root -p rootpassword --authenticationDatabase admin steam
  ```

## Verifier les donnees

```bash
docker compose exec mongodb mongosh -u root -p rootpassword --authenticationDatabase admin steam \
  --quiet --eval 'print(db.games.countDocuments()); printjson(db.games.findOne())'
```

## Structure d'un document

```json
{
  "_id": 570,
  "title": "Dota 2",
  "required_age": 0,
  "is_free": true,
  "description": "...",
  "image": "https://...",
  "website": "http://www.dota2.com/",
  "developers": ["Valve"],
  "publishers": ["Valve"],
  "platforms": ["windows", "mac", "linux"],
  "metacritic_score": 90,
  "categories": [{ "id": 1, "description": "Multi-player" }],
  "genres": [{ "id": 1, "description": "Action" }],
  "recommendations": 14301,
  "release_date": "2013-07-09",
  "release_year": 2013
}
```

`_id` est l'appid Steam, donc l'import est idempotent sur la cle.

## Reimporter le dataset

L'import est ignore si la collection contient deja des documents. Pour forcer :

```bash
SEED_FORCE=1 docker compose up mongo-seed
```

## Repartir de zero (supprime les donnees)

```bash
docker compose down -v
docker compose up -d
```
