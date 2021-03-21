package mylog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"runtime"
	"time"
)

type (
	Func   func(string, ...zap.Field)
	Field  = zap.Field
	Level  = zapcore.Level
)

type Logger struct {
	desugar *zap.Logger
	sugar   *zap.SugaredLogger
	lv      *zap.AtomicLevel
}

var DefaultLogger *Logger

func NewLogger() *Logger {
	l := &Logger{}
	logger, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		panic("init logger error.")
	}

	l.desugar = logger
	l.sugar = logger.Sugar()

	DefaultLogger = l

	return l
}

// SetLevel ...
func (logger *Logger) SetLevel(lv Level) {
	logger.lv.SetLevel(lv)
}

// Flush ...
func (logger *Logger) Flush() error {
	return logger.desugar.Sync()
}

// DefaultZapConfig ...
func DefaultZapConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lv",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.Unix())
}

// IsDebugMode ...
func (logger *Logger) IsDebugMode() bool {
	return false
}

func normalizeMessage(msg string) string {
	return fmt.Sprintf("%-32s", msg)
}

// Debug ...
func (logger *Logger) Debug(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Debug(msg, fields...)
}

// Debugw ...
func (logger *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Debugw(msg, keysAndValues...)
}

func sprintf(template string, args ...interface{}) string {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	return msg
}

// StdLog ...
func (logger *Logger) StdLog() *log.Logger {
	return zap.NewStdLog(logger.desugar)
}

// Debugf ...
func (logger *Logger) Debugf(template string, args ...interface{}) {
	logger.sugar.Debugw(sprintf(template, args...))
}

// Info ...
func (logger *Logger) Info(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Info(msg, fields...)
}

// Infow ...
func (logger *Logger) Infow(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Infow(msg, keysAndValues...)
}

// Infof ...
func (logger *Logger) Infof(template string, args ...interface{}) {
	logger.sugar.Infof(sprintf(template, args...))
}

// Warn ...
func (logger *Logger) Warn(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Warn(msg, fields...)
}

// Warnw ...
func (logger *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Warnw(msg, keysAndValues...)
}

// Warnf ...
func (logger *Logger) Warnf(template string, args ...interface{}) {
	logger.sugar.Warnf(sprintf(template, args...))
}

// Error ...
func (logger *Logger) Error(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Error(msg, fields...)
}

// Errorw ...
func (logger *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Errorw(msg, keysAndValues...)
}

// Errorf ...
func (logger *Logger) Errorf(template string, args ...interface{}) {
	logger.sugar.Errorf(sprintf(template, args...))
}

// Panic ...
func (logger *Logger) Panic(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		panicDetail(msg, fields...)
		msg = normalizeMessage(msg)
	}
	logger.desugar.Panic(msg, fields...)
}

// Panicw ...
func (logger *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Panicw(msg, keysAndValues...)
}

// Panicf ...
func (logger *Logger) Panicf(template string, args ...interface{}) {
	logger.sugar.Panicf(sprintf(template, args...))
}

// DPanic ...
func (logger *Logger) DPanic(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		panicDetail(msg, fields...)
		msg = normalizeMessage(msg)
	}
	logger.desugar.DPanic(msg, fields...)
}

// DPanicw ...
func (logger *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.DPanicw(msg, keysAndValues...)
}

// DPanicf ...
func (logger *Logger) DPanicf(template string, args ...interface{}) {
	logger.sugar.DPanicf(sprintf(template, args...))
}

// Fatal ...
func (logger *Logger) Fatal(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		panicDetail(msg, fields...)
		msg = normalizeMessage(msg)
		return
	}
	logger.desugar.Fatal(msg, fields...)
}

// Fatalw ...
func (logger *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Fatalw(msg, keysAndValues...)
}

// Fatalf ...
func (logger *Logger) Fatalf(template string, args ...interface{}) {
	logger.sugar.Fatalf(sprintf(template, args...))
}

func panicDetail(msg string, fields ...Field) {
	enc := zapcore.NewMapObjectEncoder()
	for _, field := range fields {
		field.AddTo(enc)
	}

	// 控制台输出
	fmt.Printf("%s: \n    %s: %s\n", "panic", "msg", msg)
	if _, file, line, ok := runtime.Caller(3); ok {
		fmt.Printf("    %s: %s:%d\n", "loc", file, line)
	}
	for key, val := range enc.Fields {
		fmt.Printf("    %s: %s\n", key, fmt.Sprintf("%+v", val))
	}

}

// With ...
func (logger *Logger) With(fields ...Field) *Logger {
	desugarLogger := logger.desugar.With(fields...)
	return &Logger{
		desugar: desugarLogger,
		sugar:   desugarLogger.Sugar(),
	}
}
