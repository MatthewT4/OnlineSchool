package http

import (
	"OnlineSchool/internal/structs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (rou *Router) servProm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		//w.Header().Set("Access-Control-Allow-Origin", "https://lk.lyc15.ru")
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}

// AvailablePaymentPeriods return periods user_course and periods available course
func (rou *Router) AvailablePaymentPeriods(w http.ResponseWriter, r *http.Request) {
	/*//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Origin", "https://lk.lyc15.ru")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")*/

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
	/*w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")*/

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

func (rou *Router) LinkingPaymentToUser(w http.ResponseWriter, r *http.Request) {
	cookie, er := r.Cookie("authToken")
	if er != nil {
		http.Error(w, "authToken not found", 401)
		return
	}
	uId, _, err := rou.BLogic.Authentication(cookie.Value)
	if err != nil {
		http.Error(w, "authToken is not valid", 401)
		return
	}

	cookiePayment, e := r.Cookie("PaymentID")
	if e != nil {
		http.Error(w, "Payment cookie not found", 400)
		return
	}

	code, mes := rou.BLogic.LinkingPaymentToUser(uId, cookiePayment.Value)
	fmt.Println(code, mes)
	if code != 200 {
		http.Error(w, mes, code)
		return
	}
	w.Write([]byte(mes))
}

func (rou *Router) ConnectingCourseGroups(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserId).(int64)
	code, mes := rou.BLogic.CheckConnectingCourseGroups(userId)
	if code != 200 {
		http.Error(w, string(mes), code)
		return
	}
	w.Write(mes)
}

func (rou *Router) InvitationLinkVkGroup(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserId).(int64)
	ci := r.URL.Query().Get("course_id")
	courseId, er := strconv.Atoi(ci)
	if er != nil {
		http.Error(w, "type \"course_id\" is not valid", 500)
		return
	}
	code, mes := rou.BLogic.GetInvitationLinkVkGroup(userId, courseId)
	if code != 200 {
		http.Error(w, string(mes), code)
	}
	w.Write(mes)
}
