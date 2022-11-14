package logger

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/fighthorse/readBook/common/setting"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

type (
	logger struct {

		//AppZeroLogger 业务日志句柄（非实时， 100ms刷新一次）
		appZeroLogger zerolog.Logger

		//AccessZeroLogger 访问日志句柄（非实时， 100ms刷新一次）
		accessZeroLogger zerolog.Logger
	}

	M map[string]interface{}
)

var (
	ZeroLogger = &logger{
		appZeroLogger: zerolog.New(os.Stdout).With().Timestamp().Logger(),
	}
)

func Init() {
	{
		// 初始化业务日志
		appCfg := setting.Config.APP
		fileName := appCfg.LogPath
		err := os.MkdirAll(path.Dir(fileName), 0777)
		if err != nil {
			panic(fmt.Errorf("creat log dir error: %s", err))
		}
		f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Errorf("open log file error: %s", err))
		}
		zerolog.TimeFieldFormat = time.RFC3339

		// 1000000 * 1024 /1024/1024 ~= 988M [5m]
		w := diode.NewWriter(f, 1000000, 100*time.Millisecond, func(missed int) {
			logContext := make(map[string]interface{})
			logContext["count"] = missed
			ZeroLogger.appZeroLogger.Log().Fields(logContext).Msg("app_log_miss")
		})

		// level
		l := parseLevel(appCfg.Level)
		ZeroLogger.appZeroLogger = zerolog.New(w).Level(l).With().Timestamp().Logger()
	}

	{
		// 初始化访问日志
		accessCfg := setting.Config.Access
		fileName := accessCfg.FilePath
		err := os.MkdirAll(path.Dir(fileName), 0777)
		if err != nil {
			panic(fmt.Errorf("creat log dir error: %s", err))
		}
		f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Errorf("open log file error: %s", err))
		}
		w := diode.NewWriter(f, 1000000, 100*time.Millisecond, func(missed int) {
			logContext := make(map[string]interface{})
			logContext["count"] = missed
			// Log("access_log_miss", logContext)
			ZeroLogger.appZeroLogger.Log().Fields(logContext).Msg("app_log_miss")
		})
		ZeroLogger.accessZeroLogger = zerolog.New(w).With().Timestamp().Logger()

	}
}

type Fields map[string]interface{}

func (f Fields) putInContext() Fields {
	if len(f) == 0 {
		return f
	}

	ff := make(map[string]interface{}, len(f))
	for key, value := range f {
		switch t := value.(type) {
		case error:
			if t != nil {
				ff[key] = t.Error()
			}
		default:
			ff[key] = value
		}
	}

	m := make(map[string]interface{})
	m["context"] = ff
	return m
}

func Debug(ctx context.Context, msg string, f Fields) {
	withTraceIdLogger(ctx).Debug().Fields(f.putInContext()).Msg(msg)
}

func Info(ctx context.Context, msg string, f Fields) {
	withTraceIdLogger(ctx).Info().Fields(f.putInContext()).Msg(msg)
}

func Warn(ctx context.Context, msg string, f Fields) {
	withTraceIdLogger(ctx).Warn().Fields(f.putInContext()).Msg(msg)
}

func Error(ctx context.Context, msg string, f Fields) {
	withTraceIdLogger(ctx).Error().Fields(f.putInContext()).Msg(msg)
}

func Fatal(ctx context.Context, msg string, f Fields) {
	withTraceIdLogger(ctx).Fatal().Fields(f.putInContext()).Msg(msg)
}

// Panic 期望的panic
type expectedPanic struct{}

var ExpectedPanic = expectedPanic{}

func Panic(ctx context.Context, msg string, f Fields) {
	withTraceIdLogger(ctx).Panic().Fields(f.putInContext()).Msg(msg)
	panic(ExpectedPanic)
}

func Stack(ctx context.Context, msg string, f Fields) {
	f["stacktrace"] = string(debug.Stack())
	withTraceIdLogger(ctx).Info().Fields(f.putInContext()).Msg(msg)
}

func withTraceIdLogger(ctx context.Context) *zerolog.Logger {
	traceId := TraceIdFromCtx(ctx)
	l := ZeroLogger.appZeroLogger.With().Str("trace_id", traceId).Logger()
	return &l
}

func AccessLog(ctx context.Context, fields Fields) {
	traceId := TraceIdFromCtx(ctx)
	l := ZeroLogger.accessZeroLogger.With().Str("trace_id", traceId).Fields(fields).Logger()
	l.Log().Msg("")
}

func AppLog(ctx context.Context) *zerolog.Logger {
	traceId := TraceIdFromCtx(ctx)
	l := ZeroLogger.appZeroLogger.With().Str("trace_id", traceId).Logger()
	return &l
}

func parseLevel(l string) (level zerolog.Level) {
	switch strings.ToUpper(l) {
	case "DEBUG":
		level = zerolog.DebugLevel
	case "INFO":
		level = zerolog.InfoLevel
	case "WARN", "WARNING":
		level = zerolog.WarnLevel
	case "ERROR":
		level = zerolog.ErrorLevel
	case "FATAL":
		level = zerolog.FatalLevel
	case "PANIC":
		level = zerolog.PanicLevel
	case "NIL", "NULL", "DISCARD", "NO":
		level = zerolog.Disabled
	default:
		level = zerolog.InfoLevel
	}

	return
}
