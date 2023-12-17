package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Logger *zap.Logger
	Sugar  zap.SugaredLogger
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func InitLogger() (*Logger, error) {
	config := zap.Config{
		Encoding: "console",
		Level:    zap.NewAtomicLevelAt(zap.DebugLevel),
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			TimeKey:     "timestamp",
			EncodeTime:  zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	mylogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	defer mylogger.Sync()
	sugar := *mylogger.Sugar()

	return &Logger{
		Logger: mylogger,
		Sugar:  sugar,
	}, nil
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func (l *Logger) MyLogger(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		respData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   respData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		l.Sugar.Infoln(
			"uri:", r.RequestURI,
			"method:", r.Method,
			"duration:", duration,
			"status:", respData.status,
			"size:", respData.size,
		)

	}
	return http.HandlerFunc(logFn)
}

func (l *Logger) LogInfo(args ...interface{}) {
	l.Sugar.Info(args)
}

func (l *Logger) LogWarning(args ...interface{}) {
	l.Sugar.Warn(args)
}
