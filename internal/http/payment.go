package http

import (
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
