package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

func main() {
	fmt.Printf("*** Build a logger from a json\n\n")

	rawJSONConfig := []byte(`{
	"level": "info",
	"encoding": "console",
	"outputPaths": ["stdout", "/tmp/logs"],
	"errorOutputPaths": ["/tmp/errorlogs"],
	"initialFields": {"initFieldKey": "fieldValue"},
	"encoderConfig": {
		"messageKey": "message",
		"levelKey": "level",
		"nameKey": "logger",
		"timeKey": "time",
		"callerKey": "logger",
		"stacktraceKey": "stacktrace",
		"callstackKey": "callstack",
		"errorKey": "error",
		"timeEncoder": "iso8601",
		"fileKey": "file",
		"levelEncoder": "capitalColor",
		"durationEncoder": "second",
		"callerEncoder": "full",
		"nameEncoder": "full",
		"sampling": {
			"initial": "3",
			"thereafter": "10"
		}
	}
}`)

	config := zap.Config{}
	if err := json.Unmarshal(rawJSONConfig, &config); err != nil {
		panic(err)
	}
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	logger.Debug("This is a DEBUG message")
	logger.Info("This should have an ISO8601 based time stamp")
	logger.Warn("This is a WARN message")
	logger.Error("This is an ERROR message")
	//logger.Fatal("This is a FATAL message")   // would exit if uncommented
	//logger.DPanic("This is a DPANIC message") // would exit if uncommented

	const url = "http://example.com"
	logger.Info("Failed to fetch URL.",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	fmt.Printf("\n*** Using a JSON encoder, at debug level, sending output to stdout, no key specified\n\n")

	logger, _ = zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"stdout"},
	}.Build()

	logger.Debug("This is a DEBUG message")
	logger.Info("This is an INFO message")
	logger.Info("This is an INFO message with fields", zap.String("region", "us-west"), zap.Int("id", 2))

	fmt.Printf("\n*** Using a JSON encoder, at debug level, sending output to stdout, message key only specified\n\n")

	logger, _ = zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
		},
	}.Build()

	logger.Debug("This is a DEBUG message")
	logger.Info("This is an INFO message")
	logger.Info("This is an INFO message with fields", zap.String("region", "us-west"), zap.Int("id", 2))

	fmt.Printf("\n*** Using a JSON encoder, at debug level, sending output to stdout, all possible keys specified\n\n")

	cfg := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, _ = cfg.Build()

	logger.Debug("This is a DEBUG message")
	logger.Info("This is an INFO message")
	logger.Info("This is an INFO message with fields", zap.String("region", "us-west"), zap.Int("id", 2))

	fmt.Printf("\n*** Same logger with console logging enabled instead\n\n")

	logger.WithOptions(
		zap.WrapCore(func(zapcore.Core) zapcore.Core {
			return zapcore.NewCore(zapcore.NewConsoleEncoder(cfg.EncoderConfig), zapcore.AddSync(os.Stderr), zapcore.DebugLevel)
		}),
	).Info("This is an INFO message")
}
