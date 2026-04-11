package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func LoadAWSConfig(ctx context.Context, region, profile string) (*aws.Config, error) {
	opts := []func(*config.LoadOptions) error{}
	if region == "" {
		return nil, errors.New("missing region")
	}
	opts = append(opts, config.WithRegion(region))

	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
