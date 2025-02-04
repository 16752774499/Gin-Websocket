package conf

import (
	"Gin-WebSocket/model"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	_ "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/ini.v1"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strconv"
	"strings"
)

// 读取conf.ini
var (
	MongoDBClient *mongo.Client
	AppMode       string
	HttpPort      string
	Db            string
	DbHost        string
	DbPort        string
	DbUser        string
	DbPassWord    string
	DbName        string
	RedisDb       string
	RedisAddr     string
	RedisPw       string
	RedisDbName   string
	MongoDBName   string
	MongoDBAddr   string
	MongoDBUser   string
	MongoDBPwd    string
	MongoDBPort   string
)

func Init() {
	//从本地读取环境

	file, err := ini.Load("./conf/conf.ini")
	if err != nil {
		fmt.Printf("ini load failed': %v", err)
	}
	LoadLog(file)
	wscInit(file)
	LoadServer(file)
	LoadMySql(file)
	LoadMongoDB(file)
	MongoDB() //链接MongoDB
	//mysqlPath := "gorm:ECweAtSJPaSBffd3@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	// 正确的拼接方式
	mysqlPath := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ")/", DbName, "?charset=utf8mb4&parseTime=True&loc=Local"}, "")
	model.Database(mysqlPath) //链接Mysql

}

func MongoDB() {
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s", MongoDBUser, MongoDBPwd, MongoDBAddr, MongoDBPort)
	clientOptions := options.Client().ApplyURI(connectionString)
	var err error
	MongoDBClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.Info(err)
		panic(err)
	}
	logrus.Info("MongoDB connect success")
}

func LoadMongoDB(file *ini.File) {

	MongoDBName = file.Section("MongoDB").Key("MongoDBName").String()
	MongoDBAddr = file.Section("MongoDB").Key("MongoDBAddr").String()
	MongoDBPwd = file.Section("MongoDB").Key("MongoDBPwd").String()
	MongoDBPort = file.Section("MongoDB").Key("MongoDBPort").String()
	MongoDBUser = file.Section("MongoDB").Key("MongoDBUser").String()

}

func LoadMySql(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}

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
	err := InitLogger(logConf)
	if err != nil {
		return
	}
}
