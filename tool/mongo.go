package tool

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDb struct {
	connection *mongo.Collection
}

func ConnectToDB(mongourl string) (*mongo.Collection, error) {
	url := mongourl
	name := "Battery"
	collection := "logAspect"
	maxCollection := 10
	var timeout time.Duration = 10

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	o := options.Client().ApplyURI(url)
	o.SetMaxPoolSize(uint64(maxCollection))
	client, err := mongo.Connect(ctx, o)
	if err != nil {
		return nil, err
	}
	return client.Database(name).Collection(collection), nil
}

func (m *MongoDb) jsonStr2Bson(str string) (interface{}, error) {
	// log.Println([]byte(str))
	var want interface{}
	err := bson.UnmarshalExtJSON([]byte(str), true, &want)
	if err != nil {
		return nil, err
	}
	log.Println(want)
	return want, nil
}

func (m *MongoDb) InsertToDb(wantStr string) (string, error) {
	if wantStr == "" {
		return "", errors.New("转换的字符串为空")
	}
	want, err := m.jsonStr2Bson(wantStr)
	if err != nil {
		return "", err
	}
	res, err := m.connection.InsertOne(context.TODO(), want)
	if err != nil {
		return "", err
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("断言错误")
	}
	return id.Hex(), nil
}

func NewMongoDbPool(mongourl string) (*MongoDb, error) { // 创建实例化
	pool, err := ConnectToDB(mongourl)
	if err != nil {
		return nil, err
	}
	return &MongoDb{
		connection: pool,
	}, nil
}
