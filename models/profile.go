package models

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"tetris/pkg/z"
	"tetris/pkg/zmongo"
)

type Counter struct {
	Id    string `bson:"_id"`
	Count int64  `bson:"count"`
}

type Profile struct {
	Name      string `bson:"name"`
	Coin      int32  `bson:"coin"`
	UserId    int64  `bson:"_id"`
	Pic       string `bson:"pic"`
	UpdatedAt int64  `bson:"updated_at"`
}

func GetPlayer(userId int64, fields ...string) (*Profile, error) {
	filter := bson.M{
		"_id": userId,
	}
	opts := &options.FindOneOptions{}
	projection := bson.M{}
	for _, field := range fields {
		projection[field] = 1
	}
	opts.Projection = projection
	p := &Profile{}
	err := mclient.FindOne(p, DB_NAME, COL_PROFILE, filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, z.NilError{Msg: fmt.Sprintf("%d", userId)}
		}
		return nil, err
	}
	return p, nil
}

func NewPlayer(name string, coin int32) *Profile {
	return &Profile{
		Name: name,
		Coin: coin,
	}
}

func UpdatePlayer(userId int64, fieldMap map[string]interface{}) error {
	filter := bson.M{
		"_id": userId,
	}
	fieldMap["updated_at"] = z.NowUnixMilli()
	err := mclient.UpsertOne(DB_NAME, COL_PROFILE, filter, fieldMap)
	return err
}

func CreatePlayer(player *Profile) (int64, error) {
	filter := bson.M{
		"_id": "_id",
	}
	after := options.After
	opts := &options.FindOneAndUpdateOptions{
		Upsert:         zmongo.NewTrue(),
		ReturnDocument: &after,
	}
	update := bson.M{
		"$inc": bson.M{
			"count": 1,
		},
	}
	counter := &Counter{}
	err := mclient.FindOneAndUpdate(DB_NAME, COL_COUNTER, filter, update, counter, opts)
	if err != nil {
		return 0, err
	}
	uid := counter.Count
	player.UserId = uid
	_, err = mclient.InsertOne(DB_NAME, COL_PROFILE, player)
	return uid, err
}
