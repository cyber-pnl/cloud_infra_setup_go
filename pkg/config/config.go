package config

import (
	"fmt"
	"os"
)

// Config stocke la configuration globale pour notre infrastructure
type Config struct {
	AWSAccountRegion string
	DefaultInstance  string
	DefaultAMI       string
}

// LoadConfig charge la configuration depuis l'environnement ou applique des valeurs par défaut
func LoadConfig() (*Config, error) {
	region := getEnv("AWS_REGION", "eu-west-3") // eu-west-3 = Paris

	// t2.micro ou t3.micro sont dans l'offre gratuite (Free Tier)
	instanceType := getEnv("DEFAULT_INSTANCE_TYPE", "t2.micro")

	// AMI Amazon Linux 2023 par défaut (cette ID dépend de la région)
	amiID := getEnv("DEFAULT_AMI_ID", "")

	cfg := &Config{
		AWSAccountRegion: region,
		DefaultInstance:  instanceType,
		DefaultAMI:       amiID,
	}

	// Validation minimale
	if cfg.AWSAccountRegion == "" {
		return nil, fmt.Errorf("la région AWS ne peut pas être vide")
	}

	return cfg, nil
}

// getEnv est une fonction helper interne qui lit une variable d'environnement ou renvoie une valeur fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}