package structs

import "time"

type UserCourse struct {
	CourseId    int   `bson:"course_id" json:"course_id"`
	BuyPeriod   []int `bson:"buy_period" json:"buy_period"`
	HeartCount  int   `bson:"heart_count"`
	FreezingDay int   `bson:"freezing_day"`
	Freeze      bool  `bson:"freeze"`
	Active      bool  `bson:"active"`
}
type User struct {
	UserId     int64        `bson:"user_id"`
	VkId       int64        `bson:"vk_id"`
	Avatar     string       `bson:"avatar"`
	BuyCourses []UserCourse `bson:"buy_courses"`
	FirstName  string       `bson:"first_name"`
	LastName   string       `bson:"last_name"`
}

type PayPeriod struct {
	PeriodId  int       `bson:"period_id"`
	StartDate time.Time `bson:"start_date"`
	EndDate   time.Time `bson:"end_date"`
	Price     int       `bson:"price"`
}
type Course struct {
	CourseId       int         `bson:"course_id"`
	NameCourse     string      `bson:"name_course"`
	PaymentPeriods []PayPeriod `bson:"payment_periods"`
	Teacher        string      `bson:"teacher"`
	VkChat         string      `bson:"vk_chat"`
	VkGroup        string      `bson:"vk_group"`
}
type Webinar struct {
	Name         string    `bson:"name"`
	MeetDate     time.Time `bson:"meet_date"`
	WebinarId    int       `bson:"webinar_id"`
	CourseId     int       `bson:"course_id"`
	SpeakerId    int64     `bson:"speaker_id"`
	WebLink      string    `bson:"web_link"`
	RecordLink   string    `bson:"record_link"`
	Conspect     string    `bson:"conspect"`
	Presentation string    `bson:"presentation"`
	Live         bool      `bson:"live"`
}

type Task struct {
	TaskId      int      `bson:"task_id"`
	CourseName  string   `bson:"course_name"`
	Text        string   `bson:"text"`
	File        []string `bson:"file,omitempty"`
	Answers     []string `bson:"answers"`
	Solution    string   `bson:"solution,omitempty"`
	Written     bool     `bson:"written"`
	TypeAnswers []string `bson:"type_answers,omitempty"`
	MaxPoint    int      `bson:"max_point"`
	Handler     string   `bson:"handler,omitempty"`
}

type HomeworkTask struct {
	Number     int    `bson:"number" json:"number"`
	TaskId     int    `bson:"task_id" json:"task_id"`
	UserAnswer string `bson:"user_answer" json:"user_answer"`
	MaxPoint   int    `bson:"max_point" json:"max_point"`
	Point      int    `bson:"point" json:"point,omitempty"`
}

type HomeworkSave struct {
	OwnerId    int64          `bson:"owner_id"`
	HomeworkId int            `bson:"homework_id"`
	Tasks      []HomeworkTask `bson:"tasks"`
	Result     int            `bson:"result,omitempty"`
	MaxPoints  int            `bson:"max_points"`
	Delivered  time.Time      `bson:"delivered,omitempty"`
	Handed     bool           `bson:"handed"`
}

type HomeworkTemplate struct {
	HomeworkName string         `bson:"homework_name"`
	PublicDate   time.Time      `bson:"public_date"`
	Deadline     time.Time      `bson:"deadline"`
	CourseId     int            `bson:"course_id"`
	HomeworkId   int            `bson:"homework_id"`
	Tasks        []HomeworkTask `bson:"tasks"`
	MaxPoints    int            `bson:"max_points"`
}
