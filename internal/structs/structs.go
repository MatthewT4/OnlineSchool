package structs

type Course struct {
	CourseId  int   `bson:"course_id" json:"course_id"`
	BuyPeriod []int `bson:"buy_period" json:"buy_period"`
}
type User struct {
	UserId     int      `bson:"user_id"`
	VkId       string   `bson:"vk_id"`
	BuyCourses []Course `bson:"buy_courses"`
}
