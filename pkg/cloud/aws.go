package cloud

import (
	"context"
	"fmt"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/aws"
)

// AWSProvider implémente l'interface CloudProvider pour AWS
type AWSProvider struct {
	client *ec2.Client
	region string
}

// NewAWSProvider initialise la session AWS SDK v2
func NewAWSProvider(ctx context.Context, region string) (*AWSProvider, error) {
	// Charge la configuration AWS (lit ~/.aws/credentials, variables d'env AWS_ACCESS_KEY_ID, etc.)
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("impossible de charger la config AWS : %w", err)
	}

	client := ec2.NewFromConfig(cfg)

	return &AWSProvider{
		client: client,
		region: region,
	}, nil
}

// CreateInstance implémente la création d'une instance EC2
func (p *AWSProvider) CreateInstance(ctx context.Context, opts CreateInstanceOptions) (*Instance, error) {
	input := &ec2.RunInstancesInput{
		ImageId:      aws.String(opts.ImageID),
		InstanceType: ec2types.InstanceType(opts.InstanceType),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		TagSpecifications: []ec2types.TagSpecification{
			{
				ResourceType: ec2types.ResourceTypeInstance,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(opts.Name),
					},
				},
			},
		},
	}

	result, err := p.client.RunInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du lancement de l'instance EC2 : %w", err)
	}

	if len(result.Instances) == 0 {
		return nil, fmt.Errorf("aucun résultat d'instance retourné par AWS")
	}

	awsInst := result.Instances[0]

	return &Instance{
		ID:    aws.ToString(awsInst.InstanceId),
		State: string(awsInst.State.Name),
	}, nil
}

// GetInstanceStatus récupère les infos à jour d'une instance EC2
func (p *AWSProvider) GetInstanceStatus(ctx context.Context, instanceID string) (*Instance, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	result, err := p.client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du statut : %w", err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance non trouvée : %s", instanceID)
	}

	awsInst := result.Reservations[0].Instances[0]

	return &Instance{
		ID:        aws.ToString(awsInst.InstanceId),
		State:     string(awsInst.State.Name),
		PublicIP:  aws.ToString(awsInst.PublicIpAddress),
		PrivateIP: aws.ToString(awsInst.PrivateIpAddress),
	}, nil
}

// TerminateInstance résilie l'instance
func (p *AWSProvider) TerminateInstance(ctx context.Context, instanceID string) error {
	input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceID},
	}

	_, err := p.client.TerminateInstances(ctx, input)
	if err != nil {
		return fmt.Errorf("erreur lors de la résiliation de l'instance : %w", err)
	}

	return nil
}