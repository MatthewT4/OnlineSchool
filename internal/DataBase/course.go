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
	GetAvailableCourses(ctx context.Context, typ string, removeCoursesId []int) ([]structs.Course, error)
	GetIntensive(ctx context.Context, courseTeg string) (structs.Course, error)
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

func (c *CourseDB) GetAvailableCourses(ctx context.Context, typ string, removeCoursesId []int) ([]structs.Course, error) {
	filter := bson.M{
		"type":                   typ,
		"available_registration": true,
		"course_id": bson.M{
			"$nin": removeCoursesId,
		},
	}
	if len(removeCoursesId) == 0 {
		filter = bson.M{
			"type":                   typ,
			"available_registration": true,
		}
	}
	var ret []structs.Course
	cursor, err := c.collection.Find(ctx, filter)
	if err != nil {
		return ret, err
	}

	for cursor.Next(context.TODO()) {
		var elem structs.Course
		err = cursor.Decode(&elem)
		if err != nil {
			return ret, err
		}
		ret = append(ret, elem)
	}
	err = cursor.Err()
	if err != nil {
		return ret, err
	}
	cursor.Close(context.TODO())
	return ret, nil
}

func (c *CourseDB) GetIntensive(ctx context.Context, courseTeg string) (structs.Course, error) {
	filter := bson.M{"teg_course": courseTeg}
	var course structs.Course
	err := c.collection.FindOne(ctx, filter).Decode(&course)
	if err != nil {
		fmt.Println(err.Error())
	}
	return course, err
}
