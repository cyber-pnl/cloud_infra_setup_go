package main

import (
	"fmt"
	"log"

	"cloud-provisioner/pkg/cloud"
	"cloud-provisioner/pkg/config"
)

func main() {
	fmt.Println(" Initialisation du Cloud Provisioner...")

	// 1. Chargement de la configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf(" Erreur de configuration : %v\n", err)
	}

	fmt.Printf(" Région configurée : %s\n", cfg.AWSAccountRegion)
	fmt.Printf(" Type d'instance par défaut : %s\n", cfg.DefaultInstance)

	// 2. Préparation des options d'instance avec la config
	opts := cloud.CreateInstanceOptions{
		Name:         "web-server-dev",
		InstanceType: cfg.DefaultInstance,
		ImageID:      cfg.DefaultAMI,
	}

	fmt.Printf(" Demande prête pour : %s (%s)\n", opts.Name, opts.InstanceType)
}