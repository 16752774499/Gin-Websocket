package cache

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"strconv"
)

var (
	RedisClient *redis.Client
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
)

func init() {
	//从本地读取环境
	file, err := ini.Load("./conf/conf.ini") //加载配置文件
	if err != nil {
		logrus.Info("redis ini load failed': %v", err)
	}
	LoadMyRedis(file) //读取配置信息
	Redis()           //创建reids链接
}

func LoadMyRedis(file *ini.File) {

	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}
func Redis() {
	db, _ := strconv.ParseUint(RedisDb, 10, 64) // string to uint64
	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPw,
		DB:       int(db),
	})
	_, err := client.Ping().Result()
	if err != nil {
		logrus.Info(err)
		panic(err)
	}
	RedisClient = client

}
