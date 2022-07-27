package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultLogLevel = zapcore.InfoLevel

	defaultDevelopmentLevel = "debug"
	defaultProductionLevel  = "info"

	defaultDevelopmentEncoding = "console"
	defaultProductionEncoding  = "json"

	stderr = "stderr"
	stdout = "stdout"

	Development = "development"
	Production  = "production"
)

// LevelNamesMap stores mapping from string loglevel to zapcore.Level
var LevelNamesMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// NewLogger creates new zap.Logger
// if pathToLogs is empty, stdout will be used
func NewLogger(logLevel string, pathsToLogs []string, encoding string) (*zap.Logger, error) {
	zapLogLevel, ok := LevelNamesMap[strings.ToLower(logLevel)]
	if !ok {
		// set default log level
		zapLogLevel = defaultLogLevel
	}

	if len(pathsToLogs) == 0 {
		pathsToLogs = []string{stderr}
	}

	loggerConfig := zap.Config{
		Level:    zap.NewAtomicLevelAt(zapLogLevel),
		Encoding: encoding,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths:      pathsToLogs,
		ErrorOutputPaths: pathsToLogs,
	}

	return loggerConfig.Build()
}

// NewSugarLogger creates new zap.SugaredLogger
func NewSugarLogger(logLevel string, pathsToLogs []string, encoding string) (*zap.SugaredLogger, error) {
	loggerInstance, err := NewLogger(logLevel, pathsToLogs, encoding)
	if err != nil {
		return nil, err
	}

	return loggerInstance.Sugar(), nil
}

// NewDevelopmentSugarLogger creates new zap.SugaredLogger to use it during development
func NewDevelopmentSugarLogger(pathsToLogs []string) (*zap.SugaredLogger, error) {
	// add stderr to logger paths in development logger
	isStdOutOrErr := false
	for _, path := range pathsToLogs {
		if isStdOutOrErr = path == stdout || path == stderr; isStdOutOrErr {
			break
		}
	}
	if !isStdOutOrErr {
		pathsToLogs = append(pathsToLogs, stderr)
	}

	// create base zap logger
	loggerInstance, err := NewLogger(defaultDevelopmentLevel, pathsToLogs, defaultDevelopmentEncoding)
	if err != nil {
		return nil, err
	}

	// add development option
	return loggerInstance.WithOptions(zap.Development()).Sugar(), nil
}

// NewProductionSugarLogger creates new zap.SugaredLogger to use it in production
func NewProductionSugarLogger(pathsToLogs []string) (*zap.SugaredLogger, error) {
	return NewSugarLogger(defaultProductionLevel, pathsToLogs, defaultProductionEncoding)
}

// NewEnvironmentSugarLogger creates new zap.SugaredLogger for environment
func NewEnvironmentSugarLogger(environment string, pathsToLogs []string) (*zap.SugaredLogger, error) {
	switch strings.ToLower(environment) {
	case Development:
		return NewDevelopmentSugarLogger(pathsToLogs)
	case Production:
		return NewProductionSugarLogger(pathsToLogs)
	default:
		return NewProductionSugarLogger(pathsToLogs)
	}
}
