package aws

import (
	"aws-downscaler/internal/util"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/rs/zerolog/log"
)

type DownscalerI interface {
	Downscale()
}

type Ec2Downscaler struct {
	cli    *ec2.Client
	groups []AWSResource
}

func NewEc2Downscaler(awsConfig aws.Config, awsResources []AWSResource) (DownscalerI, error) {
	log.Info().Msg("Generate ec2 client")
	svc := ec2.NewFromConfig(awsConfig)
	return Ec2Downscaler{
		cli:    svc,
		groups: awsResources,
	}, nil
}

func (ec2cli Ec2Downscaler) Downscale() {
	for _, ec2Group := range ec2cli.groups {
		log.Info().Msgf("Processing ec2 group %s ...", ec2Group.Name)
		isDowntime, err := util.IsNowInDowntime(ec2Group.Downtime)
		if err != nil {
			log.Error().Msg(err.Error())
			log.Warn().Msgf("Seems like downtime(%s) of ec2 group %s can not be proceeded, skip.", ec2Group.Downtime, ec2Group.Name)
			continue
		}
		if isDowntime == true {
			log.Info().Msg("DOWNSCALING ...")
			instanceIds := ec2cli.getEC2ByTag(ec2Group, ActionDownscale)
			for _, instanceID := range instanceIds {
				// TODO: add support for hibernate state
				log.Info().Msgf("Stopping '%s' ...", instanceID)
				_, err := ec2cli.cli.StopInstances(context.TODO(), &ec2.StopInstancesInput{
					InstanceIds: []string{instanceID},
				})
				if err != nil {
					log.Error().Msg(err.Error())
					continue
				}
			}
		} else {
			log.Info().Msg("UPSCALING ...")
			instanceIds := ec2cli.getEC2ByTag(ec2Group, ActionUpscale)
			for _, instanceID := range instanceIds {

				log.Info().Msgf("Starting '%s' ...", instanceID)
				_, err := ec2cli.cli.StartInstances(context.TODO(), &ec2.StartInstancesInput{
					InstanceIds: []string{instanceID},
				})
				if err != nil {
					log.Error().Msg(err.Error())
					continue
				}
			}
		}
	}
}

func (ec2cli Ec2Downscaler) BuildEC2Filters(tags []AwsTags) []types.Filter {
	var filters []types.Filter

	for _, tag := range tags {
		filters = append(filters, types.Filter{
			Name:   aws.String("tag:" + tag.Name),
			Values: []string{tag.Value},
		})
	}

	return filters
}

type Action string

const (
	ActionUpscale   Action = "upscale"
	ActionDownscale Action = "downscale"
)

func (ec2cli Ec2Downscaler) getEC2ByTag(group AWSResource, action Action) []string {
	input := &ec2.DescribeInstancesInput{
		Filters: ec2cli.BuildEC2Filters(group.Tags),
	}

	paginator := ec2.NewDescribeInstancesPaginator(ec2cli.cli, input)

	var filteredInstances []string
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			log.Error().Msgf("Failed to describe instances: %v", err)
			return []string{}
		}

		for _, reservation := range page.Reservations {
			for _, instance := range reservation.Instances {
				instanceId := aws.ToString(instance.InstanceId)
				if action == ActionDownscale && isInstanceRunning(instance.State.Name) {
					log.Info().Msgf("Instance '%s' is in '%s' state - %s", instanceId, instance.State.Name, action)
					filteredInstances = append(filteredInstances, instanceId)
					continue
				} else if action == ActionUpscale && isInstanceStoppe(instance.State.Name) {
					log.Info().Msgf("Instance '%s' is in '%s' state - %s", instanceId, instance.State.Name, action)
					filteredInstances = append(filteredInstances, instanceId)
					continue
				}
				log.Info().Msgf("Instance '%s' is in '%s' state - SKIP %s", instanceId, instance.State.Name, action)
			}
		}
	}
	return filteredInstances
}

func isInstanceStoppe(state types.InstanceStateName) bool {
	if state == types.InstanceStateNameStopped {
		return true
	}
	return false
}

func isInstanceRunning(state types.InstanceStateName) bool {
	if state == types.InstanceStateNameRunning {
		return true
	}
	return false
}
