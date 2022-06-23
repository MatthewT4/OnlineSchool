package DataBase

import (
	"OnlineSchool/internal/structs"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type WebinarDB struct {
	collection *mongo.Collection
}
type IWebinarDB interface {
	GetWebinars(ctx context.Context, start time.Time, end time.Time, courseId int) ([]structs.Webinar, error)
}

func NewWebinarDB(db *mongo.Database) *WebinarDB {
	return &WebinarDB{collection: db.Collection(nameWebinarDB)}
}

//Return array webinars from start to end date
func (w *WebinarDB) GetWebinars(ctx context.Context, start time.Time, end time.Time, courseId int) ([]structs.Webinar, error) {
	filter := bson.D{
		{"meet_date", bson.M{
			"$gte": start,
			"$lte": end,
		}},
		{"course_id", courseId},
	}
	var mas []structs.Webinar
	cur, err := w.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var elem structs.Webinar
		err = cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		mas = append(mas, elem)
	}
	err = cur.Err()
	if err != nil {
		return nil, err
	}
	cur.Close(context.TODO())
	return mas, err
}
