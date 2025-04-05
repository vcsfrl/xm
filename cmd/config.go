package cmd

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vcsfrl/xm/internal/config"
)

func buildConfig(logger zerolog.Logger) *config.Config {
	logger.Info().Msg("Initialize api config.")

	var newConfig config.Config
	newConfig.TracePort = viper.Get("tracePort").(string)
	newConfig.AuthUser = viper.Get("authUser").(string)
	newConfig.AuthPassword = viper.Get("authPassword").(string)
	newConfig.AuthJwtSecret = viper.Get("authJwtSecret").(string)
	newConfig.DbPath = viper.Get("dbPath").(string)
	newConfig.AppPort = viper.Get("appPort").(string)
	newConfig.RateLimit = viper.GetFloat64("rateLimit")
	newConfig.RateBurst = viper.GetInt("rateBurst")

	return &newConfig
}

func bindEnvConfig(command *cobra.Command) error {
	command.Flags().String("trace-port", "8090", "Trace port")
	if err := viper.BindPFlag("tracePort", command.Flags().Lookup("trace-port")); err != nil {
		return err
	}
	if err := viper.BindEnv("tracePort", "XM_TRACE_PORT"); err != nil {
		return err
	}

	command.Flags().String("auth-user", "", "Auth user")
	if err := viper.BindPFlag("authUser", command.Flags().Lookup("auth-user")); err != nil {
		return err
	}
	if err := viper.BindEnv("authUser", "XM_API_AUTH_USER"); err != nil {
		return err
	}

	command.Flags().String("auth-password", "", "Auth password")
	if err := viper.BindPFlag("authPassword", command.Flags().Lookup("auth-password")); err != nil {
		return err
	}
	if err := viper.BindEnv("authPassword", "XM_API_AUTH_USER"); err != nil {
		return err
	}

	command.Flags().String("auth-jwt-secret", "", "Auth auth jwt secret")
	if err := viper.BindPFlag("authJwtSecret", command.Flags().Lookup("auth-jwt-secret")); err != nil {
		return err
	}
	if err := viper.BindEnv("authJwtSecret", "XM_API_AUTH_JWT_SECRET"); err != nil {
		return err
	}

	command.Flags().String("app-port", "xm", "App port")
	if err := viper.BindPFlag("appPort", command.Flags().Lookup("app-port")); err != nil {
		return err
	}
	if err := viper.BindEnv("appPort", "XM_APP_PORT"); err != nil {
		return err
	}

	command.Flags().String("db-path", "/tmp/xm.db", "Db path")
	if err := viper.BindPFlag("dbPath", command.Flags().Lookup("db-path")); err != nil {
		return err
	}
	if err := viper.BindEnv("dbPath", "XM_DB_PATH"); err != nil {
		return err
	}

	command.Flags().Float64("rate-limit", 1.0, "Rate limit")
	if err := viper.BindPFlag("rateLimit", command.Flags().Lookup("rate-limit")); err != nil {
		return err
	}
	if err := viper.BindEnv("rateLimit", "XM_RATE_LIMIT"); err != nil {
		return err
	}

	command.Flags().Int("rate-burst", 1, "Rate burst")
	if err := viper.BindPFlag("rateBurst", command.Flags().Lookup("rate-burst")); err != nil {
		return err
	}
	if err := viper.BindEnv("rateBurst", "XM_RATE_BURST"); err != nil {
		return err
	}

	return nil
}
