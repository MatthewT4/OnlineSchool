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
	Name         string
	MeetDate     time.Time
	WebinarId    int
	CourseId     int
	SpeakerId    int
	WebLink      string
	Recordlink   string
	Conspect     string
	Presentation string
}
