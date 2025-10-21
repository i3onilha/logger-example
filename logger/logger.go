package logger

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type logEntry struct {
	level  zapcore.Level
	ctx    context.Context
	msg    string
	fields []zap.Field
}

var (
	log           *zap.Logger
	asyncCh       chan logEntry
	once          sync.Once
	config        Config
	batchInterval = 100 * time.Millisecond // flush batch every 100ms
)

// Init logger and start async worker
func Init(cfg Config) {
	once.Do(func() {
		config = cfg
		var zapCfg zap.Config
		if cfg.Encoding == "console" {
			zapCfg = zap.NewDevelopmentConfig()
		} else {
			zapCfg = zap.NewProductionConfig()
		}

		level := zapcore.InfoLevel
		if err := level.Set(cfg.Level); err != nil {
			panic(fmt.Errorf("invalid log level: %w", err))
		}
		zapCfg.Level = zap.NewAtomicLevelAt(level)
		zapCfg.OutputPaths = []string{"stdout"}
		zapCfg.InitialFields = map[string]interface{}{
			"service": cfg.Service,
			"env":     cfg.Environment,
		}

		var err error
		log, err = zapCfg.Build()
		if err != nil {
			panic(fmt.Errorf("failed to build zap logger: %w", err))
		}

		asyncCh = make(chan logEntry, cfg.AsyncBufferSize)
		go asyncWorker()
	})
}

// Async worker with batching
func asyncWorker() {
	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()
	batch := make([]logEntry, 0, config.BatchSize)

	flush := func() {
		for _, entry := range batch {
			writeLog(entry)
		}
		batch = batch[:0]
	}

	for {
		select {
		case entry := <-asyncCh:
			batch = append(batch, entry)
			if len(batch) >= config.BatchSize {
				flush()
			}
		case <-ticker.C:
			if len(batch) > 0 {
				flush()
			}
		}
	}
}

func writeLog(entry logEntry) {
	traceFields := extractTraceFields(entry.ctx)
	fields := append(entry.fields, traceFields...)

	switch entry.level {
	case zapcore.InfoLevel:
		log.Info(entry.msg, fields...)
	case zapcore.WarnLevel:
		log.Warn(entry.msg, fields...)
	case zapcore.ErrorLevel:
		log.Error(entry.msg, fields...)
	case zapcore.DebugLevel:
		log.Debug(entry.msg, fields...)
	}
}

// extract trace info from context
func extractTraceFields(ctx context.Context) []zap.Field {
	span, ok := tracer.SpanFromContext(ctx)
	if !ok || span == nil {
		return nil
	}
	return []zap.Field{
		zap.String("dd.trace_id", fmt.Sprintf("%d", span.Context().TraceID())),
		zap.String("dd.span_id", fmt.Sprintf("%d", span.Context().SpanID())),
	}
}

// Public logging functions (non-blocking)
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	asyncCh <- logEntry{level: zapcore.InfoLevel, ctx: ctx, msg: msg, fields: fields}
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	asyncCh <- logEntry{level: zapcore.WarnLevel, ctx: ctx, msg: msg, fields: fields}
}

func Error(ctx context.Context, err error, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("error", err.Error()))
	asyncCh <- logEntry{level: zapcore.ErrorLevel, ctx: ctx, msg: msg, fields: fields}
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if config.Level == "debug" {
		asyncCh <- logEntry{level: zapcore.DebugLevel, ctx: ctx, msg: msg, fields: fields}
	}
}

func GetLogger() *zap.Logger {
	return log
}
