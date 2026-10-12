#!/usr/bin/env bash
set -euo pipefail

HOST="${MONGO_HOST:-mongodb}"
PORT="${MONGO_PORT:-27017}"
USER="${MONGO_ROOT_USERNAME:-root}"
PASS="${MONGO_ROOT_PASSWORD:-rootpassword}"
DB="${MONGO_DATABASE:-steam}"
COLL="${MONGO_COLLECTION:-games}"
FILE="/seed/steam.json"
FORCE="${SEED_FORCE:-0}"

URI="mongodb://${USER}:${PASS}@${HOST}:${PORT}/${DB}?authSource=admin"

echo "[seed] cible : ${HOST}:${PORT} base=${DB} collection=${COLL}"

count="$(mongosh "$URI" --quiet --eval "db.getCollection('${COLL}').estimatedDocumentCount()")"
echo "[seed] documents deja presents : ${count}"

if [ "$count" -gt 0 ] && [ "$FORCE" != "1" ]; then
  echo "[seed] collection non vide, import ignore (SEED_FORCE=1 pour reimporter)."
  exit 0
fi

if [ "$FORCE" = "1" ]; then
  echo "[seed] SEED_FORCE=1 -> drop de la collection"
  mongosh "$URI" --quiet --eval "db.getCollection('${COLL}').drop()" || true
fi

echo "[seed] import de ${FILE} (~60 Mo, cela peut prendre une minute)..."
mongoimport \
  --uri "$URI" \
  --collection "$COLL" \
  --file "$FILE" \
  --jsonArray \
  --numInsertionWorkers 4 \
  --batchSize 1000 \
  --drop

echo "[seed] creation des index..."
mongosh "$URI" --quiet --eval "globalThis.TARGET_COLLECTION='${COLL}';" --file /seed/indexes.js

final="$(mongosh "$URI" --quiet --eval "db.getCollection('${COLL}').estimatedDocumentCount()")"
echo "[seed] termine : ${final} documents dans ${DB}.${COLL}"
