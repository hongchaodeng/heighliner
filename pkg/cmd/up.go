package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"helm.sh/helm/pkg/strvals"

	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/schema"
	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
)

var (
	upCmd = &cobra.Command{
		Use:   "up",
		Short: "Run an application",
		Args:  cobra.ArbitraryArgs,
		RunE:  upProj,
	}
)

func init() {
	upCmd.Flags().StringArray("set", []string{}, "The input values of your project")
	if err := viper.BindPFlags(upCmd.Flags()); err != nil {
		log.Fatal().Err(err).Msg("failed to bind flags")
	}
}

func upProj(c *cobra.Command, args []string) error {
	lg := logger.New()

	// Parse --set flag
	values := viper.GetStringSlice("set")
	base := map[string]interface{}{}
	for _, value := range values {
		if err := strvals.ParseInto(value, base); err != nil {
			lg.Fatal().Err(err).Msg("failed parsing --set data")
		}
	}

	// Enter project dir
	if err := state.EnterTemp(); err != nil {
		lg.Fatal().Err(err).Msg("Couldn't find project. Please run hln new to create one")
	}

	// Provide input values
	s := schema.New()
	if err := s.Load(); err != nil {
		lg.Fatal().Err(err).Msg("failed to load input schema")
	}
	if err := s.SetEnv(base); err != nil {
		lg.Fatal().Err(err).Msg("failed to set input values")
	}

	return util.Exec("dagger",
		"--log-format", viper.GetString("log-format"),
		"--log-level", viper.GetString("log-level"),
		"-p", "./plans",
		"do", "up")
}
