package cmd

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.infratographer.com/permissions-api/pkg/permissions"
	"go.infratographer.com/x/crdbx"
	"go.infratographer.com/x/echojwtx"
	"go.infratographer.com/x/echox"
	"go.infratographer.com/x/oauth2x"
	"go.infratographer.com/x/otelx"
	"go.infratographer.com/x/versionx"
	"go.infratographer.com/x/viperx"
	"go.uber.org/zap"

	metadata "go.infratographer.com/metadata-api/pkg/client"

	"go.infratographer.com/load-balancer-api/internal/config"
	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/graphapi"
	"go.infratographer.com/load-balancer-api/internal/manualhooks"

	"go.infratographer.com/x/events"
)

const (
	defaultLBAPIListenAddr = ":7608"
	shutdownTimeout        = 10 * time.Second
	defaultTimeout         = 5 * time.Second
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
	echojwtx.MustViperFlags(viper.GetViper(), serveCmd.Flags())
	events.MustViperFlags(viper.GetViper(), serveCmd.Flags(), appName)
	oauth2x.MustViperFlags(viper.GetViper(), serveCmd.Flags())
	permissions.MustViperFlags(viper.GetViper(), serveCmd.Flags())

	serveCmd.Flags().String("metadata-status-namespace-id", "", "status namespace id to update loadbalancer metadata status")
	viperx.MustBindFlag(viper.GetViper(), "metadata.status-namespace-id", serveCmd.Flags().Lookup("metadata-status-namespace-id"))

	serveCmd.Flags().Duration("supergraph-timeout", defaultTimeout, "client timeout")
	viperx.MustBindFlag(viper.GetViper(), "supergraph.timeout", serveCmd.Flags().Lookup("supergraph-timeout"))

	serveCmd.Flags().String("supergraph-url", "", "endpoint for supergraph gateway")
	viperx.MustBindFlag(viper.GetViper(), "supergraph.url", serveCmd.Flags().Lookup("supergraph-url"))

	// only available as a CLI arg because it shouldn't be something that could accidentially end up in a config file or env var
	serveCmd.Flags().BoolVar(&serveDevMode, "dev", false, "dev mode: enables playground, disables all auth checks, sets CORS to allow all, pretty logging, etc.")
	serveCmd.Flags().BoolVar(&enablePlayground, "playground", false, "enable the graph playground")
	serveCmd.Flags().StringVar(&pidFileName, "pid-file", "", "path to the pid file")
	serveCmd.Flags().IntSlice("restricted-ports", []int{}, "ports that are restricted from being used by the load balancer (e.g. 22, 8086, etc.)")
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
	var resolverOpts []graphapi.Option

	config.AppConfig.LoadBalancerLimit = viper.GetInt("load-balancer-limit")

	if serveDevMode {
		enablePlayground = true
		config.AppConfig.Logging.Debug = true
		config.AppConfig.Logging.Pretty = true
		config.AppConfig.Server.WithMiddleware(middleware.CORS())
	}

	events, err := events.NewConnection(config.AppConfig.Events, events.WithLogger(logger))
	if err != nil {
		logger.Fatalw("failed to initialize events", "error", err)
	}

	err = otelx.InitTracer(config.AppConfig.Tracing, appName, logger)
	if err != nil {
		logger.Fatalw("failed to initialize tracer", "error", err)
	}

	db, err := crdbx.NewDB(config.AppConfig.CRDB, config.AppConfig.Tracing.Enabled)
	if err != nil {
		logger.Fatalw("failed to connect to database", "error", err)
	}

	defer db.Close()

	entDB := entsql.OpenDB(dialect.Postgres, db)

	cOpts := []ent.Option{ent.Driver(entDB), ent.EventsPublisher(events)}

	if config.AppConfig.Logging.Debug {
		cOpts = append(cOpts,
			ent.Log(logger.Named("ent").Debugln),
			ent.Debug(),
		)
	}

	client := ent.NewClient(cOpts...)
	defer client.Close()

	// TODO - @rizzza - supergraph client
	var metadataClient *metadata.Client

	if config.AppConfig.Supergraph.URL != "" {
		if config.AppConfig.OIDCClient.Config.Issuer != "" {
			oidcTS, err := oauth2x.NewClientCredentialsTokenSrc(ctx, config.AppConfig.OIDCClient.Config)
			if err != nil {
				logger.Fatalw("failed to create oauth2 token source", "error", err)
			}

			oauthHTTPClient := oauth2x.NewClient(ctx, oidcTS)
			oauthHTTPClient.Timeout = config.AppConfig.Supergraph.Timeout

			metadataClient = metadata.New(config.AppConfig.Supergraph.URL,
				metadata.WithHTTPClient(oauthHTTPClient),
			)
		} else {
			metadataClient = metadata.New(config.AppConfig.Supergraph.URL)
		}
	}

	if metadataClient != nil {
		resolverOpts = append(resolverOpts, graphapi.WithMetadataClient(metadataClient))
	}

	// TODO: fix generated pubsubhooks
	// eventhooks.PubsubHooks(client)

	manualhooks.PubsubHooks(client)

	// Run the automatic migration tool to create all schema resources.
	if err := client.Schema.Create(ctx); err != nil {
		logger.Errorf("failed creating schema resources", zap.Error(err))
		return err
	}

	config.AppConfig.RestrictedPorts = viper.GetIntSlice("restricted-ports")

	var middleware []echo.MiddlewareFunc

	// jwt auth middleware
	if viper.GetBool("oidc.enabled") {
		auth, err := echojwtx.NewAuth(ctx, config.AppConfig.OIDC)
		if err != nil {
			logger.Fatalw("failed to initialize jwt authentication", zap.Error(err))
		}

		middleware = append(middleware, auth.Middleware())
	}

	srv, err := echox.NewServer(logger.Desugar(), config.AppConfig.Server, versionx.BuildDetails(), echox.WithLoggingSkipper(echox.SkipDefaultEndpoints))
	if err != nil {
		logger.Fatalw("failed to create server", zap.Error(err))
	}

	perms, err := permissions.New(config.AppConfig.Permissions,
		permissions.WithLogger(logger),
		permissions.WithDefaultChecker(permissions.DefaultAllowChecker),
		permissions.WithEventsPublisher(events),
	)
	if err != nil {
		logger.Fatal("failed to initialize permissions", zap.Error(err))
	}

	middleware = append(middleware, perms.Middleware())

	r := graphapi.NewResolver(client, logger.Named("resolvers"), resolverOpts...)
	handler := r.Handler(enablePlayground, middleware...)

	srv.AddHandler(handler)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		_ = events.Shutdown(ctx)
	}()

	go func() {
		if err := srv.Run(); err != nil {
			logger.Fatal("failed to run server", zap.Error(err))
		}
	}()

	select {
	case <-shutdown:
		logger.Info("signal caught, shutting down")
	case <-ctx.Done():
		logger.Info("context done, shutting down")
	}

	return nil
}
