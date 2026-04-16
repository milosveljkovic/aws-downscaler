package aws

import (
	"aws-downscaler/internal/util"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/rs/zerolog/log"
)

type ElasticacheDownscaler struct {
	cli    *elasticache.Client
	groups []AWSResource
}

func NewElasticacheDownscaler(awsConfig aws.Config, awsResources []AWSResource) (DownscalerI, error) {
	log.Info().Msg("Generate elasticache client")
	svc := elasticache.NewFromConfig(awsConfig)
	return ElasticacheDownscaler{
		cli:    svc,
		groups: awsResources,
	}, nil
}

func (elasticacheCli ElasticacheDownscaler) Downscale() {
	for _, elasticacheGroup := range elasticacheCli.groups {
		log.Info().Msgf("Processing elasticache group %s ...", elasticacheGroup.Name)
		isDowntime, err := util.IsNowInDowntime(elasticacheGroup.Downtime)
		if err != nil {
			log.Error().Msg(err.Error())
			// log.Warn().Msgf("Seems like downtime(%s) of ec2 group %s can not be proceeded - skip", ec2Group.Downtime, ec2Group.Name)
			continue
		}
		if isDowntime == true {
			log.Info().Msg("DOWNSCALING ...")

		} else {
			log.Info().Msg("UPSCALING ...")
		}
	}
}
