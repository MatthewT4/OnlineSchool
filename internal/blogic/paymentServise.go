package blogic

import (
	"OnlineSchool/internal/structs"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func (b *BLogic) GetActivePaymentsPeriod(userId int64) (int, []byte) {
	type course struct {
		CourseId    int       `json:"course_id"`
		PeriodStart time.Time `json:"period_start"`
		PeriodEnd   time.Time `json:"period_end"`
		Price       float64   `json:"price"`
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
						vr.PeriodId = v.PeriodId
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

/*
func (b *BLogic) YearActivePaymentsPeriod(userId int64) (int, []byte) {
	type course struct {
		CourseId    int       `json:"course_id"`
		PeriodStart time.Time `json:"period_start"`
		PeriodEnd   time.Time `json:"period_end"`
		Price       float64   `json:"price"`
		Name        string    `json:"name"`
		PeriodId    int       `json:"period_id"`
	}
	var retStruct struct {
		UserCourses      []course `json:"user_courses,omitempty"`
		AvailableCourses []course `json:"available_courses"`
	}
	if userId != -1 {
		userCourse, err := b.DBUser.GetCourses(context.TODO(), userId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return 401, []byte("Auth error")
			}
			return 500, []byte("Server error")
		}
		for _, val := range userCourse {
			course, er := b.DBCourse.GetCourse(context.TODO(), val.CourseId)
			if er != nil {
				continue
			}
		}
	}


}*/

func (b *BLogic) checkOpportunityToBuyCourse(courseId int, periodIds []int, userId int64) (bool, int /**int == result code*/) {
	if len(periodIds) == 0 || len(periodIds) > 1 {
		return false, 400
	}

	course, e := b.DBCourse.GetCourse(context.TODO(), courseId)
	if e != nil {
		if e == mongo.ErrNoDocuments {
			return false, 400
		} else {
			return false, 500
		}
	}

	if periodIds[0] == -1 { //payment all period
		if userId != -1 {
			userCourse, err := b.DBUser.GetCourses(context.TODO(), userId)
			if err != nil && err != mongo.ErrNoDocuments {
				return false, 500
			}
			searchFlag := false
			var userPer []int
			for _, val := range userCourse {
				if val.Active && val.CourseId == courseId {
					searchFlag = true
					userPer = val.BuyPeriod
					break
				}
			}

			if searchFlag {
				//смотрим есть ли следующий платёжный период
				var MaxUserPer int = -1
				for _, val := range userPer {
					if val > MaxUserPer {
						MaxUserPer = val
					}
				}
				for _, val := range course.PaymentPeriods {
					if val.PeriodId == MaxUserPer+1 {
						return true, 200
					}
				}
				return false, 400
			}
		}

		//userId == -1 or searchFlag == false

		if course.AvailableRegistration {
			return true, 200
		}
		return false, 400
	}

	if userId != -1 {
		userCourse, err := b.DBUser.GetCourses(context.TODO(), userId)
		if err != nil && err != mongo.ErrNoDocuments {
			return false, 500
		}
		searchFlag := false
		var userPer []int
		for _, val := range userCourse {
			if val.Active && val.CourseId == courseId {
				searchFlag = true
				userPer = val.BuyPeriod
				break
			}
		}

		if searchFlag {
			maxUserPeriodId := -1
			for _, val := range userPer {
				if val > maxUserPeriodId {
					maxUserPeriodId = val
				}
			}
			if periodIds[0] != maxUserPeriodId+1 {
				return false, 400
			}
			for _, val := range course.PaymentPeriods {
				if val.PeriodId == periodIds[0] {
					return true, 200
				}
			}
			return false, 400
		}
	}

	//UserId == -1 or searchFlag == false
	if !course.AvailableRegistration {
		return false, 400
	}

	for _, val := range course.PaymentPeriods {
		if val.PeriodId == periodIds[0] {
			if val.PeriodId == 1 {
				return true, 200
			}
			if val.StartDate.Before(time.Now()) && val.EndDate.After(time.Now()) {
				return true, 200
			}
			return false, 400
		}
	}
	return false, 400
}

func (b *BLogic) CreatePayment(buy []structs.PayCourseType, userId int64, promoCodes string) (int, []byte) {
	//check valid request
	fmt.Println("UID:", userId)
	for _, val := range buy {
		res, code := b.checkOpportunityToBuyCourse(val.CourseId, val.Periods, userId)
		if !res {
			if code == 400 {
				return code, []byte("request validation error")
			}
			return code, []byte("server error (validation)")
		}
	}
	if len(buy) == 0 {
		return 400, []byte("request validation error")
	}
	var payment structs.Payment
	if userId != -1 {
		payment.UserId = userId
	}

	for _, val := range buy {
		course, err := b.DBCourse.GetCourse(context.TODO(), val.CourseId)
		if err != nil {
			return 500, []byte("server error")
		}

		var courseBuy structs.PayCourseType
		courseBuy.CourseId = course.CourseId
		courseBuy.TotalPriceWithoutDiscount = 0
		//buy all periods
		if val.Periods[0] == -1 {
			if userId != -1 {
				userCourse, er := b.DBUser.GetCourses(context.TODO(), userId)
				if er != nil && er != mongo.ErrNoDocuments {
					return 500, []byte("server error")
				}
				searchFlag := false
				var userPer []int
				for _, uv := range userCourse {
					if uv.Active && uv.CourseId == val.CourseId {
						searchFlag = true
						userPer = uv.BuyPeriod
					}
				}

				if searchFlag {
					maxUserPer := -1 // максимальный период, который юзер купил ранее
					for _, v := range userPer {
						if v > maxUserPer {
							maxUserPer = v
						}
					}

					for _, v := range course.PaymentPeriods {
						if v.PeriodId > maxUserPer {
							courseBuy.Periods = append(courseBuy.Periods, v.PeriodId)
							courseBuy.TotalPriceWithoutDiscount += v.Price
						}
					}
					if len(courseBuy.Periods) == 0 {
						return 400, []byte("request error")
					}
				}
			}

			if len(courseBuy.Periods) == 0 { // исключаем случаи когда userId != -1 and searchFlag == true

				if !course.AvailableRegistration {
					return 400, []byte("request error")
				}

				for _, v := range course.PaymentPeriods {
					if v.EndDate.After(time.Now()) { //EndDate > time.Now
						courseBuy.Periods = append(courseBuy.Periods, v.PeriodId)
						courseBuy.TotalPriceWithoutDiscount += v.Price
					}
				}
				if len(courseBuy.Periods) == 0 {
					return 400, []byte("request error")
				}
			}

		} else {
			if userId != -1 {
				userCourse, er := b.DBUser.GetCourses(context.TODO(), userId)
				if er != nil && er != mongo.ErrNoDocuments {
					return 500, []byte("server error")
				}
				searchFlag := false
				var userPer []int
				for _, uv := range userCourse {
					if uv.Active && uv.CourseId == val.CourseId {
						searchFlag = true
						userPer = uv.BuyPeriod
					}
				}
				if searchFlag {
					maxUserPer := -1 // максимальный период, который юзер купил ранее
					for _, v := range userPer {
						if v > maxUserPer {
							maxUserPer = v
						}
					}

					for _, uv := range course.PaymentPeriods {
						if uv.PeriodId == maxUserPer+1 {
							courseBuy.Periods = append(courseBuy.Periods, uv.PeriodId)
							courseBuy.TotalPriceWithoutDiscount += uv.Price
							break
						}
					}
				}
			}
			if len(courseBuy.Periods) == 0 {
				if !course.AvailableRegistration {
					return 400, []byte("request error")
				}
				for _, v := range course.PaymentPeriods {
					if v.PeriodId == val.Periods[0] && v.EndDate.After(time.Now()) {
						courseBuy.Periods = append(courseBuy.Periods, v.PeriodId)
						courseBuy.TotalPriceWithoutDiscount += v.Price
					}
				}
			}
		}
		if len(courseBuy.Periods) == 0 {
			return 400, []byte("periods incorrect")
		}
		payment.TotalAmount += courseBuy.TotalPriceWithoutDiscount
		payment.PayCourses = append(payment.PayCourses, courseBuy)
	}
	if len(payment.PayCourses) == 0 {
		return 500, []byte("Server error (courses array == 0)")
	}

	var his structs.History
	his.Status = structs.Registered
	his.ChangeDate = time.Now()
	payment.ChangeHistory = append(payment.ChangeHistory, his)

	//!!!Пока платёжного шлюза нет!!!
	var vrHis structs.History
	vrHis.Status = structs.PreApproved
	vrHis.ChangeDate = time.Now()
	payment.ChangeHistory = append(payment.ChangeHistory, vrHis)
	payment.Status = structs.PreApproved

	payId, e := b.DBPayment.AddPayment(context.TODO(), payment)
	if e != nil {
		fmt.Println(e.Error())
		return 500, []byte("server error")
	}
	var data struct {
		PaymentName string  `json:"payment_name"`
		PaymentId   string  `json:"payment_id"`
		Total       float64 `json:"total"`
		Status      int     `json:"status"`
		Cookie      string  `json:"cookie"`
	}

	if payment.Status == structs.PreApproved || payment.Status == structs.PaymentApproved {
		addRes, errr := b.addUserCourse(userId, payment.PayCourses)
		if !addRes {
			fmt.Println("[Create payment] (add user course):", errr.Error())
			return 500, []byte("Server error  (add user course)")
		}
	}
	data.PaymentName = "Оплата курсов Лицей15"
	data.Total = payment.TotalAmount
	data.PaymentId = payId
	data.Status = payment.Status
	coc := http.Cookie{Name: "PaymentID", Value: payId, Expires: time.Now().Add(time.Hour * 24 * 10), Path: "/"}
	data.Cookie = coc.String()
	js, er := json.Marshal(&data)
	if er != nil {
		return 500, []byte("server error")
	}
	return 200, js
}

func (b *BLogic) LinkingPaymentToUser(userId int64, paymentId string) (int, string) {
	res, err := b.DBPayment.FindPayment(context.TODO(), paymentId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 400, "Payment not found"
		}
		return 500, "Server error (find payment)"
	}
	if res.UserId != 0 {
		if res.UserId == userId {
			return 400, "Payment linking already"
		}
		return 400, "Payment linking to other user already"
	}
	his := structs.History{
		Status:     res.Status,
		ChangeDate: time.Now(),
		Comment:    "Payment linking to the user",
	}

	if res.Status == structs.PreApproved || res.Status == structs.PaymentApproved {
		addRes, errr := b.addUserCourse(userId, res.PayCourses)
		if !addRes {
			fmt.Println("[LinkingPaymentToUser] (add user course):", errr.Error())
			return 500, "Server error (add user course)"
		}
	}

	updCound, e := b.DBPayment.EditOwnerPayment(context.TODO(), paymentId, userId, his)
	if e != nil {
		return 500, "Server Error (edit owner payment)"
	}
	if updCound == 0 {
		return 400, "Bad Request"
	}

	return 200, "ОК"
}
