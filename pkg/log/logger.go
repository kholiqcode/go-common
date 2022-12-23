package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/kholiqcode/go-common/pkg/constants"
	common_utils "github.com/kholiqcode/go-common/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	sugarLogger *zap.SugaredLogger
	logger      *zap.Logger
}

func NewLogger(level string, encoding string, opts ...string) *Logger {
	if level == "" {
		level = "debug"
	}

	if encoding == "" {
		encoding = "console"
	}

	levelmap := map[string]zapcore.Level{
		"info":  zap.InfoLevel,
		"debug": zap.DebugLevel,
		"warn":  zap.WarnLevel,
		"error": zap.ErrorLevel,
		"panic": zap.PanicLevel,
		"fatal": zap.FatalLevel,
	}

	levelKey := strings.ToLower(level)
	levelValue, ok := levelmap[levelKey]
	if !ok {
		levelValue = zapcore.DebugLevel
	}

	outputPath := ""
	if len(opts) > 0 {
		outputPath = opts[0]
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if outputPath != "" {
		// clean previous data, regardless it's file or path
		os.RemoveAll(outputPath)
		// recreate the folder (not the file) to be written
		dir := filepath.Dir(outputPath)
		if _, err := os.Stat(dir); err != nil {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				common_utils.PanicIfError(fmt.Errorf("failed to create log folder: %v", err))
			}
		}
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, outputPath)
		zapConfig.ErrorOutputPaths = append(zapConfig.ErrorOutputPaths, outputPath)
	}

	zapConfig.Level = zap.NewAtomicLevelAt(levelValue)
	zapConfig.Encoding = encoding
	zapLogger, _ := zapConfig.Build()

	return &Logger{sugarLogger: zapLogger.Sugar(), logger: zapLogger}
}

func (l *Logger) Infow(msg string, args ...interface{}) {
	l.sugarLogger.Infow(msg, args...)
	if err := l.sugarLogger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		fmt.Printf("error : %s", err.Error())
	}
}

func (l *Logger) Debugw(msg string, args ...interface{}) {
	l.sugarLogger.Debugw(msg, args...)
	if err := l.sugarLogger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		fmt.Printf("error : %s", err.Error())
	}
}

func (l *Logger) Warnw(msg string, args ...interface{}) {
	l.sugarLogger.Warnw(msg, args...)
	if err := l.sugarLogger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		fmt.Printf("error : %s", err.Error())
	}
}

func (l *Logger) Errorw(msg string, args ...interface{}) {
	l.sugarLogger.Errorw(msg, args...)
	if err := l.sugarLogger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		fmt.Printf("error : %s", err.Error())
	}
}

func (l *Logger) Panicw(msg string, args ...interface{}) {
	l.sugarLogger.Panicw(msg, args...)
	if err := l.sugarLogger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		fmt.Printf("error : %s", err.Error())
	}
}

func (l *Logger) Fatalw(msg string, args ...interface{}) {
	l.sugarLogger.Fatalw(msg, args...)
	if err := l.sugarLogger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		fmt.Printf("error : %s", err.Error())
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.sugarLogger.Infof(format, args...)
	if err := l.sugarLogger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		fmt.Printf("error : %s", err.Error())
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.sugarLogger.Debugf(format, args...)
	if err := l.sugarLogger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		fmt.Printf("error : %s", err.Error())
	}
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.sugarLogger.Warnf(format, args...)
	if err := l.sugarLogger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		fmt.Printf("error : %s", err.Error())
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)

	l.sugarLogger.Errorf(format, args...)

	if err := l.sugarLogger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		fmt.Printf("error : %s", err.Error())
	}

	return errors.New(msg)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.sugarLogger.Panicf(format, args...)
	if err := l.sugarLogger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		fmt.Printf("error : %s", err.Error())
	}
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.sugarLogger.Fatalf(format, args...)
	if err := l.sugarLogger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		fmt.Printf("error : %s", err.Error())
	}
}

func (l *Logger) HttpMiddlewareAccessLogger(method, uri string, status int, size int64, time time.Duration) {
	l.logger.Info(
		constants.HTTP,
		zap.String(constants.METHOD, method),
		zap.String(constants.URI, uri),
		zap.Int(constants.STATUS, status),
		zap.Int64(constants.SIZE, size),
		zap.Duration(constants.TIME, time),
	)
}

func (l *Logger) GrpcMiddlewareAccessLogger(method string, time time.Duration, metaData map[string][]string, err error) {
	l.logger.Info(
		constants.GRPC,
		zap.String(constants.METHOD, method),
		zap.Duration(constants.TIME, time),
		zap.Any(constants.METADATA, metaData),
		zap.Any(constants.ERROR, err),
	)
}

func (l *Logger) GrpcMiddlewareAccessLoggerErr(method string, time time.Duration, metaData map[string][]string, err error) {
	l.logger.Error(
		constants.GRPC,
		zap.String(constants.METHOD, method),
		zap.Duration(constants.TIME, time),
		zap.Any(constants.METADATA, metaData),
		zap.Any(constants.ERROR, err),
	)
}

func (l *Logger) GrpcClientInterceptorLogger(method string, req, reply interface{}, time time.Duration, metaData map[string][]string, err error) {
	l.logger.Info(
		constants.GRPC,
		zap.String(constants.METHOD, method),
		zap.Any(constants.REQUEST, req),
		zap.Any(constants.REPLY, reply),
		zap.Duration(constants.TIME, time),
		zap.Any(constants.METADATA, metaData),
		zap.Any(constants.ERROR, err),
	)
}

func (l *Logger) GrpcClientInterceptorLoggerErr(method string, req, reply interface{}, time time.Duration, metaData map[string][]string, err error) {
	l.logger.Error(
		constants.GRPC,
		zap.String(constants.METHOD, method),
		zap.Any(constants.REQUEST, req),
		zap.Any(constants.REPLY, reply),
		zap.Duration(constants.TIME, time),
		zap.Any(constants.METADATA, metaData),
		zap.Any(constants.ERROR, err),
	)
}

func (l *Logger) KafkaProcessMessage(topic string, partition int, message []byte, workerID int, offset int64, time time.Time) {
	l.logger.Debug(
		"(Processing Kafka message)",
		zap.String(constants.Topic, topic),
		zap.Int(constants.Partition, partition),
		zap.Int(constants.MessageSize, len(message)),
		zap.Int(constants.WorkerID, workerID),
		zap.Int64(constants.Offset, offset),
		zap.Time(constants.Time, time),
	)
}

func (l *Logger) KafkaLogCommittedMessage(topic string, partition int, offset int64) {
	l.logger.Debug(
		"(Committed Kafka message)",
		zap.String(constants.Topic, topic),
		zap.Int(constants.Partition, partition),
		zap.Int64(constants.Offset, offset),
	)
}

func (l *Logger) KafkaProcessMessageWithHeaders(topic string, partition int, message []byte, workerID int, offset int64, time time.Time, headers map[string]interface{}) {
	l.logger.Debug(
		"(Processing Kafka message)",
		zap.String(constants.Topic, topic),
		zap.Int(constants.Partition, partition),
		zap.Int(constants.MessageSize, len(message)),
		zap.Int(constants.WorkerID, workerID),
		zap.Int64(constants.Offset, offset),
		zap.Time(constants.Time, time),
		zap.Any(constants.KafkaHeaders, headers),
	)
}
