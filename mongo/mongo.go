package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"main/pkg/setting"
	"time"
)

var (
	MongoClient *mongo.Client
	MongoDB *mongo.Database
)

// Setup 初始化数据库
func Setup() {
	uri := fmt.Sprintf("mongodb://%s:%d", setting.MongoSetting.Host, setting.MongoSetting.Port)
	var err error
	// 建立连接
	MongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetConnectTimeout(5*time.Second))
	if err != nil {
		log.Fatalln("mongoDB建立连接失败：" + err.Error())
	}
	// 选择数据库
	MongoDB = MongoClient.Database(setting.MongoSetting.Name)
	err = MongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalln(err)
	}
}