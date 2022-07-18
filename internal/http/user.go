package http

import (
	"fmt"
	"net/http"
	"strconv"
)

func (rou *Router) GetCourses(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(UserId).(int64)

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
	userId := r.Context().Value(UserId).(int64)
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
	userId := r.Context().Value(UserId).(int64)
	code, mes := rou.BLogic.GetPastWebinars(userId, res)
	if code != 200 {
		http.Error(w, mes, code)
		return
	}
	w.Write([]byte(mes))
}
func (rou *Router) GetTodayWebinars(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserId).(int64)
	code, mes := rou.BLogic.GetTodayWebinars(userId)
	if code != 200 {
		http.Error(w, mes, code)
		return
	}
	w.Write([]byte(mes))
}

func (rou *Router) GetHomework(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserId).(int64)
	var hwI string
	hwI = r.URL.Query().Get("homework_id")
	homeworkId, err := strconv.Atoi(hwI)
	fmt.Println(homeworkId)
	if err != nil {
		http.Error(w, "type \"homework_id\" is not valid", 500)
		return
	}
	code, mes := rou.BLogic.GetHomework(userId, homeworkId)
	if code != 200 {
		http.Error(w, string(mes), code)
		return
	}
	w.Write(mes)
}

func (rou *Router) GetNextCourseHomeworks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserId).(int64)
	var ci string
	ci = r.URL.Query().Get("course_id")
	courseId, er := strconv.Atoi(ci)
	if er != nil {
		http.Error(w, "type \"course_id\" is not valid", 500)
		return
	}
	code, mes := rou.BLogic.GetNextCourseHomeworks(userId, courseId)
	if code != 200 {
		http.Error(w, string(mes), code)
		return
	}
	w.Write(mes)
}
func (rou *Router) GetPastCourseHomeworks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserId).(int64)
	var ci string
	ci = r.URL.Query().Get("course_id")
	courseId, er := strconv.Atoi(ci)
	fmt.Println(courseId)
	if er != nil {
		http.Error(w, "type \"course_id\" is not valid", 500)
		return
	}
	code, mes := rou.BLogic.GetPastCourseHomeworks(userId, courseId)
	if code != 200 {
		http.Error(w, string(mes), code)
		return
	}
	w.Write(mes)
}

func (rou *Router) GetNextHomeworks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserId).(int64)

	code, mes := rou.BLogic.GetNextHomeworks(userId)
	if code != 200 {
		http.Error(w, string(mes), code)
		return
	}
	w.Write(mes)
}

func (rou *Router) GetInfoCourse(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserId).(int64)
	var ci string
	ci = r.URL.Query().Get("course_id")
	courseId, er := strconv.Atoi(ci)
	if er != nil {
		http.Error(w, "type \"course_id\" is not valid", 500)
		return
	}
	code, mes := rou.BLogic.GetInfoCourse(userId, courseId)
	if code != 200 {
		http.Error(w, string(mes), code)
		return
	}
	w.Write(mes)
}
