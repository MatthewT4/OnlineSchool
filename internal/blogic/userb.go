package blogic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type retWebinar struct {
	Name         string    `json:"name"`
	MeetDate     time.Time `json:"meet_date"`
	WebinarId    int       `json:"webinar_id"`
	WebLink      string    `json:"web_link,omitempty"`
	RecordLink   string    `json:"record_link,omitempty"`
	Conspect     string    `json:"conspect,omitempty"`
	Presentation string    `json:"presentation,omitempty"`
}

func (b *BLogic) checkUserCourse(user_id int, course_id int) bool {
	res, err := b.DBUser.GetCourses(context.TODO(), user_id)
	if err != nil {
		return false
	}
	for i := 0; i < len(res); i++ {
		if course_id == res[i].CourseId {
			return true
		}
	}
	return false
}

func (b *BLogic) GetUserCourses(user_id int) (int, string) {
	res, err := b.DBUser.GetCourses(context.TODO(), user_id)
	if err != nil {
		fmt.Println(err.Error())
		return 404, "not found"
	}
	type resCourses struct {
		NameCourse string    `json:"name_course"`
		PaymentEnd time.Time `json:"payment_end"`
	}

	var mas []resCourses
	for i := 0; i < len(res); i++ {
		course, er := b.DBCourse.GetCourse(context.TODO(), res[i].CourseId)
		if er == nil {
			//find max payment period
			max := 0
			for j := 0; j < len(res[i].BuyPeriod); j++ {
				if max < res[i].BuyPeriod[j] {
					max = res[i].BuyPeriod[j]
				}
			}
			var c resCourses
			c.NameCourse = course.NameCourse
			c.PaymentEnd = course.PaymentPeriod[max]
			fmt.Println(c)
			mas = append(mas, c)
		} else {
			return 404, "not found"
			fmt.Println(err.Error())
		}
	}
	if len(mas) == 0 {
		return 404, "not found"
	}
	jr, err := json.Marshal(mas)
	if err != nil {
		return 404, "not found"
		log.Fatal(err)
	}
	return 200, string(jr)
}

func (b *BLogic) GetNextWebinars(userId int, courseId int) (int, string) {
	if !b.checkUserCourse(userId, courseId) {
		return 404, "not found"
	}
	start_time := time.Now()
	r := start_time.Format("2006-01-02")
	var er error
	start_time, er = time.Parse("2006-01-02", r)
	if er != nil {
		return 404, "time parse error"
	}
	end_time := time.Now().Add(time.Hour * 24 * 365 * 2)
	res, err := b.getWebinars(start_time, end_time, courseId)
	if err != nil {
		fmt.Println(err)
		return 404, "not found"
	}
	re, erro := json.Marshal(&res)
	if erro != nil {
		fmt.Println(erro)
		return 500, "json marshal fail"
	}
	return 200, string(re)
}

func (b *BLogic) GetPastWebinars(userId int, courseId int) (int, string) {
	if !b.checkUserCourse(userId, courseId) {
		return 404, "not found"
	}
	endTime := time.Now()
	startTime := time.Now().Add(-time.Hour * 24 * 365 * 2)
	res, err := b.getWebinars(startTime, endTime, courseId)
	if err != nil {
		fmt.Println(err)
		return 404, "not found"
	}
	re, erro := json.Marshal(&res)
	if erro != nil {
		fmt.Println(erro)
		return 500, "json marshal fail"
	}
	return 200, string(re)
}

func (b *BLogic) GetTodayWebinars(userId int) (int, string) {
	res, err := b.DBUser.GetCourses(context.TODO(), userId)
	if err != nil {
		fmt.Println(err.Error())
		return 404, "not found"
	}
	startTime := time.Now()
	r := startTime.Format("2006-01-02")
	var er error
	startTime, er = time.Parse("2006-01-02", r)
	endTime := startTime.Add(24 * time.Hour)
	if er != nil {
		return 404, "time parse error"
	}
	var mas []retWebinar
	for i := 0; i < len(res); i++ {
		re, er := b.getWebinars(startTime, endTime, res[i].CourseId)
		if er != nil {
			fmt.Println(er.Error())
			return 404, "not found"
		}
		mas = append(mas, re...)
	}
	re, erro := json.Marshal(&mas)
	if err != nil {
		fmt.Println(erro)
		return 404, "not found"
	}
	return 200, string(re)
}

func (b *BLogic) getWebinars(start_time time.Time, end_time time.Time, courseId int) ([]retWebinar, error) {
	res, err := b.DBWebinar.GetWebinars(context.TODO(), start_time, end_time, courseId)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var mas []retWebinar
	for i := 0; i < len(res); i++ {
		var webinar retWebinar
		webinar.Name = res[i].Name
		webinar.MeetDate = res[i].MeetDate.Local()
		webinar.WebinarId = res[i].WebinarId
		webinar.RecordLink = res[i].RecordLink
		webinar.Conspect = res[i].Conspect
		webinar.Presentation = res[i].Presentation
		if res[i].Live {
			webinar.WebLink = res[i].WebLink
		}
		mas = append(mas, webinar)
	}
	return mas, nil
}
