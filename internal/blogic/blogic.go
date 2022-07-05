package blogic

import (
	"OnlineSchool/internal/DataBase"
	"go.mongodb.org/mongo-driver/mongo"
)

type BLogic struct {
	DBUser         DataBase.IUserDB
	DBCourse       DataBase.ICourseDB
	DBWebinar      DataBase.IWebinarDB
	DBSaveHomework DataBase.ISaveHomeworkDB
	DBTempHomework DataBase.ITempHomeworkDB
}

func NewBLogic(db *mongo.Database) *BLogic {
	return &BLogic{DBUser: DataBase.NewUserDB(db), DBCourse: DataBase.NewCourseDB(db), DBWebinar: DataBase.NewWebinarDB(db), DBSaveHomework: DataBase.NewSaveHomeworkDB(db), DBTempHomework: DataBase.NewTempHomeworkDB(db)}
}

type IBLogic interface {
	GetUserCourses(user_id int) (int, string)
	GetNextWebinars(user_id int, course_id int) (int, string)
	GetPastWebinars(user_id int, course_id int) (int, string)
	GetTodayWebinars(user_id int) (int, string)
	GetHomework(userId int, courseId int, homeworkId int) (int, []byte)
	GetNextCourseHomeworks(userId, courseId int) (int, []byte)
	GetNextHomeworks(userId int) (int, []byte)
	GetInfoCourse(userId int, courseId int) (int, []byte)
}
