# =====================================================================
# python/server.py  —  SERVEUR gRPC en Python
# =====================================================================
# Rôle : attendre des appels "SayHello" sur le port 50051 et y répondre.
#
# Lancement (depuis le dossier grpc-td/) :
#     source venv/bin/activate
#     python python/server.py
#
# Qui l'appelle ?
#   -> le client Go     (go/client/main.go)   : cas 1 du TD
#   -> le client Python (python/client.py)    : cas "on mélange ?"
# Le serveur ne sait pas (et s'en fiche) dans quel langage est le client :
# il reçoit juste un message HelloRequest au format protobuf.
# =====================================================================

# "futures" sert à créer un groupe de threads (ThreadPoolExecutor) :
# chaque requête reçue sera traitée dans un thread de ce groupe.
from concurrent import futures

# La librairie gRPC pour Python (installée avec : pip install grpcio).
import grpc

# Ces 2 fichiers ne sont PAS écrits à la main : ils sont GÉNÉRÉS par
#   python -m grpc_tools.protoc ... proto/helloworld.proto
# - helloworld_pb2       : contient les classes HelloRequest et HelloReply
# - helloworld_pb2_grpc  : contient GreeterServicer (base du serveur),
#                          GreeterStub (pour le client) et
#                          add_GreeterServicer_to_server (pour brancher le service)
import helloworld_pb2
import helloworld_pb2_grpc


# ---------------------------------------------------------------------
# La classe qui CONTIENT LE CODE du service "Greeter".
# ---------------------------------------------------------------------
# On hérite de GreeterServicer (classe générée à partir du .proto).
# On doit y écrire une méthode pour CHAQUE "rpc" déclaré dans le .proto
# (ici une seule : SayHello).
class Greeter(helloworld_pb2_grpc.GreeterServicer):

    # QUAND est-ce appelé ?
    #   -> Jamais par nous directement !
    #   -> C'est la librairie gRPC qui appelle cette méthode AUTOMATIQUEMENT
    #      à chaque fois qu'un client fait "client.SayHello(...)".
    #
    # Paramètres (fournis par gRPC) :
    #   request : un objet HelloRequest déjà décodé -> request.texte, request.nombre
    #   context : infos sur l'appel (métadonnées, délai, possibilité d'annuler...)
    #             on ne s'en sert pas ici.
    #
    # Ce qu'on retourne : un objet HelloReply. gRPC se charge de l'encoder
    # et de le renvoyer au client sur le réseau.
    def SayHello(self, request, context):
        # Affichage dans le terminal du SERVEUR (pour voir qu'on a bien reçu).
        # !r affiche le texte entre guillemets (pratique pour voir les espaces).
        print(f"Requête reçue : texte={request.texte!r}, nombre={request.nombre}")

        # On construit la réponse avec les infos reçues.
        reponse = helloworld_pb2.HelloReply(
            message=(
                f"J'ai reçu le texte « {request.texte} » "
                f"et le nombre {request.nombre}. (réponse Python)"
            )
        )

        # Ce "return" renvoie la réponse au client (via gRPC).
        return reponse


# ---------------------------------------------------------------------
# Fonction qui CRÉE et DÉMARRE le serveur.
# ---------------------------------------------------------------------
# Appelée une seule fois, tout en bas du fichier, au lancement du script.
def serve():
    # 1) Création du serveur gRPC.
    #    max_workers=10 : jusqu'à 10 requêtes traitées en même temps.
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    # 2) On "branche" notre classe Greeter sur le serveur.
    #    Concrètement, on dit à gRPC : "quand une requête Greeter/SayHello
    #    arrive, appelle Greeter().SayHello(request, context)".
    helloworld_pb2_grpc.add_GreeterServicer_to_server(Greeter(), server)

    # 3) On choisit l'adresse et le port d'écoute.
    #    "[::]" = toutes les interfaces réseau (IPv4 et IPv6).
    #    "insecure" = pas de chiffrement TLS (OK pour des tests en local).
    server.add_insecure_port("[::]:50051")

    # 4) Démarrage : à partir d'ici, le serveur accepte les connexions.
    #    (start() rend la main tout de suite, le serveur tourne en arrière-plan.)
    server.start()
    print("Serveur Python démarré sur le port 50051")

    # 5) On bloque le programme ici pour qu'il ne se termine pas.
    #    Sans cette ligne, le script finirait et le serveur s'arrêterait.
    server.wait_for_termination()


# Ce bloc s'exécute UNIQUEMENT quand on lance le fichier directement
# (python python/server.py), pas quand on l'importe depuis un autre fichier.
if __name__ == "__main__":
    try:
        serve()
    except KeyboardInterrupt:
        # Quand on fait Ctrl+C dans le terminal, Python lève KeyboardInterrupt.
        # On l'attrape pour afficher un message propre au lieu d'une erreur.
        print("\nArrêt du serveur.")
