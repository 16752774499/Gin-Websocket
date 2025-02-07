package model

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func Database(connString string) {
	dsn := connString
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: logger.New(
		//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		//	logger.Config{
		//		SlowThreshold:             200,         // 慢查询阈值
		//		LogLevel:                  logger.Info, // 日志级别，设置为Info将打印详细信息
		//		IgnoreRecordNotFoundError: true,        // 忽略记录未找到的错误
		//		Colorful:                  true,        // 是否彩色打印
		//	},
		//),
	})
	if err != nil {
		logrus.Info("mysql ini load failed': %v", err)
		panic(err)
	}
	logrus.Info("Mysql connect success")
	Migration()
}
