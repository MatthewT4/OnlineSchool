package blogic

import (
	"OnlineSchool/internal/DataBase"
	"OnlineSchool/internal/structs"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type BLogic struct {
	DBUser             DataBase.IUserDB
	DBCourse           DataBase.ICourseDB
	DBWebinar          DataBase.IWebinarDB
	DBSaveHomework     DataBase.ISaveHomeworkDB
	DBTempHomework     DataBase.ITempHomeworkDB
	DBTaskBank         DataBase.ITaskBankDB
	DBPayment          DataBase.IPaymentDB
	DBPromoCode        DataBase.IPromoСodeDB
	DBAppliedPromoCode DataBase.IAppliedPromoСodeDB
	JWTManager         TokenManager
}

func NewBLogic(db *mongo.Database) *BLogic {
	return &BLogic{DBUser: DataBase.NewUserDB(db),
		DBCourse:           DataBase.NewCourseDB(db),
		DBWebinar:          DataBase.NewWebinarDB(db),
		DBSaveHomework:     DataBase.NewSaveHomeworkDB(db),
		DBTempHomework:     DataBase.NewTempHomeworkDB(db),
		DBTaskBank:         DataBase.NewTaskBankDB(db),
		DBPayment:          DataBase.NewPaymentDBDB(db),
		DBPromoCode:        DataBase.NewPromoСodeDB(db),
		DBAppliedPromoCode: DataBase.NewAppliedPromoСodeDB(db),
		JWTManager:         NewManager("dffid324jnk3"),
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
	GetIntensive(tagIntensive string, userId int64) (int, []byte)
	AddUserIntensive(tagIntensive string, userId int64) (int, string)

	GetActivePaymentsPeriod(userId int64) (int, []byte)
	CreatePayment(buy []structs.PayCourseType, userId int64, promoCodes string) (int, []byte, http.Cookie)
	LinkingPaymentToUser(userId int64, paymentId string) (int, string)
	CheckConnectingCourseGroups(userID int64) (int, []byte)
	GetInvitationLinkVkGroup(userId int64, courseId int) (int, []byte)
	CheckPayment(CPPayment structs.CloudPaymentReq, data []byte, contextHmac string) []byte
	RegisterApprovedPayment(CPPayment structs.CloudPaymentReq, data []byte, contextHmac string) []byte
	CheckAmountPromoCodes(userId int64, amount float64, promoCode string) (int, []byte)
}
