package main

import (
    "context"
    "log"
    "net"
	"time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
    "google.golang.org/grpc"
    pb "rt0805/tp_app/operation_grpc"
    "rt0805/tp_app/operation"
    "strconv"
)

const (
    port = ":50051"
    mongoDBHost     = "localhost"
    mongoDBPort     = 27017
    mongoDBUsername = "root" // Remplacez par votre nom d'utilisateur MongoDB
    mongoDBPassword = "root"
)

type server struct {
    pb.UnimplementedDeviceServiceServer
}

func (s *server) SendData(ctx context.Context, in *pb.DeviceDataRequest) (*pb.DeviceDataResponse, error) {
    mongoURI := "mongodb://" + mongoDBUsername + ":" + mongoDBPassword + "@" + mongoDBHost + ":" + strconv.Itoa(mongoDBPort)
    client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatal(err)
    }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)

    collection := client.Database("operation").Collection("devices")

    var existingDevice operation.Device
    filter := bson.M{"device_name": in.Device.Name}
    err = collection.FindOne(ctx, filter).Decode(&existingDevice)

    totalOps := len(in.Device.Operations)
    failedOps := 0
    for _, op := range in.Device.Operations {
        if !op.HasSucceeded {
            failedOps++
        }
    }

    if err != nil { // Si l'appareil n'existe pas
        _, err = collection.InsertOne(ctx, bson.M{
            "device_name":       in.Device.Name,
            "total_operations":  totalOps,
            "failed_operations": failedOps,
            "operations":        in.Device.Operations,
        })
        if err != nil {
            log.Fatalf("Erreur lors de la création de l'appareil: %v", err)
        }
    } else {
        // Mettre à jour l'appareil existant
        _, err = collection.UpdateOne(
            ctx,
            filter,
            bson.M{
                "$inc": bson.M{
                    "total_operations":  totalOps,
                    "failed_operations": failedOps,
                },
                "$push": bson.M{"operations": bson.M{"$each": in.Device.Operations}},
            },
        )
        if err != nil {
            log.Fatalf("Erreur lors de la mise à jour de l'appareil: %v", err)
        }
    }

    return &pb.DeviceDataResponse{Success: true}, nil
}

func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("Échec de l'écoute: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterDeviceServiceServer(s, &server{})
    log.Printf("Serveur gRPC démarré sur le port %s", port)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Échec du serveur: %v", err)
    }
}
