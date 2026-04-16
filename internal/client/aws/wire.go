package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/rs/zerolog/log"
)

func InitAWSDownscalers(downscalerConfig DownscalerConfig, profile string) (map[string][]DownscalerI, error) {
	d := map[string][]DownscalerI{}
	for region, services := range downscalerConfig.AwsCloud {

		if !isRegionSupported(region) {
			log.Warn().Msgf("Region '%s' is not aws region, skip client creation", region)
			continue
		}
		log.Info().Msgf("Generate aws client for %s region", region)
		cfg, err := LoadAWSConfig(context.Background(), region, profile)
		if err != nil {
			log.Fatal().Err(err)
		}
		if len(services.EC2) > 0 {
			log.Info().Msg("Detected EC2 services in config ...")
			ec2DownscalerCli, err := NewEc2Downscaler(*cfg, services.EC2)
			if err != nil {
				log.Fatal().Err(err)
			}
			d[region] = append(d[region], ec2DownscalerCli)
		}
		if len(services.Elasticache) > 0 {
			log.Info().Msg("Detected Elasticache services in config ...")
			elasticacheDownscalerCli, err := NewElasticacheDownscaler(*cfg, services.Elasticache)
			if err != nil {
				log.Fatal().Err(err)
			}
			d[region] = append(d[region], elasticacheDownscalerCli)
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
