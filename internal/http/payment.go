package http

import (
	"OnlineSchool/internal/structs"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// AvailablePaymentPeriods return periods user_course and periods available course
func (rou *Router) AvailablePaymentPeriods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	var userId int64 = -1
	cookie, er := r.Cookie("authToken")
	if er == nil {
		uId, _, err := rou.BLogic.Authentication(cookie.Value)
		if err == nil {
			userId = uId
		}
	}
	code, mes := rou.BLogic.GetActivePaymentsPeriod(userId)
	if code != 200 {
		http.Error(w, string(mes), code)
		return
	}
	w.Write(mes)
}

func (rou *Router) CreatePayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	var userId int64 = -1
	cookie, er := r.Cookie("authToken")
	if er == nil {
		uId, _, err := rou.BLogic.Authentication(cookie.Value)
		if err == nil {
			userId = uId
		}
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "data incorrect", 400)
		return
	}

	type payElem struct {
		CourseId int `json:"course_id"`
		PeriodId int `json:"period_id"`
	}
	var data struct {
		Buy       []payElem `json:"buy"`
		PromoCode string    `json:"promo_code,omitempty"`
	}

	e := json.Unmarshal(body, &data)
	if e != nil {
		http.Error(w, "data incorrect", 400)
		return
	}
	var blData []structs.PayCourseType
	for _, val := range data.Buy {
		var vr structs.PayCourseType
		vr.CourseId = val.CourseId
		vr.Periods = append(vr.Periods, val.PeriodId)
		blData = append(blData, vr)
	}

	code, mes := rou.BLogic.CreatePayment(blData, userId, "")

	if code != 200 {
		http.Error(w, string(mes), code)
		return
	}

	w.Write(mes)
}

/*
func (rou *Router) LinkingPaymentToUser(w http.ResponseWriter, r *http.Request) {
	cookie, er := r.Cookie("authToken")
	if er != nil {
		http.Error(w, "authToken not found", 401)
	}
	uId, _, err := rou.BLogic.Authentication(cookie.Value)
	if err == nil {
		http.Error(w, "authToken is not valid", 401)
	}

	cookiePayment, e := r.Cookie("PaymentID")
	if e != nil {
		http.Error(w, "Payment cookie not found", 400)
	}

}*/
