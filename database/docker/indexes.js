// Index crees apres l'import. Adaptes aux requetes typiques de l'API Steam.
const coll = db.getCollection(globalThis.TARGET_COLLECTION || "games");

const indexes = [
  [{ title: 1 }, { name: "idx_title" }],
  [{ release_year: -1 }, { name: "idx_release_year" }],
  [{ metacritic_score: -1 }, { name: "idx_metacritic_score" }],
  [{ is_free: 1 }, { name: "idx_is_free" }],
  [{ "genres.description": 1 }, { name: "idx_genres" }],
  [{ "categories.description": 1 }, { name: "idx_categories" }],
  [{ developers: 1 }, { name: "idx_developers" }],
  [{ publishers: 1 }, { name: "idx_publishers" }],
  [{ platforms: 1 }, { name: "idx_platforms" }],
  // Recherche plein texte sur titre + description
  [{ title: "text", description: "text" }, { name: "idx_text_search", weights: { title: 10, description: 1 } }],
];

for (const [keys, opts] of indexes) {
  coll.createIndex(keys, opts);
  print("  index cree : " + opts.name);
}

print("index total : " + coll.getIndexes().length);
