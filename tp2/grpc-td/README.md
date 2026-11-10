docker compose up -d --build
lance les deux serveurs

docker compose run --rm go-client
client go - serveur py

docker compose run --rm python-client
client py - serveur go

docker compose logs -f
affiche les requetes
