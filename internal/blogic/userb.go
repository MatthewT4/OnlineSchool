package blogic

import (
	"OnlineSchool/internal/structs"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
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

//return bool and int (number last payment period)
func (b *BLogic) checkUserCourse(courses []structs.UserCourse, courseId int) bool {
	for i := 0; i < len(courses); i++ {
		if courseId == courses[i].CourseId {
			return true
		}
	}
	return false
}

func (b *BLogic) GetUserCourses(userId int64) (int, string) {
	res, err := b.DBUser.GetCourses(context.TODO(), userId)
	if err != nil {
		fmt.Println(err.Error())
		return 404, "not found"
	}
	type resCourses struct {
		NameCourse string    `json:"name_course"`
		PaymentEnd time.Time `json:"payment_end"`
		CourseId   int       `json:"course_id"`
	}

	var mas []resCourses
	for i := 0; i < len(res); i++ {
		if res[i].Active == false {
			continue
		}
		course, er := b.DBCourse.GetCourse(context.TODO(), res[i].CourseId)
		if er == nil {
			var c resCourses
			c.NameCourse = course.NameCourse
			payEnd, e := b.getDateLastPaymentPeriod(res, res[i].CourseId, course)
			if e != nil {
				fmt.Println("Error parse end period getDateLastPaymentPeriod COURSE_ID=", res[i].CourseId)
				continue
			}
			c.PaymentEnd = payEnd
			c.CourseId = course.CourseId
			fmt.Println(c)
			mas = append(mas, c)
		} else {
			return 404, "not found"
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

func (b *BLogic) GetNextWebinars(userId int64, courseId int) (int, string) {
	res, err := b.DBUser.GetCourses(context.TODO(), userId)
	if err != nil {
		return 404, "not found"
	}
	if !b.checkUserCourse(res, courseId) {
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

	course, e := b.DBCourse.GetCourse(context.TODO(), courseId)
	if e != nil {
		return 500, "Server error"
	}
	dateLastPaymentPeriod, errorr := b.getDateLastPaymentPeriod(res, courseId, course)
	if errorr != nil {
		return 500, errorr.Error()
	}
	result, er := b.getWebinars(start_time, end_time, courseId, dateLastPaymentPeriod)
	if er != nil {
		fmt.Println(er.Error())
		return 404, "not found"
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].MeetDate.Before(result[j].MeetDate)
	})
	re, erro := json.Marshal(&result)
	if erro != nil {
		fmt.Println(erro)
		return 500, "json marshal fail"
	}
	return 200, string(re)
}

func (b *BLogic) GetPastWebinars(userId int64, courseId int) (int, string) {
	res, err := b.DBUser.GetCourses(context.TODO(), userId)
	if err != nil {
		return 404, "not found"
	}
	if !b.checkUserCourse(res, courseId) {
		return 404, "not found"
	}
	endTime := time.Now()
	startTime := time.Now().Add(-time.Hour * 24 * 365 * 2)

	course, e := b.DBCourse.GetCourse(context.TODO(), courseId)
	if e != nil {
		return 500, "Server error"
	}
	dateLastPaymentPeriod, errorr := b.getDateLastPaymentPeriod(res, courseId, course)
	if errorr != nil {
		return 500, errorr.Error()
	}
	result, err := b.getWebinars(startTime, endTime, courseId, dateLastPaymentPeriod)
	if err != nil {
		fmt.Println(err)
		return 404, "not found"
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].MeetDate.Before(result[j].MeetDate)
	})
	re, erro := json.Marshal(&result)
	if erro != nil {
		fmt.Println(erro)
		return 500, "json marshal fail"
	}
	return 200, string(re)
}

func (b *BLogic) GetTodayWebinars(userId int64) (int, string) {
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
		if !res[i].Active {
			continue
		}

		course, e := b.DBCourse.GetCourse(context.TODO(), res[i].CourseId)
		if e != nil {
			return 500, "Server error"
		}
		dateLastPaymentPeriod, errorr := b.getDateLastPaymentPeriod(res, res[i].CourseId, course)
		if errorr != nil {
			return 500, errorr.Error()
		}
		re, erro := b.getWebinars(startTime, endTime, res[i].CourseId, dateLastPaymentPeriod)
		if erro != nil {
			fmt.Println(erro.Error())
			return 404, "not found"
		}
		mas = append(mas, re...)
	}
	if len(mas) == 0 {
		return 404, "not found"
	}
	sort.SliceStable(mas, func(i, j int) bool {
		return mas[i].MeetDate.Before(mas[j].MeetDate)
	})
	re, erro := json.Marshal(&mas)
	if err != nil {
		fmt.Println(erro)
		return 404, "not found"
	}

	return 200, string(re)
}

func (b *BLogic) getWebinars(start_time time.Time, end_time time.Time, courseId int, endLastPaymentPeriod time.Time) ([]retWebinar, error) {
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
		if webinar.MeetDate.Before(endLastPaymentPeriod) { // not return if date webinar > date last payment period
			webinar.WebinarId = res[i].WebinarId
			webinar.RecordLink = res[i].RecordLink
			webinar.Conspect = res[i].Conspect
			webinar.Presentation = res[i].Presentation
		}
		if res[i].Live {
			webinar.WebLink = res[i].WebLink
		}
		mas = append(mas, webinar)
	}
	return mas, nil
}

func (b *BLogic) getDateLastPaymentPeriod(courses []structs.UserCourse, courseId int, course structs.Course) (time.Time, error) {
	maxx_period := 0
	for i := 0; i < len(courses); i++ {
		if courses[i].CourseId == courseId {
			for j := 0; j < len(courses[i].BuyPeriod); j++ {
				if maxx_period < courses[i].BuyPeriod[j] {
					maxx_period = courses[i].BuyPeriod[j]
				}
			}
		}
	}
	for _, val := range course.PaymentPeriods {
		if val.PeriodId == maxx_period {
			return val.EndDate, nil
		}
	}
	return time.Now(), fmt.Errorf("Period not found")
}
