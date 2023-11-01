package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// FullPathCallerEncoder encodes the full path and line number of the caller.
func FullPathCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if !caller.Defined {
		enc.AppendString("undefined")
		return
	}

	// from where go.mod is
	enc.AppendString(caller.FullPath())
}

// New constructs a Sugared Logger that writes to stdout and
// provides human-readable timestamps.
func New(service string, outputPaths ...string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeCaller = FullPathCallerEncoder // Set custom caller encoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{ // Fixed map type
		"service": service,
	}

	config.OutputPaths = []string{"stdout"}
	if outputPaths != nil {
		config.OutputPaths = outputPaths
	}

	log, err := config.Build(zap.WithCaller(true))
	if err != nil {
		return nil, err
	}

	return log.Sugar(), nil
}
