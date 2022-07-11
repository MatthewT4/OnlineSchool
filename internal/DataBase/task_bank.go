package DataBase

import (
	"OnlineSchool/internal/structs"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskBankDB struct {
	collection *mongo.Collection
}
type ITaskBankDB interface {
	GetTasks(ctx context.Context, tasksId []int) ([]structs.Task, error)
}

func NewTaskBankDB(db *mongo.Database) *TaskBankDB {
	return &TaskBankDB{collection: db.Collection(nameTaskBankDB)}
}

func (t *TaskBankDB) GetTasks(ctx context.Context, tasksId []int) ([]structs.Task, error) {
	filter := bson.M{"task_id": bson.M{
		"$in": tasksId,
	}}
	var mas []structs.Task
	cursor, err := t.collection.Find(ctx, filter)
	if err != nil {
		return mas, err
	}
	for cursor.Next(context.TODO()) {
		var elem structs.Task
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
