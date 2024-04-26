package main

import (
	"context"
	"encoding/json"
	
	"io/ioutil"
	"log"
	"os"

	pb "rt0805/tp_app/operation_grpc"
    "rt0805/tp_app/operation"
	"google.golang.org/grpc"
    "path/filepath"
    "bufio"
    "fmt"
    "strings"
)

const (
    address = "localhost:50051"
    basePath = "../donnees/"
)

func convertToDevicePB(device *operation.Device) *pb.Device {
    operations := make([]*pb.Operation, len(device.Operations))
    for i, op := range device.Operations {
        operations[i] = &pb.Operation{
            Type:          op.Type,
            HasSucceeded:  op.HasSucceeded,
        }
    }
    return &pb.Device{
        Name:       device.Name,
        Operations:       operations,
    }
}
func main() {
    // Créer un lecteur pour lire à partir de la console standard
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Print("Entrez le nom du fichier JSON: ")
	fileName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Erreur lors de la lecture de l'entrée : %v", err)
	}
	// Nettoyer l'entrée du chemin
	fileName = strings.TrimSpace(fileName)

	// Construire le chemin complet du fichier
	fullPath := filepath.Join(basePath, fileName)
    // Ouvrir le fichier JSON
    file, err := os.Open(fullPath)
    if err != nil {
        log.Fatalf("Impossible d'ouvrir le fichier JSON: %v", err)
    }
    defer file.Close()

    // Lire le contenu du fichier JSON
    data, err := ioutil.ReadAll(file)
    if err != nil {
        log.Fatalf("Erreur lors de la lecture du fichier JSON: %v", err)
    }

    // Initialiser une connexion gRPC avec le serveur
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Impossible de se connecter: %v", err)
    }
    defer conn.Close()
    c := pb.NewDeviceServiceClient(conn)

    // Convertir les données JSON en objet Device
    var devices []operation.Device
    err = json.Unmarshal(data, &devices)
    if err != nil {
        log.Fatalf("Erreur lors de la conversion des données JSON: %v", err)
    }
    
    // Envoyer chaque appareil au serveur via gRPC
    for _, device := range devices {
        pbDevice := convertToDevicePB(&device)
        _, err = c.SendData(context.Background(), &pb.DeviceDataRequest{Device: pbDevice})
        if err != nil {
            log.Fatalf("Erreur lors de l'envoi des données: %v", err)
        }
        log.Printf("Données envoyées avec succès au serveur pour l'appareil %s.", device.Name)
    }
}
