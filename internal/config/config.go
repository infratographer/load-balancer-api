// Package config provides a struct to store the applications config
package config

import (
	"go.infratographer.com/x/crdbx"
	"go.infratographer.com/x/ginx"
	"go.infratographer.com/x/loggingx"
	"go.infratographer.com/x/otelx"
	"go.infratographer.com/x/versionx"
)

// AppConfig stores all the config values for our application
var AppConfig struct {
	CRDB    crdbx.Config
	Logging loggingx.Config
	Server  ginx.Config
	Tracing otelx.Config
	Details versionx.Details
}
