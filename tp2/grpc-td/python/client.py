# =====================================================================
# python/client.py  —  CLIENT gRPC en Python
# =====================================================================
# Rôle :
#   1. demander un texte à l'utilisateur
#   2. générer un nombre aléatoire
#   3. envoyer les deux au serveur via l'appel gRPC "SayHello"
#   4. afficher la réponse du serveur
#
# Lancement (depuis le dossier grpc-td/), APRÈS avoir démarré un serveur :
#     source venv/bin/activate
#     python python/client.py
#
# Quel serveur répond ?
#   -> celui qui écoute sur localhost:50051, peu importe son langage :
#      le serveur Go   (go/server/main.go)  : cas 2 du TD
#      le serveur Python (python/server.py) : cas "on mélange ?"
# =====================================================================

# Module standard Python pour générer des nombres aléatoires.
import random

# La librairie gRPC pour Python.
import grpc

# Fichiers GÉNÉRÉS à partir de proto/helloworld.proto :
# - helloworld_pb2       -> la classe HelloRequest (ce qu'on envoie)
# - helloworld_pb2_grpc  -> la classe GreeterStub (l'objet qui fait l'appel réseau)
import helloworld_pb2
import helloworld_pb2_grpc


# Appelée une seule fois, tout en bas du fichier, au lancement du script.
def main():
    # --- Étape 1 : saisie au clavier -----------------------------------
    # input() bloque jusqu'à ce que l'utilisateur appuie sur Entrée.
    texte = input("Saisissez un texte : ")

    # --- Étape 2 : nombre aléatoire entre 0 et 999 (inclus) -----------
    nombre = random.randint(0, 999)
    print(f"Nombre aléatoire généré : {nombre}")

    # --- Étape 3 : connexion au serveur --------------------------------
    # On ouvre un "channel" (canal de communication) vers le serveur.
    # "insecure" = pas de chiffrement (OK en local).
    # Le "with" ferme automatiquement le canal à la fin du bloc.
    with grpc.insecure_channel("localhost:50051") as channel:

        # Le "stub" est un objet généré qui a les mêmes méthodes que le service
        # déclaré dans le .proto. Appeler stub.SayHello(...) ressemble à un appel
        # de fonction normal, mais en réalité ça envoie une requête réseau.
        client = helloworld_pb2_grpc.GreeterStub(channel)

        # --- Étape 4 : l'appel à distance (RPC = Remote Procedure Call) ---
        # Ce qui se passe ici :
        #   a) on construit un HelloRequest avec nos 2 valeurs
        #   b) gRPC l'encode en binaire (protobuf) et l'envoie au serveur (HTTP/2)
        #   c) côté serveur, la méthode SayHello(...) est exécutée
        #   d) le serveur renvoie un HelloReply, que gRPC décode pour nous
        # Cette ligne BLOQUE jusqu'à ce que la réponse arrive
        # (ou jusqu'à 5 secondes max grâce à timeout=5, sinon -> erreur).
        response = client.SayHello(
            helloworld_pb2.HelloRequest(texte=texte, nombre=nombre),
            timeout=5,
        )

        # --- Étape 5 : affichage de la réponse ------------------------
        # response est un HelloReply -> il a un champ "message" (cf. le .proto).
        print(f"Réponse du serveur : {response.message}")


# S'exécute seulement si on lance ce fichier directement.
if __name__ == "__main__":
    main()
