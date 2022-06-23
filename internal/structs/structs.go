package structs

import "time"

type UserCourse struct {
	CourseId  int   `bson:"course_id" json:"course_id"`
	BuyPeriod []int `bson:"buy_period" json:"buy_period"`
}
type User struct {
	UserId     int          `bson:"user_id"`
	VkId       string       `bson:"vk_id"`
	BuyCourses []UserCourse `bson:"buy_courses"`
}
type Course struct {
	CourseId      int               `bson:"course_id"`
	NameCourse    string            `bson:"name_course"`
	PaymentPeriod map[int]time.Time `bson:"payment_period"`
	Сontacts      map[string]string `bson:"сontacts"`
}
type Webinar struct {
	Name         string    `bson:"name"`
	MeetDate     time.Time `bson:"meet_date"`
	WebinarId    int       `bson:"webinar_id"`
	CourseId     int       `bson:"course_id"`
	SpeakerId    int       `bson:"speaker_id"`
	WebLink      string    `bson:"web_link"`
	RecordLink   string    `bson:"record_link"`
	Conspect     string    `bson:"conspect"`
	Presentation string    `bson:"presentation"`
	Live         bool      `bson:"live"`
}
