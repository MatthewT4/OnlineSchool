package DataBase

import (
	"OnlineSchool/internal/structs"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseDB struct {
	collection *mongo.Collection
}
type ICourseDB interface {
	GetCourse(ctx context.Context, CourseId int) (structs.Course, error)
}

func NewCourseDB(db *mongo.Database) *CourseDB {
	return &CourseDB{collection: db.Collection(nameCourseDB)}
}

func (c *CourseDB) GetCourse(ctx context.Context, CourseId int) (structs.Course, error) {
	filter := bson.M{"course_id": CourseId}
	var course structs.Course
	err := c.collection.FindOne(ctx, filter).Decode(&course)
	if err != nil {
		fmt.Println(err.Error())
	}
	return course, err
}
