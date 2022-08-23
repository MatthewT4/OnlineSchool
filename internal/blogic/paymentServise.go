package blogic

import (
	"OnlineSchool/internal/structs"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	CpSecret = "210910163250f488da40c0b4e575ed45"
)

func (b *BLogic) GetActivePaymentsPeriod(userId int64) (int, []byte) {
	type course struct {
		CourseId              int       `json:"course_id"`
		PeriodStart           time.Time `json:"period_start"`
		PeriodEnd             time.Time `json:"period_end"`
		Price                 float64   `json:"price"`
		Name                  string    `json:"name"`
		PeriodId              int       `json:"period_id"`
		AvailableBuyAllPeriod bool      `json:"buy_all_periods"`
		PriceAllPeriod        float64   `json:"price_all_periods"`
		EndAllPeriods         time.Time `json:"end_all_periods"`
		DiscountAllPers       float64   `json:"discount_all_pers"`
		DiscountThePer        float64   `json:"discount_the_per"`
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
				countAwaiPer := 0
				var PriceAllAwailPerWODiscount float64
				PriceAllAwailPerWODiscount = 0
				var vr course
				var flagSetEndAllPer = false
				vr.AvailableBuyAllPeriod = false
				for _, v := range cour.PaymentPeriods {
					if v.PeriodId > maxPeriodId {
						countAwaiPer += 1
						PriceAllAwailPerWODiscount += v.Price
						if flagSetEndAllPer == false || vr.EndAllPeriods.Before(v.EndDate) {
							flagSetEndAllPer = true
							vr.EndAllPeriods = v.EndDate
						}
					}
					if v.PeriodId == maxPeriodId+1 {
						vr.Name = cour.NameCourse
						vr.Price = v.Price
						vr.PeriodId = v.PeriodId
						vr.PeriodStart = v.StartDate
						vr.PeriodEnd = v.EndDate
						vr.CourseId = val.CourseId
						//masAddCourseId = append(masAddCourseId, val.CourseId)
					}
				}

				if vr.Name != "" {
					if countAwaiPer*100/len(cour.PaymentPeriods) >= 70 {
						vr.AvailableBuyAllPeriod = true
						vr.PriceAllPeriod = math.Floor(PriceAllAwailPerWODiscount - 0.15*PriceAllAwailPerWODiscount)
						vr.DiscountAllPers = PriceAllAwailPerWODiscount - vr.PriceAllPeriod
					}
					retStruct.UserCourses = append(retStruct.UserCourses, vr)
				}

				masAddCourseId = append(masAddCourseId, val.CourseId)
			}
		}
	}
	avalCourse, e := b.DBCourse.GetAvailableCourses(context.TODO(), "ege", masAddCourseId)
	if e == nil {
		for _, val := range avalCourse {
			var vr course
			countAwaiPer := 0
			var PriceAllAwailPerWODiscount float64
			PriceAllAwailPerWODiscount = 0
			vr.AvailableBuyAllPeriod = false
			flagSetEndAllPer := false
			for _, v := range val.PaymentPeriods {
				if flagSetEndAllPer == false || vr.EndAllPeriods.Before(v.EndDate) {
					flagSetEndAllPer = true
					vr.EndAllPeriods = v.EndDate
				}
				if v.EndDate.After(time.Now().Add(time.Hour * 24 * 5)) {
					countAwaiPer += 1
					PriceAllAwailPerWODiscount += v.Price
				}
				if (v.PeriodId == 1 && v.EndDate.After(time.Now().Add(time.Hour*24*5))) || (v.StartDate.Before(time.Now().Add(time.Hour*24*5)) && v.EndDate.After(time.Now().Add(time.Hour*24*5))) {
					vr.Name = val.NameCourse
					vr.PeriodId = v.PeriodId
					vr.PeriodStart = v.StartDate
					vr.PeriodEnd = v.EndDate
					vr.CourseId = val.CourseId
					vr.Price = v.Price
					//retStruct.AvailableCourses = append(retStruct.AvailableCourses, vr)
					//break
				}
			}
			if countAwaiPer*100/len(val.PaymentPeriods) >= 70 {
				vr.AvailableBuyAllPeriod = true
				vr.PriceAllPeriod = math.Floor(PriceAllAwailPerWODiscount - 0.15*PriceAllAwailPerWODiscount)
				vr.DiscountAllPers = PriceAllAwailPerWODiscount - vr.PriceAllPeriod
			}
			retStruct.AvailableCourses = append(retStruct.AvailableCourses, vr)
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
			if val.StartDate.Before(time.Now().Add(time.Hour*24*5)) && val.EndDate.After(time.Now().Add(time.Hour*24*5)) {
				return true, 200
			}
			return false, 400
		}
	}
	return false, 400
}

func (b *BLogic) CreatePayment(buy []structs.PayCourseType, userId int64, promoCode string) (int, []byte, http.Cookie) {
	//check valid request
	fmt.Println("UID:", userId)
	for _, val := range buy {
		res, code := b.checkOpportunityToBuyCourse(val.CourseId, val.Periods, userId)
		if !res {
			if code == 400 {
				return code, []byte("request validation error"), http.Cookie{}
			}
			return code, []byte("server error (validation)"), http.Cookie{}
		}
	}
	if len(buy) == 0 {
		return 400, []byte("request validation error"), http.Cookie{}
	}
	var payment structs.Payment
	if userId != -1 {
		payment.UserId = userId
	}
	for _, val := range buy {
		course, err := b.DBCourse.GetCourse(context.TODO(), val.CourseId)
		if err != nil {
			return 500, []byte("server error"), http.Cookie{}
		}

		var courseBuy structs.PayCourseType
		courseBuy.CourseId = course.CourseId
		courseBuy.TotalPriceWithoutDiscount = 0
		courseBuy.TotalPrice = 0
		//buy all periods
		if val.Periods[0] == -1 {
			if userId != -1 {
				userCourse, er := b.DBUser.GetCourses(context.TODO(), userId)
				if er != nil && er != mongo.ErrNoDocuments {
					return 500, []byte("server error"), http.Cookie{}
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
							courseBuy.TotalPrice += v.Price
						}
					}
					if len(courseBuy.Periods) == 0 {
						return 400, []byte("request error"), http.Cookie{}
					}
				}
			}

			if len(courseBuy.Periods) == 0 { // исключаем случаи когда userId != -1 and searchFlag == true

				if !course.AvailableRegistration {
					return 400, []byte("request error"), http.Cookie{}
				}

				for _, v := range course.PaymentPeriods {
					if v.EndDate.After(time.Now().Add(time.Hour * 24 * 5)) { //EndDate > time.Now
						courseBuy.Periods = append(courseBuy.Periods, v.PeriodId)
						courseBuy.TotalPriceWithoutDiscount += v.Price
						courseBuy.TotalPrice += v.Price
					}
				}
				if len(courseBuy.Periods) == 0 {
					return 400, []byte("request error"), http.Cookie{}
				}
			}
			courseBuy.TotalPrice = math.Floor(courseBuy.TotalPriceWithoutDiscount - 0.15*courseBuy.TotalPriceWithoutDiscount)
			courseBuy.TotalPriceWithoutDiscount = math.Floor(courseBuy.TotalPriceWithoutDiscount - 0.15*courseBuy.TotalPriceWithoutDiscount)
		} else {
			if userId != -1 {
				userCourse, er := b.DBUser.GetCourses(context.TODO(), userId)
				if er != nil && er != mongo.ErrNoDocuments {
					return 500, []byte("server error"), http.Cookie{}
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
							courseBuy.TotalPrice += uv.Price
							break
						}
					}
				}
			}
			if len(courseBuy.Periods) == 0 {
				if !course.AvailableRegistration {
					return 400, []byte("request error"), http.Cookie{}
				}
				for _, v := range course.PaymentPeriods {
					if v.PeriodId == val.Periods[0] && v.EndDate.After(time.Now().Add(time.Hour*24*5)) {
						courseBuy.Periods = append(courseBuy.Periods, v.PeriodId)
						courseBuy.TotalPriceWithoutDiscount += v.Price
						courseBuy.TotalPrice += v.Price
					}
				}
			}
		}
		if len(courseBuy.Periods) == 0 {
			return 400, []byte("periods incorrect"), http.Cookie{}
		}
		payment.TotalAmount += courseBuy.TotalPrice
		payment.PayCourses = append(payment.PayCourses, courseBuy)
	}
	if len(payment.PayCourses) == 0 {
		return 500, []byte("Server error (courses array == 0)"), http.Cookie{}
	}

	var his structs.History
	his.Status = structs.Registered
	his.ChangeDate = time.Now()
	payment.ChangeHistory = append(payment.ChangeHistory, his)

	//!!!Пока платёжного шлюза нет!!!
	/*var vrHis structs.History
	vrHis.Status = structs.PreApproved
	vrHis.ChangeDate = time.Now()
	payment.ChangeHistory = append(payment.ChangeHistory, vrHis)
	payment.Status = structs.PreApproved*/
	promoFlag := false
	if userId != -1 && promoCode != "" {
		resultCode, _, promoDisc := b.applyDiscount(promoCode, payment.TotalAmount, userId)
		if resultCode == 0 {
			/*payment.TotalAmount = promoAmount
			payment.DiscountAmount = promoDisc
			*/
			seilOneCourse := math.Ceil(promoDisc/float64(len(payment.PayCourses))*100) / 100
			sumNevostr := 0.00
			for i := 0; i < len(payment.PayCourses); i++ {
				if payment.PayCourses[i].TotalPrice >= seilOneCourse {
					payment.PayCourses[i].TotalPrice -= seilOneCourse
				} else {
					sumNevostr += seilOneCourse - payment.PayCourses[i].TotalPrice
					payment.PayCourses[i].TotalPrice = 0
				}
				if sumNevostr > 0 && payment.PayCourses[i].TotalPrice != 0 {
					if payment.PayCourses[i].TotalPrice >= sumNevostr {
						payment.PayCourses[i].TotalPrice -= sumNevostr
						sumNevostr = 0
					} else {
						sumNevostr -= payment.PayCourses[i].TotalPrice
						payment.PayCourses[i].TotalPrice = 0
					}
				}
			}
			payment.TotalAmount -= seilOneCourse * float64(len(payment.PayCourses))
			payment.DiscountAmount = seilOneCourse * float64(len(payment.PayCourses))

			var vrPromo structs.Discount
			vrPromo.TypeDiscount = 0
			vrPromo.DisAmount = payment.TotalAmount
			payment.Discounts = append(payment.Discounts, vrPromo)
			promoFlag = true
		}
	}

	payId, e := b.DBPayment.AddPayment(context.TODO(), payment)
	if e != nil {
		fmt.Println(e.Error())
		return 500, []byte("server error"), http.Cookie{}
	}

	if promoFlag {
		var vrHisPromo structs.ApplyPromoCode
		vrHisPromo.PromoCode = promoCode
		vrHisPromo.PaymentId = payId
		vrHisPromo.Owner = userId
		vrHisPromo.ApplicationTime = time.Now()

		errAddHisPromo := b.DBAppliedPromoCode.AddHistoryElem(context.TODO(), vrHisPromo)
		if errAddHisPromo != nil {
			fmt.Println("[CreatePayment] (AddHistoryElem) error:", errAddHisPromo.Error())
			return 500, []byte("server error"), http.Cookie{}
		}
	}

	type ReceiptItem struct { //товар в чеке
		Label    string `json:"Label"`
		Price    string `json:"Price"`
		Quantity string `Json:"Quantity"`
		Vat      string `json:"Vat"`
		Amount   string `json:"Amount"`
	}
	type AmountsType struct {
		Electronic string `json:"Electronic"`
	}
	type Receipt struct { //чек
		Items   []ReceiptItem `json:"Items"`
		Amounts AmountsType   `json:"Amounts"`
	}

	var data struct {
		PaymentName string  `json:"payment_name"`
		PaymentId   string  `json:"payment_id"`
		Total       float64 `json:"total"`
		Status      int     `json:"status"`
		Cookie      string  `json:"cookie"`
		UserId      int64   `json:"user_id"`
		ReceiptRet  Receipt `json:"receipt"`
	}
	if payment.UserId != 0 {
		if payment.Status == structs.PreApproved || payment.Status == structs.PaymentApproved {
			addRes, errr := b.addUserCourse(userId, payment.PayCourses)
			if !addRes {
				fmt.Println("[Create payment] (add user course):", errr.Error())
				return 500, []byte("Server error  (add user course)"), http.Cookie{}
			}
		}
	}
	monthConst := [12]string{"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь", "Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"}

	for _, rVal := range payment.PayCourses {
		course, erCo := b.DBCourse.GetCourse(context.TODO(), rVal.CourseId)
		if erCo != nil {
			return 500, []byte("er"), http.Cookie{}
		}
		var vr ReceiptItem
		vr.Price = fmt.Sprint(rVal.TotalPriceWithoutDiscount)
		vr.Amount = fmt.Sprint(rVal.TotalPrice)
		vr.Vat = "20"
		vr.Quantity = "1"

		if len(rVal.Periods) == 1 {
			month := 1
			year := 2022
			for _, per := range course.PaymentPeriods {
				if per.PeriodId == rVal.Periods[0] {
					month = int(per.StartDate.Local().Month())
					year = per.StartDate.Local().Year()
				}
			}

			vr.Label = "Оплата доступа к разделу " + course.NameCourse + " на " + monthConst[month-1] + " " + strconv.Itoa(year) + " года"
		} else {
			startMonth := 1
			startYear := 2022
			endMonth := 1
			endYear := 2022
			minPer := -100
			maxPer := -100
			for _, per := range rVal.Periods {
				if minPer == -100 || per < minPer {
					minPer = per
				}
				if maxPer == -100 || per > maxPer {
					maxPer = per
				}
			}
			for _, per := range course.PaymentPeriods {
				if per.PeriodId == minPer {
					startMonth = int(per.StartDate.Local().Month())
					startYear = per.StartDate.Local().Year()
				}
				if per.PeriodId == maxPer {
					endMonth = int(per.StartDate.Local().Month())
					endYear = per.StartDate.Local().Year()
				}
			}
			vr.Label = "Оплата доступа к разделу " + course.NameCourse + " на " + monthConst[startMonth-1] + " " + strconv.Itoa(startYear) + " года - " + monthConst[endMonth-1] + " " + strconv.Itoa(endYear) + " года"
		}
		data.ReceiptRet.Items = append(data.ReceiptRet.Items, vr)
	}
	data.ReceiptRet.Amounts.Electronic = fmt.Sprint(payment.TotalAmount)

	data.PaymentName = "Оплата курсов Лицей15"
	data.Total = payment.TotalAmount
	data.PaymentId = payId
	data.Status = payment.Status
	data.UserId = payment.UserId

	coc := http.Cookie{Name: "PaymentID", Value: payId, Expires: time.Now().Add(time.Hour * 24 * 10), Path: "/"}
	data.Cookie = coc.String()
	js, er := json.Marshal(&data)
	if er != nil {
		return 500, []byte("server error"), http.Cookie{}
	}
	return 200, js, coc
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

//CloudPayments servises

func (b *BLogic) CheckPayment(CPPayment structs.CloudPaymentReq, data []byte, contextHmac string) []byte {
	//fmt.Println("[CheckPayment] (ComputeHmac256):",ComputeHmac256(data))
	if ComputeHmac256(data) != contextHmac {
		fmt.Println("[CheckPayment] HMAC in request and HMAC calculated in server not request:", contextHmac, "calculated in server:", ComputeHmac256(data))
		return []byte("{\"code\":13}")
	}

	fmt.Println("[CheckPayment] payment in request:", CPPayment)
	payment, err := b.DBPayment.FindPayment(context.TODO(), CPPayment.InvoiceId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("[CheckPayment] (BLogic): Payment", CPPayment.InvoiceId, "not found")
			return []byte("{\"code\":10}")
		}
		fmt.Println("[CheckPayment] (BLogic): Payment", CPPayment.InvoiceId, "server error")
		return []byte("{\"code\":13}")
	}

	if CPPayment.AccountId != payment.UserId {
		fmt.Println("[CheckPayment] (BLogic): Payment", CPPayment.InvoiceId, "User not valid")
		return []byte("{\"code\":11}")
	}

	if payment.TotalAmount != CPPayment.Amount {
		fmt.Println("[CheckPayment] (BLogic): Payment", CPPayment.InvoiceId, "Price not valid")
		return []byte("{\"code\":12}")
	}

	if CPPayment.Currency != "RUB" {
		fmt.Println("[CheckPayment] (BLogic): Payment", CPPayment.InvoiceId, "Currency not valid")
		return []byte("{\"code\":12}")
	}
	fmt.Println("[CheckPayment] (BLogic): Payment", CPPayment.InvoiceId, "OK")
	return []byte("{\"code\":0}")
}

func (b *BLogic) RegisterApprovedPayment(CPPayment structs.CloudPaymentReq, data []byte, contextHmac string) []byte {
	if ComputeHmac256(data) != contextHmac {
		fmt.Println("[RegisterApprovedPayment] HMAC in request and HMAC calculated in server not request:", contextHmac, "calculated in server:", ComputeHmac256(data))
		return []byte("error")
	}

	fmt.Println("[RegisterApprovedPayment] payment in request:", CPPayment)
	payment, err := b.DBPayment.FindPayment(context.TODO(), CPPayment.InvoiceId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("[RegisterApprovedPayment] (BLogic): Payment", CPPayment.InvoiceId, "not found")
			return []byte("payment not found")
		}
		fmt.Println("[RegisterApprovedPayment] (BLogic): Payment", CPPayment.InvoiceId, "server error")
		return []byte("server error")
	}

	if payment.UserId != CPPayment.AccountId {
		fmt.Println("[RegisterApprovedPayment] (BLogic): Payment", CPPayment.InvoiceId, "userId not valid")
		return []byte("userId not valid")
	}

	if CPPayment.Amount != payment.TotalAmount {
		fmt.Println("[RegisterApprovedPayment] (BLogic): Payment", CPPayment.InvoiceId, "amount not valid")
		return []byte("amount not valid")
	}

	var his structs.History
	his.Status = structs.PaymentApproved
	his.ChangeDate = time.Now()
	his.Comment = "TransactionId in CloudPayment: " + strconv.FormatInt(CPPayment.TransactionId, 10)
	/*payment.ChangeHistory = append(payment.ChangeHistory, his)

	payment.Status = structs.PaymentApproved
	b.DBPayment.*/
	coundEdit, er := b.DBPayment.EditStatus(context.TODO(), payment.PaymentId, his, structs.PaymentApproved)
	if er != nil {
		fmt.Println("[RegisterApprovedPayment] (BLogic): Payment", CPPayment.InvoiceId, "UpdatePayment error:", er.Error())
		return []byte("UpdatePayment error")
	}
	if coundEdit == 0 {
		fmt.Println("[RegisterApprovedPayment] (BLogic): Payment", CPPayment.InvoiceId, "cound update is zero")
		return []byte("cound update is zero")
	}
	if payment.UserId != 0 {
		addRes, errr := b.addUserCourse(payment.UserId, payment.PayCourses)
		if !addRes {
			fmt.Println("[RegisterApprovedPayment] (add user course) error:", errr.Error())
			return []byte("add user course error")
		}
	}

	return []byte("{\"code\":0}")
}

func ComputeHmac256(message []byte) string {
	key := []byte(CpSecret)
	h := hmac.New(sha256.New, key)
	h.Write(message)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// applyDiscount return code, amount, discount ||  code: 0=ok 1=incorrect 2=serverError 3=alreadyUsed
func (b *BLogic) applyDiscount(promocode string, amount float64, userId int64) (int, float64, float64) {
	promoCode, err := b.DBPromoCode.GetPromoCode(context.TODO(), promocode)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1, 0, 0
		}
		fmt.Println("[applyDiscount] (GetPromoCode) err:", err.Error())
		return 2, 0, 0
	}
	if promoCode.Owner != 0 {
		if promoCode.Owner != userId {
			return 1, 0, 0
		}
	}
	if promoCode.ValidFrom.After(time.Now()) || promoCode.ValidUntil.Before(time.Now()) {
		return 1, 0, 0
	}
	// Проверка если количество использований не бесконечно!!!!!
	if !promoCode.Infinite {
		return 1, 0, 0
	}
	if !promoCode.MultipleUses {
		PCUsesArray, er := b.DBAppliedPromoCode.GetUserAppliedThePC(context.TODO(), userId, promocode)
		if er != nil {
			if er != mongo.ErrNoDocuments {
				fmt.Println("[applyDiscount] (GetUserAppliedThePC) err:", er.Error())
				return 2, 0, 0
			}
		}

		if len(PCUsesArray) > 0 {
			for _, val := range PCUsesArray {
				payment, e := b.DBPayment.FindPayment(context.TODO(), val.PaymentId)
				if e != nil {
					fmt.Println("[applyDiscount] (FindPayment) err:", e.Error())
					return 2, 0, 0
				}
				if payment.Status == structs.Registered || payment.Status == structs.PaymentRejected {
					continue
				}
				return 3, 0, 0
			}
		}
	}

	if promoCode.TypeDiscount == structs.FixedDiscount {
		if amount <= promoCode.DisAmount {
			if amount < 1 {
				return 0, 0, amount
			}
			return 0, 1, amount - 1
		}
		return 0, amount - promoCode.DisAmount, promoCode.DisAmount
	}
	return 1, 0, 0
}

func (b *BLogic) CheckAmountPromoCodes(userId int64, amount float64, promoCode string) (int, []byte) {
	code, finAmount, discount := b.applyDiscount(promoCode, amount, userId)
	var ret struct {
		Message  string  `json:"message"`
		Amount   float64 `json:"amount"`
		Discount float64 `json:"discount"`
	}
	if code == 0 {
		ret.Discount = discount
		ret.Amount = finAmount
	}
	if code == 1 {
		ret.Message = "Некорректный промокод"
	}
	if code == 2 {
		ret.Message = "Упс, похоже система промокодов сейчас недоступна"
	}
	if code == 3 {
		ret.Message = "Упс, похоже промокод уже был применён ранее"
	}

	jsonRet, er := json.Marshal(&ret)
	if er != nil {
		return 500, []byte("Server error")
	}
	return 200, jsonRet
}
