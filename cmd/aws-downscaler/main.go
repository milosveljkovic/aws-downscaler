package main

import (
	downscaler "aws-downscaler/internal/client/aws"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	downscalerConfig, err := downscaler.ReadDownscalerConfig()
	if err != nil {
		os.Exit(1)
	}
	setupLogger(downscalerConfig.Log)

	downscalers := map[string][]downscaler.DownscalerI{}

	downscalers, err = downscaler.InitAWSDownscalers(downscalerConfig)

	for true {
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
