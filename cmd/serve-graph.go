package cmd

import (
	"context"

	"entgo.io/ent/dialect"
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

const defaultLBAPIListenAddr = ":7608"

var serveGraphCmd = &cobra.Command{
	Use:   "serve-graph",
	Short: "Start the load balancer Graph API",
	RunE: func(cmd *cobra.Command, args []string) error {
		return servegraph(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(serveGraphCmd)

	echox.MustViperFlags(viper.GetViper(), serveGraphCmd.Flags(), defaultLBAPIListenAddr)
}

func servegraph(ctx context.Context) error {
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

	srv, err := echox.NewServer(
		logger.Desugar(),
		echox.Config{
			Listen:              viper.GetString("server.listen"),
			ShutdownGracePeriod: viper.GetDuration("server.shutdown-grace-period"),
		},
		versionx.BuildDetails(),
	)
	if err != nil {
		logger.Error("failed to create server", zap.Error(err))
	}

	handler := graphapi.NewHandler(client)

	srv.AddHandler(handler)

	if err := srv.Run(); err != nil {
		logger.Error("failed to run server", zap.Error(err))
	}

	return err
}
