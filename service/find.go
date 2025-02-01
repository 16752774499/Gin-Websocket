package service

import (
	"Gin-WebSocket/conf"
	"Gin-WebSocket/model/ws"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sort"
	"time"
)

type SendSortMsg struct {
	Context  string `json:"context"`
	Read     uint   `json:"read"`
	CreateAt int64  `json:"create_at"`
}

func InsertMsg(database, id string, content string, read uint, expire int64) error {
	// 插入到MongoDB
	// 获取指定数据库和集合的引用，如果不存在该集合，MongoDB会自动创建
	collection := conf.MongoDBClient.Database(database).Collection(id) // 没有id这个集合会创建

	// 创建一个ws.Trainer类型的实例，用于存储消息内容
	comment := ws.Trainer{
		Content:   content,                    // 消息内容
		StartTime: time.Now().Unix(),          // 消息开始时间，使用当前时间的Unix时间戳
		EndTime:   time.Now().Unix() + expire, // 消息结束时间，为当前时间加上过期时间
		Read:      read,                       // 消息阅读状态
	}

	// 将消息插入到MongoDB集合中
	_, err := collection.InsertOne(context.TODO(), comment)

	// 返回插入操作的结果错误
	return err
}
func FindMany(database string, sendID string, id string, time int64, pageSize int) (results []ws.Result, err error) {
	//var resultsMe []ws.Trainer
	//var resultsYou []ws.Trainer
	//sendIdCollection := conf.MongoDBClient.Database(database).Collection(sendID)
	//idCollection := conf.MongoDBClient.Database(database).Collection(id)
	//// 如果不知道该使用什么context，可以通过context.TODO() 产生context
	//sendIdTimeCursor, err := sendIdCollection.Find(context.TODO(),
	//	options.Find().SetSort(bson.D{{"startTime", 1}}), options.Find().SetLimit(int64(pageSize)))
	//idTimeCursor, err := idCollection.Find(context.TODO(),
	//	options.Find().SetSort(bson.D{{"startTime", 1}}), options.Find().SetLimit(int64(pageSize)))
	//err = sendIdTimeCursor.All(context.TODO(), &resultsYou) // sendId 对面发过来的
	//err = idTimeCursor.All(context.TODO(), &resultsMe)      // Id 发给对面的
	//results, _ = AppendAndSort(resultsMe, resultsYou)
	//return results, err
	var resultsMe []ws.Trainer
	var resultsYou []ws.Trainer

	sendIdCollection := conf.MongoDBClient.Database(database).Collection(sendID)
	idCollection := conf.MongoDBClient.Database(database).Collection(id)

	// 执行sendIdCollection的查询
	sendIdFindOptions := options.Find()
	sendIdFindOptions.SetSort(bson.D{{"startTime", -1}})
	sendIdFindOptions.SetLimit(int64(pageSize))
	sendIdTimeCursor, err := sendIdCollection.Find(context.TODO(), bson.D{}, sendIdFindOptions)
	if err != nil {
		return nil, fmt.Errorf("sendIdCollection.Find failed: %v", err)
	}
	defer sendIdTimeCursor.Close(context.TODO())
	err = sendIdTimeCursor.All(context.TODO(), &resultsYou)
	if err != nil {
		return nil, fmt.Errorf("sendIdTimeCursor.All failed: %v", err)
	}

	// 执行idCollection的查询
	idFindOptions := options.Find()
	idFindOptions.SetSort(bson.D{{"startTime", -1}})
	idFindOptions.SetLimit(int64(pageSize))
	idTimeCursor, err := idCollection.Find(context.TODO(), bson.D{}, idFindOptions)
	if err != nil {
		return nil, fmt.Errorf("idCollection.Find failed: %v", err)
	}
	defer idTimeCursor.Close(context.TODO())
	err = idTimeCursor.All(context.TODO(), &resultsMe)
	if err != nil {
		return nil, fmt.Errorf("idTimeCursor.All failed: %v", err)
	}

	results, _ = AppendAndSort(resultsMe, resultsYou)
	return results, nil
}

func AppendAndSort(resultMe []ws.Trainer, resultYou []ws.Trainer) (results []ws.Result, err error) {
	for _, r := range resultMe { //构造返回的Msg
		SendSort := SendSortMsg{
			Context:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{ //构造返回的内容，包括传送者
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", SendSort),
			From:      "me",
		}
		results = append(results, result)
	}
	for _, r := range resultYou { //构造返回的Msg
		SendSort := SendSortMsg{
			Context:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{ //构造返回的内容，包括传送者
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", SendSort),
			From:      "you",
		}
		results = append(results, result)
	}
	//根据时间排序
	sort.Slice(results, func(i, j int) bool { return results[i].StartTime < results[j].StartTime })
	return results, nil

}
