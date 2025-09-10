package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func addLogger() (*os.File, *zap.Logger) {
	file, err := os.OpenFile("snake.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	fileWriteSyncer := zapcore.AddSync(file)
	consoleWriteSyncer := zapcore.AddSync(os.Stdout)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), fileWriteSyncer, zapcore.InfoLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), consoleWriteSyncer, zapcore.InfoLevel),
	)

	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)

	return file, logger
}

func finishLogging(file *os.File, logger *zap.Logger) {
	logger.Sync()
	file.Close()
}
