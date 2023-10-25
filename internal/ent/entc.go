//go:build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"go.infratographer.com/x/entx"
)

func main() {
	xExt, err := entx.NewExtension(
		entx.WithFederation(),
		// entx.WithEventHooks(), // TODO: untangle additional subjects coupled to auth relationship
		entx.WithJSONScalar(),
	)
	if err != nil {
		log.Fatalf("creating entx extension: %v", err)
	}

	// load-balancer-api currently uses manual maintained hooks, however, if the schema changes, we'll need to enable the pubsubExt
	// to re-generate the hooks and merge the generated changes with the manual hooks.

	// pubsubExt, err := pubsubinfo.NewExtension(
	// 	pubsubinfo.WithEventHooks(),
	// )
	// if err != nil {
	// 	log.Fatalf("creating pubsubinfo extension: %v", err)
	// }

	gqlExt, err := entgql.NewExtension(
		// Tell Ent to generate a GraphQL schema for
		// the Ent schema in a file named ent.graphql.
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath("schema/ent.graphql"),
		entgql.WithConfigPath("gqlgen.yml"),
		entgql.WithWhereInputs(true),
		entgql.WithSchemaHook(xExt.GQLSchemaHooks()...),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}

	opts := []entc.Option{
		entc.Extensions(
			xExt,
			gqlExt,
			// pubsubExt,
		),
		entc.TemplateDir("./internal/ent/templates"),
		entc.FeatureNames("intercept"),
		entc.Dependency(
			entc.DependencyName("EventsPublisher"),
			entc.DependencyTypeInfo(&field.TypeInfo{
				Ident:   "events.Connection",
				PkgPath: "go.infratographer.com/x/events",
			}),
		),
	}

	if err := entc.Generate("./internal/ent/schema", &gen.Config{
		Target:   "./internal/ent/generated",
		Package:  "go.infratographer.com/load-balancer-api/internal/ent/generated",
		Header:   entx.CopyrightHeader,
		Features: []gen.Feature{gen.FeatureVersionedMigration},
	}, opts...); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
