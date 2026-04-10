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
	Name          string    `yaml:"name,omitempty"`
	Arn           string    `yaml:"arn,omitempty"`
	Tags          []AwsTags `yaml:"tags"`
	Downtime      string    `yaml:"downtime,omitempty"`
	DowntimeState string    `yaml:"downtime_state"`
}

type RegionConfig struct {
	EC2 []AWSResource `yaml:"ec2"`
}

type DownscalerConfig struct {
	AwsCloud        map[string]RegionConfig `yaml:"aws"`
	Interval        string                  `yaml:"interval"`
	Log             string                  `yaml:"log"`
	IntervalSeconds int64
}

const (
	AWS_DOWNSCALER_CONFIG_FILE = "aws-downscaler.yaml"
)

func ReadDownscalerConfig() (DownscalerConfig, error) {
	var cfg DownscalerConfig

	file, err := os.ReadFile(AWS_DOWNSCALER_CONFIG_FILE)
	if err != nil {
		log.Error().Msgf("Can not read config file: %s", AWS_DOWNSCALER_CONFIG_FILE)
		return DownscalerConfig{}, err
	}
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Error().Msgf("Can not unmarshal config file %s", AWS_DOWNSCALER_CONFIG_FILE)
		return DownscalerConfig{}, err
	}

	cfg.IntervalSeconds = util.ParseToSeconds(cfg.Interval)

	return cfg, nil
}
