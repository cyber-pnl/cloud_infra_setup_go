package cloud

import "context"

// 1. Structure représentant une Instance VM de manière générique
type Instance struct {
	ID        string
	State     string // ex: "pending", "running", "terminated"
	PublicIP  string
	PrivateIP string
}

// 2. Options nécessaires pour créer une VM
type CreateInstanceOptions struct {
	Name         string
	InstanceType string // ex: "t2.micro"
	ImageID      string // ID de l'image OS (AMI sur AWS)
}

// 3. L'Interface (Le Contrat)
// Tout fournisseur cloud (AWS, GCP, Mock) DOIT implémenter ces méthodes.
type CloudProvider interface {
	// CreateInstance lance la création d'une nouvelle VM
	CreateInstance(ctx context.Context, opts CreateInstanceOptions) (*Instance, error)

	// GetInstanceStatus récupère l'état actuel d'une VM via son ID
	GetInstanceStatus(ctx context.Context, instanceID string) (*Instance, error)

	// TerminateInstance détruit une VM
	TerminateInstance(ctx context.Context, instanceID string) error
}