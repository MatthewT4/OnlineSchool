package DataBase

import (
	"OnlineSchool/internal/structs"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type SaveHomeworkDB struct {
	collection *mongo.Collection
}
type ISaveHomeworkDB interface {
	GetNextSaveHomeworks(ctx context.Context, courseId int, userId int) ([]structs.HomeworkSave, error)
	GetHomework(ctx context.Context, userId int, courseId int, homeworkId int) (structs.HomeworkSave, error)
}

func NewSaveHomeworkDB(db *mongo.Database) *SaveHomeworkDB {
	return &SaveHomeworkDB{collection: db.Collection(nameSaveHomeworkDB)}
}

func (h *SaveHomeworkDB) GetHomework(ctx context.Context, userId int, courseId int, homeworkId int) (structs.HomeworkSave, error) {
	filter := bson.M{"owner_id": userId, "course_id": courseId, "homework_id": homeworkId}
	var hw structs.HomeworkSave
	err := h.collection.FindOne(ctx, filter).Decode(&hw)
	return hw, err
}

func (t *SaveHomeworkDB) GetNextSaveHomeworks(ctx context.Context, courseId int, userId int) ([]structs.HomeworkSave, error) {
	filter := bson.M{"course_id": courseId,
		"owner_id": userId,
		"handed":   true,
		"public_date": bson.M{
			"$lte": time.Now(),
		}, "deadline": bson.M{
			"$gte": time.Now(),
		}}
	var mas []structs.HomeworkSave
	cursor, err := t.collection.Find(ctx, filter)
	if err != nil {
		return mas, err
	}
	for cursor.Next(context.TODO()) {
		var elem structs.HomeworkSave
		err = cursor.Decode(&elem)
		if err != nil {
			return mas, err
		}
		mas = append(mas, elem)
	}
	err = cursor.Err()
	if err != nil {
		return mas, err
	}
	cursor.Close(context.TODO())
	return mas, err
}
