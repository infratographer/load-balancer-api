package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.infratographer.com/x/crdbx"
	"go.infratographer.com/x/otelx"
	"go.infratographer.com/x/viperx"

	"go.infratographer.sh/loadbalancerapi/internal/config"
	"go.infratographer.sh/loadbalancerapi/internal/srv"
)

const auditLogMode = os.FileMode(0644)

// serveCmd starts the TODO service
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "starts the " + appName + " service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return serve(cmd.Context(), viper.GetViper())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("listen", "0.0.0.0:8000", "address to listen on")
	viperx.MustBindFlag(viper.GetViper(), "listen", serveCmd.Flags().Lookup("listen"))

	// App specific flags
	// TODO - add your app specific flags here

	otelx.MustViperFlags(viper.GetViper(), serveCmd.Flags())
	crdbx.MustViperFlags(viper.GetViper(), serveCmd.Flags())

	serveCmd.Flags().String("audit-log-path", "/app-audit/audit.log", "file path to write audit logs to.")
	viperx.MustBindFlag(viper.GetViper(), "audit.log-path", serveCmd.Flags().Lookup("audit-log-path"))
	serveCmd.Flags().Bool("audit-log-disabled", false, "disable audit logging")
	viperx.MustBindFlag(viper.GetViper(), "audit.disabled", serveCmd.Flags().Lookup("audit-log-disabled"))
}

func serve(cmdCtx context.Context, v *viper.Viper) error {
	err := otelx.InitTracer(config.AppConfig.Tracing, appName, logger)
	if err != nil {
		logger.Fatalw("unable to initialize tracing system", "error", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(cmdCtx)

	go func() {
		<-c
		cancel()
	}()

	var auf *os.File

	if !viper.GetBool("audit.disabled") {
		auf, err = openAuditFile()
		if err != nil {
			return err
		}

		defer auf.Close()
	}

	server := &srv.Server{
		DB: srv.DB{
			Driver: initDB(),
			Debug:  config.AppConfig.Logging.Debug,
		},
		Debug:           viper.GetBool("logging.debug"),
		Listen:          viper.GetString("listen"),
		Logger:          logger,
		AuditFileWriter: auf,
	}

	logger.Infow("starting server",
		"address", viper.GetString("listen"),
	)

	if err := server.Run(ctx); err != nil {
		logger.Fatalw("failed starting server", "error", err)
	}

	return nil
}

func openAuditFile() (*os.File, error) {
	path := viper.GetString("audit.log-path")

	if path == "" {
		return nil, ErrAuditFilePathRequired
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, auditLogMode)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func initDB() *sqlx.DB {
	dbDriverName := "postgres"

	sqldb, err := crdbx.NewDB(config.AppConfig.CRDB, config.AppConfig.Tracing.Enabled)
	if err != nil {
		logger.Fatalw("failed to initialize database connection", "error", err)
	}

	db := sqlx.NewDb(sqldb, dbDriverName)

	return db
}
