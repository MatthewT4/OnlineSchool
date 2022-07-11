package DataBase

import (
	"OnlineSchool/internal/structs"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type TempHomeworkDB struct {
	collection *mongo.Collection
}
type ITempHomeworkDB interface {
	GetHomework(ctx context.Context, homeworkId int) (structs.HomeworkTemplate, error)
	GetNextTempHomeworks(ctx context.Context, courseId int) ([]structs.HomeworkTemplate, error)
	GetPastTempHomeworks(ctx context.Context, courseId int) ([]structs.HomeworkTemplate, error)
}

func NewTempHomeworkDB(db *mongo.Database) *TempHomeworkDB {
	return &TempHomeworkDB{collection: db.Collection(nameTempHomeworkDB)}
}

func (t *TempHomeworkDB) GetHomework(ctx context.Context, homeworkId int) (structs.HomeworkTemplate, error) {
	filter := bson.M{"homework_id": homeworkId}
	var hw structs.HomeworkTemplate
	err := t.collection.FindOne(ctx, filter).Decode(&hw)
	return hw, err
}

func (t *TempHomeworkDB) GetNextTempHomeworks(ctx context.Context, courseId int) ([]structs.HomeworkTemplate, error) {
	filter := bson.M{"course_id": courseId,
		"public": true,
		"public_date": bson.M{
			"$lte": time.Now()},
		"deadline": bson.M{
			"$gte": time.Now()},
	}
	var mas []structs.HomeworkTemplate
	cursor, err := t.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var elem structs.HomeworkTemplate
		err = cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}
		mas = append(mas, elem)
	}
	err = cursor.Err()
	if err != nil {
		return nil, err
	}
	cursor.Close(context.TODO())
	return mas, err
}
func (t *TempHomeworkDB) GetPastTempHomeworks(ctx context.Context, courseId int) ([]structs.HomeworkTemplate, error) {
	filter := bson.M{"course_id": courseId,
		"public": true,
		"public_date": bson.M{
			"$lte": time.Now()},
	}
	var mas []structs.HomeworkTemplate
	cursor, err := t.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var elem structs.HomeworkTemplate
		err = cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}
		mas = append(mas, elem)
	}
	err = cursor.Err()
	if err != nil {
		return nil, err
	}
	cursor.Close(context.TODO())
	return mas, err
}
