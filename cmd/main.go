package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"cloud-provisioner/pkg/cloud"
	"cloud-provisioner/pkg/config"
)

func main() {
	// 1. Chargement de la configuration globale
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Erreur de configuration : %v\n", err)
	}

	// Définition d'un contexte global avec timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Initialisation du fournisseur AWS
	awsProvider, err := cloud.NewAWSProvider(ctx, cfg.AWSAccountRegion)
	if err != nil {
		log.Fatalf("❌ Impossible d'initialiser AWS : %v\n", err)
	}
	var provider cloud.CloudProvider = awsProvider

	// 2. Définition des sous-commandes via FlagSet
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createName := createCmd.String("name", "go-vm", "Nom de l'instance")
	createType := createCmd.String("type", cfg.DefaultInstance, "Type d'instance EC2 (ex: t2.micro)")
	createAMI := createCmd.String("ami", cfg.DefaultAMI, "ID de l'AMI (Amazon Machine Image)")

	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	statusID := statusCmd.String("id", "", "ID de l'instance EC2 (ex: i-0123456789abcdef0)")

	destroyCmd := flag.NewFlagSet("destroy", flag.ExitOnError)
	destroyID := destroyCmd.String("id", "", "ID de l'instance EC2 à résilier")

	// Verification de l'argument principal
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// 3. Routing des sous-commandes
	switch os.Args[1] {

	case "create":
		createCmd.Parse(os.Args[2:])
		if *createAMI == "" {
			fmt.Println("❌ Erreur : Vous devez fournir une AMI via '-ami' ou la variable DEFAULT_AMI_ID.")
			os.Exit(1)
		}

		opts := cloud.CreateInstanceOptions{
			Name:         *createName,
			InstanceType: *createType,
			ImageID:      *createAMI,
		}

		fmt.Printf("📦 Création de l'instance '%s' (%s) dans la région %s...\n", opts.Name, opts.InstanceType, cfg.AWSAccountRegion)
		inst, err := provider.CreateInstance(ctx, opts)
		if err != nil {
			log.Fatalf("❌ Échec de la création : %v\n", err)
		}

		fmt.Printf("✅ Instance créée avec succès ! ID: %s (État: %s)\n", inst.ID, inst.State)

		// Attente active que l'instance soit "running" et possède une IP
		readyInst, err := provider.WaitForInstanceRunning(ctx, inst.ID)
		if err != nil {
			log.Fatalf("❌ Erreur lors de l'attente du démarrage : %v\n", err)
		}

		fmt.Println("\n --- INSTANCE PRÊTE A L'EMPLOI ---")
		fmt.Printf(" ID         : %s\n", readyInst.ID)
		fmt.Printf(" État       : %s\n", readyInst.State)
		fmt.Printf(" IP Publique : %s\n", readyInst.PublicIP)
		fmt.Printf(" IP Privée   : %s\n", readyInst.PrivateIP)

	case "status":
		statusCmd.Parse(os.Args[2:])
		if *statusID == "" {
			fmt.Println(" Erreur : L'ID de l'instance est obligatoire (-id=i-xxx).")
			os.Exit(1)
		}

		inst, err := provider.GetInstanceStatus(ctx, *statusID)
		if err != nil {
			log.Fatalf(" Impossible de récupérer le statut : %v\n", err)
		}

		fmt.Println(" --- STATUT DE L'INSTANCE ---")
		fmt.Printf(" ID         : %s\n", inst.ID)
		fmt.Printf(" État       : %s\n", inst.State)
		fmt.Printf(" IP Publique : %s\n", inst.PublicIP)
		fmt.Printf(" IP Privée   : %s\n", inst.PrivateIP)

	case "destroy":
		destroyCmd.Parse(os.Args[2:])
		if *destroyID == "" {
			fmt.Println(" Erreur : L'ID de l'instance est obligatoire (-id=i-xxx).")
			os.Exit(1)
		}

		fmt.Printf(" Résiliation de l'instance %s...\n", *destroyID)
		err := provider.TerminateInstance(ctx, *destroyID)
		if err != nil {
			log.Fatalf(" Échec de la résiliation : %v\n", err)
		}

		fmt.Printf("  Demande de suppression envoyée pour l'instance %s.\n", *destroyID)

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(" Usage: cloud-provisioner <command> [options]")
	fmt.Println("\nCommandes disponibles :")
	fmt.Println("  create   Provisionner une nouvelle machine virtuelle")
	fmt.Println("           Options : -name, -type, -ami")
	fmt.Println("  status   Afficher le statut et les adresses IP d'une instance")
	fmt.Println("           Options : -id")
	fmt.Println("  destroy  Détruire une instance existante")
	fmt.Println("           Options : -id")
}