package logger

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/colinmarc/hdfs/v2"
	"github.com/rs/zerolog"
)

type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type Logger struct {
	logger *zerolog.Logger
}

var _ Interface = (*Logger)(nil)

// New — создаёт логгер только в stdout (без HDFS).
// Используется на старте до подключения к HDFS.
func New(level string) *Logger {
	return newLogger(level, os.Stdout)
}

// NewWithHDFS — создаёт логгер с дублированием в HDFS.
// Вызывается из app.Run после подключения к HDFS.
func NewWithHDFS(level string, hdfsClient *hdfs.Client, serviceName string) *Logger {
	levels := []string{"debug", "info", "warn", "error", "fatal"}

	// Создаём writers для каждого уровня
	writers := make([]io.Writer, 0, len(levels)+1)
	writers = append(writers, os.Stdout) // всегда пишем в stdout

	for _, lvl := range levels {
		writers = append(writers, NewHDFSRotatingWriter(hdfsClient, serviceName, lvl))
	}

	// LevelWriter — роутит каждую запись в нужный HDFS writer по уровню
	multi := &levelRouter{
		stdout: os.Stdout,
		hdfsWriters: map[string]*HDFSRotatingWriter{
			"debug": NewHDFSRotatingWriter(hdfsClient, serviceName, "debug"),
			"info":  NewHDFSRotatingWriter(hdfsClient, serviceName, "info"),
			"warn":  NewHDFSRotatingWriter(hdfsClient, serviceName, "warn"),
			"error": NewHDFSRotatingWriter(hdfsClient, serviceName, "error"),
			"fatal": NewHDFSRotatingWriter(hdfsClient, serviceName, "fatal"),
		},
	}

	return newLogger(level, multi)
}

// levelRouter — io.Writer который пишет в stdout + нужный HDFS writer по уровню
type levelRouter struct {
	stdout      io.Writer
	hdfsWriters map[string]*HDFSRotatingWriter
}

// Write парсит уровень из zerolog JSON и роутит в нужный HDFSRotatingWriter
func (r *levelRouter) Write(p []byte) (int, error) {
	// Пишем в stdout всегда
	_, _ = r.stdout.Write(p)

	// Определяем уровень из JSON: {"level":"info",...}
	level := extractLevel(p)
	if w, ok := r.hdfsWriters[level]; ok {
		if _, err := w.Write(p); err != nil {
			fmt.Fprintf(os.Stderr, "[hdfs] write error [%s]: %v\n", level, err)
		}
	}

	return len(p), nil
}

// extractLevel — быстрый парсинг уровня из zerolog JSON без reflect
func extractLevel(p []byte) string {
	// zerolog всегда пишет "level":"..." одним из первых полей
	const key = `"level":"`
	s := string(p)
	idx := strings.Index(s, key)
	if idx == -1 {
		return "info"
	}
	start := idx + len(key)
	end := strings.Index(s[start:], `"`)
	if end == -1 {
		return "info"
	}
	return s[start : start+end]
}

func newLogger(level string, writer io.Writer) *Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	skipFrameCount := 3
	zl := zerolog.New(writer).
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()

	return &Logger{logger: &zl}
}

// — методы логгера без изменений —

func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.msg("debug", message, args...)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.log(message, args...)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(message, args...)
}

func (l *Logger) Error(message interface{}, args ...interface{}) {
	if l.logger.GetLevel() == zerolog.DebugLevel {
		l.Debug(message, args...)
	}
	l.msg("error", message, args...)
}

func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.msg("fatal", message, args...)
	os.Exit(1)
}

func (l *Logger) log(message string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Info().Msg(message)
	} else {
		l.logger.Info().Msgf(message, args...)
	}
}

func (l *Logger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(msg.Error(), args...)
	case string:
		l.log(msg, args...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}
