package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Parameters struct {
	Level zapcore.Level
}

func SetupLogger(params Parameters) (*zap.Logger, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("can't inititalize zap logger: %v", err)
	}
	logDir := fmt.Sprintf("%s/logs", cwd)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("can't create log directory: %v", err)
	}
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/logs.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	})
	consoleWriter := zapcore.AddSync(os.Stdout)
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), consoleWriter, params.Level),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), fileWriter, params.Level),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	logger.Info("Logger successfully started")
	return logger, nil
}
