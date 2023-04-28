package cmd

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.infratographer.com/x/crdbx"
	"go.infratographer.com/x/echojwtx"
	"go.infratographer.com/x/echox"
	"go.infratographer.com/x/otelx"
	"go.infratographer.com/x/versionx"
	"go.uber.org/zap"

	"go.infratographer.com/load-balancer-api/internal/config"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.infratographer.com/load-balancer-api/pkg/api/v1"
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

	echox.MustViperFlags(viper.GetViper(), serveCmd.Flags(), defaultLBAPIListenAddr)
}

func serve(ctx context.Context) {
	var middleware []echo.MiddlewareFunc

	err := otelx.InitTracer(config.AppConfig.Tracing, appName, logger)
	if err != nil {
		logger.Fatalw("failed to initialize tracer", "error", err)
	}

	db, err := crdbx.NewDB(config.AppConfig.CRDB, config.AppConfig.Tracing.Enabled)
	if err != nil {
		logger.Fatalw("failed to connect to database", "error", err)
	}

	dbx := sqlx.NewDb(db, "postgres")

	js, natsClose, err := newJetstreamConnection()
	if err != nil {
		logger.Fatalw("failed to create NATS jetstream connection", "error", err)
	}

	defer natsClose()

	srv, err := echox.NewServer(
		logger.Desugar(),
		echox.ConfigFromViper(viper.GetViper()),
		versionx.BuildDetails(),
	)
	if err != nil {
		logger.Fatal("failed to initialize new server", zap.Error(err))
	}

	if config, err := echojwtx.AuthConfigFromViper(viper.GetViper()); err != nil {
		logger.Fatal("failed to initialize jwt authentication", zap.Error(err))
	} else if config != nil {
		config.JWTConfig.Skipper = echox.SkipDefaultEndpoints

		auth, err := echojwtx.NewAuth(ctx, *config)
		if err != nil {
			logger.Fatal("failed to initialize jwt authentication", zap.Error(err))
		}

		middleware = append(middleware, auth.Middleware())
	}

	r := api.NewRouter(
		dbx,
		pubsub.NewClient(
			pubsub.WithJetreamContext(js),
			pubsub.WithLogger(logger),
			pubsub.WithStreamName(viper.GetString("nats.stream-name")),
			pubsub.WithSubjectPrefix(viper.GetString("nats.subject-prefix")),
		),
		api.WithLogger(logger.Named("api").Desugar()),
		api.WithMiddleware(middleware...),
	)

	srv.AddHandler(r).AddReadinessCheck("database", r.DatabaseCheck)

	if err := srv.Run(); err != nil {
		logger.Fatal("failed to run server", zap.Error(err))
	}
}

func newJetstreamConnection() (nats.JetStreamContext, func(), error) {
	opts := []nats.Option{nats.Name(appName)}

	if viper.GetBool("debug") {
		logger.Debug("enabling development settings")

		opts = append(opts, nats.Token(viper.GetString("nats.token")))
	} else {
		opts = append(opts, nats.UserCredentials(viper.GetString("nats.creds-file")))
	}

	nc, err := nats.Connect(viper.GetString("nats.url"), opts...)
	if err != nil {
		return nil, nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, nil, err
	}

	return js, nc.Close, nil
}
