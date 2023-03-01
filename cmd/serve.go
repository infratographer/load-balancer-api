package cmd

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.infratographer.com/x/crdbx"
	"go.infratographer.com/x/otelx"
	"go.infratographer.com/x/viperx"

	"go.infratographer.com/load-balancer-api/internal/config"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
	"go.infratographer.com/load-balancer-api/internal/x/echox"
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

	js, natsClose, err := newJetstreamConnection()
	if err != nil {
		logger.Fatalw("failed to create NATS jetstream connection", "error", err)
	}

	defer natsClose()

	e := echox.NewServer()
	r := api.NewRouter(
		dbx,
		logger,
		pubsub.NewClient(
			pubsub.WithJetreamContext(js),
			pubsub.WithLogger(logger),
			pubsub.WithStreamName(viper.GetString("nats.stream-name")),
			pubsub.WithSubjectPrefix(viper.GetString("nats.subject-prefix")),
		),
	)

	e.Debug = true
	r.Routes(e)

	e.Logger.Fatal(e.Start(viper.GetString("listen")))
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
