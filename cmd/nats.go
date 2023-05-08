package cmd

import (
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"

	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/pubsub"
)

func addNatsHooks(ec *ent.Client) (func(), error) {
	opts := []nats.Option{nats.Name(appName)}

	if viper.GetBool("debug") {
		logger.Debug("enabling development settings")

		opts = append(opts, nats.Token(viper.GetString("nats.token")))
	} else {
		opts = append(opts, nats.UserCredentials(viper.GetString("nats.creds-file")))
	}

	nc, err := nats.Connect(viper.GetString("nats.url"), opts...)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	client := pubsub.NewClient(pubsub.WithJetreamContext(js),
		pubsub.WithLogger(logger),
		pubsub.WithStreamName(viper.GetString("nats.stream-name")),
		pubsub.WithSubjectPrefix(viper.GetString("nats.subject-prefix")),
		pubsub.WithEntClient(ec),
	)

	ec.Use(client.Hooks)

	return nc.Close, nil
}
