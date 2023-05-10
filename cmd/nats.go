package cmd

import (
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"

	"go.infratographer.com/load-balancer-api/internal/pubsub"
)

func configureNatsClient() (*pubsub.Client, error) {
	opts := []nats.Option{nats.Name(appName)}

	if serveDevMode && viper.GetString("nats.creds-file") == "" {
		logger.Debug("enabling development settings")
	} else {
		opts = append(opts, nats.UserCredentials(viper.GetString("nats.creds-file")))
	}

	nc, err := nats.Connect(viper.GetString("nats.url"), opts...)
	if err != nil {
		return &pubsub.Client{}, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return &pubsub.Client{}, err
	}

	client := pubsub.NewClient(pubsub.WithJetreamContext(js),
		pubsub.WithLogger(logger),
		pubsub.WithStreamName(viper.GetString("nats.stream-name")),
		pubsub.WithSubjectPrefix(viper.GetString("nats.subject-prefix")),
	)

	_, err = client.AddStream()
	if err != nil {
		return &pubsub.Client{}, err
	}

	return client, nil
}
