package main

import (
    "context"
    "testing"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
    //"go.mongodb.org/mongo-driver/mongo/writeconcern"
    pb "rt0805/tp_app/operation_grpc"
	//"google.golang.org/grpc"
)

func TestSendData(t *testing.T) {
    // Créer un faux serveur MongoDB en mémoire pour les tests
    ctx := context.Background()
    mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        t.Fatalf("Erreur de connexion à la base de données MongoDB: %v", err)
    }
    defer mongoClient.Disconnect(ctx)

    // Assurer que le serveur MongoDB est disponible
    if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
        t.Fatalf("Le serveur MongoDB n'est pas disponible: %v", err)
    }

    // Créer un serveur gRPC et une instance de la structure de serveur
    
    deviceServiceServer := &server{}

    // Envoyer des données au serveur
    request := &pb.DeviceDataRequest{
        Device: &pb.Device{
            Name: "testDevice",
            Operations: []*pb.Operation{
                {Type: "CREATE", HasSucceeded: true},
                {Type: "UPDATE", HasSucceeded: false},
            },
        },
    }
    response, err := deviceServiceServer.SendData(ctx, request)
    if err != nil {
        t.Fatalf("Erreur lors de l'envoi des données: %v", err)
    }

    // Vérifier la réponse du serveur
    if !response.Success {
        t.Error("La réponse du serveur indique un échec")
    }
}
