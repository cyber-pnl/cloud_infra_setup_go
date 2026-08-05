package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud-provisioner/pkg/cloud"
	"cloud-provisioner/pkg/config"
)

func main() {
	// On définit un timeout global de 3 minutes pour toutes les opérations Cloud
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	fmt.Println(" Initialisation du Cloud Provisioner...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf(" Erreur config : %v\n", err)
	}

	var provider cloud.CloudProvider
	awsProvider, err := cloud.NewAWSProvider(ctx, cfg.AWSAccountRegion)
	if err != nil {
		log.Fatalf(" Connexion AWS impossible : %v\n", err)
	}
	provider = awsProvider

	if cfg.DefaultAMI == "" {
		fmt.Println("Aucune AMI fournie. Passe DEFAULT_AMI_ID pour créer une vraie VM.")
		return
	}

	opts := cloud.CreateInstanceOptions{
		Name:         "go-automation-vm",
		InstanceType: cfg.DefaultInstance,
		ImageID:      cfg.DefaultAMI,
	}

	// 1. Demande de création de la VM
	fmt.Printf(" Envoi de la demande de création pour '%s'...\n", opts.Name)
	inst, err := provider.CreateInstance(ctx, opts)
	if err != nil {
		log.Fatalf(" Erreur création : %v\n", err)
	}

	fmt.Printf(" VM demandée ! ID: %s | État: %s\n", inst.ID, inst.State)

	// 2. Attente active que l'IP soit attribuée et le statut à 'running'
	readyInstance, err := provider.WaitForInstanceRunning(ctx, inst.ID)
	if err != nil {
		log.Fatalf(" Erreur attente : %v\n", err)
	}

	fmt.Println("\n --- PROVISIONING RÉUSSI ---")
	fmt.Printf(" Instance ID : %s\n", readyInstance.ID)
	fmt.Printf(" État       : %s\n", readyInstance.State)
	fmt.Printf(" IP Publique : %s\n", readyInstance.PublicIP)
	fmt.Printf(" IP Privée   : %s\n", readyInstance.PrivateIP)
}