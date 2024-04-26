package main

import (
	"testing"

	//pb "rt0805/tp_app/operation_grpc"
	"rt0805/tp_app/operation"
)

func TestConvertToDevicePB(t *testing.T) {
	// Créer un exemple de device
	device := operation.Device{
		Name: "testDevice",
		Operations: []operation.Operation{
			{Type: "CREATE", HasSucceeded: true},
			{Type: "UPDATE", HasSucceeded: false},
		},
	}

	// Convertir le device en pb.Device
	pbDevice := convertToDevicePB(&device)

	// Vérifier si la conversion est correcte
	if pbDevice.Name != device.Name {
		t.Errorf("Nom incorrect, attendu: %s, obtenu: %s", device.Name, pbDevice.Name)
	}

	if len(pbDevice.Operations) != len(device.Operations) {
		t.Errorf("Nombre d'opérations incorrect, attendu: %d, obtenu: %d", len(device.Operations), len(pbDevice.Operations))
	}

	for i, op := range device.Operations {
		if pbDevice.Operations[i].Type != op.Type {
			t.Errorf("Type d'opération incorrect pour l'opération %d, attendu: %s, obtenu: %s", i, op.Type, pbDevice.Operations[i].Type)
		}
		if pbDevice.Operations[i].HasSucceeded != op.HasSucceeded {
			t.Errorf("Statut de l'opération incorrect pour l'opération %d, attendu: %t, obtenu: %t", i, op.HasSucceeded, pbDevice.Operations[i].HasSucceeded)
		}
	}
}
