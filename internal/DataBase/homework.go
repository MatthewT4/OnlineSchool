package DataBase

import (
	"OnlineSchool/internal/structs"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HomeworkDB struct {
	collection *mongo.Collection
}
type IHomeworkDB interface {
	GetHomework(ctx context.Context, userId int, courseId int, homeworkId int) (structs.Homework, error)
}

func NewHomeworkDB(db *mongo.Database) *HomeworkDB {
	return &HomeworkDB{collection: db.Collection(nameHomeworkDB)}
}

func (h *HomeworkDB) GetHomework(ctx context.Context, userId int, courseId int, homeworkId int) (structs.Homework, error) {
	filter := bson.M{"owner_id": userId, "course_id": courseId, "homework_id": homeworkId}
	var hw structs.Homework
	err := h.collection.FindOne(ctx, filter).Decode(&hw)
	return hw, err
}
