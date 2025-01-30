package conf

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strconv"
)

// Level 日志级别。建议从服务配置读取。
type LogConf struct {
	Dir     string `ini:"dir"`
	Name    string `ini:"name"`
	Level   string `ini:"level"`
	MaxSize int    `ini:"max_size"`
}

// InitLogger 初始化 logrus logger.
func InitLogger(logConf *LogConf) error {
	// 设置日志格式。
	formatter := &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		ForceColors:     true, // 强制开启颜色
		FullTimestamp:   true,
	}
	logrus.SetFormatter(formatter)

	switch logConf.Level {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	}
	logrus.SetReportCaller(true) // 打印文件、行号和主调函数。

	// 实现日志滚动。
	// Refer to https://www.cnblogs.com/jssyjam/p/11845475.html.
	logger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%v/%v", logConf.Dir, logConf.Name), // 日志输出文件路径。
		MaxSize:    logConf.MaxSize,                                 // 日志文件最大 size(MB)，缺省 100MB。
		MaxBackups: 10,                                              // 最大过期日志保留的个数。
		MaxAge:     30,                                              // 保留过期文件的最大时间间隔，单位是天。
		LocalTime:  true,                                            // 是否使用本地时间来命名备份的日志。
	}
	// 同时输出到标准输出与文件。
	logrus.SetOutput(io.MultiWriter(logger, os.Stdout))
	return nil
}
func LoadLog(file *ini.File) {
	MaxSizeStr := file.Section("LogConf").Key("max_size").String()
	MaxSize, _ := strconv.ParseUint(MaxSizeStr, 64, 10)
	logConf := &LogConf{
		Dir:     file.Section("LogConf").Key("Dir").String(),
		Name:    file.Section("LogConf").Key("Name").String(),
		Level:   file.Section("LogConf").Key("Level").String(),
		MaxSize: int(MaxSize),
	}
	InitLogger(logConf)
}
