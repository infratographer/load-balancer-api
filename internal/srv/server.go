package srv

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	gindump "github.com/tpkeeper/gin-dump"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.hollow.sh/toolbox/ginjwt"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	v1router "go.infratographer.sh/loadbalancerapi/pkg/api/v1/router"
)

// DB is a wrapper around sqlx.DB
type DB struct {
	Driver *sqlx.DB
	Debug  bool
}

// Server implements the HTTP Server
type Server struct {
	Logger          *zap.SugaredLogger
	Listen          string
	Debug           bool
	DB              DB
	AuthConfig      ginjwt.AuthConfig
	AuditFileWriter io.Writer
}

var (
	readTimeout     = 10 * time.Second
	writeTimeout    = 20 * time.Second
	corsMaxAge      = 12 * time.Hour
	shutdownTimeout = 5 * time.Second
)

func (s *Server) setup() *gin.Engine {
	var (
		authMW *ginjwt.Middleware
		err    error
	)

	authMW, err = ginjwt.NewAuthMiddleware(s.AuthConfig)
	if err != nil {
		s.Logger.Fatal("failed to initialize auth middleware", "error", err)
	}

	// Setup default gin router
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowAllOrigins:  true,
		AllowCredentials: true,
		MaxAge:           corsMaxAge,
	}))

	p := ginprometheus.NewPrometheus("gin")

	// Remove any params from the URL string to keep the number of labels down
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		return c.FullPath()
	}

	p.Use(r)

	logF := func(c *gin.Context) []zapcore.Field {
		return []zapcore.Field{
			zap.String("jwt_subject", ginjwt.GetSubject(c)),
			zap.String("jwt_user", ginjwt.GetUser(c)),
		}
	}
	loggerWithContext := s.Logger.With(zap.String("component", "httpsrv"))
	r.Use(ginzap.GinzapWithConfig(loggerWithContext.Desugar(), &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		TraceID:    true,
		Context:    logF,
	}))
	r.Use(ginzap.RecoveryWithZap(loggerWithContext.Desugar(), true))

	tp := otel.GetTracerProvider()
	if tp != nil {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}

		r.Use(otelgin.Middleware(hostname, otelgin.WithTracerProvider(tp)))
	}

	// Health endpoints
	r.GET("/healthz", s.livenessCheck)
	r.GET("/healthz/liveness", s.livenessCheck)
	r.GET("/healthz/readiness", s.readinessCheck)

	v1Rtr := v1router.New(authMW, s.DB.Driver, s.Logger)

	// Host our latest version of the API under / in addition to /api/v*
	latest := r.Group("/")
	{
		v1Rtr.Routes(latest)
		if s.Debug {
			r.Use(gindump.Dump())
		}
	}

	v1 := r.Group(v1router.V1URI)
	{
		v1Rtr.Routes(v1)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "invalid request - route not found"})
	})

	return r
}

// NewServer returns a configured server
func (s *Server) NewServer() *http.Server {
	if !s.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	return &http.Server{
		Handler:      s.setup(),
		Addr:         s.Listen,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

// Run will start the server listening on the specified address
func (s *Server) Run(ctx context.Context) error {
	httpsrv := s.NewServer()

	go func() {
		if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer func() {
		cancel()
	}()

	if err := httpsrv.Shutdown(ctxShutDown); err != nil {
		return err
	}

	s.Logger.Info("server shutdown cleanly", zap.String("time", time.Now().UTC().Format(time.RFC3339)))

	return nil
}

// livenessCheck ensures that the server is up and responding
func (s *Server) livenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}

// readinessCheck ensures that the server is up and that we are able to process requests.
func (s *Server) readinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}
