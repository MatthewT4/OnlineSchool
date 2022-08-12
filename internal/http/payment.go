package http

import (
	"OnlineSchool/internal/structs"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (rou *Router) servProm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", Domain)
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

	code, mes, cookiee := rou.BLogic.CreatePayment(blData, userId, "")

	if code != 200 {
		http.Error(w, string(mes), code)
		return
	}
	http.SetCookie(w, &cookiee)
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

//CloudPayments functions!!

func (rou *Router) CheckPayment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CheckPayment")
	body, errorr := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	//reader := strings.NewReader(string(body))
	if errorr != nil {
		fmt.Println("[HTTP CheckPayment] (readAll body error):", errorr.Error())
		w.Write([]byte("{\"code\":13}"))
		return
	}
	//fmt.Println(body)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("err ParseForm:", err.Error())
		w.Write([]byte("{\"code\":13}"))
		return
	}
	fmt.Println("FORM", r.Form)
	var payment structs.CloudPaymentReq
	transaction := r.Form.Get("TransactionId")
	if transaction == "" {
		fmt.Println("[HTTP CheckPayment]: code 13")
		w.Write([]byte("{\"code\":13}"))
		return
	}
	PayCloudPaymentsId, erro := strconv.ParseInt(transaction, 10, 64)
	if erro != nil {
		fmt.Println("[HTTP CheckPayment]: code 13")
		w.Write([]byte("{\"code\":13}"))
		return
	}
	payment.TransactionId = PayCloudPaymentsId

	amount := r.Form.Get("Amount")
	if amount == "" {
		fmt.Println("[HTTP CheckPayment]: code 12")
		w.Write([]byte("{\"code\":13}"))
		return
	}
	total, er := strconv.ParseFloat(amount, 64)
	if er != nil {
		fmt.Println("[HTTP CheckPayment]: code 12")
		w.Write([]byte("{\"code\":12}"))
		return
	}
	payment.Amount = total

	currency := r.Form.Get("Currency")
	if currency == "" {
		fmt.Println("[HTTP CheckPayment]: code 12")
		w.Write([]byte("{\"code\":13}"))
		return
	}
	payment.Currency = currency

	invoceId := r.Form.Get("InvoiceId")
	if invoceId == "" {
		fmt.Println("[HTTP CheckPayment]: code 10")
		w.Write([]byte("{\"code\":10}"))
		return
	}

	payment.InvoiceId = invoceId

	userId := r.Form.Get("AccountId")
	if userId != "" {
		uIdInt, e := strconv.ParseInt(userId, 10, 64)
		if e == nil {
			payment.AccountId = uIdInt
		}
	}

	//fmt.Println("[Check Payment] payment:", payment)
	//fmt.Println("[HTTP CheckPayment] (X-Content-HMAC):", r.Header.Get("X-Content-HMAC"))
	hmacHead := r.Header.Get("Content-HMAC")
	w.Write(rou.BLogic.CheckPayment(payment, body, hmacHead))
	return
}

func (rou *Router) RegisterApprovedPayment(w http.ResponseWriter, r *http.Request) {

	body, errorr := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	//reader := strings.NewReader(string(body))
	if errorr != nil {
		fmt.Println("[HTTP RegisterApprovedPayment] (readAll body error):", errorr.Error())
	}
	//fmt.Println(body)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("err ParseForm:", err.Error())

	}
	fmt.Println("FORM", r.Form)
	var payment structs.CloudPaymentReq
	transaction := r.Form.Get("TransactionId")
	if transaction == "" {
		fmt.Println("[HTTP RegisterApprovedPayment]: code 13")

	}
	PayCloudPaymentsId, erro := strconv.ParseInt(transaction, 10, 64)
	if erro != nil {
		fmt.Println("[HTTP RegisterApprovedPayment]: code 13")

	}
	payment.TransactionId = PayCloudPaymentsId

	amount := r.Form.Get("Amount")
	if amount == "" {
		fmt.Println("[HTTP RegisterApprovedPayment]: code 12")

	}
	total, er := strconv.ParseFloat(amount, 64)
	if er != nil {
		fmt.Println("[HTTP RegisterApprovedPayment]: code 12")

	}
	payment.Amount = total

	currency := r.Form.Get("Currency")
	if currency == "" {
		fmt.Println("[HTTP RegisterApprovedPayment]: code 12")

	}
	payment.Currency = currency

	invoceId := r.Form.Get("InvoiceId")
	if invoceId == "" {
		fmt.Println("[HTTP RegisterApprovedPayment]: code 10")

	}

	payment.InvoiceId = invoceId

	userId := r.Form.Get("AccountId")
	if userId != "" {
		uIdInt, e := strconv.ParseInt(userId, 10, 64)
		if e == nil {
			payment.AccountId = uIdInt
		}
	}
	//fmt.Println("[Check Payment] payment:", payment)
	//fmt.Println("[HTTP CheckPayment] (X-Content-HMAC):", r.Header.Get("X-Content-HMAC"))
	hmacHead := r.Header.Get("Content-HMAC")
	w.Write(rou.BLogic.RegisterApprovedPayment(payment, body, hmacHead))
	//w.Write([]byte("{\"code\":0}"))
	return
}
