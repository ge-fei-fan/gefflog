package gefflog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

const (
	DEBUG = 1
	INFO  = 2
	WARN  = 4
	ERROR = 8
)

var sugarLogger *zap.SugaredLogger
var encoder zapcore.Encoder
var logsDir = "./logs"
var logsLevel = byte(INFO | ERROR)

func init() {
	getEncoder()

	var debugcore, infocore, warncore, errorcore zapcore.Core
	var allcore []zapcore.Core
	if (logsLevel & DEBUG) != 0 {
		debugcore = initCore(zapcore.DebugLevel)
		allcore = append(allcore, debugcore)
	}
	if (logsLevel & INFO) != 0 {
		infocore = initCore(zapcore.InfoLevel)
		allcore = append(allcore, infocore)
	}
	if (logsLevel & WARN) != 0 {
		warncore = initCore(zapcore.WarnLevel)
		allcore = append(allcore, warncore)
	}
	if (logsLevel & ERROR) != 0 {
		errorcore = initCore(zapcore.ErrorLevel)
		allcore = append(allcore, errorcore)
	}

	core := zapcore.NewTee(allcore...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugarLogger = logger.Sugar()
}
func ChangeLogger(level byte) {
	//Encoder:设置编码器
	getEncoder()
	logsLevel = level

	var debugcore, infocore, warncore, errorcore zapcore.Core
	var allcore []zapcore.Core
	if (logsLevel & DEBUG) != 0 {
		debugcore = initCore(zapcore.DebugLevel)
		allcore = append(allcore, debugcore)
	}
	if (logsLevel & INFO) != 0 {
		infocore = initCore(zapcore.InfoLevel)
		allcore = append(allcore, infocore)
	}
	if (logsLevel & WARN) != 0 {
		warncore = initCore(zapcore.WarnLevel)
		allcore = append(allcore, warncore)
	}
	if (logsLevel & ERROR) != 0 {
		errorcore = initCore(zapcore.ErrorLevel)
		allcore = append(allcore, errorcore)
	}

	core := zapcore.NewTee(allcore...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugarLogger = logger.Sugar()
}

func initCore(level zapcore.Level) zapcore.Core {
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	//info和debug级别,debug级别是最低的
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.InfoLevel && lev >= zap.DebugLevel
	})
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.WarnLevel && lev >= zap.InfoLevel
	})
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= zap.WarnLevel
	})
	//error级别
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.DPanicLevel && lev >= zap.ErrorLevel
	})

	switch level {
	case zapcore.DebugLevel:
		debugSyncer := getLogWriter(zapcore.DebugLevel)
		writeSyncer := zapcore.NewMultiWriteSyncer(consoleDebugging, debugSyncer)
		return zapcore.NewCore(encoder, writeSyncer, debugPriority)
	case zapcore.InfoLevel:
		infoSyncer := getLogWriter(zapcore.InfoLevel)
		return zapcore.NewCore(encoder, infoSyncer, infoPriority)
	case zapcore.WarnLevel:
		warnSyncer := getLogWriter(zapcore.WarnLevel)
		return zapcore.NewCore(encoder, warnSyncer, warnPriority)
	case zapcore.ErrorLevel:
		errorSyncer := getLogWriter(zapcore.ErrorLevel)
		writeSyncer := zapcore.NewMultiWriteSyncer(consoleErrors, errorSyncer)
		return zapcore.NewCore(encoder, writeSyncer, errorPriority)
	}
	return nil
}

// 编码器(如何写入日志)
func getEncoder() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder = zapcore.NewConsoleEncoder(encoderConfig)
}

// 指定日志将写到哪里去
func getLogWriter(lel zapcore.Level) zapcore.WriteSyncer {
	name := "debug.log"
	switch lel {
	case zapcore.DebugLevel:
		name = "debug.log"
	case zapcore.InfoLevel:
		name = "info.log"
	case zapcore.WarnLevel:
		name = "warn.log"
	case zapcore.ErrorLevel:
		name = "error.log"
	default:
		name = "debug.log"
	}
	filename := filepath.Join(logsDir, name)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     10,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func Debug(val ...interface{}) {
	sugarLogger.Debug(val)
	sugarLogger.Sync()
}
func Info(val ...interface{}) {
	sugarLogger.Info(val)
	sugarLogger.Sync()
}
func Warn(val ...interface{}) {
	sugarLogger.Warn(val)
	sugarLogger.Sync()
}
func Err(val ...interface{}) {
	sugarLogger.Error(val)
	sugarLogger.Sync()
}
