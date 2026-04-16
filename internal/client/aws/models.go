package aws

import (
	"os"

	"aws-downscaler/internal/util"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert/yaml"
)

type AwsTags struct {
	Name  string `yaml:"name,omitempty"`
	Value string `yaml:"value,omitempty"`
}

type AWSResource struct {
	Name                 string    `yaml:"name,omitempty"`
	Arn                  string    `yaml:"arn,omitempty"`
	Tags                 []AwsTags `yaml:"tags"`
	Downtime             string    `yaml:"downtime,omitempty"`
	DowntimeState        string    `yaml:"downtime_state"`
	UptimeReplicaCount   string    `yaml:"uptime_replica_count,omitempty"`
	DowntimeReplicaCount string    `yaml:"downtime_replica_count,omitempty"`
}

type RegionConfig struct {
	EC2         []AWSResource `yaml:"ec2"`
	Elasticache []AWSResource `yaml:"elasticache"`
}

type DownscalerConfig struct {
	AwsCloud        map[string]RegionConfig `yaml:"aws"`
	Interval        string                  `yaml:"interval"`
	Log             string                  `yaml:"log"`
	IntervalSeconds int64
}

func ReadDownscalerConfig(configPath string) (DownscalerConfig, error) {
	var cfg DownscalerConfig

	log.Info().Msgf("Reading config file from %s", configPath)
	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Error().Msgf("Can not read config file: %s", configPath)
		return DownscalerConfig{}, err
	}
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Error().Msgf("Can not unmarshal config file %s", configPath)
		return DownscalerConfig{}, err
	}

	cfg.IntervalSeconds = util.ParseToSeconds(cfg.Interval)

	return cfg, nil
}
