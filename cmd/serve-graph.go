package cmd

import (
	"context"
	"log"

	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.infratographer.com/x/echox"
	"go.infratographer.com/x/versionx"
	"go.uber.org/zap"

	"go.infratographer.com/load-balancer-api/internal/config"
	ent "go.infratographer.com/load-balancer-api/internal/ent/generated"
	"go.infratographer.com/load-balancer-api/internal/graphapi"
)

var serveGraphCmd = &cobra.Command{
	Use:   "serve-graph",
	Short: "Start the load balancer Graph API",
	Run: func(cmd *cobra.Command, args []string) {
		servegraph(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(serveGraphCmd)

	echox.MustViperFlags(viper.GetViper(), serveGraphCmd.Flags(), defaultLBAPIListenAddr)
}

func servegraph(ctx context.Context) {
	cOpts := []ent.Option{}

	if config.AppConfig.Logging.Debug {
		cOpts = append(cOpts,
			ent.Log(logger.Named("ent").Debugln),
			ent.Debug(),
		)
	}

	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1", cOpts...)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	// Run the automatic migration tool to create all schema resources.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
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
		logger.Fatal("failed to create server", zap.Error(err))
	}

	handler := graphapi.NewHandler(client)

	srv.AddHandler(handler)

	if err := srv.Run(); err != nil {
		logger.Fatal("failed to run server", zap.Error(err))
	}
}
