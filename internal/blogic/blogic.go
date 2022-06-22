package blogic

import (
	"OnlineSchool/internal/DataBase"
	"go.mongodb.org/mongo-driver/mongo"
)

type BLogic struct {
	DBUser   DataBase.IUserDB
	DBCourse DataBase.ICourseDB
}

func NewBLogic(db *mongo.Database) *BLogic {
	return &BLogic{DBUser: DataBase.NewUserDB(db), DBCourse: DataBase.NewCourseDB(db)}
}

type IBLogic interface {
	GetUserCourses(user_id int) (int, string)
}
