package cmd

import (
	"context"
	"os"
	"strconv"
	"syscall"

	"entgo.io/ent/dialect"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.infratographer.com/x/echox"
	"go.infratographer.com/x/versionx"
	"go.uber.org/zap"

	"go.infratographer.com/load-balancer-api/internal/config"
	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/graphapi"
)

const (
	defaultLBAPIListenAddr = ":7608"
)

var (
	enablePlayground bool
	serveDevMode     bool
	pidFileName      = "/tmp/lba.pid"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the load balancer Graph API",
	RunE: func(cmd *cobra.Command, args []string) error {
		if pidFileName != "" {
			if err := writePidFile(pidFileName); err != nil {
				logger.Error("failed to write pid file", zap.Error(err))
				return err
			}

			defer os.Remove(pidFileName)
		}

		return serve(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	echox.MustViperFlags(viper.GetViper(), serveCmd.Flags(), defaultLBAPIListenAddr)

	// only available as a CLI arg because it shouldn't be something that could accidentially end up in a config file or env var
	serveCmd.Flags().BoolVar(&serveDevMode, "dev", false, "dev mode: enables playground, disables all auth checks, sets CORS to allow all, pretty logging, etc.")
	serveCmd.Flags().BoolVar(&enablePlayground, "playground", false, "enable the graph playground")
	serveCmd.Flags().StringVar(&pidFileName, "pid-file", "", "path to the pid file")
}

// Write a pid file, but first make sure it doesn't exist with a running pid.
func writePidFile(pidFile string) error {
	// Read in the pid file as a slice of bytes.
	if piddata, err := os.ReadFile(pidFile); err == nil {
		// Convert the file contents to an integer.
		if pid, err := strconv.Atoi(string(piddata)); err == nil {
			// Look for the pid in the process list.
			if process, err := os.FindProcess(pid); err == nil {
				// Send the process a signal zero kill.
				if err := process.Signal(syscall.Signal(0)); err == nil {
					// We only get an error if the pid isn't running, or it's not ours.
					return err
				}
			}
		}
	}

	logger.Debugw("writing pid file", "pid-file", pidFile)

	// If we get here, then the pidfile didn't exist,
	// or the pid in it doesn't belong to the user running this app.
	return os.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0o664) // nolint: gomnd
}

func serve(ctx context.Context) error {
	if serveDevMode {
		enablePlayground = true
		config.AppConfig.Logging.Debug = true
		config.AppConfig.Logging.Pretty = true
		config.AppConfig.Server.WithMiddleware(middleware.CORS())
	}

	cOpts := []ent.Option{}

	if config.AppConfig.Logging.Debug {
		cOpts = append(cOpts,
			ent.Log(logger.Named("ent").Debugln),
			ent.Debug(),
		)
	}

	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1", cOpts...)
	if err != nil {
		logger.Error("failed opening connection to sqlite", zap.Error(err))
		return err
	}
	defer client.Close()

	// Run the automatic migration tool to create all schema resources.
	if err := client.Schema.Create(ctx); err != nil {
		logger.Errorf("failed creating schema resources", zap.Error(err))
		return err
	}

	// Add hooks for nats events to the ent client. They function as middleware between mutators.
	natsClose, err := addNatsHooks(client)
	if err != nil {
		logger.Errorw("failed to add ent client hooks for nats events", "error", err)
		return err
	}

	defer natsClose()

	srv, err := echox.NewServer(logger.Desugar(), config.AppConfig.Server, versionx.BuildDetails())
	if err != nil {
		logger.Error("failed to create server", zap.Error(err))
	}

	r := graphapi.NewResolver(client, logger.Named("resolvers"))
	handler := r.Handler(enablePlayground)

	srv.AddHandler(handler)

	if err := srv.RunWithContext(ctx); err != nil {
		logger.Error("failed to run server", zap.Error(err))
	}

	return err
}
