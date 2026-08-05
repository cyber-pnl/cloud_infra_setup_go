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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	fmt.Println(" Initialisation du Cloud Provisioner...")

	// 1. Chargement de la config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf(" Erreur config : %v\n", err)
	}

	// 2. Initialisation du fournisseur via l'interface
	var provider cloud.CloudProvider
	awsProvider, err := cloud.NewAWSProvider(ctx, cfg.AWSAccountRegion)
	if err != nil {
		log.Fatalf(" Impossible d'initialiser le client AWS : %v\n", err)
	}
	provider = awsProvider

	fmt.Printf(" Connexion à AWS établie dans la région : %s\n", cfg.AWSAccountRegion)

	// Note : Pour lancer une vraie VM, spécifie une AMI valide dans ton fichier/env
	// ex: AMI Ubuntu ou Amazon Linux valide dans ta région.
	if cfg.DefaultAMI == "" {
		fmt.Println("  Aucune AMI spécifiée (DEFAULT_AMI_ID). Étape de création ignorée pour l'instant.")
		return
	}

	opts := cloud.CreateInstanceOptions{
		Name:         "go-demo-vm",
		InstanceType: cfg.DefaultInstance,
		ImageID:      cfg.DefaultAMI,
	}

	fmt.Printf("Création de la VM '%s' (%s)...\n", opts.Name, opts.InstanceType)
	inst, err := provider.CreateInstance(ctx, opts)
	if err != nil {
		log.Fatalf(" Échec de la création : %v\n", err)
	}

	fmt.Printf(" VM créée avec succès ! ID: %s | État initial: %s\n", inst.ID, inst.State)
}