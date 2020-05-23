package log

import (
	"context"

	"github.com/TheZeroSlave/zapsentry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// -----------------------------------------------------------------------------

// Options declares logger options for builder
type Options struct {
	Debug       bool
	LogLevel    string
	Encoding    string
	AppName     string
	Environment string
	Version     string
	Revision    string
	SentryDSN   string
}

// -----------------------------------------------------------------------------

// DefaultOptions defines default logger options
var DefaultOptions = &Options{
	Debug:       false,
	LogLevel:    "info",
	Encoding:    "console",
	AppName:     "changeme",
	Environment: "changeme",
	Version:     "0.0.1",
	Revision:    "123456789",
	SentryDSN:   "",
}

// -----------------------------------------------------------------------------

// Setup the logger
func Setup(ctx context.Context, opts Options) {
	// Initialize logs
	var config zap.Config

	if opts.Debug {
		config = zap.NewDevelopmentConfig()
		config.DisableCaller = true
		config.DisableStacktrace = true
	} else {
		config = zap.NewProductionConfig()
		config.DisableStacktrace = true
		config.EncoderConfig.MessageKey = "@message"
		config.EncoderConfig.TimeKey = "@timestamp"
		config.EncoderConfig.CallerKey = "@caller"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	config.Encoding = opts.Encoding

	if opts.Encoding == "console" {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Parse log level
	errLogLevel := config.Level.UnmarshalText([]byte(opts.LogLevel))
	if errLogLevel != nil {
		panic(errLogLevel)
	}

	// Build real logger
	logger, err := config.Build(
		zap.AddCallerSkip(2),
	)
	if err != nil {
		panic(err)
	}

	// Add prefix to logger
	logger = logger.With(
		zap.String("@appName", opts.AppName),
		zap.String("@version", opts.Version),
		zap.String("@revision", opts.Revision),
		zap.String("@Enviroment", opts.Environment),
		zap.Namespace("@fields"),
	)

	// sentry support
	if opts.SentryDSN != "" {
		logger.Info("Starting sentry collector", zap.String("dsn", opts.SentryDSN))

		cfg := zapsentry.Configuration{
			Level: zapcore.ErrorLevel, // when to send message to sentry
			Tags: map[string]string{
				"application.name":       opts.AppName,
				"application.enviroment": opts.Environment,
				"version":                opts.Version,
				"revision":               opts.Revision,
			},
		}
		core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromDSN(opts.SentryDSN))
		if err != nil {
			For(ctx).Warn("Unable to attach sentry to logger", zap.Error(err))
		}

		logger = zapsentry.AttachCoreToLogger(core, logger)
	} else {
		logger.Info("Sentry collector disabled")
	}

	// Prepare factory
	logFactory := NewFactory(logger)

	// Override the global factory
	SetLoggerFactory(logFactory)

	// Override zap default logger
	zap.ReplaceGlobals(logger)
}
