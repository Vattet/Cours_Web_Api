// =====================================================================
// go/server/main.go  —  SERVEUR gRPC en Go
// =====================================================================
// Rôle : attendre des appels "SayHello" sur le port 50051 et y répondre.
//
// Lancement (depuis le dossier grpc-td/go/) :
//     go run ./server
//
// Qui l'appelle ?
//   -> le client Python (python/client.py)   : cas 2 du TD
//   -> le client Go     (go/client/main.go)  : cas "on mélange ?"
// =====================================================================

// En Go, un programme exécutable est toujours dans le package "main".
package main

import (
	"context" // type context.Context, reçu en paramètre de chaque appel gRPC
	"fmt"     // pour construire des chaînes formatées (Sprintf)
	"log"     // pour afficher des messages (avec la date/heure)
	"net"     // pour ouvrir un port TCP et écouter dessus

	// La librairie gRPC officielle pour Go.
	"google.golang.org/grpc"

	// Le code GÉNÉRÉ à partir de proto/helloworld.proto (dossier go/pb/).
	// "pb" est un alias : on écrira pb.HelloRequest, pb.HelloReply, etc.
	pb "example.com/grpc-td/pb"
)

// ---------------------------------------------------------------------
// La structure qui représente NOTRE serveur Greeter.
// ---------------------------------------------------------------------
// pb.UnimplementedGreeterServer est généré par protoc. On l'"embarque"
// (on l'inclut dans notre struct) pour 2 raisons :
//   - gRPC Go l'exige (c'est une sécurité pour la compatibilité future)
//   - si un jour on ajoute un rpc dans le .proto sans l'écrire ici,
//     le programme compile quand même (et renvoie "non implémenté").
type greeterServer struct {
	pb.UnimplementedGreeterServer
}

// ---------------------------------------------------------------------
// SayHello : le CODE exécuté quand un client appelle SayHello.
// ---------------------------------------------------------------------
// QUAND est-ce appelé ?
//   -> Jamais par nous directement !
//   -> C'est gRPC qui appelle cette fonction automatiquement à chaque requête
//      reçue (depuis server.Serve(...) plus bas).
//
// "(s *greeterServer)" veut dire que c'est une MÉTHODE de greeterServer.
// La signature (paramètres + retour) est imposée par le code généré :
//   ctx : infos sur l'appel (délai max, annulation...) — pas utilisé ici
//   req : la requête reçue, déjà décodée (req.GetTexte(), req.GetNombre())
// Retour : (la réponse, une erreur). nil = "pas d'erreur".
func (s *greeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	// Affichage dans le terminal du SERVEUR.
	// %q affiche le texte entre guillemets, %d affiche un entier.
	log.Printf("Requête reçue : texte=%q, nombre=%d", req.GetTexte(), req.GetNombre())

	// On construit la réponse. Sprintf = "fabrique une chaîne" (sans l'afficher).
	message := fmt.Sprintf(
		"J'ai reçu le texte « %s » et le nombre %d. (réponse Go)",
		req.GetTexte(),
		req.GetNombre(),
	)

	// On renvoie un pointeur (&) vers un HelloReply + nil (= pas d'erreur).
	// gRPC s'occupe ensuite de l'encoder et de l'envoyer au client.
	return &pb.HelloReply{Message: message}, nil
}

// main() est le point d'entrée : c'est la 1re fonction exécutée par "go run".
func main() {
	// 1) On ouvre le port TCP 50051 sur toutes les interfaces réseau.
	//    En Go, les fonctions qui peuvent échouer renvoient (résultat, erreur).
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		// Fatalf affiche le message puis ARRÊTE le programme.
		// (ex : le port est déjà utilisé par le serveur Python !)
		log.Fatalf("écoute impossible : %v", err)
	}

	// 2) On crée un serveur gRPC vide.
	server := grpc.NewServer()

	// 3) On "branche" notre greeterServer sur le serveur gRPC.
	//    = "quand une requête Greeter/SayHello arrive, appelle
	//       greeterServer.SayHello(...)".
	//    RegisterGreeterServer est une fonction GÉNÉRÉE (go/pb/helloworld_grpc.pb.go).
	pb.RegisterGreeterServer(server, &greeterServer{})

	log.Println("Serveur Go démarré sur le port 50051")

	// 4) On démarre le serveur sur le port ouvert.
	//    Serve() BLOQUE : il tourne en boucle pour toujours en attendant
	//    des requêtes. On n'en sort que s'il y a une erreur (ou Ctrl+C).
	if err := server.Serve(listener); err != nil {
		log.Fatalf("serveur arrêté : %v", err)
	}
}
