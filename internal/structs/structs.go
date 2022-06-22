package structs

import "time"

type Course struct {
	CourseId  int       `bson:"course_id"`
	BuyPeriod int       `bson:"buy_period"`
	PeriodEnd time.Time `bson:"period_end"`
}
type User struct {
	UserId     int      `bson:"user_id"`
	VkId       string   `bson:"vk_id"`
	BuyCourses []Course `bson:"buy_courses"`
}
