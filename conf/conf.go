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
