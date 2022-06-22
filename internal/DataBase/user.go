package DataBase

import (
	"OnlineSchool/internal/structs"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const nameUserDB = "users"

type UserDB struct {
	collection *mongo.Collection
}
type IUserDB interface {
	GetCourses(ctx context.Context, userId int) ([]structs.Course, error)
}

func NewUserDB(db *mongo.Database) *UserDB {
	return &UserDB{collection: db.Collection(nameUserDB)}
}

func (u *UserDB) GetCourses(ctx context.Context, userId int) ([]structs.Course, error) {
	filter := bson.M{"user_id": userId}
	var courses struct {
		BuyCourses []structs.Course `bson:"buy_courses"`
	}
	err := u.collection.FindOne(ctx, filter).Decode(&courses)
	return courses.BuyCourses, err
}
