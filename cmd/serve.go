package cmd

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.infratographer.com/x/crdbx"
	"go.infratographer.com/x/otelx"
	"go.infratographer.com/x/viperx"

	"go.infratographer.com/loadbalancerapi/internal/config"
	"go.infratographer.com/loadbalancerapi/internal/x/echox"
	"go.infratographer.com/loadbalancerapi/pkg/api/v1"
)

var defaultLBAPIListenAddr = ":7608"

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the load balancer API",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("listen", "l", defaultLBAPIListenAddr, "The address to listen on")
	viperx.MustBindFlag(viper.GetViper(), "listen", serveCmd.Flags().Lookup("listen"))
}

func serve(ctx context.Context) {
	err := otelx.InitTracer(config.AppConfig.Tracing, appName, logger)
	if err != nil {
		logger.Fatalw("failed to initialize tracer", "error", err)
	}

	db, err := crdbx.NewDB(config.AppConfig.CRDB, config.AppConfig.Tracing.Enabled)
	if err != nil {
		logger.Fatalw("failed to connect to database", "error", err)
	}

	dbx := sqlx.NewDb(db, "postgres")

	e := echox.NewServer()
	r := api.NewRouter(dbx, logger)

	e.Debug = true
	r.Routes(e)

	e.Logger.Fatal(e.Start(viper.GetString("listen")))
}
