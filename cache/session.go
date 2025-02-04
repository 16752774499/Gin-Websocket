package cache

import (
	"github.com/gin-contrib/sessions/redis"
	"github.com/sirupsen/logrus"
)

func NewSessionStore() redis.Store {
	store, err := redis.NewStoreWithDB(10, "tcp", RedisAddr, RedisPw, SessionDbName, []byte("xaiohua"))
	if err != nil {
		logrus.Info("NewSessionStore err", err.Error())
		panic(err)
	}
	return store

}
