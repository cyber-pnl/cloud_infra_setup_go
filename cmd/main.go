package main

import (
	"fmt"
	"cloud-provisioner/pkg/cloud"
)

func main() {
	fmt.Println(" Initialisation du Cloud Provisioner...")

	// On prépare les options de test
	opts := cloud.CreateInstanceOptions{
		Name:         "web-server-test",
		InstanceType: "t2.micro",
		ImageID:      "ami-12345678",
	}

	fmt.Printf("Prêt à provisionner l'instance '%s' (%s)\n", opts.Name, opts.InstanceType)
}