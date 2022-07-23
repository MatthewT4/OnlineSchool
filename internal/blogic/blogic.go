package blogic

import (
	"OnlineSchool/internal/DataBase"
	"OnlineSchool/internal/structs"
	"go.mongodb.org/mongo-driver/mongo"
)

type BLogic struct {
	DBUser         DataBase.IUserDB
	DBCourse       DataBase.ICourseDB
	DBWebinar      DataBase.IWebinarDB
	DBSaveHomework DataBase.ISaveHomeworkDB
	DBTempHomework DataBase.ITempHomeworkDB
	DBTaskBank     DataBase.ITaskBankDB
	JWTManager     TokenManager
}

func NewBLogic(db *mongo.Database) *BLogic {
	return &BLogic{DBUser: DataBase.NewUserDB(db),
		DBCourse:       DataBase.NewCourseDB(db),
		DBWebinar:      DataBase.NewWebinarDB(db),
		DBSaveHomework: DataBase.NewSaveHomeworkDB(db),
		DBTempHomework: DataBase.NewTempHomeworkDB(db),
		DBTaskBank:     DataBase.NewTaskBankDB(db),
		JWTManager:     NewManager("dffid324jnk3"),
	}
}

type IBLogic interface {
	GetUserCourses(user_id int64) (int, string)
	GetNextWebinars(user_id int64, course_id int) (int, string)
	GetPastWebinars(user_id int64, course_id int) (int, string)
	GetTodayWebinars(user_id int64) (int, string)
	GetHomework(userId int64, homeworkId int) (int, []byte)
	GetNextCourseHomeworks(userId int64, courseId int) (int, []byte)
	GetNextHomeworks(userId int64) (int, []byte)
	GetInfoCourse(userId int64, courseId int) (int, []byte)
	GetPastCourseHomeworks(userId int64, courseId int) (int, []byte)
	Login(VKCode string, redirectUrl string) (int, []byte, string /*cookie*/)
	Authentication(token string) (int64, int, error)
	SubmitHomework(userId int64, homeworkId int, answers []structs.HomeworkTask) (int, string)
	GetActivePaymentsPeriod(userId int64) (int, []byte)
}
