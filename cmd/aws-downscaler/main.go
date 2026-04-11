package main

import (
	downscaler "aws-downscaler/internal/client/aws"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	configFile string
	awsProfile string
	version    = "dev"
)

const (
	AWS_DOWNSCALER_CONFIG_FILE = "aws-downscaler.yaml"
)

func main() {
	parseFlags()

	if err := run(); err != nil {
		log.Error().Msgf("AWS Downscaler failed: %s", err)
		os.Exit(1)
	}
}

func run() error {
	log.Info().Msgf("AWS Downscaler version: %s", version)

	downscalerConfig, err := downscaler.ReadDownscalerConfig(configFile)
	if err != nil {
		return fmt.Errorf("read config %q: %w", configFile, err)
	}
	setupLogger(downscalerConfig.Log)

	downscalers := map[string][]downscaler.DownscalerI{}

	downscalers, err = downscaler.InitAWSDownscalers(downscalerConfig, awsProfile)
	if err != nil {
		return fmt.Errorf("initialize AWS downscalers: %w", err)
	}

	if len(downscalers) == 0 {
		return fmt.Errorf("no clients defined")
	}

	for {
		for region, clients := range downscalers {
			log.Info().Msgf("Processing region %s ...", region)
			for _, c := range clients {
				c.Downscale()
			}
		}
		log.Info().Msg("---")
		log.Info().Msgf("Waiting %d seconds before next processing ... ", downscalerConfig.IntervalSeconds)
		time.Sleep(time.Second * time.Duration(downscalerConfig.IntervalSeconds))
	}
}

func parseFlags() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s [options] (version:%s)\n\n", os.Args[0], version)

		fmt.Fprintf(os.Stderr, "Description:\n")
		fmt.Fprintf(os.Stderr, "  AWS Downscaler - scale EC2 instances up or down based on config and tags.\n\n")

		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  %s -config config.yaml\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -config config.yaml -aws-profile example \n\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "AWS auth options:\n")
		fmt.Fprintf(os.Stderr, "  -aws-profile example  Use AWS profile (~/.aws/config)\n")
		fmt.Fprintf(os.Stderr, "  Env vars            	AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY\n")
		fmt.Fprintf(os.Stderr, "  Kubernetes (IRSA)   	Use service account IAM role (no profile/env)\n\n")

		fmt.Fprintf(os.Stderr, "Required IAM permissions:\n")
		fmt.Fprintf(os.Stderr, "  ec2:DescribeInstances\n")
		fmt.Fprintf(os.Stderr, "  ec2:StartInstances\n")
		fmt.Fprintf(os.Stderr, "  ec2:StopInstances\n\n")

		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.StringVar(&configFile, "config", AWS_DOWNSCALER_CONFIG_FILE, "AWS Downscaler config file path")
	flag.StringVar(&awsProfile, "aws-profile", "", "AWS profile to use")

	flag.Parse()
}

func setupLogger(logLevel string) {
	ll := getLogLevel(logLevel)
	zerolog.SetGlobalLevel(ll)
}

func getLogLevel(logLevel string) zerolog.Level {
	ll := zerolog.InfoLevel
	logLevelUpper := "INFO"

	if logLevel == "" {
		log.Info().Msgf("Setting up default log level to %s", ll)
		return ll
	}
	logLevelUpper = strings.ToUpper(logLevel)
	switch {
	case logLevelUpper == "DEBUG":
		ll = zerolog.DebugLevel
	case logLevelUpper == "INFO":
		ll = zerolog.InfoLevel
	case logLevelUpper == "WARN":
		ll = zerolog.WarnLevel
	case logLevelUpper == "ERROR":
		ll = zerolog.ErrorLevel
	default:
		log.Info().Msgf("Provided non supported log level '%s', using default %s", logLevel, zerolog.InfoLevel)
		ll = zerolog.InfoLevel
	}
	log.Info().Msgf("Setting up log level to %s", ll)
	return ll
}
