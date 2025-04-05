package cmd

import (
	"github.com/vcsfrl/xm/cmd/example"
	"github.com/vcsfrl/xm/internal/config"
	dbFactory "github.com/vcsfrl/xm/internal/db"
	"gorm.io/gorm"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logger zerolog.Logger
var loggerOutput zerolog.ConsoleWriter
var appConfig *config.Config
var db *gorm.DB

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xm",
	Short: "XM technical task",
	Long:  ``,
}

// apiCmd represents the runEvent command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Run api.",
	Long:  `Start rest API.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

var exampleCmd = &cobra.Command{
	Use:   "api-example",
	Short: "Run api example calls.",
	Long:  `Start rest API.`,
	Run: func(cmd *cobra.Command, args []string) {
		example.Run(appConfig, logger)
	},
}

func init() {
	// Init logger.
	loggerOutput = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger = zerolog.New(loggerOutput).With().Timestamp().Logger()
	logger.Info().Msg("Logger initialised.")

	viper.SetConfigName("xm_config")
	viper.SetEnvPrefix("XM")

	if err := bindEnvConfig(apiCmd); err != nil {
		logger.Error().Err(err).Msg("Bind monitor config.")
		os.Exit(1)
	}

	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(exampleCmd)

	// Init config.
	appConfig = buildConfig(logger)

	// Init database.
	var err error
	db, err = dbFactory.InitSqlite(appConfig, logger)
	if err != nil {
		logger.Error().Err(err).Msg("Init db.")
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
