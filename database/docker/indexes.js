// Index sur les champs utilisés par les filtres de l'API
const games = db.getCollection(COLL);

games.createIndex({ title: 1 });
games.createIndex({ release_year: -1 });
games.createIndex({ metacritic_score: -1 });
games.createIndex({ is_free: 1 });
games.createIndex({ "genres.description": 1 });
games.createIndex({ "categories.description": 1 });
games.createIndex({ developers: 1 });
games.createIndex({ publishers: 1 });
games.createIndex({ platforms: 1 });

// pour la recherche texte sur le titre et la description
games.createIndex({ title: "text", description: "text" });

print("Index créés : " + games.getIndexes().length);
