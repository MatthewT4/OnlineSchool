package blogic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (b *BLogic) GetActivePaymentsPeriod(userId int64) (int, []byte) {
	type course struct {
		CourseId    int       `json:"course_id"`
		PeriodStart time.Time `json:"period_start"`
		PeriodEnd   time.Time `json:"period_end"`
		Price       int       `json:"price"`
		Name        string    `json:"name"`
		PeriodId    int       `json:"period_id"`
	}
	var retStruct struct {
		UserCourses      []course `json:"user_courses,omitempty"`
		AvailableCourses []course `json:"available_courses"`
	}
	var masAddCourseId []int //saving courses from UserCourse to exclude them from AvailableCourse
	if userId != -1 {
		//get id user courses
		courses, err := b.DBUser.GetCourses(context.TODO(), userId)
		if err == nil {
			for _, val := range courses {
				maxPeriodId := 0
				for _, v := range val.BuyPeriod {
					if maxPeriodId < v {
						maxPeriodId = v
					}
				}
				cour, er := b.DBCourse.GetCourse(context.TODO(), val.CourseId)
				if er != nil {
					continue
				}

				for _, v := range cour.PaymentPeriods {
					if v.PeriodId == maxPeriodId+1 {
						var vr course
						vr.Name = cour.NameCourse
						vr.Price = v.Price
						vr.PeriodStart = v.StartDate
						vr.PeriodEnd = v.EndDate
						vr.CourseId = val.CourseId
						masAddCourseId = append(masAddCourseId, val.CourseId)
						retStruct.UserCourses = append(retStruct.UserCourses, vr)
					}
				}
			}
		}
	}
	avalCourse, e := b.DBCourse.GetAvailableCourses(context.TODO(), "ege", masAddCourseId)
	if e == nil {
		for _, val := range avalCourse {
			for _, v := range val.PaymentPeriods {
				if (v.PeriodId == 1 && v.EndDate.After(time.Now())) || (v.StartDate.Before(time.Now()) && v.EndDate.After(time.Now())) {
					var vr course
					vr.Name = val.NameCourse
					vr.PeriodId = v.PeriodId
					vr.PeriodStart = v.StartDate
					vr.PeriodEnd = v.EndDate
					vr.CourseId = val.CourseId
					vr.Price = v.Price
					retStruct.AvailableCourses = append(retStruct.AvailableCourses, vr)
					break
				}
			}
		}
	}
	re, erro := json.Marshal(&retStruct)
	if erro != nil {
		fmt.Println(erro)
		return 500, []byte("Server error")
	}
	return 200, re
}
