// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"code.cn/blog/cmd/flags"
	"code.cn/blog/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	globalLogger atomic.Value
	once         sync.Once
)

// Init initializes global logger (safe to call multiple times, only runs once)
func Init() {
	once.Do(func() {
		cfg := conf.Get().Logger

		if err := os.MkdirAll(cfg.LogsDir, 0755); err != nil {
			panic(fmt.Errorf("failed to create log directory: %w", err))
		}

		atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
		if flags.Debug {
			atomicLevel.SetLevel(zap.DebugLevel)
		}

		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

		jsonEncoder := zapcore.NewJSONEncoder(encoderCfg)

		newWriter := func(filename string) zapcore.WriteSyncer {
			return zapcore.AddSync(&lumberjack.Logger{
				Filename:   filepath.Join(cfg.LogsDir, filename),
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
				Compress:   true,
				LocalTime:  true,
			})
		}

		var cores []zapcore.Core

		// info core
		cores = append(cores, zapcore.NewCore(
			jsonEncoder,
			newWriter("info.log"),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return atomicLevel.Enabled(lvl) && lvl < zapcore.WarnLevel
			}),
		))

		// error core
		cores = append(cores, zapcore.NewCore(
			jsonEncoder,
			newWriter("error.log"),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return atomicLevel.Enabled(lvl) && lvl >= zapcore.WarnLevel
			}),
		))

		// console (debug only)
		if flags.Debug {
			consoleEnc := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
			cores = append(cores, zapcore.NewCore(
				consoleEnc,
				zapcore.AddSync(os.Stdout),
				atomicLevel,
			))
		}

		core := zapcore.NewTee(cores...)

		opts := []zap.Option{
			zap.AddCaller(),
			zap.AddCallerSkip(2),
			zap.AddStacktrace(zapcore.ErrorLevel),
		}

		if flags.Debug {
			opts = append(opts, zap.Development())
		}

		logger := zap.New(core, opts...)
		globalLogger.Store(logger)
	})
}

func L() *zap.Logger {
	if v := globalLogger.Load(); v != nil {
		return v.(*zap.Logger)
	}
	return zap.NewNop()
}

func Close() {
	if v := globalLogger.Load(); v != nil {
		_ = v.(*zap.Logger).Sync()
	}
}

func toFields(args ...any) []zap.Field {
	fields := make([]zap.Field, 0, len(args))

	for _, arg := range args {
		switch v := arg.(type) {
		case zap.Field:
			fields = append(fields, v)

		case error:
			if v != nil {
				fields = append(fields, zap.Error(v))
			}

		default:
			// ignore unknown types or extend later
		}
	}

	return fields
}

func log(level zapcore.Level, msg string, args ...any) {
	L().Check(level, msg).Write(toFields(args...)...)
}

func Error(msg string, args ...any) {
	log(zapcore.ErrorLevel, msg, args...)
}

func Warn(msg string, args ...any) {
	log(zapcore.WarnLevel, msg, args...)
}

func Info(msg string, args ...any) {
	log(zapcore.InfoLevel, msg, args...)
}
