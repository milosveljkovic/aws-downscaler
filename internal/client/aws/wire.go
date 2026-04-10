package aws

import (
	"context"

	"github.com/rs/zerolog/log"
)

func InitAWSDownscalers(downscalerConfig DownscalerConfig) (map[string][]DownscalerI, error) {
	d := map[string][]DownscalerI{}
	for region, services := range downscalerConfig.AwsCloud {

		log.Info().Msgf("Generate aws client for %s region", region)
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
