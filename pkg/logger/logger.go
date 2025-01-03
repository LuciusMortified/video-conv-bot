package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger *zap.Logger

func init() {
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.InfoLevel ||
			level == zapcore.WarnLevel ||
			level == zapcore.DebugLevel
	})

	errorFatalLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.ErrorLevel ||
			level == zapcore.FatalLevel
	})

	stdoutSyncer := zapcore.Lock(os.Stdout)
	stderrSyncer := zapcore.Lock(os.Stderr)

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			stdoutSyncer,
			infoLevel,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			stderrSyncer,
			errorFatalLevel,
		),
	)

	logger = zap.New(core)
}

func Flush() {
	_ = logger.Sync()
}

type Field struct {
	Key   string
	Value any
}

func NewField[T any](k string, v T) Field {
	return Field{Key: k, Value: v}
}

func makeZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	return zapFields
}

func Info(msg string, fields ...Field) {
	logger.Info(msg, makeZapFields(fields)...)
}

func Debug(msg string, fields ...Field) {
	logger.Debug(msg, makeZapFields(fields)...)
}

func Fatal(msg string, fields ...Field) {
	logger.Fatal(msg, makeZapFields(fields)...)
}

func Error(msg string, fields ...Field) {
	logger.Error(msg, makeZapFields(fields)...)
}
