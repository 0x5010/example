package main

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	file := "tmp.txt"
	logger := NewLogger(file)

	logger.Debug("This is a DEBUG message")
	logger.Info("This is an INFO message")
	logger.Info("This is an INFO message with fields", zap.String("region", "us-west"), zap.Int("id", 2))

	slogger := logger.Sugar()
	slogger.Info("Info() uses sprint")
	logger.Info("Info() uses sprint with fields", zap.String("region", "us-west"), zap.Int("id", 2))
	slogger.Infof("Infof() uses %s", "sprintf")
	slogger.Infow("Infow() allows tags", "name", "Legolas", "type", 1)
}

func NewLogger(filename string) *zap.Logger {
	output := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10, // megabytes
		MaxBackups: 10,
		MaxAge:     0, //days
	}
	fileWriter := zapcore.AddSync(output)
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = logTimeEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encCfg),
		fileWriter,
		zapcore.InfoLevel,
	)
	return zap.New(core)
}

func logTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.0000"))
}
