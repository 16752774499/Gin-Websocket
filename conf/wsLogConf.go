package conf

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/ini.v1"
	"os"
	"strings"
	"time"
)

var Log *zap.Logger

// getEncoderConfig 获取日志编码器配置
func getEncoderConfig() zapcore.EncoderConfig {
	config := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 日志级别大写
		EncodeTime:     customTimeEncoder,              // 自定义时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, // 日志时间戳
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短格式的调用栈信息
	}
	return config
}

// customTimeEncoder 自定义时间编码器
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// getLogWriter 获取日志文件写入器
func getLogWriter(filePath string) (zapcore.WriteSyncer, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return zapcore.AddSync(file), nil
}

// NewCustomZapLogger 创建一个带有文件输出和颜色编码的自定义zap日志记录器
func NewCustomZapLogger(filePath string) (*zap.Logger, error) {
	// 获取日志编码器
	encoder := zapcore.NewConsoleEncoder(getEncoderConfig())

	// 获取日志文件写入器
	logWriter, err := getLogWriter(filePath)
	if err != nil {
		return nil, err
	}

	// 设置不同级别的日志输出
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// 创建核心日志记录器
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), highPriority), // 错误级别日志输出到stderr
		zapcore.NewCore(encoder, logWriter, lowPriority),                   // 其他级别日志输出到文件
	)

	// 创建最终的日志记录器
	logger := zap.New(core, zap.AddCaller())
	return logger, nil
}

func wscInit(file *ini.File) {
	//
	Dir := file.Section("LogConf").Key("Dir").String()
	fileName := file.Section("LogConf").Key("WebSocketName").String()

	logFilePath := strings.Join([]string{Dir, fileName}, "/")

	// 创建自定义日志记录器
	logger, err := NewCustomZapLogger(logFilePath)
	if err != nil {
		fmt.Printf("创建日志记录器失败: %v\n", err)
		return
	}
	Log = logger
	defer logger.Sync()

}
