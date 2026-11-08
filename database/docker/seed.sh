#!/usr/bin/env bash
set -e

URI="mongodb://${MONGO_ROOT_USERNAME}:${MONGO_ROOT_PASSWORD}@${MONGO_HOST}:${MONGO_PORT}/${MONGO_DATABASE}?authSource=admin"
COLL="${MONGO_COLLECTION}"

# nombre de jeux déjà dans la base
count=$(mongosh "$URI" --quiet --eval "db.getCollection('$COLL').countDocuments()")

if [ "$count" -gt 0 ] && [ "$SEED_FORCE" != "1" ]; then
  echo "La base contient déjà $count jeux, pas d'import."
  exit 0
fi

echo "Import de steam.json..."
mongoimport --uri "$URI" --collection "$COLL" --file /seed/steam.json --jsonArray --drop

echo "Création des index..."
mongosh "$URI" --quiet --eval "var COLL='$COLL'" --file /seed/indexes.js

echo "Import terminé."
