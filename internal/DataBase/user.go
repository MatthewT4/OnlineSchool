package DataBase

import (
	"OnlineSchool/internal/structs"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDB struct {
	collection *mongo.Collection
}
type IUserDB interface {
	GetCourses(ctx context.Context, userId int64) ([]structs.UserCourse, error)
	GetUser(ctx context.Context, VKUserId int64) (structs.User, error)
	CreateUser(ctx context.Context, user structs.User) error
}

func NewUserDB(db *mongo.Database) *UserDB {
	return &UserDB{collection: db.Collection(nameUserDB)}
}

//GetCourses returns all courses owned by the user (active and inactive)
func (u *UserDB) GetCourses(ctx context.Context, userId int64) ([]structs.UserCourse, error) {
	filter := bson.M{"user_id": userId}
	var courses struct {
		BuyCourses []structs.UserCourse `bson:"buy_courses"`
	}
	err := u.collection.FindOne(ctx, filter).Decode(&courses)
	return courses.BuyCourses, err
}

func (u *UserDB) GetUser(ctx context.Context, VKUserId int64) (structs.User, error) {
	filter := bson.M{"vk_id": VKUserId}
	var user structs.User
	err := u.collection.FindOne(ctx, filter).Decode(&user)
	return user, err
}

func (u *UserDB) CreateUser(ctx context.Context, user structs.User) error {
	_, err := u.collection.InsertOne(ctx, user)
	return err
}
