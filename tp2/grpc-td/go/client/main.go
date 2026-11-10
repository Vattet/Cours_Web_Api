// =====================================================================
// go/client/main.go  —  CLIENT gRPC en Go
// =====================================================================
// Rôle :
//   1. demander un texte à l'utilisateur
//   2. générer un nombre aléatoire
//   3. envoyer les deux au serveur via l'appel gRPC "SayHello"
//   4. afficher la réponse du serveur
//
// Lancement (depuis le dossier grpc-td/go/), APRÈS avoir démarré un serveur :
//     go run ./client
//
// Quel serveur répond ?
//   -> celui qui écoute sur localhost:50051, peu importe son langage :
//      le serveur Python (python/server.py) : cas 1 du TD
//      le serveur Go     (go/server/main.go) : cas "on mélange ?"
// =====================================================================

package main

import (
	"bufio"        // pour lire une ligne entière au clavier
	"context"      // pour donner un délai max (timeout) à l'appel
	"fmt"          // pour afficher du texte (Print, Printf)
	"log"          // pour afficher les erreurs et arrêter le programme
	"math/rand/v2" // pour le nombre aléatoire
	"os"           // os.Stdin = l'entrée clavier
	"strings"      // pour enlever le "\n" à la fin du texte saisi
	"time"         // pour exprimer "5 secondes"

	// Librairie gRPC + option "insecure" (connexion sans chiffrement TLS).
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Code GÉNÉRÉ à partir du .proto (dossier go/pb/).
	pb "example.com/grpc-td/pb"
)

// main() : 1re fonction exécutée par "go run ./client".
func main() {
	// --- Étape 1 : saisie au clavier -----------------------------------
	// On crée un "lecteur" sur le clavier (os.Stdin).
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Saisissez un texte : ")
	// ReadString('\n') lit tout ce qui est tapé jusqu'à la touche Entrée.
	texte, err := reader.ReadString('\n')
	if err != nil && len(texte) == 0 {
		log.Fatalf("lecture du texte impossible : %v", err)
	}
	// Le texte lu contient le "\n" final (la touche Entrée) : on l'enlève.
	texte = strings.TrimSpace(texte)

	// --- Étape 2 : nombre aléatoire entre 0 et 999 ---------------------
	nombre := rand.IntN(1000)
	fmt.Printf("Nombre aléatoire généré : %d\n", nombre)

	// --- Étape 3 : connexion au serveur --------------------------------
	// NewClient prépare la connexion vers localhost:50051.
	// (La vraie connexion réseau se fait au 1er appel, pas ici.)
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()), // pas de TLS
	)
	if err != nil {
		log.Fatalf("création du client gRPC impossible : %v", err)
	}
	// "defer" = "exécute ça à la FIN de main()". Ici : fermer la connexion.
	defer conn.Close()

	// Le "client" (stub) est généré : il a une méthode SayHello qui,
	// quand on l'appelle, envoie la requête au serveur sur le réseau.
	client := pb.NewGreeterClient(conn)

	// On crée un "contexte" avec un délai max de 5 secondes :
	// si le serveur ne répond pas à temps, l'appel échoue au lieu d'attendre à l'infini.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // libère les ressources du contexte à la fin de main()

	// --- Étape 4 : l'appel à distance (RPC) ----------------------------
	// Ce qui se passe ici :
	//   a) on construit un HelloRequest (Texte et Nombre, avec majuscule en Go)
	//   b) gRPC l'encode en binaire (protobuf) et l'envoie au serveur (HTTP/2)
	//   c) côté serveur, la fonction SayHello(...) est exécutée
	//   d) le serveur renvoie un HelloReply, que gRPC décode dans "response"
	// int32(nombre) : le .proto dit "int32", donc on convertit notre int.
	response, err := client.SayHello(ctx, &pb.HelloRequest{
		Texte:  texte,
		Nombre: int32(nombre),
	})
	if err != nil {
		// Ex : aucun serveur lancé -> "connection refused".
		log.Fatalf("appel RPC en échec : %v", err)
	}

	// --- Étape 5 : affichage de la réponse -----------------------------
	fmt.Printf("Réponse du serveur : %s\n", response.GetMessage())
}
