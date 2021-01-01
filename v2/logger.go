package micro

import (
	"fmt"
	"os"

	"google.golang.org/grpc/grpclog"

	"github.com/rs/zerolog"
)

type logger struct {
	log zerolog.Logger
}

// NewLogger creates a grpc logger using zerolog
func NewLogger(serviceName string, level zerolog.Level) grpclog.LoggerV2 {
	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Caller().
		Str("protocol", "grpc").
		Str("service_name", serviceName).
		Logger().
		Level(level)
	return &logger{log: log}
}

func (l *logger) Info(args ...interface{}) {
	l.log.Info().Msg(fmt.Sprint(args...))
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.log.Info().Msg(fmt.Sprintf(format, args...))
}

func (l *logger) Infoln(args ...interface{}) {
	l.Info(args...)
}

func (l *logger) Warning(args ...interface{}) {
	l.log.Warn().Msg(fmt.Sprint(args...))
}

func (l *logger) Warningf(format string, args ...interface{}) {
	l.log.Warn().Msg(fmt.Sprintf(format, args...))
}

func (l *logger) Warningln(args ...interface{}) {
	l.Warning(args...)
}

func (l *logger) Error(args ...interface{}) {
	l.log.Error().Msg(fmt.Sprint(args...))
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.log.Error().Msg(fmt.Sprintf(format, args...))
}

func (l *logger) Errorln(args ...interface{}) {
	l.Error(args...)
}
func (l *logger) Fatal(args ...interface{}) {
	l.log.Fatal().Msg(fmt.Sprint(args...))
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.log.Fatal().Msg(fmt.Sprintf(format, args...))
}

func (l *logger) Fatalln(args ...interface{}) {
	l.Fatal(args...)
}

func (l *logger) V(level int) bool {
	return true
}
