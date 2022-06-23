package http

import (
	"net/http"
	"strconv"
)

func (rou *Router) GetCourses(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(UserId).(int)

	code, body := rou.BLogic.GetUserCourses(user_id)
	if code != 200 {
		http.Error(w, body, code)
		return
	}
	w.Write([]byte(body))
}

func (rou *Router) GetNextWebinars(w http.ResponseWriter, r *http.Request) {
	var courseId string
	courseId = r.URL.Query().Get("course_id")

	res, er := strconv.Atoi(courseId)
	if er != nil {
		http.Error(w, "type \"course_id\" is not valid", 500)
		return
	}
	userId := r.Context().Value(UserId).(int)
	code, mes := rou.BLogic.GetNextWebinars(userId, res)
	if code != 200 {
		http.Error(w, mes, code)
		return
	}
	w.Write([]byte(mes))
}
func (rou *Router) GetPastWebinars(w http.ResponseWriter, r *http.Request) {
	var courseId string
	courseId = r.URL.Query().Get("course_id")

	res, er := strconv.Atoi(courseId)
	if er != nil {
		http.Error(w, "type \"course_id\" is not valid", 500)
		return
	}
	userId := r.Context().Value(UserId).(int)
	code, mes := rou.BLogic.GetPastWebinars(userId, res)
	if code != 200 {
		http.Error(w, mes, code)
		return
	}
	w.Write([]byte(mes))
}
func (rou *Router) GetTodayWebinars(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserId).(int)
	code, mes := rou.BLogic.GetTodayWebinars(userId)
	if code != 200 {
		http.Error(w, mes, code)
		return
	}
	w.Write([]byte(mes))
}
