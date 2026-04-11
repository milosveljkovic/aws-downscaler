package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/rs/zerolog/log"
)

func InitAWSDownscalers(downscalerConfig DownscalerConfig, profile string) (map[string][]DownscalerI, error) {
	d := map[string][]DownscalerI{}
	loadOpts := []func(*config.LoadOptions) error{}
	for region, services := range downscalerConfig.AwsCloud {

		if isRegionSupported(region) {
			log.Warn().Msgf("Region '%s' is not aws region, skip client creation", region)
			continue
		}
		log.Info().Msgf("Generate aws client for %s region", region)
		if profile != "" {
			loadOpts = append(loadOpts, config.WithSharedConfigProfile(profile))
		}
		cfg, err := LoadAWSConfig(context.Background(), region, "dev")
		if err != nil {
			log.Fatal().Err(err)
		}
		if len(services.EC2) > 0 {
			log.Info().Msg("Detected EC2 services in config ...")
			cli, err := NewEc2Downscaler(*cfg, services.EC2)
			if err != nil {
				log.Fatal().Err(err)
			}
			d[region] = append(d[region], cli)
		}
	}
	return d, nil
}

func isRegionSupported(region string) bool {
	if region == "" {
		log.Error().Msg("Region can not be empty")
		return false
	}

	resolver := ec2.NewDefaultEndpointResolverV2()
	_, err := resolver.ResolveEndpoint(context.TODO(), ec2.EndpointParameters{Region: &region})

	if err != nil {
		log.Error().Msgf("Unsupported AWS region: %s", region)
		return false
	}

	return true
}
